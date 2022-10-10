package worker

import "time"

type Parameter struct {
	host  string
	start time.Time
	end   time.Time
}

// cpuUsage is not actually needed
type cpuUsage struct {
	hostname string
	min      float64
	max      float64
	bucket   time.Time
}

type statistics struct {
	workerId       int
	hostName       string
	processingTime time.Duration
	usages         []cpuUsage
}

type worker struct {
	id      int
	job     chan Parameter
	results *chan statistics
	err     error
}

type benchmarkStatistics struct {
	numberOfQueriesRun               int
	totalProcessingTimeForAllQueries time.Duration
	minimumQueryTime                 time.Duration
	medianQueryTime                  time.Duration
	averageQueryTime                 float64
	maximumQueryTime                 time.Duration
	allTimes                         []time.Duration
}
