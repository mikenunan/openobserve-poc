#!/usr/bin/env bash
# Adds a document affiliation for John Smith (12346) — not in the pre-seeded set.
# Triggers cross-service validation: HEAD → PartyMan, HEAD → DocumentCatalogue.

set -euo pipefail

echo "Creating affiliation: partyId=12346 (John Smith) ← \"Savings Account Ts&Cs\""

curl -s -w "\nHTTP %{http_code}\n" \
  -X POST http://localhost:8080/api/filing/affiliation \
  -H "Content-Type: application/json" \
  -d '{"partyId":12346,"documentName":"Savings Account Ts&Cs"}'
