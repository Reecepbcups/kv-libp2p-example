package main

import (
	"sync"
)

// iterate over the table

func main() {
	// fmt.Println("Hello, World!")

	// store := NewStore()
	// table := store.Table("users")

	// table.Set("name", "John")
	// table.Set("age", fmt.Sprintf("%d", 30))

	// name, _ := table.Get("name")
	// fmt.Println(name)

	// age, _ := table.Get("age")
	// fmt.Println(age)

	s := NewServer()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(s *Server) {
		defer wg.Done()
		s.Stop()
	}(s)

	s.Start()
	wg.Wait()
}
