package frontier

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/wholesome-ghoul/web-crawler-prototype/config"
)

const FRONT_PRIORITY_QUEUE_SIZE = 256

type URLFrontier struct {
	frontPriorityQueues []PriorityQueue
	backPriorityQueues  []PriorityQueue
	hostQueueMapping    map[string]int
}

func New() *URLFrontier {
	frontPriorityQueues := make([]PriorityQueue, FRONT_PRIORITY_QUEUE_SIZE)
	backPriorityQueues := []PriorityQueue{}
	hostQueueMapping := make(map[string]int)

	return &URLFrontier{
		frontPriorityQueues,
		backPriorityQueues,
		hostQueueMapping,
	}
}

func (u *URLFrontier) Prioritize(seedUrl config.SeedUrl) {
	// TODO: algorithm for calculating priority - PageRank, website traffic, update frequency, etc
	priorityIndex := seedUrl.Priority

	u.frontPriorityQueues[priorityIndex].Push(seedUrl)
}

func (u *URLFrontier) Print(pq PriorityQueue) {
	curr := pq.last

	for curr != nil {
		fmt.Println(curr.value)
		curr = curr.prev
	}
}

func (u *URLFrontier) PrintAllFront() {
	for i := 0; i < FRONT_PRIORITY_QUEUE_SIZE; i++ {
		if !u.frontPriorityQueues[i].Empty() {
			fmt.Printf("FRONT QUEUE #%d\n", i)
			u.Print(u.frontPriorityQueues[i])
			fmt.Println()
		}
	}
	fmt.Println("---")
}

func (u *URLFrontier) PrintAllBack() {
	for i := 0; i < len(u.backPriorityQueues); i++ {
		fmt.Printf("BACK QUEUE #%d\n", i)
		u.Print(u.backPriorityQueues[i])
		fmt.Println()
	}
	fmt.Println("---")
}

// randomly chooses a queue
func (u *URLFrontier) FrontQueueSelector() func() (PriorityQueue, error) {
	shuffledIndices := rand.Perm(FRONT_PRIORITY_QUEUE_SIZE)
	stoppedAt := 0
	for i := range shuffledIndices {
		stoppedAt = i
		randIndex := shuffledIndices[stoppedAt]
		if !u.frontPriorityQueues[randIndex].Empty() {
			// first non-empty front priority queue
			break
		}
	}

	return func() (PriorityQueue, error) {
		if stoppedAt >= len(u.frontPriorityQueues) {
			return PriorityQueue{}, errors.New("error: reached the end of the front queue")
		}

		stoppedAt++
		randIndex := shuffledIndices[stoppedAt-1]
		return u.frontPriorityQueues[randIndex], nil
	}
}

// ensures backQueue[i] only contains URLs from the same host
func (u *URLFrontier) BackQueueRouter(queue PriorityQueue) error {
	curr := queue.last
	for curr != nil {
		seedUrl := curr.value
		if !strings.HasPrefix(seedUrl.Url, "https://") {
			seedUrl.Url = "https://" + seedUrl.Url
		}

		parsedUrl, err := url.Parse(seedUrl.Url)
		if err != nil {
			fmt.Println("Could not parse url", err)
			return err
		}
		hostname := parsedUrl.Hostname()
		seedUrl.SetHostname(hostname)

		if index, ok := u.hostQueueMapping[hostname]; ok {
			u.backPriorityQueues[index].Push(seedUrl)
		} else {
			u.hostQueueMapping[hostname] = len(u.backPriorityQueues)
			pq := PriorityQueue{}
			pq.Push(seedUrl)
			u.backPriorityQueues = append(u.backPriorityQueues, pq)
		}

		curr = curr.prev
	}

	return nil
}

func (u *URLFrontier) UrlsToDownload() *[]PriorityQueue {
	return &u.backPriorityQueues
}
