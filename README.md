# Go Wayback!
A simple web app which will export a domains content from a given year from the [wayback machine](https://support.archive-it.org/hc/en-us/articles/115001790023-Access-Archive-It-s-Wayback-index-with-the-CDX-C-API).

# Architecture
```mermaid
flowchart LR
    prep[PreProcessor]
    p[Processor]
    postp[PostProcessor]
    BigQuery
    t1([Pub/Sub Topic])
    t2([Pub/Sub Topic])
    prep --> t1 --> p --> t2 --> postp
    t1 --> BigQuery
``````

Base URL http://web.archive.org/web/%vif_/%v

```bash
# Run the Server Local
task run
```