package repository

import (
	"context"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/domain/repository"
	"poetry/src/internal/infrastructure"

	"gorm.io/gorm"
)

type GameRoomRepository struct {
	db *gorm.DB
}

var _ repository.GameRoomRepo = &GameRoomRepository{}

func NewGameRoomRepository() *GameRoomRepository {

	return &GameRoomRepository{db: infrastructure.DB}
}

func (r *GameRoomRepository) CreateGameRoomInfo(ctx context.Context, GameRoomList []*entity.GameRoom) error {
	result := r.db.Table((&entity.GameRoom{}).TableName()).CreateInBatches(GameRoomList, len(GameRoomList))
	return result.Error
}
func (r *GameRoomRepository) DescribeGameRoomInfo(ctx context.Context, gameRoomId, playerId, gameRoomState []string, limit, offset *int) (int32, []*entity.GameRoom, error) {
	gameRoom := []*entity.GameRoom{}
	tx := r.db.Table((&entity.GameRoom{}).TableName())
	if len(gameRoomId) > 0 {
		tx = tx.Where("room_id in (?)", gameRoomId)
	}
	if len(playerId) > 0 {
		tx = tx.Where("player_list->'$[*].player_id' in (?)", playerId)
	}
	if len(gameRoomState) > 0 {
		tx = tx.Where("status in (?)", gameRoomState)
	}
	count := int64(0)
	tx.Count(&count)

	if limit != nil {
		tx.Limit(*limit)
	}
	if offset != nil {
		tx.Offset(*offset)
	}
	tx = tx.Find(&gameRoom)
	if tx.Error != nil {
		return 0, nil, tx.Error
	}
	return int32(count), gameRoom, nil
}
func (r *GameRoomRepository) ModifyGameRoomInfo(ctx context.Context, gameRoom []*entity.GameRoom) error {
	if len(gameRoom) == 0 {
		return nil
	}
	tx := r.db.Table((&entity.GameRoom{}).TableName()).Save(gameRoom)
	return tx.Error
}
