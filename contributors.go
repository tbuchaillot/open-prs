package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ContributorPRCount struct {
	Repo  string
	Login string
	Count int
	URLs  []string // Add a slice of strings to hold PR URLs
}

type Contributors []ContributorPRCount

func (c Contributors) Len() int {
	return len(c)
}

const (
	Stdout = "stdout"
	CSV    = "csv"
)

func (c Contributors) Output(outputType string) {
	switch outputType {
	case CSV:
		// Write to a CSV file
		var file *os.File
		if _, err := os.Stat("output.csv"); os.IsNotExist(err) {
			file, err = os.Create("output.csv")
			if err != nil {
				fmt.Println("Error creating CSV file:", err)
				return
			}
		} else {
			file, err = os.OpenFile("output.csv", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error opening CSV file:", err)
				return
			}
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		titleRow := []string{"\n" + fmt.Sprintf("Open prs for repository: %s", c[0].Repo) + "\n"}
		if err := writer.Write(titleRow); err != nil {
			return
		}

		// Writing the header
		header := []string{"Repository", "Username", "PR Count", "URLs"}
		if err := writer.Write(header); err != nil {
			fmt.Println("Error writing header to CSV file:", err)
			return
		}

		// Writing the data
		for _, contributor := range c {
			row := []string{contributor.Repo, contributor.Login, strconv.Itoa(contributor.Count), strings.Join(contributor.URLs, "\n")}
			if err := writer.Write(row); err != nil {
				fmt.Println("Error writing to CSV file:", err)
				return
			}
		}
	default:
		fmt.Printf("Total Open Pull Requests contributors: %d\n", c.Len())

		// Write to stdout
		fmt.Println("Open Pull Requests per Contributor (sorted by decreasing PR count):")

		for _, contributor := range c {
			fmt.Printf("%s: %d\n", contributor.Login, contributor.Count)
			for _, url := range contributor.URLs {
				fmt.Println(url)
			}
		}
	}
}
