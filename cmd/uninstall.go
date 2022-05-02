/*
Copyright © 2022 NotTimIsReal

*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall <package>",
	Short: "Uninstall A Package",
	Long:  `Uninstall A Package That Has Been Installed By Ferment`,
	Run: func(cmd *cobra.Command, args []string) {
		location, err := os.Executable()
		if err != nil {
			panic(err)
		}
		location = location[:len(location)-len("/ferment")]
		if len(args) == 0 {
			fmt.Println("Please specify a package to uninstall")
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
				os.RemoveAll(fmt.Sprintf("%s/Installed/%s", location, strings.ToLower(pkg)))
				fmt.Println(color.GreenString("Package Uninstalled Successfully"))
				continue
			}
			GetUninstallInstructions(pkg)
		}

	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func GetUninstallInstructions(pkg string) {
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	content, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, strings.ToLower(pkg)))
	if err != nil {
		fmt.Println("Package Does Not Exist")
		os.Exit(1)
	}
	cmd := exec.Command("python3")
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	cmd.Stdout = w
	cmd.Stderr = w
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", pkg))
	io.WriteString(closer, fmt.Sprintf("pkg.cwd=%s\n", fmt.Sprintf(`"%s/Installed/%s"`, location, pkg)))
	io.WriteString(closer, "pkg.uninstall()\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if strings.Contains(buf.String(), "True") {
		os.RemoveAll(fmt.Sprintf("%s/Installed/%s", location, pkg))
		fmt.Println(color.GreenString("Package Uninstalled Successfully"))
		os.Exit(0)
	} else {
		fmt.Println(color.RedString("Package Uninstall Failed"))
		os.Exit(1)
	}

}
func checkIfPackageExists(pkg string) bool {
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	_, err = os.ReadDir(fmt.Sprintf("%s/Installed/%s", location, pkg))
	if err != nil {
		return true
	}
	return false
}
