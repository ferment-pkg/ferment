/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// reinstallCmd represents the reinstall command
var reinstallCmd = &cobra.Command{
	Use:   "reinstall",
	Short: "Reinstall A Package",
	Long:  `Execute the uninstall and install commands for a package.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a package to reinstall")
			os.Exit(1)
		}
		for _, pkg := range args {
			if !IsUrl(pkg) {
				checkIfPackageExists(pkg)
			} else {
				if strings.Contains(pkg, "http://") {
					fmt.Println("Ferment Does Not Support http packages")
					continue
				}
				pkg = strings.Split(pkg, "https://")[1]
				if !checkIfPackageExists(strings.ToLower(pkg)) {
					fmt.Println("Package Does Not Exist")
					continue
				}
			}
			pkg = convertToReadableString(strings.ToLower(pkg))
			ferment, _ := os.Executable()
			cmd := exec.Command(ferment, "uninstall", pkg)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			cmd = exec.Command(ferment, "install", pkg)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}
	},
}

func init() {
	rootCmd.AddCommand(reinstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reinstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reinstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
