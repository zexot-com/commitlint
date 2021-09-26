// Package config contains helpers, defaults for linter
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/conventionalcommit/commitlint/lint"
)

const (
	// DefaultFile represent default config file name
	DefaultFile = "commitlint.yaml"
)

// GetConfig returns parses config file, validate it and returns config instance
func GetConfig(confPath string) (*lint.Config, error) {
	confFilePath, useDefault, err := getConfigPath(confPath)
	if err != nil {
		return nil, err
	}

	if useDefault {
		return defConf, nil
	}

	conf, err := Parse(confFilePath)
	if err != nil {
		return nil, err
	}

	err = Validate(conf)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return conf, nil
}

// getConfigPath returns config file path following below order
// 	1. commitlint.yaml in current directory
// 	2. confFilePath parameter
// 	3. use default config
func getConfigPath(confFilePath string) (confPath string, isDefault bool, retErr error) {
	// get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", false, err
	}

	// check if conf file exists in current directory
	currentDirConf := filepath.Join(currentDir, DefaultFile)
	if _, err1 := os.Stat(currentDirConf); !os.IsNotExist(err1) {
		return currentDirConf, false, nil
	}

	// if confFilePath empty,
	// means no config in current directory or config flag is empty
	// use default config
	if confFilePath == "" {
		return "", true, nil
	}
	return filepath.Clean(confFilePath), false, nil
}

// Parse parse Config from given file
func Parse(confPath string) (*lint.Config, error) {
	confBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	conf := &lint.Config{}
	err = yaml.Unmarshal(confBytes, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// Validate parses Config from given data
func Validate(conf *lint.Config) error {
	if conf.Formatter == "" {
		return errors.New("formatter is empty")
	}

	_, ok := globalRegistry.GetFormatter(conf.Formatter)
	if !ok {
		return fmt.Errorf("unknown formatter '%s'", conf.Formatter)
	}

	for ruleName, r := range conf.Rules {
		// Check Severity Level of rule config
		switch r.Severity {
		case lint.SeverityError:
		case lint.SeverityWarn:
		default:
			return fmt.Errorf("unknown severity level '%s' for rule '%s'", r.Severity, ruleName)
		}

		// Check if rule is registered
		_, ok := globalRegistry.GetRule(ruleName)
		if !ok {
			return fmt.Errorf("unknown rule '%s'", ruleName)
		}
	}

	return nil
}

// WriteConfToFile util func to write config object to given file
func WriteConfToFile(outFilePath string, conf *lint.Config) (retErr error) {
	file, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err1 := file.Close()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	w := bufio.NewWriter(file)
	defer func() {
		err1 := w.Flush()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	enc := yaml.NewEncoder(w)
	return enc.Encode(conf)
}
