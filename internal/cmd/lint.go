package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zexot-com/commitlint/config"
	"github.com/zexot-com/commitlint/lint"
	"github.com/urfave/cli/v2"
)

const (
	// errExitCode represents the error exit code
	errExitCode = 1
)

// lintMsg is the callback function for lint command
func lintMsg(confPath, msgPath string) error {
	// NOTE: lint should return with exit code for error case
	resStr, hasError, err := runLint(confPath, msgPath)
	if handleError(err, "Linting failed") != nil {
		return err
	}

	if hasError {
		return cli.Exit(resStr, errExitCode)
	}

	// print success message
	fmt.Println(resStr)
	return nil
}

func runLint(confFilePath, fileInput string) (lintResult string, hasError bool, err error) {
	linter, format, err := getLinter(confFilePath)
	if handleError(err, "Failed to create linter") != nil {
		return "", false, err
	}

	commitMsg, err := getCommitMsg(fileInput)
	if handleError(err, "Failed to read commit message") != nil {
		return "", false, err
	}

	result, err := linter.ParseAndLint(commitMsg)
	if handleError(err, "Linting process failed") != nil {
		return "", false, err
	}

	output, err := format.Format(result)
	if handleError(err, "Formatting result failed") != nil {
		return "", false, err
	}

	return output, hasErrorSeverity(result), nil
}

func getLinter(confParam string) (*lint.Linter, lint.Formatter, error) {
	conf, err := getConfig(confParam)
	if handleError(err, "Failed to get configuration") != nil {
		return nil, nil, err
	}

	format, err := config.GetFormatter(conf)
	if handleError(err, "Failed to get formatter") != nil {
		return nil, nil, err
	}

	linter, err := config.NewLinter(conf)
	if handleError(err, "Failed to create new linter") != nil {
		return nil, nil, err
	}

	return linter, format, nil
}

func getConfig(confParam string) (*lint.Config, error) {
	if confParam != "" {
		confParam = filepath.Clean(confParam)
		return config.Parse(confParam)
	}

	// If config param is empty, lookup for defaults
	conf, err := config.LookupAndParse()
	if handleError(err, "Failed to lookup and parse configuration") != nil {
		return nil, err
	}

	return conf, nil
}

func getCommitMsg(fileInput string) (string, error) {
	commitMsg, err := readStdInPipe()
	if handleError(err, "Failed to read commit message from stdin") != nil {
		return "", err
	}

	if commitMsg != "" {
		return commitMsg, nil
	}

	// TODO: check if currentDir is inside git repo?
	if fileInput == "" {
		fileInput = "./.git/COMMIT_EDITMSG"
	}

	fileInput = filepath.Clean(fileInput)
	inBytes, err := os.ReadFile(fileInput)
	if handleError(err, "Failed to read commit message file") != nil {
		return "", err
	}
	return string(inBytes), nil
}

func readStdInPipe() (string, error) {
	stat, err := os.Stdin.Stat()
	if handleError(err, "Failed to read stdin pipe status") != nil {
		return "", err
	}

	// user input from terminal
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// not handling this case
		return "", nil
	}

	// user input from stdin pipe
	readBytes, err := io.ReadAll(os.Stdin)
	if handleError(err, "Failed to read from stdin pipe") != nil {
		return "", err
	}
	s := string(readBytes)
	return strings.TrimSpace(s), nil
}

func hasErrorSeverity(result *lint.Result) bool {
	for _, i := range result.Issues() {
		if i.Severity() == lint.SeverityError {
			return true
		}
	}
	return false
}
