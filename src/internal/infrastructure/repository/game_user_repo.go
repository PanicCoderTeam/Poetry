package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/domain/repository"
	"poetry/src/internal/infrastructure"

	"gorm.io/gorm"
)

type GameUserRepository struct {
	db *gorm.DB
}

var _ repository.GameUserRepo = &GameUserRepository{}

func NewGameUserRepository() *GameUserRepository {
	return &GameUserRepository{db: infrastructure.DB}
}

func (r *GameUserRepository) DescribeGameUserInfo(ctx context.Context, userId []*string, username, password []*string, offset, limit *int32) (int32, []*entity.GameUser, error) {
	db := r.db.Table((&entity.GameUser{}).TableName())
	if len(username) > 0 {
		db.Where("username in (?)", username)
	}
	if len(password) > 0 {
		db.Where("password_hash in (?)", password)
	}
	if len(userId) > 0 {
		db.Where("user_id in (?)", userId)
	}
	count := int64(0)
	db.Count(&count)
	if offset != nil {
		db.Offset(int(*offset))
	}
	if limit != nil {
		db.Limit(int(*limit))
	}
	gameUserList := []*entity.GameUser{}
	db.Find(&gameUserList)
	return int32(count), gameUserList, nil
}
func (r *GameUserRepository) CreateGameUserInfo(ctx context.Context, gameUser []*entity.GameUser) error {
	result := r.db.Table((&entity.GameUser{}).TableName()).CreateInBatches(gameUser, len(gameUser))
	return result.Error
}
