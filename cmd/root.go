/*
Copyright © 2022 NotTimIsReal
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/qeesung/image2ascii/convert"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ferment",
	Short: "A faster and more efficient way to install packages",
	Long: fmt.Sprintf(`%s
Ferment is a faster and more efficient way to install packages.
Uses a similar concept to brew but much faster as it uses a
compiled language rather than Ruby.
Run ferment install your first package.`, convertAscii(location+"/images/logo.png")),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		v, err := cmd.Flags().GetBool("version")
		if err != nil {
			panic(err)
		}
		if v {
			resp, _ := http.Get("https://raw.githubusercontent.com/ferment-pkg/ferment/main/ferment.config.json")
			resp.Request.Header.Set("Cache-Control", "private, no-store, max-age=0")
			git, _ := io.ReadAll(resp.Body)

			var gitConfig Config
			json.Unmarshal(git, &gitConfig)

			location, err := os.Executable()
			if err != nil {
				panic(err)
			}
			location = location[:len(location)-len("/ferment")]
			config := getConfig(location)
			fmt.Printf("Ferment: %s\n", config.Version)
			if gitConfig.Version == "" {
				gitConfig.Version = "unknown"
			}
			fmt.Printf("Latest Version: %s\n", gitConfig.Version)
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
func convertAscii(imageFilename string) string {
	// Create convert options
	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = 100
	convertOptions.FixedHeight = 40
	convertOptions.Colored = true
	convertOptions.FitScreen = true
	if _, err := os.Open(imageFilename); err != nil {
		return ""
	}
	// Create the image converter
	converter := convert.NewImageConverter()
	return converter.ImageFile2ASCIIString(imageFilename, &convertOptions)
}
