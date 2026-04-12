package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetTemplateFiles(t *testing.T) {
	dir := t.TempDir()

	mustWriteFile(t, filepath.Join(dir, "b.tmpl"), "b")
	mustWriteFile(t, filepath.Join(dir, "a.tmpl"), "a")
	mustWriteFile(t, filepath.Join(dir, "ignore.txt"), "x")
	mustWriteFile(t, filepath.Join(dir, "nested", "c.tmpl"), "c")

	files, err := getTemplateFiles(dir)
	if err != nil {
		t.Fatalf("getTemplateFiles() error = %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("expected 3 template files, got %d", len(files))
	}

	if filepath.Base(files[0]) != "a.tmpl" {
		t.Fatalf("expected sorted files, first file = %s", files[0])
	}
}

func TestParseLatexTemplates(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.tmpl"), "hello")

	tmpl, err := parseLatexTemplates(dir)
	if err != nil {
		t.Fatalf("parseLatexTemplates() error = %v", err)
	}
	if tmpl.Lookup("main.tmpl") == nil {
		t.Fatal("expected template main.tmpl to be parsed")
	}
}

func TestParseLatexTemplates_NoTemplates(t *testing.T) {
	dir := t.TempDir()

	_, err := parseLatexTemplates(dir)
	if err == nil {
		t.Fatal("expected error when no templates are present")
	}
}

func TestRenderResume(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.tex.tmpl"), "Name={{ .Identity.Name }} Skills={{ .Skills }}")

	tmpl, err := parseLatexTemplates(dir)
	if err != nil {
		t.Fatalf("parseLatexTemplates() error = %v", err)
	}

	rendered, err := renderResume(tmpl, Resume{
		Identity: Identity{Name: "Ujaan Das"},
		Skills:   SkillsetTmplName("skills_default"),
	})
	if err != nil {
		t.Fatalf("renderResume() error = %v", err)
	}

	if rendered != "Name=Ujaan Das Skills=skills_default" {
		t.Fatalf("unexpected render output: %q", rendered)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write failed: %v", err)
	}
}
