package unit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloTom(t *testing.T) {
	output := HelloTom()

	expectOutput := "Tom"
	if output != expectOutput {
		t.Errorf("Expected %s do not match actual %s", expectOutput, output)
	}
}
func TestJudgePassLine(t *testing.T) {
	isPass := JudgePassline(70)
	assert.Equal(t, true, isPass)

}
