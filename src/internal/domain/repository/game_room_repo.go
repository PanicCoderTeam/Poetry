package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
)

type GameRoomRepo interface {
	CreateGameRoomInfo(ctx context.Context, gameRoomList []*entity.GameRoom) error
	DescribeGameRoomInfo(ctx context.Context, gameRoomID, playerId, state []string, limit, offset *int) (int32, []*entity.GameRoom, error)
	ModifyGameRoomInfo(ctx context.Context, gameRoom []*entity.GameRoom) error
}
