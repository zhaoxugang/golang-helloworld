package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main3() {
	var list []string
	list = []string{"http://www.baidu.com", "http://yar999.gitbook.io/gopl-zh/ch1/ch1-05"}
	for _, url := range list[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "getch: reading %s:%v\n", url, err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}
