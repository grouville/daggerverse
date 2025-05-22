// A generated module for ClaudetteWorkspace functions
//
// This module provides tools for interacting with a workspace,
// including file operations and command execution, similar to the
package main

import (
	"context"
	"dagger/claudette-workspace/internal/dagger"
	"fmt"
	"strings"
)

// ClaudetteWorkspace provides tools for interacting with a workspace directory.
type ClaudetteWorkspace struct {
	// The working directory containing the source code.
	//+defaultPath "/"
	Workdir *dagger.Directory

	// The workspace container tool in which the source code is mounted.
	// This container includes common utilities like 'sh' and 'ripgrep'.
	Container *dagger.Container
}

// BashResult holds the output of a Bash command execution.
type BashResult struct {
	Stdout   string // Standard output of the command.
	Stderr   string // Standard error output of the command.
	ExitCode int    // Exit code of the command.
}

// ClaudetteWorkspace provides tools for interacting with a workspace directory.
// Initializes the toolset with a starting working directory, inside an alpine container
// with ripgrep installed, and mounts the working directory to /src in the container.
func New(workdir *dagger.Directory) *ClaudetteWorkspace {
	// Use a base image with common shell tools + install ripgrep for Grep tool
	ctr := dag.Container().From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "bash", "ripgrep"}). // Add bash and ripgrep
		WithMountedDirectory("/src", workdir).
		WithWorkdir("/src")

	return &ClaudetteWorkspace{
		Workdir:   workdir,
		Container: ctr,
	}
}

// Bash executes a shell command within the workspace container.
// Note: `cd` commands executed here will only affect the single execution context
// of this command, not the persistent state of the workspace Container or Workdir.
func (w *ClaudetteWorkspace) Bash(ctx context.Context, command string) (*BashResult, error) {
	// Basic safety check
	banned := []string{"curl", "wget", "rm -rf /"}
	for _, ban := range banned {
		if strings.Contains(command, ban) {
			// Or return a BashResult with non-zero exit code and error message
			return nil, fmt.Errorf("command '%s' is not allowed for security reasons", command)
		}
	}

	// Execute using 'bash -c' for better shell compatibility
	ctr := w.Container.WithExec([]string{"bash", "-c", command}, dagger.ContainerWithExecOpts{
		UseEntrypoint:            false,
		RedirectStdout:           "/stdout", // Capture stdout to file
		RedirectStderr:           "/stderr", // Capture stderr to file
		InsecureRootCapabilities: false,     // Usually false unless needed
	})

	// Get exit code first. If this errors, it's a Dagger infra error.
	exitCode, err := ctr.ExitCode(ctx)
	if err != nil && exitCode == 0 { // Dagger error, not command failure
		return nil, fmt.Errorf("dagger exec error getting exit code: %w", err)
	}

	// Read stdout/stderr, ignoring errors if the command failed (exitCode != 0)
	stdout, _ := ctr.File("/stdout").Contents(ctx)
	stderr, _ := ctr.File("/stderr").Contents(ctx)

	return &BashResult{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
	}, nil
}

// ReadFile reads the content of a file within the working directory.
// The path should be relative to the root of the Workdir.
func (w *ClaudetteWorkspace) ReadFile(ctx context.Context, path string) (string, error) {
	// Basic path validation
	if strings.HasPrefix(path, "/") || strings.Contains(path, "..") {
		return "", fmt.Errorf("invalid path: '%s'. Path must be relative to the workspace root and cannot contain '..'", path)
	}
	return w.Workdir.File(path).Contents(ctx)
}

// WriteFile writes content to a file within the working directory (overwriting if it exists).
// The path should be relative to the root of the Workdir.
// It returns an *updated* ClaudetteWorkspace instance reflecting the change.
func (w *ClaudetteWorkspace) WriteFile(ctx context.Context, path string, contents string) (*ClaudetteWorkspace, error) {
	// Basic path validation
	if strings.HasPrefix(path, "/") || strings.Contains(path, "..") {
		return nil, fmt.Errorf("invalid path: '%s'. Path must be relative to the workspace root and cannot contain '..'", path)
	}

	// Create the new directory state
	updatedWorkdir := w.Workdir.WithNewFile(path, contents)

	// Return a *new* instance with the updated state
	// Important: Also update the container to mount the *new* directory
	return &ClaudetteWorkspace{
		Workdir: updatedWorkdir,
		Container: dag.Container().From("alpine:latest").
			WithExec([]string{"apk", "add", "--no-cache", "bash", "ripgrep"}).
			WithMountedDirectory("/src", updatedWorkdir).
			WithWorkdir("/src"),
	}, nil
}

// EditFile replaces the first occurrence of oldString with newString in the specified file.
// The path should be relative to the root of the Workdir.
// It returns an *updated* ClaudetteWorkspace instance reflecting the change.
func (w *ClaudetteWorkspace) EditFile(ctx context.Context, path string, oldString string, newString string) (*ClaudetteWorkspace, error) {
	// Basic path validation
	if strings.HasPrefix(path, "/") || strings.Contains(path, "..") {
		return nil, fmt.Errorf("invalid path: '%s'. Path must be relative to the workspace root and cannot contain '..'", path)
	}

	// 1. Read current content
	currentContents, err := w.Workdir.File(path).Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s' for editing: %w", path, err)
	}

	// 2. Perform basic checks and replacement
	if !strings.Contains(currentContents, oldString) {
		return nil, fmt.Errorf("oldString not found in file '%s'", path)
	}
	if strings.Count(currentContents, oldString) > 1 {
		// A more robust solution might involve line/context matching or diff application.
		return nil, fmt.Errorf("oldString found multiple times in file '%s'; edit aborted for safety. Provide more context.", path)
	}

	updatedContents := strings.Replace(currentContents, oldString, newString, 1)

	// 3. Create the new directory state
	updatedWorkdir := w.Workdir.WithNewFile(path, updatedContents)

	// 4. Return a *new* instance with the updated state
	return &ClaudetteWorkspace{
		Workdir: updatedWorkdir,
		Container: dag.Container().From("alpine:latest"). // Recreate container base (to ensure having a fresh state)
									WithExec([]string{"apk", "add", "--no-cache", "bash", "ripgrep"}).
									WithMountedDirectory("/src", updatedWorkdir).
									WithWorkdir("/src"),
	}, nil
}

// Glob finds files matching the specified glob pattern within the working directory.
func (w *ClaudetteWorkspace) Glob(ctx context.Context, pattern string) ([]string, error) {
	// Ensure the pattern doesn't try to escape the workdir
	if strings.HasPrefix(pattern, "/") || strings.Contains(pattern, "..") {
		return nil, fmt.Errorf("invalid glob pattern: '%s'. Pattern cannot start with '/' or contain '..'", pattern)
	}
	return w.Workdir.Glob(ctx, pattern)
}

// Grep searches for a pattern within files in the workspace.
// Uses ripgrep (rg) for efficient searching.
// Returns a string containing matching lines (or filenames if pattern is complex/binary).
func (w *ClaudetteWorkspace) Grep(
	ctx context.Context,
	pattern string,
	// +optional
	path string, // Path relative to workspace root, or empty for root.
	// +optional
	includeGlob string, // Glob pattern to include files (e.g., "*.go")
) (string, error) {

	// Build rg arguments
	args := []string{"rg", "--color=never", "--no-heading", "--with-filename"} // Base args
	if includeGlob != "" {
		args = append(args, "-g", includeGlob)
	}
	args = append(args, pattern) // The search pattern
	if path != "" {
		// Basic path validation
		if strings.HasPrefix(path, "/") || strings.Contains(path, "..") {
			return "", fmt.Errorf("invalid path for grep: '%s'. Path must be relative to the workspace root and cannot contain '..'", path)
		}
		args = append(args, path) // Specific path to search within
	} else {
		args = append(args, ".") // Search current dir (container's /src)
	}

	// Execute rg in the container
	ctr := w.Container.WithExec(args, dagger.ContainerWithExecOpts{
		UseEntrypoint: false,
	})

	// Ripgrep exits 1 if no matches found, 0 on success, >1 on error
	exitCode, err := ctr.ExitCode(ctx)
	if err != nil && exitCode <= 1 { // Ignore dagger exec error if rg exit code is 0 or 1
		err = nil
	} else if err != nil { // Genuine Dagger error
		return "", fmt.Errorf("dagger exec error during grep: %w", err)
	} else if exitCode > 1 { // Ripgrep error
		stderr, _ := ctr.Stderr(ctx)
		return "", fmt.Errorf("grep command failed with exit code %d: %s", exitCode, stderr)
	}

	// If exit code is 1 (no matches), stdout will be empty.
	// If exit code is 0, stdout contains matches.
	stdout, _ := ctr.Stdout(ctx) // Ignore error reading stdout if command potentially failed (exit 1)

	return stdout, nil
}

// Ls lists directory contents using the container's 'ls' command.
// The path should be relative to the root of the Workdir.
func (w *ClaudetteWorkspace) Ls(ctx context.Context, path string) (string, error) {
	// Default to current directory if path is empty
	targetPath := path
	if targetPath == "" {
		targetPath = "."
	}
	// Basic path validation
	if strings.HasPrefix(targetPath, "/") || strings.Contains(targetPath, "..") {
		return "", fmt.Errorf("invalid path for ls: '%s'. Path must be relative to the workspace root and cannot contain '..'", targetPath)
	}

	ctr := w.Container.WithExec([]string{"ls", "-la", targetPath}, dagger.ContainerWithExecOpts{
		UseEntrypoint: false,
	})

	// Capture stdout and stderr, handle potential errors
	stdout, err := ctr.Stdout(ctx)
	if err != nil {
		stderr, _ := ctr.Stderr(ctx) // Attempt to get stderr for context
		return "", fmt.Errorf("ls command failed: %w\nstderr: %s", err, stderr)
	}

	return stdout, nil
}
