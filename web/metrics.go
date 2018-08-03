package web

import (
	"net/http"
	"fmt"
	"html/template"
	"time"
)

type MetricPageData struct {
	MetricPageSource string
	Timestamp 		 string
}

func MakeMetricHandler() http.HandlerFunc {
	settings := GetEnvSettings()

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/metrics.html")

		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, MetricPageData{
			MetricPageSource: settings.MetricPageSource,
			Timestamp:        time.Now().Format("Mon Jan _2 15:04:05"),
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}
