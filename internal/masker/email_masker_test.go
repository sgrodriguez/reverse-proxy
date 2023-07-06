package masker_test

import (
	"context"
	"testing"

	"reverseproxy/internal/masker"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailMasker_Mask(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected string
	}{
		// Standard Cases
		"Single Email": {
			input:    "Contact us at test@example.com for assistance.",
			expected: "Contact us at ****@example.com for assistance.",
		},
		"Multiple Emails": {
			input:    "Send an email to john@example.com or jane@example.com for more information.",
			expected: "Send an email to ****@example.com or ****@example.com for more information.",
		},
		"Email in Parentheses": {
			input:    "For inquiries, email (support@example.com).",
			expected: "For inquiries, email (****@example.com).",
		},
		"Mixed Text": {
			input:    "Contact john@example.com or visit our website at www.example.com for more details.",
			expected: "Contact ****@example.com or visit our website at www.example.com for more details.",
		},
		"NoEmail": {
			input:    "This is a sample text without any email addresses.",
			expected: "This is a sample text without any email addresses.",
		},
		"EmptyText": {
			input:    "",
			expected: "",
		},
		"OnlyEmail": {
			input:    "user@example.com",
			expected: "****@example.com",
		},
		"InvalidEmail": {
			input:    "This is not a valid email address: test@",
			expected: "This is not a valid email address: test@",
		},
		"EmailWithNumbersAndUnderscore": {
			input:    "Contact 123_asd@gmail.com for assistance.",
			expected: "Contact ****@gmail.com for assistance.",
		},
		// Strange Email Domains
		"EmailWithNon-StandardDomain": {
			input:    "Contact user@example.co.uk for support.",
			expected: "Contact ****@example.co.uk for support.",
		},
		"EmailWithNon-StandardTLD": {
			input:    "Send an email to info@company.xyz for inquiries.",
			expected: "Send an email to ****@company.xyz for inquiries.",
		},
	}
	m := masker.NewEmailMasker()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual, err := m.Mask(context.TODO(), []byte(tc.input))
			require.NoError(t, err)
			assert.Equal(t, tc.expected, string(actual))
		})
	}
}

func TestEmailMasker_Name(t *testing.T) {
	m := masker.NewEmailMasker()
	actual := m.Name()
	expected := "Email Masker"
	if actual != expected {
		t.Errorf("Expected: %s, Got: %s", expected, actual)
	}

}
