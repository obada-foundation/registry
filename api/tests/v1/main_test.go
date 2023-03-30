package tests

import (
	"fmt"
	"testing"

	"github.com/obada-foundation/registry/testutil"
)

// nolint
var c *testutil.Container

type JSON map[string]interface{}

func TestMain(m *testing.M) {
	var err error

	c, err = testutil.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer testutil.StopDB(c)

	m.Run()
}
