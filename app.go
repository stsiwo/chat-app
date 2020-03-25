package main

import (
	"log"
	"net/http"
  "github.com/stsiwo/chat-app/infra/net"
  "github.com/stsiwo/chat-app/ui"
)

type App struct {
}

func (App) Run() {
	log.Println("start running ws app...")

	log.Println("start running admin pool & user pool ...")
  var adminPool net.IPool = net.NewPool()
  var userPool net.IPool = net.NewPool()
  go adminPool.Run()
  go userPool.Run()

  log.Println("start setup controllers ...")
  wsController := ui.NewController(
    adminPool,
    userPool,
  )

	log.Println("start listening websocket endpoint ...")
	http.ListenAndServe(":8088", wsController)
}
