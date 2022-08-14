/*
Copyright © 2022 NotTimIsReal
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var c = make(chan Dep)
var uninstallCmd = &cobra.Command{
	Use:   "uninstall <package>",
	Short: "Uninstall A Package",
	Long:  `Uninstall A Package That Has Been Installed By Ferment`,
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
			if e, _ := exists(fmt.Sprintf("%s/Installed/%s", location, pkg)); !e {
				color.Red("Package Not %s installed\n", pkg)
				continue
			}
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
	//uninstallCmd.Flags().Bool("y", false, "Help message for toggle")
}
func GetUninstallInstructions(pkg string) {
	pkg = convertToReadableString(strings.ToLower(pkg))
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
	go getDepInfo(pkg)
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
		fmt.Println(color.GreenString("Package %s Uninstalled Successfully", pkg))
		removePkg(pkg)
		color.Yellow("Removing Dependencies...")
		dep := <-c
		for _, d := range dep.Deps {
			removeDep(d.Name)
			GetUninstallInstructions(d.Name)

		}

	} else {
		fmt.Println(color.RedString("Package Uninstall Failed"))
		os.Exit(1)
	}

}
func checkIfPackageExists(pkg string) bool {
	pkg = convertToReadableString(strings.ToLower(pkg))
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	_, err = os.ReadDir(fmt.Sprintf("%s/Installed/%s", location, pkg))
	return err == nil
}
func getDepInfo(pkg string) {
	pkg = convertToReadableString(strings.ToLower(pkg))
	location, err := os.Executable()
	if err != nil {
		fmt.Printf("%s: %s\n", color.RedString("ERROR"), err.Error())
		os.Exit(1)
	}
	location = location[:len(location)-len("/ferment")]
	os.Chdir(location)
	content, err := os.ReadFile("dependencies.json")
	if err != nil {
		fmt.Printf("%s: %s\n", color.RedString("ERROR"), err.Error())
		os.Exit(1)
	}
	var dependencies Dep
	json.Unmarshal(content, &dependencies)
	var realDeps Dep
	for _, dep := range dependencies.Deps {
		if strings.Contains(dep.ReliedBy, pkg) && !dep.InstalledByUser {
			s := strings.Split(dep.ReliedBy, " ")
			var requiredByDiff bool = false
			for _, d := range s {
				if d != "" && d != pkg {
					requiredByDiff = true
				}
			}
			if requiredByDiff {
				continue
			}
			realDeps.Deps = append(realDeps.Deps, dep)

		}
	}
	c <- realDeps

}
func removeDep(pkg string) {
	pkg = convertToReadableString(strings.ToLower(pkg))
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	os.Chdir(location)
	content, err := os.ReadFile("dependencies.json")
	if err != nil {
		fmt.Printf("%s: %s\n", color.RedString("ERROR"), err.Error())
		os.Exit(1)
	}
	var dependencies Dep
	json.Unmarshal(content, &dependencies)
	for i, dep := range dependencies.Deps {
		if strings.Contains(dep.ReliedBy, pkg) || dep.Name == pkg {
			dependencies.Deps = append(dependencies.Deps[:i], dependencies.Deps[i+1:]...)
		}
	}
	dependencies.LastUpdated = time.Now().UTC().Unix()
	content, err = json.Marshal(dependencies)
	if err != nil {
		fmt.Printf("%s: %s\n", color.RedString("ERROR"), err.Error())
		os.Exit(1)
	}
	err = os.WriteFile("dependencies.json", content, 0644)
	if err != nil {
		fmt.Printf("%s: %s\n", color.RedString("ERROR"), err.Error())
		os.Exit(1)
	}
}
func removePkg(pkg string) {
	pkg = convertToReadableString(strings.ToLower(pkg))
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	os.Chdir(location)
	content, err := os.ReadFile("dependencies.json")
	if err != nil {
		fmt.Printf("%s: %s\n", color.RedString("ERROR"), err.Error())
		os.Exit(1)
	}
	var dependencies Dep
	json.Unmarshal(content, &dependencies)
	dependencies.LastUpdated = time.Now().UTC().Unix()
	for i, dep := range dependencies.Deps {
		if dep.Name == pkg {
			dependencies.Deps = append(dependencies.Deps[:i], dependencies.Deps[i+1:]...)
		}
	}
	content, err = json.Marshal(dependencies)
	if err != nil {
		fmt.Printf("%s: %s\n", color.RedString("ERROR"), err.Error())
		os.Exit(1)
	}
	err = os.WriteFile("dependencies.json", content, 0644)
	if err != nil {
		fmt.Printf("%s: %s\n", color.RedString("ERROR"), err.Error())
		os.Exit(1)
	}
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
