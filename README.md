# Elasticsearch to GELF Relay

*es-to-gelf-relay* is a small web server application that redirect ElasticSearch log events to Graylog GELF protocol .

The main use case for this project is to send logs from *[AWS EKS Fargate using FluentBit output ES](https://docs.aws.amazon.com/eks/latest/userguide/fargate-logging.html)* to Graylog

## Build

```
docker build -t felixgarciaborrego/es-to-gelf-relay:0.0.2 . 
docker push felixgarciaborrego/es-to-gelf-relay:0.0.2
```
# Run it locally:

```
docker run -p 8000:8000 felixgarciaborrego/es-to-gelf-relay:0.0.1
```

## AWS EKS Fargate config


```
kind: Namespace
apiVersion: v1
metadata:
  name: aws-observability
  labels:
    aws-observability: enabled

---
kind: ConfigMap
apiVersion: v1
metadata:
  name: aws-logging
  namespace: aws-observability
  labels:
data:
  output.conf: |
    [OUTPUT]
        Name es
        Match  *
        Host <your-host>
        Port 80
        Index <index-name>
        Type  aks

```