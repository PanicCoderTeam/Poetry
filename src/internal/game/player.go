package game

import (
	"poetry/src/internal/domain/entity"

	"github.com/lonng/nano/session"
)

type Player struct {
	*entity.GameUser
	IP string // ip地址
	// 玩家数据
	session *session.Session
	State   string //状态
	Room    *GameRoom
}

func NewPlayer(gameUser *entity.GameUser, ip string, state string) *Player {
	return &Player{
		GameUser: gameUser,
		IP:       ip,
		State:    state,
	}
}

func (p *Player) bindSession(s *session.Session) {
	p.session = s
	p.session.Set(KEY_CUR_PLAYER, p)
}
