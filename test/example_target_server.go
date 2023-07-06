package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var listOfEmailsAndCC = `This is the list of emails ["jorge@gmail.com", "pepe_232@hotmail.com", "fake@yahoo.com.ar"] and the admin mail is admin@argentina.gov   this is a credit card 5105105105105100 4012888888881881 this is no cc 1234-5678-9012-3456 this is 3530 1113 3330 0000 4012-8888-8888-1881 40128888888818814012888888881881`

func main() {
	srv := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Printf("[target server] received %s request at: %s\n", req.Method, time.Now())
		// reading headers
		for name, values := range req.Header {
			// Loop over all values for the name.
			for _, value := range values {
				fmt.Printf("[target server] header key: %s value: %s \n", name, value)
			}
		}
		_, _ = fmt.Fprint(rw, listOfEmailsAndCC)
	})

	log.Fatal(http.ListenAndServe(":8080", srv))
}
