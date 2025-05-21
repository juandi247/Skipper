package utils

import (
	"fmt"
	"net/http"
)

func Ping(url string, c *http.Client) (int, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("error en request", err)
		return 0, err
	}
	resp.Body.Close()
	fmt.Println("estatus code", resp.StatusCode)
	return resp.StatusCode, nil
}
