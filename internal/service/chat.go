package service

type Message string

type ChatService interface {
	Connect(userId string) error
	Disconnect(userId string) error
	SendMessage(fromId, toId string, message Message) error
	IsClientConnected(userId string) bool
}
