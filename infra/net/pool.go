package net

import (
	"log"
)

type Pool struct {
	pool map[string]*Client

  /**
   * broadcast message to all clients in this pool
   **/
	broadcast chan *Message

  unicast chan *Message

	register chan *Client

	unregister chan *Client
}

func NewPool() *Pool {
	return &Pool{
		pool:       make(map[string]*Client),
		broadcast:  make(chan *Message),
		unicast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (p *Pool) find(key string) *Client {
	return p.pool[key]
}

func (p *Pool) Register(client *Client) {
  p.register <- client
}

func (p *Pool) Unregister(client *Client) {
  p.unregister <- client
}

func (p *Pool) Run() {
  log.Println("start running pool ...")
	for {
		select {
		case c := <-p.register:
			log.Println("start adding client from register channel")
			p.pool[c.id] = c
      log.Printf("register client to pool: %v", c.id)

		case c := <-p.unregister:
			log.Println("start removing client from unregister channel")
      delete(p.pool, c.id)

		case m := <-p.unicast:
			log.Println("start receiving message from unicast channel")
      log.Printf("received message: %v \n", m)
      log.Printf("target client id to be found: %v \n", m.receiver.id)
      targetClient := p.find(m.receiver.id)
      log.Printf("located target client to unicast: %v \n", targetClient)
      targetClient.send <-m

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
