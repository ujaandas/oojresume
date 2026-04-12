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

func formatVSpace(v *int, fallback int) string {
	value := fallback
	if v != nil {
		value = *v
	}
	return fmt.Sprintf("\\vspace{%dpt}", value)
}

func processEntries(tmpl *template.Template, r Resume) (Resume, error) {
	processed := r

	for i := range processed.Sections {
		processed.Sections[i].EntryVSpaceTex = formatVSpace(processed.Sections[i].EntryVSpace, -5)
		processed.Sections[i].SectionVSpaceTex = formatVSpace(processed.Sections[i].SectionVSpace, -10)

		var rendered []string
		for _, entry := range processed.Sections[i].Entries {
			var buf bytes.Buffer
			if err := tmpl.ExecuteTemplate(&buf, entry, nil); err != nil {
				return Resume{}, fmt.Errorf("failed to render %q: %v", entry, err)
			}
			rendered = append(rendered, buf.String())
		}
		processed.Sections[i].Entries = rendered
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

	for _, section := range r.Sections {
		for _, entry := range section.Entries {
			if !tmplNames[entry] {
				missing = append(missing, entry)
			}
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing templates: %v", missing)
	}

	return nil
}
