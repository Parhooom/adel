package judge

import (
	"adel/internal/service/postgres"
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type JudgeService struct {
	dockerClient *DockerClient
	submission   *postgres.Submission
	problem      *postgres.Problem
	logger       *log.Logger
}

func NewJudgeService(dockerClient *DockerClient, submission *postgres.Submission, problem *postgres.Problem, logger *log.Logger) *JudgeService {
	return &JudgeService{
		dockerClient: dockerClient,
		submission:   submission,
		problem:      problem,
		logger:       logger,
	}
}

func (j *JudgeService) Judge() error {
	j.submission.Status = "running"
	var pidLimit int64 = 100

	languageConfig := GetLanguageConfig(j.submission.Language)
	if languageConfig == nil {
		j.submission.Status = "internal_error"
		j.submission.ErrorMessage = "language not supported"
		return fmt.Errorf("language `%s` not supported", j.submission.Language)
	}

	ctx := context.Background()
	res, err := j.dockerClient.client.ContainerCreate(
		ctx,
		&container.Config{
			Image:           languageConfig.Image,
			WorkingDir:      "/app",
			AttachStdin:     true,
			AttachStdout:    true,
			AttachStderr:    true,
			OpenStdin:       true,
			NetworkDisabled: true,
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:    int64(j.problem.MemoryLimit * 1024 * 1024),
				NanoCPUs:  1_000_000_000,
				PidsLimit: &pidLimit,
			},
			AutoRemove:     true,
			ReadonlyRootfs: true,
			CapDrop:        []string{"ALL"},
			SecurityOpt:    []string{"no-new-privileges"},
			Tmpfs: map[string]string{
				"/app": "rw,exec,size=100m",
			},
		},
		&network.NetworkingConfig{},
		&v1.Platform{},
		"",
	)
	if err != nil {
		j.submission.Status = "internal_error"
		j.submission.ErrorMessage = fmt.Sprintf("failed to create container: %v", err)
		return fmt.Errorf("failed to create container: %w", err)
	}

	codeTar, err := makeTar(j.submission.Code, languageConfig.FileName)
	if err != nil {
		j.submission.Status = "internal_error"
		j.submission.ErrorMessage = fmt.Sprintf("failed to make code tar: %v", err)
		return fmt.Errorf("failed to make code tar: %w", err)
	}

	err = j.dockerClient.client.CopyToContainer(ctx, res.ID, "/app", codeTar, container.CopyToContainerOptions{})
	if err != nil {
		j.submission.Status = "internal_error"
		j.submission.ErrorMessage = fmt.Sprintf("failed to copy code to container: %v", err)
		return fmt.Errorf("failed to copy code to container: %w", err)
	}

	err = j.dockerClient.client.ContainerStart(ctx, res.ID, container.StartOptions{})
	if err != nil {
		j.submission.Status = "internal_error"
		j.submission.ErrorMessage = fmt.Sprintf("failed to start container: %v", err)
		return fmt.Errorf("failed to start container: %w", err)
	}

	if languageConfig.IsCompiled {
		compileCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		compileResult, err := j.runCommand(compileCtx, res.ID, languageConfig.CompileCommand, nil)
		if err != nil {
			j.submission.Status = "internal_error"
			j.submission.ErrorMessage = fmt.Sprintf("failed to compile code: %v", err)
			return fmt.Errorf("failed to compile code: %w", err)
		}
		if compileResult.ExitCode != 0 {
			j.submission.Status = "compile_error"
			j.submission.ErrorMessage = compileResult.Stderr
			return fmt.Errorf("compile error: %s", compileResult.Stderr)
		}
	}

	j.submission.ExecutionTimeMs = 0
	j.submission.MemoryUsageMB = 0
	j.submission.Status = "accepted"

	for _, testCase := range j.problem.TestCases {
		if !testCase.IsActive {
			continue
		}

		timeout := time.Duration(j.problem.TimeLimit) * time.Millisecond
		testCtx, cancel := context.WithTimeout(ctx, timeout+100*time.Millisecond)
		defer cancel()

		result, err := j.runCommand(testCtx, res.ID, languageConfig.RunCommand, []byte(testCase.InputData))
		if err != nil {
			if testCtx.Err() == context.DeadlineExceeded {
				j.submission.Status = "time_limit_exceeded"
				j.submission.ErrorMessage = "execution time limit exceeded"
				return nil
			}

			j.submission.Status = "internal_error"
			j.submission.ErrorMessage = fmt.Sprintf("failed to run test case %d: %v", testCase.ID, err)
			return fmt.Errorf("failed to run test case %d: %w", testCase.ID, err)
		}

		// max execution time for testcases
		executionTimeMs := result.Duration.Milliseconds()
		if executionTimeMs > j.submission.ExecutionTimeMs {
			j.submission.ExecutionTimeMs = executionTimeMs
		}

		if result.ExitCode != 0 {
			j.submission.Status = "runtime_error"
			j.submission.ErrorMessage = fmt.Sprintf("test case %d failed with exit code %d", testCase.ID, result.ExitCode)
			return nil
		}

		actualOutput := strings.TrimSpace(result.Stdout)
		expectedOutput := strings.TrimSpace(testCase.OutputData)
		if actualOutput != expectedOutput {
			j.submission.Status = "wrong_answer"
			j.submission.ErrorMessage = fmt.Sprintf("test case %d failed: expected %q, got %q", testCase.ID, expectedOutput, actualOutput)
			return nil
		}

		memoryUsageMB, err := j.getPeakMemoryUsage(testCtx, res.ID)
		if err != nil {
			j.submission.Status = "internal_error"
			j.submission.ErrorMessage = fmt.Sprintf("failed to get memory usage: %v for test case %d", err, testCase.ID)
			return fmt.Errorf("failed to get memory usage: %w", err)
		}
		if memoryUsageMB > j.submission.MemoryUsageMB {
			j.submission.MemoryUsageMB = memoryUsageMB
		}
		if memoryUsageMB > int64(j.problem.MemoryLimit) {
			j.submission.Status = "memory_limit_exceeded"
			j.submission.ErrorMessage = fmt.Sprintf("test case %d exceeded memory limit", testCase.ID)
			return nil
		}
	}

	return nil
}

type runCommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

func (j *JudgeService) runCommand(ctx context.Context, containerID string, cmd []string, input []byte) (runCommandResult, error) {
	result := runCommandResult{}

	exec, err := j.dockerClient.client.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd:          cmd,
		WorkingDir:   "/app",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return result, fmt.Errorf("failed to create exec: %w", err)
	}

	res, err := j.dockerClient.client.ContainerExecAttach(ctx, exec.ID, container.ExecStartOptions{})
	if err != nil {
		return result, fmt.Errorf("failed to attach to exec: %w", err)
	}
	defer res.Close()

	if input != nil {
		_, err = res.Conn.Write(input)
		if err != nil {
			return result, fmt.Errorf("failed to write input to exec: %w", err)
		}
		err = res.CloseWrite()
		if err != nil {
			return result, fmt.Errorf("failed to close stdin: %w", err)
		}
	}

	var stdout, stderr bytes.Buffer
	startTime := time.Now()
	err = demuxDockerStream(res.Reader, &stdout, &stderr)
	if err != nil {
		return result, fmt.Errorf("failed to read output: %w", err)
	}
	result.Duration = time.Since(startTime)

	inspect, err := j.dockerClient.client.ContainerExecInspect(ctx, exec.ID)
	if err != nil {
		return result, fmt.Errorf("failed to inspect exec: %w", err)
	}
	result.ExitCode = inspect.ExitCode
	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	return result, nil
}

func demuxDockerStream(reader io.Reader, stdout io.Writer, stderr io.Writer) error {
	// _, err := io.Copy(io.MultiWriter(stdout, stderr), reader)
	// if err != nil && err != io.EOF {
	// 	return fmt.Errorf("failed to demux stream: %w", err)
	// }
	_, err := stdcopy.StdCopy(stdout, stderr, reader)
	return err
}

func (j *JudgeService) getPeakMemoryUsage(ctx context.Context, containerID string) (int64, error) {
	stats, err := j.dockerClient.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return 0, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer stats.Body.Close()

	var statsJSON container.StatsResponse
	err = json.NewDecoder(stats.Body).Decode(&statsJSON)
	if err != nil {
		return 0, fmt.Errorf("failed to decode stats: %w", err)
	}

	memoryUsage := int64(statsJSON.MemoryStats.MaxUsage) / (1024 * 1024)
	return memoryUsage, nil
}

func makeTar(code string, fileName string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	tw := tar.NewWriter(&buf)

	header := &tar.Header{
		Name: fileName, // code.c
		Mode: 0744,     // rwxr--r--
		Size: int64(len(code)),
	}

	err := tw.WriteHeader(header)
	if err != nil {
		return nil, fmt.Errorf("failed to write tar header: %w", err)
	}
	_, err = tw.Write([]byte(code))
	if err != nil {
		return nil, fmt.Errorf("failed to write code to tar: %w", err)
	}
	err = tw.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close tar: %w", err)
	}

	return &buf, nil
}
