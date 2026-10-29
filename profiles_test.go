package nats

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestProfileConfigYamlTimeout(t *testing.T) {
	cfg := ProfileConfig{
		ProfileID:   12,
		Exclusions:  []string{"ex1"},
		Deployments: []string{"dep1"},
		Timeout:     420,
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal ProfileConfig: %v", err)
	}

	var parsed ProfileConfig
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal ProfileConfig: %v", err)
	}

	if parsed.Timeout != 420 {
		t.Fatalf("expected Timeout 420, got %d", parsed.Timeout)
	}
	if parsed.ProfileID != 12 {
		t.Fatalf("expected ProfileID 12, got %d", parsed.ProfileID)
	}
}
