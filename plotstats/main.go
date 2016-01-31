package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var count, freq int

	file := flag.String("infile", "input.stats", "Input file")
	outfile := flag.String("outfile", "stats.svg", "Output file")
	statKeys := flag.String("keys", "", "Stats keys")
	flag.Parse()

	p := NewPlot("stats")
	keys := strings.Split(*statKeys, ",")
	stats := make([]*Line, len(keys))
	for i, key := range keys {
		stats[i] = p.NewLine(key)
	}

	fd, err := os.Open(*file)
	defer fd.Close()

	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(fd)
	line, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	fmt.Sscanf(line, "Count:%d Frequency:%d", &count, &freq)
	for i := 0; i < count; i++ {
		var l int
		var v map[string]int64
		line, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}

		fmt.Sscanf(line, "Len:%d", &l)
		buf := make([]byte, l)
		l, err = io.ReadFull(r, buf)
		if err != nil {
			panic(err)
		}

		json.Unmarshal(buf, &v)
		for x, stat := range stats {
			stat.AddPoint(float64(i*freq), float64(v[keys[x]]))
		}
	}

	p.Write(*outfile)
}
