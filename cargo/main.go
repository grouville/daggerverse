package main

import (
	"context"
	"fmt"
)

const (
	DEFAULT_RUST = "1.73"
	PROJ_MOUNT   = "/src"
)

type Cargo struct {
	Ctr  *Container
	Proj *Directory
}

// Build the project
func (c *Cargo) Build(args []string) *Directory {
	command := append([]string{"cargo", "build"}, args...)
	return c.prepare().
		WithExec(command).
		Directory(PROJ_MOUNT)
}

// Format the project
func (c *Cargo) Fmt(ctx context.Context, args []string) (string, error) {
	command := append([]string{"cargo", "fmt"}, args...)
	return c.prepare().WithExec(command).Stdout(ctx)
}

// Test the project
func (c *Cargo) Test(ctx context.Context, args []string) (string, error) {
	command := append([]string{"cargo", "test"}, args...)
	return c.prepare().WithExec(command).Stdout(ctx)
}

// Lint the project
func (c *Cargo) Clippy(ctx context.Context) (string, error) {
	return c.prepare().
		WithExec([]string{"cargo", "clippy"}).
		Stdout(ctx)
}

// Sets up the Container with a rust image
func (c *Cargo) Base(version string) *Cargo {
	image := fmt.Sprintf("rust:%s", version)
	c.Ctr = dag.Container().From(image)
	return c
}

// Install a package from a git repository
func (c *Cargo) InstallFromGit(url string, branch string, bin string, pkg string) *Cargo {
	command := []string{"cargo", "install", "--git", url, "--branch", branch, "--bin", bin, pkg}
	c.Ctr = c.prepare().
		WithMountedCache("/usr/local/cargo/registry", dag.CacheVolume("cargoregistry")).
		WithExec(command)
	return c
}

// Accessor for the Container
func (c *Cargo) Container() *Container {
	return c.Ctr
}

// Accessor for the Project
func (c *Cargo) Project() *Directory {
	return c.Ctr.Directory(PROJ_MOUNT)
}

// Specify the Project to use in the module
func (c *Cargo) WithProject(dir *Directory) *Cargo {
	c.Proj = dir
	return c
}

// Bring your own container
func (c *Cargo) WithContainer(ctr *Container) *Cargo {
	c.Ctr = ctr
	return c
}

// Private func to check readiness and prepare the container for build/test/lint
func (c *Cargo) prepare() *Container {
	if c.Proj == nil {
		c.Proj = dag.Directory() // Unsure about this. Maybe want to error
	}

	if c.Ctr == nil {
		cd := c.Base(DEFAULT_RUST)
		c.Ctr = cd.Ctr
	}

	c.Ctr = c.Ctr.
		WithDirectory(PROJ_MOUNT, c.Proj).
		WithWorkdir(PROJ_MOUNT)
	return c.Ctr
}
