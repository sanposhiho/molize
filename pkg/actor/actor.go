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
	if parent != nil {
		parent.AddChild(a)
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
	fmt.Printf("start: %s\n", a.address)
	defer func() {
		a.notifyPanic()
		if err := recover(); err != nil {
			fmt.Println("Recovered. err:", err)
		}
	}()
	for {
		select {
		case message := <-a.mailbox:
			fmt.Printf("message received: %s\n", message)
			if message == "This is crush message!" {
				panic("crush!")
			}
		case address := <-a.childPanicReceiver:
			a.receiveChildPanic(address)
		}
	}
}

func (a *Actor) notifyPanic() {
	if a.parent != nil {
		a.parent.childPanicReceiver <- a.address
	} else {
		fmt.Print("no")
	}
}

func (a *Actor) receiveChildPanic(childAddress string) {
	fmt.Printf("recover child: %s\n", a.address)
	a.children[childAddress].Run()
}
