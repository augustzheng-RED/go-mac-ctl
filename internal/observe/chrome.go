package observe

import "go-mac-ctl/internal/executor"

func ChromeCurrentURL() (string, error) {
	return executor.Output(
		"osascript",
		"-e",
		`tell application "Google Chrome"
if (count of windows) is 0 then
	return ""
end if
return URL of active tab of front window
end tell`,
	)
}
