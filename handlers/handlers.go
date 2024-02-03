package handlers

import (
	"es-job-search/connectors"
	"log"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	es connectors.IEsConnector
}

func NewJobHandler(esConnector connectors.IEsConnector) *JobHandler {
	return &JobHandler{es: esConnector}
}

const indexName = "jobs_search"

func (h *JobHandler) IndexJobsHandler(c *gin.Context) {

	err := h.es.CleanIndex(indexName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Println("Cleaned index:", indexName)

	err = h.es.IndexJobs(indexName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

func (h *JobHandler) SearchJobsHandler(c *gin.Context) {

	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(400, gin.H{"error": "Missing keyword query parameter"})
		return
	}

	jobs, err := h.es.SearchJobs(indexName, keyword)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"total": len(jobs), "jobs": jobs})
}

func (h *JobHandler) JobsByDepartmentHandler(c *gin.Context) {
	isocode := c.Query("isocode")
	if isocode == "" {
		c.JSON(400, gin.H{"error": "isocode is required"})
		return
	}

	results, err := h.es.GetJobsByDepartment(indexName, isocode)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, results)
}
