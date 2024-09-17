package main

import (
	"context"
	"dagger/ancient-plotter/internal/dagger"
	"fmt"
	"strings"
)

type AncientPlotter struct{}

const (
	GoVersion = "1.23"

	Registry = "dirkfaust/ancientplotter"
	Version  = "1.1.0"
)

func (*AncientPlotter) Build(_ context.Context, src *dagger.Directory) *dagger.Directory {
	gooses := []string{"linux", "darwin"}
	goarches := []string{"amd64", "arm64"}

	// create empty directory to put build artifacts
	outputs := dag.Directory()

	golang := dag.Container().
		From("golang:"+GoVersion).
		WithDirectory("/src", src).
		WithWorkdir("/src")

	for _, goos := range gooses {
		for _, goarch := range goarches {
			// create directory for each OS and architecture
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)

			// build artifact
			build := golang.
				WithEnvVariable("CGO_ENABLED", "0").
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithExec([]string{"go", "build", "-o", path})

			// add build to outputs
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}

	// return build directory
	return outputs
}

func (*AncientPlotter) Publish(ctx context.Context,
	src *dagger.Directory,
	actor string,
	token *dagger.Secret,
) (string, error) {
	// platforms to build for and push in a multi-platform image
	var platforms = []dagger.Platform{
		"linux/amd64",
		"linux/arm64",
	}

	// container registry for the multi-platform image
	imageRepo := fmt.Sprintf("%s:%s", Registry, Version)

	platformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {
		temp := strings.Split(string(platform), "/")
		if len(temp) != 2 {
			return "", fmt.Errorf("invalid platform: %s", platform)
		}
		platformArch := temp[1]

		ctr := dag.Container().
			From("golang:1.23-alpine").
			WithDirectory("/src", src).
			WithDirectory("/output", dag.Directory()).
			WithEnvVariable("CGO_ENABLED", "0").
			WithEnvVariable("GOOS", "linux").
			WithEnvVariable("GOARCH", platformArch).
			WithWorkdir("/src").
			WithExec([]string{"go", "build", "-o", "/bin/ancientplotter", "-ldflags", "-s -w"})

		outputDir := ctr.Directory("/bin")

		// wrap the output directory in the new empty container marked with the same platform
		binaryCtr := dag.Container(dagger.ContainerOpts{Platform: platform}).
			WithRootfs(outputDir).
			WithEntrypoint([]string{"/ancientplotter"})

		platformVariants = append(platformVariants, binaryCtr)
	}

	// publish to registry
	imageDigest, err := dag.Container().
		WithRegistryAuth("docker.io", actor, token).
		Publish(ctx, imageRepo, dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})

	if err != nil {
		return "", err
	}

	// return build directory
	return imageDigest, nil
}
