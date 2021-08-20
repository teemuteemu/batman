package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/teemuteemu/batman/pkg/env"
	"github.com/teemuteemu/batman/pkg/files"
	"github.com/teemuteemu/batman/pkg/output"
	"github.com/teemuteemu/batman/pkg/runner"
)

const (
	APP_NAME = "batman"
	VERSION  = "0.1"
)

var cli struct {
	Env string `short:"e" help:"Custom .env file location"`
	Run struct {
		Collection string   `arg:"" required:"" help:"Collection YAML"`
		Requests   []string `arg:"" required:"" name:"requests" help:"Requests to execute"`
	} `cmd:"" help:"Run one or more requests from the given collection"`
	Script struct {
		Script string `arg:"" required:"" help:"Script YAML"`
	} `cmd:"" help:"Run a script YAML"`
	Output struct {
		Collection string   `arg:"" required:"" help:"Collection YAML"`
		Requests   []string `arg:"" optional:"" name:"requests" help:"Requests to execute"`
		Format     string   `short:"o" default:"curl" enum:"curl" help:"Output format"`
	} `cmd:"" help:"Output collection or requests"`
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

	case "output <collection> <requests>":
		err := outputRequests(cli.Output.Collection, cli.Env, cli.Output.Requests, cli.Output.Format)
		if err != nil {
			exit(&err)
		}

	case "output <collection>":
		err := outputAllRequests(cli.Output.Collection, cli.Env, cli.Output.Format)
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

func outputRequests(collectionFile string, envFile string, requestNames []string, format string) error {
	c, err := files.GetCollection(collectionFile)
	if err != nil {
		return err
	}

	e, err := env.GetEnv(envFile)
	if err != nil {
		return err
	}

	o := output.New(c, &e)

	err = o.OutputRequests(requestNames, format)
	if err != nil {
		return err
	}

	return nil
}

func outputAllRequests(collectionFile string, envFile string, format string) error {
	c, err := files.GetCollection(collectionFile)
	if err != nil {
		return err
	}

	e, err := env.GetEnv(envFile)
	if err != nil {
		return err
	}

	requestNames := make([]string, 0)

	for _, request := range c.Requests {
		requestNames = append(requestNames, request.Name)
	}

	o := output.New(c, &e)

	err = o.OutputRequests(requestNames, format)
	if err != nil {
		return err
	}

	return nil
}
