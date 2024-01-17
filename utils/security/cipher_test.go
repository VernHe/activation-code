package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCipherText(t *testing.T) {
	// MD5:7721ed22f6921b7be8a8e81f585b5ab4，数字部分之和是 7+7+2+1+2+2+6+9+2+1+7+8+8+8+1+5+8+5+5+4=98
	value := "1eds23s4w567g89f01"

	t.Run("test total count", func(t *testing.T) {
		text := GenerateCipherText(value, 6, 0)
		assert.Equal(t, 98, getNumberSum(text))
	})

	t.Run("test hour count", func(t *testing.T) {
		text := GenerateCipherText(value, 5, 4)
		assert.Equal(t, 98, getNumberSum(text))
	})

	t.Run("test nomal case", func(t *testing.T) {
		text := GenerateCipherText(value, 5, 3)
		assert.NotEqual(t, 98, getNumberSum(text))
	})
}

func TestGetNumberSum(t *testing.T) {
	assert.Equal(t, 98, getNumberSum("7721ed22f6921b7be8a8e81f585b5ab4"))
}

func TestGenerateString(t *testing.T) {
	s := generateString(98)
	assert.Equal(t, 98, getNumberSum(s))
}
