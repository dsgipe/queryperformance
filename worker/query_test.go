package worker

import (
	"testing"
	"time"
)

// Benchmark tests to confirm workers are actually running in parallel and to see time improvement
func BenchmarkTest1Worker(b *testing.B) {
	numberOfWorkers := 1
	csv, _ := ImportCsv("../resources/query_params.csv")
	importedCsv, _ := processImportedCsv(csv)
	_, err := submitJobsToWorkers(importedCsv, numberOfWorkers)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkTest10Workers(b *testing.B) {
	numberOfWorkers := 10
	csv, _ := ImportCsv("../resources/query_params.csv")
	importedCsv, _ := processImportedCsv(csv)
	_, err := submitJobsToWorkers(importedCsv, numberOfWorkers)
	if err != nil {
		b.Fatal(err)
	}
}

func Test_importToGenerateQueryJourney(t *testing.T) {
	csv, err := ImportCsv("../resources/query_params.csv")
	if err != nil {
		t.Fatal(err)
	}
	importedCsv, err := processImportedCsv(csv)
	if err != nil {
		t.Fatal(err)
	}
	if len(importedCsv) == 0 {
		t.Fatal("Imported empty csv")
	}
	if len(importedCsv) != 200 {
		t.Fatal("Issue importing data", len(importedCsv))
	}

	numberOfWorkers := 1
	compiledResults, err := submitJobsToWorkers(importedCsv, numberOfWorkers)
	if len(compiledResults) != 200 {
		t.Fatal(len(compiledResults))
	}

	hostNames := map[string]int{}

	for _, stats := range compiledResults {
		if workerId, ok := hostNames[stats.hostName]; ok {
			if workerId != stats.workerId {
				t.Fatal("A host was run on multiple workers!")
			}
		} else {
			hostNames[stats.hostName] = workerId
		}
	}

	benchmarkStat := calculateBenchmarkStatistics(compiledResults)
	if benchmarkStat.numberOfQueriesRun != 200 {
		t.Fatal("incorrect number of queries run", benchmarkStat.numberOfQueriesRun)
	}

	standardMetrics(t, benchmarkStat)
	benchmarkStat.report()
}

func Test_importNoHeader(t *testing.T) {
	csv, err := ImportCsv("../resources/query_params_no_header.csv")
	if err != nil {
		t.Fatal(err)
	}
	importedCsv, err := processImportedCsv(csv)
	if len(importedCsv) != 5 {
		t.Fatal(len(importedCsv))
	}
	numberOfWorkers := 2
	compiledResults, err := submitJobsToWorkers(importedCsv, numberOfWorkers)
	if len(compiledResults) != 5 {
		t.Fatal(len(compiledResults))
	}
	benchmarkStat := calculateBenchmarkStatistics(compiledResults)
	if benchmarkStat.numberOfQueriesRun != 5 {
		t.Fatal("incorrect number of queries run", benchmarkStat.numberOfQueriesRun)
	}
	standardMetrics(t, benchmarkStat)
	benchmarkStat.report()

}

func standardMetrics(t *testing.T, stat benchmarkStatistics) {
	if stat.minimumQueryTime == 0 {
		t.Fatal("error in calculating standard metrics", stat.minimumQueryTime)
	}
	if stat.minimumQueryTime > stat.maximumQueryTime {
		t.Fatal("error in calculating standard metrics", stat.maximumQueryTime)
	}
	if stat.averageQueryTime > float64(stat.maximumQueryTime.Milliseconds()) {
		t.Fatal("error in calculating standard metrics", stat.averageQueryTime)
	}
	if stat.medianQueryTime > stat.maximumQueryTime {
		t.Fatal("error in calculating standard metrics", stat.medianQueryTime)
	}
}

func Test_importNoData(t *testing.T) {
	csv, err := ImportCsv("../resources/query_params_header_no_data.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, err = processImportedCsv(csv)
	if err == nil {
		t.Fatal("error not found")
	}

}

func Test_query(t *testing.T) {
	cpuUsages, _, err := queryUsages(Parameter{
		host:  "host_000008",
		start: time.Date(2017, 01, 01, 8, 59, 22, 0, time.UTC),
		end:   time.Date(2017, 01, 01, 9, 59, 22, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	//TODO currently assume buckets are set to 0 seconds in the interval
	bucketTime := time.Date(2017, 01, 01, 8, 59, 0, 0, time.UTC)
	usage := cpuUsages[0]
	if usage.min != 27.54 || usage.max != 51.01 || usage.bucket.Unix() != bucketTime.Unix() {
		t.Fatal(usage)
	}
}
