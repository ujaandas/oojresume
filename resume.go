package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

/*
Define the fields we can modify (or are expected to be modified) in the template files.

All templates are referenced by template name (ie; "edu_hkust" or "proj_dissertation" as
defined in their respective .tmpl files)
*/

type Identity struct {
	Name     string
	Email    string
	Phone    string
	LinkedIn string
	Github   string
	Website  string
}

type EducationTmplName string
type ExperienceTmplName string
type ProjectTmplName string
type SkillsetTmplName string
type AwardTmplName string

type Section struct {
	Title            string
	Entries          []string
	EntryVSpace      *int
	SectionVSpace    *int
	EntryVSpaceTex   string
	SectionVSpaceTex string
}

type Resume struct {
	Identity Identity
	Sections []Section
}

func getResumeConfig(dir, cfgFilename string) ([]Resume, error) {
	cfg := []Resume{}

	cfgFile, err := os.Open(filepath.Join(dir, cfgFilename))
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	jsonParser := json.NewDecoder(cfgFile)
	err = jsonParser.Decode(&cfg)
	if err != nil {
		if err == io.EOF {
			return cfg, nil
		}
		return nil, err
	}

	return cfg, nil
}
