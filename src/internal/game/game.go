package game

import (
	"context"
	"net/http"
	"poetry/pb/game_room"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/domain/service"
	"poetry/src/internal/infrastructure/repository"
	"poetry/src/pkg/async"
	"poetry/src/pkg/basic"
	"poetry/src/pkg/log"
	"poetry/src/pkg/trpc/codec/capi_error"
	"poetry/src/pkg/trpc/filter"
	"poetry/src/pkg/utils"
	"strings"
	"sync"
	"time"

	serviceimpl "poetry/src/internal/domain/service/service_impl"

	ejson "encoding/json"

	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/pipeline"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/serialize/json"
	"github.com/lonng/nano/session"
)

func NewGameHandler() *GameHandler {
	return &GameHandler{}
}

type GameHandler struct {
	component.Base
	GameRoomService service.GameRoomService
	GameUserService service.GameUserService
	PoetryService   service.PoetryService
}

var WORDS = []string{"春", "秋", "月", "风", "花", "雪", "云", "雨", "日", "夜", "江", "山", "天", "星", "霜", "露", "霞", "烟", "河", "海", "东", "南", "西", "北", "中", "上", "下", "前", "后", "左", "右", "内", "外", "边", "关", "塞", "野", "原", "峰", "谷", "朝", "夕", "晨", "昏", "晓", "午", "晚", "年", "节", "时", "刻", "分", "秒", "今", "古", "昨", "明", "夏", "冬", "昼", "柳", "梅", "桃", "李", "荷", "菊", "松", "竹", "兰", "草", "木", "叶", "枝", "莺", "雁", "燕", "鹤", "蝉", "马", "鱼", "情", "愁", "恨", "泪", "思", "念", "心", "魂", "梦", "酒", "诗", "书", "剑", "歌", "舞", "笑", "泣", "醉", "醒", "家"}

func Startup() {
	components := &component.Components{}
	// 流量统计
	pip := pipeline.New()
	var stats = &Stats{}
	// 入队 Outbound pipeline
	pip.Outbound().PushBack(stats.Outbound)
	// 入队 Inbound pipeline
	pip.Inbound().PushBack(stats.Inbound)
	// crypto := game.NewCrypto()
	// pip.Inbound().PushBack(crypto.Inbound)
	// pip.Outbound().PushBack(crypto.Outbound)
	// 注册下流量统计组件
	components.Register(stats, component.WithName("stats"))
	roomGame := &GameHandler{
		GameRoomService: serviceimpl.NewGameRoomServiceImpl(serviceimpl.GameRoomServiceImplOption{
			GameRoomRepo: repository.NewGameRoomRepository(),
		}),
		GameUserService: serviceimpl.NewGameUserServiceImpl(serviceimpl.GameUserServiceImplOption{
			GameUserRepo: repository.NewGameUserRepository(),
		}),
		PoetryService: serviceimpl.NewPoetryServiceImpl(serviceimpl.PoetryServiceImplOption{
			PoetryRepo: repository.NewPoetryRepository(),
		}),
	}
	components.Register(roomGame, component.WithName("game"), component.WithNameFunc(func(s string) string {
		return strings.ToLower(s)
	}))
	session.Lifetime.OnClosed(roomGame.UserDisconnected)
	nano.Listen(":3250", // 端口号
		nano.WithIsWebsocket(true), // 是否使用 websocket
		nano.WithPipeline(pip),     // 是否使用 pipeline
		nano.WithCheckOriginFunc(func(_ *http.Request) bool { return true }), // 允许跨域
		nano.WithWSPath("/nano"),                  // websocket 连接地址
		nano.WithDebugMode(),                      // 开启 debug 模式
		nano.WithSerializer(json.NewSerializer()), // 使用 json 序列化器
		nano.WithComponents(components),           // 加载组件
	)

}

// 握手时检查一下Token，确认用户已经登录，并校验房间号和密码，直接加入房间
func (h *GameHandler) OnHandshake(s *session.Session, data map[string]interface{}) error {
	// 提取头部信息（例如从 data["user"] 中解析）
	headers := data["user"].(map[string]interface{})
	_, err := filter.CheckToken(headers["Authorization"].(string))
	if err != nil {
		return err
	}
	return nil
}

func (h *GameHandler) Init() {
	ctx := context.Background()
	// 初始化游戏处理器
	log.DebugEx(ctx, "初始化游戏处理器")
}
func (h *GameHandler) AfterInit() {
	ctx := context.Background()
	log.DebugEx(ctx, "AfterInit add scheduler")
	// 初始化后处理
	scheduler.NewTimer(24*time.Hour, func() {
		log.DebugEx(ctx, "开始定时任务执：")
		async.Run(func() {
			log.DebugEx(ctx, "定时任务")
			// 1. 获取所有房间
			_, rooms, err := h.GameRoomService.DescribeGameRoom(ctx, nil, nil, []string{entity.ROOMSTATE_PLAYING, entity.PLAYER_STATE_WAITING}, nil, nil)
			if err != nil {
				log.ErrorEx(ctx, "获取所有房间失败", "err", err)
				return
			}
			// 2. 遍历房间
			outtimeRoomList := make([]*entity.GameRoom, 0)
			for _, room := range rooms {
				if room.CreatedAt.Before(time.Now().Add(-24 * time.Hour)) {
					outtimeRoomList = append(outtimeRoomList, room)
					defaultGameRoomManagerInfo.DeleteGameRoom(room.RoomID)
				}
			}

			err = h.GameRoomService.DeleteGameRoom(ctx, outtimeRoomList)
			if err != nil {
				log.ErrorEx(ctx, "删除过期房间失败", "err", err)
				return
			}
		})
	})
}
// func (h *GameHandler) OnClosed(s *session.Session, reason error) {
// 	ctx := context.Background()
// 	log.DebugEx(ctx, "OnClosed:", "reason", reason, "uid", s.UID())
// }

func (h *GameHandler) UserDisconnected(s *session.Session) {

	// 服务关闭前处理
	// if err := h.group.Leave(s); err != nil {
	// 	log.ErrorEx(ctx,"Remove user from group failed", "err", err, "uid", s.UID())
	// }
	ctx := context.Background()
	log.DebugEx(ctx, "UserDisconnected:", "session", s, "info", s.State(), "uid", s.UID())
	playerInfo, err := playerWithSession(s)
	if err != nil {
		log.ErrorEx(ctx, "PlayerWithSessionERRor", "err", err)
		return
	}
	playerInfo, exist := defaultPlayerManager.player(playerInfo.UserID)
	if !exist {
		log.ErrorEx(ctx, "玩家不存在", "玩家不存在", playerInfo.UserID)
		return
	}
	round := playerInfo.Room.Rounds[len(playerInfo.Room.Rounds)-1]

	defaultPlayerManager.group.Leave(s)
	defaultPlayerManager.offline(playerInfo.UserID)
	playerInfo.Room.RemovePlayer(playerInfo.UserID)
	playerInfo.Room.group.Leave(s)
	s.Clear()
	s.Close()
	if round.State == ROUND_STATE_PLAYING {
		h.updateGameRoomInfo(ctx, playerInfo, round)
	}
	playerInfo.Room.group.Broadcast("game.userdisconnect", playerInfo.Room)
}

// func (h *GameHandler) SyncMessage(s *session.Session, msg []byte) error {
// 	if err := s.RPC("GameHandler.Stats", &protocol.MasterStats{Uid: s.UID()}); err != nil {
// 		return err
// 	}
// 	// Sync message to all members in this room
// 	return h.group.Broadcast("onMessage", msg)
// 	// 服务关闭前处理
// }

// 加入游戏
func (h *GameHandler) Joingame(s *session.Session, req *game_room.NanoJoinGameRoomRequest) error {
	// 实现加入游戏逻辑
	// 1. 创建选手
	// 1.1 判断选手是否已经在游戏房间中
	// 1.2 添加选手到player
	// 1.3 添加选手到广播中
	// 1.3 添加选手到房间
	// 获取房间号以及密码
	ctx := context.Background()
	// 校验用户信息
	userClaimInfo, err := filter.CheckToken(req.Token)
	if err != nil {
		utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "Token无效", nil, nil)
		return err
	}
	//校验房间信息
	roomId := req.RoomId
	password := req.Password
	if roomId == "" || password == "" {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "缺少必要参数", nil, nil)
	}

	_, gameRoomList, err := h.GameRoomService.DescribeGameRoom(ctx, []string{roomId}, nil, nil, nil, nil)
	if err != nil {
		return err
	}
	if len(gameRoomList) == 0 {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "房间不存在", nil, nil)
	}

	if gameRoomList[0].Password != password {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "密码错误", nil, nil)
	}
	player, exist := defaultPlayerManager.player(userClaimInfo.UserID)
	if exist {
		log.DebugEx(ctx, "玩家已经存在，正在覆盖玩家， UID", "userId", userClaimInfo.UserID)
		if preSession := player.session; preSession != nil && preSession != s {
			defaultPlayerManager.group.Leave(preSession)
			playerInfo, err := playerWithSession(preSession)
			if err == nil && playerInfo != nil && playerInfo.Room != nil && playerInfo.Room.group != nil {
				playerInfo.Room.group.Leave(preSession)
			}
			preSession.Clear()
			preSession.Close()
		}
	} else {
		_, users, err := h.GameUserService.DescribeGameUserInfo(ctx, []*string{&userClaimInfo.UserID}, nil, nil, nil, nil)
		if err != nil {
			log.ErrorEx(ctx, "DescribeGameUserInfo Error", "err", err)
			return err
		}
		if len(users) == 0 {
			log.ErrorEx(ctx, "用户不存在", "UserInfo", userClaimInfo)
			return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "用户不存在", "UserInfo", userClaimInfo)
		}
		player = NewPlayer(users[0].GameUser, s.RemoteAddr().String(), entity.PLAYER_STATE_READY)
	}
	rounds := []*Round{}
	ejson.Unmarshal([]byte(gameRoomList[0].Rounds), &rounds)
	//获取房间信息
	gameRoomInfo := defaultGameRoomManagerInfo.GetGameRoom(roomId)
	if gameRoomInfo == nil {
		gameRoomInfo = &GameRoom{
			GameRoom: gameRoomList[0],
			Rounds:   rounds,
			Players:  []*game_room.Player{},
			mu:       sync.Mutex{},
			group:    nano.NewGroup(roomId),
		}
		defaultGameRoomManagerInfo.PutGameRoom(roomId, gameRoomInfo)
	}
	player.Room = gameRoomInfo
	log.DebugEx(ctx, "join Game", "player", player)
	player.bindSession(s)
	defaultPlayerManager.setPlayer(userClaimInfo.UserID, player)
	gameRoomInfo.group.Add(s)
	gameRoomInfo.AddPlayer(&game_room.Player{
		UserId:   player.UserID,
		Username: player.Username,
		State:    player.State,
	})
	log.DebugEx(ctx, "加入房间成功", "gameRoomInfo", gameRoomInfo)
	//给客户端推送用户的加入信息, 用户更新房间信息
	gameRoomInfo.group.Broadcast("game.joingame", gameRoomInfo)
	return nil
}

func (h *GameHandler) ReadyGame(s *session.Session, readyGameReq *game_room.ReadyGameRequest) error {
	ctx := context.Background()
	// 实现加入游戏逻辑
	log.DebugEx(ctx, "ReadyGame:", "req", readyGameReq)
	playerInfo, err := playerWithSession(s)
	if err != nil {
		return err
	}
	playerInfo, exist := defaultPlayerManager.player(playerInfo.UserID)
	if !exist {
		log.ErrorEx(ctx, "玩家不存在", "玩家不存在", playerInfo.UserID)
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "玩家不存在", "玩家不存在", playerInfo.UserID)
	}
	log.DebugEx(ctx, "ReadyGame:", "playerInfo", playerInfo)
	playerInfo.State = readyGameReq.State
	if playerInfo.Room != nil && playerInfo.Room.Status == entity.ROOMSTATE_PLAYING {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "游戏已经开始", "游戏已经开始", playerInfo)
	}
	
	playerInfo.Room.GetPlayer(playerInfo.UserID).State = readyGameReq.State
	if playerInfo.Room != nil && playerInfo.Room.group != nil {
		playerInfo.Room.group.Broadcast("game.readygame", playerInfo)
	}
	return nil
}

// 开始游戏
func (h *GameHandler) Startgame(s *session.Session, msg []byte) error {
	ctx := context.Background()
	// 实现开始游戏逻辑
	log.DebugEx(ctx, "开始游戏:", "msg", string(msg))
	playerInfo, err := playerWithSession(s)
	if err != nil {
		return err
	}
	playerInfo, exist := defaultPlayerManager.player(playerInfo.UserID)
	if !exist {
		log.ErrorEx(ctx, "玩家不存在", "玩家不存在", playerInfo.UserID)
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "玩家不存在", "玩家不存在", playerInfo.UserID)
	}
	if playerInfo.Room == nil {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "房间不存在", "用户没有加入任何房间", playerInfo)
	}
	if playerInfo.Room.OwnerID != playerInfo.UserID {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "只有房主才能开始游戏", "用户不是房主", playerInfo)
	}
	if len(playerInfo.Room.Players) < int(playerInfo.Room.MaxPlayers) {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "人数不够不能开始游戏", playerInfo)
	}

	for _, player := range playerInfo.Room.Players {
		if player.State != entity.PLAYER_STATE_READY {
			return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "还有玩家没有准备", "还有玩家没有准备", player.UserId)
		}
	}

	playerInfo.Room.Status = entity.ROOMSTATE_PLAYING
	isNewRound := true
	for _, round := range playerInfo.Room.Rounds {
		if round.State == ROUND_STATE_PLAYING {
			isNewRound = false
		}
	}
	if isNewRound {
		word := ""
		for {
			randNum := basic.RandInt(0, len(WORDS))
			isBreak := true
			for _, round := range playerInfo.Room.Rounds {
				if round.CurrentWord == WORDS[randNum] {
					isBreak = false
					break
				}
			}
			if isBreak {
				word = WORDS[randNum]
				break
			}
		}
		round := &Round{
			CurrentWord:   word,
			CurrentPlayer: playerInfo.Room.Players[0],
			UsedPoems:     []*PlayPoemInfo{},
			Winner:        nil,
			State:         ROUND_STATE_PLAYING,
		}
		playerInfo.Room.Rounds = append(playerInfo.Room.Rounds, round)
	}
	playerInfo.Room.group.Broadcast("game.startgame", playerInfo.Room)
	return nil
}

// 出题
func (h *GameHandler) Submitpoem(s *session.Session, req *game_room.SubmitPoetryRequest) error {
	req.Poetry = utils.ConvertChinsesTraditional2S(req.Poetry)
	// 实现提交诗句逻辑
	ctx := context.Background()
	log.DebugEx(ctx, "Submitpoem:", "req", req)
	playerInfo, err := playerWithSession(s)
	if err != nil {
		return err
	}
	playerInfo, exist := defaultPlayerManager.player(playerInfo.UserID)
	if !exist {
		log.ErrorEx(ctx, "玩家不存在", "玩家不存在", playerInfo.UserID)
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "玩家不存在", "玩家不存在", playerInfo.UserID)
	}
	if playerInfo.Room == nil {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "房间不存在", "用户没有加入任何房间", playerInfo)
	}
	if playerInfo.Room.Status != entity.ROOMSTATE_PLAYING {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "游戏还没有开始", "游戏还没有开始", playerInfo)
	}
	round := playerInfo.Room.Rounds[len(playerInfo.Room.Rounds)-1]
	if round.State != ROUND_STATE_PLAYING {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "游戏已经结束", "游戏已经结束", playerInfo)
	}
	log.DebugEx(ctx, "Submitpoem:", "round", round, "req", req)
	if !strings.Contains(req.Poetry, round.CurrentWord) {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "不能使用该诗句", "该诗句未包含该题目", playerInfo)
	}

	for _, poem := range round.UsedPoems {
		if poem.Poem == req.Poetry {
			return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "已经使用过该诗句", "已经使用过该诗句", playerInfo)
		}
	}
	count, poetryList, err := h.PoetryService.DescribePoetryInfo(ctx, []string{}, []string{}, []string{req.Poetry}, []string{}, []string{}, []int64{}, 1, 0)
	if err != nil {
		return err
	}
	if count == 0 || len(poetryList) == 0 {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "该诗句不存在", "该诗句不存在", playerInfo)
	}
	poetryStr := req.Poetry
	paraList := []string{}
	err = ejson.Unmarshal([]byte(poetryList[0].Paragraphs), &paraList)
	if err != nil {
		return utils.NewErrAndResponse(s, capi_error.INTERNAL_ERROR_CODE, "诗句解析失败", "诗句解析失败", err)
	}
	for _, graph := range paraList {
		if strings.Contains(graph, poetryStr) {
			poetryStr = graph
		}
	}
	round.UsedPoems = append(round.UsedPoems, &PlayPoemInfo{
		Poem:       poetryStr,
		GamePoetry: poetryList[0],
		Player: &game_room.Player{
			UserId:   playerInfo.UserID,
			Username: playerInfo.Username,
			State:    playerInfo.State,
		},
	})
	players := playerInfo.Room.Players
	for i := 0; i < len(players); i++ {
		player := players[i]
		if player.UserId == playerInfo.UserID {
			for j := 1; j <= len(players)+1; j++ {
				curPlayerInfo := players[(i+j)%len(players)]
				if curPlayerInfo.State == entity.PLAYER_STATE_READY&& curPlayerInfo.UserId != playerInfo.UserID {
					log.DebugEx(ctx, "nextAnsPlayer:", "nextAnsPlayerInfo", curPlayerInfo, "round", round, "curplayerInfo", playerInfo, "roomPlayers", players)
					*round.CurrentPlayer = *curPlayerInfo
					break
				}
			}
			break
		}
	}
	playerInfo.Room.Rounds[len(playerInfo.Room.Rounds)-1] = round
	playerInfo.Room.group.Broadcast("game.submitpoem", playerInfo.Room)
	return nil
}

func (h *GameHandler) PlayerFailed(s *session.Session, msg []byte) error {
	ctx := context.Background()
	// 实现失败逻辑
	playerInfo, err := playerWithSession(s)
	if err != nil {
		return err
	}
	playerInfo, exist := defaultPlayerManager.player(playerInfo.UserID)
	if !exist {
		log.ErrorEx(ctx, "玩家不存在", "玩家不存在", playerInfo.UserID)
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "玩家不存在", "玩家不存在", playerInfo.UserID)
	}
	if playerInfo.Room == nil {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "房间不存在", "用户没有加入任何房间", playerInfo)
	}
	if playerInfo.Room.Status != entity.ROOMSTATE_PLAYING {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "游戏还没有开始", "游戏还没有开始", playerInfo)
	}
	round := playerInfo.Room.Rounds[len(playerInfo.Room.Rounds)-1]
	if round.State != ROUND_STATE_PLAYING {
		return utils.NewErrAndResponse(s, capi_error.INVAILD_PARAM_CODE, "游戏已经结束", "游戏已经结束", playerInfo)
	}
	h.updateGameRoomInfo(ctx, playerInfo, round)
	return nil
}

func (h *GameHandler) updateGameRoomInfo(ctx context.Context, playerInfo *Player, round *Round) {
	existPlayerNum := 0
	winnerPlayerInfo := &game_room.Player{}
	for _, player := range playerInfo.Room.Players {
		if player.State == entity.PLAYER_STATE_READY {
			existPlayerNum++
			winnerPlayerInfo = player
		}
		if player.UserId == playerInfo.UserID {
			player.State = entity.PLAYER_STATE_FAILED
		}
	}
	if existPlayerNum <= 1 {
		if round.State == ROUND_STATE_PLAYING {
			round.Winner = winnerPlayerInfo
			round.State = ROUND_STATE_FINISHED
			//reset player state
			for _, player := range playerInfo.Room.Players {
				player.State = entity.PLAYER_STATE_WAITING
			}

		}
		playerListStr, _ := ejson.Marshal(playerInfo.Room.Players)
		roundStr, _ := ejson.Marshal(playerInfo.Room.Rounds)
		playerInfo.Room.GameRoom.PlayerList = playerListStr
		playerInfo.Room.GameRoom.Rounds = roundStr
		playerInfo.Room.GameRoom.Status = entity.ROOMSTATE_WAITING
		//如果房间没人了，则直接更新房间无人状态
		if existPlayerNum == 0 {
			playerInfo.Room.Status = entity.ROOMSTATE_CLOSED
		}
		playerInfo.Room.group.Broadcast("game.finish", playerInfo.Room)
		async.Run(func() {
			h.GameRoomService.FinishRound(ctx, playerInfo.Room.GameRoom)
		})
	}
}
