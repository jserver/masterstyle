package main

import (
	"fmt"
	"strings"
)

type Row []string

type Table struct {
	Header []string
	Rows   []Row
}

func PrintTable(table *Table) {
	cols := len(table.Header)
	lens := make([]int, cols)

	for idx, hdr := range table.Header {
		lens[idx] = len(hdr)
	}

	for _, row := range table.Rows {
		for idx, col := range row {
			if len(col) > lens[idx] {
				lens[idx] = len(col)
			}
		}
	}

	total := 2
	for _, itm := range lens {
		total += itm + 2
	}

	dashLine := "+"
	for _, size := range lens {
		dashLine += fmt.Sprintf("%s+", strings.Repeat("-", size + 2))
	}

	fmt.Println(dashLine)
	fmt.Print("|")
	for idx, hdr := range table.Header {
		str := hdr + strings.Repeat(" ", lens[idx] - len(hdr))
		fmt.Printf(" %s |", str)
	}
	fmt.Println()

	fmt.Println(dashLine)
	for _, row := range table.Rows {
		fmt.Print("|")
		for idx, col := range row {
			str := col + strings.Repeat(" ", lens[idx] - len(col))
			fmt.Printf(" %s |", str)
		}
		fmt.Println()
	}
	fmt.Println(dashLine)
}

/*
func main() {
	header := []string{"ID", "Name", "URL", "Description"}
	rowA := Row{"1234", "Joe", "/home/joe", "Here it is cocky"}
	rowB := Row{"123456", "Joe", "/home/joe", "Here it is"}
	rowC := Row{"1234", "Joe Server", "/home/joe", "Here it is"}
	rowD := Row{"1234", "Joe", "/home/joe grab some", "Here it is"}
	rows := []Row{rowA, rowB, rowC, rowD}
	table := Table{header, rows}
	PrintTable(&table)
}
*/
