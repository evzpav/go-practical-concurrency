package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Download struct {
	Url           string
	TargetPath    string
	TotalSections int
}

var fileToDownload string = "https://golang.org/dl/go1.15.linux-amd64.tar.gz"
var fileName string = "go1.15.linux-amd64.tar.gz"

func main() {
	startTime := time.Now()
	d := Download{
		// Provide the URL to download,
		//	example: "https://golang.org/dl/go1.15.linux-amd64.tar.gz"
		Url: fileToDownload,
		// Provide the target file path with extension, example: sample.mp4
		TargetPath: fileName,
		// Number of sections/connections to make to the server
		TotalSections: 20,
	}
	err := d.Do()
	if err != nil {
		log.Printf("An error occured while downloading the file: %s\n", err)
	}
	fmt.Printf("Download completed in %v seconds\n", time.Now().Sub(startTime).Seconds())
}

func (d Download) Do() error {
	fmt.Println("Checking URL")
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	fmt.Printf("Got %v\n", resp.StatusCode)

	if resp.StatusCode > 299 {
		return fmt.Errorf("Can't process, response is %v", resp.StatusCode)
	}

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	fmt.Printf("Size is %v bytes\n", size)

	sections := d.splitSessions(size)

	fmt.Printf("Sections: %v\n", sections)

	var wg sync.WaitGroup
	for i, s := range sections {
		wg.Add(1)
		go func(i int, s [2]int) {
			defer wg.Done()
			err = d.downloadSection(i, s)
			if err != nil {
				panic(err)
			}
		}(i, s)
	}
	wg.Wait()

	return d.mergeFiles(sections)
}

func (d Download) downloadSection(i int, c [2]int) error {
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", c[0], c[1]))

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("Can't process, response is %v", resp.StatusCode)
	}
	fmt.Printf("Downloaded %v bytes for section %v\n", resp.Header.Get("Content-Length"), i)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(fmt.Sprintf("section-%v.tmp", i), b, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// Get a new http request
func (d Download) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		d.Url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	r.Header.Set("User-Agent", "Silly Download Manager v001")
	return r, nil
}

func (d Download) mergeFiles(sections [][2]int) error {
	f, err := os.OpenFile(d.TargetPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	for i := range sections {
		tmpFileName := fmt.Sprintf("section-%v.tmp", i)
		b, err := ioutil.ReadFile(tmpFileName)
		if err != nil {
			return err
		}

		n, err := f.Write(b)
		if err != nil {
			return err
		}

		if err = os.Remove(tmpFileName); err != nil {
			return err
		}

		fmt.Printf("%v bytes merged\n", n)
	}

	return nil
}

func (d Download) splitSessions(size int) [][2]int {
	var sections = make([][2]int, d.TotalSections)
	eachSize := size / d.TotalSections
	fmt.Printf("Each size is %v bytes\n", eachSize)

	// example: if file size is 100 bytes, our section should like:
	// [[0 10] [11 21] [22 32] [33 43] [44 54] [55 65] [66 76] [77 87] [88 98] [99 99]]
	for i := range sections {
		if i == 0 {
			// starting byte of first section
			sections[i][0] = 0
		} else {
			// starting byte of other sections
			sections[i][0] = sections[i-1][1] + 1
		}

		if i < d.TotalSections-1 {
			// ending byte of other sections
			sections[i][1] = sections[i][0] + eachSize
		} else {
			// ending byte of other sections
			sections[i][1] = size - 1
		}
	}

	return sections
}
