package serviceimpl

import (
	"context"
	"poetry/pb/tag"
	"poetry/src/internal/domain/repository"
	"poetry/src/internal/domain/service"
)

type TagServiceImpl struct {
	tagRepo repository.TagRepo
}

type TagServiceImplOption struct {
	TagRepo repository.TagRepo
}

var _ service.TagService = &TagServiceImpl{}

func NewTagServiceImpl(option TagServiceImplOption) *TagServiceImpl {
	return &TagServiceImpl{
		tagRepo: option.TagRepo,
	}
}

func (psi *TagServiceImpl) DescribeTagInfo(ctx context.Context, name, category []string, level []int64, limit, offset int) (int64, []*tag.TagInfo, error) {
	count, tagDaoList, err := psi.tagRepo.DescribeTag(ctx, name, category, level, limit, offset)
	tagList := []*tag.TagInfo{}
	if err != nil {
		return count, tagList, err
	}
	for _, tagInfo := range tagDaoList {
		tagList = append(tagList, &tag.TagInfo{
			Id:       tagInfo.ID,
			Name:     tagInfo.Name,
			Category: tagInfo.Category,
			Level:    int64(tagInfo.Level),
		})
	}
	return count, tagList, nil
}
