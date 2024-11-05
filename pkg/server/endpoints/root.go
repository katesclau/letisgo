package endpoints

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var Root = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		ret := "Received request: \n"

		// Echo headers
		for key, vals := range r.Header {
			for _, val := range vals {
				ret = fmt.Sprintf("%sHeader: %s, Value: %s\n", ret, key, val)
			}
		}

		// Echo method
		ret = fmt.Sprintf("%sMethod: %s\n", ret, r.Method)

		// Echo params
		values := r.URL.Query()
		for key, vals := range values {
			for _, val := range vals {
				ret = fmt.Sprintf("%sQuery param: %s, Value: %s\n", ret, key, val)
			}
		}

		// Echo body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			ret = fmt.Sprintf("%sFailed to read request body: %v\n", ret, err)
		} else {
			ret = fmt.Sprintf("%sRequest body: %s\n", ret, string(body))
		}

		time.Sleep(1 * time.Second)

		io.WriteString(w, ret)
	})
