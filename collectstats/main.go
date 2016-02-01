package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	endpoint := flag.String("endpoint", "http://localhost:9102/stats", "Stats http endpoint")
	frequency := flag.Int("frequency", 1, "Sampling frequency in seconds")
	duration := flag.Int("seconds", 30, "Collection duration")
	file := flag.String("file", "out.stats", "Output file")

	flag.Parse()

	fd, err := os.OpenFile(*file, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	count := *duration / *frequency
	fd.Write([]byte(fmt.Sprintf("Count:%d Frequency:%d\n", count, *frequency)))
	for i := 0; i < count; i++ {
		resp, err := http.Get(*endpoint)
		if err != nil {
			panic(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		fd.Write([]byte(fmt.Sprintf("Len:%d\n", len(body)+1)))
		fd.Write(body)
		fd.Write([]byte("\n"))

		fmt.Print(".")
		time.Sleep(time.Second * time.Duration(*frequency))
	}

	fd.Close()
}
