package main

import (
	"encoding/json"
	"os"
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

type AwardTemplName string

type ResumeConfig struct {
	Identity   Identity
	Education  []EducationTmplName
	Experience []ExperienceTmplName
	Project    []ProjectTmplName
	Skills     SkillsetTmplName
	Awards     []AwardTemplName
}

func getResumeConfig(dir, cfgFilename string) ([]ResumeConfig, error) {
	cfg := []ResumeConfig{}

	cfgFile, err := os.Open(dir + cfgFilename)
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	jsonParser := json.NewDecoder(cfgFile)
	jsonParser.Decode(&cfg)

	return cfg, nil
}
