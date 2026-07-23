package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/CGuiho/runx/pkg/update"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pageFlag        int
	perPageFlag     int
	preReleasesFlag bool
	upgradeArch     string
	upgradeVariant  string
)

var upgradeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List RunX releases newest first.",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentVer := "0.8.0-dev"

		goarch := runtime.GOARCH
		if upgradeArch != "" {
			goarch = upgradeArch
		}

		platform, err := update.ResolveUpgradePlatform(runtime.GOOS, goarch)
		if err != nil {
			return err
		}

		if upgradeVariant != "" {
			platform.Variant = upgradeVariant
		}

		catalog, err := update.FetchReleaseCatalog("", platform, currentVer, nil)
		if err != nil {
			return fmt.Errorf("failed to list releases: %w", err)
		}

		if pageFlag > 0 || perPageFlag > 0 {
			p := pageFlag
			if p <= 0 {
				p = 1
			}
			pp := perPageFlag
			if pp <= 0 {
				pp = 20
			}
			start := (p - 1) * pp
			if start < len(catalog.Releases) {
				end := start + pp
				if end > len(catalog.Releases) {
					end = len(catalog.Releases)
				}
				catalog.Releases = catalog.Releases[start:end]
			} else {
				catalog.Releases = []update.ReleaseCatalogEntry{}
			}
		}

		if jsonFmt || viper.GetBool("json") {
			data, err := json.MarshalIndent(catalog, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Println("AVAILABLE RUNX VERSIONS")
		fmt.Println()
		fmt.Printf("%-15s  %-10s  %-12s  %-7s  %-7s  %-5s\n", "VERSION", "CHANNEL", "PUBLISHED", "CURRENT", "LATEST", "ASSET")
		for _, rel := range catalog.Releases {
			pub := "-"
			if rel.PublishedAt != nil && len(*rel.PublishedAt) >= 10 {
				pub = (*rel.PublishedAt)[:10]
			}
			curr := ""
			if rel.Current {
				curr = "yes"
			}
			latest := ""
			if rel.LatestStable {
				latest = "yes"
			}
			asset := "no"
			if rel.CompatibleAsset != nil {
				asset = "yes"
			}

			fmt.Printf("%-15s  %-10s  %-12s  %-7s  %-7s  %-5s\n", rel.Version, rel.Channel, pub, curr, latest, asset)
		}

		return nil
	},
}

func init() {
	upgradeListCmd.Flags().IntVar(&pageFlag, "page", 0, "Select result page.")
	upgradeListCmd.Flags().IntVar(&perPageFlag, "per-page", 0, "Select page size.")
	upgradeListCmd.Flags().BoolVar(&preReleasesFlag, "pre-releases", false, "Accepted explicitly; prereleases are always included.")
	upgradeListCmd.Flags().StringVar(&upgradeArch, "arch", "", "Select target architecture.")
	upgradeListCmd.Flags().StringVar(&upgradeVariant, "variant", "", "Select x64 variant.")

	upgradeCmd.AddCommand(upgradeListCmd)
}
