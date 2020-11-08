package main

import (
	"bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

func testPublish(base string, dataId string, group string, content string) {
    // Publish Config
    publishUrl, err := url.Parse(base)
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return
    }

    params := url.Values{}
	params.Add("dataId", dataId)
    params.Add("group", group)
    params.Add("content", content)
    publishUrl.RawQuery = params.Encode()

    jsonData := map[string]string{}
    jsonValue, _ := json.Marshal(jsonData)

    publishRsp, err := http.Post(publishUrl.String(), "application/json", bytes.NewBuffer(jsonValue))
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(publishRsp.Body)
        fmt.Printf("Publish Result is status %s , value %s\n", publishRsp.Status, string(data))
    }
}

func testGet(base string, dataId string, group string) {
    // Get Config
    getUrl, err := url.Parse(base)
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return
    }

    getParams := url.Values{}
	getParams.Add("dataId", dataId)
    getParams.Add("group", group)

    getUrl.RawQuery = getParams.Encode()

    getRsp, err := http.Get(getUrl.String())
    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(getRsp.Body)
        fmt.Printf("Get config Result is status %s , value %s\n", getRsp.Status, string(data))
    }
}

func testDelete(base string, dataId string, group string) {
    // Get Config
    deleteUrl, err := url.Parse(base)
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return
    }

    deleteParams := url.Values{}
	deleteParams.Add("dataId", dataId)
    deleteParams.Add("group", group)
    deleteUrl.RawQuery = deleteParams.Encode()

    jsonData := map[string]string{}
    jsonValue, _ := json.Marshal(jsonData)

    // Delete Config
    req, err := http.NewRequest(http.MethodDelete, deleteUrl.String(), bytes.NewBuffer(jsonValue))
    client := &http.Client{}
    deleteResp, err := client.Do(req)

    if err != nil {
        fmt.Printf("The HTTP request failed with error %s\n", err)
    } else {
        data, _ := ioutil.ReadAll(deleteResp.Body)
        fmt.Printf("Delete config Result is status %s , value %s\n", deleteResp.Status, string(data))
    }
}

func main() {
    fmt.Println("Starting the application...")

    base := "http://localhost:8080/sweets/v1/cs/configs"
    dataId := "mydata"
    group := "mygroup"
    content := "hello world!"

    testPublish(base, dataId, group, content) 
    testGet(base, dataId, group)
    testDelete(base, dataId, group)
    testPublish(base, dataId, group, content)

	fmt.Println("Terminating the application...")
}