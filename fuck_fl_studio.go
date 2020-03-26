package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Welcome you beautiful soul.\n" +
		"Let's \"convert\" some FL studio \"wav\" samples to .ogg\n")

	// We want parameters
	if len(os.Args) < 2 {
		fmt.Println("Supply .wav filenames as parameters, or drag and drop files on the executable")
		os.Exit(2)
	}

	// loop over all the params
	for _, name := range os.Args[1:] {
		// wrap each item in a func so we can defer closing
		if filepath.Ext(name) != ".wav" {
			fmt.Println("not .wav file extension, exiting.")
			os.Exit(3)
		}
		fin, err := os.Open(name)
		check(err)

		// read file and close
		content, err := ioutil.ReadAll(fin)
		check(err)


		// first 4 bytes should be RIFF
		if !bytes.Equal(content[:4], []byte{82, 73, 70, 70}) ||
			// 0x36 to 0x3A should be "OggS"
			!bytes.Equal(content[54:58], []byte{79, 103, 103, 83}) {
			fmt.Println("not a OGG within a WAV..")
			os.Exit(4)
		}

		// create new file
		newfile := strings.TrimSuffix(name, "wav") + "ogg"
		fout, err := os.Create(newfile)
		check(err)

		// delete the first 54 bytes
		_, err = fin.Seek(54, io.SeekStart)
		check(err)
		_, err = io.Copy(fout, fin)
		check(err)
		err = fin.Close()
		check(err)
		err = fout.Close()
		check(err)
	}

}
