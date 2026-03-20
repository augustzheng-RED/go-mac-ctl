package observe

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

func WaitForChromeReady(target string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		currentURL, err := ChromeCurrentURL()
		if err == nil && chromeMatchesTarget(currentURL, target) {
			return nil
		}

		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("timed out waiting for Google Chrome to reach %q", target)
}

func chromeMatchesTarget(currentURL, target string) bool {
	if currentURL == "" {
		return false
	}

	currentParsed, currentErr := url.Parse(currentURL)
	targetParsed, targetErr := url.Parse(target)
	if currentErr == nil && targetErr == nil {
		if hostsMatch(currentParsed.Host, targetParsed.Host) {
			return true
		}
	}

	return strings.Contains(currentURL, target)
}

func hostsMatch(currentHost, targetHost string) bool {
	currentHost = normalizeHost(currentHost)
	targetHost = normalizeHost(targetHost)

	return currentHost != "" && targetHost != "" && currentHost == targetHost
}

func normalizeHost(host string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	host = strings.TrimPrefix(host, "www.")
	return host
}
