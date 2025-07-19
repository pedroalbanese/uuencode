package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	
	"github.com/pedroalbanese/uuencode
)

var (
	dec = flag.Bool("d", false, "Decode instead of Encode")
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage of", os.Args[0]+":")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !*dec {
		// Modo codificação
		mw := uuencode.NewMultiWriter(os.Stdout)
		defer mw.Close()

		for _, filename := range flag.Args() {
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}

			info, err := file.Stat()
			if err != nil {
				file.Close()
				log.Fatal(err)
			}

			if err := mw.WriteFile(filename, info.Mode(), file); err != nil {
				file.Close()
				log.Fatal(err)
			}
			file.Close()
		}
	} else {
		// Modo decodificação - lê do primeiro argumento não-flag
		if flag.NArg() == 0 {
			log.Fatal("No input file specified for decoding")
		}
		inputFile := flag.Arg(0)

		file, err := os.Open(inputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		mr := uuencode.NewMultiReader(file)

		for {
			fileInfo, reader, err := mr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			outFile, err := os.Create(fileInfo.Name)
			if err != nil {
				log.Fatal(err)
			}

			if _, err := io.Copy(outFile, reader); err != nil {
				outFile.Close()
				log.Fatal(err)
			}

			if err := outFile.Chmod(fileInfo.Mode); err != nil {
				log.Printf("Warning: could not set mode for %s: %v", fileInfo.Name, err)
			}
			outFile.Close()
			fmt.Printf("file: %s mode: %03o\n", fileInfo.Name, fileInfo.Mode)
		}
	}
}
