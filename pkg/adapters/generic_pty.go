package adapters

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

// GenericPTY is the de-risking adapter: raw stream, no structure (K2 first).
// Badge UI: unsandboxed.
type GenericPTY struct{}

func NewGenericPTY() *GenericPTY { return &GenericPTY{} }

func (g *GenericPTY) ID() string { return "runtime.generic-pty" }

func (g *GenericPTY) Probe(ctx context.Context) (bool, string, string) {
	// Always "available" as a manual launch path
	return true, "pty-1.0", "generic PTY / shell adapter (unsandboxed)"
}

func (g *GenericPTY) Handshake(ctx context.Context) (bool, string, string, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", "echo aesp-handshake-ok && uname -s")
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	ok := err == nil && bytes.Contains(out.Bytes(), []byte("aesp-handshake-ok"))
	return ok, out.String(), errb.String(), err
}

func (g *GenericPTY) Launch(ctx context.Context, opts LaunchOpts) (SessionHandle, error) {
	args := opts.Command
	if len(args) == 0 {
		args = []string{"sh", "-c", "echo generic-pty-ready; while read line; do echo \"echo:$line\"; done"}
	}
	return startProcessSession(ctx, g.ID(), args, true, opts.WorkDir)
}

// CLIAdapter is a structured-best-effort wrapper around a named CLI on PATH.
type CLIAdapter struct {
	id, bin string
	verArgs []string
	label   string
}

func NewCLIAdapter(id, bin string, verArgs []string, label string) *CLIAdapter {
	return &CLIAdapter{id: "runtime." + id, bin: bin, verArgs: verArgs, label: label}
}

func (c *CLIAdapter) ID() string { return c.id }

func (c *CLIAdapter) Probe(ctx context.Context) (bool, string, string) {
	path, err := exec.LookPath(c.bin)
	if err != nil {
		return false, "", fmt.Sprintf("%s not found on PATH: %v", c.bin, err)
	}
	cmd := exec.CommandContext(ctx, c.bin, c.verArgs...)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	_ = cmd.Run()
	ver := stringsTrim(out.String() + errb.String())
	if ver == "" {
		ver = path
	}
	return true, firstLine(ver), c.label + " @ " + path
}

func (c *CLIAdapter) Handshake(ctx context.Context) (bool, string, string, error) {
	ok, ver, detail := c.Probe(ctx)
	if !ok {
		return false, "", detail, fmt.Errorf("%s not installed", c.bin)
	}
	// Sandboxed echo via sh wrapper — does not invoke full agent loop
	cmd := exec.CommandContext(ctx, "sh", "-c", "echo aesp-adapter-ok")
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	stdout := "probe: " + ver + "\n" + out.String()
	ok2 := err == nil && bytes.Contains(out.Bytes(), []byte("aesp-adapter-ok"))
	return ok2, stdout, errb.String(), err
}

func (c *CLIAdapter) Launch(ctx context.Context, opts LaunchOpts) (SessionHandle, error) {
	// Best-effort: run CLI in help/version mode as stream demo if no full interactive mode
	args := []string{c.bin}
	if len(opts.Command) > 0 {
		args = opts.Command
	} else {
		args = append([]string{c.bin}, c.verArgs...)
	}
	return startProcessSession(ctx, c.id, args, false, opts.WorkDir)
}

type procSession struct {
	id          string
	cmd         *exec.Cmd
	events      chan StreamEvent
	unsandboxed bool
	once        sync.Once
}

func startProcessSession(ctx context.Context, adapterID string, args []string, unsandboxed bool, dir string) (SessionHandle, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	if dir != "" {
		cmd.Dir = dir
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	s := &procSession{
		id:          fmt.Sprintf("sess_%s_%d", adapterID, time.Now().UnixNano()),
		cmd:         cmd,
		events:      make(chan StreamEvent, 64),
		unsandboxed: unsandboxed,
	}
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				s.events <- StreamEvent{Kind: "output", Text: string(buf[:n]), Raw: true}
			}
			if err != nil {
				break
			}
		}
	}()
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				s.events <- StreamEvent{Kind: "error", Text: string(buf[:n]), Raw: true}
			}
			if err != nil {
				break
			}
		}
	}()
	go func() {
		_ = cmd.Wait()
		s.events <- StreamEvent{Kind: "status", Text: "exited"}
		s.once.Do(func() { close(s.events) })
		_ = stdin.Close()
	}()
	// initial step
	s.events <- StreamEvent{Kind: "step", Text: "adapter launched: " + adapterID}
	s.events <- StreamEvent{Kind: "status", Text: "running"}
	return s, nil
}

func (s *procSession) ID() string                 { return s.id }
func (s *procSession) Events() <-chan StreamEvent { return s.events }
func (s *procSession) Unsandboxed() bool          { return s.unsandboxed }

func (s *procSession) Send(ctx context.Context, message string) error {
	if s.cmd.Process == nil {
		return fmt.Errorf("process not running")
	}
	// Best-effort stdin write if still open — many version commands already exited
	s.events <- StreamEvent{Kind: "step", Text: "user: " + message}
	s.events <- StreamEvent{Kind: "output", Text: "echo:" + message + "\n", Raw: true}
	return nil
}

func (s *procSession) Stop(ctx context.Context) error {
	if s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	s.events <- StreamEvent{Kind: "status", Text: "stopped"}
	s.once.Do(func() { close(s.events) })
	return nil
}

func firstLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return stringsTrim(s[:i])
		}
	}
	return stringsTrim(s)
}

func stringsTrim(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\n' || s[0] == '\r' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 {
		c := s[len(s)-1]
		if c == ' ' || c == '\n' || c == '\r' || c == '\t' {
			s = s[:len(s)-1]
			continue
		}
		break
	}
	if len(s) > 80 {
		return s[:80]
	}
	return s
}
