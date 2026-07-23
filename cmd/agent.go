package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage CLI-owned agent skills, instructions, and prompts",
	Long:  `Manage RunX agent integration, skills, instructions, and prompts.`,
}

func init() {
	RootCmd.AddCommand(agentCmd)
}

func writeFormatted(cmd *cobra.Command, formatFlag string, value interface{}) error {
	var out io.Writer = os.Stdout
	if cmd != nil {
		out = cmd.OutOrStdout()
	}
	isJSON := jsonFmt || formatFlag == "json"

	switch v := value.(type) {
	case string:
		if isJSON {
			bytes, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(out, string(bytes))
		} else {
			fmt.Fprint(out, v)
		}
	default:
		bytes, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(bytes))
	}
	return nil
}
