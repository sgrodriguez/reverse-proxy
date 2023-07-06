package masker

import (
	"context"
	"regexp"
)

const emailPattern = `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`

var emailRegexp = regexp.MustCompile(emailPattern)

// EmailMasker ...
type EmailMasker struct{}

// NewEmailMasker creates a email masker
func NewEmailMasker() *EmailMasker {
	return &EmailMasker{}
}

// Mask every email found in text, left domain intact, replace localPart with '*'
func (em *EmailMasker) Mask(ctx context.Context, text []byte) ([]byte, error) {
	// Replace the matched email addresses with masked values
	maskedText := emailRegexp.ReplaceAllStringFunc(string(text), func(email string) string {
		// Generate the masked email address
		maskedEmail := maskEmail(email)
		return maskedEmail
	})

	return []byte(maskedText), nil
}

// Name ...
func (em *EmailMasker) Name() string {
	return "Email Masker"
}

func maskEmail(email string) string {
	// We let a fix amount of '*' so we don't share any extra info about email length
	return "****@" + getDomainFromEmail(email)
}

func getDomainFromEmail(email string) string {
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}
	if atIndex == -1 || atIndex == len(email)-1 {
		return ""
	}
	return email[atIndex+1:]
}
