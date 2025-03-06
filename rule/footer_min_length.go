package rule

import "github.com/zexot-com/commitlint/lint"

var _ lint.Rule = (*FooterMinLenRule)(nil)

// FooterMinLenRule to validate min length of footer
type FooterMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *FooterMinLenRule) Name() string { return "footer-min-length" }

// Apply sets the needed argument for the rule
func (r *FooterMinLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates FooterMinLenRule
func (r *FooterMinLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMinLen("footer", r.CheckLen, msg.Footer())
}
