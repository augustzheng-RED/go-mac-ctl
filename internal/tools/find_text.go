package tools

import (
	"fmt"
	"time"

	"go-mac-ctl/internal/act"
	"go-mac-ctl/internal/executor"
	"go-mac-ctl/internal/observe"
)

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

type Target struct {
	Text       string
	Confidence float64
	Bounds     Rect
}

type FindTextResult struct {
	Query      string
	ObservedAt time.Time
	Screen     Rect
	Screenshot string
	Targets    []Target
	Usage      map[string]any
}

func FindText(query string) (FindTextResult, error) {
	width, height, err := observe.ScreenSize()
	if err != nil {
		return FindTextResult{}, err
	}

	screenshot, err := observe.Screenshot()
	if err != nil {
		return FindTextResult{}, err
	}

	return findTextInImageFile(screenshot, query, Rect{
		X:      0,
		Y:      0,
		Width:  width,
		Height: height,
	})
}

func findTextInImageFile(path, query string, screen Rect) (FindTextResult, error) {
	return findTextInImageFileWithScale(path, query, screen, 2.0)
}

func findTextInImageFileWithScale(path, query string, screen Rect, scale float64) (FindTextResult, error) {
	ocrTargets, err := executor.OCRFindTextWithScale(path, query, scale)
	if err != nil {
		return FindTextResult{}, err
	}

	return FindTextResult{
		Query:      query,
		ObservedAt: time.Now(),
		Screen:     screen,
		Screenshot: path,
		Targets:    fromOCRTargets(ocrTargets),
	}, nil
}

func ClickText(query string, index int) (Target, error) {
	result, err := FindText(query)
	if err != nil {
		return Target{}, err
	}

	return clickTarget(result.Targets, index)
}

func clickTarget(targets []Target, index int) (Target, error) {
	target, err := selectTarget(targets, index)
	if err != nil {
		return Target{}, err
	}

	centerX := target.Bounds.X + target.Bounds.Width/2
	centerY := target.Bounds.Y + target.Bounds.Height/2

	if err := act.LeftClick(centerX, centerY); err != nil {
		return Target{}, err
	}

	return target, nil
}

func selectTarget(targets []Target, index int) (Target, error) {
	if index < 0 || index >= len(targets) {
		return Target{}, fmt.Errorf("target index %d out of range", index)
	}

	return targets[index], nil
}

func fromOCRTargets(ocrTargets []executor.OCRTarget) []Target {
	targets := make([]Target, 0, len(ocrTargets))
	for _, item := range ocrTargets {
		targets = append(targets, Target{
			Text:       item.Text,
			Confidence: item.Confidence,
			Bounds: Rect{
				X:      item.Bounds.X,
				Y:      item.Bounds.Y,
				Width:  item.Bounds.Width,
				Height: item.Bounds.Height,
			},
		})
	}

	return targets
}
