package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"es-job-search/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

type IEsConnector interface {
	IndexJobs(indexName string) error
	CleanIndex(indexName string) error
	SearchJobs(indexName string, keyword string) ([]models.Job, error)
	GetJobsByDepartment(indexName string, isocode string) (map[string]interface{}, error)
}

type EsConnector struct {
	Client *elastic.Client
}

const (
	jobsFileName = "resources/jobs.json"
	esUrl        = "http://localhost:9200"
)

func NewEsConnector() (IEsConnector, error) {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetSniff(false),                              // disable sniffing in docker environment
		elastic.SetHealthcheckInterval(10*time.Second),       // set healthcheck interval
		elastic.SetHealthcheckTimeoutStartup(30*time.Second), // set startup healthcheck timeout
	)
	if err != nil {
		return nil, err
	}

	info, code, err := client.Ping(esUrl).Do(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	return &EsConnector{
		Client: client,
	}, nil
}

func (es *EsConnector) CleanIndex(indexName string) error {
	exists, err := es.Client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		return err
	}

	if exists {
		_, err := es.Client.DeleteIndex(indexName).Do(context.Background())
		if err != nil {
			return err
		}
	}

	_, err = es.Client.CreateIndex(indexName).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (es *EsConnector) IndexJobs(indexName string) error {
	var jobs []models.Job

	file, err := os.Open(jobsFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jobs); err != nil {
		return err
	}

	log.Println("Starting indexing jobs...")

	for _, job := range jobs {
		_, err := es.Client.Index().
			Index(indexName).
			BodyJson(job).
			Do(context.Background())
		if err != nil {
			return err
		}
	}

	log.Println("Finished indexing jobs.")
	return nil
}

func (es *EsConnector) SearchJobs(indexName string, keyword string) ([]models.Job, error) {
	query := elastic.NewMultiMatchQuery(keyword, "title").Fuzziness("AUTO")
	searchResult, err := es.Client.Search().
		Index(indexName).
		Query(query).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var jobs []models.Job
	for _, hit := range searchResult.Hits.Hits {
		var job models.Job
		err := json.Unmarshal(hit.Source, &job)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (es *EsConnector) GetJobsByDepartment(indexName string, isocode string) (map[string]interface{}, error) {
	termAgg := elastic.NewTermsAggregation().Field("location.department.isoCode.keyword")
	searchResult, err := es.Client.Search().
		Index(indexName).
		Aggregation("jobs_by_department", termAgg).
		Query(elastic.NewTermQuery("location.department.isoCode.keyword", isocode)).
		Size(1000).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	agg, found := searchResult.Aggregations.Terms("jobs_by_department")
	if !found {
		return nil, errors.New("aggregation not found")
	}

	var jobs []models.Job
	for _, hit := range searchResult.Hits.Hits {
		var job models.Job
		err := json.Unmarshal(hit.Source, &job)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	results := map[string]interface{}{
		"total": agg.Buckets[0].DocCount,
		"jobs":  jobs,
	}

	return results, nil
}
