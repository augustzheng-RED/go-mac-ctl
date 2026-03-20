package observe

import "go-mac-ctl/internal/executor"

func ScreenSize() (int, int, error) {
	return executor.ScreenSize()
}
