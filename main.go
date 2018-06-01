package main

import (
	"net/http"
	"swarm-status/web"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", web.RootHandler)

	http.ListenAndServe(":8080", nil)
	//
	//var services, err = docker.ReadServiceList()
	//
	//if err != nil {
	//	fmt.Print(err)
	//} else {
	//	fmt.Print(time.Now().UTC(), "\n\n\n")
	//	for _, s := range services {
	//		fmt.Print(s, "\n")
	//	}
	//}
}
