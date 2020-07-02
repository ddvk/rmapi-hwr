package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ddvk/rmapi-hwr/models"
	"github.com/juruen/rmapi/archive"
	"github.com/juruen/rmapi/encoding/rm"
)

const url = "https://cloud.myscript.com/api/v4.0/iink/batch"

type config struct {
	page           int
	applicationKey string
	hmacKey        string
	lang           string
	inputType      string
	outputType     string
}

func sendRequest(key, hmackey string, data []byte, mimeType string) (body []byte, err error) {
	fullkey := key + hmackey
	mac := hmac.New(sha512.New, []byte(fullkey))
	mac.Write(data)
	result := hex.EncodeToString(mac.Sum(nil))

	client := http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Set("Accept", mimeType+", application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("applicationKey", key)
	req.Header.Set("hmac", result)

	res, err := client.Do(req)

	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Not ok, Status: %d", res.StatusCode)
		return
	}

	return body, nil
}

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

var noContent = errors.New("no page content")

func getJson(zip *archive.Zip, contenttype string, lang string, pageNumber int) (r []byte, err error) {
	numPages := len(zip.Pages)

	if pageNumber >= numPages || pageNumber < 0 {
		err = fmt.Errorf("page %d outside range, max: %d", numPages)
		return
	}

	batch := models.BatchInput{
		Configuration: &models.Configuration{
			Lang: lang,
		},
		StrokeGroups: []*models.StrokeGroup{
			&models.StrokeGroup{},
		},
		ContentType: &contenttype,
		Width:       14040,
		Height:      18720,
		XDPI:        2280,
		YDPI:        2280,
	}

	sg := batch.StrokeGroups[0]

	page := zip.Pages[pageNumber]

	if page.Data == nil {
		return nil, noContent
	}

	for _, layer := range page.Data.Layers {
		for _, line := range layer.Lines {
			pointerType := ""
			if line.BrushType == rm.EraseArea {
				continue
			}
			if line.BrushType == rm.Eraser {
				pointerType = "ERASER"
			}
			stroke := models.Stroke{
				X:           make([]float32, 0),
				Y:           make([]float32, 0),
				PointerType: pointerType,
			}
			sg.Strokes = append(sg.Strokes, &stroke)

			for _, point := range line.Points {
				x := point.X * 10
				y := point.Y * 10
				stroke.X = append(stroke.X, x)
				stroke.Y = append(stroke.Y, y)
			}
		}
	}

	r, err = batch.MarshalBinary()
	if err != nil {
		return
	}
	return
}

func main() {
	applicationKey := os.Getenv("RMAPI_HWR_APPLICATIONKEY")
	if applicationKey == "" {
		log.Fatal("provide the myScript applicationKey in: RMAPI_HWR_APPLICATIONKEY")
	}
	hmacKey := os.Getenv("RMAPI_HWR_HMAC")
	if applicationKey == "" {
		log.Fatal("provide the myScript hmac in: RMAPI_HWR_HMAC")
	}

	flag.Usage = func() {
		exec := os.Args[0]
		output := flag.CommandLine.Output()
		fmt.Fprintf(output, "Usage: %s [options] somefile.zip\n", exec)
		fmt.Fprintln(output, "\twhere somefile.zip is what you got with rmapi get")
		fmt.Fprintln(output, "\tOutputs: Text->text, Math->LaTex, Diagram->svg")
		fmt.Fprintln(output, "Options:")
		flag.PrintDefaults()
	}
	var textType = flag.String("type", "Text", "type of the content: Text, Math, Diagram")
	var lang = flag.String("lang", "en_US", "language culture")
	//todo: page range, all pages etc
	var page = flag.Int("page", 0, "page to convert (default lastopened)")
	// var outputFile = flag.String("o", "-", "output default stdout, wip")
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

	pageNumber := *page

	if pageNumber == 0 {
		pageNumber = zip.Content.LastOpenedPage
	} else if pageNumber < 0 {
		pageNumber = 0
	} else {
		pageNumber -= 1
	}

	contenttype, output := setContentType(*textType)

	//loop over all pages, some scatter-gather
	js, err := getJson(zip, contenttype, *lang, pageNumber)
	if err != nil {
		log.Fatal(err)
	}

	body, err := sendRequest(applicationKey, hmacKey, js, output)
	if err != nil {
		if body != nil {
			log.Println(string(body))
		}
		log.Fatal(err)
	}

	//todo: file output
	fmt.Println(string(body))
}

func setContentType(requested string) (contenttype string, output string) {
	switch strings.ToLower(requested) {
	case "math":
		contenttype = "Math"
		output = "application/x-latex"
	case "text":
		contenttype = "Text"
		output = "text/plain"
	case "diagram":
		contenttype = "Diagram"
		output = "image/svg+xml"
	default:
		log.Fatal("unsupported content type: " + contenttype)
	}
	return
}
