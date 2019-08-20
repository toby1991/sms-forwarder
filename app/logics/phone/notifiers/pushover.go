package notifiers

import (
	"github.com/totoval/framework/biu"
	"github.com/totoval/framework/config"
)

const (
	PushOverMessageUrl = "https://api.pushover.net/1/messages.json"
)

type Pushover struct {
}

func (po *Pushover) Notify(sender, content string) error {
	_, err := biu.Ready(biu.MethodPost, PushOverMessageUrl, &biu.Options{
		Body: &biu.Body{
			"token":  config.GetString("pushover.token"),
			"user":   config.GetString("pushover.user"),
			"device": config.GetString("pushover.device"),

			"title":   sender,
			"message": content,
		},
	}).Biu().Status()
	return err
}