package masker_test

import (
    "context"
    "testing"

    "reverseproxy/internal/masker"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCreditCardMasker_Mask(t *testing.T) {
    testCases := map[string]struct {
        input    []byte
        expected []byte
    }{
        "ValidCreditNoSpaces": {
            input:    []byte(" 5105105105105100 "),
            expected: []byte(" **************** "),
        },
        "ValidCreditCardDashes": {
            input:    []byte("4012-8888-8888-1881"),
            expected: []byte("****-****-****-****"),
        },
        "ValidCreditCardSpaces": {
            input:    []byte("4012 8888 8888 1881"),
            expected: []byte("**** **** **** ****"),
        },
        "InvalidCreditCardLunh": {
            input:    []byte("1234-5678-9012-3454"),
            expected: []byte("1234-5678-9012-3454"),
        },
        "InvalidCreditCardRandomNumbers": {
            input:    []byte("3215754745"),
            expected: []byte("3215754745"),
        },
        "ValidCreditCardMixedText": {
            input:    []byte("Some text 3530-1113-3330-0000 and more text"),
            expected: []byte("Some text ****-****-****-**** and more text"),
        },
        "ValidMultipleCreditCard": {
            input:    []byte("3530-1113-3330-0000 4012 8888 8888 18815105105105105100 "),
            expected: []byte("****-****-****-**** **** **** **** ******************** "),
        },
    }
    m := masker.NewCreditCardMasker()
    for name, tc := range testCases {
        t.Run(name, func(t *testing.T) {
            actual, err := m.Mask(context.TODO(), tc.input)
            require.NoError(t, err)
            assert.Equal(t, tc.expected, actual)
        })
    }
}

func TestNewCreditCardMasker_Name(t *testing.T) {
    m := masker.NewCreditCardMasker()
    if m.Name() != "Credit Card Masker" {
        t.Errorf("Expected: %s, Got: %s", "credit_card", m.Name())
    }
}
