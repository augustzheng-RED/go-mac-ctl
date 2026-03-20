package tools

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"go-mac-ctl/internal/act"
	"go-mac-ctl/internal/executor"
)

type AIClickResult struct {
	Instruction string
	ObservedAt  time.Time
	Screenshot  string
	Candidate   Target
	Reason      string
	Usage       map[string]any
}

func AIClick(instruction string) (AIClickResult, error) {
	result, err := FindText("")
	if err != nil {
		return AIClickResult{}, err
	}

	candidates := candidateSubsetForInstruction(result.Targets, instruction)
	if len(candidates) == 0 {
		return AIClickResult{}, fmt.Errorf("no OCR candidates available for instruction %q", instruction)
	}

	selection, err := executor.OpenAIPickTarget(result.Screenshot, instruction, candidates)
	if err != nil {
		return AIClickResult{}, err
	}

	chosen, ok := resolveSelectedCandidate(candidates, selection.CandidateIndex)
	if !ok {
		return AIClickResult{}, fmt.Errorf("OpenAI did not return a valid candidate index")
	}
	target := Target{
		Text:       chosen.Text,
		Confidence: chosen.Confidence,
		Bounds: Rect{
			X:      chosen.Bounds.X,
			Y:      chosen.Bounds.Y,
			Width:  chosen.Bounds.Width,
			Height: chosen.Bounds.Height,
		},
	}

	centerX := target.Bounds.X + target.Bounds.Width/2
	centerY := target.Bounds.Y + target.Bounds.Height/2
	if err := act.LeftClick(centerX, centerY); err != nil {
		return AIClickResult{}, err
	}

	return AIClickResult{
		Instruction: instruction,
		ObservedAt:  time.Now(),
		Screenshot:  result.Screenshot,
		Candidate:   target,
		Reason:      selection.Reason,
		Usage:       selection.Usage,
	}, nil
}

func candidateSubsetForInstruction(targets []Target, instruction string) []executor.OpenAITargetCandidate {
	terms := instructionTerms(instruction)
	prioritized := make([]executor.OpenAITargetCandidate, 0, len(targets))
	fallback := make([]executor.OpenAITargetCandidate, 0, len(targets))

	for i, target := range targets {
		candidate := executor.OpenAITargetCandidate{
			Index:      i,
			Text:       target.Text,
			Confidence: target.Confidence,
			Bounds: executor.OCRBounds{
				X:      target.Bounds.X,
				Y:      target.Bounds.Y,
				Width:  target.Bounds.Width,
				Height: target.Bounds.Height,
			},
		}

		if matchesInstructionTerm(target.Text, terms) {
			prioritized = append(prioritized, candidate)
			continue
		}

		fallback = append(fallback, candidate)
	}

	sort.SliceStable(prioritized, func(i, j int) bool {
		return prioritized[i].Confidence > prioritized[j].Confidence
	})
	sort.SliceStable(fallback, func(i, j int) bool {
		return fallback[i].Confidence > fallback[j].Confidence
	})

	if len(prioritized) > 0 {
		return limitCandidates(prioritized, 40)
	}

	return limitCandidates(fallback, 80)
}

func instructionTerms(instruction string) []string {
	normalized := strings.ToLower(instruction)
	replacer := strings.NewReplacer(",", " ", ".", " ", ":", " ", ";", " ", "-", " ", "_", " ", "/", " ", "\n", " ", "\t", " ")
	normalized = replacer.Replace(normalized)

	parts := strings.Fields(normalized)
	terms := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	stopWords := map[string]struct{}{
		"a": {}, "an": {}, "and": {}, "button": {}, "click": {}, "go": {}, "item": {}, "link": {},
		"menu": {}, "on": {}, "open": {}, "select": {}, "tab": {}, "the": {}, "to": {},
	}

	for _, part := range parts {
		if len(part) < 2 {
			continue
		}
		if _, skip := stopWords[part]; skip {
			continue
		}
		if _, exists := seen[part]; exists {
			continue
		}
		seen[part] = struct{}{}
		terms = append(terms, part)
	}

	return terms
}

func matchesInstructionTerm(text string, terms []string) bool {
	if len(terms) == 0 {
		return false
	}

	lowerText := strings.ToLower(text)
	for _, term := range terms {
		if strings.Contains(lowerText, term) {
			return true
		}
	}

	return false
}

func limitCandidates(candidates []executor.OpenAITargetCandidate, limit int) []executor.OpenAITargetCandidate {
	if len(candidates) <= limit {
		return candidates
	}

	return candidates[:limit]
}

func resolveSelectedCandidate(candidates []executor.OpenAITargetCandidate, candidateIndex int) (executor.OpenAITargetCandidate, bool) {
	if candidateIndex >= 0 && candidateIndex < len(candidates) {
		return candidates[candidateIndex], true
	}

	for _, candidate := range candidates {
		if candidate.Index == candidateIndex {
			return candidate, true
		}
	}

	return executor.OpenAITargetCandidate{}, false
}
