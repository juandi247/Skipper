package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewHttpClient() {
	c := http.Client{
		Timeout: time.Duration(2) * time.Second,
	}

	resp, err := c.Get("http://127.0.0.1:5000")

	if err != nil {
		fmt.Println("error", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}
