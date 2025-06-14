package judge

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/image"
)

func (docker *DockerClient) ImageList() ([]string, error) {
	res, err := docker.client.ImageList(context.Background(), image.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	var images []string
	for _, img := range res {
		images = append(images, img.RepoTags...)
	}

	return images, nil
}

func (docker *DockerClient) DownloadAllImages() error {
	languageConfigs := GetLanguageConfigs()

	for _, config := range languageConfigs {
		err := docker.ImagePull(config.Image)
		if err != nil {
			return fmt.Errorf("failed to pull image: %w", err)
		}
	}

	docker.logger.Println("all images downloaded")
	return nil
}

func (docker *DockerClient) ImagePull(imageTag string) error {
	res, err := docker.client.ImagePull(context.Background(), imageTag, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer res.Close()

	_, err = io.Copy(io.Discard, res) // could use any io.Writer, for now i'll use io.Discard to avoid noisy logs
	if err != nil {
		return fmt.Errorf("failed to copy image: %w", err)
	}

	docker.logger.Printf("successfully pulled docker image: %s", imageTag)
	return nil
}
