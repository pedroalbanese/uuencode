package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pedroalbanese/uuencode"
)

var (
	dec   = flag.Bool("d", false, "Decode instead of Encode")
	ifile = flag.String("f", "", "Target file")
	ofile = flag.String("o", "", "Output file")
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
		infile, err := os.Open(*ifile)
		if err != nil {
			log.Println(err)
		}
		var outfile *os.File
		if *ofile == "-" || *ofile == "" {
			outfile = os.Stdout
		} else {
			outfile, err = os.Create(*ofile)
			if err != nil {
				log.Println(err)
			}
		}
		uw := uuencode.NewWriter(outfile, *ifile, 0664)
		if _, err = io.Copy(uw, infile); err != nil {
			return
		}
		if err := uw.Flush(); err != nil {
			return
		}
	} else {
		file, err := os.Open(*ifile)
		if err != nil {
			log.Println(err)
		}
		var outfile *os.File
		if *ofile == "-" || *ofile == "" {
			outfile = os.Stdout
		} else {
			outfile, err = os.Create(*ofile)
			if err != nil {
				log.Println(err)
			}
		}
		ur := uuencode.NewReader(file, nil)
		_, err = io.Copy(outfile, ur)
		if err != nil {
			log.Fatal(err)
		}
		f, _ := ur.File()
		m, _ := ur.Mode()
		fmt.Fprintf(os.Stderr, "file: %s, mode: %03o\n", f, m)
	}
}