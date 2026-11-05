package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	usernamePtr := flag.String("github-activity", "", "github user name")

	flag.Parse()

	if usernamePtr == nil {
		fmt.Fprintln(os.Stderr, "flag github-activity is empty or missing")
		return
	}

	url := fmt.Sprintf("https://api.github.com/users/%s/events", *usernamePtr)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading response error: %v", err)
	}

	fmt.Println(string(data))
}
