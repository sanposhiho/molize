package actor

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sanposhiho/mol/pkg/message"
)

type Actor struct {
	mailbox            chan message.Message
	childPanicReceiver chan string // send child's address
	parent             *Actor
	children           map[string]*Actor // A Key is child's address
	address            string
}

func NewActor(parent *Actor) *Actor {
	a := &Actor{
		mailbox:            make(chan message.Message),
		childPanicReceiver: make(chan string),
		parent:             parent,
		children:           map[string]*Actor{},
		address:            uuid.New().String(),
	}

	a.Run()

	return a
}

func (a *Actor) AddChild(child *Actor) {
	a.children[child.address] = child
	return
}

func (a *Actor) Tell(message message.Message) {
	a.mailbox <- message
}

func (a *Actor) Run() {
	go a.processLoop()
	return
}

func (a *Actor) processLoop() {
	defer func() {
		a.notifyPanic()
		if err := recover(); err != nil {
			fmt.Println("Recovered. err:", err)
		}
	}()
	for {
		select {
		case message := <-a.mailbox:
			fmt.Print(message)
		case address := <-a.childPanicReceiver:
			a.receiveChildPanic(address)
		}

		return
	}
}

func (a *Actor) notifyPanic() {
	a.parent.childPanicReceiver <- a.address
}

func (a *Actor) receiveChildPanic(childAddress string) {
	a.children[childAddress].Run()
}
