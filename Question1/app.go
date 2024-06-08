package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	windowSize  int
	window      []int
	windowMutex *sync.Mutex
}

func NewApp(windowSize int) *App {
	return &App{
		windowSize:  windowSize,
		window:      make([]int, 0, windowSize),
		windowMutex: &sync.Mutex{},
	}
}

func (a *App) fetchNumbers() ([]int, error) {
	// Simulate fetching numbers from the test server
	time.Sleep(time.Second) // Simulate some latency
	return []int{1, 3, 5}, nil

	// Replace with actual logic to fetch numbers from test server API
}

func (a *App) updateWindow(numbers []int) {
	a.windowMutex.Lock()
	defer a.windowMutex.Unlock()

	if len(a.window)+len(numbers) > a.windowSize {
		a.window = a.window[len(numbers):] // Remove the oldest elements if exceeding window size
	}
	a.window = append(a.window, numbers...)
}

func (a *App) calculateAverage() float64 {
	a.windowMutex.Lock()
	defer a.windowMutex.Unlock()

	var sum float64
	for _, num := range a.window {
		sum += float64(num)
	}
	return sum / float64(len(a.window))
}

func (a *App) handleRequest(c *gin.Context) {
	numbers, err := a.fetchNumbers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	a.updateWindow(numbers)
	avg := a.calculateAverage()

	c.JSON(http.StatusOK, gin.H{
		"numbers":         numbers,
		"windowPrevState": a.window[:len(a.window)-len(numbers)],
		"windowCurrState": a.window,
		"avg":             avg,
	})
}
