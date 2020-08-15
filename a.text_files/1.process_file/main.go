package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type task struct {
	id         int
	lineNumber int
	lineValues []string
}

var startTime time.Time

func init() {
	startTime = time.Now()
}

func main() {
	fmt.Println("[main] main() started")

	numberOfWorkers := 10
	fileName := "companies.csv"
	outputFileName := "companies_output.csv"

	records, err := ReadCsvFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[main] File has %v lines\n", len(records))

	var wg sync.WaitGroup

	tasks := make(chan task, len(records))
	results := make(chan task, len(records))

	for i := 1; i <= numberOfWorkers; i++ {
		wg.Add(1)
		go processLineWorker(&wg, tasks, results, i)
	}

	for i, line := range records {
		tasks <- task{lineNumber: i, lineValues: line}
	}

	fmt.Println("[main] Wrote tasks")

	close(tasks)

	wg.Wait()

	// receving results from all workers
	endResult := make([][]string, len(records))
	for i := 0; i < len(records); i++ {
		result := <-results // non-blocking because buffer is non-empty
		endResult[result.lineNumber] = result.lineValues
		fmt.Println("[main] Result", i, ":", result)
	}

	if err = WriteCsvFile(outputFileName, endResult); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[main] main() stopped. Duration: %v\n", time.Since(startTime))
}

func ReadCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return [][]string{}, fmt.Errorf("Unable to read input file %v: %v", filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return [][]string{}, fmt.Errorf("Unable to parse file as CSV for %v: %v ", filePath, err)
	}

	return records, nil
}

func WriteCsvFile(name string, records [][]string) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	writer := csv.NewWriter(f)

	return writer.WriteAll(records)
}

func processLineWorker(wg *sync.WaitGroup, tasks <-chan task, results chan<- task, workerIndex int) {
	for t := range tasks {
		fmt.Printf("[worker %v] Processing line %v \n", workerIndex, t.lineNumber)
		t.lineValues = doProcess(t.lineValues, "suffix")
		results <- t
	}
	wg.Done()
}

func doProcess(line []string, suffix string) []string {
	line = append(line, suffix)
	time.Sleep(1 * time.Millisecond)
	return line
}
