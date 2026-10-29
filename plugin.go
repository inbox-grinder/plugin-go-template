package main

import (
	"encoding/json"
	"strings"
)

// EmailInput is the JSON structure the host sends to this plugin.
type EmailInput struct {
	Sender   SenderInfo `json:"sender"`
	Subject  string     `json:"subject"`
	TextBody *string    `json:"textBody"`
}

// SenderInfo holds the sender's address.
type SenderInfo struct {
	Address string `json:"address"`
}

// ExtractionResult is the JSON structure this plugin returns to the host.
type ExtractionResult struct {
	SchemaVersion string            `json:"schemaVersion"`
	Type          string            `json:"type"`
	Template      string            `json:"template"`
	Confidence    float64           `json:"confidence"`
	Summary       string            `json:"summary"`
	Fields        map[string]any    `json:"fields"`
	Actions       []ExtractionAction `json:"actions"`
}

// ExtractionAction represents a suggested action the user can take.
type ExtractionAction struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Label string `json:"label"`
}

const previewLength = 200

// run parses the EmailInput JSON, extracts data, and returns ExtractionResult JSON.
func run(input []byte) ([]byte, error) {
	var email EmailInput
	if err := json.Unmarshal(input, &email); err != nil {
		return nil, err
	}

	textBody := ""
	if email.TextBody != nil {
		textBody = *email.TextBody
	}

	preview := textBody
	if len([]rune(preview)) > previewLength {
		preview = string([]rune(preview)[:previewLength]) + "…"
	}
	wordCount := len(strings.Fields(textBody))

	result := ExtractionResult{
		SchemaVersion: "1",
		Type:          "go_template.text_preview",
		Template:      "text-preview",
		Confidence:    0.8,
		Summary:       email.Sender.Address + " — " + email.Subject,
		Fields: map[string]any{
			"sender":    email.Sender.Address,
			"subject":   email.Subject,
			"preview":   preview,
			"wordCount": wordCount,
		},
		Actions: []ExtractionAction{{
			ID:    "archive",
			Kind:  "suggest_archive",
			Label: "Archive",
		}},
	}

	return json.Marshal(result)
}
