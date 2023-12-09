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

type RunFunction func(m *module.Module)

type HandshakeHandler struct {
    RegisteredCommands []*command.Command
    Module             *module.Module
    Run                RunFunction
    RPCPort            module.Port
    MainPort           module.Port
}

type HandshakeResponse struct {
    RegisteredCommands []*command.Command
}

func (h *HandshakeHandler) OnHandshake(config module.Configuration) HandshakeResponse {
    h.Module.Port = config.Port
    h.Module.Manifest = config.Manifest
    h.Module.HostPort = config.HostPort
    h.Module.HandshakeHandler = h

    defer h.Run(h.Module)

    return HandshakeResponse{
        RegisteredCommands: h.RegisteredCommands,
    }
}

func InitModule(run RunFunction) {
    rpcPort, _ := strconv.Atoi(os.Getenv("MODULE_RPC_PORT"))
    mainPort, _ := strconv.Atoi(os.Getenv("MODULE_MAIN_PORT"))

    rpcServer := jsonrpc.NewServer()

    mod := new(module.Module)

    handshakeHandler := new(HandshakeHandler)
    handshakeHandler.RPCPort = module.Port(rpcPort)
    handshakeHandler.MainPort = module.Port(mainPort)
    handshakeHandler.Module = mod

    hostInvokeHandler := new(command.HostInvokeHandler)
    hostInvokeHandler.Module = mod

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

    if err = server.Serve(listener); err != nil {
        log.Fatalf("Failed starting RPC server: %e", err)
    }
}
