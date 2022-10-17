# UUE
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/pedroalbanese/uuencode/blob/master/LICENSE.md) 
[![GoDoc](https://godoc.org/github.com/pedroalbanese/uuencode?status.png)](http://godoc.org/github.com/pedroalbanese/uuencode)
[![GitHub downloads](https://img.shields.io/github/downloads/pedroalbanese/uuencode/total.svg?logo=github&logoColor=white)](https://github.com/pedroalbanese/uuencode/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/pedroalbanese/uuencode)](https://goreportcard.com/report/github.com/pedroalbanese/uuencode)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/pedroalbanese/uuencode)](https://golang.org)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/pedroalbanese/uuencode)](https://github.com/pedroalbanese/uuencode/releases)  
UUEncode is a tool that converts to and from uuencoding
<pre>     permission mode _______       ______ file name to be given to decoded file
                            |     |
     begin line ____ begin 644 filename
                     M;2XN+BXN+R\N+B\O+BXN+BXN+R\N+B\O+BXO+RXO+RXN+B\ON+B\O+BXN
     encoded data __ M"AM;-#LV2"`@("`@+R`@7`H;6S$[,3%("AM;,CLQ,4@@("`@<("\*&ULS
                     `
     end line ______ end
</pre>
## Usage
<pre>
Usage of uuencode:
  -d    Decode instead of Encode
  -f string
        Target file
  -o string
        Output file</pre>
 
## License
This project is licensed under the MIT License. 
 
**Copyright (c) 2020 Blake Williams <code@shabbyrobe.org>**  
**Copyright (c) 2022 Pedro Albanese <pedroalbanese@hotmail.com>**
