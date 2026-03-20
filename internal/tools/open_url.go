package tools

import "go-mac-ctl/internal/executor"

func OpenURL(url string) error {
	return executor.Run("open", url)
}
