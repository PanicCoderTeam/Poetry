package main

import (
	"fmt"
	"poetry/pb/game_room"
	"poetry/pb/game_user"
	"poetry/pb/poetry"
	"poetry/pb/tag"
	"poetry/src/config"
	"poetry/src/internal/application/handler"
	serviceimpl "poetry/src/internal/domain/service/service_impl"
	"poetry/src/internal/game"
	"poetry/src/internal/infrastructure"
	"poetry/src/internal/infrastructure/repository"
	_ "poetry/src/pkg/trpc/codec/capi"
	_ "poetry/src/pkg/trpc/filter"

	_ "github.com/Andrew-M-C/trpc-go-utils/plugin"
	trpc "trpc.group/trpc-go/trpc-go"
	_ "trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/log"
)

func main() {
	// 初始化nano应用

	// 初始化TRPC服务
	s := trpc.NewServer()
	fmt.Printf("db Cofnig%+v\n", config.DBConfig)
	infrastructure.InitDB()
	phOption := &handler.PoetryHandlerOption{
		PeotryService: serviceimpl.NewPoetryServiceImpl(serviceimpl.PoetryServiceImplOption{
			PoetryRepo: repository.NewPoetryRepository(),
		}),
	}
	grOption := &handler.GameRoomHandlerOption{
		GameRoomService: serviceimpl.NewGameRoomServiceImpl(serviceimpl.GameRoomServiceImplOption{
			GameRoomRepo: repository.NewGameRoomRepository(),
		}),
	}
	guOption := &handler.GameUserHandlerOption{
		GameUserService: serviceimpl.NewGameUserServiceImpl(serviceimpl.GameUserServiceImplOption{
			GameUserRepo: repository.NewGameUserRepository(),
		}),
	}
	tagOption := &handler.TagHandlerOption{
		TagService: serviceimpl.NewTagServiceImpl(serviceimpl.TagServiceImplOption{
			TagRepo: repository.NewTagRepository(),
		}),
	}
	poetry.RegisterPoetryService(s.Service("trpc.poetry.http.Poetry"), handler.NewPoetryHandler(phOption))
	poetry.RegisterPoetryService(s.Service("trpc.poetry.trpc.Poetry"), handler.NewPoetryHandler(phOption))
	game_room.RegisterGameRoomService(s.Service("trpc.poetry.http.Poetry"), handler.NewGameRoomHandler(grOption))
	game_room.RegisterGameRoomService(s.Service("trpc.poetry.trpc.Poetry"), handler.NewGameRoomHandler(grOption))
	game_user.RegisterUserService(s.Service("trpc.poetry.http.Poetry"), handler.NewGameUserHandler(guOption))
	game_user.RegisterUserService(s.Service("trpc.poetry.trpc.Poetry"), handler.NewGameUserHandler(guOption))
	tag.RegisterTagService(s.Service("trpc.poetry.http.Poetry"), handler.NewTagHandler(tagOption))
	tag.RegisterTagService(s.Service("trpc.poetry.trpc.Poetry"), handler.NewTagHandler(tagOption))
	fmt.Printf("%+v\n", config.DBConfig)
	// 启动服务
	go func() {
		game.Startup()
	}()

	// 启动TRPC服务
	if err := s.Serve(); err != nil {
		log.Error("trpc服务启动失败: ", err)
	}
}
