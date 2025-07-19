package game

import (
	"poetry/pb/game_room"
	"poetry/src/internal/domain/entity"
	"sync"

	"github.com/lonng/nano"
)

type GameRoom struct {
	*entity.GameRoom
	Players []*game_room.Player
	Rounds  []*Round
	mu      sync.Mutex
	group   *nano.Group
}

func (r *GameRoom) AddPlayer(p *game_room.Player) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p1 := range r.Players {
		if p1.UserId == p.UserId {
			return
		}
	}
	r.Players = append(r.Players, p)
}

func (r *GameRoom) RemovePlayer(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, p := range r.Players {
		if p.UserId == id {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			break
		}
	}
}

func (r *GameRoom) GetPlayer(userId string) *game_room.Player {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.Players {
		if p.UserId == userId {
			return p
		}
	}
	return nil
}
