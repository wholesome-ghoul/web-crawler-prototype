package html_downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	contentParser "github.com/wholesome-ghoul/web-crawler-prototype/content-parser"
	logger "github.com/wholesome-ghoul/web-crawler-prototype/custom-logger"
	"github.com/wholesome-ghoul/web-crawler-prototype/frontier"
)

const WAIT_PER_REQUEST = time.Second
const TIMEOUT = 5 * time.Second
const ROOT_DATA_DIR = "html-data"

func sanitize(name string) string {
	invalidChars := "/\\:*?\"<>|"
	for _, c := range invalidChars {
		name = strings.ReplaceAll(name, string(c), "_")
	}
	return name
}

func Download(id int,
	jobs <-chan *frontier.PriorityQueue,
	parserChannel chan<- contentParser.Content,
) {
	for job := range jobs {
		if job.Empty() {
			parserChannel <- contentParser.Content{}
			return
		}

		curr := job.Pop()
		client := &http.Client{
			Timeout: TIMEOUT,
		}

		hostnameDir := sanitize(curr.Hostname())
		rootDir := path.Join(ROOT_DATA_DIR, hostnameDir)
		os.MkdirAll(rootDir, 0777)

		for curr != nil {
			url := curr.Url()
			filename := path.Join(rootDir, sanitize(url)) + ".html"
			request, _ := http.NewRequest("GET", url, nil)

			response, err := client.Do(request)
			// logger.Log().Printf("WORKER %d started fetching (priority: %d) %s\n", id, curr.Priority(), url)
			if err != nil {
				fmt.Println("something went wrong. ", err)
				parserChannel <- contentParser.Content{}
				return
			}

			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("could not read body", err)
			}
			defer response.Body.Close()

			logger.Log().Printf("WORKER %d finished fetching %s status: %d\n", id, url, response.StatusCode)
			if response.StatusCode != 200 {
				parserChannel <- contentParser.Content{}
				return
			}

			err = os.WriteFile(filename, responseBody, 0666)
			if err != nil {
				fmt.Println("could not write to file", filename, "reason: ", err)
			}

			parserChannel <- contentParser.Content{Url: url, Path: filename}

			time.Sleep(WAIT_PER_REQUEST)
			curr = job.Pop()
		}
	}
}

func _Download(urls []frontier.PriorityQueue, jobs chan *frontier.PriorityQueue, wg *sync.WaitGroup) {
	numJobs := len(urls)
	results := make(chan int, numJobs)

	numWorkers := numJobs
	fmt.Printf("Number of jobs: %d\n", numJobs)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		// go worker(w, jobs, results, wg)
	}

	// priorities may not be distributed evenly, thus, some of the workers will
	// do more work than the others
	for j := 1; j <= numJobs; j++ {
		jobs <- &urls[j-1]
	}
	close(jobs)

	for r := 1; r <= numJobs+3; r++ {
		fmt.Println(<-results)
	}

	wg.Wait()
}
