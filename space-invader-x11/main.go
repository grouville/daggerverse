package main

import (
	"context"
)

type SpaceInvader struct{}

func (m *SpaceInvader) RunSpaceInvader(ctx context.Context, stringArg string) (*Container, error) {
	_ = stringArg
	project := dag.
		Git("https://github.com/x-hgg-x/space-invaders-go").
		Branch("master").
		Tree()
	return dag.Container().From("golang:latest").
		WithDirectory("/src", project).
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "libgl1-mesa-dev", "libasound2-dev", "libx11-dev", "libxcursor-dev", "libxrandr-dev", "libxinerama-dev", "libxi-dev", "libxxf86vm-dev"}).
		WithEnvVariable("DISPLAY", "host.docker.internal:0").
		WithWorkdir("/src").
		WithExec([]string{"go", "run", "."}).Sync(ctx)
}
