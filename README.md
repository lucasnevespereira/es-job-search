# es-job-search

Sample project to practice elastic search indexing, document search and bucket aggregations.

## Usage
```
make es
```

```
make run
```

## Curl Endpoints Commands

Index jobs from `resources/jobs.json` file

```shell
curl http://localhost:8080/index-jobs
```

Search jobs by keyword

```shell
curl "http://localhost:8080/search-jobs?keyword=engineer" | jq
```

Get jobs of a department

```shell
curl "http://localhost:8080/jobs-by-department?isocode=FR-92" | jq
```