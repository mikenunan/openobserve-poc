# A-Bank OpenObserve PoC — Demo Script

> Guided walkthrough for presenting the PoC to the engineering team.
> Estimated time: **15–20 minutes**

---

## Setup (before the demo)

```bash
cd openobserve-poc
docker compose up -d
# Wait ~15 seconds for health checks
```

Verify all services are healthy:

```bash
docker compose ps
```

Open two browser tabs:

1. **Frontend:** http://localhost:8080
2. **OpenObserve:** http://localhost:5080 (login: admin@example.com / `Complexpass#123`)

---

## Part 1 — The Application (5 min)

### 1.1 Customer List

- Open the frontend at http://localhost:8080
- You'll see the **Customer List** table with 5 pre-loaded UK bank customers
- Point out: _"This is a React 19 + TypeScript SPA served through an nginx reverse proxy. All API calls go via /api/ routes — no CORS configuration needed."_

### 1.2 Customer Detail

- Click on **Jane Doe** (Party ID 12345)
- Show the editable form with personal details and address
- Point out the **Affiliated Documents** section — two documents are pre-seeded:
    - "A-Bank customer Ts&Cs"
    - "FSCS Info Sheet"

### 1.3 Document Viewer

- Click on **"FSCS Info Sheet"** in the affiliations list
- The document text appears inline below — this is fetched from the DocumentCatalogue service
- Point out: _"This single click triggered calls across three microservices — the frontend called DocumentFiling for the affiliation list, and DocumentCatalogue for the document content."_

### 1.4 Edit a Customer

- Change Jane's last name to **"Smith"**
- Click **Amend** — the success message confirms the update
- Click **← Back to List** — the table now shows "Jane Smith"
- Point out: _"That PUT request went through nginx → PartyMan, with payload validation and structured logging."_

---

## Part 2 — Observability in OpenObserve (10 min)

### 2.1 Traces

- Switch to the **OpenObserve** tab (http://localhost:5080)
- Navigate to **Traces** in the left sidebar
- You should see traces from the recent interactions
- Click on a trace to expand it — show the spans:
    - The HTTP request entering via the service
    - Internal processing spans
- Point out: _"Each service is instrumented with OpenTelemetry. Traces propagate across service boundaries using W3C Trace Context headers."_

### 2.2 Cross-Service Trace

- From a terminal, create a new affiliation to generate a cross-service trace:

```bash
curl -X POST http://localhost:8080/api/filing/affiliation \
  -H "Content-Type: application/json" \
  -d '{"partyId":12346,"documentName":"Savings Account Ts&Cs"}'
```

- Refresh traces in OpenObserve
- Find the POST /affiliation trace — it should show:
    1. DocumentFiling receives the request
    2. HEAD call to PartyMan (customer validation)
    3. HEAD call to DocumentCatalogue (document validation)
    4. Affiliation created
- Point out: _"This is the key value of distributed tracing — we can see the full request flow across all three services in a single trace view. In production, this helps us diagnose latency, failures, and bottlenecks."_

### 2.3 Logs

- Navigate to **Logs** in OpenObserve
- Search for recent log entries
- Point out the structured JSON format with fields like `service`, `level`, `msg`, and contextual attributes (partyId, documentName)
- Point out: _"All services use Go's slog library with JSON output. The OTel Collector forwards these to OpenObserve, giving us logs and traces in one unified platform."_

### 2.4 Validation Failure Trace (optional)

- Trigger a validation failure:

```bash
curl -X POST http://localhost:8080/api/filing/affiliation \
  -H "Content-Type: application/json" \
  -d '{"partyId":99999,"documentName":"Savings Account Ts&Cs"}'
```

- This returns 422 — the customer doesn't exist
- Show the trace in OpenObserve — the HEAD call to PartyMan returned 404
- Point out: _"Even failure paths are fully traced. This is invaluable for debugging in production."_

---

## Part 3 — Architecture Discussion (5 min)

### Key Talking Points

1. **Zero-config CORS** — nginx reverse proxy means the frontend, API, and observability all share a single origin
2. **Lightweight validation** — HEAD requests for existence checks (no data transfer, just 200/404)
3. **Structured logging** — `slog` with JSON output, zero third-party dependencies
4. **OTel is vendor-agnostic** — we could swap OpenObserve for Jaeger, Grafana Tempo, or Datadog with only collector config changes
5. **Docker Compose today → Kubernetes tomorrow** — every service has health checks, graceful shutdown, and OTel instrumentation. The path to GKE/EKS is straightforward

### Questions to Prompt

- _"What observability gaps do we have in our current platform?"_
- _"How does this compare to our existing logging/tracing setup?"_
- _"What would we need to add for a production deployment?"_ (Auth, persistent storage, alerting, dashboards)

---

## Cleanup

```bash
docker compose down -v
```
