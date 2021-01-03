package notifier

import (
	"context"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// Notifier is a notification interface
type Notifier interface {
	Notify(ctx context.Context, body string) (exit int, err error)
	Exec(ctx context.Context, param ParamExec) error
}

type ParamExec struct {
	Stdout         string
	Stderr         string
	CombinedOutput string
	Args           cli.Args
	Cmd            *exec.Cmd
}
