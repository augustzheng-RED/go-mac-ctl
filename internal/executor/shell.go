package executor

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command %s: %w", name, err)
	}
	return nil
}

func Output(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message != "" {
			return "", fmt.Errorf("failed to run command %s: %w: %s", name, err, message)
		}
		return "", fmt.Errorf("failed to run command %s: %w", name, err)
	}

	return strings.TrimSpace(string(output)), nil
}

func Screenshot(path string) error {
	return Run("screencapture", "-x", path)
}

func MoveMouse(x, y int) error {
	script := strings.Join([]string{
		"import CoreGraphics",
		"CGWarpMouseCursorPosition(CGPoint(x: " + strconv.Itoa(x) + ", y: " + strconv.Itoa(y) + "))",
	}, "\n")

	return Run("swift", "-e", script)
}

func MouseLocation() (int, int, error) {
	script := strings.Join([]string{
		"import CoreGraphics",
		"let point = CGEvent(source: nil)!.location",
		`print("\(Int(point.x)),\(Int(point.y))")`,
	}, "\n")

	output, err := Output("swift", "-e", script)
	if err != nil {
		return 0, 0, err
	}

	parts := strings.Split(output, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected mouse location output: %q", output)
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid mouse x coordinate: %w", err)
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid mouse y coordinate: %w", err)
	}

	return x, y, nil
}

func MoveMouseSmooth(x, y int, steps int, delay time.Duration) error {
	startX, startY, err := MouseLocation()
	if err != nil {
		return err
	}

	if steps < 1 {
		steps = 1
	}

	for i := 1; i <= steps; i++ {
		progress := float64(i) / float64(steps)
		nextX := startX + int(math.Round(float64(x-startX)*progress))
		nextY := startY + int(math.Round(float64(y-startY)*progress))

		if err := MoveMouse(nextX, nextY); err != nil {
			return err
		}

		time.Sleep(delay)
	}

	return nil
}

func LeftClick(x, y int) error {
	script := strings.Join([]string{
		"import CoreGraphics",
		"let point = CGPoint(x: " + strconv.Itoa(x) + ", y: " + strconv.Itoa(y) + ")",
		"let mouseDown = CGEvent(mouseEventSource: nil, mouseType: .leftMouseDown, mouseCursorPosition: point, mouseButton: .left)",
		"let mouseUp = CGEvent(mouseEventSource: nil, mouseType: .leftMouseUp, mouseCursorPosition: point, mouseButton: .left)",
		"mouseDown?.post(tap: .cghidEventTap)",
		"mouseUp?.post(tap: .cghidEventTap)",
	}, "\n")

	return Run("swift", "-e", script)
}

func TypeText(text string) error {
	return Run(
		"osascript",
		"-e",
		`tell application "System Events" to keystroke `+strconv.Quote(text),
	)
}

func PressKey(key string) error {
	keyCode, ok := map[string]int{
		"enter":  36,
		"tab":    48,
		"escape": 53,
	}[key]
	if !ok {
		return fmt.Errorf("unsupported key %q", key)
	}

	return Run(
		"osascript",
		"-e",
		fmt.Sprintf(`tell application "System Events" to key code %d`, keyCode),
	)
}

func ScreenSize() (int, int, error) {
	script := strings.Join([]string{
		"import AppKit",
		"let frame = NSScreen.main?.frame",
		`print("\(Int(frame?.width ?? 0)),\(Int(frame?.height ?? 0))")`,
	}, "\n")

	output, err := Output("swift", "-e", script)
	if err != nil {
		return 0, 0, err
	}

	parts := strings.Split(output, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected screen size output: %q", output)
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid screen width: %w", err)
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid screen height: %w", err)
	}

	return width, height, nil
}

func OCRTSV(path string) (string, error) {
	return Output("tesseract", path, "stdout", "tsv")
}
