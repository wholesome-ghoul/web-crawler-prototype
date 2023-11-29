package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/wholesome-ghoul/web-crawler-prototype/config"
	"github.com/wholesome-ghoul/web-crawler-prototype/frontier"
	html_downloader "github.com/wholesome-ghoul/web-crawler-prototype/html-downloader"
)

func main() {
	var configFilepath string
	var url cli.StringSlice

	app := &cli.App{
		Name:  "crawler-proto1",
		Usage: "Web crawler prototype 1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `FILE`",
				Destination: &configFilepath,
			},
			&cli.StringSliceFlag{
				Name:        "url",
				Aliases:     []string{"u"},
				Usage:       "Pass multiple urls per flag",
				Destination: &url,
			},
			&cli.BoolFlag{
				Name:    "locality",
				Aliases: []string{"l"},
				Usage:   "Crawl based on client's location",
			},
			// TODO: category; subcommand of --config?
			&cli.StringFlag{
				Name:  "category",
				Usage: "Crawl based on comma separetd categories",
			},
		},
		Action: func(ctx *cli.Context) error {
			if configFilepath != "" {
				c := config.NewConfig()
				c.ParseFile(configFilepath)
				// c.Print(true)

				urlFrontier := frontier.New()
				for _, seedUrl := range c.SeedUrls {
					urlFrontier.Prioritize(seedUrl)
				}
				// urlFrontier.PrintAllFront()

				frontQueueSelector := urlFrontier.FrontQueueSelector()
				for {
					frontQueue, err := frontQueueSelector()
					if err != nil {
						break
					}

					if frontQueue.Empty() {
						continue
					}

					urlFrontier.BackQueueRouter(frontQueue)
				}

				urlsToDownload := *urlFrontier.UrlsToDownload()
				totalUrls := urlFrontier.TotalUrls()
				numWorkers := len(urlsToDownload)
				var wg sync.WaitGroup
				downloaderChannel := make(chan *frontier.PriorityQueue, numWorkers)
				parserChannel := make(chan int, totalUrls)

				fmt.Printf("Number of urls: %d\n", totalUrls)
				fmt.Printf("Number of jobs: %d\n", numWorkers)
				for w := 1; w <= numWorkers; w++ {
					wg.Add(1)
					go html_downloader.Download(w, urlsToDownload, downloaderChannel, parserChannel, &wg)
				}

				for r := 1; r <= numWorkers; r++ {
					downloaderChannel <- &urlsToDownload[r-1]
				}

				for r := 0; r < totalUrls; r++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						fmt.Println(<-parserChannel)
					}()
				}

				close(downloaderChannel)
				wg.Wait()

				close(parserChannel)
				wg.Wait()

				// fileReader, _ := os.Open("./html-data/google.com/https___google.com.html")
				// content_parser.Parse(fileReader, results)

				return nil
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Couldn't run app. ", err)
	}
}
