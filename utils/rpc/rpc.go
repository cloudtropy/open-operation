package rpc

import (
  "net/rpc"
  "net/rpc/jsonrpc"
  "sync"
  "time"

  "github.com/pkg/errors"
)

type SingleConnRpcClient struct {
  sync.Mutex
  RpcServer string
  Timeout   time.Duration
}

func (this *SingleConnRpcClient) Get() (*rpc.Client, error) {
  client, err := jsonrpc.Dial("tcp", this.RpcServer)
  if err != nil {
    return nil, errors.WithMessage(err, "NewRpcClientFail")
  }
  return client, nil
}

func (this *SingleConnRpcClient) Call(method string, args interface{}, reply interface{}) error {

  done := make(chan error, 1)
  var rpcClient *rpc.Client

  go func() {
    var err error
    for i := 0; i < 3; i++ {
      rpcClient, err = this.Get()
      if err != nil {
        continue
      }
      err = rpcClient.Call(method, args, reply)
      closeRpcTcp(rpcClient)
      if err == rpc.ErrShutdown {
        continue
      } else {
        break
      }
    }
    done <- err
  }()

  select {
  case <-time.After(this.Timeout):
    if rpcClient != nil {
      closeRpcTcp(rpcClient)
    }
    return errors.New("RpcTimeout")
  case err := <-done:
    if err != nil {
      return err
    }
  }

  return nil
}


func closeRpcTcp(rc *rpc.Client) {
  if rc == nil {
    return
  }
  rc.Close()
  rc = nil
}
