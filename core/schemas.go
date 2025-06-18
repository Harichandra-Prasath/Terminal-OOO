package core

type CreateRoomRequest struct {
	PlayerName string `json:"player_name"`
	MaxSize    int    `json:"max_size"`
	MinSize    int    `json:"min_size"`
	RoundTime  int    `json:"round_time"`
}
