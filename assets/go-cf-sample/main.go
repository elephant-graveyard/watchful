// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	go spamLog(time.Second)

	fmt.Println("Starting watchful sample app for curl merkhets!")
	http.HandleFunc("/", renderIndexPage)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func renderIndexPage(out http.ResponseWriter, in *http.Request) {
	out.Header().Add("content-type", "application/json")
	out.WriteHeader(200)

	_, _ = fmt.Fprint(out, `{"totally-random-number-without-any-meaning":949207500}`)
}

func spamLog(iteration time.Duration) {
	ticker := time.NewTicker(iteration)
	for {
		select {
		case t := <-ticker.C:
			fmt.Println(t.Unix())
		}
	}
}
