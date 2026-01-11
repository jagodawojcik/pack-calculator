# 📦 ✈️ Product Shipping Pack Calculator Service

This repository hosts a solution for the Gymshark coding challenge. The task involved building an application that determines the minimum number of packs needed to fulfill customer orders, given fixed pack sizes: 250, 500, 1000, 2000, and 5000.

**Key Task Requirements**

✅ Calculate optimal pack distributions for a specified order quantity following the rules in the provided priority order:

| Priority | Rule                                                       | Rationale                  |
| -------- | ---------------------------------------------------------- | -------------------------- |
| 1        | Only **whole packs** allowed                               | Cannot split pack contents |
| 2        | Send **no more items than necessary**                      | Minimize waste             |
| 3        | Use **as few packs as possible** respecting previous rules | Minimize operational cost  |

✅ App must be built in Go

✅ API must be HTTP-accessible

**Optional requirements**

✅ UI is provided to interact with the service

✅ Pack sizes are configurable without modifying the source code

## 💡 Solution

This solution meets all the key and optional requirements. The service can be accessed online via a basic UI:

https://jagodawojcik.github.io/pack-calculator/

The backend URL:

https://packs-service-253100512672.europe-west2.run.app/health

See the following sections for more detailed app documentation.

## 🔎 Table of Contents

- [Project Structure](#project-structure)
- [API Reference](#api-reference)
- [Algorithm: Pack Calculation Strategy](#algorithm-pack-calculation-strategy)
- [Infrastructure & Deployment](#infrastructure--deployment)
- [Configuration](#configuration)
- [CI/CD Pipeline](#cicd-pipeline)
- [Future Improvements & Considerations](#future-improvements--considerations)
- [Local Development](#local-development)

---

## Project Structure

```
.
├── cmd/
│   ├── cli/           # CLI tool for quick local testing of pack calculations
│   └── server/        # HTTP server with a single GET endpoint responsible for calculation
├── internal/
│   └── calculatepacks/  # Core pack calculation logic
├── terraform/         # IaC
├── Dockerfile         # Container image
├── go.mod            # Go dependencies (currently no external dependencies)
├── index.html        # UI frontend
├── styles.css        # UI styling
└── README.md
```

## API Reference

### Health Check Endpoint

```bash
GET /health
```

Basic health check. Returns `200 OK` if the service is up and running.

### Pack Calculation Endpoint

```bash
GET /packs?quantity=<number>
```

**Parameters:**

- `quantity` (required): Number of items to order

**Response:**

```json
{
  "packs": {
    "500": 1
  }
}
```

**Example Requests:**

```bash
# Order 251 items
curl "https://packs-service-253100512672.europe-west2.run.app/packs?quantity=251"
# Response: {"packs": {"500": 1}}

# Order 1501 items
curl "https://packs-service-253100512672.europe-west2.run.app/packs?quantity=1501"
# Response: {"packs": {"2000": 1}}

# Order 2501 items
curl "https://packs-service-253100512672.europe-west2.run.app/packs?quantity=2501"
# Response: {"packs": {"2000": 1, "500": 1}}
```

**Error Responses:**

- `400 Bad Request`: Missing or invalid item quantity parameter
- `5xx Internal Server Error`: Unexpected server error

---

## 🧮 Algorithm: Pack Calculation Strategy

The pack calculation algorithm uses **dynamic programming** to find the optimal solution:

### Algorithm Steps

1. **Exact Match Check**: If the order quantity exactly matches a pack size, return 1 pack of that size immediately.

2. **Dynamic Programming Table**: Build a DP table where `best[target]` stores the best solution for reaching that target quantity:

   - Iterates from 1 up to the requested quantity
   - For each target considers all available pack sizes
   - Selects the pack that satisfies constraints in priority order:
     - Minimizes total items shipped (primary constraint)
     - Then minimizes number of packs (secondary constraint)

3. **Solution Reconstruction**: Backtrack from `best[target]` through the DP table to construct the final pack distribution.

### Time & Memory Complexity

Where N = order quantity and M = number of pack sizes:

- **Time Complexity**: O(NxM)
- **Space Complexity**: O(N) for the DP table

### Maximum Order Limit

The API restricts orders to a maximum of **10,000,000 items** to prevent excessive memory consumption (DP table size).

For most practical e-commerce scenarios, this limit should be sufficient.

---

## ☁️ Infrastructure & Deployment

The solution is deployed on **Google Cloud** using the $300 free tier allowance. The following components were selected to deploy:

1. **Google Cloud Run** - Serverless container platform hosting the API
2. **Google Artifact Registry** - Docker image repository
3. **GitHub Pages** - Static hosting for the UI (index.html, styles.css)
4. **Terraform** - Infrastructure as Code
5. **GitHub Actions** - CI/CD pipeline

The service is configured with:

- **Min Instances**: 0 → Zero cost when idle
- **Max Instances**: 5 → Automatically scales to handle traffic spikes
- **Resources per Instance**: 2 CPU cores, 4GB memory

## ⚙️ Configuration

### Pack Sizes

Pack sizes are read from the `PACK_SIZES` environment variable, allowing changes **without modifying source code**, but requires redeployment.

Pack sizes are configured in [terraform/variables.tf](terraform/variables.tf#L14). The Terraform passes this to Cloud Run via the `PACK_SIZES` environment variable in [terraform/main.tf](terraform/main.tf#L69).

---

## CI/CD Pipeline

An automated CI/CD pipeline is available.
On every pull request the pipeline will perform docker build and terraform plan.
On merge to main terraform apply is performed.

## 🚀 Future Improvements & Considerations

- **Unit Tests** - Currently not included due to time constraints, but essential for production
- **Docker Image Optimization** - Current image is very small ~32.4 MB; possibly could use `distroless` images to further reduce image size and cold start latency
- **Caching** - Add in-memory cache for frequently requested quantities (popular order sizes). Also we could add caching for the DP table.
- **Authentication & Rate Limiting** - Protect API from abuse
- **Use Workload Identity Federation** - Use WIF instead of Service Account Key for Terraform deployments
- **If staying single-endpoint**: Could migrate core logic to Cloud Functions
- **If expanding features**: Currently changing pack size requirest redeployment, if we wanted to avoid that we could add metadata service for retrieving pack sizes dynamically (trade-off: adds latency)

## 💻 Local Development

### Running the Server

```bash
# Run with custom pack sizes
PACK_SIZES="250,500,1000,2000,5000" go run ./cmd/server
# Server runs on http://localhost:8080
```

### Testing the API

```bash
# Health check
curl "http://localhost:8080/health"

# Pack calculation
curl "http://localhost:8080/packs?quantity=251"
curl "http://localhost:8080/packs?quantity=1501"
```

### Running with Docker

```bash
docker build -t packs-server:latest .

docker run -p 8080:8080 -e PACK_SIZES="250,500,1000,2000,5000" packs-server:latest
```

### Using the CLI Tool

A CLI tool is available for quick testing without the server:

```bash
go run ./cmd/cli 251
go run ./cmd/cli 1501
```

---
