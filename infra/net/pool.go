package net

import (
	"log"
)

type IPool interface {
  find(key string) IClient
  Register(client IClient)
  Unregister(client IClient)
  Broadcast(message *Message)
  Unicast(message *Message)
  Size() int
  Run()
}

type Pool struct {
	pool map[string]IClient

  /**
   * broadcast message to all clients in this pool
   **/
	broadcast chan *Message

  unicast chan *Message

	register chan IClient

	unregister chan IClient
}

func NewPool() *Pool {
	return &Pool{
		pool:       make(map[string]IClient),
		broadcast:  make(chan *Message),
		unicast:  make(chan *Message),
		register:   make(chan IClient),
		unregister: make(chan IClient),
	}
}

func (p *Pool) find(key string) IClient {
	return p.pool[key]
}

func (p *Pool) Register(client IClient) {
  p.register <- client
}

func (p *Pool) Unregister(client IClient) {
  p.unregister <- client
}

func (p *Pool) Broadcast(message *Message) {
  p.broadcast <- message
}

func (p *Pool) Unicast(message *Message) {
  p.unicast <- message
}

func (p *Pool) Size() int {
  return len(p.pool)
}

func (p *Pool) Run() {
  log.Println("start running pool ...")
	for {
		select {
		case c := <-p.register:
			log.Println("start adding client from register channel")
			p.pool[c.Id()] = c
      log.Printf("register client to pool: %v", c.Id())

		case c := <-p.unregister:
			log.Println("start removing client from unregister channel")
      delete(p.pool, c.Id())

		case m := <-p.unicast:
			log.Println("start receiving message from unicast channel")
      log.Printf("received message: %v \n", m)
      log.Printf("target client id to be found: %v \n", m.receiver.Id())
      targetClient := p.find(m.receiver.Id())
      log.Printf("located target client to unicast: %v \n", targetClient)
      targetClient.Send(m)

		case m := <-p.broadcast:
			log.Println("start receiving message from broadcast channel")
      log.Printf("received message: %v \n", m)
      log.Printf("size of pool: %v", len(p.pool))
			for _, c := range p.pool {
        log.Printf("sending message to client (%v) in loop", c)
			  c.Send(m)
			}
		}
	}
}
