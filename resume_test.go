package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) (string, string) {
	t.Helper()

	dir := t.TempDir()
	filename := "test.json"
	fullPath := filepath.Join(dir, filename)

	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	return dir + string(os.PathSeparator), filename
}

func TestGetResumeConfigSingle(t *testing.T) {
	json := `[
		{
			"Identity": {
				"Name": "John Doe",
				"Email": "john@example.com",
				"Phone": "123456",
				"LinkedIn": "linkedin.com/john",
				"Github": "github.com/john",
				"Website": "john.dev"
			},
			"Education": ["bsc_cs"],
			"Experience": ["backend_dev"],
			"Project": ["proj1"],
			"Skills": "golang",
			"Awards": ["award1"]
		}
	]`

	dir, file := writeTempFile(t, json)

	cfgs, err := getResumeConfig(dir, file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfgs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(cfgs))
	}

	if cfgs[0].Identity.Name != "John Doe" {
		t.Errorf("expected Name 'John Doe', got %s", cfgs[0].Identity.Name)
	}
}

func TestGetResumeConfigMultiple(t *testing.T) {
	json := `[
		{
			"Identity": { "Name": "A" },
			"Education": [],
			"Experience": [],
			"Project": [],
			"Skills": "x",
			"Awards": []
		},
		{
			"Identity": { "Name": "B" },
			"Education": [],
			"Experience": [],
			"Project": [],
			"Skills": "y",
			"Awards": []
		}
	]`

	dir, file := writeTempFile(t, json)

	cfgs, err := getResumeConfig(dir, file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfgs) != 2 {
		t.Fatalf("expected 2 configs, got %d", len(cfgs))
	}
}

func TestGetResumeConfigFileNotFound(t *testing.T) {
	_, err := getResumeConfig("/nonexistent/", "file.json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetResumeConfigInvalidJSON(t *testing.T) {
	json := `[{ invalid json }]`

	dir, file := writeTempFile(t, json)

	_, err := getResumeConfig(dir, file)

	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestGetResumeConfigEmptyFile(t *testing.T) {
	dir, file := writeTempFile(t, "")

	cfgs, err := getResumeConfig(dir, file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cfgs) != 0 {
		t.Fatalf("expected empty result, got %d", len(cfgs))
	}
}
