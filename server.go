package main

/**
	Server is the middle man between web/mobile to database
*/
import (
	"net/http"
)


func main() {

	http.ListenAndServe(":8080", nil)

}