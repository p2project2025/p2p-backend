package chat

import (
	"p2p/models"
	"p2p/repo/chats"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatServiceInterface interface {
	CreateChat(chatData *models.Chat) error
	GetChatsBetweenUsers(userA, userB primitive.ObjectID) ([]models.Chatres, error)
	GetUniqueChatUsers(userID primitive.ObjectID) ([]models.User, error)
}

type ChatService struct{}

func (s *ChatService) CreateChat(chatData *models.Chat) error {
	repo := chats.ChatRepoInterface(&chats.ChatRepo{})
	return repo.CreateChat(chatData)
}

func (s *ChatService) GetChatsBetweenUsers(userA, userB primitive.ObjectID) ([]models.Chatres, error) {
	repo := chats.ChatRepoInterface(&chats.ChatRepo{})
	return repo.GetChatsBetweenUsers(userA, userB)
}

func (s *ChatService) GetUniqueChatUsers(userID primitive.ObjectID) ([]models.User, error) {
	repo := chats.ChatRepoInterface(&chats.ChatRepo{})
	return repo.GetUniqueChatUsers(userID)
}
