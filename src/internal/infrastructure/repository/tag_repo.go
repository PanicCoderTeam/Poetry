package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/domain/repository"
	"poetry/src/internal/infrastructure"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

var _ repository.TagRepo = &TagRepository{}

func NewTagRepository() *TagRepository {
	return &TagRepository{db: infrastructure.DB}
}

func (r *TagRepository) DescribeTag(ctx context.Context, name, category []string, parentTagId []int64, limit, offset int) (int64, []*entity.Tag, error) {

	db := r.db.Table((&entity.Tag{}).TableName())
	if len(name) > 0 {
		db.Where("name = ?", name)
	}
	if len(category) == 1 && category[0] == "main" {
		db.Where("level = 1")
	}
	if len(parentTagId) > 0 {
		db.Where("parent_tag_id in (?)", parentTagId)
	}
	count := int64(0)
	db.Count(&count)
	db.Limit(limit)
	db.Offset(offset)
	tagList := []*entity.Tag{}
	if count > 0 {
		db.Find(&tagList)
	}
	return count, tagList, nil
}
func (r *TagRepository) DeleteTagInfo(ctx context.Context, tagIdList []int64) error {
	r.db.Table(((&entity.Tag{}).TableName())).Delete(&entity.Tag{}, "id in (?)", tagIdList)
	return nil
}

func (r *TagRepository) CreateTag(ctx context.Context, tagList []*entity.Tag) error {
	result := r.db.Table((&entity.Tag{}).TableName()).CreateInBatches(tagList, len(tagList))
	return result.Error
}
