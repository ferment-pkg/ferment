/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ownCmd represents the own command
var ownCmd = &cobra.Command{
	Use:   "own",
	Short: "Take ownership of a installed package",
	Long:  `Changes the value of InstalledByUser to true for a dependency`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please specify a package")
			os.Exit(1)
		}
		location, _ := os.Executable()
		location = location[:len(location)-len("/ferment")]
		os.Chdir(location)
		c, err := os.ReadFile("dependencies.json")
		if err != nil {
			fmt.Println("Error reading dependencies.json")
			os.Exit(1)
		}
		var dependencies Dep
		json.Unmarshal(c, &dependencies)
		for _, pkg := range args {
			for i, dep := range dependencies.Deps {
				if dep.Name == pkg {
					fmt.Println(dep.InstalledByUser)
					dependencies.Deps[i].InstalledByUser = true
					break
				}
			}

		}
		c, err = json.Marshal(dependencies)
		if err != nil {
			fmt.Println("Error marshalling dependencies.json")
			os.Exit(1)
		}
		os.WriteFile("dependencies.json", c, 0644)
	},
}

func init() {
	rootCmd.AddCommand(ownCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ownCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ownCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
