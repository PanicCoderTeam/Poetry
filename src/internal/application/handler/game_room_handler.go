package handler

import (
	"context"
	"encoding/json"
	"poetry/pb/game_room"
	"poetry/src/internal/domain/model"
	"poetry/src/internal/domain/service"
	"poetry/src/pkg/log"
	"poetry/src/pkg/trpc/codec/capi_error"
)

type GameRoomHandler struct {
	gameRoomeService service.GameRoomService
}
type GameRoomHandlerOption struct {
	GameRoomService service.GameRoomService
}

func NewGameRoomHandler(option *GameRoomHandlerOption) *GameRoomHandler {
	return &GameRoomHandler{
		gameRoomeService: option.GameRoomService,
	}
}
func (gameRoomHandler *GameRoomHandler) CreateGameRoom(ctx context.Context, req *game_room.CreateGameRoomRequest) (*game_room.CreateGameRoomResp, error) {
	gameRoom, err := gameRoomHandler.gameRoomeService.CreateGameRoom(ctx, int64(req.MaxPlayers), req.Password)
	if err != nil {
		return nil, err
	}
	return &game_room.CreateGameRoomResp{
		RoomId: gameRoom.RoomID,
	}, nil
}

func (gameRoomHandler *GameRoomHandler) JoinGameRoom(ctx context.Context, req *game_room.JoinGameRoomRequest) (*game_room.JoinGameRoomResp, error) {
	gameRoom, err := gameRoomHandler.gameRoomeService.JoinGameRoom(ctx, req.RoomId, req.Password)
	if err != nil {
		return &game_room.JoinGameRoomResp{}, err
	}
	internalPlayers := []*model.Player{}
	err = json.Unmarshal(gameRoom.PlayerList, &internalPlayers)
	if err != nil {
		return nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", gameRoom)
	}
	players := []*game_room.Player{}
	for _, player := range internalPlayers {
		players = append(players, &game_room.Player{
			UserId:   player.UserID,
			Username: player.UserName,
			State:    player.State,
		})
	}
	return &game_room.JoinGameRoomResp{
		RoomId:  gameRoom.RoomID,
		Players: players,
	}, nil
}

func (GameRoomHandler *GameRoomHandler) DescribeGameRoom(ctx context.Context, req *game_room.DescribeGameRoomRequest) (*game_room.DescribeGameRoomResp, error) {
	roomIdList := []string{}
	playIdList := []string{}
	stateList := []string{}
	for _, filter := range req.Filter {
		if filter.Name == "room-id" {
			roomIdList = filter.Value
		}
		if filter.Name == "player-id" {
			playIdList = filter.Value
		}
		if filter.Name == "state" {
			stateList = filter.Value
		}
	}
	limit := 20
	if req.Limit != 0 {
		limit = int(req.Limit)
	}
	offset := 0
	if req.Offset != 0 {
		offset = int(req.Offset)
	}
	count, gameRoomList, err := GameRoomHandler.gameRoomeService.DescribeGameRoom(ctx, roomIdList, playIdList, stateList, &limit, &offset)
	if err != nil {
		return nil, err
	}
	gameRoomInfoList := []*game_room.GameRoomInfo{}
	for _, gameRoom := range gameRoomList {
		internalPlayers := []*model.Player{}
		err = json.Unmarshal(gameRoom.PlayerList, &internalPlayers)
		if err != nil {
			log.ErrorContextEx(ctx, "序列化失败", err, "playerList", gameRoom.PlayerList)
			return nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", gameRoom)
		}
		players := []*game_room.Player{}
		for _, player := range internalPlayers {
			players = append(players, &game_room.Player{
				UserId:   player.UserID,
				Username: player.UserName,
				State:    player.State,
			})
		}
		gameRoomInfoList = append(gameRoomInfoList, &game_room.GameRoomInfo{
			RoomId:         gameRoom.RoomID,
			Status:         gameRoom.Status,
			MaxPlayers:     int32(gameRoom.MaxPlayers),
			CurrentPlayers: int32(gameRoom.CurrentPlayers),
			GameMode:       gameRoom.GameMode,
			Slogan:         gameRoom.Slogan,
			OwnerId:        gameRoom.OwnerID,
			PlayerList:     players,
		})
	}
	return &game_room.DescribeGameRoomResp{
		TotalCount:   count,
		GameRoomList: gameRoomInfoList,
	}, nil
}

func (GameRoomHandler *GameRoomHandler) LeaveGameRoom(ctx context.Context, req *game_room.LeaveGameRoomRequest) (*game_room.LeaveGameRoomResp, error) {
	// // gameRoom, err := GameRoomHandler.gameRoomeService.LeaveGameRoom(ctx, &req.RoomId, &req.PlayerId)
	// // return &game_room.LeaveGameRoomResp{
	// // 	RoomId: gameRoom.RoomID,
	// }, err
	return nil, nil
}

// func (GameRoomHandler *GameRoomHandler) ReadyGame(ctx context.Context, req *game_room.ReadyGameRequest) (*game_room.ReadyGameResp, error) {

// }

// func (GameRoomHandler *GameRoomHandler) CancelReadyGame(ctx context.Context, req *game_room.CancelReadyGameRequest) (*game_room.CancelReadyGameResp, error) {

// }
