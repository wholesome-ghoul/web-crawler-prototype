package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/urfave/cli/v2"
	"github.com/wholesome-ghoul/web-crawler-prototype/config"
	contentParser "github.com/wholesome-ghoul/web-crawler-prototype/content-parser"
	"github.com/wholesome-ghoul/web-crawler-prototype/frontier"
	htmlDownloader "github.com/wholesome-ghoul/web-crawler-prototype/html-downloader"
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
				downloaderChannel := make(chan *frontier.PriorityQueue, numWorkers)
				parserChannel := make(chan contentParser.Content, totalUrls)

				fmt.Printf("Number of urls: %d\n", totalUrls)
				fmt.Printf("Number of jobs: %d\n", numWorkers)
				var wg sync.WaitGroup
				for w := 1; w <= numWorkers; w++ {
					wg.Add(1)
					w := w
					go func() {
						defer wg.Done()
						htmlDownloader.Download(w, downloaderChannel, parserChannel)
					}()

					for i := 0; i < urlsToDownload[w-1].Size(); i++ {
						wg.Add(1)
						go func() {
							defer wg.Done()
							contentParser.Parse(<-parserChannel)
						}()
					}

					downloaderChannel <- &urlsToDownload[w-1]
				}
				close(downloaderChannel)
				wg.Wait()

				return nil
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Couldn't run app. ", err)
	}
}
