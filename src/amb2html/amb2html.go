package main

import (
	"ambutils/amb"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var f *amb.AmbFile

func entryToHtml(entry []byte) []byte {
	var rv []rune
	var src = []rune(string(entry))

	var readlink bool = false
	var escnow bool = false
	var opentag string = ""
	for _, c := range src {
		if c == 0x000D { // \r
			continue
		}

		if readlink {
			if c == 0x000A || c == ':' {

				linkIndex := strings.LastIndex(string(rv), "\"") + 1
				linkValue := string(rv)[linkIndex:]
				if f.HasEntry(linkValue) {
					rv = []rune(strings.TrimSuffix(string(rv), ".ama"))
					rv = append(rv, []rune(".htm")...)
				}

				rv = append(rv, []rune("\">")...)
				opentag = "a"
				readlink = false
			} else {
				rv = append(rv, c) // May be need url encode?
			}
			continue
		}

		if escnow {
			if c == '%' {
				rv = append(rv, '%')
			} else if c == 'l' {
				rv = append(rv, []rune("<a href=\"")...)
				readlink = true
			} else if c == 'h' {
				rv = append(rv, []rune("<h1>")...)
				opentag = "h1"
			} else if c == '!' {
				rv = append(rv, []rune("<div style=\"color:#333300\">")...)
				opentag = "div"

			} else if c == 'b' {
				rv = append(rv, []rune("<div style=\"color:#888888\">")...)
				opentag = "div"
			}

			escnow = false
			continue
		}

		if opentag != "" && (c == '%' || c == 0x000A) {
			// == chop start
			var pos = len(rv) - 1
			for j := len(rv) - 1; j > 0; j-- {
				if rv[j] != 0x0020 {
					pos = j
					break
				}
			}
			var rv_tmp []rune
			for j, ch := range []rune(rv) {
				rv_tmp = append(rv_tmp, ch)
				if j == pos {
					rv_tmp = append(rv_tmp, []rune("</"+opentag+">")...)
				}
			}
			rv = rv_tmp
			// == chop end

			opentag = ""
		}

		if c == '%' {
			escnow = true
		} else if c == 0x000A {
			rv = append(rv, []rune("<br />\n")...)
		} else {
			rv = append(rv, c)
		}

	}

	if opentag != "" {
		rv = append(rv, []rune("</"+opentag+">")...)
	}

	return []byte(string(rv))
}

func main() {
	var fileNameIn, dirNameOut string
	flag.StringVar(&fileNameIn, "i", "", "AMB-file to read")
	flag.StringVar(&dirNameOut, "o", "./", "Destination directory")
	flag.Parse()

	if fileNameIn == "" || dirNameOut == "" {
		fmt.Println("Usage:")
		fmt.Println("ambunpack -i inputfile [-o dirname]")
		return
	}
	fmt.Printf("Unpacking file \"%s\" into directory \"%s\"...\n", fileNameIn, dirNameOut)
	if dirNameOut[len(dirNameOut)-1] != '/' {
		dirNameOut += "/"
	}

	f = amb.New(fileNameIn)
	e := f.LoadFile()
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	entryNames := f.ListNames()
	for _, name := range entryNames {
		// Process only AMA-files
		if !strings.HasSuffix(name, ".ama") {
			continue
		}
		fmt.Printf("Processing file %12s: ", name)
		dstFile := dirNameOut + strings.TrimSuffix(name, ".ama") + ".htm"
		entry, _ := f.GetEntry(name)
		var htmlEntry []rune
		htmlEntry = append(htmlEntry, []rune("<html>\n<head>\n")...)
		htmlEntry = append(htmlEntry, []rune("<style>\n")...)
		htmlEntry = append(htmlEntry, []rune("body{color:#000000; background-color:#E0E0E0; font-family:monospace;white-space:pre;}\n")...)
		htmlEntry = append(htmlEntry, []rune("a{color: #0000AA;}\n")...)
		htmlEntry = append(htmlEntry, []rune("\n")...)
		htmlEntry = append(htmlEntry, []rune("\n")...)
		htmlEntry = append(htmlEntry, []rune("</style>\n")...)
		htmlEntry = append(htmlEntry, []rune("</head>\n<body>\n")...)
		htmlEntry = append(htmlEntry, []rune(string(entryToHtml(entry)))...)
		htmlEntry = append(htmlEntry, []rune("</body>\n</html>\n")...)
		err := ioutil.WriteFile(dstFile, []byte(string(htmlEntry)), 0644)
		if err != nil {
			fmt.Print("Can't write file!\n")
			return
		}
		fmt.Print("[DONE]\n")
	}
	fmt.Println("Unpacked successfuly!")
}
