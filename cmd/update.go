/*
Copyright © 2022 Nottimisreal
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/theckman/yacspin"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates All Packages Or A List Of Packages",
	Long:  `Finds the latest supported version of the packages provided`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		command := exec.Command("ferment", "list")
		var out bytes.Buffer
		command.Stdout = &out
		command.Run()
		pkgs := out.String()
		pkgArr := strings.Split(pkgs, "\n")
		//remove first element
		pkgArr = pkgArr[1:]
		return pkgArr, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		type summary struct {
			pkg        string
			updated    bool
			oldVersion string
			newVersion string
		}
		ferment, _ := os.Executable()
		if len(args) == 0 {
			command := exec.Command(ferment, "list")
			var out bytes.Buffer
			command.Stdout = &out
			command.Run()
			pkgs := out.String()
			args = strings.Split(pkgs, "\n")
			//remove first element
			args = args[1:]
			//remove any args that only contain spaces or \n
			for i, arg := range args {
				if arg == "" || arg == "\n" {
					args = append(args[:i], args[i+1:]...)
					continue
				}
				arg = strings.Replace(arg, "\n", "", -1)
				args[i] = arg

			}
		}
		spinner, err := yacspin.New(yacspin.Config{
			Frequency:         100 * time.Millisecond,
			CharSet:           yacspin.CharSets[57],
			Suffix:            color.GreenString(" Update"),
			SuffixAutoColon:   true,
			StopCharacter:     "✓",
			StopColors:        []string{"fgGreen"},
			StopFailCharacter: "✗",
			StopFailColors:    []string{"fgRed"},
		})
		if err != nil {
			panic(err)
		}
		var summaryArr []summary
		spinner.Start()
		for _, pkg := range args {

			version, err := getVersion(pkg)
			if err != nil {
				spinner.StopFailMessage(fmt.Sprintf("%s: %s", pkg, err))
				spinner.StopFail()
				os.Exit(1)
			}
			spinner.Message(fmt.Sprintf("Latest Version Of %s: %s", pkg, version))
			spinner.Message(fmt.Sprintf("Retrieveing Current Version Of %s", pkg))
			currentVersion, err := getCurrentVersion(pkg)
			if err != nil {
				spinner.Message("Error Retrieving Current Version")
				summaryArr = append(summaryArr, summary{pkg: pkg, updated: false, oldVersion: version, newVersion: version})
				continue
			}
			spinner.Message(fmt.Sprintf("Current Version Of %s: %s", pkg, currentVersion))
			if version == currentVersion {
				spinner.Message(fmt.Sprintf("%s is up to date", pkg))
				summaryArr = append(summaryArr, summary{pkg: pkg, updated: false, oldVersion: currentVersion, newVersion: version})
				continue
			}
			spinner.Message(fmt.Sprintf("%s is not up to date", pkg))
			spinner.Message(fmt.Sprintf("Updating %s", pkg))
			cmd := exec.Command(ferment, "uninstall", pkg)
			err = cmd.Run()
			if err != nil {
				spinner.StopFailMessage("Failed Uninstalling " + pkg)
				spinner.StopFail()
				os.Exit(1)
			}
			cmd = exec.Command(ferment, "install", pkg)
			err = cmd.Run()
			if err != nil {
				spinner.StopFailMessage("Failed Reinstalling " + pkg)
				spinner.StopFail()
				os.Exit(1)
			}
			summaryArr = append(summaryArr, summary{pkg: pkg, updated: true, oldVersion: currentVersion, newVersion: version})
		}
		spinner.Stop()
		color.HiBlue("Summary:")
		for _, summary := range summaryArr {
			if summary.updated {
				color.Green("%s: %s -> %s", summary.pkg, summary.oldVersion, summary.newVersion)
			} else {
				color.Yellow("%s: %s (Unchanged)", summary.pkg, summary.oldVersion)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func getVersion(pkg string) (string, error) {
	version, err := executeQuickPython(fmt.Sprintf("from %s import %s;pkg=%s();print(pkg.version)", pkg, pkg, pkg))
	if err != nil {
		return "", err
	}
	version = strings.Replace(version, "\n", "", -1)
	return version, nil
}
func getCurrentVersion(pkg string) (string, error) {
	//go into the VERSION.meta file of the package and get the version
	content, err := os.ReadFile(fmt.Sprintf("%s/Installed/%s/VERSION.meta", location, pkg))
	if err != nil {
		return "", err
	}
	return string(content), nil

}
