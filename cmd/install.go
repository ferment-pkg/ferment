/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"net/url"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <package>",
	Short: "Install Packages",
	Long:  `Install Official Packages or Custom Packages From Git Repositories From GitLab Or Github`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a package to install, it can either be a custom package from github, gitlab, etc or a official package")
			os.Exit(1)
		}
		var foundPkg bool = false
		verbose, err := cmd.Flags().GetString("verbose")
		pyexec, err := cmd.Flags().GetString("python-exec")

		if err != nil {

			panic(err)
		}
		location, err := os.Executable()
		//redefine location so that it is the directory of the executable
		location = location[:len(location)-len("ferment")]
		if err != nil {

			panic(err)
		}
		if verbose == "" {
			verbose = "false"
		}
		if !IsUrl(args[0]) {
			//search for package in default list
			if verbose == "true" {
				fmt.Println("Searching for package in default list")
			}
			files, err := os.ReadDir(fmt.Sprintf("%s/Barrells", location))
			if err != nil {
				panic(err)
			}
			for _, v := range files {
				if strings.Split(v.Name(), ".")[0] == args[0] {
					if verbose == "true" {
						fmt.Println("Found package in default list")
					}
					foundPkg = true
					break
				}
			}
			if !foundPkg {
				fmt.Println("Package not found in default list or https URL invalid")
				os.Exit(1)
			}

		}
		s := spinner.New(spinner.CharSets[2], 100*time.Millisecond) // Build our new spinner
		s.Suffix = color.GreenString(" Downloading Source...")
		s.Start()
		time.Sleep(4 * time.Second) // Run for some time to simulate work
		s.Stop()

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringP("verbose", "v", "", "Log All Output")
	installCmd.PersistentFlags().String("python-exec", "/usr/bin/python3", "Python Executable Location")
	installCmd.Flag("verbose").NoOptDefVal = "true"

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return (err == nil && u.Scheme != "" && u.Host != "")
}
