/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"net/url"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <package>",
	Short: "Install Packages",
	Long:  `Install Official Packages or Custom Packages From Git Repositories From GitLab Or Github`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args[0]) == 0 {
			fmt.Println("Please provide a package to install, it can either be a custom package from github, gitlab, etc or a official package")
			os.Exit(1)
		}
		verbose, err := cmd.Flags().GetString("verbose")
		if err != nil {

			panic(err)
		}
		location, err := os.Executable()
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

		} else {
			args[0] = strings.Split(args[0], "https://")[1]
		}
		location = strings.Split(location, "/ferment")[0]
		location = fmt.Sprintf("%s/PKG/%s/", location, args[0])
		_, err = git.PlainClone(location, false, &git.CloneOptions{
			URL: "https://" + args[0],
		})
		if verbose == "true" {
			fmt.Println("Cloned Repository")
		}
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				fmt.Println("Package already exists")
				os.Exit(1)
			}
			fmt.Println(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringP("verbose", "v", "", "Log All Output")
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
