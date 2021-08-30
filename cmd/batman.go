package cmd

import (
	"fmt"
	"path"

	"github.com/teemuteemu/batman/pkg/client"
	"github.com/teemuteemu/batman/pkg/env"
	"github.com/teemuteemu/batman/pkg/files"
	"github.com/teemuteemu/batman/pkg/output"
	"github.com/teemuteemu/batman/pkg/runner"

	"github.com/spf13/cobra"
)

const (
	AppName        = "batman"
	AppDescription = "Scriptable HTTP client for command line."
	Version        = "0.0.6"

	HelpRun    = "batman run <collection file> [requests]\n"
	HelpPrint  = "batman print <collection file> [requests]\n"
	HelpScript = "batman script <script file>\n"
)

var EnvFlag string
var RunOutputFlag string
var PrintOutputFlag string

var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: AppDescription,
	Long:  fmt.Sprintf("%s v%s - %s", AppName, Version, AppDescription),
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", AppName, Version)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run requests in collection",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		collectionFile := args[0]
		requests := args[1:]
		return processRequests(EnvFlag, collectionFile, requests, RunOutputFlag, true)
	},
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "prints the request in given format",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		collectionFile := args[0]
		requests := args[1:]
		return processRequests(EnvFlag, collectionFile, requests, PrintOutputFlag, false)
	},
}

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "run a script file",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptFile := args[0]
		return runScript(EnvFlag, scriptFile)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	runCmd.Flags().StringVarP(&EnvFlag, "env", "e", "", "Custom environment file (.env)")
	runCmd.Flags().StringVarP(&RunOutputFlag, "output", "o", "console", "Output format")
	runCmd.SetHelpTemplate(HelpRun)
	rootCmd.AddCommand(runCmd)

	printCmd.Flags().StringVarP(&EnvFlag, "env", "e", "", "Custom environment file (.env)")
	printCmd.Flags().StringVarP(&PrintOutputFlag, "output", "o", "curl", "Output format")
	printCmd.SetHelpTemplate(HelpPrint)
	rootCmd.AddCommand(printCmd)

	scriptCmd.Flags().StringVarP(&EnvFlag, "env", "e", "", "Custom environment file (.env)")
	scriptCmd.SetHelpTemplate(HelpScript)
	rootCmd.AddCommand(scriptCmd)
}

func Execute() error {
	return rootCmd.Execute()
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
			fmt.Println(err)
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

	collectionFile := path.Join(path.Dir(scriptFile), s.Collection)

	c, err := files.GetCollection(collectionFile)
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
