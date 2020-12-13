package main

import (
	"time"

	"github.com/sanposhiho/mol/pkg/actor"
)

func main() {
	supervisor := actor.NewActor(nil)
	supervisor.Tell("started supervisor!")
	child := actor.NewActor(supervisor)
	child.Tell("started child actor!")
	child.Tell("This is crush message!")
	child.Tell("Child restarted successfully")
	time.Sleep(2 * time.Second)
	return
}
