package tests

import (
	"fmt"
	"testing"

	"github.com/1core-dev/go-service/business/data/dbtest"
	"github.com/1core-dev/go-service/pkg/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}
