// // package main

// // import (
// // 	"fmt"
// // 	"math/rand"
// // 	"net/http"
// // 	"time"

// // 	"github.com/gin-gonic/gin"
// // )

// // // Number represents a number received from the test server
// // type Number struct {
// // 	Value float64 `json:"value"`
// // }

// // // Window represents the window of numbers stored in the microservice
// // type Window struct {
// // 	maxSize int
// // 	data    []Number
// // }

// // // NewWindow creates a new Window with a specified size
// // func NewWindow(size int) *Window {
// // 	return &Window{
// // 		maxSize: size,
// // 		data:    make([]Number, 0, size),
// // 	}
// // }

// // // AddNumber adds a number to the window, maintaining its size
// // func (w *Window) AddNumber(number Number) {
// // 	if len(w.data) == w.maxSize {
// // 		w.data = append(w.data[1:], number)
// // 	} else {
// // 		w.data = append(w.data, number)
// // 	}
// // }

// // // GetAverage calculates the average of the numbers in the window
// // func (w *Window) GetAverage() float64 {
// // 	total := 0.0
// // 	for _, num := range w.data {
// // 		total += num.Value
// // 	}
// // 	if len(w.data) == 0 {
// // 		return 0.0
// // 	}
// // 	return total / float64(len(w.data))
// // }

// // // Microservice represents the average calculator microservice
// // type Microservice struct {
// // 	window  *Window
// // 	timeout time.Duration
// // }

// // // NewMicroservice creates a new Microservice instance
// // func NewMicroservice(windowSize int, timeout time.Duration) *Microservice {
// // 	return &Microservice{
// // 		window:  NewWindow(windowSize),
// // 		timeout: timeout,
// // 	}
// // }

// // // fetchNumbers simulates fetching numbers from a test server based on a qualifier
// // func (m *Microservice) fetchNumbers(qualifier string) ([]Number, error) {
// // 	// Replace this with actual logic to fetch numbers from your test server
// // 	// based on the qualifier (e.g., even, prime, fibonacci, random)
// // 	// Simulate a delay and potential errors
// // 	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
// // 	if rand.Intn(10) == 0 {
// // 		return nil, fmt.Errorf("error fetching numbers")
// // 	}
// // 	numbers := make([]Number, 0)
// // 	for i := 0; i < 10; i++ {
// // 		numbers = append(numbers, Number{Value: float64(i * 2)})
// // 	}
// // 	return numbers, nil
// // }

// // // handleAverage calculates and returns the average of the numbers in the window
// // func (m *Microservice) handleAverage(c *gin.Context) {
// // 	qualifier := c.Param("numberid")

// // 	numbers, err := m.fetchNumbers(qualifier)
// // 	if err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// // 		return
// // 	}

// // 	for _, num := range numbers {
// // 		m.window.AddNumber(num)
// // 	}

// // 	response := gin.H{
// // 		"numbersFetched": numbers,
// // 		"previousWindow": m.window.data[:len(m.window.data)-len(numbers)],
// // 		"currentWindow":  m.window.data,
// // 		"average":        m.window.GetAverage(),
// // 	}
// // 	c.JSON(http.StatusOK, response)
// // }

// // func main() {
// // 	router := gin.Default()

// // 	microservice := NewMicroservice(10, 500*time.Millisecond)

// // 	router.GET("/numbers/:numberid", microservice.handleAverage)

// // 	router.Run(":9876")
// // }
// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"sync"
// 	"time"
// )

// // Number represents a number received from the test server
// type Number struct {
// 	Value float64 `json:"value"`
// }

// // ConcurrentList is a thread-safe linked list for storing numbers
// type ConcurrentList struct {
// 	mutex sync.Mutex
// 	head  *Node
// 	tail  *Node
// }

// // Node represents a node in the linked list
// type Node struct {
// 	Value Number
// 	next  *Node
// }

// // NewConcurrentList creates a new ConcurrentList
// func NewConcurrentList() *ConcurrentList {
// 	return &ConcurrentList{
// 		mutex: sync.Mutex{},
// 	}
// }

// // PushBack adds a new element to the back of the list
// func (l *ConcurrentList) PushBack(value Number) {
// 	l.mutex.Lock()
// 	defer l.mutex.Unlock()

// 	newNode := &Node{Value: value}
// 	if l.head == nil {
// 		l.head = newNode
// 		l.tail = newNode
// 		return
// 	}
// 	l.tail.next = newNode
// 	l.tail = newNode
// }

// // PopFront removes and returns the first element from the list
// func (l *ConcurrentList) PopFront() (*Number, bool) {
// 	l.mutex.Lock()
// 	defer l.mutex.Unlock()

// 	if l.head == nil {
// 		return nil, false
// 	}
// 	value := l.head.Value
// 	l.head = l.head.next
// 	if l.head == nil {
// 		l.tail = nil
// 	}
// 	return &value, true
// }

// // Len returns the number of elements in the list
// func (l *ConcurrentList) Len() int {
// 	l.mutex.Lock()
// 	defer l.mutex.Unlock()

// 	count := 0
// 	node := l.head
// 	for node != nil {
// 		count++
// 		node = node.next
// 	}
// 	return count
// }

// // ToSlice creates a copy of the entire list as a slice
// func (l *ConcurrentList) ToSlice() []Number {
// 	l.mutex.Lock()
// 	defer l.mutex.Unlock()

// 	slice := make([]Number, 0)
// 	node := l.head
// 	for node != nil {
// 		slice = append(slice, node.Value)
// 		node = node.next
// 	}
// 	return slice
// }

// // Microservice represents the average calculator microservice
// type Microservice struct {
// 	windowSize int
// 	timeout    time.Duration
// 	window     *ConcurrentList
// }

// // NewMicroservice creates a new Microservice instance
// func NewMicroservice(windowSize int, timeout time.Duration) *Microservice {
// 	return &Microservice{
// 		windowSize: windowSize,
// 		timeout:    timeout,
// 		window:     NewConcurrentList(),
// 	}
// }

// // fetchNumbers simulates fetching numbers from a test server based on a qualifier
// func (m *Microservice) fetchNumbers(qualifier string) ([]Number, error) {
// 	// Replace this with actual logic to fetch numbers from your test server
// 	// based on the qualifier (e.g., even, prime, fibonacci, random)
// 	// Simulate a delay and potential errors
// 	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
// 	if rand.Intn(10) == 0 {
// 		return nil, fmt.Errorf("error fetching numbers")
// 	}
// 	numbers := make([]Number, 0)
// 	for i := 0; i < 10; i++ {
// 		numbers = append(numbers, Number{Value: float64(i * 2)})
// 	}
// 	return numbers, nil
// }

// // handleAverage calculates and returns the average of the numbers in the window
// func (m *Microservice) handleAverage(c *gin.Context) {
// 	qualifier := c.Param("numberid")

// 	numbers, err := m.fetchNumbers(qualifier)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

//	for _, num := range numbers {
//		m.window.PushBack(num)
//		if m.window.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Number represents a number received from the API
type Number struct {
	Value float64 `json:"value"`
}
type NumberResponse struct {
	Numbers []float64 `json:"numbers"` // Adjust the field names based on their API response format
}

// ConcurrentList is a thread-safe linked list for storing numbers
type ConcurrentList struct {
	mutex sync.Mutex
	head  *Node
	tail  *Node
}

// Node represents a node in the linked list
type Node struct {
	Value Number
	next  *Node
}

// NewConcurrentList creates a new ConcurrentList
func NewConcurrentList() *ConcurrentList {
	return &ConcurrentList{
		mutex: sync.Mutex{},
	}
}

// PushBack adds a new element to the back of the list
func (l *ConcurrentList) PushBack(value Number) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	newNode := &Node{Value: value}
	if l.head == nil {
		l.head = newNode
		l.tail = newNode
		return
	}
	l.tail.next = newNode
	l.tail = newNode
}

// PopFront removes and returns the first element from the list
func (l *ConcurrentList) PopFront() (*Number, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.head == nil {
		return nil, false
	}
	value := l.head.Value
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	}
	return &value, true
}

// Len returns the number of elements in the list
func (l *ConcurrentList) Len() int {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	count := 0
	node := l.head
	for node != nil {
		count++
		node = node.next
	}
	return count
}

// ToSlice creates a copy of the entire list as a slice
func (l *ConcurrentList) ToSlice() []Number {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	slice := make([]Number, 0)
	node := l.head
	for node != nil {
		slice = append(slice, node.Value)
		node = node.next
	}
	return slice
}

// Microservice represents the average calculator microservice
type Microservice struct {
	windowSize int
	timeout    time.Duration
	window     *ConcurrentList
	apiUrl     string // Store the API URL for retrieving numbers
}

// NewMicroservice creates a new Microservice instance
func NewMicroservice(windowSize int, timeout time.Duration, apiUrl string) *Microservice {
	return &Microservice{
		windowSize: windowSize,
		timeout:    timeout,
		window:     NewConcurrentList(),
		apiUrl:     apiUrl,
	}
}

// fetchNumbers retrieves numbers from their API based on qualifier (replace with actual logic)
func (m *Microservice) fetchNumbers(qualifier string) ([]Number, error) {
	// Replace this with actual logic to call their API based on qualifier
	apiUrl := m.apiUrl + "/" + qualifier // Construct the API URL with qualifier

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Retrieve bearer token securely (e.g., from environment variables)
	bearerToken := "<your_bearer_token>"
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var numberResponse NumberResponse // Adjust the struct based on their API response format
	err = json.Unmarshal(body, &numberResponse)
	if err != nil {
		return nil, err
	}

	return numberResponse.Numbers, nil // Assuming the response contains an array of numbers
}

// handleAverage calculates and returns the average of the numbers in the window
func (m *Microservice) handleAverage(c *gin.Context) {
	qualifier := c.Param("numberid")

	numbers, err := m.fetchNumbers(qualifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, num := range numbers {
		m.window.PushBack(num)
		if m.window.Len() > m.windowSize {
			_, _ = m.window.PopFront() // Remove the oldest element while keeping the lock
		}
	}

	previousWindow := make([]Number, 0)
	node := m.window.head
	for node != nil {
		previousWindow = append(previousWindow, node.Value)
		node = node.next
	}

	response := gin.H{
		"numbersFetched": numbers,
		"previousWindow": previousWindow,
		"currentWindow":  m.window.ToSlice(), // Get a copy of the window for response
		"average":        m.calculateAverage(),
	}
	c.JSON(http.StatusOK, response)
}

// calculateAverage calculates the average of the numbers in the window
func (m *Microservice) calculateAverage() float64 {
	total := 0.0
	for _, num := range m.window.ToSlice() {
		total += num.Value
	}
	if m.window.Len() == 0 {
		return 0.0
	}
	return total / float64(m.window.Len())
}

func main() {
	router := gin.Default()

	// Replace these with actual values
	windowSize := 10
	timeout := 500 * time.Millisecond
	apiUrl := "http://<actual_api_url>" // Replace with their actual API URL

	microservice := NewMicroservice(windowSize, timeout, apiUrl)

	router.GET("/numbers/:numberid", microservice.handleAverage)

	router.Run(":9876")
}

