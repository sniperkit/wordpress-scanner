package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jmhobbs/wordpress-scanner/shared"
)

var scanDirCmd = &cobra.Command{
	Use:   "scan-dir <plugin_name> <plugin_version> <plugin_directory>",
	Short: "Scans a plugin directory for corruption",
	Long:  "",
	Run: func (cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("You must provide the plugin name, the plugin version, and the directory to scan")
			os.Exit(1)
		} else if len(args) > 3 {
			fmt.Println("You gave too many arguments")
			os.Exit(1)
		}

		plugin := args[0]
		version := args[1]
		directory := args[2]

		scan := shared.NewScan(plugin, version)

		err := filepath.Walk(directory, scanFile(scan))
		if err != nil {
			log.Fatal(err)
		}

		/*
			bytes, err := scan.MarshalToBinary()
			if err != nil {
				log.Fatal(err)
			}

			err = ioutil.WriteFile("example.bin", bytes, 0644)
			if err != nil {
				log.Fatal(err)
			}
		*/

		json.NewEncoder(os.Stdout).Encode(scan)
	},
}

func scanFile(scan *shared.Scan) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".php") {
			f, err := os.Open(path)
			if err != nil {
				scan.AddErrored(path, err)
				return nil
			}

			hash, err := shared.GetHash(f)
			if err != nil {
				scan.AddErrored(path, err)
				return nil
			}

			scan.AddHashed(path, hash)
		}

		return nil
	}
}