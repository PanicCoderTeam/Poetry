package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
)

type PeotryRepo interface {
	DescribePeotryInfo(ctx context.Context, title, author, paragraphs, dynasty, poetry_type []string, tagId []int64, limit, offset int) (int64, []*entity.Poetry, error)
	CreatePoetryInfo(ctx context.Context, poetryList []*entity.Poetry) error
}
