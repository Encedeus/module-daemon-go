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
    log.Println("Hands have been shook")
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
    log.Println("Hands have been shook 1")
    backendPort, _ := strconv.Atoi(os.Getenv("MODULE_BACKEND_PORT"))
    frontendPort, _ := strconv.Atoi(os.Getenv("MODULE_FRONTEND_PORT"))
    // rpcPort, _ := strconv.Atoi(os.Getenv("MODULE_RPC_PORT"))

    log.Println("Hands have been shook 2")
    rpcServer := jsonrpc.NewServer()

    handshakeHandler := new(HandshakeHandler)
    handshakeHandler.Run = run
    handshakeHandler.BackendPort = module.Port(backendPort)
    handshakeHandler.FrontendPort = module.Port(frontendPort)

    log.Println("Hands have been shook 3")
    rpcServer.Register("HandshakeHandler", handshakeHandler)

    listener, err := wasip1.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", 8092))
    if err != nil {
        log.Fatalf("Failed creating TCP listener: %e", err)
        return err
    }
    log.Println("Hands have been shook 4")

    server := http.Server{
        Handler:      rpcServer,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
        Addr:         fmt.Sprintf(":%v", 8092),
    }

    if err = server.Serve(listener); err != nil {
        log.Fatalf("Failed starting RPC server: %e", err)
        return err
    }
    log.Println("Hands have been shook 5")

    return nil
}
