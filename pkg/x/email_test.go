package x

import (
	"strings"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		// Valid email addresses
		{"Valid: simple email", "user@example.com", true},
		{"Valid: with subdomain", "user@mail.example.com", true},
		{"Valid: with plus sign", "user+tag@example.com", true},
		{"Valid: with dots", "first.last@example.com", true},
		{"Valid: with underscore", "user_name@example.com", true},
		{"Valid: with hyphen in local", "user-name@example.com", true},
		{"Valid: with hyphen in domain", "user@my-company.com", true},
		{"Valid: with numbers", "user123@example123.com", true},
		{"Valid: RFC 3696 example 1", "user+mailbox@example.com", true},
		{"Valid: RFC 3696 example 2", "customer/department=shipping@example.com", true},
		{"Valid: RFC 3696 example 3", "$A12345@example.com", true},
		{"Valid: RFC 3696 example 4", "!def!xyz%abc@example.com", true},
		{"Valid: with special chars", "test!#$%&'*+/=?^_`{|}~@example.com", true},
		{"Valid: long TLD", "user@example.museum", true},
		{"Valid: numeric TLD", "user@example.co", true},

		// Invalid email addresses
		{"Invalid: empty string", "", false},
		{"Invalid: no @ sign", "userexample.com", false},
		{"Invalid: multiple @ signs", "user@@example.com", false},
		{"Invalid: no domain", "user@", false},
		{"Invalid: no local part", "@example.com", false},
		{"Invalid: no TLD", "user@example", false},
		{"Invalid: starts with dot", ".user@example.com", false},
		{"Invalid: ends with dot", "user.@example.com", false},
		{"Invalid: consecutive dots in local", "user..name@example.com", false},
		{"Invalid: consecutive dots in domain", "user@example..com", false},
		{"Invalid: domain starts with hyphen", "user@-example.com", false},
		{"Invalid: domain ends with hyphen", "user@example-.com", false},
		{"Invalid: domain starts with dot", "user@.example.com", false},
		{"Invalid: domain ends with dot", "user@example.com.", false},
		{"Invalid: space in email", "user name@example.com", false},
		{"Invalid: missing domain part", "user@.com", false},
		{"Invalid: TLD too short", "user@example.c", false},
		{"Invalid: invalid characters", "user@exa mple.com", false},
		{"Invalid: double @", "user@domain@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.want {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantError error
	}{
		// Valid emails should return nil
		{"Valid email", "user@example.com", nil},
		{"Valid with subdomain", "user@mail.example.com", nil},

		// Specific error cases
		{"Empty email", "", ErrEmailEmpty},
		{"No @ sign", "userexample.com", ErrEmailMissingAtSign},
		{"Multiple @ signs", "user@@example.com", ErrEmailMissingAtSign},
		{"Local part starts with dot", ".user@example.com", ErrEmailLocalPartInvalid},
		{"Local part ends with dot", "user.@example.com", ErrEmailLocalPartInvalid},
		{"Consecutive dots in local", "user..name@example.com", ErrEmailLocalPartInvalid},
		{"Domain starts with dot", "user@.example.com", ErrEmailDomainPartInvalid},
		{"Domain ends with dot", "user@example.com.", ErrEmailDomainPartInvalid},
		{"Consecutive dots in domain", "user@example..com", ErrEmailDomainPartInvalid},
		{"Domain starts with hyphen", "user@-example.com", ErrEmailDomainPartInvalid},
		{"Domain ends with hyphen", "user@example-.com", ErrEmailDomainPartInvalid},
		{"No domain part", "user@", ErrEmailDomainPartInvalid},
		{"No local part", "@example.com", ErrEmailLocalPartInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if err != tt.wantError {
				t.Errorf("ValidateEmail(%q) error = %v, want %v", tt.email, err, tt.wantError)
			}
		})
	}
}

func TestEmailLengthLimits(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantError error
	}{
		{
			name:      "Local part exactly 64 chars",
			email:     strings.Repeat("a", 64) + "@example.com",
			wantError: nil,
		},
		{
			name:      "Local part exceeds 64 chars",
			email:     strings.Repeat("a", 65) + "@example.com",
			wantError: ErrEmailLocalPartTooLong,
		},
		{
			name:      "Domain part long but valid (multiple labels)",
			email:     "user@abc.def.ghi.jkl.mno.pqr.stu.vwx.example.com",
			wantError: nil,
		},
		{
			name: "Domain part exceeds 255 chars",
			// Create a domain with 256 characters (exceeds limit)
			// Using short local part to avoid hitting total length limit first
			email:     "u@" + strings.Repeat("abcdefgh.", 28) + "example.com",
			wantError: ErrEmailDomainPartTooLong,
		},
		{
			name:      "Total length near limit",
			email:     strings.Repeat("a", 64) + "@subdomain.example.com",
			wantError: nil,
		},
		{
			name:      "Total length exceeds limit (255 chars)",
			email:     strings.Repeat("a", 64) + "@" + strings.Repeat("ab", 93) + ".com",
			wantError: ErrEmailTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if err != tt.wantError {
				t.Errorf("ValidateEmail() error = %v, want %v", err, tt.wantError)
			}
		})
	}
}

func TestEmailWithWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"Leading whitespace", "  user@example.com", true},
		{"Trailing whitespace", "user@example.com  ", true},
		{"Both leading and trailing", "  user@example.com  ", true},
		{"Whitespace in middle", "user @example.com", false},
		{"Tab character", "user\t@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.want {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestRFC3696SpecialCharacters(t *testing.T) {
	// Test special characters allowed in local part per RFC 3696
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"Exclamation mark", "test!test@example.com", true},
		{"Hash", "test#test@example.com", true},
		{"Dollar sign", "test$test@example.com", true},
		{"Percent", "test%test@example.com", true},
		{"Ampersand", "test&test@example.com", true},
		{"Apostrophe", "test'test@example.com", true},
		{"Asterisk", "test*test@example.com", true},
		{"Plus", "test+test@example.com", true},
		{"Forward slash", "test/test@example.com", true},
		{"Equals", "test=test@example.com", true},
		{"Question mark", "test?test@example.com", true},
		{"Caret", "test^test@example.com", true},
		{"Underscore", "test_test@example.com", true},
		{"Backtick", "test`test@example.com", true},
		{"Left brace", "test{test@example.com", true},
		{"Pipe", "test|test@example.com", true},
		{"Right brace", "test}test@example.com", true},
		{"Tilde", "test~test@example.com", true},
		{"Hyphen", "test-test@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.want {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func BenchmarkIsValidEmail(b *testing.B) {
	email := "user@example.com"
	for i := 0; i < b.N; i++ {
		IsValidEmail(email)
	}
}

func BenchmarkValidateEmail(b *testing.B) {
	email := "user@example.com"
	for i := 0; i < b.N; i++ {
		ValidateEmail(email)
	}
}
