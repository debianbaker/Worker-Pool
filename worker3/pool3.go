package main

import (
	"fmt"
	"sync"
	"time"
	"net/http"
	"context"
	"io"
)

type Job struct{
	Id 	int
	Url string
}

type Result struct{
	JobId  	int
	Url 	string
	Latency  time.Duration
	StatusCode int
	Err 	error
}

func main(){
	jobs    := make(chan Job)
    results := make(chan Result)

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++{
		wg.Add(1)
		go worker(jobs, results, &wg, 3*time.Second)
	}

	go func(){
		urls := []string{
		"https://httpbin.org/status/200",
        "https://httpbin.org/status/404",
        "https://httpbin.org/status/500",
        "https://httpbin.org/status/500",
		"https://httpbin.org/delay/5",
        "https://httpbin.org/delay/1",
        "https://httpbin.org/delay/2",
        "https://httpbin.org/delay/3",
        "https://httpbin.org/delay/4",
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
		fmt.Printf("job-%d %s → %d, latency-%f\n", result.JobId, result.Url, result.StatusCode, float64(result.Latency)/float64(time.Second))
	}

	fmt.Println("Main ended")
}

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup, td time.Duration) {
    defer wg.Done()

    for job := range jobs {
        result := executeJob(job, td)
        results <- result
    }
    fmt.Println("G ended")
}

func executeJob(job Job, td time.Duration) Result {
    ctx, cancel := context.WithTimeout(context.Background(), td)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.Url, nil)
    if err != nil {
        return Result{JobId: job.Id, Url: job.Url, Err: err}
    }

    start := time.Now()
    resp, err := http.DefaultClient.Do(req)
    latency := time.Since(start)
    if err != nil {
        return Result{JobId: job.Id, Url: job.Url, Latency: latency, Err: err}
    }
    defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

    return Result{
        JobId:      job.Id,
        Url:        job.Url,
        Latency:    latency,
        StatusCode: resp.StatusCode,
    }
}