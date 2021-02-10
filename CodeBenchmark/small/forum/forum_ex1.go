package main

import (
	"fmt"
	"strconv"
	"time"
)

type task struct {
	name      string
	path      string
	startTime time.Time
}

func downloader(id int, toDownload chan task, toProcess chan task) {
	for {
		task := <-toDownload
		fmt.Printf("[downloader_%v]: downloading %v \n", id, task.name)
		time.Sleep(1 * time.Second)
		toProcess <- task
	}
}

func transformer(id int, toProcess chan task, toFinish chan task) {
	for {
		task := <-toProcess
		fmt.Printf("[processor_%v]: processing %v \n", id, task.name)
		time.Sleep(1 * time.Second)
		toFinish <- task
	}
}

func finisher(id int, toFinish chan task, finished chan task) {
	for {
		task := <-toFinish
		fmt.Printf("[finisher_%v]: uploading %v \n", id, task.name)
		time.Sleep(1 * time.Second)
		finished <- task
	}
}

func fetchJobs(n int) *[]task {
	tasks := make([]task, n)
	for i := 0; i < n; i++ {
		name := "job_" + strconv.Itoa(i)
		tasks[i] = task{name: name}
	}
	return &tasks
}

func main() {
	start := time.Now()

	n := 10

	toDownload := make(chan task)
	toProcess := make(chan task)
	toFinish := make(chan task)
	finished := make(chan task)

	for i := 0; i < 4; i++ {
		go downloader(i, toDownload, toProcess)
	}

	for i := 0; i < 10; i++ {
		go transformer(i, toProcess, toFinish)
	}

	for i := 0; i < 4; i++ {
		go finisher(i, toFinish, finished)
	}

	jobs := fetchJobs(n)
	for _, j := range *jobs {
		toDownload <- j
	}

	for task := range finished {
		fmt.Printf("finished %v in %s \n", task.name, time.Since(task.startTime)/time.Second)
	}

	fmt.Printf("all tasks finished in %s", time.Since(start)/time.Second)
}
