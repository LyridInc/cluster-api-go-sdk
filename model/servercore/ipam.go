package servercore

type IP struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	IP      string `json:"ip"`
	UserID  int    `json:"user_id"`
	TTL     int    `json:"ttl"`
}
