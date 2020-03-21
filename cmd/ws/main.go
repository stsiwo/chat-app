package main

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
  "encoding/json"
	"net/http"
)

func main() {
	log.Println("start running ws app...")
	http.ListenAndServe(":8088", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		log.Println("hey")
		if err != nil {
			// handle error
		}
		go func() {
			defer conn.Close()

			var (
				r       = wsutil.NewReader(conn, ws.StateServerSide)
				w       = wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
				decoder = json.NewDecoder(r)
				encoder = json.NewEncoder(w)
			)
			for {
        hdr, err := r.NextFrame()
				if err != nil {
          log.Println("error")
          return
				}
				if hdr.OpCode == ws.OpClose {
					//io.EOF
          return
				}
				var req http.Request
				if err := decoder.Decode(&req); err != nil {
          log.Println("error")
          return
				}
				var resp http.Response
				if err := encoder.Encode(&resp); err != nil {
          log.Println("error")
          return
				}
				if err = w.Flush(); err != nil {
          log.Println("error")
          return
				}
			}
		}()
	}))
}
