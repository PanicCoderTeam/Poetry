package service

import (
	"context"
	"poetry/src/internal/domain/model"
)

type GameUserService interface {
	DescribeGameUserInfo(ctx context.Context, userId []*string, username, password []*string, offset, limit *int32) (int32, []*model.UserWithToken, error)
	CreateGameUserInfo(ctx context.Context, username, password string) error
}
