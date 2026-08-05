package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"charm.land/huh/v2"
	"charm.land/huh/v2/spinner"
	"github.com/usdogu/zkitap2pdf/core"
	"github.com/usdogu/zkitap2pdf/data"
	"github.com/usdogu/zkitap2pdf/types"
	"github.com/usdogu/zkitap2pdf/util"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "zkitap2pdf",
	Run: func(cmd *cobra.Command, args []string) {

		var publisher types.Publisher
		var shelfId *int
		var key string
		var publisherData types.PublisherData
		var book *types.Book

		publisherOptions := []huh.Option[types.Publisher]{}
		for _, p := range data.Publishers {
			publisherOptions = append(publisherOptions, huh.NewOption(p.Name, p.Publisher))
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[types.Publisher]().
					Options(publisherOptions...).
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
				huh.NewSelect[*int]().OptionsFunc(func() []huh.Option[*int] {
					options := []huh.Option[*int]{}
					for _, shelf := range publisherData.Shelfs {
						options = append(options, huh.NewOption(shelf.Name, &shelf.ID))
					}
					return options
				}, &publisher).Title("Raf Seç").Value(&shelfId),
				huh.NewSelect[*types.Book]().OptionsFunc(func() []huh.Option[*types.Book] {
					options := []huh.Option[*types.Book]{}
					if shelfId == nil {
						return options
					}
					for _, shelf := range publisherData.Shelfs {
						if shelf.ID == *shelfId {
							for _, book := range shelf.Books {
								if len(book.List) > 0 {
									for _, subBook := range book.List {
										options = append(options, huh.NewOption(subBook.Name, &subBook))
									}
								} else {
									options = append(options, huh.NewOption(book.Name, &book))
								}
							}
							break
						}
					}
					return options
				}, &shelfId).Title("Kitap Seç").Value(&book),
			),
		)
		err := form.Run()
		if err == huh.ErrUserAborted {
			os.Exit(0)
		}

		if err := spinner.New().Title("Kitap indiriliyor...").ActionWithErr(func(_ context.Context) error {
			return core.ProcessBook(publisher, *book)
		}).Run(); err != nil {
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
