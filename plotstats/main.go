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
	useRate := flag.Bool("rate", false, "Plot the rate of growth")
	statKeys := flag.String("keys", "", "Stats keys")
	plotFreq := flag.Int("freq", 1, "Plot frequency")
	doSum := flag.Bool("sum", false, "Plot the sum of all stats")
	flag.Parse()

	p := NewPlot("stats")
	keys := strings.Split(*statKeys, ",")
	prevsY := make([]float64, len(keys))
	prevsX := make([]float64, len(keys))
	stats := make([]*Line, len(keys))
	for i, key := range keys {
		if *useRate {
			key = key + "/s"
		}

		if *doSum {
			if i == len(keys)-1 {
				stats[i] = p.NewLine(key)
			}

		} else {
			stats[i] = p.NewLine(key)
		}
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

		if i%*plotFreq == 0 {
			json.Unmarshal(buf, &v)
			var val int64

			for x, stat := range stats {
				val += v[keys[x]]
				if *doSum && x < len(stats)-1 {
					continue
				}

				xval := float64(i * freq)
				if *useRate {
					if i > 0 {
						rate := (float64(val) - prevsY[x]) / (xval - prevsX[x])
						stat.AddPoint(xval, rate)
					}
					prevsY[x] = float64(val)
					prevsX[x] = float64(xval)

				} else {
					stat.AddPoint(float64(i*freq), float64(val))
				}

				val = 0
			}
		}
	}

	p.Write(*outfile)
}
