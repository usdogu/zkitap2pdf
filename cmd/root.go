package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"zkitap2pdf/types"
	"zkitap2pdf/util"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/grafana/gofpdf"
	"github.com/h2non/filetype"
	"github.com/spf13/cobra"
	"golang.org/x/image/webp"
)

var publisherData types.PublisherData
var publisher types.Publisher
var shelfId int = -1
var book types.Book

func downloadBook() ([]byte, error) {
	splitted := strings.Split(book.DataFile, "/")
	zipUrl := strings.Join(splitted[:len(splitted)-1], "/") + "/" + book.DataZip
	resp, err := http.Get(zipUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func decryptBook(keyPrefix string, bookZip []byte) (map[string][]byte, error) {
	reader := bytes.NewReader(bookZip)
	zipReader, err := zip.NewReader(reader, int64(len(bookZip)))
	if err != nil {
		return nil, err
	}
	zipDecryptedContents := make(map[string][]byte)
	for _, file := range zipReader.File {
		if !strings.HasPrefix(file.Name, "p-") || strings.HasPrefix(file.Name, "p-l-") {
			continue
		}
		encryptedContents, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer encryptedContents.Close()
		contents, err := io.ReadAll(encryptedContents)
		if err != nil {
			return nil, err
		}
		decryptedContents, err := util.DecryptFile(keyPrefix, contents)
		if err != nil {
			return nil, err
		}
		zipDecryptedContents[file.Name] = decryptedContents
	}
	return zipDecryptedContents, nil
}

func convertBook(files map[string][]byte) error {
	fileNames := make([]string, 0, len(files))

	for f := range files {
		fileNames = append(fileNames, f)
	}
	slices.SortFunc(fileNames, func(a, b string) int {
		getPageNum := func(name string) int {
			lastDash := strings.LastIndex(name, "-")
			if lastDash == -1 {
				return 0
			}
			dot := strings.Index(name[lastDash:], ".")
			if dot == -1 {
				return 0
			}
			numStr := name[lastDash+1 : lastDash+dot]
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return 0
			}
			return num
		}
		return getPageNum(a) - getPageNum(b)
	})
	pdf := gofpdf.New("P", "mm", "A4", "")
	pageWidth, pageHeight := pdf.GetPageSize()

	for _, fileName := range fileNames {
		content := files[fileName]
		reader := bytes.NewReader(content)
		var imageType string
		kind, _ := filetype.MatchReader(reader)
		reader.Seek(0, io.SeekStart)
		switch kind.MIME.Value {
		case "image/webp":
			imageType = "webp"
		case "image/jpeg":
			imageType = "jpg"
		case "image/png":
			imageType = "png"
		default:
			log.Fatalf("Error: %s is not a webp or jpg or png file", fileName)
		}
		var imageContent []byte
		if imageType == "webp" {
			img, err := webp.Decode(reader)
			if err != nil {
				return err
			}
			var buf bytes.Buffer
			if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
				return err
			}
			imageContent = buf.Bytes()
			imageType = "jpg"
		} else {
			var err error
			imageContent, err = io.ReadAll(reader)
			if err != nil {
				return err
			}
		}

		pdf.RegisterImageOptionsReader(
			fileName,
			gofpdf.ImageOptions{ImageType: imageType, ReadDpi: true},
			bytes.NewReader(imageContent),
		)
		pdf.AddPage()
		pdf.ImageOptions(fileName, 0, 0, pageWidth, pageHeight, false,
			gofpdf.ImageOptions{ImageType: imageType, ReadDpi: true}, 0, "")
	}

	if err := pdf.OutputFileAndClose(book.Name + ".pdf"); err != nil {
		return err
	}

	return nil
}

func handleBook() {
	bookZip, err := downloadBook()
	if err != nil {
		fmt.Println(err)
		return
	}
	decryptedZip, err := decryptBook(publisher.KeyPrefix, bookZip)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = convertBook(decryptedZip)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var rootCmd = &cobra.Command{
	Use: "zkitap2pdf",
	Run: func(cmd *cobra.Command, args []string) {

		var key string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[types.Publisher]().Options(
					huh.NewOption("Benim Hocam", types.Publisher{
						DataURL:   "http://www.benimhocamdijital.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "23_11_22",
					}),
				).
					Title("Yayınevi Seç").
					Value(&publisher),
				huh.NewInput().Title("Anahtar").Validate(func(s string) error {
					if len(s) != 8 {
						return fmt.Errorf("anahtar 8 karakter olmalıdır")
					}
					return nil
				}).Value(&key),
			),
			huh.NewGroup(
				huh.NewSelect[int]().OptionsFunc(func() []huh.Option[int] {
					var err error
					publisherData, err = util.GetPublisherData(publisher, key)
					if err != nil {
						fmt.Println(err)
						return []huh.Option[int]{}
					}
					options := []huh.Option[int]{}
					for _, shelf := range publisherData.Shelfs {
						options = append(options, huh.NewOption(shelf.Name, shelf.ID))
					}
					return options
				}, &publisher).Title("Raf Seç").Value(&shelfId),
				huh.NewSelect[types.Book]().OptionsFunc(func() []huh.Option[types.Book] {
					options := []huh.Option[types.Book]{}
					if shelfId == -1 {
						return options
					}
					for _, book := range publisherData.Shelfs[shelfId-1].Books {
						options = append(options, huh.NewOption(book.Name, book))
					}
					return options
				}, &shelfId).Title("Kitap Seç").Value(&book),
			),
		)
		form.Run()

		if err := spinner.New().Title("Kitap indiriliyor...").Action(handleBook).Run(); err != nil {
			fmt.Println("Failed:", err)
			return
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
