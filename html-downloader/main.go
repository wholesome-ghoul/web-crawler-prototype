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

func worker(id int, jobs <-chan *frontier.PriorityQueue) {
	for job := range jobs {
		if job.Empty() {
			return
		}

		curr := job.Pop()
		client := &http.Client{}

		hostnameDir := sanitize(curr.Hostname())
		rootDir := path.Join(ROOT_DATA_DIR, hostnameDir)
		os.MkdirAll(rootDir, 0777)

		for curr != nil {
			url := curr.Url()
			filename := path.Join(rootDir, sanitize(url)) + ".html"
			request, _ := http.NewRequest("GET", url, nil)

			response, err := client.Do(request)
			logger.Log().Printf("WORKER %d started fetching (priority: %d) %s\n", id, curr.Priority(), url)
			if err != nil {
				fmt.Println("something went wrong. ", err)
			}

			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println("could not read body", err)
			}
			defer response.Body.Close()

			err = os.WriteFile(filename, responseBody, 0444)
			if err != nil {
				fmt.Println("could not write to file ", filename)
			}

			time.Sleep(WAIT_PER_REQUEST)
			logger.Log().Printf("WORKER %d finished fetching %s status: %d\n", id, url, response.StatusCode)

			curr = job.Pop()
		}
	}
}

func Download(urls []frontier.PriorityQueue) {
	numJobs := len(urls)
	jobs := make(chan *frontier.PriorityQueue, numJobs)
	var wg sync.WaitGroup

	numWorkers := numJobs
	fmt.Printf("Number of jobs: %d\n", numJobs)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		w := w

		go func() {
			defer wg.Done()
			worker(w, jobs)
		}()
	}

	// priorities may not be distributed evenly, thus, some of the workers will
	// do more work than the others
	for j := 1; j <= numJobs; j++ {
		jobs <- &urls[j-1]
	}
	close(jobs)

	wg.Wait()
}
