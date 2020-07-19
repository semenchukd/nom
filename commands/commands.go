package commands

import (
	"errors"
	"fmt"
	"path"
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
	defaultContainers := []string{"core", "dhcp", "data", "dns"}
	switch operation {
	case "start":
		execShell(fmt.Sprintf("docker-compose start %s", strings.Join(defaultContainers, " ")), nil)
		watchCloser := initWatchCloser(len(defaultContainers))
		execShell(`watch docker ps --format \"table {{.ID}}\\t{{.Names}}\\t{{.Status}}\"`, watchCloser)
	case "stop":
		execShell(fmt.Sprintf("docker-compose stop %s", strings.Join(defaultContainers, " ")), nil)
	}
	return nil
}

var Build = func(c *cli.Context) error {
	target := c.Args().First()
	if len(target) == 0 {
		return errors.New("target argument required")
	}

	_, d := path.Split(strings.TrimRight(execAndReturn("pwd"), "/"))
	d = strings.TrimSpace(d)
	var cd string
	switch d {
	case "platform":
		cd = "golang/"
	case "golang":
	default:
		return errors.New("wrong current location")
	}

	switch target {
	case "gatewayd":
		cmd1 := fmt.Sprintf("GOOS=linux make -C %sgatewayd build", cd)
		cmd := fmt.Sprintf("%s && docker cp %sgatewayd/build_output/gatewayd platform_core_1:/usr/local/bin/ && docker-compose exec core sv restart gatewayd", cmd1, cd)
		execShell(cmd, nil)
	case "nexusd":
		cmd1 := fmt.Sprintf("GOOS=linux make -C %snexusd nexusd", cd)
		cmd := fmt.Sprintf("%s && docker cp %snexusd/build_output/nexusd platform_core_1:/usr/local/bin/ && docker-compose exec core sv restart nexusd", cmd1, cd)
		execShell(cmd, nil)
	case "keadatad":
		cmd1 := fmt.Sprintf("GOOS=linux make -C %skeadatad build", cd)
		cmd := fmt.Sprintf("%s && docker cp %skeadatad/build_output/keadatad platform_dhcp_1:/usr/local/bin/ && docker-compose exec dhcp sv restart keadatad", cmd1, cd)
		execShell(cmd, nil)
	default:
		return errors.New("unknown target")
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
