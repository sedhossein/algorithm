package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// readCsvFile consider line id as a transaction id (T_ID)
// => horizontal data format
// TID   |   items
//  1    |    word1 word2 word3
//  2    |    word4
//  3    |    word1 word3
func readCsvFile(path string) (map[int][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(file))
	records := make(map[int][]string)
	step := 0
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil && !strings.Contains(err.Error(), csv.ErrFieldCount.Error()) {
			return nil, err
		}

		records[step] = append(records[step], line...)
		step++
	}

	return records, nil
}

// => vertical data format
//  item     |   T_ID s
//  word1    |   1 3
//  word2    |    1
//  word3    |   1 3
func normalize(records map[int][]string) (items map[string][]int) {
	items = make(map[string][]int)
	// `booked` help to prevent inserting duplicate T_ID in each item
	booked := make(map[string]bool)
	for i, record := range records {
		for _, item := range record {
			item = strings.TrimSpace(item)
			key := item + strconv.Itoa(i)
			if ok := booked[key]; !ok {
				items[item] = append(items[item], i)
				booked[key] = true
			}
		}
	}

	return
}

func main() {
	records, err := readCsvFile("test.csv")
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	items := normalize(records)
	table := MakeTable(items)
	res, err := Run(table, 2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("elcat response time: ", time.Now().Sub(start).String())
	fmt.Println(res)
}
