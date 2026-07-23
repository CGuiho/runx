package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/CGuiho/runx/pkg/update"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type UpdateCheckResult struct {
	CurrentVersion  string `json:"currentVersion"`
	LatestVersion   string `json:"latestVersion"`
	UpdateAvailable bool   `json:"updateAvailable"`
	URL             string `json:"url,omitempty"`
}

var upgradeCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if a new RunX version is available",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentVer := "0.8.0-dev"

		platform, err := update.ResolveUpgradePlatform(runtime.GOOS, runtime.GOARCH)
		if err != nil {
			return err
		}

		catalog, err := update.FetchReleaseCatalog("", platform, currentVer, nil)
		latestVer := currentVer
		updateAvailable := false
		url := ""

		if err == nil && catalog != nil {
			if catalog.LatestStableVersion != "" {
				latestVer = catalog.LatestStableVersion
			}
			updateAvailable = update.CompareVersions(latestVer, currentVer) > 0
			if updateAvailable {
				for _, rel := range catalog.Releases {
					if rel.Version == latestVer && rel.CompatibleAsset != nil {
						url = rel.CompatibleAsset.URL
						break
					}
				}
			}
		}

		res := UpdateCheckResult{
			CurrentVersion:  currentVer,
			LatestVersion:   latestVer,
			UpdateAvailable: updateAvailable,
			URL:             url,
		}

		if jsonFmt || viper.GetBool("json") {
			data, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		if updateAvailable {
			fmt.Printf("⚠ New version available: v%s\n  Run 'runx upgrade' to update.\n", latestVer)
		} else {
			fmt.Printf("RunX is currently at version v%s (up to date).\n", currentVer)
		}

		return nil
	},
}

func init() {
	upgradeCmd.AddCommand(upgradeCheckCmd)
}
