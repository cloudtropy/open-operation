package router

import (
	"cloudtropy.com/alert/controller"
	"cloudtropy.com/alert/g"
	"log"
	"net/http"
	"strconv"
	"time"
)

func RunServer() {
	// ctl.Init()

	http.HandleFunc("/alert", controller.PostAlert)

	http.HandleFunc("/email", controller.PostEmailWarning)
	http.HandleFunc("/wechat", controller.PostWechatWarning)

	// ws
	// g.HandleWithMid("/wsapi/msg", ctl.HandleWsMsg)
	// g.HandleWithMid("/report", ctl.HandleReport)

	var listenPort = g.Config().Self.ListenPort
	s := &http.Server{
		Addr:           ":" + strconv.Itoa(listenPort),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, //if not set, use (DefaultMaxHeaderBytes = 1 << 20) // 1 MB
		//ErrorLog *log.Logger
	}

	log.Println("Server listen on:", listenPort)
	log.Fatal(s.ListenAndServe())
}
