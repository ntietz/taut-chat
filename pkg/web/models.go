package web

type User struct {
	ID           string   `json:"id"`
	Handle       string   `json:"handle"`
	PasswordHash string   `json:"password_hash"`
	Following    []string `json:"following"`
	Credits      int64    `json:"credits"`
}

type Channel struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type Message struct {
	ID string `json:"id"`

	// either UserID or ChannelID must be present
	UserID    *string `json:"user_id"`
	ChannelID *string `json:"channel_id"`

	Content string `json:"content"`
}
