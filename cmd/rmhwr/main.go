package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/ddvk/rmapi-hwr/hwr"
	"github.com/juruen/rmapi/archive"
	"github.com/juruen/rmapi/encoding/rm"
)

func loadRmPage(filename string) (zip *archive.Zip, err error) {
	zip = archive.NewZip()
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pageData, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal("cant read fil")
		return
	}
	page := archive.Page{}
	page.Data = rm.New()
	err = page.Data.UnmarshalBinary(pageData)
	if err != nil {
		return nil, err
	}

	zip.Pages = append(zip.Pages, page)

	return zip, nil

}

func loadRmZip(filename string) (zip *archive.Zip, err error) {
	zip = archive.NewZip()
	file, err := os.Open(filename)

	if err != nil {
		return
	}

	defer file.Close()

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
	var addPages = flag.Bool("a", false, "add page headers")
	cfg := hwr.Config{
		Page:      *page,
		Lang:      *lang,
		InputType: *inputType,
		AddPages:  *addPages,
		BatchSize: *flag.Int64("b", 3, "batch size"),
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("no file specified")
	}

	filename := args[0]
	ext := path.Ext(filename)
	cfg.OutputFile = strings.TrimSuffix(filename, ext)

	var err error
	var z *archive.Zip

	switch ext {
	case ".zip":
		z, err = loadRmZip(filename)
	case ".rm":
		z, err = loadRmPage(filename)
	default:
		log.Fatal("Unsupported file")

	}
	if err != nil {
		log.Fatalln(err, "Can't read file ", filename)

	}
	hwr.Hwr(z, cfg)
}
