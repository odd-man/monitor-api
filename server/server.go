/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package server

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/seeleteam/monitor-api/config"
	"github.com/seeleteam/monitor-api/core"
	"github.com/seeleteam/monitor-api/log"
	"github.com/seeleteam/monitor-api/rpc"
	"github.com/seeleteam/monitor-api/ws"
)

// Start WithErrorGroup
func Start(g *errgroup.Group) {

	monitorServer := core.GetServer(g)
	monitorServer.RunServer()

	// start RPCService, if enableWs = true
	enableWs := config.SeeleConfig.ServerConfig.EnableWebSocket
	if enableWs {
		fmt.Println("will start web socket")
		time.Sleep(5 * time.Second)
		startWsService()
	} else {
		panic("web socket start failed, EnableWebSocket is false")
	}

}

func startWsService() {
	enableRPC := config.SeeleConfig.ServerConfig.EnableRPC
	if !enableRPC {
		fmt.Println("start RPC Service failed, EnableRPC is false")
		return
	}

	rpcURL := config.SeeleConfig.ServerConfig.RPCConfig.URL
	rpcSeeleRPC := rpc.NewSeeleRPC(rpcURL)

	wsURL := config.SeeleConfig.ServerConfig.WebSocketConfig.WsURL

	wsLogger := log.GetLogger("ws", false)
	service, err := ws.New(wsURL, rpcSeeleRPC, wsLogger)
	if err != nil {
		fmt.Println(err)
	}
	//go service.Start()
	service.Start()
}
