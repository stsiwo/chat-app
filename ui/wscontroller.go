package ui

import (
	"github.com/stsiwo/chat-app/domain/user"
	"github.com/stsiwo/chat-app/infra/net"
	"github.com/stsiwo/chat-app/infra/wsutil"
	"github.com/stsiwo/chat-app/infra/ws"
	"net/http"
	"strconv"
  "log"
)

type WsController struct {

	adminPool net.IPool

	userPool net.IPool

  wsutilHandler wsutil.IWsutilHandler

  wsHandler ws.IWsHandler
}

func NewController(adminPool net.IPool, userPool net.IPool) *WsController {
  return &WsController{
    adminPool: adminPool,
    userPool: userPool,
    wsutilHandler: &wsutil.Handler{},
    wsHandler: &ws.Handler{},
  }
}

func (wc *WsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	conn, _, _, err := wc.wsHandler.UpgradeHTTP(r, w)
	if err != nil {
		log.Fatalf("error during upgrading to ws protocol: %v \n", err)
	}

	// read user info from 'conn'
	// use query string to identify type of client
	roleQueryParam := r.URL.Query().Get("role")
	role, err := strconv.Atoi(roleQueryParam)
	if err != nil {
		log.Fatalf("role type query string converting error: %v", err)
	}
	userRole := user.Role(role)
	id := r.URL.Query().Get("id")
  if id == "" {
    log.Fatal("id is nil; id is mandatory query string")
  }

	var newClient net.IClient
	if userRole == user.Admin {
    log.Printf("requested client is admin: %v\n", userRole)
		newUser := user.NewAdminUser(id, conn.RemoteAddr().String())
		newClient = net.NewClient(id, conn, newUser, wc.adminPool, wc.userPool, wc.wsutilHandler)
		wc.adminPool.Register(newClient)
	} else if userRole == user.Member {
    log.Printf("requested client is member: %v\n", userRole)
		newUser := user.NewMemberUser(id, conn.RemoteAddr().String())
		newClient = net.NewClient(id, conn, newUser, wc.adminPool, wc.userPool, wc.wsutilHandler)
		wc.userPool.Register(newClient)
	} else {
    log.Printf("requested client is guest: %v\n", userRole)
		newUser := user.NewGuestUser(id, conn.RemoteAddr().String())
		newClient = net.NewClient(id, conn, newUser, wc.adminPool, wc.userPool, wc.wsutilHandler)
		wc.userPool.Register(newClient)
	}

  log.Printf("start running Read & Write operation for this client: %v \n", newClient)
	go newClient.Read()
	go newClient.Write()
}

