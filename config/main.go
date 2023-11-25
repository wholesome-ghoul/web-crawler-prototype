package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	SeedUrls []SeedUrl `json:"seed-urls"`
}

type SeedUrl struct {
	Category string `json:"category"`
	Url      string `json:"url"`
	Priority uint8  `json:"priority"`
	hostname string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) ParseFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	size := stat.Size()

	content := make([]byte, size)
	file.Read(content)

	err = json.Unmarshal(content, c)
	if err != nil {
		log.Fatal("Couldn't unmarshal config file. ", err)
	}
}

func (c *Config) Print(pretty bool) {
	if pretty {
		fmt.Println("Category\tURL\tPriority")
		fmt.Println("--------------------------------")
		for _, seedUrl := range c.SeedUrls {
			fmt.Printf("%s\t\t%12s\t%d\n", seedUrl.Category, seedUrl.Url, seedUrl.Priority)
		}

		return
	}

	fmt.Println(c)
}

func (s *SeedUrl) SetHostname(hostname string) {
	s.hostname = hostname
}

func (s *SeedUrl) Hostname() string {
	return s.hostname
}
