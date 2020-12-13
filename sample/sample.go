package sample

import "github.com/sanposhiho/mol/pkg/actor"

func main() {
	supervisor := actor.NewActor(nil)
	child := actor.NewActor(supervisor)
	child.Tell("hoge")
	return
}
