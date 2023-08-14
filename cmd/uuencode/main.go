package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pedroalbanese/uuencode"
)

var (
	dec   = flag.Bool("d", false, "Decode instead of Encode")
	ifile = flag.String("f", "", "Target file")
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage of", os.Args[0]+":")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *dec == false {
		var err error
		var infile *os.File
		infile, err = os.Open(*ifile)
		if err != nil {
			log.Println(err)
		}
		info, err := infile.Stat()
		if err != nil {
			log.Fatal(err)
		}
		uw := uuencode.NewWriter(os.Stdout, *ifile, info.Mode())
		if _, err = io.Copy(uw, infile); err != nil {
			return
		}
		if err := uw.Flush(); err != nil {
			return
		}
	} else {
		var err error
		var infile *os.File
		if *ifile == "-" || *ifile == "" {
			infile = os.Stdin
		} else {
			infile, err = os.Open(*ifile)
			if err != nil {
				log.Println(err)
			}
		}
		ur := uuencode.NewReader(infile, nil)

		outFileName, err := GetFileName(*ifile)
		if err != nil {
			fmt.Println("Error:", err)
		}
		outfile, err := os.Create(outFileName)
		if err != nil {
			log.Println(err)
		}
		_, err = io.Copy(outfile, ur)
		if err != nil {
			log.Fatal(err)
		}
		f, _ := ur.File()
		m, _ := ur.Mode()
		if err := outfile.Chmod(m); err != nil {
			log.Println(err)
		}
		fmt.Fprintf(os.Stderr, "file: %s, mode: %03o\n", f, m)
	}
}

func GetFileName(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "begin ") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[2], nil
			} else {
				return "", fmt.Errorf("Invalid 'begin' line format")
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("No 'begin' line found")
}
