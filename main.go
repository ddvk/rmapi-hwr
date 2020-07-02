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

func sendApi(key, hmackey string, data []byte, mimeType string) ([]byte, error) {
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
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		strBody := string(body)
		log.Print(strBody)
		return nil, errors.New("not OK")
	}

	return body, nil
}

func getJson(filename, contenttype string, lang string, pageNumber int) (r []byte, err error) {
	zip := archive.NewZip()
	file, err := os.Open(filename)
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
	if pageNumber > numPages {
		err = errors.New(fmt.Sprintf("max pages %d", numPages))
		return
	}
	if pageNumber == 0 {
		pageNumber = zip.Content.LastOpenedPage
	} else if pageNumber < 0 {
		pageNumber = 0
	} else {
		pageNumber -= 1
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

	filename := ""
	var textType = flag.String("type", "Text", "type of the content: Text, Math, Diagram")
	var lang = flag.String("lang", "en_US", "language culture")
	var page = flag.Int("page", 0, "page to convert")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("no file specified")
	}
	filename = args[0]

	var contenttype = ""
	var output = ""

	switch strings.ToLower(*textType) {
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

	js, err := getJson(filename, contenttype, *lang, *page)
	if err != nil {
		log.Fatal(err)
	}

	body, err := sendApi(applicationKey, hmacKey, js, output)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}
