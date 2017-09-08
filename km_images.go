package main

import (
	"errors"
	"fmt"
	"github.com/sclevine/agouti"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var urls = []string{}
var savePathBase = "/Users/akiraak/Downloads/**********/"

func getUrlContent(url string) ([]byte, error) {
	//fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	if resp.StatusCode != http.StatusOK {
		errorString := "Not StatusOK. URL:" + url
		return []byte{}, errors.New(errorString)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

func saveFile(data []byte, filePath string) error {
	dir, _ := path.Split(filePath)
	os.MkdirAll(dir, os.FileMode(0700))
	err := ioutil.WriteFile(filePath, data, os.FileMode(0600))
	return err
}

func main() {
	driver := agouti.ChromeDriver()
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Fatalf("Failed to open page:%v", err)
	}

	for i, url := range urls {
		fmt.Println(url)
		if err := page.Navigate(url); err != nil {
			log.Fatalf("Failed to navigate:%v", err)
		}
		if i == 0 {
			time.Sleep(time.Second * 5)
		}
		var imageUrls []string
		page.RunScript("return lstImages;", map[string]interface{}{}, &imageUrls)
		for id, url := range imageUrls {
			content, err := getUrlContent(url)
			if err == nil {
				path := fmt.Sprintf("%s/%02d/%03d.jpg", savePathBase, i+1, id+1)
				fmt.Println(path)
				saveFile(content, path)
			}
		}
	}
}
