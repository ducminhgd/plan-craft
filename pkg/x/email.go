package x

import (
	"errors"
	"regexp"
	"strings"
)

// Email validation errors based on RFC 3696
var (
	ErrEmailEmpty             = errors.New("email address is empty")
	ErrEmailTooLong           = errors.New("email address exceeds 254 characters (RFC 5321 limit)")
	ErrEmailMissingAtSign     = errors.New("email address must contain exactly one @ sign")
	ErrEmailLocalPartTooLong  = errors.New("email local part exceeds 64 characters (RFC 3696)")
	ErrEmailDomainPartTooLong = errors.New("email domain part exceeds 255 characters (RFC 3696)")
	ErrEmailLocalPartInvalid  = errors.New("email local part contains invalid characters or format")
	ErrEmailDomainPartInvalid = errors.New("email domain part is invalid")
	ErrEmailInvalidFormat     = errors.New("email address has invalid format")
)

// RFC 3696 compliance constants
const (
	// MaxEmailLength is the maximum total length of an email address per RFC 5321
	MaxEmailLength = 254
	// MaxLocalPartLength is the maximum length of the local part (before @) per RFC 3696
	MaxLocalPartLength = 64
	// MaxDomainPartLength is the maximum length of the domain part (after @) per RFC 3696
	MaxDomainPartLength = 255
)

var (
	// localPartRegex validates the local part (before @) of an email address
	// Allows: alphanumeric, dots, hyphens, underscores, plus signs, and other RFC 3696 allowed characters
	// Dots cannot be consecutive, start, or end the local part
	localPartRegex = regexp.MustCompile(`^[a-zA-Z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(\.[a-zA-Z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*$`)

	// domainPartRegex validates the domain part (after @) of an email address
	// Must follow DNS naming conventions (LDH rule: Letters, Digits, Hyphens)
	// Each label must start/end with alphanumeric, hyphens only in middle
	// Must have at least one dot with valid TLD
	domainPartRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)
)

// IsValidEmail checks if the given string is a valid email address according to RFC 3696
// Returns true if valid, false otherwise
func IsValidEmail(email string) bool {
	err := ValidateEmail(email)
	return err == nil
}

// ValidateEmail performs comprehensive email validation according to RFC 3696
// Returns a specific error describing the validation failure, or nil if valid
func ValidateEmail(email string) error {
	// Check if empty
	if email == "" {
		return ErrEmailEmpty
	}

	// Trim whitespace
	email = strings.TrimSpace(email)

	// Split into local and domain parts
	atIndex := strings.LastIndex(email, "@")
	if atIndex == -1 || strings.Count(email, "@") != 1 {
		return ErrEmailMissingAtSign
	}

	localPart := email[:atIndex]
	domainPart := email[atIndex+1:]

	// Validate local part length (RFC 3696: max 64 characters)
	if len(localPart) > MaxLocalPartLength {
		return ErrEmailLocalPartTooLong
	}

	// Validate domain part length (RFC 3696: max 255 characters)
	if len(domainPart) > MaxDomainPartLength {
		return ErrEmailDomainPartTooLong
	}

	// Check total length (RFC 5321 limit is 254 characters)
	if len(email) > MaxEmailLength {
		return ErrEmailTooLong
	}

	// Validate local part format
	// Cannot start or end with a dot, no consecutive dots
	if localPart == "" || strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") || strings.Contains(localPart, "..") {
		return ErrEmailLocalPartInvalid
	}

	if !localPartRegex.MatchString(localPart) {
		return ErrEmailLocalPartInvalid
	}

	// Validate domain part format
	if domainPart == "" || !domainPartRegex.MatchString(domainPart) {
		return ErrEmailDomainPartInvalid
	}

	// Additional domain validation: check for consecutive dots, leading/trailing dots or hyphens
	if strings.Contains(domainPart, "..") || strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") ||
		strings.HasPrefix(domainPart, "-") || strings.HasSuffix(domainPart, "-") {
		return ErrEmailDomainPartInvalid
	}

	return nil
}
