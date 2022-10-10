package worker

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"time"
)

const dateFormatString = "2006-01-02 15:04:05"

var emptyParameter = Parameter{}

// ProcessInput converts input stream to a matrix ready, similar to csv.NewReader
func ProcessInput(stdin io.Reader) ([][]string, error) {
	scanner := bufio.NewScanner(stdin)
	records := make([][]string, 0)
	for scanner.Scan() {
		records = append(records, strings.Split(scanner.Text(), ","))
	}

	if len(records) == 1 && len(records[0]) == 1 {
		flags := strings.Split(records[0][0], " ")
		if flags[0] != "-f" || len(flags) == 1 {
			return nil, errors.New("issue with input")
		}

		return ImportCsv(flags[1])
	}
	return records, nil
}

// ImportCsv preps the input file for processing
func ImportCsv(pathName string) ([][]string, error) {
	file, err := os.Open(pathName)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// processImportedCsv takes data that has been already imported and assigned to a matrix to prepare for passing to workers
func processImportedCsv(importedCSV [][]string) ([]Parameter, error) {
	if len(importedCSV) == 0 || isRowAHeader(importedCSV[0]) && len(importedCSV) == 1 {
		return nil, errors.New("empty file")
	}

	startIndex := 0
	if isRowAHeader(importedCSV[0]) {
		startIndex = 1
	}

	params := make([]Parameter, 0)
	for _, row := range importedCSV[startIndex:] {
		param, err := convertArrayToParameter(row)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}
	return params, nil
}

// isRowAHeader attempts to parse the first incoming row. Any error results in assuming it is a header
func isRowAHeader(row []string) bool {
	_, err := convertArrayToParameter(row)
	if err != nil {
		return true
	}
	return false
}

func convertArrayToParameter(column []string) (Parameter, error) {
	if len(column) != 3 {
		return emptyParameter, errors.New("parsed string did not have the expected number of columns")
	}
	host := column[0]
	start, err := time.ParseInLocation(dateFormatString, column[1], time.UTC)
	if err != nil {
		return emptyParameter, err
	}

	end, err := time.ParseInLocation(dateFormatString, column[2], time.UTC)
	if err != nil {
		return emptyParameter, err
	}

	return Parameter{
		host:  host,
		start: start,
		end:   end,
	}, nil
}
