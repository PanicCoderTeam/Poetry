package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
)

type GameUserRepo interface {
	DescribeGameUserInfo(ctx context.Context, userId []*string, username, password []*string, offset, limit *int32) (int32, []*entity.GameUser, error)
	CreateGameUserInfo(ctx context.Context, gameUser []*entity.GameUser) error
}
