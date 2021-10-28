package main

import (
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)



func Ping(site string) string {
	var client = &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	// 1. HTTP клиент по умолчанию устанавливает ограничения на количество одновременных запросов.
	// Это стоит учитывать когда делаешь канал на 10000 элементов и 70 000 горутин
	// 2. Для тебя это ограничение не срабатывает потому что ты каждый раз при запуске функци создаешь новый клиент,
	// а не переиспользуешь то что есть, что должно приводить к богатому выделению памяти и неэффективности работы приложения.
	if !strings.Contains(site, "https://") {
		site = "https://" + site
	}
	re := New(site, "", 0, "-")
	resp, err := client.Get(site)
	resp.Body.Close()
	if err != nil {
		re.statusCode = 500
		re.ip = site
		return re.ToString()
	}
	re.title = getTitle(resp)
	re.statusCode = resp.StatusCode
	re.ip = getIP(site)
	if re.statusCode == 302 {
		re.ip = resp.Header.Get("Location")
	}

	return re.ToString()
}

func getIP(site string) string {
	ip, err := net.ResolveIPAddr("ip4", site)
	if err != nil {
		defer recovery()
		panic(err)
	}
	return ip.String()
}

func getTitle(resp *http.Response) string {
	tkn := html.NewTokenizer(resp.Body)
	var isTitle bool
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return ""
		case tt == html.StartTagToken:
			t := tkn.Token()
			isTitle = t.Data == "title"
		case tt == html.TextToken:
			t := tkn.Token()
			if isTitle {
				return t.Data
			}
		}

	}
}
