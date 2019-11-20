package chat

// Message struct
type Message struct {
	// Email    string `json:"email"`
	Username string `json:"username"`
	Command  string `json:"command"` // The command to run, eg "nick", "exit"
	Message  string `json:"message"` // The message body
}
