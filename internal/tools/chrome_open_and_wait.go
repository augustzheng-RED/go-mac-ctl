package tools

import (
	"time"

	"go-mac-ctl/internal/observe"
)

func ChromeOpenAndWait(url string, timeout time.Duration) error {
	if err := OpenAppURL("Google Chrome", url); err != nil {
		return err
	}

	return observe.WaitForChromeReady(url, timeout)
}
