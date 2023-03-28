package tests

import (
	"testing"

	"github.com/obada-foundation/registry/testutil"
)

// nolint
var c *testutil.Container

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	m.Run()
}
