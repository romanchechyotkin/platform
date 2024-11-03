# platform

Library contains 
 - logger which sends logs into [Loki](https://github.com/grafana/loki)
 - [Prometheus](https://prometheus.io/) metrics 
 - tracing into [Jaeger](https://www.jaegertracing.io/)
 - reusable code for [PostgreSQL](https://www.postgresql.org), [MinIO](https://min.io/), HTTP Server, gRPC and etc 

### usage
```env
    export GOPRIVATE=github.com/TakeAway-Inc/*
    export GONOSUMDB=github.com/TakeAway-Inc/*
    go get github.com/TakeAway-Inc/platform
```