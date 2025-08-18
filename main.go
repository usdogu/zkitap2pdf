/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import "zkitap2pdf/cmd"

func main() {
	cmd.Execute()
}

// import (
// 	"bytes"
// 	"fmt"
// 	"image/jpeg"
// 	"io"
// 	"log"
// 	"os"
// 	"slices"
// 	"strconv"
// 	"strings"

// 	"github.com/grafana/gofpdf"
// 	"github.com/h2non/filetype"
// 	"golang.org/x/image/webp"
// )
// func main() {
// 	// List of .webp images in order
// 	webpFiles := []string{}
// 	files, err := os.ReadDir("deneme/decrypted")
// 	if err != nil {
// 		log.Fatalf("Error reading directory: %v", err)
// 	}
// 	for _, file := range files {
// 		webpFiles = append(webpFiles, "deneme/decrypted/"+file.Name())
// 	}
// 	slices.SortFunc(webpFiles, func(a, b string) int {
// 		aSplit := strings.Split(a, "-")
// 		aInt, _ := strconv.Atoi(strings.TrimSuffix(aSplit[len(aSplit)-1], ".webp")) // .webp
// 		bSplit := strings.Split(b, "-")
// 		bInt, _ := strconv.Atoi(strings.TrimSuffix(bSplit[len(bSplit)-1], ".webp"))
// 		fmt.Println(aSplit[len(aSplit)-1], bSplit[len(bSplit)-1])
// 		return aInt - bInt
// 	})
// 	pdf := gofpdf.New("P", "mm", "A4", "")
// 	pageWidth, pageHeight := pdf.GetPageSize()

// 	for _, file := range webpFiles {
// 		fmt.Printf("%s ", file)
// 		// Open WebP file
// 		f, err := os.Open(file)
// 		if err != nil {
// 			log.Fatalf("Error opening %s: %v", file, err)
// 		}
// 		var imageType string
// 		kind, _ := filetype.MatchReader(f)
// 		f.Seek(0, io.SeekStart)
// 		var imageContent []byte
// 		switch kind.MIME.Value {
// 		case "image/webp":
// 			imageType = "webp"
// 		case "image/jpeg":
// 			imageType = "jpg"
// 		case "image/png":
// 			imageType = "png"
// 		default:
// 			log.Fatalf("Error: %s is not a webp or jpg or png file", file)
// 		}
// 		if imageType == "webp" {
// 			img, err := webp.Decode(f)
// 			if err != nil {
// 				log.Fatalf("Error decoding %s: %v", file, err)
// 			}
// 			var buf bytes.Buffer
// 			if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
// 				log.Fatalf("Error encoding JPEG: %v", err)
// 			}
// 			imageContent = buf.Bytes()
// 			imageType = "jpg"
// 		} else {
// 			imageContent, err = io.ReadAll(f)
// 			if err != nil {
// 				log.Fatalf("Error reading %s: %v", file, err)
// 			}
// 		}
// 		f.Close()

// 		// Give each image a unique name in the PDF registry
// 		imageName := file // could also use fmt.Sprintf("img%d", idx)

// 		// Register image from bytes buffer
// 		pdf.RegisterImageOptionsReader(
// 			imageName,
// 			gofpdf.ImageOptions{ImageType: imageType, ReadDpi: true},
// 			bytes.NewReader(imageContent),
// 		)

// 		// Add page and draw image full-page
// 		pdf.AddPage()
// 		pdf.ImageOptions(imageName, 0, 0, pageWidth, pageHeight, false,
// 			gofpdf.ImageOptions{ImageType: imageType, ReadDpi: true}, 0, "")
// 	}

// 	// Output the PDF file
// 	if err := pdf.OutputFileAndClose("output.pdf"); err != nil {
// 		log.Fatalf("Error saving PDF: %v", err)
// 	}

// 	log.Println("PDF created: output.pdf")
// }
