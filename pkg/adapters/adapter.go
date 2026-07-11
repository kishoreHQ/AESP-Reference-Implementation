package adapters

import (
	"context"
	"io"
)

// RuntimeAdapter is the contract for external CLI agents (K2).
// Vendor-specific behavior lives only in plugin adapters + runtime.yaml.
type RuntimeAdapter interface {
	ID() string
	// Probe returns whether this adapter's CLI is present and a version string.
	Probe(ctx context.Context) (ok bool, version string, detail string)
	// Handshake runs a minimal smoke task (e.g. version + echo).
	Handshake(ctx context.Context) (ok bool, stdout, stderr string, err error)
	// Launch starts or attaches; returns a live session handle.
	Launch(ctx context.Context, opts LaunchOpts) (SessionHandle, error)
}

type LaunchOpts struct {
	WorkDir  string
	Env      map[string]string
	Command  []string // override for generic-pty
	Unsandboxed bool  // true for raw PTY adapters
}

// SessionHandle is a live adapter session stream.
type SessionHandle interface {
	ID() string
	Send(ctx context.Context, message string) error
	// Events yields structured or raw stream events until closed.
	Events() <-chan StreamEvent
	Stop(ctx context.Context) error
	Unsandboxed() bool
}

type StreamEvent struct {
	Kind    string // step|tool_call|model_switch|output|status|raw|error
	Text    string
	Tool    string
	Input   any
	Output  any
	Model   string
	Provider string
	Tokens  int64
	CostUSD float64
	Raw     bool
}

// Catalog lists built-in adapters (plugins can register more).
func Catalog() []RuntimeAdapter {
	return []RuntimeAdapter{
		NewGenericPTY(),
		NewCLIAdapter("claude-code", "claude", []string{"--version"}, "Claude Code CLI"),
		NewCLIAdapter("codex-cli", "codex", []string{"--version"}, "Codex CLI"),
		NewCLIAdapter("gemini-cli", "gemini", []string{"--version"}, "Gemini CLI"),
		NewCLIAdapter("opencode", "opencode", []string{"--version"}, "OpenCode"),
		NewCLIAdapter("aider", "aider", []string{"--version"}, "Aider"),
	}
}

// Discard is a helper writer.
var Discard = io.Discard
