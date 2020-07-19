package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

// Create and bootstrap new evironment
var Newenv = func(c *cli.Context) error {
	build := c.Args().First()
	if len(build) == 0 {
		return errors.New("build argument is required")
	}

	execShell("make clean-env", nil)
	execShell(fmt.Sprintf("BUILD=%s make localize-standalone-only", build), nil)

	if c.Bool("noup") {
		return nil
	}

	defaultUp := []string{"core", "dhcp", "data", "dns"}
	execShell(fmt.Sprintf("docker-compose up -d %s", strings.Join(defaultUp, " ")), nil) // TODO: configurable params

	if c.Bool("nobs") {
		return nil
	}

	watchCloser := initWatchCloser(len(defaultUp))

	execShell(`watch docker ps --format \"table {{.ID}}\\t{{.Names}}\\t{{.Status}}\"`, watchCloser)
	execShell("make bootstrap-api", nil)

	return nil
}

// docker-compose operations
var DC = func(c *cli.Context) error {
	operation := c.Args().First()
	if len(operation) == 0 {
		return errors.New("operation argument required")
	}
	switch operation {
	case "start":
		defaultStart := []string{"core", "dhcp", "data", "dns"}
		execShell(fmt.Sprintf("docker-compose start %s", strings.Join(defaultStart, " ")), nil)
		watchCloser := initWatchCloser(len(defaultStart))
		execShell(`watch docker ps --format \"table {{.ID}}\\t{{.Names}}\\t{{.Status}}\"`, watchCloser)
	case "stop":
		execShell("docker-compose stop", nil)
	}
	return nil
}

func initWatchCloser(numToWait int) *termCloser {
	var watchCloser termCloser
	go func() {
		i := 200
		for i > 0 {
			res := execAndReturn("docker ps | grep -c '(healthy)'")
			if strings.TrimSpace(res) == strconv.Itoa(numToWait) {
				if err := watchCloser.Close(); err != nil {
					panic(err)
				}
				return
			}
			time.Sleep(1 * time.Second)
			i--
		}
		if err := watchCloser.Close(); err != nil {
			panic(err)
		}
	}()
	return &watchCloser
}
