package models

type BotRequest struct {
	Content string `json:"content"`
	RoomID string `json:"room_id"`
	RoomName string `json:"room_name"`
	UserID string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type BotResponse struct {
	Content string `json:"content"`
	RoomID string `json:"room_id"`
	RoomName string `json:"room_name"`
	UserID string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}
