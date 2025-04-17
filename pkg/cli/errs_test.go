package cli

import (
	"testing"
)

func TestErrNotMatched(t *testing.T) {
	assert := func(condition bool, msg string) {
		if !condition {
			t.Errorf(msg)
		}
	}

	assert(ErrNotMatched.Error() == "not matched",
		"expected 'not matched', got '"+ErrNotMatched.Error()+"'")
}

func TestDuplicateCommandError(t *testing.T) {
	assert := func(condition bool, msg string) {
		if !condition {
			t.Errorf(msg)
		}
	}

	err := DuplicateCommandError("test-command")
	expected := "duplicate command: test-command"
	assert(err.Error() == expected,
		"expected '"+expected+"', got '"+err.Error()+"'")
}

func TestMixOfArgumentAndCommandError(t *testing.T) {
	assert := func(condition bool, msg string) {
		if !condition {
			t.Errorf(msg)
		}
	}

	err := MixOfArgumentAndCommandError("test-argument")
	expected := "mix of argument and command: test-argument"
	assert(err.Error() == expected,
		"expected '"+expected+"', got '"+err.Error()+"'")
}
