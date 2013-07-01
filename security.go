package main

import (
	"fmt"
)

func ListSecurityGroups() {
	resp, err := conn.SecurityGroups(nil, nil)
	if err != nil {
		fmt.Println("Unable to get SecurityGroups", err)
		return
	}
	header := []string{"Name", "Description"}

	rows := make([]Row, len(resp.Groups))
	for idx, group := range resp.Groups {
		rows[idx] = Row{
			group.Name,
			group.Description,
		}
	}
	table := Table{header, rows}
	PrintTable(&table)
}
