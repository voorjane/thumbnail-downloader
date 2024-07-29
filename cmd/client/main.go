package main

import (
	"flag"
	"getThumb/internal"
	"github.com/sirupsen/logrus"
)

func main() {
	async := flag.Bool("async", false, "Download thumbnails asynchronously")
	flag.Parse()

	urls := flag.Args()
	if len(urls) == 0 {
		logrus.Fatal("no URLs provided")
	}

	err := internal.RunClient(urls, *async)
	if err != nil {
		logrus.Fatal(err)
	}
}
