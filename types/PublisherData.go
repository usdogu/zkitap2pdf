package types

type PublisherData struct {
	Shelfs []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Books []Book `json:"books"`
	} `json:"shelfs"`
}

type Book struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	DataFile string `json:"data_file"`
	DataZip  string `json:"data_zip"`
}
