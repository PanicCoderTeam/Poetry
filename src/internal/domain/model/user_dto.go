package model

import "poetry/src/internal/domain/entity"

// UserWithToken 带Token的用户信息
type UserWithToken struct {
	*entity.GameUser
	Token string `json:"token"`
}
