package command

import (
    "github.com/Encedeus/module-daemon-go/module"
)

type Result any
type Parameters []string
type Arguments map[string]any
type Executor func(m *module.Module, args Arguments) (Result, error)

type Command struct {
    Name   string
    Params Parameters
    Exec   Executor
}

type InvokeFunc func(command string, args Arguments) (Result, error)

type HostInvokeHandler struct {
    Module *module.Module
}

func (h *HostInvokeHandler) HostInvoke(command string, args Arguments) (Result, error) {
    for _, cmd := range h.Module.Commands {
        if cmd.Name == command {
            result, err := cmd.Exec(h.Module, args)
            if err != nil {
                return nil, err
            }

            return result, nil
        }
    }

    return nil, nil
}
