package daemon

import (
    "fmt"
    "github.com/Encedeus/module-daemon-go/module"
    "github.com/filecoin-project/go-jsonrpc"
    "github.com/stealthrocket/net/wasip1"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
)

func InitModule(run module.RunFunction) {
    rpcPort, _ := strconv.Atoi(os.Getenv("MODULE_RPC_PORT"))
    mainPort, _ := strconv.Atoi(os.Getenv("MODULE_MAIN_PORT"))
    fmt.Printf("Perceived RPC port: %v\n", rpcPort)

    rpcServer := jsonrpc.NewServer()

    mod := new(module.Module)

    handshakeHandler := new(module.HandshakeHandler)
    handshakeHandler.RPCPort = module.Port(rpcPort)
    handshakeHandler.MainPort = module.Port(mainPort)
    handshakeHandler.Module = mod
    handshakeHandler.Run = run

    hostInvokeHandler := new(module.HostInvokeHandler)
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
