# Real-time Chat Platform Implementation Roadmap

**Goal:** Build a hyper-scale Real-time Chat Platform (1M concurrent connections) in Go, focusing on zero-allocation WebSocket parsing, NATS Pub/Sub fanout, DynamoDB micro-batching, and ECS Fargate deployment via NLB.

---

## Phase 1: Monorepo & Tooling Setup

**Goal:** Establish the project skeleton, linting, and basic CI.
**Risk/Cost:** Low. Sets the foundation for team collaboration.

- [ ] **DoD:** `make lint` and `make test` run successfully on empty modules.

| #   | Service/Component | File path                        | Goal                                      | Verify commands               | Depends on |
| --- | ----------------- | -------------------------------- | ----------------------------------------- | ----------------------------- | ---------- |
| 1.1 | Root              | `go.mod`                         | Initialize Go workspace/module            | `go mod init chat-platform`   | None       |
| 1.2 | Root              | `Makefile`                       | Add targets for build, test, lint         | `make help`                   | 1.1        |
| 1.3 | Root              | `.golangci.yml`                  | Configure strict linting rules            | `golangci-lint run`           | 1.1        |
| 1.4 | CI                | `.github/workflows/ci.yml`       | Create basic GitHub Actions lint/test job | Push to GitHub, check Actions | 1.3        |
| 1.5 | All               | `src/api-core`, `src/edge-ws`... | Create empty directories for 4 services   | `tree src/`                   | 1.1        |

## Phase 2: DevOps - Docker & Local Compose Foundation

**Goal:** Setup multi-stage Dockerfiles and the local infrastructure (Kafka, Redis, NATS, Postgres, DynamoDB).
**Risk/Cost:** Medium. Ensuring all services can run locally before coding.

- [ ] **DoD:** `docker compose up` brings up all databases/brokers and placeholder app containers.

| #   | Service/Component | File path            | Goal                                              | Verify commands                             | Depends on |
| --- | ----------------- | -------------------- | ------------------------------------------------- | ------------------------------------------- | ---------- |
| 2.1 | All               | `Dockerfile` (x4)    | Create multi-stage Alpine Dockerfiles             | `docker build -t api-core src/api-core`     | 1.5        |
| 2.2 | Root              | `.dockerignore`      | Exclude vendor, local env files                   | `cat .dockerignore`                         | 2.1        |
| 2.3 | Infra             | `docker-compose.yml` | Add Postgres 15, Redis 7 (for presence), and NATS | `docker compose up -d postgres redis nats`  | 2.1        |
| 2.4 | Infra             | `docker-compose.yml` | Add Kafka (Zookeeper/Kraft) & DynamoDB Local      | `docker compose up -d kafka dynamodb-local` | 2.3        |
| 2.5 | App               | `docker-compose.yml` | Add app services with `depends_on`                | `docker compose up -d`                      | 2.4        |

## Phase 3: Shared Contracts & Configuration

**Goal:** Define the API models, WebSocket payloads, and environment configuration loaders.
**Risk/Cost:** Low.

- [ ] **DoD:** Structs defined, config loads from env vars.

| #   | Service/Component | File path                            | Goal                                | Verify commands        | Depends on |
| --- | ----------------- | ------------------------------------ | ----------------------------------- | ---------------------- | ---------- |
| 3.1 | All               | `src/pkg/config/`                    | Viper wrapper for env var loading   | `go test ./pkg/config` | 1.1        |
| 3.2 | API               | `src/api-core/internal/domain`       | Define User, Group entities         | `go build ./...`       | 3.1        |
| 3.3 | Edge              | `src/edge-ws/internal/models`        | Define WS JSON payloads (Auth, Msg) | `go build ./...`       | 3.1        |
| 3.4 | Router            | `src/message-router/internal/events` | Define Kafka/NATS message structs   | `go build ./...`       | 3.3        |

## Phase 4: API Core - Domain & Database Layer

**Goal:** Implement the REST API's Postgres persistence using Clean Architecture.
**Risk/Cost:** Medium.

- [ ] **DoD:** CRUD operations on Users and Groups work locally with Postgres.

| #   | Service/Component | File path                          | Goal                                     | Verify commands                   | Depends on |
| --- | ----------------- | ---------------------------------- | ---------------------------------------- | --------------------------------- | ---------- |
| 4.1 | API               | `db/migrations/`                   | Create Flyway/Golang-migrate SQL scripts | `make migrate-up`                 | 2.3        |
| 4.2 | API               | `src/api-core/internal/adapter/db` | Postgres connection & connection pool    | `go run cmd/api/main.go`          | 4.1        |
| 4.3 | API               | `src/api-core/internal/domain`     | Define Repository interfaces             | `golangci-lint run`               | 3.2        |
| 4.4 | API               | `src/api-core/internal/adapter/db` | Implement Postgres Repositories          | `go test -tags=integration ./...` | 4.3        |
| 4.5 | API               | `src/api-core/internal/usecase`    | Implement User/Group business logic      | `go test -v ./...`                | 4.4        |

## Phase 5: API Core - HTTP Edge & Auth

**Goal:** Expose REST endpoints with JWT authentication.
**Risk/Cost:** Low.

- [ ] **DoD:** `curl` commands can fetch tokens and create groups.

| #   | Service/Component | File path                            | Goal                             | Verify commands               | Depends on |
| --- | ----------------- | ------------------------------------ | -------------------------------- | ----------------------------- | ---------- |
| 5.1 | API               | `src/api-core/internal/adapter/http` | Setup Chi router & middleware    | `curl localhost:8081/health`  | 4.5        |
| 5.2 | API               | `src/pkg/auth/`                      | JWT generation & parsing utility | `go test ./pkg/auth`          | 3.1        |
| 5.3 | API               | `src/api-core/internal/adapter/http` | Implement Auth & User endpoints  | `curl -X POST /v1/users/auth` | 5.2        |

## Phase 6: Edge WS - WebSocket Connection Management

**Goal:** Implement the high-performance WebSocket edge using `epoll` / `gobwas/ws`.
**Risk/Cost:** High. Crucial for handling 1M connections.

- [ ] **DoD:** 50K simulated connections hold stable memory per node.

| #   | Service/Component | File path                         | Goal                                        | Verify commands                                         | Depends on |
| --- | ----------------- | --------------------------------- | ------------------------------------------- | ------------------------------------------------------- | ---------- |
| 6.1 | Edge              | `src/edge-ws/cmd/edge/main.go`    | HTTP server for WS upgrade endpoint         | `curl -i -N -H "Connection: Upgrade" localhost:8080/ws` | 3.3        |
| 6.2 | Edge              | `src/edge-ws/internal/connection` | Implement `gobwas/ws` zero-alloc upgrade    | `go run cmd/edge/main.go`                               | 6.1        |
| 6.3 | Edge              | `src/edge-ws/internal/connection` | Implement Epoll loop (linux) / Kqueue (mac) | `go build ./...`                                        | 6.2        |
| 6.4 | Edge              | `src/edge-ws/internal/pool`       | `sync.Pool` for reading WS frames           | `go test -bench . ./internal/pool`                      | 6.3        |
| 6.5 | Edge              | `src/edge-ws/internal/connection` | Parse JWT auth on first WS frame            | Connect via Postman WS                                  | 5.2        |

## Phase 7: Edge WS - Inbound Message Ingestion (Kafka)

**Goal:** Forward validated WebSocket messages to Kafka.
**Risk/Cost:** Medium.

- [ ] **DoD:** Messages sent via WS appear in Kafka topic.

| #   | Service/Component | File path                     | Goal                                                | Verify commands                     | Depends on |
| --- | ----------------- | ----------------------------- | --------------------------------------------------- | ----------------------------------- | ---------- |
| 7.1 | Edge              | `src/pkg/kafka`               | Sarama Kafka Producer wrapper                       | `go test ./pkg/kafka`               | 2.4        |
| 7.2 | Edge              | `src/edge-ws/internal/router` | Read frame -> validate -> publish to `chat.inbound` | WS send -> `kafka-console-consumer` | 6.5        |

## Phase 8: Message Router - DynamoDB Micro-batching

**Goal:** Consume Kafka messages and persist to DynamoDB efficiently using micro-batching.
**Risk/Cost:** High. Micro-batching is critical to prevent astronomical DynamoDB costs at 50K writes/sec.

- [ ] **DoD:** DynamoDB `BatchWriteItem` reduces write requests by 90% in logs.

| #   | Service/Component | File path                            | Goal                                              | Verify commands                     | Depends on |
| --- | ----------------- | ------------------------------------ | ------------------------------------------------- | ----------------------------------- | ---------- |
| 8.1 | Router            | `src/pkg/dynamo`                     | AWS SDK v2 DynamoDB client                        | `go test ./pkg/dynamo`              | 2.4        |
| 8.2 | Router            | `scripts/dynamodb_schema.sh`         | Create DynamoDB local tables                      | `./scripts/dynamodb_schema.sh`      | 8.1        |
| 8.3 | Router            | `src/message-router/internal/worker` | Create buffered channel for incoming messages     | Unit test channel buffer limits     | 8.1        |
| 8.4 | Router            | `src/message-router/internal/worker` | Implement `time.Ticker` (100ms) for flush trigger | Unit test flush timing              | 8.3        |
| 8.5 | Router            | `src/message-router/internal/db`     | Implement DynamoDB `BatchWriteItem`               | `go test -tags=integration ./...`   | 8.4        |
| 8.6 | Router            | `src/pkg/kafka`                      | Kafka Consumer Group wrapper                      | `go test ./pkg/kafka`               | 7.1        |
| 8.7 | Router            | `src/message-router/internal/worker` | Consume `chat.inbound` -> send to batch channel   | Run edge + router, send 100 WS msgs | 8.5        |

## Phase 9: Presence Tracking & NATS Fanout

**Goal:** Map users to Edge nodes via Redis and route messages via NATS.
**Risk/Cost:** Medium. NATS replaces Redis PubSub for extreme throughput.

- [ ] **DoD:** Alice sends message to Bob, Bob receives it on a different WS node via NATS.

| #   | Service/Component | File path                            | Goal                                                | Verify commands                      | Depends on |
| --- | ----------------- | ------------------------------------ | --------------------------------------------------- | ------------------------------------ | ---------- |
| 9.1 | Edge              | `src/edge-ws/internal/connection`    | On connect: Set `user:{id}:node = node_id` in Redis | `redis-cli GET user:X:node`          | 6.5        |
| 9.2 | Router            | `src/message-router/internal/worker` | After batching, lookup recipient node in Redis      | Check logs during routing            | 8.7        |
| 9.3 | Router            | `src/pkg/nats`                       | Implement NATS Go client connection                 | `go test ./pkg/nats`                 | 2.3        |
| 9.4 | Router            | `src/message-router/internal/worker` | Publish to NATS topic `node.{node_id}`              | Subscribe via NATS CLI               | 9.3        |
| 9.5 | Edge              | `src/edge-ws/internal/router`        | Edge subscribes to its NATS `node.{node_id}` topic  | `go run cmd/edge/main.go`            | 9.4        |
| 9.6 | Edge              | `src/edge-ws/internal/router`        | NATS msg -> send to specific local Epoll FD         | End-to-end WS chat between 2 clients | 9.5        |

## Phase 10: E2E Testing & API History Integration

**Goal:** Allow users to fetch chat history from DynamoDB via REST API.
**Risk/Cost:** Low.

- [ ] **DoD:** `/v1/channels/{id}/messages` returns sorted history.

| #    | Service/Component | File path                            | Goal                                          | Verify commands                              | Depends on |
| ---- | ----------------- | ------------------------------------ | --------------------------------------------- | -------------------------------------------- | ---------- |
| 10.1 | API               | `src/api-core/internal/adapter/db`   | Implement DynamoDB FetchMessages (Pagination) | `go test -tags=integration ./...`            | 8.2        |
| 10.2 | API               | `src/api-core/internal/adapter/http` | REST endpoint for history                     | `curl localhost:8081/v1/channels/1/messages` | 10.1       |

## Phase 11: Observability, Metrics & SLO

**Goal:** Instrument all services with RED metrics using Prometheus/OTel.
**Risk/Cost:** Medium.

- [ ] **DoD:** Grafana dashboard shows WS connection count and NATS message throughput.

| #    | Service/Component | File path            | Goal                                         | Verify commands                     | Depends on |
| ---- | ----------------- | -------------------- | -------------------------------------------- | ----------------------------------- | ---------- |
| 11.1 | All               | `src/pkg/metrics/`   | Prometheus SDK wrapper & `/metrics` endpoint | `curl localhost:8080/metrics`       | 5.1        |
| 11.2 | Infra             | `docker-compose.yml` | Add Prometheus & Grafana                     | `docker compose up -d prom grafana` | 11.1       |

## Phase 12: Load Harness & Resilience

**Goal:** Prove the system can handle 100K+ connections locally and test failure modes.
**Risk/Cost:** High.

- [ ] **DoD:** `tcpkali` or `k6` sustains 100K WS connections locally without OOM.

| #    | Service/Component | File path                      | Goal                                               | Verify commands      | Depends on |
| ---- | ----------------- | ------------------------------ | -------------------------------------------------- | -------------------- | ---------- |
| 12.1 | Tests             | `scripts/load/k6_ws.js`        | Write k6 script to auth and hold WS connection     | `k6 run k6_ws.js`    | 6.5        |
| 12.2 | Edge              | `src/edge-ws/cmd/edge/main.go` | Tune `ulimit -n` and OS limits locally for testing | Check limits, run k6 | 12.1       |

## Phase 13: Terraform Modules & AWS IaC

**Goal:** Provision AWS infrastructure including NLBs for the Edge.
**Risk/Cost:** High.

- [ ] **DoD:** `terraform apply` successfully creates the staging environment.

| #    | Service/Component | File path                       | Goal                                       | Verify commands                        | Depends on |
| ---- | ----------------- | ------------------------------- | ------------------------------------------ | -------------------------------------- | ---------- |
| 13.1 | Infra             | `infra/terraform/modules/vpc`   | VPC, Subnets, NAT Gateway, Security Groups | `terraform plan`                       | 1.1        |
| 13.2 | Infra             | `infra/terraform/modules/rds`   | Postgres RDS (Multi-AZ) & DynamoDB         | `terraform apply -target=module.rds`   | 13.1       |
| 13.3 | Infra             | `infra/terraform/modules/redis` | ElastiCache Redis cluster & NATS setup     | `terraform apply -target=module.redis` | 13.1       |
| 13.4 | Infra             | `infra/terraform/modules/ecs`   | ECS Cluster, NLB for Edge, ALB for API     | `terraform apply -target=module.ecs`   | 13.1       |

## Phase 14: AWS Compute - ECS Task Definitions (NLB)

**Goal:** Deploy the Go services to ECS Fargate behind an NLB.
**Risk/Cost:** Medium.

- [ ] **DoD:** ECS services are GREEN and NLB routes raw TCP traffic correctly.

| #    | Service/Component | File path                     | Goal                                            | Verify commands                         | Depends on |
| ---- | ----------------- | ----------------------------- | ----------------------------------------------- | --------------------------------------- | ---------- |
| 14.1 | Infra             | `infra/terraform/modules/ecs` | Task Def for `api-core` & ALB Target Group      | AWS Console: Task Running               | 13.4       |
| 14.2 | Infra             | `infra/terraform/modules/ecs` | Task Def for `edge-ws` & NLB Target Group (TCP) | AWS Console: Target Group Healthy       | 14.1       |
| 14.3 | Infra             | `infra/terraform/modules/ecs` | Task Def for `message-router` (Worker, no LB)   | AWS Console: Logs show Kafka connection | 14.2       |

## Phase 15: CI/CD Pipeline - Deploy to Staging

**Goal:** Automate the build and deployment process via GitHub Actions.
**Risk/Cost:** Low.

- [ ] **DoD:** Merging to `main` auto-deploys to AWS ECS.

| #    | Service/Component | File path                      | Goal                                          | Verify commands             | Depends on |
| ---- | ----------------- | ------------------------------ | --------------------------------------------- | --------------------------- | ---------- |
| 15.1 | CI                | `.github/workflows/deploy.yml` | Build image and push to Amazon ECR            | Verify image in ECR console | 2.1        |
| 15.2 | CI                | `.github/workflows/deploy.yml` | Update ECS Task Definition and force redeploy | `aws ecs describe-services` | 14.3       |

## Phase 16: Prod Readiness & Final Verification

**Goal:** Prepare the system for production-like load.
**Risk/Cost:** Low.

- [ ] **DoD:** Load test against the live AWS NLB succeeds.

| #    | Service/Component | File path                | Goal                                          | Verify commands          | Depends on |
| ---- | ----------------- | ------------------------ | --------------------------------------------- | ------------------------ | ---------- |
| 16.1 | Docs              | `docs/runbook.md`        | Document alert responses (Kafka lag, RDS CPU) | Read document            | 11.1       |
| 16.2 | Load              | `scripts/load/k6_aws.js` | Run k6 against the NLB                        | `k6 run` -> View Grafana | 12.1       |

---

## Session Log

| Date       | Phase | Step      | Status  | Note                                                 |
| ---------- | ----- | --------- | ------- | ---------------------------------------------------- |
| 2024-04-10 | 1     | 1.1       | ✅ Done | Initialized module as `pulse`                        |
| 2024-04-10 | 1     | 1.2       | ✅ Done | Added `Makefile`                                     |
| 2024-04-10 | 1     | 1.3 - 1.5 | ✅ Done | CI/CD, Linters, Directories                          |
| 2024-04-17 | 2     | 2.1 - 2.2 | ✅ Done | Added Multi-stage Dockerfiles & Ignore               |
| 2024-04-17 | 2     | 2.3 - 2.5 | ✅ Done | Setup Local Infra (Kafka, Postgres, NATS, DynamoDB)  |
| 2024-04-24 | 3     | 3.1       | ✅ Done | Added Config loader with Viper                       |
| 2024-04-24 | 3     | 3.2       | ✅ Done | Defined API Core Entities                            |
| 2024-04-24 | 3     | 3.3       | ✅ Done | Defined WS JSON payloads with json.RawMessage        |
| 2024-04-24 | 3     | 3.4       | ✅ Done | Defined Kafka/NATS RoutedMessage event               |
| 2024-05-02 | 4     | 4.1       | ✅ Done | Created Postgres migration scripts with UUID         |
| 2024-05-10 | 4     | 4.2       | ✅ Done | Implemented Postgres connection pool with pgx        |
| 2024-05-17 | 4     | 4.3       | ✅ Done | Defined Repository interfaces for Clean Architecture |
| 2024-05-18 | 4     | 4.4       | ✅ Done | Implemented Postgres Repository for Users            |
| 2024-05-18 | 4     | 4.5       | ✅ Done | Implemented User business logic (UseCase)            |
| 2024-05-25 | 4     | 4.4       | ✅ Done | Implemented Postgres Repository for Groups           |
| 2024-05-25 | 4     | 4.5       | ✅ Done | Implemented Group UseCase with Clean Architecture    |
