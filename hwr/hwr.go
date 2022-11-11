package hwr

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sync/semaphore"

	"github.com/ddvk/rmapi-hwr/hwr/client"
	"github.com/ddvk/rmapi-hwr/hwr/models"
	"github.com/juruen/rmapi/archive"
	"github.com/juruen/rmapi/encoding/rm"
)

var NoContent = errors.New("no page content")

type Config struct {
	Page           int
	applicationKey string
	hmacKey        string
	Lang           string
	InputType      string
	OutputType     string
	OutputFile     string
	AddPages       bool
	BatchSize      int64
}

func getJson(zip *archive.Zip, contenttype string, lang string, pageNumber int) (r []byte, err error) {
	numPages := len(zip.Pages)

	if pageNumber >= numPages || pageNumber < 0 {
		err = fmt.Errorf("page %d outside range, max: %d", pageNumber, numPages)
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
		return nil, NoContent
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

func Hwr(zip *archive.Zip, cfg Config) {
	applicationKey := os.Getenv("RMAPI_HWR_APPLICATIONKEY")
	if applicationKey == "" {
		log.Fatal("provide the myScript applicationKey in: RMAPI_HWR_APPLICATIONKEY")
	}
	hmacKey := os.Getenv("RMAPI_HWR_HMAC")
	if applicationKey == "" {
		log.Fatal("provide the myScript hmac in: RMAPI_HWR_HMAC")
	}

	capacity := 1
	start := 0
	var end int

	if cfg.Page == 0 {
		start = zip.Content.LastOpenedPage
		end = start
	} else if cfg.Page < 0 {
		capacity = len(zip.Pages)
		end = capacity - 1
	} else {
		start = cfg.Page - 1
		end = start
	}
	result := make([][]byte, capacity)

	contenttype, output := setContentType(cfg.InputType)

	ctx := context.TODO()
	sem := semaphore.NewWeighted(cfg.BatchSize)
	for p := start; p <= end; p++ {
		log.Println("Page: ", p)
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}
		go func(p int) {
			defer sem.Release(1)
			js, err := getJson(zip, contenttype, cfg.Lang, p)
			if err != nil {
				log.Fatalf("Can't get page: %d %v\n", p, err)
			}
			log.Println("sending request: ", p)

			body, err := client.SendRequest(applicationKey, hmacKey, js, output)
			if err != nil {
				if body != nil {
					log.Println(string(body))
				}
				log.Fatal(err)
			}
			result[p] = body
			log.Println("converted page ", p)
		}(p)
	}
	log.Println("wating for all to finish")
	if err := sem.Acquire(ctx, cfg.BatchSize); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
	}

	if cfg.OutputFile == "-" {
		dump(result, cfg.AddPages)
	} else {
		//text file
		f, err := os.Create(cfg.OutputFile + ".txt")
		if err != nil {
			dump(result, cfg.AddPages)
			log.Fatal(err)
		}

		for _, c := range result {
			f.Write(c)
			f.Write([]byte("\n"))
		}
		f.Close()
	}
}

func dump(result [][]byte, addPages bool) {
	for p, c := range result {
		if addPages {
			fmt.Printf("=== Page %d ===\n", p)

		}
		fmt.Println(string(c))
	}
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
	case "jiix":
		contenttype = "Text"
		output = "application/vnd.myscript.jiix"
	default:
		log.Fatal("unsupported content type: " + contenttype)
	}
	return
}
