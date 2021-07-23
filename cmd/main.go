package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/teemuteemu/batman/pkg/env"
	"github.com/teemuteemu/batman/pkg/files"
	"github.com/teemuteemu/batman/pkg/runner"
)

const (
	APP_NAME = "batman"
	VERSION  = "0.1"
)

var cli struct {
	Env string `short:"e" help:"Custom .env file location"`
	Run struct {
		Collection string   `arg:"" help:"Collection YAML" required:""`
		Requests   []string `arg:"" name:"requests" help:"Requests to execute"`
	} `cmd:"" help:"Run one or more requests from the given collection"`
	Script struct {
		Script string `arg:"" help:"Script YAML" required:""`
	} `cmd:"" help:"Run a script YAML"`
	Version struct{} `cmd:"version" help:"Print version number"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name(APP_NAME),
		kong.Description("Scriptable HTTP client for command line."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	switch ctx.Command() {
	case "run <collection> <requests>":
		err := runRequests(cli.Run.Collection, cli.Env, cli.Run.Requests)
		if err != nil {
			exit(&err)
		}

	case "script <script>":
		err := runScript(cli.Script.Script, cli.Env)
		if err != nil {
			exit(&err)
		}
	case "version":
		fmt.Printf("%s v%s\n", APP_NAME, VERSION)
		exit(nil)
	}
}

func exit(err *error) {
	if err == nil {
		os.Exit(0)
	} else {
		fmt.Println(*err)
		os.Exit(1)
	}
}

func runRequests(collectionFile string, envFile string, requestNames []string) error {
	c, err := files.GetCollection(collectionFile)
	if err != nil {
		return err
	}

	e, err := env.GetEnv(envFile)
	if err != nil {
		return err
	}

	run, err := runner.New(c, &e)
	if err != nil {
		return err
	}

	err = run.RunRequests(requestNames)
	if err != nil {
		return err
	}

	return nil
}

func runScript(scriptFile string, envFile string) error {
	e, err := env.GetEnv(envFile)
	if err != nil {
		return err
	}

	s, err := files.GetScript(scriptFile)
	if err != nil {
		return err
	}

	c, err := files.GetCollection(s.Collection)
	if err != nil {
		return err
	}

	run, err := runner.New(c, &e)
	if err != nil {
		return err
	}

	err = run.RunScript(s)
	if err != nil {
		return err
	}

	return nil
}
