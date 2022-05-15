/*
Copyright © 2022 NotTimIsReal

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

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
		content, err := os.ReadFile("VERSION.meta")
		if err != nil {
			panic(err)
		}
		fmt.Println("Local Version: " + string(content))
		fmt.Println("Getting Latest Version On GitHub...")
		repo, err := git.PlainOpen(".")
		if err != nil {
			panic(err)
		}
		w, err := repo.Worktree()
		if err != nil {
			panic(err)
		}
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			if strings.Contains(err.Error(), "already up-to-date") {
				fmt.Println("Already Up To Date")
				os.Exit(0)
			}
			panic(err)
		}
		fmt.Println("Version Updated To Latest")

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
