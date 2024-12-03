package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	})
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			_, err := fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
		}
	})
	err := r.Run(":9999")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
}
