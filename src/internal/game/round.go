package game

import (
	"poetry/pb/game_room"
	"poetry/pb/poetry"
)

type PlayPoemInfo struct {
	Poem       string
	GamePoetry *poetry.PoetryInfo
	Player     *game_room.Player
}

const (
	ROUND_STATE_PLAYING  = "PLAYING"
	ROUND_STATE_FINISHED = "FINISHED"
)

type Round struct {
	CurrentWord   string
	CurrentPlayer *game_room.Player
	UsedPoems     []*PlayPoemInfo
	Winner        *game_room.Player
	State         string
}
