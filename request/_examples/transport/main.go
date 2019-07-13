package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"go-sdk/request"
	"go-sdk/util"
)

func main() {
	//create external transport reference
	transport := &http.Transport{}

	// pass to the request
	req := request.New().AsGet().MustWithRawURL("https://google.com/robots.txt").WithTransport(transport)

	var meta *request.ResponseMeta
	var err error
	for x := 0; x < 10; x++ {
		// re-use it a whole bunch.
		meta, err = req.ExecuteWithMeta()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		} else {
			fmt.Fprintf(os.Stdout, util.JSON.SerializePretty(meta, "", "  "))
		}
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Done")
}
