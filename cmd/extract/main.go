package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/juruen/rmapi/annotations"
	"github.com/juruen/rmapi/archive"
)

func doJson(input, output, token string) error {
	jsn, err := jsonannotations(input)
	if err != nil {
		return err
	}
	if token != "" {
		SendRequest(jsn, token)
	} else {
		if output == "" {
			nameOnly := strings.TrimSuffix(input, filepath.Ext(input))
			output = nameOnly + ".json"
		}
		return ioutil.WriteFile(output, jsn, 0700)
	}
	return nil
}

func main() {
	inputName := flag.String("i", "", "file to convert")
	outputName := flag.String("o", "", "outpufilename")
	extract := flag.String("e", "", "extract, a - annotations, p - pdf (default)")
	format := flag.String("f", "json", "format (json, txt)")
	token := flag.String("t", "", "readwise token")
	flag.Parse()
	var err error

	switch *extract {

	case "a":
		if *format == "json" {
			err = doJson(*inputName, *outputName, *token)
		} else {
			err = txtannotations(*inputName, *outputName)
		}
	case "", "p":
		err = convert(*inputName, *outputName)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func jsonannotations(inputName string) (output []byte, err error) {
	file, err := os.Open(inputName)
	if err != nil {
		return
	}
	defer file.Close()

	zip := archive.NewZip()

	fi, err := file.Stat()
	if err != nil {
		return
	}

	err = zip.Read(file, fi.Size())
	if err != nil {
		return
	}

	rw := ReadWise{}

	docId := zip.UUID
	nameOnly := strings.TrimSuffix(inputName, filepath.Ext(inputName))

	for index, p := range zip.Pages {
		if p.Highlights == nil {
			continue
		}
		high := p.Highlights.Highlights[0]
		sort.Slice(high, func(i, j int) bool {
			y1 := high[i].Rects[0].Y
			y2 := high[j].Rects[0].Y
			x1 := high[i].Rects[0].X
			x2 := high[j].Rects[0].X
			if math.Abs(y1-y2) < 5 {
				return x1 < x2
			}
			return y1 < y2
		})
		for _, h := range high {
			rw.Highlights = append(rw.Highlights,
				Highlight{
					DocumentID:   docId,
					Text:         h.Text,
					Title:        nameOnly,
					LocationType: "page",
					Location:     index + 1,
				})
		}
	}

	return json.Marshal(rw)
}
func txtannotations(inputName, outputName string) error {
	if outputName == "" {
		nameOnly := strings.TrimSuffix(inputName, filepath.Ext(inputName))
		outputName = nameOnly + ".txt"
	}
	file, err := os.Open(inputName)
	if err != nil {
		return err
	}
	defer file.Close()

	zip := archive.NewZip()

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	err = zip.Read(file, fi.Size())
	if err != nil {
		return err
	}
	f, err := os.Open(outputName)
	if err != nil {
		return err
	}
	defer f.Close()

	for index, p := range zip.Pages {
		if p.Highlights == nil {
			continue
		}
		f.WriteString(fmt.Sprintf("Page %d\n", index))
		high := p.Highlights.Highlights[0]
		sort.Slice(high, func(i, j int) bool {
			y1 := high[i].Rects[0].Y
			y2 := high[j].Rects[0].Y
			x1 := high[i].Rects[0].X
			x2 := high[j].Rects[0].X
			if math.Abs(y1-y2) < 5 {
				return x1 < x2
			}
			return y1 < y2
		})
		for _, h := range high {
			f.WriteString(fmt.Sprintf(" X:%d Y:%d\t %s\n", int(h.Rects[0].X), int(h.Rects[0].Y), h.Text))
		}
	}

	return nil
}

func convert(inputName, outputName string) (err error) {
	if inputName == "" {
		return errors.New("missing input file")
	}

	if outputName == "" {
		nameOnly := strings.TrimSuffix(inputName, filepath.Ext(inputName))
		outputName = nameOnly + ".pdf"
	}

	outputFile, err := os.Create(outputName)
	if err != nil {
		return fmt.Errorf("can't create outputfile %w", err)
	}
	defer outputFile.Close()

	reader, err := zip.OpenReader(inputName)
	if err != nil {
		return fmt.Errorf("can't open file %w", err)
	}
	defer reader.Close()

	options := annotations.PdfGeneratorOptions{
		AllPages: true,
	}
	gen := annotations.CreatePdfGenerator(inputName, outputName, options)
	return gen.Generate()
}
