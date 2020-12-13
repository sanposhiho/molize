# molize

molize is library to realize actor model based applications with Go.

Inspired by Erlang and Akka

## What is Actor Model

**TDB**

## Features

**TDB**

## Sample

```
func main() {
	// create supervisor
	supervisor := actor.NewActor(nil, actor.Options{})

	// create child
	child := actor.NewActor(supervisor, actor.Options{
		// Wait this time when retry
		RetryWaitTime: time.Duration(1 * time.Second),
		// Retry processing a message, when panic
		MaxRetryCount: 2,
		// After retry, restart actor
		Restart: true,
	})

	// define child's behavior
	child.React(func(message message.Message) {
		switch message {
		case "This is crush message!":
			panic("crush!")
		default:
			fmt.Printf("message received: %s\n", message)
		}
	})

	// send message to child
	child.Tell("Started child actor!")

	// send message make panic on actor.
	// with some options, retry automatically.
	child.Tell("This is crush message!")

	// send message to child successfully, because of supervisor's recovering
	child.Tell("Check actor restarted successfully")
	time.Sleep(6 * time.Second)
	return
}
```
