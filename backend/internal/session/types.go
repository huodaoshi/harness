package session

// Input is the graph input for a single streamed turn (Spike S1).
type Input struct {
	Message string `json:"message"`
	Mode    string `json:"mode"`
}
