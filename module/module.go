package module

import (
    "context"
    "fmt"
    "github.com/Encedeus/module-daemon-go/command"
    "github.com/filecoin-project/go-jsonrpc"
)

type Manifest struct {
    Name             string   `hcl:"name"`
    Authors          []string `hcl:"authors"`
    Version          string   `hcl:"version"`
    FrontendMainFile string   `hcl:"frontend_main"`
    // BackendMainFile  string   `hcl:"backend_main"`
}

type Module struct {
    Port     Port
    Manifest Manifest
    HostPort Port
}

func (m *Module) Invoke(cmd string, args command.Arguments) (command.Result, error) {
    var client struct {
        ModuleInvoke command.InvokeFunc
    }

    closer, err := jsonrpc.NewClient(context.Background(), fmt.Sprintf("http://localhost:%v", m.HostPort), "ModuleInvokeHandler", &client, nil)
    if err != nil {
        return nil, err
    }
    defer closer()

    result, err := client.ModuleInvoke(cmd, args)
    if err != nil {
        return nil, err
    }

    return result, nil
}

type Port uint16

type Configuration struct {
    Port     Port
    HostPort Port
    Manifest Manifest
}
