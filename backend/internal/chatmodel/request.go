package chatmodel

// Request is input for ChatModelGateway (pass path only).
type Request struct {
	Mode        string
	Message     string
	InjectBlock string
}
