package service

import (
	"context"
	"poetry/pb/poetry"
	"poetry/src/internal/domain/entity"
)

type PoetryService interface {
	DescribePoetryInfo(ctx context.Context, title, author, paragraphs, dynasty, poetryType []string, tagId []int64, limit, offset int) (int64, []*poetry.PoetryInfo, error)
	CreatePoetryInfo(ctx context.Context, poetry []*entity.Poetry) error
}
