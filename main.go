package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Event struct {
	Id   int
	Task string
}

func main() {
	ch := make(chan Event)

	http.HandleFunc("/consume", func(w http.ResponseWriter, r *http.Request) {
		for {
			select {
			case e := <-ch:
				fmt.Println("consuming", e)
				fmt.Fprintf(w, "%v\n", e)
			case <-time.After(1 * time.Second):
				w.WriteHeader(204)
				return
			}
		}
	})

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query()["id"][0])
		task := r.URL.Query()["task"][0]
		go func() {
			e := Event{Id: id, Task: task}
			fmt.Println(e)
			ch <- e
		}()
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
