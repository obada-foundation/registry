package testutil

import (
	"fmt"
	"time"
)

// StartDB starts a immudb instance.
func StartDB() (*Container, error) {
	image := "codenotary/immudb:1.4.1-full"
	port := "3322"
	args := []string{}

	c, err := StartContainer(image, port, args...)
	if err != nil {
		return nil, fmt.Errorf("starting container: %w", err)
	}

	time.Sleep(1 * time.Second)

	fmt.Printf("Image:       %s\n", image)
	fmt.Printf("ContainerID: %s\n", c.ID)
	fmt.Printf("Host:        %s\n", c.Host)

	return c, nil
}

// StopDB stops a running immudb instance.
func StopDB(c *Container) {
	// nolint
	StopContainer(c.ID)
	fmt.Println("Stopped:", c.ID)
}
