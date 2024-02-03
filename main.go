package main

import (
	"es-job-search/connectors"
	"es-job-search/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	es, err := connectors.NewEsConnector()
	if err != nil {
		log.Fatal(err)
	}
	h := handlers.NewJobHandler(es)
	r.GET("/index-jobs", h.IndexJobsHandler)
	r.GET("/search-jobs", h.SearchJobsHandler)
	r.GET("/jobs-by-department", h.JobsByDepartmentHandler)
	r.Run()
}
