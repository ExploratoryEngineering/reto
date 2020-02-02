package release

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/ExploratoryEngineering/reto/pkg/toolbox"
)

// File is the configuration setting for a single file
type File struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Target string `json:"target"`
}

// Actions on templates
const (
	// IncludeAction means that the file will be included in the archive
	// without any other actions
	IncludeAction = "include"
	// ConcatenateAction will concatenate across all previous releases into
	// one big file before it is included in the archive.
	ConcatenateAction = "concatenate"
)

// Templates are template files that are included in the release archive. They
// are Go templates
type Template struct {
	Name           string `json:"name"`
	TemplateAction string `json:"action"`
}

// Config is the tool configuration
type Config struct {
	SourceRoot     string     `json:"sourceRoot"`
	Name           string     `json:"name"`
	CommitterEmail string     `json:"committerEmail"`
	CommitterName  string     `json:"committerName"`
	Targets        []string   `json:"targets"`
	Files          []File     `json:"files"`
	Templates      []Template `json:"templates"`
}

// ConfigPath is the path to the configuration file
const ConfigPath = "release/config.json"

// WriteSampleConfig writes a sample configuration to the release directory
func WriteSampleConfig() error {
	_, err := os.Stat(ConfigPath)
	if !os.IsNotExist(err) {
		toolbox.PrintError("Configuration file already exists")
		return err
	}

	c := defaultConfig()
	buf, err := json.MarshalIndent(&c, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(ConfigPath, buf, toolbox.DefaultFilePerm); err != nil {
		toolbox.PrintError("Could not write sample config: %v", err)
		return err
	}
	return nil
}

// sampleConfig is the sample configuration file.
func defaultConfig() Config {
	return Config{
		SourceRoot:     ".",
		Name:           "TODO set your product name",
		CommitterName:  "TODO set name for git commits",
		CommitterEmail: "TODO set email for git commits",
		Targets:        []string{"TODO: set target (amd64-darwin, arm-linux, mips-plan9...)"},
		Templates: []Template{Template{
			Name:           "changelog.md",
			TemplateAction: ConcatenateAction,
		}},
		Files: []File{
			File{
				ID:     "TODO: set ID for file",
				Name:   "TODO: Add your built files here",
				Target: "TODO: Set target for file here, '-' if it doesn't apply",
			},
		},
	}
}

func readConfig() (Config, error) {
	buf, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		toolbox.PrintError("Could not read configuration: %v", err)
		return Config{}, err
	}
	ret := Config{}
	if err := json.Unmarshal(buf, &ret); err != nil {
		toolbox.PrintError("Configuration file format error: %v", err)
		return Config{}, err
	}
	return ret, nil
}

const anyTarget = "-"

// VerifyConfig verifies that the artifact config is correct
func VerifyConfig(config Config, printErrors bool) error {
	if err := toolbox.CheckForTODO(ConfigPath, true); err != nil {
		return err
	}
	if len(config.CommitterEmail) == 0 || len(config.CommitterName) == 0 {
		if printErrors {
			toolbox.PrintError("Committer email and name must be set")
		}
		return errors.New("no committer")
	}
	if len(config.Targets) == 0 {
		if printErrors {
			toolbox.PrintError("There are no targets in the configuration file")
		}
		return errors.New("no targets")
	}
	if len(config.Files) == 0 {
		if printErrors {
			toolbox.PrintError("There are no output files in the configuration file")
		}
		return errors.New("no targets")
	}

	fileTargets := make(map[string]map[string]bool)
	for _, v := range config.Files {
		if v.Target == anyTarget {
			continue
		}
		existing, ok := fileTargets[v.ID]
		if !ok {
			existing = make(map[string]bool)
		}
		existing[v.Target] = true
		fileTargets[v.ID] = existing
	}

	errs := 0
	for id, v := range fileTargets {
		targets := make(map[string]bool)
		for _, t := range config.Targets {
			if t == anyTarget {
				continue
			}
			targets[t] = true
		}
		for target := range v {
			if target == anyTarget {
				continue
			}
			_, ok := targets[target]
			if !ok {
				if printErrors {
					toolbox.PrintError("File with ID '%s' have unknown target %s", id, target)
				}
				errs++
			}
			delete(targets, target)
		}
		if len(targets) > 0 {
			for target := range targets {
				if printErrors {
					toolbox.PrintError("File with ID '%s' is missing target %s", id, target)
				}
				errs++
			}
		}
	}
	if errs > 0 {
		return errors.New("target missing")
	}

	for _, file := range config.Files {
		if _, err := os.Stat(file.Name); err != nil {
			if os.IsNotExist(err) {
				if printErrors {
					toolbox.PrintError("The file %s does not exist", file.Name)
				}
				return err
			}
			toolbox.PrintError("Could not access %s: %v", file.Name, err)
			return err
		}
	}
	return nil
}
