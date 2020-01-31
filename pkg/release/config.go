package release

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/ExploratoryEngineering/releasetool/pkg/toolbox"
)

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type Config struct {
	SourceRoot    string   `json:"sourceRoot"`
	Architectures []string `json:"architectures"`
	OSes          []string `json:"oses"`
	Files         []File   `json:"files"`
}

const ConfigPath = "release/config.json"

// WriteSampleConfig writes a sample configuration to the release directory
func WriteSampleConfig() error {
	_, err := os.Stat(ConfigPath)
	if !os.IsNotExist(err) {
		toolbox.PrintError("Configuration file already exists")
		return err
	}

	c := sampleConfig()
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
func sampleConfig() Config {
	return Config{
		SourceRoot:    ".",
		Architectures: []string{"TODO: set architecture (amd64, arm, 386, mips...)"},
		OSes:          []string{"TODO: set operating system (darwin, linux, netbsd, openbsp, plan9, windows...)"},
		Files: []File{
			File{
				ID:   "TODO: set ID for file",
				Name: "TODO: Add your built files here",
				OS:   "TODO: Set OS for file here, '-' if it doesn't apply",
				Arch: "TODO: Set architecture for file here, '-' if it doesn't apply",
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

// VerifyConfig verifies that the artifact config is correct
func VerifyConfig(config Config) error {

	return errors.New("not implemented")
}
