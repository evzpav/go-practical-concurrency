package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func main() {
	fmt.Println("read file line by line - started")

	filePath := "companies.csv"
	outputFilePath := "companies_output.csv"
	totalLinesToProcess := 100
	numberOfWorkers := 10
	buffer := 0

	outputChan := make(chan []string, buffer)
	done := make(chan bool)

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Unable to read input file %v: %v", filePath, err)
		os.Exit(1)
	}
	defer f.Close()

	// READ LINE
	inputChan := readLines(totalLinesToProcess, f)

	// PROCESS
	process(inputChan, numberOfWorkers, outputChan)

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// WRITE OUTPUT TO FILE
	writeLines(outputChan, outputFile, done)

	<-done
	fmt.Printf("read file line by line - stopped. Duration: %v\n", time.Since(startTime))
}

func readLines(totalLinesToProcess int, f *os.File) <-chan []string {
	out := make(chan []string)
	csvReader := csv.NewReader(f)

	go func() {
		for i := 0; i < totalLinesToProcess; i++ {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Printf("Unable to parse file as CSV for %v: %v ", f.Name(), err)
				continue
			}

			out <- line
		}
		close(out)
	}()

	return out
}

func process(in <-chan []string, numberOfWorkers int, out chan<- []string) {
	var wg sync.WaitGroup

	for i := 0; i < numberOfWorkers; i++ {

		wg.Add(1)
		go func(i int) {
			suffix := fmt.Sprintf("worker_%d", i)

			for line := range in {
				out <- doProcess(line, suffix)
			}
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
}

func doProcess(line []string, suffix string) []string {
	line = append(line, suffix)
	time.Sleep(10 * time.Millisecond)
	return line
}

func writeLines(out chan []string, outputFile *os.File, done chan bool) {
	writer := csv.NewWriter(outputFile)

	go func() {
		for line := range out {
			fmt.Printf("out: %v\n", line)

			if err := writer.Write(line); err != nil {
				fmt.Printf("err: %v", err)
				return
			}
			writer.Flush()
		}
		done <- true
	}()

}
