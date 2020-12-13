package message

type Message interface{}

func NewMessage(contents interface{}) Message {
	return Message(contents)
}
