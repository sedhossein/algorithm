package main

import (
	"fmt"
	"strconv"
	"strings"
)

var itemSetSeparator = "//"

func MakeTable(records map[string][]int) *Table {
	newTable := &Table{Records: make(map[string][]Item)}
	for itemset, items := range records {
		for _, item := range items {
			newTable.Records[itemset] = append(newTable.Records[itemset], Item(strconv.Itoa(item)))
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

func (e Eclat) eclat(table *Table) *Table {
	newTable := &Table{Records: make(map[string][]Item)}
	var firstIndex = 0
	var secondIndex int
	for firstIS, firstItems := range table.Records {
		secondIndex = firstIndex + 1
		c := -1
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
			secondIndex++
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
