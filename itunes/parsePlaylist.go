package itunes

import (
	"encoding/csv"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/jeremyroberts0/itunes-to-spotify/types"
)

func ParsePaylist(target io.Reader) (songs []types.ItunesSong) {
	reader := csv.NewReader(target)
	reader.LazyQuotes = true
	reader.Comma = '\t'

	// var colPositions ColPositions
	data := []string{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, record...)
	}

	rows := [][]string{[]string{}}

	for _, value := range data {
		rowIndex := len(rows) - 1
		if strings.Contains(value, "\r") {
			split := strings.Split(value, "\r")
			rows[rowIndex] = append(rows[rowIndex], split[0])
			rows = append(rows, []string{split[1]})
		} else {
			rows[rowIndex] = append(rows[rowIndex], value)
		}
	}

	var colPositions ColPositions

	for index, row := range rows {
		if index == 0 {
			colPositions = getColPositions(row)
		} else {
			song := parseRow(row, colPositions)
			songs = append(songs, song)
		}
	}

	return songs
}

func parseRow(row []string, colPositions ColPositions) types.ItunesSong {
	return types.ItunesSong{
		Name:   getCol(row, colPositions["Name"]),
		Artist: getCol(row, colPositions["Artist"]),
		Album:  getCol(row, colPositions["Album"]),
	}
}

func getColPositions(row []string) ColPositions {
	cp := ColPositions{}
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	for pos, name := range row {
		scrubbed := reg.ReplaceAllString(name, "")
		cp[scrubbed] = pos
	}

	return cp
}

type ColPositions map[string]int

func getCol(row []string, index int) string {
	if len(row)-1 >= index {
		return row[index]
	}

	return ""
}
