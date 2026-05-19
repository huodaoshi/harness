package session

import "github.com/huodaoshi/harness/backend/internal/safety"

// Input is the graph input for a single streamed turn.
type Input struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Mode    string `json:"mode"`
}

// RoutedInput carries gate outcome and original input past SafetyGate.
type RoutedInput struct {
	Input Input
	Gate  safety.Result
}

// EnrichedChatInput is pass-path state after ProfileInject.
type EnrichedChatInput struct {
	Routed      RoutedInput
	InjectBlock string
}

// TurnOutput is the graph result for one turn (crisis branch or chat branch).
type TurnOutput struct {
	Crisis      *safety.CrisisPayload
	Chat        string
	ChatUsed    bool
	InjectBlock string
}

// TurnOutcome is returned to the HTTP layer after one turn.
type TurnOutcome struct {
	Crisis      *safety.CrisisPayload
	Chat        string
	ChatCalls   int64
	InjectBlock string
}
