package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/domain/repository"
	"poetry/src/internal/infrastructure"
	"strings"

	"gorm.io/gorm"
)

type PoetryRepository struct {
	db *gorm.DB
}

var _ repository.PeotryRepo = &PoetryRepository{}

func NewPoetryRepository() *PoetryRepository {
	return &PoetryRepository{db: infrastructure.DB}
}

func (r *PoetryRepository) DescribePeotryInfo(ctx context.Context, title, author, paragraphs, dynasty, poetry_type []string, tagId []int64, limit, offset int) (int64, []*entity.Poetry, error) {

	db := r.db.Table((&entity.Poetry{}).TableName())
	if len(tagId) > 0 {
		db.Joins("left join poetry_tag on poetry_tag.poetry_id = poetry.id").Where("poetry_tag.tag_id in (?)", tagId)
	}
	if len(title) > 0 {
		conditions := []string{}
		for _, t := range title {
			conditions = append(conditions, "title like '%"+t+"%'")
		}
		db.Where(strings.Join(conditions, " OR "))
	}
	if len(author) > 0 {
		db.Where("author in (?)", author)
	}
	if len(paragraphs) > 0 {
		db.Where("paragraphs like ? or paragraphs_tradition like ?", "%"+paragraphs[0]+"%", "%"+paragraphs[0]+"%")
	}
	if len(dynasty) > 0 {
		db.Where("dynasty in (?) ", dynasty)
	}
	if len(poetry_type) > 0 {
		db.Where("poetry_type in (?)", poetry_type)
	}
	count := int64(0)
	db.Count(&count)
	db.Limit(limit)
	db.Offset(offset)
	poetryList := []*entity.Poetry{}
	if count > 0 {
		db.Find(&poetryList)
	}
	// log.DebugContextEx(ctx, "poetryList", poetryList, "count", count)
	return count, poetryList, nil
}

func (r *PoetryRepository) UpdatePoetryInfo(ctx context.Context, poetryList []*entity.Poetry) error {
	r.db.Table(((&entity.Poetry{}).TableName())).Save(poetryList)
	return nil
}

func (r *PoetryRepository) CreatePoetryInfo(ctx context.Context, poetryList []*entity.Poetry) error {
	result := r.db.Table((&entity.Poetry{}).TableName()).CreateInBatches(poetryList, len(poetryList))
	return result.Error
}
