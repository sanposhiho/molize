package actor

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sanposhiho/molize/pkg/message"
)

type Actor struct {
	mailbox            chan message.Message
	childPanicReceiver chan string // receive child's address
	behavior           func(message message.Message)
	state              *state
	parent             *Actor
	children           map[string]*Actor // A Key is child's address
	address            string
	options            Options
}

type state struct {
	// recentMessage has a message processing recently
	recentMessage message.Message
	// retryCount has the count of retry processing the message.
	retryCount int
}

type Options struct {
	// Wait this time when retry
	RetryWaitTime time.Duration
	// Retry processing a message, when panic
	MaxRetryCount int
	// After retry, restart actor
	Restart bool
}

func NewActor(parent *Actor, options Options) *Actor {
	a := &Actor{
		mailbox:            make(chan message.Message),
		childPanicReceiver: make(chan string),
		parent:             parent,
		children:           map[string]*Actor{},
		address:            uuid.New().String(),
		state:              &state{},
		options:            options,
	}
	if parent != nil {
		parent.AddChild(a)
	}

	a.Run()

	return a
}

func (a *Actor) React(behavior func(message message.Message)) {
	a.behavior = behavior
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
			if a.behavior == nil {
				fmt.Print("no behaviour defined")
			}
			a.state.setRecentMessage(message)
			a.behavior(message)
		case address := <-a.childPanicReceiver:
			a.receiveChildPanic(address)
		}
	}
}

func (a *Actor) notifyPanic() {
	if a.parent != nil {
		a.parent.childPanicReceiver <- a.address
	}
}

func (a *Actor) receiveChildPanic(childAddress string) {
	fmt.Printf("start recovering child...: %s\n", a.address)
	a.children[childAddress].restart()
}

func (s *state) setRecentMessage(m message.Message) {
	s.recentMessage = m
}

func (s *state) incrementRetryCount() {
	s.retryCount++
}

func (a *Actor) restart() {
	for ; a.state.retryCount < a.options.MaxRetryCount; a.state.incrementRetryCount() {
		fmt.Printf("retry message:%s, retry count:%d\n", a.address, a.state.retryCount)
		time.Sleep(a.options.RetryWaitTime)
		a.retryRecentMessage()
	}
	if a.options.Restart {
		a.Run()
	}
}

func (a *Actor) retryRecentMessage() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered. err:", err)
		}
	}()
	a.behavior(a.state.recentMessage)
}
