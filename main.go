package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

var worklst = make(chan string)

func main() {
	var httpTimeout int
	var funcNumber int
	var filename string

	flag.IntVar(&httpTimeout, "timeout", 6, "Timeout in seconds for connection")
	flag.IntVar(&funcNumber, "funcs", 20, "Number of goroutines")
	flag.StringVar(&filename, "file", " ", "Filename")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(funcNumber)

	client := http.Client{
		Timeout: time.Duration(httpTimeout) * time.Second,
	}

	go readLines(filename)

	for x := 0; x < funcNumber; x++ {
		go func() {
			defer wg.Done()
			for v := range worklst {
				test(client, v)
			}

		}()

	}
	wg.Wait()

}

func test(client http.Client, v string) {
	resp, err := client.Get("http://" + v)
	if err != nil {

		fmt.Println("error ", v)

	} else {

		if resp.StatusCode == 200 && resp.Request.URL.String() == "http://test.com" {
			fmt.Println("OK ", v)
		} else if resp.StatusCode == 200 && resp.Request.URL.String() != "http://test.com" {
			fmt.Println(v, " ", resp.Request.URL)
		} else {
			fmt.Println(resp.Request.URL, " ", resp.StatusCode)
		}
		resp.Body.Close()
	}

}

func readLines(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		worklst <- scanner.Text()
	}
	close(worklst)
}
