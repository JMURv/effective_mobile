[![Go Coverage](https://github.com/JMURv/effective_mobile/wiki/coverage.svg)](https://raw.githack.com/wiki/JMURv/effective_mobile/coverage.html)

---

## Stack
|         | Technology |
|---------|------------|
| Backend | Golang     |
| DB      | Postgres   |
| Cache   | Redis      |

---
## Configuration
### App
Configuration files placed in `/build/configs/envs`
- Create your own `.env` based on `.env.example`:
```shell
cp build/configs/envs/.env.example build/configs/envs/.env.dev && \
cp build/configs/envs/.env.example build/configs/envs/.env.prod
```

### Docker Compose
Run dev (requires `build/configs/envs/.env.dev`):
```shell
task dc-dev-build
```

Run dev with observation containers:
```shell
task dc-dev-obs
```

Observe profile starts svcs like: `prometheus`, `jaeger`, `node-exporter`, `grafana` and etc.

Services are available at:

| Сервис           | Адрес                  |
|------------------|------------------------|
| App (HTTP)       | http://localhost:8080  |
| App (PROMETHEUS) | http://localhost:8085  |
| Prometheus       | http://localhost:9090  |
| Node-exporter    | http://localhost:9100  |
| Jaeger           | http://localhost:16686 |
| Loki             | http://localhost:3100  |
| Grafana          | http://localhost:3000  |

More information could be found inside `compose.yaml`.

___


## Tests
```shell
task t
```
Will run all unit and integration tests.
It will spin up all containers for integration testing automatically using `testcontainers`.