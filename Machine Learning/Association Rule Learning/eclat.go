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

var itemSetSeparator = "//"

type Eclat struct {
	table      *Table
	MinSupport int
}

type Table struct {
	Records map[string][]Item
}

// Item is transaction ID (item)
type Item string

func (e Eclat) itemSetNormalizer(s ...string) string {
	return strings.Trim(strings.Join(s, itemSetSeparator), itemSetSeparator)
}

func (e Eclat) GetItems(table Table, itemset ...string) []Item {
	return table.Records[e.itemSetNormalizer(itemset...)]
}

func MakeTable(records map[string][]int) *Table {
	newTable := &Table{}
	for itemset, items := range records {
		for _, item := range items {
			newTable.Records[itemset] = append(newTable.Records[itemset], Item(item))
		}
	}

	return newTable
}

func Run(table *Table, ms int) (*Table, error) {
	if table == nil {
		return nil, fmt.Errorf("empty table received")
	}

	e := Eclat{
		table:      table,
		MinSupport: ms,
	}

	return e.eclat(e.table), nil
}

func (e Eclat) eclat(table *Table) *Table {
	newTable := new(Table)
	var firstIndex = 0
	var secondIndex int
	for firstIS, firstItems := range table.Records {
		secondIndex = firstIndex + 1
		c := 0
		for secondIS, secondItems := range table.Records {
			c++
			if c < secondIndex {
				continue
			}

			var newItems []Item
			for _, firstItem := range firstItems {
				for _, secondItem := range secondItems {
					if firstItem == secondItem {
						newItems = append(newItems, firstItem) // or secondItem! both of them are equal
					}
				}
			}

			newIS := firstIS + itemSetSeparator + secondIS
			newTable.Records[newIS] = newItems
		}

		firstIndex++
	}

	for i, itemSet := range newTable.Records {
		// ignore item-sets that they length are smaller than the MinSupport
		if len(itemSet) >= e.MinSupport {
			newTable.Records[i] = itemSet
		}
	}

	if len(newTable.Records) == 0 {
		return table
	}

	return e.eclat(newTable)
}

// =================================================================
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
	//fmt.Println(items)
	fmt.Println(res)
}
