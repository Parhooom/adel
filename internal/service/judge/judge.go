package judge

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"adel/internal/service/postgres"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
)

type JudgeService struct {
	dockerClient *DockerClient
	logger       *log.Logger
}

func NewJudge(dockerClient *DockerClient, logger *log.Logger) *JudgeService {
	return &JudgeService{
		dockerClient: dockerClient,
		logger:       logger,
	}
}

func (j *JudgeService) Judge(submission *postgres.Submission, problem *postgres.Problem) error {
	submission.Status = "running"
	var pidLimit int64 = 100

	languageConfig := GetLanguageConfig(submission.Language)
	if languageConfig == nil {
		submission.Status = "internal_error"
		submission.ErrorMessage = "language not supported"
		return fmt.Errorf("language `%s` not supported", submission.Language)
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
				Memory:    int64(problem.MemoryLimit * 1024 * 1024),
				NanoCPUs:  1_000_000_000,
				PidsLimit: &pidLimit,
			},
			AutoRemove:  true,
			CapDrop:     []string{"ALL"},
			SecurityOpt: []string{"no-new-privileges"},
		},
		&network.NetworkingConfig{},
		nil,
		"",
	)
	if err != nil {
		submission.Status = "internal_error"
		submission.ErrorMessage = fmt.Sprintf("failed to create container: %v", err)
		return fmt.Errorf("failed to create container: %w", err)
	}

	codeTar, err := makeTar(submission.Code, languageConfig.FileName)
	if err != nil {
		submission.Status = "internal_error"
		submission.ErrorMessage = fmt.Sprintf("failed to make code tar: %v", err)
		return fmt.Errorf("failed to make code tar: %w", err)
	}

	err = j.dockerClient.client.CopyToContainer(ctx, res.ID, "/app", codeTar, container.CopyToContainerOptions{})
	if err != nil {
		submission.Status = "internal_error"
		submission.ErrorMessage = fmt.Sprintf("failed to copy code to container: %v", err)
		return fmt.Errorf("failed to copy code to container: %w", err)
	}

	err = j.dockerClient.client.ContainerStart(ctx, res.ID, container.StartOptions{})
	if err != nil {
		submission.Status = "internal_error"
		submission.ErrorMessage = fmt.Sprintf("failed to start container: %v", err)
		return fmt.Errorf("failed to start container: %w", err)
	}

	defer j.dockerClient.client.ContainerRemove(ctx, res.ID, container.RemoveOptions{Force: true})

	if len(languageConfig.PreCompileCommand) > 0 {
		_, err := j.runCommand(ctx, res.ID, languageConfig.PreCompileCommand, nil)
		if err != nil {
			submission.Status = "internal_error"
			submission.ErrorMessage = fmt.Sprintf("failed to run pre-compile command: %v", err)
			return fmt.Errorf("failed to run pre-compile command: %w", err)
		}
	}

	if languageConfig.IsCompiled {
		compileCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
		defer cancel()

		compileResult, err := j.runCommand(compileCtx, res.ID, languageConfig.CompileCommand, nil)
		if err != nil {
			submission.Status = "internal_error"
			submission.ErrorMessage = fmt.Sprintf("failed to compile code: %v", err)
			return fmt.Errorf("failed to compile code: %w", err)
		}
		if compileResult.ExitCode != 0 {
			submission.Status = "compile_error"
			submission.ErrorMessage = fmt.Sprintf("compile error: %s", compileResult.Stderr)
			return nil
		}
	}

	submission.ExecutionTimeMs = 0
	submission.MemoryUsageMB = 0
	submission.Status = "accepted"

	for _, testCase := range problem.TestCases {
		if !testCase.IsActive {
			continue
		}

		timeout := time.Duration(problem.TimeLimit) * time.Millisecond

		testCtx, cancel := context.WithTimeout(ctx, timeout+2*time.Second)
		defer cancel()

		result, err := j.runCommand(testCtx, res.ID, languageConfig.RunCommand, []byte(testCase.InputData))

		if err != nil {
			if testCtx.Err() == context.DeadlineExceeded {
				submission.Status = "time_limit_exceeded"
				submission.ErrorMessage = fmt.Sprintf("time limit exceeded: %v", err)
				submission.ExecutionTimeMs = int64(problem.TimeLimit)
				return nil
			}

			submission.Status = "internal_error"
			submission.ErrorMessage = fmt.Sprintf("failed to run test case: %v", err)
			return fmt.Errorf("failed to run test case: %w", err)
		}

		// maximum execution time between all testcases
		executionTimeMs := result.Duration.Milliseconds()
		if executionTimeMs > submission.ExecutionTimeMs {
			submission.ExecutionTimeMs = executionTimeMs
		}

		if result.ExitCode != 0 {
			submission.Status = "runtime_error"
			submission.ErrorMessage = fmt.Sprintf("test case %d failed with exit code %d. Stdout: %q, Stderr: %q", testCase.ID, result.ExitCode, result.Stdout, result.Stderr)
			return nil
		}

		actualOutput := strings.TrimSpace(result.Stdout)
		expectedOutput := strings.TrimSpace(testCase.OutputData)
		if actualOutput != expectedOutput {
			submission.Status = "wrong_answer"
			submission.ErrorMessage = fmt.Sprintf("test case %d failed: expected %q, got %q", testCase.ID, expectedOutput, actualOutput)
			return nil
		}

		memoryUsageMB, err := j.getPeakMemoryUsage(testCtx, res.ID)
		if err != nil {
			submission.Status = "internal_error"
			submission.ErrorMessage = fmt.Sprintf("failed to get memory usage: %v for test case %d", err, testCase.ID)
			return fmt.Errorf("failed to get memory usage: %w", err)
		}
		if memoryUsageMB > submission.MemoryUsageMB {
			submission.MemoryUsageMB = memoryUsageMB
		}
		if memoryUsageMB > int64(problem.MemoryLimit) {
			submission.Status = "memory_limit_exceeded"
			submission.ErrorMessage = fmt.Sprintf("test case %d exceeded memory limit", testCase.ID)
			return nil
		}

		cancel()
	}

	return nil
}

type RunCommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

func (j *JudgeService) runCommand(ctx context.Context, containerID string, cmd []string, input []byte) (*RunCommandResult, error) {
	result := &RunCommandResult{}
	startTime := time.Now()
	defer func() {
		result.Duration = time.Since(startTime)
	}()

	exec, err := j.dockerClient.client.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd:          cmd,
		WorkingDir:   "/app",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create exec: %w", err)
	}

	res, err := j.dockerClient.client.ContainerExecAttach(ctx, exec.ID, container.ExecAttachOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to attach to exec: %w", err)
	}
	defer res.Close()

	deadline, ok := ctx.Deadline()
	if ok {
		res.Conn.SetDeadline(deadline)
	}

	if input != nil {
		_, err = res.Conn.Write(input)
		if err != nil {
			return nil, fmt.Errorf("failed to write input to exec: %w", err)
		}
		err = res.CloseWrite()
		if err != nil {
			return nil, fmt.Errorf("failed to close write: %w", err)
		}
	}

	inspect, err := j.dockerClient.client.ContainerExecInspect(ctx, exec.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect exec: %w", err)
	}

	var stdout, stderr bytes.Buffer

	err = demuxDockerStream(res.Reader, &stdout, &stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to demux stream: %w", err)
	}

	result.Stdout = stdout.String()
	result.Stderr = stderr.String()
	result.ExitCode = inspect.ExitCode

	return result, nil
}

func demuxDockerStream(reader io.Reader, stdout io.Writer, stderr io.Writer) error {
	_, err := stdcopy.StdCopy(stdout, stderr, reader)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to demux stream: %w", err)
	}

	return nil
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

	memoryUsage := int64(statsJSON.MemoryStats.Usage) / (1024 * 1024)
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
