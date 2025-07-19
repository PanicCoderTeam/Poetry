package serviceimpl

import (
	"context"
	"encoding/json"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/domain/model"
	"poetry/src/internal/domain/repository"
	"poetry/src/internal/domain/service"
	"poetry/src/pkg/basic"
	"poetry/src/pkg/log"
	"poetry/src/pkg/trpc/codec/capi_error"
	"poetry/src/pkg/trpc/filter"
	"time"

	"github.com/google/uuid"
)

type GameRoomServiceImpl struct {
	GameRoomRepo repository.GameRoomRepo
}

type GameRoomServiceImplOption struct {
	GameRoomRepo repository.GameRoomRepo
}

var _ service.GameRoomService = &GameRoomServiceImpl{}

func NewGameRoomServiceImpl(option GameRoomServiceImplOption) *GameRoomServiceImpl {
	return &GameRoomServiceImpl{
		GameRoomRepo: option.GameRoomRepo,
	}
}

func (psi *GameRoomServiceImpl) CreateGameRoom(ctx context.Context, maxPlayer int64, gameRoomPassword string) (*entity.GameRoom, error) {
	gameRoomList := []*entity.GameRoom{}
	var ml = int32(1)
	slogan := ""
	playerList := []model.Player{}
	authUserInfo := ctx.Value(basic.ClaimsKeyVal).(filter.AuthUserInfo)
	playerList = append(playerList, model.Player{
		UserID:   authUserInfo.UserId,
		UserName: authUserInfo.UserName,
		State:    "READY",
	})
	playListStr, err := json.Marshal(playerList)
	if err != nil {
		return nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", err)
	}
	authUserInfo, ok := ctx.Value(basic.ClaimsKeyVal).(filter.AuthUserInfo)
	if !ok {
		return nil, capi_error.NewErr(capi_error.REQUEST_NOT_AUTH_CODE, "获取用户信息失败", err)
	}
	gameRoomList = append(gameRoomList, &entity.GameRoom{
		RoomID:         "room-" + uuid.New().String()[:12],
		Status:         "WAITING",
		MaxPlayers:     int32(maxPlayer),
		CurrentPlayers: 1,
		Password:       gameRoomPassword,
		MinLevel:       ml,
		CostItemID:     ml,
		CostAmount:     ml,
		CreatedAt:      time.Now(),
		GameMode:       "CLASSIC",
		Slogan:         slogan,
		OwnerID:        authUserInfo.UserId,
		GatewayIP:      "127.0.0.1", //默认是当前机器，后续需要通过逻辑改成集群gatewayIp，用于客户端建联
		PlayerList:     playListStr,
	})
	err = psi.GameRoomRepo.CreateGameRoomInfo(ctx, gameRoomList)
	if err != nil {
		return nil, err
	}
	return gameRoomList[0], nil
}

func (psi *GameRoomServiceImpl) JoinGameRoom(ctx context.Context, gameRoomId string, gameRoomPassword string) (*entity.GameRoom, error) {
	authUserInfo := ctx.Value(basic.ClaimsKeyVal).(filter.AuthUserInfo)
	stateList := []string{entity.ROOMSTATE_WAITING}
	if len(authUserInfo.UserId) == 0 {
		return nil, capi_error.NewError(capi_error.REQUEST_NOT_AUTH_CODE, "user is not login", nil)
	}
	_, gameRoomList, err := psi.GameRoomRepo.DescribeGameRoomInfo(ctx, nil, []string{authUserInfo.UserId}, stateList, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(gameRoomList) > 0 {
		return nil, capi_error.NewErr(capi_error.USER_ALREADY_IN_GAME_CODE, "用户已加入房间", nil)
	}
	//fixme 查询用户等级，进行用户相关配置校验
	_, gameRoomList, err = psi.GameRoomRepo.DescribeGameRoomInfo(ctx, []string{gameRoomId}, nil, stateList, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(gameRoomList) == 0 {
		return nil, capi_error.NewErr(capi_error.RESOURCE_NOT_FOUND_CODE, "房间不存在", nil)
	}
	gameRoom := gameRoomList[0]
	if gameRoom.Password != gameRoomPassword {
		return nil, capi_error.NewError(capi_error.REQUEST_NOT_AUTH_CODE, "password is wrong", nil)
	}
	log.DebugContextEx(ctx, "gameRoom", gameRoom, "gameRoom.PlayerList", gameRoom.PlayerList)
	playerList := []*model.Player{}
	err = json.Unmarshal(gameRoom.PlayerList, &playerList)
	if err != nil {
		log.ErrorContextEx(ctx, "序列化失败", err)
		return nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", gameRoom)
	}
	for _, player := range playerList {
		log.DebugContextEx(ctx, "ctxValue", ctx.Value(basic.ClaimsKeyVal))
		if player.UserID == authUserInfo.UserId {
			return gameRoom, nil
		}
	}
	if gameRoom.CurrentPlayers >= gameRoom.MaxPlayers {
		return nil, capi_error.NewError(capi_error.RESOURCE_OUT_OF_LIMIT, "房间已满", nil)
	}
	playerList = append(playerList, &model.Player{
		UserID:   authUserInfo.UserId,
		UserName: authUserInfo.UserName,
		State:    "READY",
	})
	playListStr, err := json.Marshal(playerList)
	if err != nil {
		return nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", err)
	}
	gameRoom.CurrentPlayers = int32(len(playerList))
	gameRoom.PlayerList = playListStr
	err = psi.GameRoomRepo.ModifyGameRoomInfo(ctx, []*entity.GameRoom{gameRoom})
	if err != nil {
		return nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "modify game room info failed", err)
	}
	return gameRoom, nil

}

func (psi *GameRoomServiceImpl) DescribeGameRoom(ctx context.Context, gameRoomId, playerId, stateList []string, limit, offset *int) (int32, []*entity.GameRoom, error) {
	count, gameRoom, err := psi.GameRoomRepo.DescribeGameRoomInfo(ctx, gameRoomId, playerId, stateList, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	return count, gameRoom, nil
}

func (psi *GameRoomServiceImpl) LeaveGameRoom(ctx context.Context, gameRoomId, playerId *string) (*entity.GameRoom, error) {
	if gameRoomId == nil || playerId == nil {
		return nil, capi_error.NewError(capi_error.REQUEST_NOT_AUTH_CODE, "gameRoomId or playerId is nil", nil)
	}
	stateList := []string{entity.ROOMSTATE_WAITING, entity.ROOMSTATE_PLAYING}

	playerList := []*model.Player{}
	_, gameRoomList, err := psi.GameRoomRepo.DescribeGameRoomInfo(ctx, []string{*gameRoomId}, nil, stateList, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(gameRoomList) == 0 {
		return nil, capi_error.NewErr(capi_error.RESOURCE_NOT_FOUND_CODE, "房间不存在", nil)
	}
	gameRoom := gameRoomList[0]
	err = json.Unmarshal(gameRoom.PlayerList, &playerList)
	if err != nil {
		log.ErrorContextEx(ctx, "序列化失败", err)
		return nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", gameRoom)
	}
	for index, player := range playerList {
		if player.UserID == *playerId {
			playerList = append(playerList[:index], playerList[index+1:]...)
			break
		}
	}
	playListStr, err := json.Marshal(playerList)
	if err != nil {
		return nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", err)
	}
	gameRoom.CurrentPlayers = int32(len(playerList))
	gameRoom.PlayerList = playListStr
	gameRoom.Status = entity.ROOMSTATE_WAITING
	err = psi.GameRoomRepo.ModifyGameRoomInfo(ctx, []*entity.GameRoom{gameRoom})
	if err != nil {
		return nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "modify game room info failed", err)
	}
	return gameRoom, nil
}

func (gameRoomService *GameRoomServiceImpl) ReadyGameRoom(ctx context.Context, playerId *string) error {
	_, gameRoomList, err := gameRoomService.GameRoomRepo.DescribeGameRoomInfo(ctx, nil, []string{*playerId}, nil, nil, nil)
	if err != nil {
		return err
	}
	if len(gameRoomList) == 0 {
		return capi_error.NewErr(capi_error.RESOURCE_NOT_FOUND_CODE, "该选手不在任何房间中", nil)
	}
	gameRoom := gameRoomList[0]
	playerList := []*model.Player{}
	err = json.Unmarshal(gameRoom.PlayerList, &playerList)
	if err != nil {
		log.ErrorContextEx(ctx, "序列化失败", err)
		return capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", gameRoom)
	}
	for index, player := range playerList {
		if player.UserID == *playerId {
			playerList[index].State = "READY"
			playListStr, err := json.Marshal(playerList)
			if err != nil {
				return capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", err)
			}
			gameRoom.PlayerList = playListStr
			break
		}
	}
	err = gameRoomService.GameRoomRepo.ModifyGameRoomInfo(ctx, []*entity.GameRoom{gameRoom})
	return err
}
func (gameRoomService *GameRoomServiceImpl) CancelReadyGameRoom(ctx context.Context, playerId *string) error {

	_, gameRoomList, err := gameRoomService.GameRoomRepo.DescribeGameRoomInfo(ctx, nil, []string{*playerId}, nil, nil, nil)
	if err != nil {
		return err
	}
	if len(gameRoomList) == 0 {
		return capi_error.NewErr(capi_error.RESOURCE_NOT_FOUND_CODE, "该选手不在任何房间中", nil)
	}
	gameRoom := gameRoomList[0]
	playerList := []*model.Player{}
	err = json.Unmarshal(gameRoom.PlayerList, &playerList)
	if err != nil {
		log.ErrorContextEx(ctx, "序列化失败", err)
		return capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", gameRoom)
	}
	for index, player := range playerList {
		if player.UserID == *playerId {
			playerList[index].State = "NOTREADY"
			playListStr, err := json.Marshal(playerList)
			if err != nil {
				return capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "序列化失败", err)
			}
			gameRoom.PlayerList = playListStr
			break
		}
	}
	err = gameRoomService.GameRoomRepo.ModifyGameRoomInfo(ctx, []*entity.GameRoom{gameRoom})
	return err
}

func (gameRoomService *GameRoomServiceImpl) StartGame(ctx context.Context, gameRoomID *string) error {
	authUserInfo := ctx.Value(basic.ClaimsKeyVal).(filter.AuthUserInfo)
	if len(authUserInfo.UserId) == 0 {
		return capi_error.NewError(capi_error.REQUEST_NOT_AUTH_CODE, "user is not login", nil)
	}
	_, gameRoomList, err := gameRoomService.GameRoomRepo.DescribeGameRoomInfo(ctx, []string{*gameRoomID}, nil, nil, nil, nil)
	if err != nil {
		return err
	}
	if len(gameRoomList) == 0 {
		return capi_error.NewErr(capi_error.RESOURCE_NOT_FOUND_CODE, "房间不存在", nil)
	}
	gameRoom := gameRoomList[0]
	if gameRoom.OwnerID != authUserInfo.UserId {
		return capi_error.NewErr(capi_error.REQUEST_NOT_AUTH_CODE, "用户不是房主", nil)
	}
	gameRoom.Status = entity.ROOMSTATE_PLAYING
	return gameRoomService.GameRoomRepo.ModifyGameRoomInfo(ctx, []*entity.GameRoom{gameRoom})
}

func (gameRoomService *GameRoomServiceImpl) FinishRound(ctx context.Context, gameRoom *entity.GameRoom) error {
	return gameRoomService.GameRoomRepo.ModifyGameRoomInfo(ctx, []*entity.GameRoom{gameRoom})
}

func (gameRoomService *GameRoomServiceImpl) DeleteGameRoom(ctx context.Context, gameRoomList []*entity.GameRoom) error {
	for _, gameRoom := range gameRoomList {
		gameRoom.Status = entity.ROOMSTATE_CLOSED
	}
	return gameRoomService.GameRoomRepo.ModifyGameRoomInfo(ctx, gameRoomList)
}
