package async

import "poetry/src/pkg/log"

func pcall(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("aync/pcall: Error=%v", err)
		}
	}()

	fn()
}

func Run(fn func()) {
	go pcall(fn)
}
