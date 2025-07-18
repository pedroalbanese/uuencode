package uuencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	uuBaseChar  = ' ' // Primeiro caractere da codificação UU
	uuMaxLine   = 45  // Máximo de bytes por linha (codificados em 60 caracteres)
	uuLineLen   = 61  // Comprimento da linha codificada (incluindo o caractere de comprimento)
	uuBlockSize = 3   // Tamanho do bloco de decodificação
)

// FileInfo contém metadados do arquivo do cabeçalho
type FileInfo struct {
	Name string
	Mode os.FileMode
}

// NewWriter cria um novo codificador UUencode
func NewWriter(w io.Writer, filename string, mode os.FileMode) *Writer {
	ww := &Writer{
		w:        w,
		filename: filename,
		mode:     mode,
		buf:      make([]byte, uuMaxLine),
	}
	fmt.Fprintf(w, "begin %o %s\n", mode.Perm(), filename)
	return ww
}

// Writer implementa a codificação UUencode
type Writer struct {
	w        io.Writer
	filename string
	mode     os.FileMode
	buf      []byte
	n        int
}

func (w *Writer) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		m := copy(w.buf[w.n:], p)
		w.n += m
		p = p[m:]
		n += m

		if w.n == uuMaxLine {
			if err := w.encodeLine(); err != nil {
				return n, err
			}
			w.n = 0
		}
	}
	return n, nil
}

func (w *Writer) encodeLine() error {
	if w.n == 0 {
		return nil
	}

	// O primeiro byte indica o comprimento da linha (bytes decodificados)
	lengthChar := byte(uuBaseChar + w.n)
	if _, err := w.w.Write([]byte{lengthChar}); err != nil {
		return err
	}

	encoded := make([]byte, 0, uuLineLen)
	for i := 0; i < w.n; i += uuBlockSize {
		// Pegar até 3 bytes
		var b0, b1, b2 byte
		b0 = w.buf[i]
		if i+1 < w.n {
			b1 = w.buf[i+1]
		}
		if i+2 < w.n {
			b2 = w.buf[i+2]
		}

		// Codificar em 4 caracteres
		encoded = append(encoded,
			uuBaseChar+((b0>>2)&0x3F),
			uuBaseChar+(((b0&0x03)<<4)|(b1>>4)),
			uuBaseChar+(((b1&0x0F)<<2)|(b2>>6)),
			uuBaseChar+(b2&0x3F),
		)
	}

	if _, err := w.w.Write(encoded); err != nil {
		return err
	}
	_, err := w.w.Write([]byte{'\n'})
	return err
}

func (w *Writer) Flush() error {
	if w.n > 0 {
		if err := w.encodeLine(); err != nil {
			return err
		}
	}
	_, err := w.w.Write([]byte("`\nend\n"))
	return err
}

// NewReader cria um novo decodificador UUencode
func NewReader(r io.Reader, fileInfo *FileInfo) *Reader {
	return &Reader{
		r:        bufio.NewReader(r),
		fileInfo: fileInfo,
	}
}

// Reader implementa a decodificação UUencode
type Reader struct {
	r          *bufio.Reader
	fileInfo   *FileInfo
	headerRead bool
	eof        bool
	buf        []byte
	pos        int
}

func (r *Reader) readHeader() error {
	if r.headerRead {
		return nil
	}

	for {
		line, err := r.r.ReadString('\n')
		if err != nil {
			return err
		}

		line = strings.TrimRight(line, "\r\n")
		if strings.HasPrefix(line, "begin ") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				if r.fileInfo == nil {
					r.fileInfo = &FileInfo{}
				}
				r.fileInfo.Name = fields[2]
				mode, err := strconv.ParseUint(fields[1], 8, 32)
				if err != nil {
					return fmt.Errorf("invalid mode: %v", err)
				}
				r.fileInfo.Mode = os.FileMode(mode)
				r.headerRead = true
				return nil
			}
		}
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if !r.headerRead {
		if err := r.readHeader(); err != nil {
			return 0, err
		}
	}

	if r.eof {
		return 0, io.EOF
	}

	if r.pos >= len(r.buf) {
		line, err := r.r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				r.eof = true
				return 0, io.EOF
			}
			return 0, err
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "`" || line == "end" {
			r.eof = true
			return 0, io.EOF
		}

		if len(line) == 0 {
			return 0, nil
		}

		// O primeiro caractere indica o comprimento
		length := int(line[0] - uuBaseChar)
		if length <= 0 || length > uuMaxLine {
			return 0, errors.New("invalid length character")
		}

		data := line[1:]
		r.buf = make([]byte, length)
		bytesDecoded := 0
		requiredChars := ((length + 2) / 3) * 4

		// Verificação rigorosa do comprimento dos dados
		if len(data) < requiredChars {
			return 0, fmt.Errorf("incomplete data: expected %d chars, got %d", requiredChars, len(data))
		}

		for i := 0; i < requiredChars && bytesDecoded < length; i += 4 {
			c0 := safeUUIndex(data, i)
			c1 := safeUUIndex(data, i+1)
			c2 := safeUUIndex(data, i+2)
			c3 := safeUUIndex(data, i+3)

			if c0 == -1 || c1 == -1 {
				return 0, fmt.Errorf("invalid encoding at position %d-%d", i, i+1)
			}

			// Primeiro byte
			r.buf[bytesDecoded] = byte((c0 << 2) | (c1 >> 4))
			bytesDecoded++

			if bytesDecoded >= length {
				break
			}

			if c2 == -1 {
				continue
			}

			// Segundo byte
			r.buf[bytesDecoded] = byte((c1 << 4) | (c2 >> 2))
			bytesDecoded++

			if bytesDecoded >= length {
				break
			}

			if c3 == -1 {
				continue
			}

			// Terceiro byte
			r.buf[bytesDecoded] = byte((c2 << 6) | c3)
			bytesDecoded++
		}

		r.pos = 0
	}

	n = copy(p, r.buf[r.pos:])
	r.pos += n
	return n, nil
}

// safeUUIndex retorna o valor decodificado do caractere ou -1 se for inválido
func safeUUIndex(s string, pos int) int {
	if pos >= len(s) {
		return -1
	}
	c := s[pos] - uuBaseChar
	if c > 64 {
		return -1
	}
	return int(c & 0x3F)
}

func (r *Reader) File() (string, error) {
	if !r.headerRead {
		if err := r.readHeader(); err != nil {
			return "", err
		}
	}
	return r.fileInfo.Name, nil
}

func (r *Reader) Mode() (os.FileMode, error) {
	if !r.headerRead {
		if err := r.readHeader(); err != nil {
			return 0, err
		}
	}
	return r.fileInfo.Mode, nil
}
