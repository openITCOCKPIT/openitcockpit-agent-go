package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-viper/encoding/ini"
	"github.com/openITCOCKPIT/openitcockpit-agent-go/basiclog"
	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Definitions (custom checks, logfiles, certificates, ...) can be split into several files. Next to each ini file a
// drop-in directory with the same name and the extension .d is read, e.g. customchecks.ini and customchecks.d/*.ini.
// This allows configuration management tools (Ansible, Puppet, SCCM, ...) to deploy single checks to selected hosts.

// dropInDirectory returns the drop-in directory of a configuration file, e.g. /etc/openitcockpit-agent/customchecks.d
// for /etc/openitcockpit-agent/customchecks.ini
func dropInDirectory(configPath string) string {
	return strings.TrimSuffix(configPath, filepath.Ext(configPath)) + ".d"
}

// definitionFiles returns the existing files of a definition type: the configuration file itself and the *.ini files
// of its drop-in directory in alphabetical order
func definitionFiles(configPath string) []string {
	files := make([]string, 0)
	if configPath == "" {
		return files
	}
	if utils.FileExists(configPath) {
		files = append(files, configPath)
	}
	return append(files, dropInFiles(dropInDirectory(configPath))...)
}

// dropInFiles returns the *.ini files of a drop-in directory in alphabetical order. Hidden files (e.g. editor swap
// files) and sub directories are ignored
func dropInFiles(directory string) []string {
	files := make([]string, 0)
	entries, err := os.ReadDir(directory)
	if err != nil {
		// A missing drop-in directory is not an error
		return files
	}
	for _, entry := range entries {
		if isDropInFile(entry) {
			files = append(files, filepath.Join(directory, entry.Name()))
		}
	}
	return files
}

func isDropInFile(entry os.DirEntry) bool {
	if entry.IsDir() {
		return false
	}
	name := entry.Name()
	if strings.HasPrefix(name, ".") {
		return false
	}
	return strings.EqualFold(filepath.Ext(name), ".ini")
}

// hasDefinitionFiles returns true if the configuration file or at least one file of its drop-in directory exists
func hasDefinitionFiles(configPath string) bool {
	return len(definitionFiles(configPath)) > 0
}

// readIniSections reads the ini file and the files of its drop-in directory with one section per definition.
// A section of a later file replaces the section with the same name of an earlier file completely, so a definition
// of the main file can be changed or disabled (enabled = False) by a drop-in file.
// Files which can not be read are skipped, an error is only returned if no file could be read.
func readIniSections[T any](configPath string) (map[string]*T, error) {
	sections := map[string]*T{}
	origins := map[string]string{}
	files := definitionFiles(configPath)
	errs := make([]error, 0)

	for _, file := range files {
		fileSections, err := readIniFile[T](file)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", file, err))
			continue
		}
		for name, section := range fileSections {
			if origin, exists := origins[name]; exists {
				logDefinitionReplaced(name, origin, file)
			}
			sections[name] = section
			origins[name] = file
		}
	}

	if len(errs) > 0 && len(errs) == len(files) {
		return nil, errors.Join(errs...)
	}
	for _, err := range errs {
		logger, _ := basiclog.New()
		logger.Errorln("Configuration: skipped file: ", err)
	}
	return sections, nil
}

// readIniFile reads one ini file with one section per definition
func readIniFile[T any](configPath string) (map[string]*T, error) {
	// Register the INI encoding
	// As ini support got removed from viper with v1.20.0
	// https://github.com/spf13/viper/releases/tag/v1.20.0
	// https://github.com/spf13/viper/blob/master/UPGRADE.md#v120x
	codecRegistry := viper.NewCodecRegistry()
	codecRegistry.RegisterCodec("ini", ini.Codec{})

	v := viper.NewWithOptions(
		viper.WithCodecRegistry(codecRegistry),
	)

	v.SetConfigFile(configPath)
	v.SetConfigType("ini")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := map[string]*T{}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	delete(cfg, "default")
	return cfg, nil
}

// logDefinitionReplaced logs a replaced section as information, as replacing definitions with drop-in files is intended
func logDefinitionReplaced(name, origin, file string) {
	log.Infof("Configuration: [%s] of %s replaces the definition of %s", name, file, origin)
}
