package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
)

type TagRepo interface {
	DescribeTag(ctx context.Context, name, category []string, parentTag []int64, limit, offset int) (int64, []*entity.Tag, error)
	DeleteTagInfo(ctx context.Context, tagIdList []int64) error
	CreateTag(ctx context.Context, poetryList []*entity.Tag) error
}
