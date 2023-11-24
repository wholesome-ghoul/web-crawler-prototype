package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/wholesome-ghoul/web-crawler-prototype/config"
	"github.com/wholesome-ghoul/web-crawler-prototype/frontier"
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
			// TODO: locality
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

				urlFrontier.PrintAllBack()
				urlFrontier.Crawl()
				urlFrontier.Crawl()
				urlFrontier.Crawl()

				return nil
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Couldn't run app. ", err)
	}
}
