package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/ddvk/rmapi-hwr/hwr"
	"github.com/juruen/rmapi/archive"
)

func loadRmZip(filename string) (zip *archive.Zip, err error) {
	zip = archive.NewZip()
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return
	}
	fi, err := file.Stat()
	if err != nil {
		return
	}
	err = zip.Read(file, fi.Size())
	if err != nil {
		return
	}
	numPages := len(zip.Pages)

	if numPages == 0 {
		err = errors.New("no pages")
		return
	}
	return zip, nil
}

func main() {

	flag.Usage = func() {
		exec := os.Args[0]
		output := flag.CommandLine.Output()
		fmt.Fprintf(output, "Usage: %s [options] somefile.zip\n", exec)
		fmt.Fprintln(output, "\twhere somefile.zip is what you got with rmapi get")
		fmt.Fprintln(output, "\tOutputs: Text->text, Math->LaTex, Diagram->svg")
		fmt.Fprintln(output, "Options:")
		flag.PrintDefaults()
	}
	var inputType = flag.String("type", "Text", "type of the content: Text, Math, Diagram")
	var lang = flag.String("lang", "en_US", "language culture")
	//todo: page range, all pages etc
	var page = flag.Int("page", -1, "page to convert (default all)")
	//var outputFile = flag.String("o", "-", "output default stdout, wip")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("no file specified")
	}

	filename := args[0]
	zip, err := loadRmZip(filename)
	if err != nil {
		log.Fatal(err)
	}

	output := strings.TrimSuffix(filename, path.Ext(filename))

	cfg := hwr.Config{
		Page:       *page,
		Lang:       *lang,
		InputType:  *inputType,
		OutputFile: output,
	}

	hwr.Hwr(zip, cfg)

}
