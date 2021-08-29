package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/teemuteemu/batman/pkg/client"
	"github.com/teemuteemu/batman/pkg/env"
	"github.com/teemuteemu/batman/pkg/files"
	"github.com/teemuteemu/batman/pkg/output"
	"github.com/teemuteemu/batman/pkg/runner"
)

const (
	AppName        = "batman"
	AppDescription = "Scriptable HTTP client for command line."
	Version        = "0.1"
)

var cli struct {
	Env string `short:"e" help:"Custom .env file location"`
	Run struct {
		Collection string   `arg:"" required:"" help:"Collection YAML"`
		Requests   []string `arg:"" required:"" name:"requests" help:"Requests to execute"`
		Output     string   `short:"o" defalt:"console" default:"console" enum:"console, json, yaml" help:"Output format"`
	} `cmd:"" help:"Run one or more requests from the given collection"`
	Script struct {
		Script string `arg:"" required:"" help:"Script YAML"`
	} `cmd:"" help:"Run a script YAML"`
	Print struct {
		Collection string   `arg:"" required:"" help:"Collection YAML"`
		Requests   []string `arg:"" optional:"" name:"requests" help:"Requests to execute"`
		Output     string   `short:"o" default:"curl" enum:"curl" help:"Output format"`
	} `cmd:"" help:"Output collection or requests"`
	Version struct{} `cmd:"version" help:"Print version number"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name(AppName),
		kong.Description(AppDescription),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	switch ctx.Command() {
	case "run <collection> <requests>":
		err := processRequests(cli.Env, cli.Run.Collection, cli.Run.Requests, cli.Run.Output, true)
		if err != nil {
			exit(&err)
		}

	case "script <script>":
		err := runScript(cli.Env, cli.Script.Script)
		if err != nil {
			exit(&err)
		}

	case "print <collection> <requests>":
		err := processRequests(cli.Env, cli.Print.Collection, cli.Print.Requests, cli.Print.Output, false)
		if err != nil {
			exit(&err)
		}

		/*
			case "print <collection>":
				err := outputAllRequests(cli.Env, cli.Print.Collection, cli.Print.Output)
				if err != nil {
					exit(&err)
				}
		*/

	case "version":
		fmt.Printf("%s v%s\n", AppName, Version)
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

func processRequests(envFile string, collectionFile string, requestNames []string, format string, executeRequests bool) error {
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

	formatters := client.Formaters{
		"curl":    output.CurlFormatter{},
		"yaml":    output.YAMLFormatter{},
		"json":    output.JSONFormatter{},
		"console": output.ConsoleFormatter{},
	}

	if formatter, ok := formatters[format]; ok {
		err = run.ProcessRequests(requestNames, formatter, executeRequests)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf(`Could not found formatter "%s"`, format)
}

func runScript(envFile string, scriptFile string) error {
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
