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
	mustWriteFile(t, filepath.Join(dir, "main.tex.tmpl"), "Name={{ .Identity.Name }} Sections={{ len .Sections }}")

	tmpl, err := parseLatexTemplates(dir)
	if err != nil {
		t.Fatalf("parseLatexTemplates() error = %v", err)
	}

	rendered, err := renderResume(tmpl, "main.tex.tmpl", Resume{
		Identity: Identity{Name: "Ujaan Das"},
		Sections: []Section{
			{
				Title:   "Education",
				Entries: []string{"edu_warwick", "edu_hkust"},
			},
		},
	})
	if err != nil {
		t.Fatalf("renderResume() error = %v", err)
	}

	if rendered != "Name=Ujaan Das Sections=1" {
		t.Fatalf("unexpected render output: %q", rendered)
	}
}

func TestValidateResume_Valid(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.tex.tmpl"), "test")
	mustWriteFile(t, filepath.Join(dir, "edu_warwick.tex.tmpl"), `{{ define "edu_warwick" }}education{{ end }}`)

	tmpl, err := parseLatexTemplates(dir)
	if err != nil {
		t.Fatalf("parseLatexTemplates() error = %v", err)
	}

	err = validateResume(tmpl, "main.tex.tmpl", Resume{
		Identity: Identity{Name: "Test"},
		Sections: []Section{
			{
				Title:   "Education",
				Entries: []string{"edu_warwick.tex.tmpl"},
			},
		},
	})
	if err != nil {
		t.Fatalf("validateResume() expected no error, got %v", err)
	}
}

func TestValidateResume_MissingMain(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "other.tex.tmpl"), "test")

	tmpl, err := parseLatexTemplates(dir)
	if err != nil {
		t.Fatalf("parseLatexTemplates() error = %v", err)
	}

	err = validateResume(tmpl, "main.tex.tmpl", Resume{})
	if err == nil {
		t.Fatal("validateResume() expected error for missing main template")
	}
}

func TestValidateResume_MissingEntry(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.tex.tmpl"), "test")

	tmpl, err := parseLatexTemplates(dir)
	if err != nil {
		t.Fatalf("parseLatexTemplates() error = %v", err)
	}

	err = validateResume(tmpl, "main.tex.tmpl", Resume{
		Sections: []Section{
			{
				Title:   "Education",
				Entries: []string{"missing_edu"},
			},
		},
	})
	if err == nil {
		t.Fatal("validateResume() expected error for missing entry template")
	}
}

func TestProcessEntries(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.tex.tmpl"), "main")
	mustWriteFile(t, filepath.Join(dir, "edu_test.tex.tmpl"), `{{ define "edu_test" }}CONTENT_EDU{{ end }}`)
	mustWriteFile(t, filepath.Join(dir, "work_test.tex.tmpl"), `{{ define "work_test" }}CONTENT_WORK{{ end }}`)

	tmpl, err := parseLatexTemplates(dir)
	if err != nil {
		t.Fatalf("parseLatexTemplates() error = %v", err)
	}

	input := Resume{
		Identity: Identity{Name: "Test User"},
		Sections: []Section{
			{
				Title:   "Education",
				Entries: []string{"edu_test", "work_test"},
			},
		},
	}

	processed, err := processEntries(tmpl, input)
	if err != nil {
		t.Fatalf("processEntries() error = %v", err)
	}

	if len(processed.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(processed.Sections))
	}

	if len(processed.Sections[0].Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(processed.Sections[0].Entries))
	}

	if processed.Sections[0].Entries[0] != "CONTENT_EDU" {
		t.Fatalf("expected CONTENT_EDU, got %q", processed.Sections[0].Entries[0])
	}

	if processed.Sections[0].Entries[1] != "CONTENT_WORK" {
		t.Fatalf("expected CONTENT_WORK, got %q", processed.Sections[0].Entries[1])
	}
}

func TestProcessEntries_FormatsNumericVSpace(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, filepath.Join(dir, "main.tex.tmpl"), "main")
	mustWriteFile(t, filepath.Join(dir, "edu_test.tex.tmpl"), `{{ define "edu_test" }}CONTENT_EDU{{ end }}`)

	tmpl, err := parseLatexTemplates(dir)
	if err != nil {
		t.Fatalf("parseLatexTemplates() error = %v", err)
	}

	entrySpace := -2
	sectionSpace := -7

	processed, err := processEntries(tmpl, Resume{
		Sections: []Section{
			{
				Title:         "Education",
				Entries:       []string{"edu_test"},
				EntryVSpace:   &entrySpace,
				SectionVSpace: &sectionSpace,
			},
		},
	})
	if err != nil {
		t.Fatalf("processEntries() error = %v", err)
	}

	if processed.Sections[0].EntryVSpaceTex != "\\vspace{-2pt}" {
		t.Fatalf("expected entry spacing to be formatted, got %q", processed.Sections[0].EntryVSpaceTex)
	}
	if processed.Sections[0].SectionVSpaceTex != "\\vspace{-7pt}" {
		t.Fatalf("expected section spacing to be formatted, got %q", processed.Sections[0].SectionVSpaceTex)
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
