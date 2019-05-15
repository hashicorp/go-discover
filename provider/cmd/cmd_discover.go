// Package cmd provides node discovery via customizable cmd.
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `cmd:

    provider:    "cmd"
    cmdline:     cmdline to excute to discover node.
	timeout:	 Timeout to wait cmdline to execute in seconds (Default 5).

    The cmdline should output the node addresses line by line through stdout.
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	var res bytes.Buffer
	var err error
	var timeout uint64 = 5

	if args["provider"] != "cmd" {
		return nil, fmt.Errorf("discover-cmd: invalid provider " + args["provider"])
	}
	if _, ok := args["cmdline"]; !ok {
		return nil, fmt.Errorf("discover-cmd: missing arg cmdline")
	}
	if t, ok := args["timeout"]; ok {
		timeout, err = strconv.ParseUint(t, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("discover-cmd: invalid timeout: " + t)
		}
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	cmd := exec.CommandContext(ctx, "sh", "-c", args["cmdline"])

	cmd.Stdout = &res
	err = cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("discover-cmd: excute cmd error: %s", err)
	}

	return strings.Split(res.String(), "\n"), nil
}
