package zerolog_context_issue

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/rs/zerolog"
)

func TestZerologChildContext(t *testing.T) {
	expected := []struct {
		Level     string           `json:"level"`
		Parent    optional[string] `json:"parent"`
		ParentKey optional[string] `json:"parent_key"`
		Child     optional[string] `json:"child"`
		Message   string           `json:"message"`
	}{
		// {"level":"info","parent":"parent","parent_key":"parent_key","message":"parent"}
		{
			Level:     "info",
			Parent:    optional[string]{"parent", true},
			ParentKey: optional[string]{"parent_key", true},
			Message:   "parent",
		},
		// {"level":"info","parent":"parent","child":"child","message":"child"}
		{
			Level:   "info",
			Parent:  optional[string]{"parent", true},
			Child:   optional[string]{"child", true},
			Message: "child",
		},
	}

	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	parentContext := logger.With().Str("parent", "parent")
	childContext := parentContext.Str("child", "child")

	parentLogger := parentContext.Logger()
	childLogger := childContext.Logger()

	parentLogger.Info().Str("parent_key", "parent_key").Msg("parent")
	childLogger.Info().Msg("child")

	assertLines(t, expected, buf)
}

func TestZerologChildrenContexts(t *testing.T) {
	expected := []struct {
		Level   string           `json:"level"`
		Parent  optional[string] `json:"parent"`
		Child1  optional[string] `json:"child_1"`
		Child2  optional[string] `json:"child_2"`
		Message string           `json:"message"`
	}{
		// {"level":"info","parent":"parent","message":"parent"}
		{
			Level:   "info",
			Parent:  optional[string]{"parent", true},
			Message: "parent",
		},
		// {"level":"info","parent":"parent","child_1":"child_1","message":"child_1"}
		{
			Level:   "info",
			Parent:  optional[string]{"parent", true},
			Child1:  optional[string]{"child_1", true},
			Message: "child_1",
		},
		// {"level":"info","parent":"parent","child_2":"child_2","message":"child_2"}
		{
			Level:   "info",
			Parent:  optional[string]{"parent", true},
			Child2:  optional[string]{"child_2", true},
			Message: "child_2",
		},
	}

	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	parentContext := logger.With().Str("parent", "parent")
	child1Context := parentContext.Str("child_1", "child_1")
	child2Context := parentContext.Str("child_2", "child_2")

	parentLogger := parentContext.Logger()
	child1Logger := child1Context.Logger()
	child2Logger := child2Context.Logger()

	parentLogger.Info().Msg("parent")
	child1Logger.Info().Msg("child_1")
	child2Logger.Info().Msg("child_2")

	assertLines(t, expected, buf)
}

type optional[T any] struct {
	value    T
	hasValue bool
}

func (m *optional[T]) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &m.value); err != nil {
		return err
	}
	m.hasValue = true
	return nil
}

func assertEqual[T comparable](t *testing.T, expected *T, actual string) {
	var v T
	err := json.Unmarshal([]byte(actual), &v)
	if err != nil {
		t.Error(err)
		return
	}
	if v != *expected {
		t.Errorf("expected: %+v, actual: %+v", expected, &v)
	}
}

func assertLines[T comparable](t *testing.T, expected []T, buf bytes.Buffer) {
	for _, expectedLine := range expected {
		line, err := buf.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			t.Fatal(err)
		}
		assertEqual(t, &expectedLine, line)
	}
}
