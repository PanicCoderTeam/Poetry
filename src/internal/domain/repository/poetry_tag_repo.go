package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
)

type PoetryTagRepo interface {
	DeletePoetryTagInfo(ctx context.Context, tagIdList []int64) error
	CreatePoetryTagInfo(ctx context.Context, poetryList []*entity.PoetryTag) error
}
