package service

import (
	"context"
	"poetry/pb/tag"
)

type TagService interface {
	DescribeTagInfo(ctx context.Context, name, category []string, parentTagId []int64, limit, offset int) (int64, []*tag.TagInfo, error)
}
