package main

import "strconv"

type response struct { //дабы не создавать кучу лишних переменных и не возвращать одну и ту же строку, решил результирующие данные выделить в отдельную структуру
	url        string
	ip         string
	statusCode int
	title      string
}

func New(url string, ip string, statusCode int, title string) *response {
	return &response{
		url:        url,
		ip:         ip,
		statusCode: statusCode,
		title:      title,
	}
}

func (r *response) ToString() string {
	return r.url + " | " + r.ip + " | " + strconv.Itoa(r.statusCode) + " | " + r.title + "\n"
}
