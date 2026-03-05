package tests

import (
	"strings"
	"testing"

	"github.com/andreagrandi/mb-cli/internal/cli"
)

func TestContextContent(t *testing.T) {
	content := cli.ContextContent()

	if content == "" {
		t.Fatal("context content is empty")
	}

	requiredSections := []string{
		"# mb-cli - Agent Context",
		"## Authentication",
		"## Commands",
		"## Global Flags",
		"## Flags That Do NOT Exist",
		"## Database Name Resolution",
		"## Output Formats",
		"## Examples",
	}

	for _, section := range requiredSections {
		if !strings.Contains(content, section) {
			t.Errorf("context content missing section: %s", section)
		}
	}
}

func TestContextContentContainsKeyCommands(t *testing.T) {
	content := cli.ContextContent()

	commands := []string{
		"database list",
		"table list",
		"field get",
		"query sql",
		"card list",
		"search",
		"context",
		"version",
	}

	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Errorf("context content missing command: %s", cmd)
		}
	}
}

func TestContextContentContainsEnvVars(t *testing.T) {
	content := cli.ContextContent()

	envVars := []string{"MB_HOST", "MB_API_KEY"}
	for _, v := range envVars {
		if !strings.Contains(content, v) {
			t.Errorf("context content missing env var: %s", v)
		}
	}
}
