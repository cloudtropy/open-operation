package ctx

import (
  "fmt"
  "os"
  "os/signal"
  "syscall"

  log "github.com/cloudtropy/open-operation/utils/logger"
)


var signalMap = map[os.Signal]func() {
  os.Interrupt:     handleInt,  //`ctrl-c` / `kill -2 $pid`
  syscall.SIGHUP:   handleHup,  //shutdown shell window
  syscall.SIGTERM:  handleTerm, //`kill $pid` / `pkill agent`
  syscall.SIGQUIT:  handleQuit, //`ctrl-\`
  syscall.SIGUSR1:  handleUsr1, //`kill -USR1 $pid`
}

var ctrl_c = 0
func handleInt() {
  ctrl_c++
  if ctrl_c < 2 {
    fmt.Println("(To exit, press ^C again, or press ctrl-\\)")
    return
  }
  // DelShmData()
  os.Exit(0)
}

func handleHup() {
  return
}

func handleTerm() {
  // DelShmData()
  os.Exit(0)
}

func handleQuit() {
  // DelShmData()
  os.Exit(0)
}

func handleUsr1() {
  os.Exit(0)
}



func Handle(sig os.Signal) (err error) {
  if sig != os.Interrupt {
    ctrl_c = 0
  }
  if _, found := signalMap[sig]; found {
    signalMap[sig]()
    return nil
  } else {
    return fmt.Errorf("No handler available for signal %v", sig)
  }
}


func HandleSignals() {
  sigCount := len(signalMap)
  if sigCount == 0 {
    panic("No signal handler")
  }

  sigs := make([]os.Signal, sigCount)
  i := 0
  for sig, _ := range signalMap {
    sigs[i] = sig
    i++
  }

  signals := make(chan os.Signal)
  signal.Notify(signals, sigs...)

  for sig := range signals {
    log.Printf("Signal captured: %v\n", sig)
    err := Handle(sig)
    if err != nil {
      log.Println("Signal error: ", err)
    }
  }
}

