package chatmodel

import (
	"context"
	"strings"
)

// FakeGateway implements Gateway without calling external APIs (tests / dev default).
type FakeGateway struct{}

func NewFakeGateway() *FakeGateway {
	return &FakeGateway{}
}

func (f *FakeGateway) Generate(ctx context.Context, req Request) (string, error) {
	return f.reply(req), nil
}

func (f *FakeGateway) Stream(ctx context.Context, req Request, onToken func(text string) error) error {
	for _, r := range f.reply(req) {
		if err := onToken(string(r)); err != nil {
			return err
		}
	}
	return nil
}

func (f *FakeGateway) reply(req Request) string {
	msg := strings.TrimSpace(req.Message)
	var base string
	switch req.Mode {
	case "distress":
		if msg == "" {
			base = "我在这里，你愿意多说一点吗？"
		} else {
			base = "【洪峰】我听到了。此刻很难受是正常的，我们先慢慢说。"
		}
	default:
		if msg == "" {
			base = "我在这里，你愿意多说一点吗？"
		} else {
			base = "【普聊】我听到你了：" + msg
		}
	}
	if block := strings.TrimSpace(req.InjectBlock); block != "" {
		return "【已读上下文】\n" + block + "\n---\n" + base
	}
	return base
}
