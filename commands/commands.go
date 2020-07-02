package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var Newenv = func(c *cli.Context) error {
	build := c.Args().First()
	if len(build) == 0 {
		return errors.New("build argument required")
	}

	execShell("make clean-env", nil)
	execShell(fmt.Sprintf("BUILD=%s make localize-standalone-only", build), nil)

	if !c.Bool("noup") {
		defaultUp := []string{"core", "dhcp", "data", "dns"}
		execShell(fmt.Sprintf("docker-compose up -d %s", strings.Join(defaultUp, " ")), nil) // TODO: configurable params

		if c.Bool("nobs") {
			return nil
		}

		var watchCloser termCloser
		go func() {
			i := 200
			for i > 0 {
				res := execAndReturn("docker ps | grep -c '(healthy)'")
				if strings.TrimSpace(res) == strconv.Itoa(len(defaultUp)) {
					watchCloser.Close()
					return
				}
				time.Sleep(1 * time.Second)
				i--
			}
			watchCloser.Close()
		}()

		execShell(`watch docker ps --format \"table {{.ID}}\\t{{.Names}}\\t{{.Status}}\"`, &watchCloser)
		execShell("clear", nil)

		execShell("make bootstrap-api", nil)
	}

	return nil
}
