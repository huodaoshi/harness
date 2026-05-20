package application

import "sync/atomic"

// FakeChatCallCounter tracks fake ChatModel invocations (Spike S1/S2).
type FakeChatCallCounter struct {
	n atomic.Int64
}

func (c *FakeChatCallCounter) Inc() {
	c.n.Add(1)
}

func (c *FakeChatCallCounter) Load() int64 {
	return c.n.Load()
}

func (c *FakeChatCallCounter) Reset() {
	c.n.Store(0)
}
