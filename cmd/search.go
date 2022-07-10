/*
Copyright © 2022 NotTimIsReal

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [package]",
	Short: "Search For Default Packages In Ferment",
	Long:  `Search For Default Packages In Ferment OR List All Packages Installable In Ferment`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(color.BlueString("===> ") + color.GreenString("Looking For Package In Barrells"))
		if len(args) == 0 {
			args = append(args, "")
		}
		SearchForPackages(args[0])
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func SearchForPackages(pkg string) {
	pkg = convertToReadableString(strings.ToLower(pkg))
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	pkgs, err := os.ReadDir(fmt.Sprintf("%s/Barrells", location))
	if err != nil {
		panic(err)
	}
	//convert pkgs to array
	var pkgsArray []string
	for _, pkg := range pkgs {
		name := pkg.Name()
		name = strings.ReplaceAll(name, ".py", "")
		if name == "__init__" {
			continue
		}
		if name == "__pycache__" {
			continue
		}
		if name == "index" {
			continue
		}
		if name == ".git" || name == ".gitignore" || name == ".github" || name == "helpers" {
			continue
		}
		if name == ".DS_STORE" {
			continue
		}
		pkgsArray = append(pkgsArray, name)
	}
	//filter pkgsArray to only pkgs that contain pkg
	var filteredPkgsArray []string
	for _, pkgStr := range pkgsArray {
		if strings.Contains(pkgStr, pkg) {
			filteredPkgsArray = append(filteredPkgsArray, pkgStr)
		}
	}
	for _, pkgStr := range filteredPkgsArray {
		if pkgStr == pkg {
			fmt.Println(color.BlueString("===> ") + color.GreenString("Exact Match: ") + color.GreenString(pkgStr))
			os.Exit(0)
		} else {
			fmt.Println(color.YellowString(pkgStr))
		}

	}

}
