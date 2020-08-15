package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brianvoe/gofakeit"
)

func GenerateSampleCSV(numberOfLines int, outputFilePath string) {

	csvOptions := gofakeit.CSVOptions{
		Delimiter: ",",
		RowCount:  numberOfLines,
		Fields: []gofakeit.Field{
			{Name: "id", Function: "autoincrement"},
			{Name: "company", Function: "company"},
			{Name: "street", Function: "street"},
			{Name: "city", Function: "city"},
			{Name: "country", Function: "country"},
		},
	}

	bs, err := gofakeit.CSV(&csvOptions)
	if err != nil {
		log.Println("err:", err)
	}

	f, e := os.Create(outputFilePath)
	if e != nil {
		fmt.Println(e)
	}
	_, err = f.Write(bs)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	GenerateSampleCSV(1000, "companies.csv")
}
