/*
Copyright © 2022 NotTimIsReal
*/
package cmd

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"net/http"
	"net/url"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/radovskyb/watcher"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/theckman/yacspin"
	spinner "github.com/theckman/yacspin"
	"github.com/ulikunitz/xz"
)

var l, _ = os.Executable()
var location = l[:len(l)-len("/ferment")]

type pkg struct {
	name         string
	version      string
	alias        []string
	desc         string
	dependencies []string
	Dbuild       []string
	arch         []string
	source       []string
	build        string
	install      string
	test         string
	caveats      string
	license      string
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <package>",
	Short: "Install Packages",
	Long:  `Install Official Barrells `,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		command := exec.Command("ferment", "search")
		var out bytes.Buffer
		command.Stdout = &out
		command.Run()
		pkgs := out.String()
		pkgArr := strings.Split(pkgs, "\n")
		pkgArr = pkgArr[1:]

		return pkgArr, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		buildfromsource, err := cmd.Flags().GetBool("build-from-source")
		if err != nil {
			panic(err)
		}
		nocache, err := cmd.Flags().GetBool("no-cache")
		os.Setenv("FERMENT_NO_CACHE", fmt.Sprintf("%v", nocache))
		if err != nil {
			panic(err)
		}
		if len(args) == 0 {
			fmt.Println("Please provide a package to install, it can either be a custom package from github, gitlab, etc or a official package")
			os.Exit(1)
		}
		for _, pkg := range args {
			var foundPkg bool = false
			verbose, err := cmd.Flags().GetBool("verbose")

			if err != nil {

				panic(err)
			}
			location, err := os.Executable()
			//redefine location so that it is the directory of the executable
			location = location[:len(location)-len("ferment")]
			if err != nil {

				panic(err)
			}
			type pkgI struct {
				LatestVersion string   `json:"latestVersion"`
				AllFiles      []string `json:"allFiles"`
			}
			pkg = convertToReadableString(strings.ToLower(pkg))
			if strings.Contains(pkg, "@") {
				//set env
				os.Setenv("FERMENT_PKG_VERSION", strings.Split(pkg, "@")[1])
				pkg = strings.Split(pkg, "@")[0]
			} else {
				resp, err := http.Get(fmt.Sprintf("https://api.fermentpkg.tech/barrells/info/%s", pkg))
				if err != nil {
					color.Red("Error getting package info")
					os.Exit(1)
				}
				defer resp.Body.Close()

				var pkgInfo pkgI
				json.NewDecoder(resp.Body).Decode(&pkgInfo)
				os.Setenv("FERMENT_PKG_VERSION", pkgInfo.LatestVersion)
			}
			//search for package in default list
			if verbose {
				fmt.Println("Searching for package in default list")
			}
			files, err := os.ReadDir(fmt.Sprintf("%s/Barrells", location))
			if err != nil {
				panic(err)
			}

			for _, v := range files {
				name := strings.ToLower(v.Name())
				if strings.Split(name, ".")[0] == strings.ToLower(strings.Join([]string{pkg, os.Getenv("FERMENT_PKG_VERSION")}, "@")) {
					if verbose {
						fmt.Println("Found package in default list")
					}
					foundPkg = true
					break
				}
			}

			if e, _ := exists(fmt.Sprintf("%s/Installed/%s", location, pkg)); e {
				color.Red("Package %s already installed", pkg)
				continue
			}
			c, err := os.ReadFile(location + "ferment.lock")
			if err == nil {
				content := string(c)
				color.Red("Lock file exists, PID %s is using the current lockfile", content)
				os.Exit(1)
			}
			err = os.WriteFile(location+"ferment.lock", []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
			if !foundPkg {
				s, err := spinner.New(spinner.Config{
					Frequency:         100 * time.Millisecond,
					CharSet:           spinner.CharSets[57],
					Suffix:            color.GreenString(" Prebuild"),
					SuffixAutoColon:   true,
					Message:           "Download",
					StopCharacter:     "✓",
					StopColors:        []string{"fgGreen"},
					StopFailCharacter: "✗",
					StopFailColors:    []string{"fgRed"},
				}) // Build our new spinner
				if err != nil {
					panic(err)
				}
				s.Start()
				resp, err := http.Get(fmt.Sprintf("https://api.fermentpkg.tech/barrells/info/%s/", pkg))
				if err != nil {
					color.Red("Error getting package info")
					os.Exit(1)
				}
				defer resp.Body.Close()
				type pkgI struct {
					LatestVersion string   `json:"latestVersion"`
					AllFiles      []string `json:"allFiles"`
					Universal     bool     `json:"universal"`
				}
				var pkgInfo pkgI
				json.NewDecoder(resp.Body).Decode(&pkgInfo)
				if pkgInfo.Universal {
					s.Message("Downloading From API (universal)")
					prebuildDownloadFromAPI(pkg, fmt.Sprintf("%s@%s.ferment", pkg, os.Getenv("FERMENT_PKG_VERSION")), s)

				} else {
					s.Message("Downloading From API (single-arch)")
					prebuildDownloadFromAPI(pkg, fmt.Sprintf("%s@%s.%s.ferment", pkg, os.Getenv("FERMENT_PKG_VERSION"), runtime.GOARCH), s)
				}

				s.Stop()
			}

			if err != nil {
				panic(err)
			}
			if !buildfromsource && checkifPrebuildSuitable(pkg) {

				installPackages(pkg, verbose, false, "", false)
				os.Exit(0)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().BoolP("verbose", "v", false, "Log All Output")
	installCmd.PersistentFlags().Bool("no-cache", false, "Not to use cached downloads")
	installCmd.PersistentFlags().BoolP("build-from-source", "b", false, "Build From Source or use an available prebuild")

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
func DownloadFromGithub(url string, path string, verbose bool) error {
	if verbose {
		fmt.Println("Downloading from Github")
	}
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			return fmt.Errorf("package already exists")
		}
		return err
	}
	if verbose {
		fmt.Println("Downloaded from Github")
	}
	return nil
}
func UsingGit(pkg string, verbose bool) bool {
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
	if verbose {
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

// TODO: Create the ExtractToMemory function
func extractFerment(path string) (filesystem afero.Fs, err error) {
	fs := afero.NewMemMapFs()
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	xzReader, err := xz.NewReader(file)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(xzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := fs.MkdirAll(fmt.Sprintf(header.Name), 0755); err != nil {
				return nil, err
			}
		case tar.TypeReg:
			f, err := fs.Create(header.Name)
			if err != nil {
				return nil, err
			}
			if _, err := io.Copy(f, tarReader); err != nil {
				return nil, err
			}
			f.Close()
		default:
			return nil, fmt.Errorf("unable to untar type: %c in file %s", header.Typeflag, header.Name)
		}
	}

	return fs, nil
}
func installPackages(packageName string, verbose bool, isDep bool, installedBy string, buildFromSource bool) {
	if verbose {
		fmt.Println("Looking For Dependencies")
	}
	//Variables
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}
	location = location[:len(location)-len("/ferment")]
	name := strings.Split(convertToReadableString(strings.ToLower(packageName)), "@")[0]
	version := strings.Split(convertToReadableString(strings.ToLower(packageName)), "@")[1]
	_, err = os.ReadFile(fmt.Sprintf("%s/Barrells/%s.ferment", location, name))

	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println(color.RedString("Reinstall ferment, Barrells is missing"))
			os.Exit(1)
		}
		panic(err)
	}

	// pkg.ferment is a tar.xz file
	// extract to fs
	fs, err := extractFerment(fmt.Sprintf("%s/Barrells/%s.ferment", location, name))
	if err != nil {
		panic(err)
	}
	fpkg, err := fs.Open(fmt.Sprintf("%s.pkg", name))
	if err != nil {
		panic(err)
	}
	defer fpkg.Close()
	os.MkdirAll(fmt.Sprintf("/tmp/ferment/%s/%s", name, version), 0755)
	fpkgOut, err := os.OpenFile(fmt.Sprintf("/tmp/ferment/%s/%s/%s.fpkg", name, version, name), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	io.Copy(fpkgOut, fpkg)
	pkg := parseFpkg(fmt.Sprintf("/tmp/ferment/%s/%s/%s.fpkg", name, version, name))
	for _, dep := range pkg.dependencies {
		color.Yellow("Package %s depends on %s", pkg.name, dep)
		if exec.Command("which", dep).Wait() == nil {
			color.Yellow("%s is already installed or already exists", dep)
			EditDepsTracker(dep, pkg.name)
			color.Yellow("Skipping")
			continue
		}
		if IsLib(dep) && checkIfPackageExists(dep) {
			color.Yellow("%s is a library and is already installed", dep)
			EditDepsTracker(dep, pkg.name)
			color.Yellow("Skipping")
			continue
		}
		s, err := spinner.New(spinner.Config{
			Frequency:       100 * time.Millisecond,
			CharSet:         spinner.CharSets[57],
			Suffix:          color.GreenString(" %s", dep),
			SuffixAutoColon: true,
			Message:         "Getting Download Info",
			StopCharacter:   "✓",
			StopColors:      []string{"fgGreen"},
			StopFailMessage: "Failed",
			StopFailColors:  []string{"fgRed"},
		})
		if err != nil {
			panic(err)
		}
		s.Start()
		// TODO Fix the build from source
		if !buildFromSource {
			prebuildDownloadFromAPI(dep, getFileFromLink(*getPrebuildURL(dep)), s)
			s.Stop()
			installPackages(dep, verbose, true, pkg.name, buildFromSource)
		} else {
			GetDownloadUrl(dep, verbose, s)
			s.Stop()
			installPackages(dep, verbose, true, pkg.name, buildFromSource)
			// TestInstallationScript(dep, verbose)
		}

	}

	if e, _ := checkIfDepExists(packageName); e {
		EditDepsTracker(packageName, installedBy)
	} else {
		DepTrackerAdd(packageName, isDep, installedBy)
	}
	s, err := spinner.New(spinner.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         spinner.CharSets[57],
		Suffix:          " Installing " + packageName,
		SuffixAutoColon: true,
		Message:         "Install",
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
		StopFailMessage: "Failed",
		StopFailColors:  []string{"fgRed"},
	})
	if err != nil {
		panic(err)
	}
	s.Message(" Installing " + packageName)
	s.StopMessage(color.GreenString("Installed %s", packageName))

	s.Start()
	if checkifPrebuildSuitable(packageName) && !buildFromSource {
		InstallPrebuilds(packageName)
	} else {
		RunInstallationScript(convertToReadableString(strings.ToLower(packageName)), verbose, convertToReadableString(strings.ToLower(packageName)))
	}
	version, _ = getVersion(packageName)
	writeVersionFile(packageName, version)
	s.Stop()
	s, err = spinner.New(spinner.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           spinner.CharSets[57],
		Suffix:            " Binary",
		SuffixAutoColon:   true,
		Message:           "Install",
		StopColors:        []string{"fgGreen"},
		StopCharacter:     "✓",
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
	})
	if err != nil {
		panic(err)
	}
	s.Message(" Installing Binaries For " + packageName)
	s.Start()
	msg := InstallBinary(packageName, verbose)
	if msg == "No Binary" {
		s.StopMessage(color.GreenString("No Binaries To Be Installed"))

	} else {
		s.StopMessage(color.GreenString("Installed Binary for %s", packageName))
	}
	s.Stop()

	s, err = spinner.New(spinner.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           spinner.CharSets[57],
		Suffix:            " Testing",
		SuffixAutoColon:   true,
		Message:           " Testing " + packageName,
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
	})
	if err != nil {
		panic(err)
	}
	s.Start()
	result := TestInstallationScript(packageName, verbose)
	if result {
		s.Stop()
	} else {
		s.StopFail()
	}
	if caveats := getCaveats(packageName); caveats != nil {
		fmt.Println(color.YellowString("Caveats:\n"), *caveats)
	}
	os.Remove(location + "/ferment.lock")

}
func GetGitURL(pkg string, verbose bool) string {
	if verbose {
		fmt.Println("Looking For GitURl")
	}
	//Variables
	location, err := os.Executable()
	if err != nil {
		panic(err)
	}

	location = location[:len(location)-len("/ferment")]
	if path := checkifFpkgExists(pkg); pkg != "" {
		pkg := parseFpkg(path)
		var git string
		for _, source := range pkg.source {
			//check if source ends with .git
			if strings.HasSuffix(source, ".git") {
				git = source
				break
			}
		}
		return git

	} else {
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
		if verbose {
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
}
func DownloadFromTar(pkg string, url string, verbose bool, spinner *spinner.Spinner) string {
	resp, err := http.Get(url)
	if err != nil {
		spinner.StopFailMessage(color.WhiteString(err.Error()))
		spinner.StopFail()
		os.Exit(1)
	}
	//dont keep-alive
	resp.Header.Set("Connection", "close")
	location, err := os.Executable()
	location = location[:len(location)-len("/ferment")]
	if err != nil {
		spinner.StopFailMessage(color.WhiteString(err.Error()))
		spinner.StopFail()
		os.Exit(1)
	}
	fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	defer resp.Body.Close()
	if verbose {
		fmt.Printf("Downloading Tar From %s\n", url)
	}
	//check if file exists in /tmp/ferment
	if os.Getenv("FERMENT_NO_CACHE") == "true" {
		os.Remove(fmt.Sprintf("%s/Barrells/%s", location, fileName))
		spinner.Message("Removed Old Cache File")
	}
	if _, err := os.Stat(fmt.Sprintf("%s/Barrells/%s", location, fileName)); err == nil {
		spinner.Message(fmt.Sprintf("Using Cached %s", pkg))

	} else {
		os.MkdirAll(fmt.Sprintf("%s/Barrells/", location), 0777)
		f, err := os.OpenFile(fmt.Sprintf("%s/Barrells/%s", location, fileName), os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			spinner.StopFailMessage("Failed Creating Tar")
			spinner.StopFail()
			os.Exit(1)
		}
		progress := make(chan float64)
		go func() {
			for {
				p := getDownloadProgress(fmt.Sprintf("%s/Barrells/%s", location, fileName), resp.ContentLength, resp)
				progress <- p
				if p == 100 {
					break
				}

			}
		}()
		go func() {
			_, err = io.Copy(f, resp.Body)
			if err != nil {
				spinner.StopFailMessage("Failed Writing To Tar")
				spinner.StopFail()
			}
		}()
		sigChan := make(chan os.Signal)
		go func() {
			signal.Notify(sigChan, syscall.SIGINT)
			<-sigChan
			spinner.Message("Download Cancelled Cleaning Up...")
			os.Remove(fmt.Sprintf("%s/Barrells/%s", location, fileName))
			spinner.StopFailMessage("Download Cancelled")
			spinner.StopFail()
			os.Exit(1)
		}()

		for {

			p := <-progress
			if p == -1 {
				spinner.Message(fmt.Sprintf("Downloading %s Of Unknown Size", pkg))
			} else {
				spinner.Message(fmt.Sprintf("Downloading %s %d%s/100%s", pkg, int(p), "%", "%"))
			}
			if p == 100 {
				break
			}
			time.Sleep(500 * time.Millisecond)

		}

	}
	if verbose {
		fmt.Println("Extracting Tar")
	}
	spinner.Message("Extracting Tar... (This Might Take a While)")
	err = Untar(fmt.Sprintf("%s/Installed/", location), fmt.Sprintf("%s/Barrells/%s", location, fileName), pkg)

	if err != nil {
		spinner.StopFailMessage(err.Error())
		spinner.StopFail()
		os.Exit(1)
	}
	return fmt.Sprintf("%s/Installed/%s", location, pkg)
}
func GetDownloadUrl(pkg string, verbose bool, s *spinner.Spinner) string {
	if verbose {
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
	if verbose {
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
	path := DownloadFromTar(convertToReadableString(strings.ToLower(pkg)), strings.Replace(buf.String(), "\n", "", -1), verbose, s)
	return path
}
func Untar(dst string, downloadedFile string, pkg string) error {
	os.Mkdir(dst, 0777)
	//list dst
	oldentries, err := os.ReadDir(dst)
	if err != nil {
		return err
	}
	cmd := exec.Command("tar", "-xvf", downloadedFile, "--directory", dst)

	var bytes bytes.Buffer
	cmd.Stderr = &bytes
	err = cmd.Run()

	if err != nil {
		return errors.New(bytes.String())
	}
	newentries, err := os.ReadDir(dst)
	if err != nil {
		return err
	}
	//find the difference between the two
	if len(oldentries) == 0 && len(newentries) > 0 {
		os.Rename(fmt.Sprintf("%s/%s", dst, newentries[0].Name()), fmt.Sprintf("%s/%s", dst, pkg))
	} else {
		//Using the old entries, find the first one that is not in the old entries
		for _, entry := range newentries {
			found := false
			for _, oldentry := range oldentries {
				if entry.Name() == oldentry.Name() {
					found = true
					break
				}
			}
			if !found {
				os.Rename(fmt.Sprintf("%s/%s", dst, entry.Name()), fmt.Sprintf("%s/%s", dst, pkg))
				break
			}
		}
	}

	return nil
}
func RunInstallationScript(pkg string, verbose bool, cwd string) {
	if verbose {
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
	if verbose {
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
	done := make(chan bool)
	go magicWatcher(pkg, done)
	io.WriteString(closer, "pkg.install()\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	done <- true
	var buf bytes.Buffer
	io.Copy(&buf, r)
	//print nothing
	buf.Reset()

}
func TestInstallationScript(pkg string, verbose bool) bool {
	if verbose {
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
	if verbose {
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
func InstallBinary(pkg string, verbose bool) string {
	if verbose {
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
	if verbose {
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
		binary = strings.Replace(binary, `"`, "", -1)
		binary = strings.Replace(binary, "\n", "", -1)
		os.Remove(fmt.Sprintf("/usr/local/bin/%s", binary))
		os.Symlink(fmt.Sprintf("%s/Installed/%s/%s", location, pkg, binary), fmt.Sprintf("/usr/local/bin/%s", binary))
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
		os.WriteFile("dependencies.json", []byte(c), 0777)
	} else {
		d := deps{Name: pkg, ReliedBy: "", InstalledByUser: true}
		dependencies.Deps = append(dependencies.Deps, d)
		c, err = json.Marshal(dependencies)
		if err != nil {
			panic(err)
		}
		os.WriteFile("dependencies.json", []byte(c), 0777)
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
	os.WriteFile("dependencies.json", []byte(c), 0777)
}
func InstallPrebuilds(pkg string) {
	os.Chdir(location)
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
	io.WriteString(closer, fmt.Sprintf("from %s import %s\n", pkg, pkg))
	io.WriteString(closer, fmt.Sprintf("pkg=%s()\n", convertToReadableString(strings.ToLower(pkg))))
	io.WriteString(closer, fmt.Sprintf(`pkg.prebuild.cwd="%s/Installed/%s"`+"\n", location, pkg))
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
func checkIfSetupPkg(pkg string) bool {
	c, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(pkg)))
	if err != nil {
		color.RedString("ERROR: %s", err)
		os.Exit(1)
	}
	content := string(c)
	lines := strings.Split(content, "\n")
	var isSetup bool = false
	for _, line := range lines {
		if !strings.ContainsAny(line, "=") {
			continue
		}
		l := strings.Split(line, "=")
		if strings.Contains(l[0], "self.setup") && strings.Contains(l[1], "True") {
			isSetup = true
			break
		}

	}
	return isSetup
}
func installPackageWithSetup(pkg string, spinner *yacspin.Spinner) {
	c, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(pkg)))
	if err != nil {
		spinner.Message(color.RedString("ERROR: %s", err))
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
	io.WriteString(closer, `print(f"url:{pkg.url}")`+"\n")
	closer.Close()
	w.Close()
	cmd.Wait()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	content := strings.Split(buf.String(), "\n")
	var i int
	for index, line := range content {
		s := strings.Replace(line, "url:", "", 1)
		if IsUrl(s) {
			i = index
			break
		}
	}
	url := content[i]
	var sh bool
	var interactive bool
	if strings.Contains(url, "sh:") {
		sh = true
	}
	if strings.Contains(url, "interactive:") {
		interactive = true
	}
	url = strings.Replace(url, "sh:", "", 1)
	url = strings.Replace(url, "interactive:", "", 1)
	url = strings.Replace(url, "url:", "", 1)
	os.Mkdir(fmt.Sprintf("%s/Installed/%s", location, convertToReadableString(pkg)), 0755)
	if sh {
		spinner.Message(color.GreenString("Package %s is to be setup with sh and curl...", pkg))

		cmd = exec.Command("curl", "-sSLo", fmt.Sprintf("/tmp/%s-setup.sh", pkg), url)
		cmd.Dir = fmt.Sprintf("%s/Installed/%s", location, pkg)
		err = cmd.Run()
		if err != nil {
			color.Red("ERROR - CURL: %s", err)
			os.Exit(1)
		}
		cmd = exec.Command("sh", fmt.Sprintf("/tmp/%s-setup.sh", pkg))
		cmd.Dir = fmt.Sprintf("%s/Installed/%s", location, pkg)
		if interactive {
			spinner.Pause()
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
		}
		err = cmd.Run()
		if err != nil {
			spinner.Message(color.RedString("ERROR - SH: %s", err))
			os.Exit(1)
		}

	}
	spinner.Stop()

}
func prebuildDownloadFromAPI(pkg string, file string, s *spinner.Spinner) {
	url := fmt.Sprintf("https://api.fermentpkg.tech/barrells/download/%s/%s", pkg, file)
	DownloadFromTar(pkg, url, false, s)

}

type body struct {
	LatestVersion string   `json:"latestVersion"`
	AllFiles      []string `json:"allFiles"`
}

func getLatestVersionOfPrebuild(pkg string) body {
	res, err := http.Get(fmt.Sprintf("https://api.fermentpkg.tech/barrells/info/%s", pkg))
	if err != nil {
		color.Red("ERROR: %s", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	var body body
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		color.Red("ERROR: %s", err)
		os.Exit(1)
	}
	return body
}

// Returns the links that wants to be downloaded or nil if is not prebuildapi
// Works well with the prebuildDownloadFromAPI function
func checkIfPrebuildApi(pkg string) (data *struct {
	arm64           string
	amd64           string
	usingFermentTag struct {
		amd64 bool
		arm64 bool
	}
}, Error error) {
	c, err := os.ReadFile(fmt.Sprintf("%s/Barrells/%s.py", location, convertToReadableString(pkg)))
	if err != nil {
		return nil, err
	}
	content := string(c)
	lines := strings.Split(content, "\n")
	var amd64Download string
	var arm64Download string
	for _, line := range lines {
		if strings.Contains(line, "self.arm64") {
			f := strings.Split(line, "=")
			arm64Download = strings.Replace(f[1], `"`, "", -1)
			arm64Download = strings.Replace(arm64Download, `'`, "", -1)

		}
		if strings.Contains(line, "self.amd64") {
			f := strings.Split(line, "=")
			amd64Download = strings.Replace(f[1], `"`, "", -1)
			amd64Download = strings.Replace(amd64Download, `'`, "", -1)

		}
	}
	if amd64Download != "" || arm64Download != "" {
		d := struct {
			amd64 bool
			arm64 bool
		}{}
		for i, v := range []string{amd64Download, arm64Download} {
			if strings.Contains(v, "ferment://") {
				if i == 0 {
					d.amd64 = true
				} else {
					d.arm64 = true
				}
			}
		}
		return &struct {
			arm64           string
			amd64           string
			usingFermentTag struct {
				amd64 bool
				arm64 bool
			}
		}{arm64Download, amd64Download, d}, nil
	}

	return nil, errors.New("no prebuild api found")
}

// This function allows for you to extract the tar from the download link
// Example Download Link: ferment://<pkg>@<file>
func getFileFromLink(link string) string {
	l := strings.Split(link, "@")
	return l[1]
}
func getCaveats(pkg string) *string {
	if caveats, err := executeQuickPython(fmt.Sprintf("import %s;pkg=%s.%s();print(pkg.caveats)", convertToReadableString(pkg), convertToReadableString(pkg), convertToReadableString(pkg))); err != nil {
		return nil
	} else {
		return &caveats
	}
}
func executeQuickPython(code string) (string, error) {
	cmd := exec.Command("python3", "-c", code)
	cmd.Dir = fmt.Sprintf("%s/Barrells", location)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", errors.New(out.String())
	}
	return out.String(), nil

}

// work on later
func getFileSize(pkg string) int {
	version := os.Getenv("FERMENT_PKG_VERSION")

	if version == "" {
		type returnBody struct {
			StatusCode int
			Body       struct {
				LatestVersion string   `json:"latesVersion"`
				AllFiles      []string `json:"allFiles"`
			}
		}
		res, err := http.Get(fmt.Sprintf("https://api.fermentpkg.tech/barrells/info/%s", pkg))
		if err != nil {
			panic(err)
		}
		var body returnBody
		res.StatusCode = body.StatusCode
		//un marshall json
		json.NewDecoder(res.Body).Decode(&body.Body)
		version = body.Body.LatestVersion

	}

	res, err := http.Get(fmt.Sprintf("https://api.fermentpkg.tech/barrells/info/%s/%s", pkg, version))
	if err != nil {
		panic(err)
	}
	var body struct {
		StatusCode int
		Body       struct {
			Size          int      `json:"fileSize"`
			LatestVersion string   `json:"latestVersion"`
			AllFiles      []string `json:"allFiles"`
		}
	}
	res.StatusCode = body.StatusCode
	//un marshall json
	if body.StatusCode > 200 {
		return 0
	}
	json.NewDecoder(res.Body).Decode(&body.Body)

	return body.Body.Size
}
func getDownloadProgress(file string, total int64, r *http.Response) float64 {
	fileO, err := os.Open(file)
	if err != nil {
		return 0
	}
	defer fileO.Close()
	fi, err := fileO.Stat()
	if err != nil {
		return 0
	}
	//check if total is -1
	if total == -1 {
		select {
		case <-r.Request.Context().Done():
			return 100
		default:
			return -1
		}
	}
	return float64(fi.Size()) / float64(total) * 100
}
func writeVersionFile(pkg string, version string) {
	f, err := os.OpenFile(fmt.Sprintf("%s/Installed/%s/VERSION.meta", location, convertToReadableString(pkg)), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(version)
	if err != nil {
		panic(err)
	}
}
func magicWatcher(pkg string, done chan bool) {

	watch := watcher.New()
	dirsWatched := []string{"bin", "share", "include", "lib"}
	for _, dir := range dirsWatched {
		watch.Add(fmt.Sprintf("/usr/local/%s", dir))
	}
	watch.Add("/usr/local")
	go watch.Start(10 * time.Millisecond)
	watcherfile, err := os.OpenFile(fmt.Sprintf("%s/Installed/%s/.ferment_watcher", location, pkg), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case event := <-watch.Event:
				//ger file name
				if event.Op == watcher.Create || event.Op == watcher.Write && strings.Contains(event.Path, pkg) {
					watcherfile.WriteString(event.Path + "\n")
				}
			case <-watch.Closed:
				return
			}
		}
	}()
	for {
		d := <-done
		if d {
			break
		}
	}
	watch.Close()
	watcherfile.Close()

}
func parseFpkg(file string) pkg {
	c, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	//parse the file

	content := string(c)
	var pkg pkg
	arrayContent := strings.Split(content, "\n")
	var skipToIndex int = 0
	for i, line := range arrayContent {
		if i < skipToIndex {
			continue
		}
		lineCut := strings.Split(line, "=")
		if len(lineCut) > 1 {
			lineCut[1] = strings.ReplaceAll(lineCut[1], "\"", "")
			switch lineCut[0] {
			case "pkgname":
				pkg.name = lineCut[1]
			case "version":
				pkg.version = lineCut[1]
			case "desc":
				pkg.desc = lineCut[1]
			case "alias":
				pkg.alias = strings.Split(lineCut[1], ",")
			case "arch":
				pkg.arch = strings.Split(lineCut[1], ",")
			case "dependencies":
				pkg.dependencies = strings.Split(lineCut[1], ",")
			case "Dbuild":
				pkg.Dbuild = strings.Split(lineCut[1], ",")
			case "source":
				pkg.source = strings.Split(lineCut[1], ",")
			case "license":
				pkg.license = lineCut[1]
			case "caveats":
				pkg.caveats = lineCut[1]

			}
		} else if strings.Contains(line, "()") {
			//starting from index i look for a }
			var indexOfEnd int
			for i, line := range arrayContent[i:] {
				if strings.Contains(line, "}") {
					indexOfEnd = i
					break
				}

			}
			//get the function name
			functionName := strings.Split(line, "(")[0]
			//get the function body
			functionBody := strings.Join(arrayContent[i:i+indexOfEnd], "\n")
			//add the function to the pkg
			switch functionName {
			case "build":
				pkg.build = functionBody
			case "install":
				pkg.install = functionBody
			case "test":
				pkg.test = functionBody

			}
		}
		skipToIndex++
	}
	//check if every field in pkg is filled
	if pkg.name == "" || pkg.version == "" || pkg.desc == "" || pkg.arch == nil || pkg.source == nil || pkg.build == "" || pkg.install == "" || pkg.test == "" {
		panic("Invalid fpkg file")
	}
	return pkg
}
func checkIfFpkg(filename string) bool {
	if strings.Contains(filename, ".fpkg") {
		return true
	}
	return false
}
func checkifFpkgExists(pkg string) string {
	//look through barrells directory
	files, err := os.ReadDir(fmt.Sprintf("%s/Barrels", location))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		//check if the file exists
		if _, err := os.Stat(fmt.Sprintf("%s/Barrels/%s/%s.fpkg", location, file.Name(), pkg)); err == nil {
			return fmt.Sprintf("%s/Barrels/%s/%s.fpkg", location, file.Name(), pkg)
		}
	}
	return ""
}
func getPrebuildUrlFromFpkg(pkg string) string {
	//look through barrells directory
	files, err := os.ReadDir(fmt.Sprintf("%s/Barrels", location))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		//check if the file exists
		if _, err := os.Stat(fmt.Sprintf("%s/Barrels/%s/%s.fpkg", location, file.Name(), pkg)); err == nil {
			//get the fpkg file
			fpkgFile := parseFpkg(fmt.Sprintf("%s/Barrels/%s/%s.fpkg", location, file.Name(), pkg))
			//get the url
			return fpkgFile.source[0]
		}
	}
	return ""
}
