package stdio

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"path/filepath"
	"sync"
)

type Executor interface {
	Execute(ctx context.Context, program string, input io.Reader, output io.Writer) error
	Close() error
}

func NewCommandExecutor(path string) Executor {
	return &commandExecutor{
		path:      path,
		processes: make(map[string]*managedProcess),
	}
}

type managedProcess struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	sync.Mutex
}

type commandExecutor struct {
	path      string
	processes map[string]*managedProcess
	sync.Mutex
}

func (e *commandExecutor) startProcess(program string) (*managedProcess, error) {
	path := filepath.Join(e.path, program)

	cmd := exec.Command(path)
	if err := cmd.Err; err != nil {
		return nil, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		return nil, err
	}

	return &managedProcess{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
	}, nil
}

func (e *commandExecutor) Execute(ctx context.Context, program string, input io.Reader, output io.Writer) error {
	e.Lock()
	defer e.Unlock()

	proc, ok := e.processes[program]
	if !ok {
		p, err := e.startProcess(program)
		if err != nil {
			return err
		}

		proc = p
		e.processes[program] = proc
	}

	proc.Lock()
	defer proc.Unlock()

	if _, err := io.Copy(proc.stdin, input); err != nil {
		return err
	}

	scanner := bufio.NewScanner(proc.stdout)
	if scanner.Scan() {
		line := scanner.Text()
		_, err := output.Write([]byte(line))
		return err
	}

	return scanner.Err()
}

func (e *commandExecutor) Close() error {
	e.Lock()
	defer e.Unlock()

	for _, proc := range e.processes {
		proc.stdin.Close()
		proc.stdout.Close()
		proc.cmd.Process.Kill()
	}

	e.processes = make(map[string]*managedProcess)
	return nil
}

func NewTestableExecutor(handler ExecuteHandler) Executor {
	return &testableExecutor{
		handler: handler,
	}
}

type ExecuteHandler func(ctx context.Context, program string, input io.Reader, output io.Writer) error

type testableExecutor struct {
	handler ExecuteHandler
}

func (e *testableExecutor) Execute(ctx context.Context, program string, input io.Reader, output io.Writer) error {
	return e.handler(ctx, program, input, output)
}

func (e *testableExecutor) Close() error {
	return nil
}
