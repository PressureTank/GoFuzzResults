package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func doFileWork(path string, info os.FileInfo) error {
	fmt.Println(path + "  --  " + info.Name())
	return nil
}

func find(root, ext string) ([]string, []string) {
	var a []string
	var b []string

	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
			otherPath := strings.TrimSuffix(s, ext)
			_, err := os.Stat(otherPath)
			// errors.Is(err, os.ErrNotExist)
			if err != nil {
				log.Fatalf("expected file doesn't exist %s", otherPath)
			}
			b = append(b, otherPath)

		}

		return nil
	})
	return a, b
}

func main() {
	// create output file
	report, err := os.Create("crashers.csv")
	if err != nil {
		log.Fatalf("unable to create output file: %s", err)
	}
	defer report.Close()

	stacktraceFiles, crashValueFiles := find(os.Args[1], ".output")
	for i, s := range stacktraceFiles {
		// open s and copy the contents out of the file and add them to our
		// report file.
		crasherFile := crashValueFiles[i]
		file, err := os.Open(crasherFile)
		if err != nil {
			log.Fatalf("failed openeing file: %s", err)
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var txtlines []string

		for scanner.Scan() {
			txtlines = append(txtlines, scanner.Text())
		}

		file.Close()
		// report.WriteString("---------------")
		file, err = os.Open(s)
		if err != nil {
			log.Fatalf("failed openeing file: %s", err)
		}

		scanner = bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			txtlines = append(txtlines, scanner.Text())
		}

		file.Close()

		_, err = report.WriteString("\n\n-------Crasher Value Object: ")

		if err != nil {
			log.Fatalf("Failed to write string to file; %s", err)
		}

		for _, eachline := range txtlines {

			_, err = report.WriteString(eachline)

			if err != nil {
				log.Fatalf("Failed to write string to file; %s", err)
			}

			_, err = report.WriteString("\n")

			if err != nil {
				log.Fatalf("Failed to write string to file; %s", err)
			}
		}
		_ = report.Close
	}
}
