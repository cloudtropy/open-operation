package router

import (
  "net"
  "time"
  "net/rpc"
  "net/rpc/jsonrpc"
  "strconv"

  log "github.com/cloudtropy/open-operation/utils/logger"
  "github.com/cloudtropy/open-operation/modules/monitor/ctx"
  "github.com/cloudtropy/open-operation/modules/monitor/controller"
)

func rpcRegister(s *rpc.Server) {
  s.Register(new(controller.Template))
  s.Register(new(controller.Agent))
}


func ListenJsonRpc() {
  port := ctx.Cfg().RpcListenPort
  rpcHost := "0.0.0.0:" + strconv.Itoa(port)

  rpcServer := rpc.NewServer()
  // register rpc api
  rpcRegister(rpcServer)

  // /_goRPC_    /debug/rpc
  rpcServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

  tcpAdd, err := net.ResolveTCPAddr("tcp", rpcHost)
  if err != nil {
    log.Fatal("Config error: invalid rpc host.")
  }

  listener, err := net.ListenTCP("tcp", tcpAdd)
  if err != nil {
    log.Fatalf("Listen tcp %s failed: %s\n", tcpAdd, err.Error())
  } else {
    log.Info("Server listen rpc on:", rpcHost)
  }

  for {
    conn, err := listener.Accept()
    if err != nil {
      log.Error("Tcp listener.Accept error", err)
      time.Sleep(time.Second)
      continue
    }
    go rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
  }
}
