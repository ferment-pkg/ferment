/*
Copyright © 2022 NotTimIsReal

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dchest/validator"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Long:  `Shows All Installed Packages`,
	Short: "Shows All Installed Packages Excluding Git Packages",
	Run: func(cmd *cobra.Command, args []string) {
		ListAllNormalPackages()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func ListAllNormalPackages() {
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	pkgs, err := os.ReadDir(fmt.Sprintf("%s/Installed", location))
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println(color.RedString("===> ") + color.RedString("No Packages Installed"))
			os.Exit(1)
		}
		panic(err)
	}
	//convert pkgs to array
	var pkgsArray []string
	for _, pkg := range pkgs {
		name := pkg.Name()
		//check if name is a domain name
		domain := validator.IsValidDomain(name)
		if domain {
			pkgsArray = append(pkgsArray, fmt.Sprintf("%s (Is Repo Not Indexed)", name))
			continue
		}
		pkgsArray = append(pkgsArray, name)
	}
	fmt.Println(color.BlueString("===> ") + color.GreenString("Listing All Installed Packages"))
	for _, pkg := range pkgsArray {
		fmt.Println(pkg)
	}

}
