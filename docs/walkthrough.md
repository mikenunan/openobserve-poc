# A-Bank OpenObserve PoC — Implementation Walkthrough

> Technical summary of what was built, how it was verified, and key design decisions.

---

## Architecture Overview

```
Browser → nginx (:8080) → PartyMan (:8081) / DocumentCatalogue (:8082) / DocumentFiling (:8083)
                        → Frontend (:5173)
                        → OTel Collector (:4317/4318) → OpenObserve (:5080)
```

All 7 containers orchestrated via Docker Compose with health checks, graceful shutdown, and automatic restart.

---

## Services Built

| Service               | Language              | Endpoints                                                      | Key Features                                                                                                     |
| --------------------- | --------------------- | -------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| **PartyMan**          | Go 1.25               | GET/POST/PUT/HEAD `/customer`, GET `/customers`, GET `/health` | 5 seed customers, field validation, slog structured logging                                                      |
| **DocumentCatalogue** | Go 1.25               | GET/HEAD `/document`, GET `/documents`, GET `/health`          | 3 UK banking documents served from embedded text files                                                           |
| **DocumentFiling**    | Go 1.25               | POST `/affiliation`, GET `/documents`, GET `/health`           | Cross-service validation (HEAD to PartyMan + DocumentCatalogue), duplicate prevention, 3 pre-seeded affiliations |
| **Frontend**          | React 19 + TypeScript | SPA at `/`                                                     | Customer list, editable detail view, inline document viewer, dark theme                                          |

---

## Observability Stack

| Component          | Role                                                 | Configuration                                                                   |
| ------------------ | ---------------------------------------------------- | ------------------------------------------------------------------------------- |
| **OTel SDK (Go)**  | Instruments each Go service with distributed tracing | gRPC export to OTel Collector, W3C Trace Context propagation                    |
| **OTel Web SDK**   | Instruments browser-side `fetch()` calls             | OTLP/HTTP export via `/api/otel/` nginx proxy                                   |
| **OTel Collector** | Receives, batches, and forwards telemetry            | OTLP gRPC receiver (:4317) + HTTP (:4318), exports to OpenObserve via OTLP/HTTP |
| **OpenObserve**    | Unified traces, logs, and metrics UI                 | Basic auth, persistent volume, OTLP/HTTP ingestion on :5080                     |

### Trace Propagation Flow

```
Browser fetch()  →  nginx  →  Go service (receives traceparent header)
                                    ↓
                              OTel SDK creates child span
                                    ↓
                              Cross-service call (HEAD) carries traceparent
                                    ↓
                              Downstream service creates child span
                                    ↓
                              All spans export to OTel Collector → OpenObserve
```

---

## Verification Results

| Test                                              | Result  |
| ------------------------------------------------- | ------- |
| Health checks (all 3 services)                    | ✅ Pass |
| Customer CRUD (GET/HEAD/POST/PUT)                 | ✅ Pass |
| Document retrieval (GET/HEAD, all 3 documents)    | ✅ Pass |
| Affiliations (GET pre-seeded + POST new)          | ✅ Pass |
| Cross-service validation (HEAD 200/404)           | ✅ Pass |
| Frontend renders at `:8080`                       | ✅ Pass |
| OTel proxy `/api/otel/v1/traces` accepts traces   | ✅ Pass |
| OTel Collector receives spans from Go services    | ✅ Pass |
| OTel Collector forwards to OpenObserve (HTTP 200) | ✅ Pass |
| Traces visible in OpenObserve UI                  | ✅ Pass |
| `docker compose up/down` clean lifecycle          | ✅ Pass |

---

## Key Design Decisions

| Decision                                                            | Rationale                                                                         |
| ------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| **gRPC for service→collector, OTLP/HTTP for collector→OpenObserve** | gRPC is the standard OTel transport; OpenObserve's native ingestion is HTTP-based |
| **nginx reverse proxy**                                             | Single origin for frontend + API + OTel — eliminates CORS entirely                |
| **HEAD requests for validation**                                    | Lightweight existence checks (200/404) with no data transfer                      |
| **`slog` for structured logging**                                   | Go stdlib, zero dependencies, JSON output for machine parsing                     |
| **In-memory stores with pre-seeded data**                           | PoC scope — demonstrates the full flow without database complexity                |
| **React Router v6**                                                 | v7 has an incompatible API for the BrowserRouter/Routes/Route pattern             |
| **OTel init wrapped in try-catch**                                  | Prevents telemetry failures from crashing the React application                   |

---

## Issues Encountered & Resolved

| Issue                               | Root Cause                                                                                                  | Fix                                    |
| ----------------------------------- | ----------------------------------------------------------------------------------------------------------- | -------------------------------------- |
| Blank frontend page                 | React Router v7 API incompatibility                                                                         | Pinned `react-router-dom` to `^6.28.0` |
| OTel init could crash React         | `initTelemetry()` called without error handling                                                             | Wrapped in try-catch                   |
| Go services failed to export traces | `OTEL_EXPORTER_OTLP_ENDPOINT` set as URL (`http://host:port`) but gRPC `WithEndpoint()` expects `host:port` | Removed `http://` prefix               |
| OTel Collector deprecation warning  | `otlphttp` exporter alias deprecated                                                                        | Renamed to `otlp_http`                 |

---

## Commit History

| Commit                                                   | Description                      |
| -------------------------------------------------------- | -------------------------------- |
| `chore: add .gitattributes and update .gitignore`        | Cross-platform line endings      |
| `feat(infra): add docker-compose, nginx, otel-collector` | 7-container orchestration        |
| `feat(partyman): implement customer management service`  | CRUD, validation, seed data      |
| `feat(documentcatalogue): implement document retrieval`  | 3 UK banking documents           |
| `feat(documentfiling): implement affiliation service`    | Cross-service validation         |
| `feat(frontend): implement React customer management UI` | List, detail, document viewer    |
| `fix(docker): update Go base image to 1.25-alpine`       | Match go.mod requirement         |
| `docs: add README`                                       | Project overview and quick start |
| `feat(frontend): add OTel Web SDK + nginx proxy`         | Browser-side tracing             |
| `fix(frontend): pin react-router-dom to v6`              | API compatibility                |
| `fix(frontend): guard OTel init with try-catch`          | Defensive error handling         |
| `chore(frontend): add package-lock.json`                 | Reproducible builds              |
| `fix(infra): remove http:// from OTLP endpoint`          | gRPC address format              |
| `fix(otel): rename otlphttp to otlp_http`                | Deprecation fix                  |
