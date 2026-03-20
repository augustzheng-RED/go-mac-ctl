package act

import (
	"time"

	"go-mac-ctl/internal/executor"
)

func MoveMouse(x, y int) error {
	return executor.MoveMouseSmooth(x, y, 40, 2*time.Millisecond)
}

func MouseLocation() (int, int, error) {
	return executor.MouseLocation()
}

func LeftClick(x, y int) error {
	if err := MoveMouse(x, y); err != nil {
		return err
	}

	return executor.LeftClick(x, y)
}
