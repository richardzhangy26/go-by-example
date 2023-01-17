package unit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloV2(t *testing.T) {
	output := HelloTom()
	expectOutput := "Tom"
	assert.Equal(t, expectOutput, output)
}
