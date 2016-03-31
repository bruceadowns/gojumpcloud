package main

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	sleepDurationMs = 5000
	formName        = "password"
)

func usage(w http.ResponseWriter, err error) {
	htmlHeader := `<html>
<body>`
	htmlForm := `<form action="/" method="post">
  Enter Password:<br>
  <input type="text" name="password" value="angryMonkey">
  <input type="submit" value="Submit">
</form>`
	htmlFooter := `</body>
</html>`

	fmt.Fprintf(w, "%s", htmlHeader)

	if err != nil {
		fmt.Fprintf(w, "Error occurred: %v<br><br>", err)
	}

	fmt.Fprintf(w, "%s", htmlForm)
	fmt.Fprintf(w, "%s", htmlFooter)
}

func hashEncodePassword(p string) string {
	h := sha512.New()
	h.Write([]byte(p))
	sum := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(sum)
}

// rootHandler handles all requests
//
// hash password iif all conditions are met
// * this is an http post
// * the form contains a single element named "password"
// * the password element is not empty
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// looking for POST only
	if r.Method != "POST" {
		usage(w, fmt.Errorf("Expect POST, received %v", r.Method))
		return
	}

	// looking for / only
	if r.URL.Path != "/" {
		usage(w, fmt.Errorf("Expect /, received %v", r.URL.Path))
		return
	}

	// looking for form with single element password=%s only
	r.ParseForm()

	if len(r.Form) != 1 {
		usage(w, fmt.Errorf("Expect single password field. Received %v", r.Form))
		return
	}

	var password string
	if v, ok := r.Form[formName]; ok {
		if len(v) != 1 {
			usage(w, fmt.Errorf("Expect single password field. Received %v", r.Form[formName]))
			return
		}

		password = v[0]
	} else {
		usage(w, fmt.Errorf("Expect password form field"))
		return
	}

	d := time.Duration(sleepDurationMs) * time.Millisecond
	log.Printf("Sleep for %s", d)
	time.Sleep(d)

	log.Printf("Encode '%s'", password)
	hashedPassword := hashEncodePassword(password)
	fmt.Fprintf(w, "%s", hashedPassword)
	log.Printf("%s -> %s", password, hashedPassword)
}

func main() {
	http.HandleFunc("/", rootHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
