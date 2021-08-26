// Package config contains helpers, defaults for linter
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/conventionalcommit/commitlint/formatter"
	"github.com/conventionalcommit/commitlint/lint"
	"github.com/conventionalcommit/commitlint/rule"
)

const (
	// ConfFileName represent config file name
	ConfFileName = "commitlint.yaml"
)

var allFormatters = []lint.Formatter{
	&formatter.DefaultFormatter{},
	&formatter.JSONFormatter{},
}

var allRules = []lint.Rule{
	&rule.BodyMinLenRule{},
	&rule.BodyMaxLenRule{},

	&rule.FooterMinLenRule{},
	&rule.FooterMaxLenRule{},

	&rule.HeadMaxLenRule{},
	&rule.HeadMinLenRule{},

	&rule.TypeEnumRule{},
	&rule.ScopeEnumRule{},

	&rule.BodyMaxLineLenRule{},
	&rule.FooterMaxLineLenRule{},

	&rule.TypeCharsetRule{},
	&rule.ScopeCharsetRule{},

	&rule.TypeMaxLenRule{},
	&rule.ScopeMaxLenRule{},
	&rule.DescriptionMaxLenRule{},

	&rule.TypeMinLenRule{},
	&rule.ScopeMinLenRule{},
	&rule.DescriptionMinLenRule{},
}

// GetConfig returns conf and follows the belwo precedence
// 	1. check for conf in current directory
// 	2. check for conf flag
// 	3. load default conf
func GetConfig(confFilePath string) (*lint.Config, error) {
	confFilePath = filepath.Clean(confFilePath)

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// check if conf file exists in current directory
	currentDirConf := filepath.Join(currentDir, ConfFileName)
	if _, err1 := os.Stat(currentDirConf); !os.IsNotExist(err1) {
		confFilePath = currentDirConf
	}

	if confFilePath == "" {
		return defConf, nil
	}

	conf, err := Parse(confFilePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// GetFormatter returns the formatter as defined in conf
func GetFormatter(c *lint.Config) (lint.Formatter, error) {
	for _, f := range allFormatters {
		if f.Name() == c.Formatter {
			return f, nil
		}
	}
	return nil, fmt.Errorf("%s formatter not found", c.Formatter)
}

// GetRules forms Rule object for rules which are enabled in config
func GetRules(conf *lint.Config) ([]lint.Rule, error) {
	// rules lookup map
	rulesMap := map[string]lint.Rule{}
	for _, r := range allRules {
		rulesMap[r.Name()] = r
	}

	enabledRules := make([]lint.Rule, 0, len(conf.Rules))

	for ruleName, ruleConfig := range conf.Rules {
		r, ok := rulesMap[ruleName]
		if !ok {
			return nil, fmt.Errorf("unknown rule: %s", ruleName)
		}
		if ruleConfig.Enabled {
			err := r.SetAndCheckArgument(ruleConfig.Argument)
			if err != nil {
				return nil, err
			}
			enabledRules = append(enabledRules, r)
		}
	}

	return enabledRules, nil
}

// GetLinter returns Linter for given confFilePath
func GetLinter(conf *lint.Config) (*lint.Linter, error) {
	rules, err := GetRules(conf)
	if err != nil {
		return nil, err
	}
	return lint.NewLinter(conf, rules)
}
