package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/CGuiho/runx/embed"
	"github.com/CGuiho/runx/pkg/installstate"
	"github.com/CGuiho/runx/pkg/maintenance"
	"github.com/CGuiho/runx/pkg/update"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"buildDate"`
	Target  string `json:"buildTarget"`
}

type Dependencies struct {
	In         io.Reader
	Out        io.Writer
	Err        io.Writer
	HTTPClient *http.Client
	Now        func() time.Time
	Executable func() (string, error)
	Getwd      func() (string, error)
	HomeDir    func() (string, error)
	Spawn      func(string, ...string) error
	IsTerminal func(io.Reader) bool
}

type exitError struct {
	code int
	err  error
}

func (value *exitError) Error() string { return value.err.Error() }
func (value *exitError) Unwrap() error { return value.err }
func withExitCode(code int, err error) error {
	if err == nil {
		return nil
	}
	return &exitError{code: code, err: err}
}

var errDeveloperHelp = errors.New("developer context rendered")

func DefaultDependencies() Dependencies {
	return Dependencies{
		In: os.Stdin, Out: os.Stdout, Err: os.Stderr,
		HTTPClient: &http.Client{Timeout: 15 * time.Second}, Now: time.Now,
		Executable: os.Executable, Getwd: os.Getwd, HomeDir: os.UserHomeDir,
		Spawn: update.SpawnUpdateWorker, IsTerminal: isTerminalReader,
	}
}

func Execute(info BuildInfo) error {
	root := NewRootCommand(DefaultDependencies(), info)
	err := root.Execute()
	if errors.Is(err, errDeveloperHelp) {
		return nil
	}
	return err
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var coded *exitError
	if errors.As(err, &coded) {
		return coded.code
	}
	var child *ChildExitError
	if errors.As(err, &child) {
		return child.Code
	}
	message := err.Error()
	for _, fragment := range []string{"unknown command", "unknown flag", "requires at least", "requires exactly", "accepts ", "required flag"} {
		if strings.Contains(message, fragment) {
			return 2
		}
	}
	return 1
}

func NewRootCommand(deps Dependencies, info BuildInfo) *cobra.Command {
	deps = normalizeDependencies(deps)
	if strings.TrimSpace(info.Version) == "" {
		info.Version = "dev"
	}
	if strings.TrimSpace(info.Target) == "" {
		info.Target = "development"
	}
	var helpTree bool
	var helpTreeDepth int
	var helpTreeGlobalFlags bool
	var helpDocs bool
	var showVersion bool

	root := &cobra.Command{
		Use: "runx", Short: "A language-agnostic, documented command catalog and local CLI executor",
		Long:         "RunX is a strict YAML command catalog and local CLI executor.",
		SilenceUsage: true, SilenceErrors: true, Args: cobra.NoArgs,
		PersistentPreRunE: func(command *cobra.Command, _ []string) error {
			if command.Flags().Changed("help-tree-depth") && helpTreeDepth < 1 {
				return withExitCode(2, fmt.Errorf("--help-tree-depth must be a positive integer"))
			}
			if showVersion {
				fmt.Fprintln(command.OutOrStdout(), info.Version)
				return errDeveloperHelp
			}
			if helpTree || command.Flags().Changed("help-tree-depth") {
				fmt.Fprint(command.OutOrStdout(), RenderCommandTree(command, helpTreeDepth, helpTreeGlobalFlags))
				return errDeveloperHelp
			}
			if helpDocs {
				markdown, err := RenderHelpDocs(command)
				if err != nil {
					return err
				}
				fmt.Fprint(command.OutOrStdout(), markdown)
				return errDeveloperHelp
			}
			scheduleLifecycle(command, deps, info)
			return nil
		},
		RunE: func(command *cobra.Command, _ []string) error {
			if os.Getenv("RUNX_DISABLE_AGENT_MAINTENANCE_WORKER") != "1" {
				skill, err := bundledSkill()
				if err != nil {
					return withExitCode(5, err)
				}
				home, err := deps.HomeDir()
				if err != nil {
					return withExitCode(5, err)
				}
				cwd, err := deps.Getwd()
				if err != nil {
					return withExitCode(5, err)
				}
				if _, err := maintenance.MaintainAgentIntegration(cwd, home, skill); err != nil {
					return withExitCode(5, err)
				}
			}
			fmt.Fprintf(command.OutOrStdout(), "Hello Windows - runx v%s\n", info.Version)
			if notice, ok := update.ReadCachedUpdateNotice(update.GetDefaultCachePath(), info.Version); ok {
				fmt.Fprintln(command.OutOrStdout(), notice)
			}
			return nil
		},
	}
	root.SetIn(deps.In)
	root.SetOut(deps.Out)
	root.SetErr(deps.Err)
	root.CompletionOptions.DisableDefaultCmd = true
	root.SetHelpCommand(&cobra.Command{Hidden: true})
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return withExitCode(2, err) })
	root.Flags().BoolVarP(&showVersion, "version", "v", false, "Show the RunX version.")
	root.PersistentFlags().BoolVar(&helpTree, "help-tree", false, "Show this command hierarchy.")
	root.PersistentFlags().IntVar(&helpTreeDepth, "help-tree-depth", 0, "Limit help-tree recursion depth.")
	root.PersistentFlags().BoolVar(&helpTreeGlobalFlags, "help-tree-global-flags", false, "Repeat inherited global flags under every descendant in the tree.")
	root.PersistentFlags().BoolVar(&helpDocs, "help-docs", false, "Emit Markdown documentation for this command.")

	root.AddCommand(
		newListCommand(deps), newDescribeCommand(deps), newRevealCommand(deps), newRunCommand(deps), newCheckCommand(deps), newInitCommand(deps, info),
		newAgentCommand(deps), newUpgradeCommand(deps, info), newUninstallCommand(deps), newUpdateWorkerCommand(deps, info),
		newMaintenanceWorkerCommand(deps), newSelfTestCommand(deps, info),
	)
	configureDeveloperContext(root, func(command *cobra.Command) bool {
		return showVersion || helpTree || helpDocs || command.Flags().Changed("help-tree-depth")
	})
	return root
}

func normalizeDependencies(deps Dependencies) Dependencies {
	defaults := DefaultDependencies()
	if deps.In == nil {
		deps.In = defaults.In
	}
	if deps.Out == nil {
		deps.Out = defaults.Out
	}
	if deps.Err == nil {
		deps.Err = defaults.Err
	}
	if deps.HTTPClient == nil {
		deps.HTTPClient = defaults.HTTPClient
	}
	if deps.Now == nil {
		deps.Now = defaults.Now
	}
	if deps.Executable == nil {
		deps.Executable = defaults.Executable
	}
	if deps.Getwd == nil {
		deps.Getwd = defaults.Getwd
	}
	if deps.HomeDir == nil {
		deps.HomeDir = defaults.HomeDir
	}
	if deps.Spawn == nil {
		deps.Spawn = defaults.Spawn
	}
	if deps.IsTerminal == nil {
		deps.IsTerminal = defaults.IsTerminal
	}
	return deps
}

func isTerminalReader(reader io.Reader) bool {
	file, ok := reader.(*os.File)
	if !ok {
		return false
	}
	return isatty.IsTerminal(file.Fd()) || isatty.IsCygwinTerminal(file.Fd())
}

func scheduleLifecycle(command *cobra.Command, deps Dependencies, info BuildInfo) {
	top := command
	for top.Parent() != nil && top.Parent().Parent() != nil {
		top = top.Parent()
	}
	name := top.Name()
	if strings.HasPrefix(name, "__") || name == "uninstall" {
		return
	}
	executable, err := deps.Executable()
	if err != nil {
		return
	}
	_ = deps.Spawn(executable, "__update-worker", "--version", info.Version, "--target", info.Target)
	if name != "runx" && name != "agent" && name != "upgrade" {
		cwd, err := deps.Getwd()
		if err == nil {
			_ = deps.Spawn(executable, "__maintenance-worker", "--cwd", cwd)
		}
	}
}

func bundledSkill() (string, error) {
	value, err := embed.FS.ReadFile("skills/guiho-s-runx.SKILL.md")
	return string(value), err
}

func bundledInstruction() (string, error) {
	value, err := embed.FS.ReadFile("prompts/guiho-i-runx.md")
	return string(value), err
}

func effectiveCWD(deps Dependencies, requested string) (string, error) {
	if requested == "" {
		var err error
		requested, err = deps.Getwd()
		if err != nil {
			return "", err
		}
	}
	absolute, err := filepath.Abs(requested)
	if err != nil {
		return "", err
	}
	return absolute, nil
}

func newMaintenanceWorkerCommand(deps Dependencies) *cobra.Command {
	cwd := os.Getenv("RUNX_MAINTENANCE_CWD")
	command := &cobra.Command{Use: "__maintenance-worker", Hidden: true, Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		skill, err := bundledSkill()
		if err != nil {
			return err
		}
		home, err := deps.HomeDir()
		if err != nil {
			return err
		}
		_, err = maintenance.MaintainAgentIntegration(cwd, home, skill)
		return err
	}}
	command.Flags().StringVar(&cwd, "cwd", "", "worker cwd")
	return command
}

func newUpdateWorkerCommand(deps Dependencies, info BuildInfo) *cobra.Command {
	var version, target string
	command := &cobra.Command{Use: "__update-worker", Hidden: true, Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		_, err := update.RunUpdateWorker(update.WorkerOptions{CurrentVersion: version, BuildTarget: target, HTTPClient: deps.HTTPClient, Now: deps.Now})
		return err
	}}
	command.Flags().StringVar(&version, "version", info.Version, "worker version")
	command.Flags().StringVar(&target, "target", info.Target, "worker build target")
	return command
}

// newSelfTestCommand implements the mandatory hidden installation self-test:
// it proves the payload starts, its embedded release resources are readable,
// and its build target/version are coherent without mutating the installation.
func newSelfTestCommand(deps Dependencies, info BuildInfo) *cobra.Command {
	return &cobra.Command{Use: "__self-test", Hidden: true, Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		skill, err := bundledSkill()
		if err != nil {
			return withExitCode(5, fmt.Errorf("embedded skill unreadable: %w", err))
		}
		instruction, err := bundledInstruction()
		if err != nil {
			return withExitCode(5, fmt.Errorf("embedded instruction unreadable: %w", err))
		}
		result := struct {
			SelfTest         string `json:"selfTest"`
			Version          string `json:"version"`
			BuildTarget      string `json:"buildTarget"`
			Protocol         int    `json:"protocol"`
			SkillBytes       int    `json:"skillBytes"`
			InstructionBytes int    `json:"instructionBytes"`
		}{
			SelfTest: "ok", Version: info.Version, BuildTarget: info.Target,
			Protocol:   installstate.ProtocolVersion,
			SkillBytes: len(skill), InstructionBytes: len(instruction),
		}
		return writeJSON(command, result)
	}}
}

type ChildExitError struct{ Code int }

func (value *ChildExitError) Error() string {
	return fmt.Sprintf("configured command exited with code %d", value.Code)
}

func runtimePlatform() string { return runtime.GOOS }
