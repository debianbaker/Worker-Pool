/*
	Main is both the Producer and Consumer - First Produces jobs and then Consumes results.
	Using buffer, delays deadlock therefore - Tis Not a correction mechanism, only performance optimization
	Correction - Decoupling Prod and Cons
*/

package main

import (
	"fmt"
	"sync"
	"strings"
	"strconv"
)

type Job struct{
	Id 	int
	Url string
}

type Result struct{
	JobId  	int
	Url 	string
	StatusCode int
	Err 	error
}

func main(){
	/*
		Unbuffered channels - jobs and results
		jobs := make(chan Job)
	 	results := make(chan Result)
		No. of jobs = No. of workers o/w DEADLOCK!
	*/
	/*
		Unbuffered - Jobs
		Buffered - Results
		No. of max. Jobs that main can produce = results chan size + no. of workers
	*/
	/*
		Buffered - Jobs
		Unbuffered - Results
		No. of max. Jobs that main can produce = jobs chan size + no. of workers
	*/
	/*
		Buffered - Jobs
		Buffered - Results
		No. of max. Jobs that main can produce = jobs chan size + results chan size +  no. of workers
	*/

	// in this example, no of jobs produced max = 5+5+3 = 13 
	jobs    := make(chan Job, 5)
    results := make(chan Result, 5)

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++{
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	urls := [13]string{
		"https://httpbin.org/status/200",
        "https://httpbin.org/status/404",
        "https://httpbin.org/status/500",
        "https://httpbin.org/delay/1",
        "https://httpbin.org/status/200",
        "https://httpbin.org/status/200",
        "https://httpbin.org/status/404",
        "https://httpbin.org/status/500",
        "https://httpbin.org/delay/1",
        "https://httpbin.org/status/200",
        "https://httpbin.org/status/200",
        "https://httpbin.org/status/404",
        "https://httpbin.org/status/404",
	}
	for i, url := range urls{
		jobs <- Job{Id: i+1, Url: url}
	}
	close(jobs)

	go func(){
		wg.Wait()
		close(results)
	}()

	for result := range results{
		fmt.Printf("job-%d %s → %d\n", result.JobId, result.Url, result.StatusCode)
	}
	fmt.Println("Main ended")
}

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup){
	defer wg.Done()
    for job := range jobs {
        tp := strings.Split(job.Url, "/")
        statusCode, _ := strconv.Atoi(tp[len(tp)-1])
        results <- Result{
            JobId:      job.Id,
            Url:        job.Url,
            StatusCode: statusCode,
        }
    }
    fmt.Println("G ended")
}