package service

import (
	"context"
	"poetry/src/internal/domain/entity"
)

type GameRoomService interface {
	CreateGameRoom(ctx context.Context, maxPlayer int64, gameRoomPassword string) (*entity.GameRoom, error)
	JoinGameRoom(ctx context.Context, gameRoomID string, gameRoomPassword string) (*entity.GameRoom, error)
	DescribeGameRoom(ctx context.Context, gameRoomID, playId, stateList []string, limit, offset *int) (int32, []*entity.GameRoom, error)
	LeaveGameRoom(ctx context.Context, gameRoomID, playId *string) (*entity.GameRoom, error)
	ReadyGameRoom(ctx context.Context, playerId *string) error
	CancelReadyGameRoom(ctx context.Context, playerId *string) error
	StartGame(ctx context.Context, gameRoomID *string) error
	FinishRound(ctx context.Context, gameRoom *entity.GameRoom) error
	DeleteGameRoom(ctx context.Context, gameRoomList []*entity.GameRoom) error
}
