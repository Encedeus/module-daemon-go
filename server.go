package daemon

import (
    "fmt"
    "github.com/Encedeus/module-daemon-go/command"
    "github.com/Encedeus/module-daemon-go/module"
    "github.com/filecoin-project/go-jsonrpc"
    "github.com/stealthrocket/net/wasip1"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
)

type RunFunction func(m *module.Module, backendPort module.Port)

type HandshakeHandler struct {
    RegisteredCommands []*command.Command
    Module             *module.Module
    Run                RunFunction
    BackendPort        module.Port
    FrontendPort       module.Port
}

type HandshakeResponse struct {
    RegisteredCommands []*command.Command
}

func (h *HandshakeHandler) OnHandshake(config module.Configuration) HandshakeResponse {
    go log.Println("Hands have been shook")
    h.Module = &module.Module{
        Port:     config.Port,
        Manifest: config.Manifest,
    }

    defer h.Run(h.Module, h.BackendPort)

    return HandshakeResponse{
        RegisteredCommands: h.RegisteredCommands,
    }
}

func InitModule(run RunFunction) error {
    backendPort, _ := strconv.Atoi(os.Getenv("MODULE_BACKEND_PORT"))
    frontendPort, _ := strconv.Atoi(os.Getenv("MODULE_FRONTEND_PORT"))
    rpcPort, _ := strconv.Atoi(os.Getenv("MODULE_RPC_PORT"))

    rpcServer := jsonrpc.NewServer()

    handshakeHandler := new(HandshakeHandler)
    handshakeHandler.Run = run
    handshakeHandler.BackendPort = module.Port(backendPort)
    handshakeHandler.FrontendPort = module.Port(frontendPort)

    rpcServer.Register("HandshakeHandler", handshakeHandler)

    listener, err := wasip1.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", rpcPort))
    if err != nil {
        log.Fatalf("Failed creating TCP listener: %e", err)
    }

    server := http.Server{
        Handler:      rpcServer,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
        Addr:         fmt.Sprintf(":%v", rpcPort),
    }

    go func() {
        if err = server.Serve(listener); err != nil {
            log.Fatalf("Failed starting RPC server: %e", err)
        }
    }()

    return nil
}
