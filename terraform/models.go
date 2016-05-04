package terraform

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Model struct {
	Source              string                 `json:"source"`
	Vars                map[string]interface{} `json:"vars,omitempty"`     // optional
	VarFile             string                 `json:"var_file,omitempty"` // optional
	StateFileLocalPath  string                 `json:"-"`                  // not specified pipeline
	StateFileRemotePath string                 `json:"-"`                  // not specified pipeline
}

func (m Model) Validate() error {
	missingFields := []string{}
	if m.StateFileLocalPath == "" {
		missingFields = append(missingFields, "state_file_local_path")
	}
	if m.StateFileRemotePath == "" {
		missingFields = append(missingFields, "state_file_remote_path")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("Missing required terraform fields: %s", strings.Join(missingFields, ", "))
	}
	return nil
}

func (m Model) Merge(other Model) Model {
	mergedVars := map[string]interface{}{}
	for key, value := range m.Vars {
		mergedVars[key] = value
	}
	for key, value := range other.Vars {
		mergedVars[key] = value
	}
	m.Vars = mergedVars
	if other.Source != "" {
		m.Source = other.Source
	}
	if other.VarFile != "" {
		m.VarFile = other.VarFile
	}
	if other.StateFileLocalPath != "" {
		m.StateFileLocalPath = other.StateFileLocalPath
	}
	if other.StateFileRemotePath != "" {
		m.StateFileRemotePath = other.StateFileRemotePath
	}

	return m
}

func (m *Model) ParseVarsFromFile() error {
	terraformVars := map[string]interface{}{}
	for key, value := range m.Vars {
		terraformVars[key] = value
	}

	if m.VarFile != "" {
		fileContents, readErr := ioutil.ReadFile(m.VarFile)
		if readErr != nil {
			return fmt.Errorf("Failed to read TerraformVarFile at '%s': %s", m.VarFile, readErr)
		}

		fileVars := map[string]interface{}{}
		readErr = yaml.Unmarshal(fileContents, &fileVars)
		if readErr != nil {
			return fmt.Errorf("Failed to parse TerraformVarFile at '%s': %s", m.VarFile, readErr)
		}

		for key, value := range fileVars {
			terraformVars[key] = value
		}
	}

	m.Vars = terraformVars

	return nil
}
