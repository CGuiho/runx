package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	jsonFmt bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "runx",
	Short: "A language-agnostic, documented command catalog and local CLI executor",
	Long: `RunX is an open-source, schema-backed CLI task runner and command catalog
executor that operates over runx.yaml manifest files.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if helpTree {
			RenderHelpTree(cmd.OutOrStdout(), helpTreeDepth)
			if exitOnHelpTree {
				os.Exit(0)
			}
			return nil
		}
		if helpDocs {
			RenderHelpDocs(cmd.OutOrStdout(), cmd)
			if exitOnHelpTree {
				os.Exit(0)
			}
			return nil
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		renderWelcome()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./runx.yaml)")
	RootCmd.PersistentFlags().BoolVar(&jsonFmt, "json", false, "output in JSON format")

	RootCmd.PersistentFlags().BoolVar(&helpTree, "help-tree", false, "Show this command hierarchy.")
	RootCmd.PersistentFlags().IntVar(&helpTreeDepth, "help-tree-depth", 0, "Limit help-tree recursion depth.")
	RootCmd.PersistentFlags().BoolVar(&helpDocs, "help-docs", false, "Emit Markdown documentation for this command.")

	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("json", RootCmd.PersistentFlags().Lookup("json"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("runx")
	}

	viper.SetEnvPrefix("RUNX")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}

func renderWelcome() {
	fmt.Printf("Hello %s - runx v0.8.0-dev (%s/%s)\n", getPlatformLabel(), runtime.GOOS, runtime.GOARCH)
	fmt.Println("Type 'runx --help' for available commands and usage.")
}

func getPlatformLabel() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	default:
		return runtime.GOOS
	}
}
