/*
	Decoupled Prod and Cons
	Main is result Consumer, GoRoutine spawned as jobs Producer.
	** Why not main producer? Because after producing main will end without waiting for consumer to finish.
	Unbuffered Channels - Jobs and Results
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
	jobs    := make(chan Job)
    results := make(chan Result)

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++{
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	go func(){
		urls := [15]string{
		"https://httpbin.org/status/200",
        "https://httpbin.org/status/404",
        "https://httpbin.org/status/500",
        "https://httpbin.org/status/500",
        "https://httpbin.org/delay/1",
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
	}()
	
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