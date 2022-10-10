package worker

import (
	"io"
	"strings"
	"testing"
)

func Test_processInput(t *testing.T) {
	input := `hostname,start_time,end_time
host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000001,2017-01-02 13:02:02,2017-01-02 14:02:02
host_000008,2017-01-02 18:50:28,2017-01-02 19:50:28`
	var reader io.Reader = strings.NewReader(input)
	records, err := ProcessInput(reader)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 4 {
		t.Fatal()
	}
	if len(records[0]) != 3 {
		t.Fatal()
	}

}
func Test_processInputAsCSV(t *testing.T) {
	input := `-f ../resources/query_params_no_header.csv`
	var reader io.Reader = strings.NewReader(input)
	records, err := ProcessInput(reader)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 5 {
		t.Fatal(len(records))
	}
	if len(records[0]) != 3 {
		t.Fatal()
	}
}

func Test_processBadInput(t *testing.T) {
	input := `-f`
	var reader io.Reader = strings.NewReader(input)
	_, err := ProcessInput(reader)
	if err == nil {
		t.Fatal("There should be an error")
	}

}
