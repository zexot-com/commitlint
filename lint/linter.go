// Package lint provides a simple linter for conventional commits
package lint

// Linter is linter for commit message
type Linter struct {
	conf  *Config
	rules []Rule
}

// New returns a new Linter instance with given config and rules
func New(conf *Config, rules []Rule) (*Linter, error) {
	return &Linter{conf: conf, rules: rules}, nil
}

// Lint checks the given commitMsg string against rules
func (l *Linter) Lint(commitMsg string) (*Failure, error) {
	msg, err := Parse(commitMsg)
	if err != nil {
		return l.parserErrorRule(commitMsg, err)
	}
	return l.LintCommit(msg)
}

// LintCommit checks the given Commit against rules
func (l *Linter) LintCommit(msg *Commit) (*Failure, error) {
	res := newFailure(msg.FullCommit)

	for _, rule := range l.rules {
		currentRule := rule
		ruleConf := l.conf.GetRule(currentRule.Name())
		ruleRes, isValid := l.runRule(currentRule, ruleConf.Severity, msg)
		if !isValid {
			res.add(ruleRes)
		}
	}

	return res, nil
}

func (l *Linter) runRule(rule Rule, severity Severity, msg *Commit) (*RuleFailure, bool) {
	failMsg, isOK := rule.Validate(msg)
	if isOK {
		return nil, true
	}
	res := newRuleFailure(rule.Name(), failMsg, severity)
	return res, false
}

func (l *Linter) parserErrorRule(commitMsg string, err error) (*Failure, error) {
	res := newFailure(commitMsg)

	var errMsg string
	if isHeaderErr(err) {
		// TODO: show more information
		errMsg = "commit header is not in valid format"
	} else {
		errMsg = err.Error()
	}

	ruleFail := newRuleFailure("parser", errMsg, SeverityError)
	res.add(ruleFail)

	return res, nil
}
