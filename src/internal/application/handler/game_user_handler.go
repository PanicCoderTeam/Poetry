package handler

import (
	"context"
	"poetry/pb/game_user"
	"poetry/src/internal/domain/service"
	"poetry/src/pkg/trpc/codec/capi_error"
)

type GameUserHandler struct {
	GameUsereService service.GameUserService
}
type GameUserHandlerOption struct {
	GameUserService service.GameUserService
}

// NewGameUserHandler 创建并返回一个新的 GameUserHandler 实例
// 参数 option 包含 GameUserHandler 所需的依赖项
func NewGameUserHandler(option *GameUserHandlerOption) *GameUserHandler {
	return &GameUserHandler{
		GameUsereService: option.GameUserService,
	}
}

// CreateGameUser 创建游戏用户
// 参数:
//
//	ctx: 上下文
//	req: 创建用户请求，包含用户名和密码
//
// 返回值:
//
//	*game_user.CreateUserResp: 创建成功的用户信息
//	error: 创建过程中发生的错误
func (gameUserHandler *GameUserHandler) CreateUser(ctx context.Context, req *game_user.CreateUserRequest) (*game_user.CreateUserResp, error) {
	err := gameUserHandler.GameUsereService.CreateGameUserInfo(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &game_user.CreateUserResp{
		Username: req.Username,
	}, nil
}

// Login 处理用户登录请求
// 参数:
//
//	ctx: 上下文
//	req: 包含用户名和密码的登录请求
//
// 返回值:
//
//	*game_user.LoginResp: 登录成功返回用户信息
//	error: 登录失败返回错误信息
//
// 错误码:
//
//	RESOURCE_NOT_FOUND_CODE: 用户不存在
func (gameUserHandler *GameUserHandler) Login(ctx context.Context, req *game_user.LoginRequest) (*game_user.LoginResp, error) {
	_, gameUserList, err := gameUserHandler.GameUsereService.DescribeGameUserInfo(ctx, nil, []*string{&req.Username}, []*string{&req.Password}, nil, nil)
	if err != nil {
		return nil, err
	}
	if len(gameUserList) == 0 {
		return nil, capi_error.NewError(capi_error.RESOURCE_NOT_FOUND_CODE, "用户不存在", nil)
	}
	return &game_user.LoginResp{
		Username: gameUserList[0].Username,
		Token:    gameUserList[0].Token,
		UserId:   gameUserList[0].UserID,
	}, nil
}
