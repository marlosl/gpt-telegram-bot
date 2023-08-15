package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Command
	}{
		{
			name:     "Create Image Command",
			input:    "/createimage some more text here",
			expected: CreateImageCommand,
		},
		{
			name:     "Edit Command",
			input:    "/edit some more text here",
			expected: EditCommand,
		},
		{
			name:     "No Command",
			input:    "random text without command",
			expected: None,
		},
		{
			name:     "Short Text",
			input:    "/edit",
			expected: EditCommand,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetCommand(&tt.input)
			assert.Equal(t, tt.expected, cmd)
		})
	}
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name         string
		cmd          Command
		input        string
		expectedText string
		expectedInst string
	}{
		{
			name:         "Parse Edit Message with Instruction",
			cmd:          EditCommand,
			input:        "/edit rotate: 90 degrees",
			expectedText: "90 degrees",
			expectedInst: "rotate",
		},
		{
			name:         "Parse Edit Message without Instruction",
			cmd:          EditCommand,
			input:        "/edit 90 degrees",
			expectedText: "90 degrees",
			expectedInst: "",
		},
		{
			name:         "Parse CreateImage Message",
			cmd:          CreateImageCommand,
			input:        "/createimage create a landscape",
			expectedText: "create a landscape",
			expectedInst: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, inst := ParseMessage(tt.cmd, &tt.input)
			assert.Equal(t, tt.expectedText, *text)
			assert.Equal(t, tt.expectedInst, *inst)
		})
	}
}
