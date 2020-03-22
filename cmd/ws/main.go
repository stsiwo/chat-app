package main

import (
	"github.com/gobwas/ws"
	"log"
	"net/http"
  "github.com/stsiwo/chat-app/infra/net"
  "github.com/stsiwo/chat-app/domain/user"
  "strconv"
)

func main() {
	log.Println("start running ws app...")

	log.Println("start running admin pool & user pool ...")
  adminPool := net.NewPool()
  userPool := net.NewPool()

  go adminPool.Run()
  go userPool.Run()

	log.Println("start listening websocket endpoint ...")
	http.ListenAndServe(":8088", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
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

    var newClient *net.Client
    if userRole == user.Admin {
      newUser := user.NewAdminUser(conn.RemoteAddr().String())
      newClient := net.NewClient(conn, newUser, adminPool, userPool)
      adminPool.Register(newClient)
    } else if userRole == user.Member {
      newUser := user.NewMemberUser(conn.RemoteAddr().String())
      newClient := net.NewClient(conn, newUser, adminPool, userPool)
      userPool.Register(newClient)
    } else {
      newUser := user.NewGuestUser(conn.RemoteAddr().String())
      newClient := net.NewClient(conn, newUser, adminPool, userPool)
      userPool.Register(newClient)
    }

    go newClient.Read()
    go newClient.Write()
	}))
}
