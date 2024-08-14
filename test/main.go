package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

var wg sync.WaitGroup

type StressTest struct {
	Url         string
	Requests    int
	Concurrency int
	Token       string
	Begin       time.Time
	Final       time.Time
	Duration    time.Duration
}

type Results struct {
	TotalRequests   int
	SuccessRequests int
	StatusCodes     map[int]int
}

func (r Results) GenerateReport(duration time.Duration) {
	fmt.Println()
	red := color.New(color.FgRed)
	red.Println("TEST RESULT")
	fmt.Printf("Test duration: %.2fs\n", duration.Seconds())
	fmt.Printf("Successful Requests: %d\n", r.SuccessRequests)
	fmt.Println()
	headerFmt := color.New(color.FgGreen).SprintfFunc()
	tbl2 := table.New("Status", "Total")
	tbl2.WithHeaderFormatter(headerFmt)
	for code, count := range r.StatusCodes {
		tbl2.AddRow(code, count)
	}
	tbl2.Print()
	fmt.Println()
}

func (t *StressTest) testing(token string) Results {
	m := sync.Mutex{}
	result := Results{
		TotalRequests:   t.Requests,
		SuccessRequests: 0,
		StatusCodes:     make(map[int]int),
	}
	concurrency := make(chan struct{}, t.Concurrency)
	for i := 0; i < t.Requests; i++ {
		wg.Add(1)
		concurrency <- struct{}{}
		go func() {
			defer wg.Done()
			code, err := HttpRequest(t.Url, token)
			m.Lock()
			result.StatusCodes[code]++
			if err == nil {
				result.SuccessRequests++
			}
			m.Unlock()
			<-concurrency
		}()
	}
	wg.Wait()
	t.Final = time.Now()
	t.Duration = t.Final.Sub(t.Begin)
	return result
}

func (t *StressTest) Execute() {
	blue := color.New(color.FgHiBlue)
	yellow := color.New(color.FgHiYellow)
	t.Begin = time.Now()
	fmt.Println()
	blue.Println("-- RUN TEST -->")
	fmt.Printf("Url: %s\n", t.Url)
	fmt.Printf("Requests: %d\n", t.Requests)
	// fmt.Printf("Concurrency: %d\n", t.Concurrency)
	fmt.Println("")

	// TEST IP LIMITATION
	t.Begin = time.Now()
	blue.Println("IP LIMITATION: First Testing")
	fmt.Println(" - default settings: 5 RPS and 20 seconds of block")
	fmt.Println(" - expected result: status 200 = 5 and 429 = 15")
	result := t.testing("")
	result.GenerateReport(t.Duration)
	yellow.Println("Wait 10 seconds of blocking and retests")
	time.Sleep(10 * time.Second)
	fmt.Println("")

	t.Begin = time.Now()
	blue.Println("IP LIMITATION: Second Testing (blocking time)")
	fmt.Println(" - default settings: 5 RPS and 20 seconds of block")
	fmt.Println(" - expected result: status 429 = 20")
	result = t.testing("")
	result.GenerateReport(t.Duration)
	yellow.Println("Wait 10 seconds of blocking and retests")
	time.Sleep(10 * time.Second)
	fmt.Println("")

	t.Begin = time.Now()
	blue.Println("IP LIMITATION: Third Testing (block released)")
	fmt.Println(" - default settings: 5 RPS and 20 seconds of block")
	fmt.Println(" - expected result: status 200 = 5 and 429 = 15")
	result = t.testing("")
	result.GenerateReport(t.Duration)
	yellow.Println("Wait 10 seconds of new test")
	time.Sleep(10 * time.Second)
	fmt.Println("")

	// TEST TOKEN LIMITATION
	t.Begin = time.Now()
	blue.Println("TOKEN LIMITATION: First Testing")
	fmt.Println(" - default settings: token bxp2sl28mv78p => 10 RPS and 10 seconds of block")
	fmt.Println(" - expected result: status 200 = 10 and 429 = 10")
	result = t.testing(t.Token)
	result.GenerateReport(t.Duration)
	yellow.Println("Wait 7 seconds of blocking and retests")
	time.Sleep(7 * time.Second)
	fmt.Println("")

	t.Begin = time.Now()
	blue.Println("TOKEN LIMITATION: Second Testing (blocking time)")
	fmt.Println(" - default settings: token bxp2sl28mv78p => 10 RPS and 10 seconds of block")
	fmt.Println(" - expected result: status 429 = 20")
	result = t.testing(t.Token)
	result.GenerateReport(t.Duration)
	yellow.Println("Wait 3 seconds of blocking and retests")
	time.Sleep(3 * time.Second)
	fmt.Println("")

	t.Begin = time.Now()
	blue.Println("TOKEN LIMITATION: Third Testing (block released)")
	fmt.Println(" - default settings: token bxp2sl28mv78p => 10 RPS and 10 seconds of block")
	fmt.Println(" - expected result: status 200 = 5 and 429 = 15")
	result = t.testing(t.Token)
	result.GenerateReport(t.Duration)
	yellow.Println("-- FINISH TESTS -->")
	fmt.Println("")

}

func HttpRequest(url string, token string) (int, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	if token != "" {
		req.Header.Set("API_KEY", token)
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
		return resp.StatusCode, err
	}
	return 0, err
}

func main() {
	test := StressTest{
		Url:         "http://localhost:8080",
		Requests:    20,
		Concurrency: 1,
		Token:       "bxp2sl28mv78p",
	}
	test.Execute()
}
