package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"text/template"
)

func getTemplateFiles(dir string) ([]string, error) {
	tmplPaths := []string{}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}

		tmplPaths = append(tmplPaths, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(tmplPaths)
	return tmplPaths, nil
}

func parseLatexTemplates(dir string) (*template.Template, error) {
	tmplPaths, err := getTemplateFiles(dir)
	if err != nil {
		return nil, err
	}
	if len(tmplPaths) == 0 {
		return nil, fmt.Errorf("no templates found in %s", dir)
	}

	tmpl, err := template.New("resume").ParseFiles(tmplPaths...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func renderResume(tmpl *template.Template, mainTmplName string, r Resume) (string, error) {
	var out bytes.Buffer

	err := tmpl.ExecuteTemplate(&out, mainTmplName, r)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func renderEntries(tmpl *template.Template, entries []string) ([]string, error) {
	rendered := make([]string, 0, len(entries))

	for _, entry := range entries {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, entry, nil); err != nil {
			return nil, fmt.Errorf("failed to render %q: %v", entry, err)
		}
		rendered = append(rendered, buf.String())
	}

	return rendered, nil
}

func processListSection(tmpl *template.Template, section *ListSection) (*ListSection, error) {
	if section == nil {
		return nil, nil
	}

	processed := *section
	rendered, err := renderEntries(tmpl, section.Entries)
	if err != nil {
		return nil, err
	}
	processed.Entries = rendered

	return &processed, nil
}

func processSkillsSection(tmpl *template.Template, section *SkillsSection) (*SkillsSection, error) {
	if section == nil {
		return nil, nil
	}

	processed := *section
	rendered, err := renderEntries(tmpl, section.Entries)
	if err != nil {
		return nil, err
	}
	processed.Entries = rendered

	return &processed, nil
}

func processEntries(tmpl *template.Template, r Resume) (Resume, error) {
	processed := r

	var err error
	if processed.Education, err = processListSection(tmpl, processed.Education); err != nil {
		return Resume{}, fmt.Errorf("failed to render education section: %v", err)
	}
	if processed.Experience, err = processListSection(tmpl, processed.Experience); err != nil {
		return Resume{}, fmt.Errorf("failed to render experience section: %v", err)
	}
	if processed.Projects, err = processListSection(tmpl, processed.Projects); err != nil {
		return Resume{}, fmt.Errorf("failed to render projects section: %v", err)
	}
	if processed.Awards, err = processListSection(tmpl, processed.Awards); err != nil {
		return Resume{}, fmt.Errorf("failed to render awards section: %v", err)
	}
	if processed.Skills, err = processSkillsSection(tmpl, processed.Skills); err != nil {
		return Resume{}, fmt.Errorf("failed to render skills section: %v", err)
	}

	return processed, nil
}

func validateResume(tmpl *template.Template, mainTmplName string, r Resume) error {
	if tmpl.Lookup(mainTmplName) == nil {
		return fmt.Errorf("main template %q not found", mainTmplName)
	}

	var missing []string
	tmplNames := make(map[string]bool)
	for _, t := range tmpl.Templates() {
		tmplNames[t.Name()] = true
	}

	collectMissing := func(entries []string) {
		for _, entry := range entries {
			if !tmplNames[entry] {
				missing = append(missing, entry)
			}
		}
	}

	if r.Education != nil {
		collectMissing(r.Education.Entries)
	}
	if r.Experience != nil {
		collectMissing(r.Experience.Entries)
	}
	if r.Projects != nil {
		collectMissing(r.Projects.Entries)
	}
	if r.Awards != nil {
		collectMissing(r.Awards.Entries)
	}
	if r.Skills != nil {
		collectMissing(r.Skills.Entries)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing templates: %v", missing)
	}

	return nil
}
