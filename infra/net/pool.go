package net

import (
	"log"
)

type Pool struct {
	pool map[string]*Client

  /**
   * broadcast message to all clients in this pool
   **/
	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

func newPool() *Pool {
	return &Pool{
		pool:       make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (p *Pool) find(key string) *Client {
	return p.pool[key]
}

func (p *Pool) run() {
  log.Println("start running pool ...")
	for {
		select {
		case c := <-p.register:
			log.Println("start adding client from register channel")
			p.pool[c.id] = c

		case c := <-p.unregister:
			log.Println("start removing client from unregister channel")
      delete(p.pool, c.id)

		case m := <-p.broadcast:
			log.Println("start receiving message from broadcast channel")
      log.Printf("received message: %v \n", m)
      log.Printf("size of pool: %v", len(p.pool))
			for _, c := range p.pool {
        log.Printf("sending message to client (%v) in loop", c)
			  c.send <-m
			}
		}
	}
}
