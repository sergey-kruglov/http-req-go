package main

import (
	"net/http"
)

const Url = "http://localhost:8080/api/public"
const Method = http.MethodGet
const Threads = 110
const PerThread = 1

var Headers = map[string]string{
	"Content-Type": "application/json",
}
var Body = map[string]string{}

func main() { run() }
