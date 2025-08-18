package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"zkitap2pdf/types"
)

func GetPublisherData(publisher types.Publisher, key string) (types.PublisherData, error) {
	link := fmt.Sprintf(publisher.DataURL, key)
	fmt.Println(link)
	resp, err := http.Get(link)
	if err != nil {
		return types.PublisherData{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.PublisherData{}, err
	}
	var publisherData types.PublisherData
	err = json.Unmarshal(body, &publisherData)
	if err != nil {
		return types.PublisherData{}, err
	}

	return publisherData, nil
}
