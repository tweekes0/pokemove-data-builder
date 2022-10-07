package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)


func getResponse(url string) ([]byte, error) {
	var data []byte

	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	return data, nil
}

func getBasicMoveUrls(limit int, endpoint string) ([]string, error) {
	var moveUrls []string
	var basicResp BasicMoveResponse

	url := fmt.Sprintf("%v?limit=%v",endpoint, limit)
	
	data, err := getResponse(url)
	if err != nil {
		return moveUrls, err
	}
	json.Unmarshal(data, &basicResp)

	for _, res := range basicResp.Results {
		moveUrls = append(moveUrls, res.Url)
	}	

	return moveUrls, nil
}