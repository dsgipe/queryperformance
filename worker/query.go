package worker

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

// numberOfWorkers can be updated with env variable NUMBER_OF_WORKERS
var numberOfWorkers = 8

// EvaluateTimescaleQuerySpeeds is the main control code for handling data and managing workers
func EvaluateTimescaleQuerySpeeds(csv [][]string) error {
	importedCsv, err := processImportedCsv(csv)
	if err != nil {
		return err
	}
	if len(importedCsv) == 0 {
		return err
	}

	if envConnStr, exists := os.LookupEnv("NUMBER_OF_WORKERS"); exists {
		if envNumberOfWorkers, err := strconv.Atoi(envConnStr); err == nil {
			numberOfWorkers = envNumberOfWorkers
		} else {
			fmt.Println("issue converting NUMBER_OF_WORKERS", err)
		}

	}

	compiledResults, err := submitJobsToWorkers(importedCsv, numberOfWorkers)
	if err != nil {
		return err
	}

	benchmarkStat := calculateBenchmarkStatistics(compiledResults)
	benchmarkStat.report()
	return nil
}

// Creates query for the max cpu usage and min cpu usage of the
// given hostname for every minute in the time range specified by the start time and end time.
func queryUsages(param Parameter) ([]cpuUsage, time.Duration, error) {
	ctx := context.Background()

	startTime := time.Now()
	rows, err := TimescaleConn().Query(ctx, `SELECT time_bucket('1 minute', ts) AS bucket,
				   min(usage),
				   max(usage)
			FROM cpu_usage
			where host=$1 AND ts between $2 and $3
			GROUP BY bucket
			ORDER BY bucket ASC;`,
		param.host, param.start, param.end)
	defer rows.Close()

	if err != nil {
		return nil, 0, err
	}
	cpuUsages := make([]cpuUsage, 0)
	for rows.Next() {

		var bucket time.Time
		var min, max float64
		err = rows.Scan(&bucket, &min, &max)
		if err != nil {
			return nil, 0, err
		}
		cpuUsages = append(cpuUsages, cpuUsage{
			hostname: param.host,
			min:      min,
			max:      max,
			bucket:   bucket,
		})
		if rows.Err() != nil {
			return nil, 0, rows.Err()
		}
	}

	return cpuUsages, time.Since(startTime), nil
}

func calculateBenchmarkStatistics(stats []statistics) benchmarkStatistics {
	benchmarkStats := benchmarkStatistics{}
	for i := range stats {
		benchmarkStats.addCurrent(stats[i].processingTime)
	}
	benchmarkStats.processAverages()
	return benchmarkStats
}

// processAverages must be run before returning in order to capture the averages for the market
func (stat *benchmarkStatistics) processAverages() {
	stat.averageQueryTime = float64(stat.totalProcessingTimeForAllQueries.Milliseconds()) / float64(stat.numberOfQueriesRun)

	stat.medianQueryTime = findMedian(stat.allTimes, stat.numberOfQueriesRun)
}

func findMedian(times []time.Duration, numberOfQueriesRun int) time.Duration {

	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
	if numberOfQueriesRun%2 == 0 {
		return (times[numberOfQueriesRun/2] + times[numberOfQueriesRun/2+1]) / 2
	} else {
		return times[numberOfQueriesRun/2]
	}

}

// addCurrent adds to the total of TotalAssets
func (stat *benchmarkStatistics) addCurrent(processingTime time.Duration) {
	stat.numberOfQueriesRun++
	stat.allTimes = append(stat.allTimes, processingTime)
	stat.totalProcessingTimeForAllQueries += processingTime

	if processingTime < stat.minimumQueryTime || stat.minimumQueryTime == 0 {
		stat.minimumQueryTime = processingTime
	}
	if processingTime > stat.maximumQueryTime {
		stat.maximumQueryTime = processingTime
	}
}

// report all requested information
func (stat *benchmarkStatistics) report() {
	fmt.Println("====================== RESULTS ======================")
	fmt.Printf("● total # of queries processed: %v\n",
		stat.numberOfQueriesRun)
	fmt.Printf("● total processing time across all queries: %.2fms\n",
		float64(stat.totalProcessingTimeForAllQueries.Nanoseconds())/1e6)
	fmt.Printf("● the minimum query time: %.2fms\n",
		float64(stat.minimumQueryTime.Nanoseconds())/1e6)
	fmt.Printf("● the median query time: %.2fms\n",
		float64(stat.medianQueryTime.Nanoseconds())/1e6)
	fmt.Printf("● the average query time: %.2fms\n",
		stat.averageQueryTime)
	fmt.Printf("● the maximum query time: %.2fms\n",
		float64(stat.maximumQueryTime.Nanoseconds())/1e6)
	fmt.Println("=====================================================")
}
