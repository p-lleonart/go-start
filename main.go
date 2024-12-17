package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Todo struct {
	Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// converts a todo to JSON
func (t *Todo) WriteJson(w io.Writer) error {
	asJson, err := json.Marshal(*t)

	if err != nil {
		return err
	}

	_, err = w.Write(asJson)

	return err
}

// converts an array of todos to JSON
func WriteJson(todos []Todo, w io.Writer) error {
	asJson, err := json.Marshal(todos)

	if err != nil {
		return err
	}

	_, err = w.Write(asJson)

	return err
}

func fetchTodos(link string) ([]Todo, error) {
	var todos []Todo
	res, err := http.Get(link)

	if err != nil {
		return todos, fmt.Errorf("error: %w", err)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&todos)

	if err != nil {
		return todos, fmt.Errorf("error: %w", err)
	}

	return todos, nil
}

func HttpViewHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	todos, err := fetchTodos("https://jsonplaceholder.typicode.com/todos/?_limit=3")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = WriteJson(todos, w)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

}

func main() {
	http.HandleFunc("/", HttpViewHandler)
	http.ListenAndServe(":8000", nil)
}

// ch := make(chan []Todo)
// 	errCh := make(chan error)

// 	go func() {
// 		todos, err := fetchTodos("https://jsonplaceholder.typicode.com/todos/?_limit=3")

// 		if err != nil {
// 			errCh <- err
// 		} else {
// 			ch <- todos
// 		}
// 	}()

// 	defer close(ch)
// 	defer close(errCh)

// 	fmt.Println("Welcome")

// 	select {
// 	case err := <-errCh:
// 		panic(err)
// 	case todos := <-ch:
// 		fmt.Printf("Results: %#v", todos)
// 	}
