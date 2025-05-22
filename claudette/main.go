package main

import (
	"dagger/claudette/internal/dagger"
)

// Claudette orchestrates interaction with an LLM and a tools workspace.
type Claudette struct {
	// The directory the agent operates on.
	// This directory state is passed to the tools module for each interaction turn.
	WorkspaceDir *dagger.Directory
	APIKey       *dagger.Secret
	Model        string
}

// New initializes the Claude agent.
func New(
	// The host directory to mount as the agent's workspace.
	workdir *dagger.Directory,

	// Anthropic API Key secret.
	apiKey *dagger.Secret,

	// +optional
	// The Claude model to use (e.g., claude-3-5-sonnet-20240620).
	model string,
) *Claudette {
	return &Claudette{
		WorkspaceDir: workdir,
		APIKey:       apiKey,
		Model:        model,
	}
}

// // Chat sends a prompt to the agent and returns the final response after any tool executions.
// func (a *Claudette) Chat(ctx context.Context, prompt string) (string, error) {

// 	// --- Initialize Tool Workspace ---
// 	// Create an instance of the tools module, passing the *current* state of the workspace directory.
// 	// This ensures tools operate on the latest version of the files for this specific chat turn.
// 	tools := dag.ClaudetteWorkspace(a.WorkspaceDir)

// 	systemPrompt := `You are Claudette, an AI assistant powered by Dagger, designed to help with software development tasks in a given workspace.
// You can interact with the workspace using the provided tools module named 'tools'.

// Available tools in the 'tools' module:
// - tools.Bash(command string) (*BashResult): Executes a shell command. Returns stdout, stderr, and exit code. Use cautiously.
// - tools.ReadFile(path string) (string): Reads the content of a file relative to the workspace root.
// - tools.WriteFile(path string, contents string) (*ClaudetteWorkspace): Writes content to a file (overwrites), relative to the workspace root. Returns an updated workspace tool instance.
// - tools.EditFile(path string, oldString string, newString string) (*ClaudetteWorkspace): Replaces the first occurrence of oldString with newString. Returns an updated workspace tool instance.
// - tools.Glob(pattern string) ([]string): Finds files matching a glob pattern relative to the workspace root.
// - tools.Grep(pattern string, path? string, includeGlob? string) (string): Searches file contents using regex (ripgrep). Returns matching lines.
// - tools.Ls(path string) (string): Lists directory contents relative to the workspace root.

// IMPORTANT NOTES:
// - All file paths for tools (ReadFile, WriteFile, EditFile, Ls, Grep path) MUST be relative to the workspace root (e.g., "src/main.go", not "/src/main.go"). Do not use absolute paths or "..".
// - When using WriteFile or EditFile, the tool modifies the workspace. The LLM runtime handles using the *updated* workspace state for subsequent reasoning or tool calls *within this single Chat execution*. Plan your actions accordingly. If you need to ensure a file is read *after* a modification in a completely separate interaction, perform the read in a new Chat call.
// - Keep responses concise and focused on the user's request.
// - If a command execution fails (non-zero exit code from Bash), analyze the stderr output provided in the BashResult.
// `

// 	env := dag.Env().
// 		WithClaudetteWorkspaceInput("claudette-workspace", tools, "tools").
// 		WithDirectoryInput("workdir", a.WorkspaceDir, "current workdir directory")

// 	// --- Initialize LLM Runner ---
// 	response, err := dag.LLM().
// 		WithPrompt(systemPrompt).
// 		WithEnv(env).LastReply(ctx)

// 	return response, err
// }


