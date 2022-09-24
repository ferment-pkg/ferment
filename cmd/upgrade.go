/*
Copyright © 2022 NotTimIsReal
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Updates The Package Manager",
	Long:  `Updates ferment to the latest version on github.`,
	Run: func(cmd *cobra.Command, args []string) {
		location, _ := os.Executable()
		location = location[:len(location)-len("ferment")]
		os.Chdir(location)
		fmt.Println("Getting Local Version On System...")
		config := getConfig(location)
		fmt.Println("Local Version: " + string(config.Version))
		fmt.Println("Getting Latest Version On GitHub...")
		repo, err := git.PlainOpen(".")
		if err != nil {
			panic(err)
		}
		w, err := repo.Worktree()
		if err != nil {
			panic(err)
		}
		w.Pull(&git.PullOptions{RemoteName: "origin"})
		fmt.Println("Resetting Repository...")
		w.Reset(&git.ResetOptions{Mode: git.HardReset})
		fmt.Println("Version Updated To Latest")
		fmt.Println("Updating Packages...")
		err = os.Chdir("Barrells")
		if err != nil {
			color.RedString("ERRORS: %s", err)
		}
		repo, err = git.PlainOpen(".")
		if err != nil {
			panic(err)
		}
		w, err = repo.Worktree()
		if err != nil {
			panic(err)
		}
		w.Pull(&git.PullOptions{RemoteName: "origin"})
		fmt.Println("Successfully Updated All Packages!")
		fmt.Println("Downloading Binary For Ferment...")
		resp, _ := http.Get(fmt.Sprintf("https://github.com/ferment-pkg/ferment/releases/download/v%s/ferment", config.Version))
		binary, _ := io.ReadAll(resp.Body)
		os.WriteFile(location+"ferment", binary, 0777)
		fmt.Println("Successfully Updated Ferment!")

	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upgradeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upgradeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Config struct {
	Version string `json:"version"`
}

func getConfig(location string) Config {
	config, err := os.ReadFile(location + "/ferment.config.json")
	if err != nil {
		resp, _ := http.Get("https://raw.githubusercontent.com/NotTimIsReal/ferment/main/ferment.config.json")
		config, _ = io.ReadAll(resp.Body)
		os.WriteFile(location+"ferment.config.json", config, 0777)
		configReturnValue := Config{}
		json.Unmarshal(config, &configReturnValue)
		return configReturnValue
	}
	configReturnValue := Config{}
	err = json.Unmarshal(config, &configReturnValue)
	if err != nil {
		panic(err)
	}
	return configReturnValue
}
