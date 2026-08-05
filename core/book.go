package core

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/usdogu/zkitap2pdf/types"
	"github.com/usdogu/zkitap2pdf/util"

	"github.com/grafana/gofpdf"
	"golang.org/x/sync/errgroup"
)

func downloadBook(book types.Book) ([]byte, error) {
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

func convertBook(book types.Book, files map[string][]byte) error {
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
	images := make([]Image, len(fileNames))

	g := new(errgroup.Group)
	g.SetLimit(runtime.NumCPU())

	for i, name := range fileNames {

		g.Go(func() error {
			img, err := convertImage(name, files[name])
			if err != nil {
				return err
			}

			images[i] = img
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	w, h := pdf.GetPageSize()

	for _, img := range images {
		pdf.RegisterImageOptionsReader(
			img.Name,
			gofpdf.ImageOptions{
				ImageType: img.Type,
				ReadDpi:   true,
			},
			bytes.NewReader(img.Data),
		)

		pdf.AddPage()
		pdf.ImageOptions(
			img.Name,
			0, 0,
			w, h,
			false,
			gofpdf.ImageOptions{
				ImageType: img.Type,
				ReadDpi:   true,
			},
			0,
			"",
		)
	}
	if err := pdf.OutputFileAndClose(book.Name + ".pdf"); err != nil {
		return err
	}

	return nil
}

func ProcessBook(publisher types.Publisher, book types.Book) error {
	bookZip, err := downloadBook(book)
	if err != nil {
		return err
	}
	decryptedZip, err := decryptBook(publisher.KeyPrefix, bookZip)
	if err != nil {
		return err
	}
	err = convertBook(book, decryptedZip)
	if err != nil {
		return err
	}
	return nil
}
