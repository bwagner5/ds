package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"log"

	"github.com/bwagner5/ds/pkg/dataset"
	"github.com/spf13/cobra"
)

const version = "version"

// This is injected via ldflags on build in the Makefile
var versionID string

type userArgs struct {
	file string
}

func main() {
	args := userArgs{}
	rootCmd := &cobra.Command{
		Use:   "ds",
		Short: "ds is a CLI tool to compute stats for data sets",
		Run: func(cmd *cobra.Command, _ []string) {
			if f, _ := cmd.Flags().GetBool(version); f {
				printAndExitVersion()
			}
			ds(args)
		},
	}
	rootCmd.PersistentFlags().StringVarP(&args.file, "file", "f", "", "Input file to compute statistics for")
	rootCmd.PersistentFlags().BoolP(version, "v", false, "the version")
	rootCmd.Execute()
}

func ds(args userArgs) {
	file := os.Stdin
	if args.file != "" {
		var err error
		file, err = os.Open(args.file)
		if err != nil {
			fmt.Printf("Could not open the specified file: %v\n", err)
			os.Exit(-1)
		}
	}
	data, err := loadData(file)
	if err != nil {
		log.Printf("Error summarizing the numbers: %v", err)
		os.Exit(1)
	}
	fmt.Println(data.SummaryString(true))
}

func printAndExitVersion() {
	fmt.Println(versionID)
	os.Exit(0)
}

func loadData(file *os.File) (*dataset.DataSet, error) {
	r := bufio.NewReader(file)
	d := dataset.New()
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return d, err
		}
		// trim commas out of numbers
		line = strings.ReplaceAll(line, ",", "")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		num, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return d, err
		}
		d.Put(num)
	}
	return d, nil
}
