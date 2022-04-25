/*
Copyright © 2022 NotTimIsReal

*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
const repo string = "https://github.com/NotTimIsReal/ferment"

var rootCmd = &cobra.Command{
	Use:   "ferment",
	Short: "A faster and more efficient way to install packages",
	Long: `
Ferment is a faster and more efficient way to install packages.
Uses a similar concept to brew but much faster as it uses a 
compiled language rather than Ruby.
Run ferment install your first package.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		v, err := cmd.Flags().GetBool("version")
		if err != nil {
			panic(err)
		}
		if v {
			fs := memfs.New()
			_, err = git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
				URL: repo,
			})
			if err != nil {
				panic(err)
			}
			f, err := fs.Open("VERSION.meta")
			if err != nil {
				panic(err)
			}
			var buf bytes.Buffer
			io.Copy(&buf, f)
			location, err := os.Executable()
			if err != nil {
				panic(err)
			}
			location = location[:len(location)-len("/ferment")]
			content, err := os.ReadFile(fmt.Sprintf("%s/VERSION.meta", location))
			if err != nil {
				panic(err)
			}
			version := string(content)
			fmt.Printf("Ferment %s\n", version)
			fmt.Printf("Latest Version %s\n", buf.String())
			os.Exit(0)
		}
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ferment.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Prints the version of the Ferment")
}
