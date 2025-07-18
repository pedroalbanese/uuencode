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
	ifile = flag.String("f", "", "Target file (use '-' or empty for stdin)")
)

func main() {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage of", os.Args[0]+":")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !*dec {
		// Modo de codificação (UUencode)
		infile, err := os.Open(*ifile)
		if err != nil {
			log.Fatal(err)
		}
		defer infile.Close()

		info, err := infile.Stat()
		if err != nil {
			log.Fatal(err)
		}

		// Usa o nome do arquivo de entrada no cabeçalho UUencode
		uuw := uuencode.NewWriter(os.Stdout, *ifile, info.Mode())
		if _, err := io.Copy(uuw, infile); err != nil {
			log.Fatal(err)
		}
		if err := uuw.Flush(); err != nil {
			log.Fatal(err)
		}
	} else {
		// Modo de decodificação (UUdecode)
		var infile *os.File
		var err error
		
		if *ifile == "-" || *ifile == "" {
			infile = os.Stdin
		} else {
			infile, err = os.Open(*ifile)
			if err != nil {
				log.Fatal(err)
			}
			defer infile.Close()
		}

		// Cria o decoder UUencode
		uur := uuencode.NewReader(infile, nil)

		// Obtém o nome do arquivo do cabeçalho UUencode
		filename, err := uur.File()
		if err != nil {
			log.Fatal("Error getting filename from UUencoded data:", err)
		}

		// Cria o arquivo de saída
		outfile, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer outfile.Close()

		// Decodifica e escreve no arquivo de saída
		if _, err := io.Copy(outfile, uur); err != nil {
			log.Fatal(err)
		}

		// Obtém e define as permissões do arquivo
		mode, err := uur.Mode()
		if err != nil {
			log.Println("Warning: could not get file mode from UUencoded data:", err)
		} else {
			if err := outfile.Chmod(mode); err != nil {
				log.Println("Warning: could not set file mode:", err)
			}
		}

		// Exibe informações no stderr
		fmt.Fprintf(os.Stderr, "file: %s, mode: %03o\n", filename, mode)
	}
}
