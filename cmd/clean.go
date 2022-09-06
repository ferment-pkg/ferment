/*
Copyright © 2022 Nottimisreal
*/
package cmd

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean The Caches For Ferment",
	Long:  `Removes downloaded caches and similar`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat("/tmp/ferment")
		if err != nil {
			color.Red("No Caches Found")
			os.Exit(0)
		}
		dirSize, _ := DirSize("/tmp/ferment")
		size, prefix := convertTolarger(dirSize)
		if dirSize < 1024 {
			color.Red("No Caches Found")
			os.Exit(0)
		}
		if dirSize/(1<<20) > 25 {

			color.Yellow("Are You Sure You Want To Delete %d%s of Caches? (y/n)", size, prefix)
			buf := bufio.NewReader(os.Stdin)
			sentence, err := buf.ReadString('\n')
			if err != nil {
				color.Red("Error Reading Input")
				os.Exit(1)
			}
			if sentence != "y\n" {
				os.Exit(1)
			}
		}

		os.RemoveAll("/tmp/ferment/")
		color.Green("Done. %d%s of Caches Deleted", size, prefix)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func convertTolarger(bytes int64) (int64, string) {
	if bytes < 1024 {
		return bytes, "B"
	} else if bytes < 1024*1024 {
		return bytes / 1024, "KB"
	} else if bytes < 1024*1024*1024 {
		return bytes / (1024 * 1024), "MB"
	} else {
		return bytes / (1024 * 1024 * 1024), "GB"
	}
}
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
