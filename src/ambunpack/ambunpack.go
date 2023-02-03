package main

import (
	"ambutils/amb"
	"flag"
	"fmt"
	"io/ioutil"
)

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

	f := amb.New(fileNameIn)
	e := f.LoadFile()
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	entryNames := f.ListNames()
	for _, name := range entryNames {
		fmt.Printf("Processing file %12s: ", name)
		dstFile := dirNameOut + name
		entry, _ := f.GetEntryRaw(name)
		err := ioutil.WriteFile(dstFile, entry, 0644)
		if err != nil {
			fmt.Print("Can't write file!\n")
			return
		}
		fmt.Print("[DONE]\n")
	}
	fmt.Println("Unpacked successfuly!")
}
