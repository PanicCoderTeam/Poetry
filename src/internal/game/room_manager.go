package game

import "sync"

type GameRoomManager struct {
	GameRooms sync.Map
}

func NewGameRoomManager() *GameRoomManager {
	return &GameRoomManager{
		GameRooms: sync.Map{},
	}
}

var defaultGameRoomManagerInfo = NewGameRoomManager()

func (r *GameRoomManager) GetGameRoom(id string) *GameRoom {
	gameRoom, err := r.GameRooms.Load(id)
	if !err {
		return nil
	}
	gameRoomInfo := gameRoom.(*GameRoom)
	return gameRoomInfo
}

func (r *GameRoomManager) PutGameRoom(id string, gameRoom *GameRoom) {
	r.GameRooms.Store(id, gameRoom)
}

func (r *GameRoomManager) DeleteGameRoom(id string) {
	r.GameRooms.Delete(id)
}
