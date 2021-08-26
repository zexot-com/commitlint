package rule

import (
	"fmt"

	"github.com/conventionalcommit/commitlint/message"
)

// HeadMaxLenRule to validate max length of header
type HeadMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *HeadMaxLenRule) Name() string { return "header-max-length" }

// Validate validates HeadMaxLenRule
func (r *HeadMaxLenRule) Validate(msg *message.Commit) (string, bool) {
	return checkMaxLen(r.CheckLen, msg.Header.FullHeader)
}

// SetAndCheckArgument sets the needed argument for the rule
func (r *HeadMaxLenRule) SetAndCheckArgument(arg interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// BodyMaxLenRule to validate max length of body
type BodyMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *BodyMaxLenRule) Name() string { return "body-max-length" }

// Validate validates BodyMaxLenRule
func (r *BodyMaxLenRule) Validate(msg *message.Commit) (string, bool) {
	return checkMaxLen(r.CheckLen, msg.Body)
}

// SetAndCheckArgument sets the needed argument for the rule
func (r *BodyMaxLenRule) SetAndCheckArgument(arg interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// FooterMaxLenRule to validate max length of footer
type FooterMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *FooterMaxLenRule) Name() string { return "footer-max-length" }

// Validate validates FooterMaxLenRule
func (r *FooterMaxLenRule) Validate(msg *message.Commit) (string, bool) {
	return checkMaxLen(r.CheckLen, msg.Footer.FullFooter)
}

// SetAndCheckArgument sets the needed argument for the rule
func (r *FooterMaxLenRule) SetAndCheckArgument(arg interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// TypeMaxLenRule to validate max length of type
type TypeMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *TypeMaxLenRule) Name() string { return "type-max-length" }

// Validate validates TypeMaxLenRule
func (r *TypeMaxLenRule) Validate(msg *message.Commit) (string, bool) {
	return checkMaxLen(r.CheckLen, msg.Header.Type)
}

// SetAndCheckArgument sets the needed argument for the rule
func (r *TypeMaxLenRule) SetAndCheckArgument(arg interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// ScopeMaxLenRule to validate max length of type
type ScopeMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *ScopeMaxLenRule) Name() string { return "scope-max-length" }

// Validate validates ScopeMaxLenRule
func (r *ScopeMaxLenRule) Validate(msg *message.Commit) (string, bool) {
	return checkMaxLen(r.CheckLen, msg.Header.Scope)
}

// SetAndCheckArgument sets the needed argument for the rule
func (r *ScopeMaxLenRule) SetAndCheckArgument(arg interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// DescriptionMaxLenRule to validate max length of type
type DescriptionMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *DescriptionMaxLenRule) Name() string { return "description-max-length" }

// Validate validates DescriptionMaxLenRule
func (r *DescriptionMaxLenRule) Validate(msg *message.Commit) (string, bool) {
	return checkMaxLen(r.CheckLen, msg.Header.Description)
}

// SetAndCheckArgument sets the needed argument for the rule
func (r *DescriptionMaxLenRule) SetAndCheckArgument(arg interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

func checkMaxLen(checkLen int, toCheck string) (string, bool) {
	if checkLen < 0 {
		return "", true
	}
	actualLen := len(toCheck)
	if actualLen > checkLen {
		errMsg := fmt.Sprintf("length is %d, should have less than %d chars", actualLen, checkLen)
		return errMsg, false
	}
	return "", true
}
