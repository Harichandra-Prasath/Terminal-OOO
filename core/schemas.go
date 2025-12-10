package core

type CreateRoomRequest struct {
	PlayerName string `json:"player_name"`
	MaxSize    int    `json:"max_size"`
	MinSize    int    `json:"min_size"`
	RoundTime  int    `json:"round_time"`
}

type JoinRoomRequest struct {
	PlayerName string `json:"player_name"`
	RoomId     string `json:"room_id"`
}

type StartRoomRequest struct {
	HostId string `json:"host_id"`
	RoomId string `json:"room_id"`
}
