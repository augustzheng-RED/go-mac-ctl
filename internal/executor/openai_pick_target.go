package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type OpenAITargetCandidate struct {
	Index      int
	Text       string
	Confidence float64
	Bounds     OCRBounds
}

type OpenAITargetSelection struct {
	CandidateIndex int
	Reason         string
	Usage          map[string]any
}

func OpenAIPickTarget(path, instruction string, candidates []OpenAITargetCandidate) (OpenAITargetSelection, error) {
	payload, err := json.Marshal(candidates)
	if err != nil {
		return OpenAITargetSelection{}, fmt.Errorf("failed to encode candidates: %w", err)
	}

	cmd := exec.Command("python3", "scripts/openai_pick_target.py", path, instruction)
	cmd.Stdin = bytes.NewReader(payload)

	output, err := cmd.CombinedOutput()
	if err != nil {
		message := string(bytes.TrimSpace(output))
		if message != "" {
			return OpenAITargetSelection{}, fmt.Errorf("failed to pick target with OpenAI: %w: %s", err, message)
		}

		return OpenAITargetSelection{}, fmt.Errorf("failed to pick target with OpenAI: %w", err)
	}

	var selection OpenAITargetSelection
	if err := json.Unmarshal(bytes.TrimSpace(output), &selection); err != nil {
		return OpenAITargetSelection{}, fmt.Errorf("failed to decode OpenAI target selection: %w", err)
	}

	return selection, nil
}
