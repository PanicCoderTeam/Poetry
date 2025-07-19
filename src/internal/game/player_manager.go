package game

import (
	"poetry/src/pkg/log"
	"time"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/session"
)

const (
	exitChannelNum  = 16
	resetChannelNum = 16
)

var defaultPlayerManager = NewPlayerManager()

type PlayerManager struct {
	component.Base
	group   *nano.Group        // 广播channel
	players map[string]*Player // 所有的玩家
	chKick  chan string        // 退出队列
	chReset chan string        // 重置队列
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		group:   nano.NewGroup("_SYSTEM_MESSAGE_BROADCAST"),
		players: map[string]*Player{},
		chKick:  make(chan string, exitChannelNum),
		chReset: make(chan string, resetChannelNum),
	}
}

func (m *PlayerManager) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		m.group.Leave(s)
	})

	// 处理踢出玩家和重置玩家消息(来自http)
	scheduler.NewTimer(time.Second, func() {
	ctrl:
		for {
			select {
			case uid := <-m.chKick:
				p, ok := defaultPlayerManager.player(uid)
				if !ok || p.session == nil {
					log.Errorw("玩家不在线", "uid", uid)
				}
				p.session.Close()
				log.Infow("踢出玩家", "uid", uid)

			case uid := <-m.chReset:
				p, ok := defaultPlayerManager.player(uid)
				if !ok {
					return
				}
				if p.session != nil {
					log.Errorw("玩家正在游戏中，不能重置: ", "uid", uid)
					return
				}
				p.Room = nil
				log.Infow("重置玩家, UID", "uid", uid)

			default:
				break ctrl
			}
		}
	})
}

// func (m *PlayerManager) Login(s *session.Session, req *protocol.LoginToGameServerRequest) error {
// 	uid := req.Uid
// 	s.Bind(uid)

// 	log.Infof("玩家: %d登录: %+v", uid, req)
// 	if p, ok := m.player(uid); !ok {
// 		log.Infof("玩家: %d不在线，创建新的玩家", uid)
// 		p = newPlayer(s, uid, req.Name, req.HeadUrl, req.IP, req.Sex)
// 		m.setPlayer(uid, p)
// 	} else {
// 		log.Infof("玩家: %d已经在线", uid)
// 		// 重置之前的session
// 		if prevSession := p.session; prevSession != nil && prevSession != s {
// 			// 移除广播频道
// 			m.group.Leave(prevSession)

// 			// 如果之前房间存在，则退出来
// 			if p, err := playerWithSession(prevSession); err == nil && p != nil && p.desk != nil && p.desk.group != nil {
// 				p.desk.group.Leave(prevSession)
// 			}

// 			prevSession.Clear()
// 			prevSession.Close()
// 		}

// 		// 绑定新session
// 		p.bindSession(s)
// 	}

// 	// 添加到广播频道
// 	m.group.Add(s)

// 	res := &protocol.LoginToGameServerResponse{
// 		Uid:      s.UID(),
// 		Nickname: req.Name,
// 		Sex:      req.Sex,
// 		HeadUrl:  req.HeadUrl,
// 		FangKa:   req.FangKa,
// 	}

// 	return s.Response(res)
// }

func (m *PlayerManager) player(uid string) (*Player, bool) {
	p, ok := m.players[uid]

	return p, ok
}

func (m *PlayerManager) setPlayer(uid string, p *Player) {
	if _, ok := m.players[uid]; ok {
		log.Errorw("玩家已经存在，正在覆盖玩家， UID", "uid", uid)
	}
	m.players[uid] = p
}

// func (m *PlayerManager) CheckOrder(s *session.Session, msg *protocol.CheckOrderReqeust) error {
// 	log.Infof("%+v", msg)

// 	return s.Response(&protocol.CheckOrderResponse{
// 		FangKa: 20,
// 	})
// }

// func (m *PlayerManager) sessionCount() int {
// 	return len(m.players)
// }

func (m *PlayerManager) offline(uid string) {
	delete(m.players, uid)
	log.Debugw("从在线列表中删除", "uid", uid, "在线人数", len(m.players))
}
