/*
Copyright © 2022 NotTimIsReal

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"archive/tar"
	"compress/gzip"
	"net/http"
	"net/url"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var l, _ = os.Executable()
var location = l[:len(l)-len("/ferment")]

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
		for _, pkg := range args {
			var foundPkg bool = false
			verbose, err := cmd.Flags().GetString("verbose")

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
			pkg = convertToReadableString(strings.ToLower(pkg))
			if !IsUrl(pkg) {
				//search for package in default list
				if verbose == "true" {
					fmt.Println("Searching for package in default list")
				}
				files, err := os.ReadDir(fmt.Sprintf("%s/Barrells", location))
				if err != nil {
					panic(err)
				}

				for _, v := range files {
					name := strings.ToLower(v.Name())
					if strings.Split(name, ".")[0] == strings.ToLower(pkg) {
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

			if !foundPkg {
				color.RedString("ERROR: %s", errors.New("downloading from url is not supported anymore"))
				os.Exit(1)
			}
			if checkifPrebuildSuitable(pkg) {
				s := spinner.New(spinner.CharSets[2], 100*time.Millisecond) // Build our new spinner
				s.Suffix = color.GreenString(" Downloading Prebuild...")
				s.Start()
				DownloadFromTar(pkg, *getPrebuildURL(pkg), verbose)
				s.Stop()
				installPackages(pkg, verbose, false, "")
				os.Exit(0)
			}
			if UsingGit(pkg, verbose) {
				s := spinner.New(spinner.CharSets[2], 100*time.Millisecond) // Build our new spinner
				s.Suffix = color.GreenString(" Downloading Source...")
				s.Start()
				url := GetGitURL(pkg, verbose)
				err := DownloadFromGithub(url, fmt.Sprintf("%s/Installed/%s", location, pkg), verbose)
				if err != nil {
					s.Stop()
					fmt.Println(color.RedString(err.Error()))
					fmt.Println(color.RedString("Aborting"))
					os.Exit(1)
				}
				s.Stop()
				installPackages(pkg, verbose, false, "")

			} else {
				_, err := os.ReadDir(fmt.Sprintf("%s/Installed/%s", location, pkg))
				if err == nil {
					fmt.Println(color.GreenString("Package Already Installed"))
					continue
				}
				s := spinner.New(spinner.CharSets[2], 100*time.Millisecond) // Build our new spinner
				s.Suffix = color.GreenString(" Downloading Source...")
				s.Start()
				GetDownloadUrl(pkg, verbose)
				DownloadInstructions(pkg)
				s.Stop()
				s = spinner.New(spinner.CharSets[36], 100*time.Millisecond) // Build our new spinner
				s.Suffix = color.GreenString(" Installing Package...")
				installPackages(pkg, verbose, false, "")
			}
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
func DownloadFromGithub(url string, path string, verbose string) error {
	if verbose == "true" {
		fmt.Println("Downloading from Github")
	}
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			return fmt.Errorf("package already exists")
		}
		panic(err)
	}
	if verbose == "true" {
		fmt.Println("Downloaded from Github")
	}
	return nil
}
func UsingGit(pkg string, verbose string) bool {
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("ferment")]
	content, err := os.ReadFile(fmt.Sprintf("%sBarrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println(color.RedString("Reinstall ferment, Barrells is missing"))
			os.Exit(1)
		}
		panic(err)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	if verbose == "true" {
		fmt.Println("Starting STDIN pipe")
	}
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "print(pkg.git)\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String() == "True\n"
	//fmt.Println(out)

}
func installPackages(pkg string, verbose string, isDep bool, installedBy string) {
	if verbose == "true" {
		fmt.Println("Looking For Dependencies")
	}
	//Variables
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	content, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println(color.RedString("Reinstall ferment, Barrells is missing"))
			os.Exit(1)
		}
		panic(err)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	if verbose == "true" {
		fmt.Println("Starting STDIN pipe")
	}
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "print(pkg.dependencies)\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	dependencies := strings.Split(buf.String(), "\n")[0]
	dependencies = strings.Replace(dependencies, "[", "", 1)
	dependencies = strings.Replace(dependencies, "]", "", 1)
	dependencies = strings.Replace(dependencies, " ", "", -1)
	dependencies = strings.Replace(dependencies, "'", "", -1)
	dependenciesArr := strings.Split(dependencies, ",")
	//run a function for each dependecy in dependenciesArr
	if !strings.Contains(dependencies, "Traceback(mostrecentcalllast)") {
		for _, dep := range dependenciesArr {
			fmt.Printf(color.YellowString("Package %s depends on %s\n"), pkg, dep)
			cmd := exec.Command("which", strings.ReplaceAll(dep, "'", ""))
			r, w, err := os.Pipe()
			if err != nil {
				panic(err)
			}
			cmd.Stdout = w
			cmd.Start()
			w.Close()
			cmd.Wait()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			if buf.String() != "" {
				fmt.Printf(color.YellowString("%s is already installed\n"), dep)
				EditDepsTracker(dep, pkg)
				fmt.Println(color.YellowString("Skipping"))
				continue
			}
			if IsLib(dep) && checkIfPackageExists(dep) {
				fmt.Printf(color.YellowString("%s is a library and is already installed\n"), dep)
				EditDepsTracker(dep, pkg)
				fmt.Println(color.YellowString("Skipping"))
				continue
			}
			_, err = os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
			if err != nil {
				fmt.Println("Not Downloadable By Ferment, Skipping")
				continue
			}
			fmt.Printf(color.YellowString("Now Downloading %s\n"), dep)
			if UsingGit(dep, verbose) {
				url := GetGitURL(dep, verbose)
				err := DownloadFromGithub(url, fmt.Sprintf("%s/Installed/%s", location, dep), verbose)
				if err != nil {
					panic(err)
				}
				installPackages(dep, verbose, true, pkg)
				//TestInstallationScript(dep, verbose)
			} else {
				GetDownloadUrl(dep, verbose)
				installPackages(dep, verbose, true, pkg)
				// TestInstallationScript(dep, verbose)
			}

		}
	}
	if e, _ := checkIfDepExists(pkg); e {
		EditDepsTracker(pkg, installedBy)
	} else {
		DepTrackerAdd(pkg, isDep, installedBy)
	}
	fmt.Printf("Installing %s", pkg)
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Suffix = " Installing " + pkg
	s.FinalMSG = color.GreenString("Installed " + pkg + "\n")
	s.Start()
	if checkifPrebuildSuitable(pkg) {
		InstallPrebuilds(pkg)
	} else {
		RunInstallationScript(convertToReadableString(strings.ToLower(pkg)), verbose, convertToReadableString(strings.ToLower(pkg)))
	}
	s.Stop()
	s = spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Suffix = " Installing Binaries For " + pkg
	s.FinalMSG = color.GreenString("Installed Binary For " + pkg + "\n")
	s.Start()
	msg := InstallBinary(pkg, verbose)
	if msg == "No Binary" {
		s.FinalMSG = color.GreenString("No Binaries Needed To Be Installed")

	}
	s.Stop()

	s = spinner.New(spinner.CharSets[36], 100*time.Millisecond)
	s.Suffix = " Testing " + pkg
	s.Start()
	result := TestInstallationScript(pkg, verbose)
	if result {
		s.FinalMSG = color.GreenString("Successfully Tested Installed " + pkg + "\n")
	} else {
		s.FinalMSG = color.RedString("Installed " + pkg + " Failed Test\n")
	}
	s.Stop()

}
func GetGitURL(pkg string, verbose string) string {
	if verbose == "true" {
		fmt.Println("Looking For GitURl")
	}
	//Variables
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	content, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println(color.RedString("Reinstall ferment, Barrells is missing"))
			os.Exit(1)
		}
		panic(err)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	if verbose == "true" {
		fmt.Println("Starting STDIN pipe")
	}
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "print(pkg.url)\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()

}
func DownloadFromTar(pkg string, url string, verbose string) string {
	var isGZ bool
	if strings.Contains(url, ".gz") {
		isGZ = true
	} else {
		isGZ = false
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(color.RedString("Unable to download %s", pkg))
		panic(err)
	}
	location, err := os.Executable()
	location = location[:len(location)-len("/ferment")]
	if err != nil {
		fmt.Println(color.RedString("Unable to download %s", pkg))
		panic(err)
	}
	defer resp.Body.Close()
	if verbose == "true" {
		fmt.Printf("Downloading Tar From %s\n", url)
	}
	if verbose == "true" {
		fmt.Println("Extracting Tar")
	}
	path, err := Untar(fmt.Sprintf("%s/Installed/", location), resp.Body, convertToReadableString(strings.ToLower(pkg)), isGZ)
	if err != nil {
		fmt.Println(color.RedString("Unable to extract %s", pkg))
		panic(err)
	}
	return path
}
func GetDownloadUrl(pkg string, verbose string) string {
	if verbose == "true" {
		fmt.Println("Looking For GitURl")
	}
	//Variables
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	content, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println(color.RedString("Reinstall ferment, Barrells is missing"))
			os.Exit(1)
		}
		panic(err)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	if verbose == "true" {
		fmt.Println("Starting STDIN pipe")
	}
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "print(pkg.url)\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	path := DownloadFromTar(convertToReadableString(strings.ToLower(pkg)), strings.Replace(buf.String(), "\n", "", -1), verbose)
	return path
}
func Untar(dst string, r io.Reader, pkg string, isGz bool) (string, error) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return "", nil

		// return any other error
		case err != nil:
			return "", err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		header.Name = fmt.Sprintf("%s/%s", pkg, strings.Join(strings.Split(header.Name, "/")[1:], "/"))
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return "", err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return "", err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
func RunInstallationScript(pkg string, verbose string, cwd string) {
	if verbose == "true" {
		fmt.Println("Running Installation Script")
	}
	location, err := os.Executable()
	location = location[:len(location)-len("/ferment")]
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	if verbose == "true" {
		fmt.Println("Starting STDIN pipe")
	}
	r, w, _ := os.Pipe()
	defer w.Close()
	cmd.Stdout = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))

	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, fmt.Sprintf(`pkg.cwd="%s/Installed/%s"`+"\n", location, cwd))
	io.WriteString(closer, "pkg.install()\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	//print nothing
	buf.Reset()
	fmt.Println(buf.String())

}
func TestInstallationScript(pkg string, verbose string) bool {
	if verbose == "true" {
		fmt.Println("Running Installation Script")
	}
	location, err := os.Executable()
	location = location[:len(location)-len("/ferment")]
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	if verbose == "true" {
		fmt.Println("Starting STDIN pipe")
	}
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "pkg.test()\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if strings.Contains(buf.String(), "no attribute") {
		return true
	}
	return strings.Contains(buf.String(), "True")

}
func InstallBinary(pkg string, verbose string) string {
	if verbose == "true" {
		fmt.Println("Installing Binary")
	}
	location, err := os.Executable()
	location = location[:len(location)-len("/ferment")]
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	if verbose == "true" {
		fmt.Println("Starting STDIN pipe")
	}
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "print(pkg.binary)\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if strings.Contains(buf.String(), "no attribute") {
		return "No Binary"
	} else {
		binary := strings.Replace(buf.String(), "'", "", -1)
		binary = strings.Replace(buf.String(), `"`, "", -1)
		binary = strings.Replace(binary, "\n", "", -1)
		err := os.Symlink(fmt.Sprintf("%s/Installed/%s/%s", location, pkg, binary), fmt.Sprintf("/usr/local/bin/%s", binary))
		if err != nil {
			fmt.Println(err)
		}
		return "Binary Installed"
	}
}
func DownloadInstructions(pkg string) {
	location, err := os.Executable()
	location = location[:len(location)-len("/ferment")]
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	_, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, fmt.Sprintf("pkg.cwd=%s\n", fmt.Sprintf(`"%s/Installed/%s"`, location, convertToReadableString(strings.ToLower(pkg)))))
	io.WriteString(closer, "pkg.download()\n")
	closer.Close()
	w.Close()
	cmd.Wait()

}
func IsLib(pkg string) bool {
	location, _ := os.Executable()
	location = location[:len(location)-len("/ferment")]
	content, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(content)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "print(pkg.lib)\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if strings.Contains(buf.String(), "no attribute") {
		return false
	} else {
		return strings.Contains(buf.String(), "True")
	}

}
func convertToReadableString(pkg string) string {
	pkg = strings.Replace(pkg, "-", "", -1)
	pkg = strings.Replace(pkg, "_", "", -1)
	pkg = strings.Replace(pkg, ".", "", -1)
	pkg = strings.Replace(pkg, " ", "", -1)
	return pkg
}

type deps struct {
	Name            string
	ReliedBy        string
	InstalledByUser bool
}
type Dep struct {
	LastUpdated int64
	Deps        []deps
}

func DepTrackerAdd(pkg string, isDep bool, installedBy string) {
	location, _ := os.Executable()
	location = location[:len(location)-len("/ferment")]
	os.Chdir(location)
	c, err := os.ReadFile("dependencies.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			os.WriteFile("dependencies.json", []byte("{}"), 0777)
			c, err = os.ReadFile("dependencies.json")
			if err != nil {
				panic(err)
			}
		}
	}
	var dependencies Dep
	err = json.Unmarshal(c, &dependencies)
	if err != nil {
		panic(err)
	}
	dependencies.LastUpdated = time.Now().Unix()
	if isDep {
		dependencies.Deps = append(dependencies.Deps, deps{Name: pkg, ReliedBy: installedBy, InstalledByUser: false})
		c, err = json.Marshal(dependencies)
		if err != nil {
			panic(err)
		}
		os.WriteFile("dependencies.json", []byte(c), 0644)
	} else {
		d := deps{Name: pkg, ReliedBy: "", InstalledByUser: true}
		dependencies.Deps = append(dependencies.Deps, d)
		c, err = json.Marshal(dependencies)
		if err != nil {
			panic(err)
		}
		os.WriteFile("dependencies.json", []byte(c), 0644)
	}

}
func checkIfDepExists(pkg string) (bool, *deps) {
	location, _ := os.Executable()
	location = location[:len(location)-len("/ferment")]
	os.Chdir(location)
	c, err := os.ReadFile("dependencies.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			os.WriteFile("dependencies.json", []byte("{}"), 0777)
			_, err = os.ReadFile("dependencies.json")
			if err != nil {
				panic(err)
			}
			return false, nil
		} else {
			panic(err)
		}
	}
	var dependencies Dep
	err = json.Unmarshal(c, &dependencies)
	if err != nil {
		panic(err)
	}
	for _, d := range dependencies.Deps {
		if d.Name == pkg {
			return true, &d
		}
	}
	return false, nil
}
func EditDepsTracker(pkg string, roots string) {
	location, _ := os.Executable()
	location = location[:len(location)-len("/ferment")]
	os.Chdir(location)
	c, err := os.ReadFile("dependencies.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			os.WriteFile("dependencies.json", []byte("{}"), 0777)
			_, err = os.ReadFile("dependencies.json")
			if err != nil {
				panic(err)
			}
			return
		} else {
			panic(err)
		}
	}
	var dependencies Dep
	err = json.Unmarshal(c, &dependencies)
	if err != nil {
		panic(err)
	}
	for i, d := range dependencies.Deps {
		if d.Name == pkg {
			old := d.ReliedBy
			dependencies.Deps[i].ReliedBy = fmt.Sprintf("%s %s", old, roots)
		}
	}
	c, err = json.Marshal(dependencies)
	if err != nil {
		panic(err)
	}
	os.WriteFile("dependencies.json", []byte(c), 0644)
}
func InstallPrebuilds(pkg string) {
	os.Chdir(location)
	c, err := os.ReadFile(fmt.Sprintf("Barrells/%s.py", convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		color.RedString("ERROR: %s", err)
		os.Exit(1)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(c)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "pkg.prebuild.install()\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	f, err := os.OpenFile("/tmp/ferment-install.log", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		color.RedString("ERROR: %s\n", err)
	}
	io.Copy(f, r)
}
func getSaidDeps(pkg string) []string {
	os.Chdir(location)
	c, err := os.ReadFile(fmt.Sprintf("Barrells/%s.py", convertToReadableString(strings.ToLower(pkg))))
	if err != nil {
		color.RedString("ERROR: %s", err)
		os.Exit(1)
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(c)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "print(pkg.dependencies)\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	dependencies := strings.Split(buf.String(), "\n")[0]
	dependencies = strings.Replace(dependencies, "[", "", 1)
	dependencies = strings.Replace(dependencies, "]", "", 1)
	dependencies = strings.Replace(dependencies, " ", "", -1)
	dependencies = strings.Replace(dependencies, "'", "", -1)
	dependencies = strings.Replace(dependencies, `"`, "", -1)
	return strings.Split(dependencies, ",")

}
func checkifPrebuildSuitable(pkg string) bool {
	arch := strings.ToLower(runtime.GOARCH)
	c, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(pkg)))
	if err != nil {
		color.RedString("ERROR: %s", err)
		os.Exit(1)
	}
	content := string(c)
	lines := strings.Split(content, "\n")
	var isPrebuild bool
	for _, line := range lines {
		if !strings.ContainsAny(line, "=") {
			continue
		}
		l := strings.Split(line, "=")
		if strings.Contains(l[0], "self.prebuild") && !strings.Contains(l[0], "None") {
			isPrebuild = true
			break
		}

	}
	if !isPrebuild {
		return false
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(c)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "print(pkg.prebuild.amd64)\n")
	io.WriteString(closer, "print(pkg.prebuild.arm64)\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	arm64 := strings.Split(buf.String(), "\n")[1]
	amd64 := strings.Split(buf.String(), "\n")[0]
	if arch == "arm64" && strings.Contains(arm64, "no attribute") {
		return false
	}

	if arch == "amd64" && strings.Contains(amd64, "no attribute") {
		return false
	}
	return true

}
func getPrebuildURL(pkg string) *string {
	arch := strings.ToLower(runtime.GOARCH)
	c, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(pkg)))
	if err != nil {
		color.RedString("ERROR: %s", err)
		os.Exit(1)
	}
	content := string(c)
	lines := strings.Split(content, "\n")
	var isPrebuild bool
	for _, line := range lines {
		if !strings.ContainsAny(line, "=") {
			continue
		}
		l := strings.Split(line, "=")
		if strings.Contains(l[0], "self.prebuild") && !strings.Contains(l[0], "None") {
			isPrebuild = true
			break
		}

	}
	if !isPrebuild {
		return nil
	}
	cmd := exec.Command("python3")
	closer, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	r, w, _ := os.Pipe()
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	cmd.Start()
	closer.Write(c)
	closer.Write([]byte("\n"))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, "print(pkg.prebuild.amd64)\n")
	io.WriteString(closer, "print(pkg.prebuild.arm64)\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	arm64 := strings.Split(buf.String(), "\n")[1]
	amd64 := strings.Split(buf.String(), "\n")[0]
	if arch == "arm64" && strings.Contains(arm64, "no attribute") {
		return nil
	}

	if arch == "amd64" && strings.Contains(amd64, "no attribute") {
		return nil
	}
	if arch == "amd64" {
		return &amd64
	}
	if arch == "arm64" {
		return &arm64
	}
	return nil
}
