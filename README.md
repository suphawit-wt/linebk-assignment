# 🏦 LINE BK - Software Engineer Assignment

A Banking service delivering high-concurrency performance (capable of **3,700+ RPS** with **<400ms p95 latency** and **0.00% error rate**), secure authentication, JWT-based role authorization, and core banking operations. Built with strict code quality standards, resilient caching architecture, and full container orchestration.

---

## 🌟 Overview & Architecture

* **Performance Benchmark:** Successfully validated under aggressive stress testing to process **1,024,800+ requests in under 5 minutes** (4m 30s) at a sustained throughput of **3,795 requests/sec**, maintaining an elite **p95 latency of ~390ms** and a **0.00% error rate** across millions of operations.
* **Core Engine:** Written in **Go (Golang)** using clean hexagonal architecture patterns, explicit error handling, high-performance connection pooling, and strict linting compliance (`golangci-lint`).
* **Storage & Caching:** Uses **MySQL** for primary relational storage paired with **Redis cache** to optimize high-volume core read endpoints and protect database threads from connection starvation.
* **Security & Auth:** Secures endpoint routes via RSA-signed JWT tokens and Argon2id password hashing, with configuration secrets initialized through **HashiCorp Vault**.
* **Containerization:** Deployed using multi-stage Docker builds into minimal **Distroless** static images running under non-root security contexts.

---

## 🚀 Setup & Getting Started

To initialize the project environment on a freshly cloned repository, run:

```bash
make init
```

Once initialized, spin up the full stack containers in the background by running:
```bash
make up
```

### Accessing the Application
* **React Frontend:** Open your browser and navigate to **[http://localhost:3000](http://localhost:3000)** to access the LINE BK web interface. (Passcode/PIN: 123456)
* **API / Swagger UI:** Available at **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**.

### Prerequisites
* **Homebrew** (`brew`) — *(used during `make init` to install `k6` if missing)*
* **Go**
* **Docker & Docker Compose**
* **Node.js & pnpm** (for local hooks like Lefthook)
* **OpenSSL** (for automated RSA key pair generation)

---

## 🛠️ Makefile Command Reference

Below is the summary of available commands in the Makefile:

| Command | Category | Description |
| :--- | :--- | :--- |
| `make init` | Setup | Initializes dependencies, configures `.env` files, generates Vault JWT RSA keys, and installs `k6` via Homebrew if needed. |
| `make dev-up` | Development | Starts development infrastructure containers in the background. |
| `make dev-down` | Development | Stops running development infrastructure containers. |
| `make dev-clean` | Development | Stops dev containers and removes volumes, orphans, and local dev images. |
| `make up` | Production Stack | Starts the full production-grade stack containers in the background. |
| `make down` | Production Stack | Stops full production-grade stack containers. |
| `make clean` | Production Stack | Stops production containers and cleans volumes, orphans, and local build images. |
| `make lint` | Quality | Executes `golangci-lint` with auto-fix enabled across the Go codebase. |
| `make test` | Quality | Runs Go unit tests with coverage metrics enabled. |
| `make load-test` | Performance | Runs k6 load test script (`load_test.js`). |
| `make stress-test` | Performance | Runs k6 stress test script (`stress_test.js`). |
| `make swag` | Documentation | Re-generates OpenAPI Swagger documentation files. |

---

## 📚 API Documentation & Swagger UI

Interactive OpenAPI/Swagger documentation is available when the service is running.

**Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Key Endpoints Overview

#### 🔑 Authentication (`/v1/auth`)
* `POST /v1/auth/signin/password` — Authenticate via ID and password to receive JWT token pair.
* `POST /v1/auth/signin/passcode` — Authenticate via passcode.
* `POST /v1/auth/refresh` — Issue a new access token using a valid HTTP-only refresh token cookie.
* `POST /v1/auth/signout` — Invalidate current session and clear authentication cookies.

#### 💳 Core Banking (`/v1`)
* `GET /v1/user` — Fetch user profile info and greeting *(Requires Bearer Token)*.
* `GET /v1/accounts` — Retrieve user account list and balances *(Requires Bearer Token)*.
* `GET /v1/debit-cards` — Retrieve debit cards and card status *(Requires Bearer Token)*.
* `GET /v1/transactions/recent` — Fetch recent ledger transaction history *(Requires Bearer Token)*.
* `GET /v1/banners` — Fetch active app display banners *(Requires Bearer Token)*.