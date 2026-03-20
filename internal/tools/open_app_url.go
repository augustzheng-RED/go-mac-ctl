package tools

import "go-mac-ctl/internal/executor"

func OpenAppURL(appName, url string) error {
	return executor.Run("open", "-a", appName, url)
}
