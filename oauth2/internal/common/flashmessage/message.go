package flashmessage

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/common/session"
)

type MessageType string

const (
	Success MessageType = "success"
	Notice  MessageType = "notice"
	Error   MessageType = "error"
)

type Messages struct {
	Success []string
	Notice  []string
	Error   []string
}

func getFlashMessages(c *gin.Context, s session.SessionClient) (Messages, error) {
	messages, ok, err := session.Load[Messages](c, s, "flashMessage")
	if err != nil {
		return Messages{}, err
	}
	if !ok {
		return Messages{}, nil
	}
	return messages, nil
}

func setFlashMessages(c *gin.Context, s session.SessionClient, messages Messages) error {
	return session.Save(c, s, "flashMessage", messages)
}

func AddMessage(c *gin.Context, s session.SessionClient, t MessageType, message string) error {
	messages, err := getFlashMessages(c, s)
	if err != nil {
		return err
	}

	switch t {
	case Success:
		messages.Success = append(messages.Success, message)
	case Notice:
		messages.Notice = append(messages.Notice, message)
	case Error:
		messages.Error = append(messages.Error, message)
	}

	return setFlashMessages(c, s, messages)
}

func Flash(c *gin.Context, s session.SessionClient) (*Messages, error) {
	messages, ok, err := session.Pop[Messages](c, s, "flashMessage")
	if err != nil {
		return nil, err
	}
	if !ok {
		return &Messages{}, nil
	}
	return &messages, nil
}
