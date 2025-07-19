package serviceimpl

import (
	"context"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/domain/model"
	"poetry/src/internal/domain/repository"
	"poetry/src/internal/domain/service"
	"poetry/src/pkg/basic"
	"poetry/src/pkg/log"
	"poetry/src/pkg/trpc/auth"
	"poetry/src/pkg/trpc/codec/capi_error"
	"time"

	"github.com/google/uuid"
)

type GameUserServiceImpl struct {
	GameUserRepo repository.GameUserRepo
}

type GameUserServiceImplOption struct {
	GameUserRepo repository.GameUserRepo
}

var _ service.GameUserService = &GameUserServiceImpl{}

func NewGameUserServiceImpl(option GameUserServiceImplOption) *GameUserServiceImpl {
	return &GameUserServiceImpl{
		GameUserRepo: option.GameUserRepo,
	}
}

func (psi *GameUserServiceImpl) DescribeGameUserInfo(ctx context.Context, userId []*string, username, password []*string, offset, limit *int32) (int32, []*model.UserWithToken, error) {
	if len(userId) == 0 && (len(username) == 0 || len(password) == 0) {
		return -1, nil, nil
	}
	for i := range password {
		*password[i] = basic.HashPassword(*(password[i]))
	}
	count, GameUserDaoList, err := psi.GameUserRepo.DescribeGameUserInfo(ctx, userId, username, password, offset, limit)
	if err != nil {
		return -1, nil, err
	}
	if len(GameUserDaoList) == 0 {
		return 0, []*model.UserWithToken{}, nil
	}
	// 转换为UserWithToken并添加token
	usersWithToken := []*model.UserWithToken{}
	for _, GameUserDao := range GameUserDaoList {
		user := &model.UserWithToken{GameUser: GameUserDao}
		token, err := auth.GenerateToken(user.UserID, user.Username)
		if err != nil {
			return -1, nil, capi_error.NewErr(capi_error.INTERNAL_ERROR_CODE, "生成Token失败")
		}
		user.Token = token
		usersWithToken = append(usersWithToken, user)
	}
	log.DebugContextEx(ctx, "userWithToken", usersWithToken)
	return count, usersWithToken, nil
}

func (psi *GameUserServiceImpl) CreateGameUserInfo(ctx context.Context, username, password string) error {
	if len(username) == 0 || len(password) == 0 {
		return capi_error.NewErr(capi_error.INVAILD_PARAM_CODE, "用户名或密码不能为空")
	}
	t2, _ := time.ParseInLocation("2006-01-02 15:04:05", "1969-07-20 20:17:40", time.UTC)

	gameUserList := []*entity.GameUser{
		{
			UserID:        "user-" + uuid.New().String()[:12],
			Username:      username,
			PasswordHash:  basic.HashPassword(password),
			ThirdPartyID:  "",
			GameRounds:    0,
			Score:         0,    // 当前积分
			Coins:         0,    // 游戏币（购买道具用）
			GameRank:      "青铜", // 段位等级
			CurrentRoomID: "",
			HighestStreak: 0, // 最高连击次数
			Achievements:  "[]",
			FriendList:    "[]", // 好友ID列表
			CreatedAt:     time.Now(),
			LastLogin:     t2,
			AccountStatus: "正常",
			ThemeStats:    "[]",
			Items:         "[]",
			GameHistory:   "[]",
		},
	}
	return psi.GameUserRepo.CreateGameUserInfo(ctx, gameUserList)
}
