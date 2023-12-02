package command

import "github.com/Encedeus/module-daemon-go/module"

type Result any
type Parameters []string
type Arguments map[string]any
type Executor func(m *module.Module, args Arguments) (Result, error)

type Command struct {
    Name   string
    Params Parameters
    Exec   Executor
}

type InvokeHandler struct{}

func Invoke(command string) (Result, error) {
    return 0, nil
}
