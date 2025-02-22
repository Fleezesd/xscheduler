package main

import (
	"fmt"
	"os"

	"github.com/fleezesd/xscheduler/cmd/xscheduler/app"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
)

func main() {
	if err := runXshedulerCmd(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func runXshedulerCmd() error {
	// WordSepNormalizeFunc changes all flags that contain "_" separators
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	command := app.NewXschedulerCmd()

	logs.InitLogs()
	defer logs.FlushLogs()

	// ParseFlags parses persistent flag tree and local flags.
	err := command.ParseFlags(os.Args[1:])
	if err != nil {
		return fmt.Errorf("%v\n%s", err, command.UsageString())
	}
	return command.Execute()
}
