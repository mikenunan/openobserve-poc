# A-Bank Customer & Document Management PoC

> **OpenObserve + OpenTelemetry Observability Demonstration**

## Purpose

This proof-of-concept demonstrates how **OpenObserve** and **OpenTelemetry** can provide unified observability — distributed traces, structured logs, and metrics — across a microservices architecture representative of internal banking systems.

The application is a simplified customer record and document management system built with Go microservices and a React frontend, fully containerised with Docker Compose.

| Item              | Detail                                                    |
| ----------------- | --------------------------------------------------------- |
| **Audience**      | Engineering team                                          |
| **Stack**         | Go 1.25, React 19, TypeScript, nginx, Docker Compose      |
| **Observability** | OpenTelemetry SDK (Go + Web), OTel Collector, OpenObserve |

---

## Architecture

```
Browser → nginx (:8080)
             ├── /                   → React Frontend (:5173)
             ├── /api/partyman/      → PartyMan (:8081)
             ├── /api/catalogue/     → DocumentCatalogue (:8082)
             └── /api/filing/        → DocumentFiling (:8083)

Go Services → OTel Collector (:4317) → OpenObserve (:5080)
```

### Services

| Service               | Description                                                      |
| --------------------- | ---------------------------------------------------------------- |
| **PartyMan**          | Customer CRUD — GET/POST/PUT/HEAD endpoints, 5 seed records      |
| **DocumentCatalogue** | Serves 3 banking document texts (Ts&Cs, Savings, FSCS)           |
| **DocumentFiling**    | Associates documents with customers, cross-service validation    |
| **Frontend**          | React SPA — customer list, editable detail view, document viewer |

For the full architecture plan with diagrams, see [`docs/architecture-plan.md`](docs/architecture-plan.md).

---

## Quick Start

### Prerequisites

- **Docker** with Docker Compose v2
- ~2 GB free disk space (for container images)

### Run

```bash
# Clone and start everything
git clone https://github.com/mikenunan/openobserve-poc.git
cd openobserve-poc
docker compose up -d

# Wait for health checks (~15 seconds), then open:
#   Frontend:     http://localhost:8080
#   OpenObserve:  http://localhost:5080
```

### Stop

```bash
docker compose down        # Stop containers
docker compose down -v     # Stop and remove data volumes
```

### Rebuild after changes

```bash
docker compose up -d --build
```

---

## Access Points

| Component                 | URL                                  | Credentials                           |
| ------------------------- | ------------------------------------ | ------------------------------------- |
| **Frontend**              | http://localhost:8080                | —                                     |
| **OpenObserve**           | http://localhost:5080                | admin@example.com / `Complexpass#123` |
| **PartyMan API**          | http://localhost:8080/api/partyman/  | —                                     |
| **DocumentCatalogue API** | http://localhost:8080/api/catalogue/ | —                                     |
| **DocumentFiling API**    | http://localhost:8080/api/filing/    | —                                     |

---

## API Quick Reference

```bash
# List all customers
curl http://localhost:8080/api/partyman/customers

# Get a single customer
curl http://localhost:8080/api/partyman/customer?partyId=12345

# Update a customer
curl -X PUT http://localhost:8080/api/partyman/customer \
  -H "Content-Type: application/json" \
  -d '{"partyId":12345,"firstName":"Jane","lastName":"Smith","addressLine1":"10 Downing Street","city":"London","postcode":"SW1A 2AA","country":"United Kingdom"}'

# Get affiliated documents for a customer
curl http://localhost:8080/api/filing/documents?partyId=12345

# Create a new affiliation (validates customer + document exist)
curl -X POST http://localhost:8080/api/filing/affiliation \
  -H "Content-Type: application/json" \
  -d '{"partyId":12346,"documentName":"A-Bank customer Ts&Cs"}'

# View a document
curl "http://localhost:8080/api/catalogue/document?name=FSCS%20Info%20Sheet"

# Health checks
curl http://localhost:8080/api/partyman/health
curl http://localhost:8080/api/catalogue/health
curl http://localhost:8080/api/filing/health
```

---

## Pre-Seeded Data

**Customers** (loaded at startup by PartyMan):

| PartyID | Name             | City       |
| ------- | ---------------- | ---------- |
| 12345   | Jane Doe         | London     |
| 12346   | John Smith       | London     |
| 12347   | Anya Patel       | Edinburgh  |
| 12348   | Liam O'Brien     | Cardiff    |
| 12349   | Fatima Al-Hassan | Birmingham |

**Affiliations** (pre-seeded in DocumentFiling):

| PartyID | Document              |
| ------- | --------------------- |
| 12345   | A-Bank customer Ts&Cs |
| 12345   | FSCS Info Sheet       |
| 12347   | Savings Account Ts&Cs |

---

## Known Limitations (PoC Scope)

- **No authentication/authorisation** — all endpoints are unauthenticated
- **No DELETE endpoints** — create and update only
- **In-memory storage** — data resets on container restart
- **No affiliation UI** — affiliations created via curl/Postman (mirrors real-world mobile app flow)
- **PUT replaces entire record** — no partial updates (PATCH)

---

## Project Structure

```
openobserve-poc/
├── docker-compose.yml
├── nginx/nginx.conf
├── observability/otel-collector-config.yaml
├── docs/
│   ├── architecture-plan.md
│   └── demo-script.md
├── services/
│   ├── partyman/
│   ├── documentcatalogue/
│   └── documentfiling/
└── frontend/
```

---

## Further Reading

- [Architecture Plan](docs/architecture-plan.md) — full design, diagrams, and task breakdown
- [Demo Script](docs/demo-script.md) — guided walkthrough for presenting the demo
- [Implementation Walkthrough](docs/walkthrough.md) — what was built, how it was verified, and design decisions
