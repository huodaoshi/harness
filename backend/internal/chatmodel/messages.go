package chatmodel

import (
	"strings"

	"github.com/cloudwego/eino/schema"
)

func buildMessages(req Request) []*schema.Message {
	system := SystemPromptForMode(req.Mode)
	if block := strings.TrimSpace(req.InjectBlock); block != "" {
		system += "\n\n" + block
	}
	return []*schema.Message{
		{Role: schema.System, Content: system},
		{Role: schema.User, Content: strings.TrimSpace(req.Message)},
	}
}
