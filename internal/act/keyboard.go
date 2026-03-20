package act

import "go-mac-ctl/internal/executor"

func TypeText(text string) error {
	return executor.TypeText(text)
}

func PressKey(key string) error {
	return executor.PressKey(key)
}
