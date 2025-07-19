package game

import (
	"poetry/src/pkg/trpc/codec/capi_error"
	"runtime"
	"strings"

	"github.com/lonng/nano/session"
)

const (
	ModeTrios      = 3 // 三人模式
	ModeFours      = 4 // 四人模式
	KEY_CUR_PLAYER = "cur_player"
)

// func verifyOptions(opts *protocol.DeskOptions) bool {
// 	if opts == nil {
// 		return false
// 	}

// 	if opts.Mode != ModeTrios && opts.Mode != 4 {
// 		return false
// 	}

// 	if opts.MaxRound != 1 && opts.MaxRound != 4 && opts.MaxRound != 8 && opts.MaxRound != 16 {
// 		return false
// 	}

// 	return true
// }

func playerWithSession(s *session.Session) (*Player, error) {
	p, ok := s.Value(KEY_CUR_PLAYER).(*Player)
	if !ok {
		return nil, capi_error.NewErr(capi_error.INVAILD_PARAM_CODE, "invalid player")
	}
	return p, nil
}

func stack() string {
	buf := make([]byte, 10000)
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	s := string(buf)

	// skip nano frames lines
	const skip = 7
	count := 0
	index := strings.IndexFunc(s, func(c rune) bool {
		if c != '\n' {
			return false
		}
		count++
		return count == skip
	})
	return s[index+1:]
}
