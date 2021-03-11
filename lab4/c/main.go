package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var counter int
var m sync.RWMutex

func randName(n int) string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphabet[b%byte(len(alphabet))]
	}
	return string(bytes)
}

func randNumber(n int) string {
	const alphabet = "0123456789"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphabet[b%byte(len(alphabet))]
	}
	return string(bytes)
}

func reader(file *os.File) {
	buffer := make([]byte, 42)

	for {
		m.RLock()

		if counter <= 1 {
			m.RUnlock()
			time.Sleep(time.Second)
			continue
		}

		pos := rand.Intn(counter - 1)
		file.Seek(int64(pos*42), 0)
		_, err := file.Read(buffer)

		if err != nil {
			// log.Fatal(counter, pos)
		}

		m.RUnlock()

		fmt.Println(string(buffer))
		time.Sleep(time.Second)
	}
}

func writer(file *os.File) {
	for {

		m.Lock()

		name := randName(20)
		number := randNumber(20)
		data := []byte(fmt.Sprintf("%s %s\n", name, number))

		task := rand.Intn(3)
		if counter < 1 {
			task = 1
		}

		switch task {
		case 0:
			file.WriteAt(data, int64(rand.Intn(counter)*42))
			fmt.Println("Modified")
		case 1:
			file.WriteAt(data, int64(counter)*42)
			counter++
			fmt.Println("Added")
		case 2:
			pos := rand.Intn(counter)
			for i := pos + 1; i < counter; i++ {
				var buffer = make([]byte, 42)
				file.Seek(int64(i*42), 0)
				file.Read(buffer)
				file.WriteAt(buffer, int64(i-1)*42)
			}
			counter--
			file.Truncate(int64(counter * 42))
			fmt.Println("Deleted")
		}

		m.Unlock()

		time.Sleep(time.Second)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	file, _ := os.OpenFile("database.txt", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer file.Close()

	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	go reader(file)
	// go reader(file)

	go writer(file)
	go writer(file)
	go writer(file)

	wg.Wait()
}
