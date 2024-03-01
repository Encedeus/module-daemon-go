package daemon

import (
	"fmt"
	"github.com/Encedeus/module-daemon-go/module"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/labstack/echo/v4"
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

	mod := new(module.Module)

	//go SetupEchoServer(mod)
	InitRPCServer(module.Port(rpcPort), module.Port(mainPort), mod, run)
}

func SetupEchoServer(mod *module.Module) {
	e := echo.New()

	listener, err := wasip1.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", mod.Port))
	if err != nil {
		log.Fatal(err)
	}
	e.Listener = listener

	mod.Echo = e

	err = module.StartCraterEchoServer(mod)
	if err != nil {
		log.Fatal(err)
	}
}

func InitRPCServer(rpcPort, mainPort module.Port, mod *module.Module, run module.RunFunction) {
	rpcServer := jsonrpc.NewServer()

	handshakeHandler := new(module.HandshakeHandler)
	handshakeHandler.RPCPort = rpcPort
	handshakeHandler.MainPort = mainPort
	handshakeHandler.Module = mod
	handshakeHandler.Run = run

	hostInvokeHandler := new(module.HostInvokeHandler)
	hostInvokeHandler.Module = mod

	cratersHandler := new(module.CraterHandler)
	cratersHandler.RegisteredCraters = &mod.Craters

	rpcServer.Register("HandshakeHandler", handshakeHandler)
	rpcServer.Register("HostInvokeHandler", hostInvokeHandler)
	rpcServer.Register("CratersHandler", cratersHandler)

	rpcListener, err := wasip1.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", rpcPort))
	if err != nil {
		log.Fatalf("Failed creating TCP rpcListener: %v", err)
	}

	server := http.Server{
		Handler:      rpcServer,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Addr:         fmt.Sprintf(":%v", rpcPort),
	}

	if err = server.Serve(rpcListener); err != nil {
		log.Fatalf("Failed starting RPC server: %v", err)
	}
}
