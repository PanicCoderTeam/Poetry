package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/domain/repository"
	"poetry/src/internal/infrastructure"

	"gorm.io/gorm"
)

type PoetryTagRepository struct {
	db *gorm.DB
}

var _ repository.PoetryTagRepo = &PoetryTagRepository{}

func NewPoetryTagRepository() *PoetryTagRepository {
	return &PoetryTagRepository{db: infrastructure.DB}
}

func (r *PoetryTagRepository) DeletePoetryTagInfo(ctx context.Context, tagIdList []int64) error {
	r.db.Table(((&entity.Tag{}).TableName())).Delete(&entity.PoetryTag{}, "id in (?)", tagIdList)
	return nil
}

func (r *PoetryTagRepository) CreatePoetryTagInfo(ctx context.Context, tagList []*entity.PoetryTag) error {
	result := r.db.Table((&entity.PoetryTag{}).TableName()).CreateInBatches(tagList, len(tagList))
	return result.Error
}
