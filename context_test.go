package zerolog_context_issue

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
)

func TestZerologChildContext(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	parentContext := logger.With().Bool("parent", true)
	childContext := parentContext.Bool("child", true)

	parentLogger := parentContext.Logger()
	childLogger := childContext.Logger()

	parentLogger.Info().Bool("parent_key", true).Msg("parent")
	childLogger.Info().Msg("child")

	parentLoggerOutput, err := buf.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	childLoggerOutput, err := buf.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}

	// parentLoggerOutput: {"level":"info","parent":true,"parent_key":true,"message":"parent"}
	assertEqual(t, &struct {
		Level     string `json:"level"`
		Parent    bool   `json:"parent"`
		ParentKey bool   `json:"parent_key"`
		Child     bool   `json:"child"`
		Message   string `json:"message"`
	}{
		Level:     "info",
		Parent:    true,
		ParentKey: true,
		Child:     false,
		Message:   "parent",
	}, parentLoggerOutput)

	// childLoggerOutput: {"level":"info","parent":true,"child":true,"message":"child"}
	assertEqual(t, &struct {
		Level     string `json:"level"`
		Parent    bool   `json:"parent"`
		ParentKey bool   `json:"parent_key"`
		Child     bool   `json:"child"`
		Message   string `json:"message"`
	}{
		Level:     "info",
		Parent:    true,
		ParentKey: false,
		Child:     true,
		Message:   "child",
	}, childLoggerOutput)
}

func TestZerologChildrenContexts(t *testing.T) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)

	parentContext := logger.With().Bool("parent", true)
	child1Context := parentContext.Bool("child_1", true)
	child2Context := parentContext.Bool("child_2", true)

	parentLogger := parentContext.Logger()
	child1Logger := child1Context.Logger()
	child2Logger := child2Context.Logger()

	parentLogger.Info().Msg("parent")
	child1Logger.Info().Msg("child_1")
	child2Logger.Info().Msg("child_2")

	parentLoggerOutput, err := buf.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	child1LoggerOutput, err := buf.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	child2LoggerOutput, err := buf.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}

	// parentLoggerOutput: {"level":"info","parent":true,"message":"parent"}
	assertEqual(t, &struct {
		Level   string `json:"level"`
		Parent  bool   `json:"parent"`
		Child1  bool   `json:"child_1"`
		Child2  bool   `json:"child_2"`
		Message string `json:"message"`
	}{
		Level:   "info",
		Parent:  true,
		Message: "parent",
	}, parentLoggerOutput)

	// child1LoggerOutput: {"level":"info","parent":true,"child_1":true,"message":"child_1"}
	assertEqual(t, &struct {
		Level   string `json:"level"`
		Parent  bool   `json:"parent"`
		Child1  bool   `json:"child_1"`
		Child2  bool   `json:"child_2"`
		Message string `json:"message"`
	}{
		Level:   "info",
		Parent:  true,
		Child1:  true,
		Message: "child_1",
	}, child1LoggerOutput)

	// child2LoggerOutput: {"level":"info","parent":true,"child_2":true,"message":"child_2"}
	assertEqual(t, &struct {
		Level   string `json:"level"`
		Parent  bool   `json:"parent"`
		Child1  bool   `json:"child_1"`
		Child2  bool   `json:"child_2"`
		Message string `json:"message"`
	}{
		Level:   "info",
		Parent:  true,
		Child2:  true,
		Message: "child_2",
	}, child2LoggerOutput)
}

func assertEqual[T comparable](t *testing.T, expected *T, actual string) {
	var v T
	err := json.Unmarshal([]byte(actual), &v)
	if err != nil {
		t.Fatal(err)
	}
	if v != *expected {
		t.Errorf("expected: %+v, actual: %+v", expected, &v)
	}
}
