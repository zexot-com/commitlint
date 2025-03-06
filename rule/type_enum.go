package rule

import (
	"fmt"
	"sort"

	"github.com/zexot-com/commitlint/lint"
)

var _ lint.Rule = (*TypeEnumRule)(nil)

// TypeEnumRule to validate types
type TypeEnumRule struct {
	Types []string
}

// Name return name of the rule
func (r *TypeEnumRule) Name() string { return "type-enum" }

// Apply sets the needed argument for the rule
func (r *TypeEnumRule) Apply(setting lint.RuleSetting) error {
	err := setStringArrArg(&r.Types, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	// sorting the string elements for binary search
	sort.Strings(r.Types)
	return nil
}

// Validate validates TypeEnumRule
func (r *TypeEnumRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	isFound := search(r.Types, msg.Type())
	if isFound {
		return nil, true
	}
	desc := fmt.Sprintf("type '%s' is not allowed, you can use one of %v", msg.Type(), r.Types)
	return lint.NewIssue(desc), false
}
