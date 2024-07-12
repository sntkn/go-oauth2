package flashmessage

import (
	"github.com/gin-gonic/gin"
	"github.com/sntkn/go-oauth2/oauth2/internal/session"
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

func getFlashMessages(c *gin.Context, s *session.Session) (Messages, error) {
	var messages Messages
	if err := s.GetNamedSessionData(c, "flashMessage", &messages); err != nil {
		return messages, err
	}
	return messages, nil
}

func setFlashMessages(c *gin.Context, s *session.Session, messages Messages) error {
	if err := s.SetNamedSessionData(c, "flashMessage", messages); err != nil {
		return err
	}
	return nil
}

func Add(c *gin.Context, s *session.Session, t MessageType, m string) error {

	messages, err := getFlashMessages(c, s)
	if err != nil {
		return err
	}

	switch t {
	case Success:
		messages.Success = append(messages.Success, m)
	case Notice:
		messages.Notice = append(messages.Notice, m)
	case Error:
		messages.Error = append(messages.Error, m)
	}

	return setFlashMessages(c, s, messages)
}

func Flash(c *gin.Context, s *session.Session) (Messages, error) {
	messages, err := getFlashMessages(c, s)
	if err != nil {
		return Messages{}, err
	}

	setFlashMessages(c, s, Messages{})

	return messages, nil
}
