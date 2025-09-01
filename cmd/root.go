package cmd

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/usdogu/zkitap2pdf/types"
	"github.com/usdogu/zkitap2pdf/util"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/grafana/gofpdf"
	"github.com/h2non/filetype"
	"github.com/spf13/cobra"
	"golang.org/x/image/webp"
)

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

func handleBook(_ context.Context) error {
	bookZip, err := downloadBook()
	if err != nil {
		return err
	}
	decryptedZip, err := decryptBook(publisher.KeyPrefix, bookZip)
	if err != nil {
		return err
	}
	err = convertBook(decryptedZip)
	if err != nil {
		return err
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use: "zkitap2pdf",
	Run: func(cmd *cobra.Command, args []string) {

		var key string
		var publisherData types.PublisherData
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[types.Publisher]().Options(
					huh.NewOption("Açı Yayınları", types.Publisher{
						DataURL:   "http://dijital.aciyayinlari.com.tr/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "21_11_22",
					}),

					huh.NewOption("Ankara Yayıncılık", types.Publisher{
						DataURL:   "https://onlineankarayayincilik.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "6_13_26",
					}),

					huh.NewOption("Arı Yayıncılık", types.Publisher{
						DataURL:   "https://akillidefter.com.tr/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "7_14_21",
					}),

					huh.NewOption("ATA", types.Publisher{
						DataURL:   "http://www.ataogretmen.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "1_2_6",
					}),

					huh.NewOption("Benim Hocam", types.Publisher{
						DataURL:   "http://www.benimhocamdijital.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "23_11_22",
					}),

					huh.NewOption("Berkay Yayıncılık", types.Publisher{
						DataURL:   "http://berkayokul.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "26_7_14",
					}),

					huh.NewOption("Çalışkan Arı", types.Publisher{
						DataURL:   "http://www.caliskanarizkitap.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "6_15_7",
					}),

					huh.NewOption("Çanta Yayıncılık", types.Publisher{
						DataURL:   "http://www.cantadaakillitahta.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "2_12_24",
					}),

					huh.NewOption("Endemik Yayınları", types.Publisher{
						DataURL:   "http://akillitahta.endemikyayinlari.com.tr/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "22_1_2",
					}),

					huh.NewOption("Fi Yayınları", types.Publisher{
						DataURL:   "http://www.fiakillitahta.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "7_3_6",
					}),

					huh.NewOption("Final Kurumsal Yayinlari", types.Publisher{
						DataURL:   "http://finaldijital.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "2_1_21",
					}),

					huh.NewOption("Gizli Yayınları", types.Publisher{
						DataURL:   "http://gizliokul.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "1_21_6",
					}),

					huh.NewOption("Günay Yayınları", types.Publisher{
						DataURL:   "http://www.gunayportal.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "9_4_8",
					}),

					huh.NewOption("Hiper Zeka", types.Publisher{
						DataURL:   "http://www.hiperzekadijital.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "8_4_8",
					}),

					huh.NewOption("İşler", types.Publisher{
						DataURL:   "http://yeni.isleronline.com/controller/zkitap/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "28_15_5",
					}),

					huh.NewOption("Limit", types.Publisher{
						DataURL:   "http://kurumsal.limitheryerde.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "4_17_5",
					}),

					huh.NewOption("Model Eğitim", types.Publisher{
						DataURL:   "http://www.modelogretmen.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "7_14_22",
					}),

					huh.NewOption("More And More", types.Publisher{
						DataURL:   "http://zkitap.kurmayokul.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "1_4_8",
					}),

					huh.NewOption("Mozaik Yayınları", types.Publisher{
						DataURL:   "http://mozaikakillitahta.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "23_12_24",
					}),

					huh.NewOption("Orijinal Matematik", types.Publisher{
						DataURL:   "http://yeni.isleronline.com/controller/zkitap/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "28_15_5",
					}),

					huh.NewOption("Paraf Akademi", types.Publisher{
						DataURL:   "https://kutuphane.parafakademi.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "2_26_3",
					}),

					huh.NewOption("PRF Yayınları", types.Publisher{
						DataURL:   "http://prfkutuphane.prfyayinlari.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "8_1_2",
					}),

					huh.NewOption("Toprak Yayıncılık", types.Publisher{
						DataURL:   "http://toprakogretmen.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "9_5_10",
					}),

					huh.NewOption("ÜçDörtBeş Yayınları", types.Publisher{
						DataURL:   "http://www.345dijitalicerik.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "1_10_20",
					}),

					huh.NewOption("Üçgen", types.Publisher{
						DataURL:   "http://www.e-ucgen.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "17_18_20",
					}),

					huh.NewOption("Ünlüler Karması", types.Publisher{
						DataURL:   "http://unlulerakillitahta.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "26_12_24",
					}),

					huh.NewOption("X Yayıncılık", types.Publisher{
						DataURL:   "http://onlinexyayincilik.com/data_json.php?v=zkitapx&action=check_key&type=0&key=%s&subPrefix",
						KeyPrefix: "2_1_3",
					}),
				).
					Title("Yayınevi Seç").
					Value(&publisher),
				huh.NewInput().Title("Anahtar").Validate(func(s string) error {
					if len(s) != 8 {
						return fmt.Errorf("anahtar 8 karakter olmalıdır")
					}
					var err error
					publisherData, err = util.GetPublisherData(publisher, key)
					if err != nil {
						return err
					}
					if !publisherData.Status {
						return fmt.Errorf("anahtar yanlış")
					}
					return nil
				}).Value(&key),
			),

			huh.NewGroup(
				huh.NewSelect[int]().OptionsFunc(func() []huh.Option[int] {
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
		err := form.Run()
		if err == huh.ErrUserAborted {
			os.Exit(0)
		}

		if err := spinner.New().Title("Kitap indiriliyor...").ActionWithErr(handleBook).Run(); err != nil {
			fmt.Println("Hata:", err)
			fmt.Scanln()
			return
		}
		cwd, _ := os.Getwd()
		fmt.Printf("Kitap %s konumuna kaydedildi. Çıkmak için bir tuşa basın.\n", path.Join(cwd, book.Name+".pdf"))
		fmt.Scanln()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
