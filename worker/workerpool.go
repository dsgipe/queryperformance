package worker

import (
	"sync"
)

func task(wg *sync.WaitGroup, w *worker) {
	for job := range w.job {
		usages, processingTime, err := queryUsages(job)
		if err != nil {
			w.err = err
		}
		*w.results <- statistics{
			workerId:       w.id,
			hostName:       job.host,
			processingTime: processingTime,
			usages:         usages,
		}
	}
	wg.Done()

}

func submitJobsToWorkers(parameters []Parameter, numberOfWorkers int) ([]statistics, error) {
	wg := sync.WaitGroup{}
	wg.Add(numberOfWorkers)
	results := make(chan statistics, len(parameters))

	workers := initializeWorkers(parameters, numberOfWorkers, results, &wg)

	var err error
	for i := range workers {
		close(workers[i].job)
		if workers[i].err != nil {
			err = workers[i].err
		}
	}

	wg.Wait()
	close(results)

	jobResults := make([]statistics, 0)
	for result := range results {
		jobResults = append(jobResults, result)
	}

	return jobResults, err
}

func initializeWorkers(parameters []Parameter, numberOfWorkers int, results chan statistics, wg *sync.WaitGroup) []worker {
	workers := make([]worker, numberOfWorkers)
	for i := range workers {
		workers[i] = worker{
			id:      i,
			job:     make(chan Parameter),
			results: &results,
		}
		go task(wg, &workers[i])
	}

	whichWorker := 0
	hostToWorker := map[string]*worker{}
	for _, job := range parameters {
		//Give worker the right job
		if _, ok := hostToWorker[job.host]; !ok {
			hostToWorker[job.host] = &workers[whichWorker]
			whichWorker++
			if whichWorker >= numberOfWorkers {
				whichWorker = 0
			}
		}
		hostToWorker[job.host].job <- job
	}
	return workers
}
