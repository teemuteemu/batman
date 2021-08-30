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

var Version = "unset"

const (
	AppName        = "batman"
	AppDescription = "Scriptable HTTP client for command line."

	UsageRun    = "batman run <collection file> [<request>]\n"
	UsagePrint  = "batman print <collection file> [<request>]\n"
	UsageScript = "batman script <script file>\n"
)

var EnvFlag string
var RunOutputFlag string
var PrintOutputFlag string

var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: AppDescription,
	Long:  fmt.Sprintf("%s %s - %s", AppName, Version, AppDescription),
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", AppName, Version)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run requests in collection",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		collectionFile := args[0]
		requests := args[1:]
		return processRequests(EnvFlag, collectionFile, requests, RunOutputFlag, true)
	},
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the request in given format",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		collectionFile := args[0]
		requests := args[1:]
		return processRequests(EnvFlag, collectionFile, requests, PrintOutputFlag, false)
	},
}

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "Run a script file",
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
	runCmd.SetUsageTemplate(UsageRun)
	rootCmd.AddCommand(runCmd)

	printCmd.Flags().StringVarP(&EnvFlag, "env", "e", "", "Custom environment file (.env)")
	printCmd.Flags().StringVarP(&PrintOutputFlag, "output", "o", "curl", "Output format")
	printCmd.SetUsageTemplate(UsagePrint)
	rootCmd.AddCommand(printCmd)

	scriptCmd.Flags().StringVarP(&EnvFlag, "env", "e", "", "Custom environment file (.env)")
	scriptCmd.SetUsageTemplate(UsageScript)
	rootCmd.AddCommand(scriptCmd)
}

func Execute(version string) error {
	Version = version
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
		"curl":     output.CurlFormatter{},
		"yaml":     output.YAMLFormatter{},
		"json":     output.JSONFormatter{},
		"console":  output.ConsoleFormatter{},
		"response": output.ResponseFormatter{},
		"request":  output.RequestFormatter{},
		"mute":     output.MuteFormatter{},
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
