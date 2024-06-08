package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Window struct {
	numbers []int
	size    int
	mu      sync.Mutex
}

type NumbersResponse struct {
	Numbers []int `json:"numbers"`
}

func NewWindow(size int) *Window {
	return &Window{
		numbers: make([]int, 0, size),
		size:    size,
	}
}

func (w *Window) Add(numbers []int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, number := range numbers {
		if len(w.numbers) == w.size {
			copy(w.numbers, w.numbers[1:])
			w.numbers[w.size-1] = number
		} else {
			w.numbers = append(w.numbers, number)
		}
	}
}

func (w *Window) Get() []int {
	w.mu.Lock()
	defer w.mu.Unlock()

	return append([]int(nil), w.numbers...)
}

func getNumbersFromAPI(url string) ([]int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzE3ODI0NTcxLCJpYXQiOjE3MTc4MjQyNzEsImlzcyI6IkFmZm9yZG1lZCIsImp0aSI6ImJlMDMyMTQ1LWJkNDQtNDRiNS1iYmI1LWNjNDA1ODhlNmE1MCIsInN1YiI6IjIxMzExQTY2MDJAc3JlZW5pZGhpLmVkdS5pbiJ9LCJjb21wYW55TmFtZSI6IklNNDUxNDVWIiwiY2xpZW50SUQiOiJiZTAzMjE0NS1iZDQ0LTQ0YjUtYmJiNS1jYzQwNTg4ZTZhNTAiLCJjbGllbnRTZWNyZXQiOiJXclJPS3BrdUNjYldYa0ZIIiwib3duZXJOYW1lIjoiQXNoaXNoIE1hbGxhIiwib3duZXJFbWFpbCI6IjIxMzExQTY2MDJAc3JlZW5pZGhpLmVkdS5pbiIsInJvbGxObyI6IjIxMzExQTY2MDIifQ.Ayc8IkmA6hRHUjNLRj4BY9GbuObuyMl4FLShD2UCUqM")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var numbersResponse NumbersResponse
	err = json.Unmarshal(body, &numbersResponse)
	if err != nil {
		return nil, err
	}

	return numbersResponse.Numbers, nil
}

func main() {
	window := NewWindow(10)

	router := gin.Default()

	router.GET("/numbers/:type", func(c *gin.Context) {
		typ := c.Param("type")

		var url string
		switch typ {
		case "e":
			url = "http://20.244.56.144/test/even"
		case "f":
			url = "http://20.244.56.144/test/fibonacci"
		case "p":
			url = "http://20.244.56.144/test/primes"
		case "r":
			url = "http://20.244.56.144/test/random"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
			return
		}

		numbers, err := getNumbersFromAPI(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		prevState := window.Get()
		window.Add(numbers)
		currState := window.Get()

		avg := 0
		if len(currState) > 0 {
			sum := 0
			for _, num := range currState {
				sum += num
			}
			avg = sum / len(currState)
		}

		c.JSON(http.StatusOK, gin.H{
			"numbers":         numbers,
			"windowPrevState": prevState,
			"windowCurrState": currState,
			"avg":             avg,
		})
	})

	router.Run(":9876")
}
