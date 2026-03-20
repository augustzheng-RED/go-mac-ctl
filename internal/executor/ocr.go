package executor

import (
	"encoding/json"
	"fmt"
)

type OCRBounds struct {
	X      int
	Y      int
	Width  int
	Height int
}

type OCRTarget struct {
	Text       string
	Confidence float64
	Bounds     OCRBounds
}

type ocrFindTextResult struct {
	Targets []OCRTarget
}

func OCRFindText(path, query string) ([]OCRTarget, error) {
	return OCRFindTextWithScale(path, query, 2.0)
}

func OCRFindTextWithScale(path, query string, scale float64) ([]OCRTarget, error) {
	output, err := Output("python3", "scripts/ocr_find_text.py", path, query, fmt.Sprintf("%g", scale))
	if err != nil {
		return nil, err
	}

	var result ocrFindTextResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to decode OCR output: %w", err)
	}

	return result.Targets, nil
}
