package masker

import (
	"context"
	"regexp"
	"strconv"
	"strings"
)

const creditCardBasePattern = `(?:\d[ -]*?){13,16}`

var creditCardBaseRegexp = regexp.MustCompile(creditCardBasePattern)

// CreditCardMasker credit card masker
type CreditCardMasker struct{}

// NewCreditCardMasker creates a cc masker
func NewCreditCardMasker() *CreditCardMasker {
	return &CreditCardMasker{}
}

// Mask every cc found in text, replace every cc digit with '*'
func (ccm *CreditCardMasker) Mask(ctx context.Context, text []byte) ([]byte, error) {
	// Replace the matched credit card numbers with masked values
	maskedText := creditCardBaseRegexp.ReplaceAllStringFunc(string(text), func(cc string) string {
		// Replace - and space characters with empty strings to check luhn
		var onlyDigits strings.Builder
		for _, r := range cc {
			if r >= '0' && r <= '9' {
				onlyDigits.WriteRune(r)
			}
		}
		if luhn(onlyDigits.String()) {
			return maskCreditCard(cc)
		}
		return cc
	})

	return []byte(maskedText), nil
}

// Name ...
func (ccm *CreditCardMasker) Name() string {
	return "Credit Card Masker"
}

func maskCreditCard(cc string) string {
	maskedCC := ""
	for i := 0; i < len(cc); i++ {
		if cc[i] < '0' || cc[i] > '9' {
			maskedCC += string(cc[i])
			continue
		}
		maskedCC += "*"
	}
	return maskedCC
}

func luhn(s string) bool {
	var sum int
	var alternate bool
	numberLen := len(s)
	if numberLen < 13 || numberLen > 19 {
		return false
	}
	for i := numberLen - 1; i > -1; i-- {
		mod, _ := strconv.Atoi(string(s[i]))
		if alternate {
			mod *= 2
			if mod > 9 {
				mod = (mod % 10) + 1
			}
		}
		alternate = !alternate
		sum += mod
	}
	return sum%10 == 0
}
