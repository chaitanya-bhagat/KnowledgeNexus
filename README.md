# KnowledgeNexus
# Enterprise Knowledge & AI Research Platform

A multi-tenant, multi-domain enterprise knowledge platform for securely ingesting organizational documents, performing semantic search, and generating evidence-grounded AI research answers with citations.

> **V1 focuses on the Legal domain while keeping the core platform domain-neutral for future domains such as Finance, HR, Compliance, Healthcare, and Engineering.**

## 1. Problem Statement

Organizations accumulate large amounts of valuable knowledge across documents, policies, reports, judgments, regulations, contracts, manuals, and internal repositories.

This knowledge is often:

- fragmented across teams and systems;
- difficult to discover using exact keyword search;
- expensive and time-consuming to research manually;
- difficult to use safely with generic LLMs;
- difficult to trace back to authoritative evidence;
- sensitive and subject to organizational access boundaries.

Generic LLMs may not have access to an organization's private knowledge and may generate unsupported answers.

This platform addresses that problem by securely ingesting organizational knowledge, making it semantically searchable, retrieving relevant evidence, and generating grounded AI responses with traceable citations.

## 2. What the Platform Does

A customer organization joins the platform as a **Tenant**.

Each tenant can organize knowledge into one or more **Domains** and **Knowledge Bases**.

```text
Tenant
  ↓
Domain
  ↓
Knowledge Base
  ↓
Document
  ↓
Document Version
  ↓
Chunk
```

Users can:

- create and manage knowledge bases;
- upload organizational documents;
- asynchronously process and index documents;
- perform semantic search;
- filter results using domain-specific metadata;
- ask natural-language research questions;
- receive RAG-based answers grounded in retrieved evidence;
- inspect citations;
- save useful research.

## 3. V1 — Legal Domain

V1 implements the **Legal** domain.

Example:

```text
ABC Law Firm
└── Legal
    ├── Supreme Court Judgments
    ├── Contracts
    └── Regulations
```

Legal-specific capabilities include:

- court metadata;
- case numbers;
- judgment dates;
- judges;
- categories;
- Acts and Sections;
- Legal-specific search filters;
- Legal document classification.

The underlying tenant, knowledge, ingestion, retrieval, and AI infrastructure remains reusable for additional domains.

## 4. Core User Flow

```text
Company onboarding
      ↓
Create Tenant
      ↓
Invite Users
      ↓
Create Knowledge Base
      ↓
Upload Document
      ↓
Asynchronous Ingestion
      ↓
Extract
      ↓
Process / Classify
      ↓
Chunk
      ↓
Embed
      ↓
Index
      ↓
READY
      ↓
Semantic Search
      ↓
Ask Research Question
      ↓
Retrieve Evidence
      ↓
Generate Grounded Answer
      ↓
Citations
```

## 5. Architecture

V1 uses:

- **Modular Monolith** for the main Go application.
- **Hexagonal / Ports-and-Adapters Architecture** internally.
- **Constructor-based Dependency Injection**.
- **Separate asynchronous ingestion worker**.

```text
                         Client
                           │
                           ▼
                  Go Modular Monolith
              ┌────────────┼────────────┐
              │            │            │
          Identity       Knowledge     Search
          Tenant         Research       AI
              │            │            │
              └────────────┼────────────┘
                           │
       ┌───────────────────┼──────────────────┐
       ▼                   ▼                  ▼
   PostgreSQL            MinIO              Qdrant
                                               │
                           NATS                │
                            │                  │
                            ▼                  │
                     Ingestion Worker ─────────┘
                            │
                            ▼
                  Embedding Provider

AI Module
    │
    ▼
LLM Provider Interface
    │
    ▼
OpenRouter
```

## 6. Core Modules

```text
identity
tenant
knowledge
ingestion
search
ai
research
domain/legal
```

Responsibilities:

**Identity / Tenant**
Authentication, membership, roles, authorization, and tenant isolation.

**Knowledge**
Domains, knowledge bases, documents, document versions, and metadata.

**Ingestion**
Extraction, Legal processing, chunking, embeddings, indexing, retries, and ingestion lifecycle.

**Search**
Tenant-aware semantic retrieval, filters, Top-K results, scores, and evidence retrieval.

**AI / RAG**
Retrieval orchestration, context selection, prompt construction, LLM invocation, streaming, citation mapping, and usage tracking.

**Research**
Saved searches, bookmarks, conversations, and saved research.

**Domain / Legal**
Legal metadata extraction, classification, taxonomy, filters, and Legal-specific processing.

## 7. Technology Stack

| Area | Technology |
|---|---|
| Backend | Go |
| Primary Database | PostgreSQL |
| Object Storage | MinIO / S3-compatible storage |
| Vector Database | Qdrant |
| Messaging | NATS |
| LLM Gateway | OpenRouter |
| Local Infrastructure | Docker Compose |
| Architecture | Modular Monolith + Hexagonal |
| Dependency Injection | Constructor DI |

## 8. Data Ownership

```text
PostgreSQL
→ authoritative structured application data

MinIO / S3
→ original document files

Qdrant
→ derived chunk embeddings and retrieval index

Redis
→ optional temporary/cache data

NATS
→ asynchronous events/jobs
```

## 9. Multi-Tenancy

Tenant isolation is a core architectural invariant.

```text
Request
  ↓
Authentication
  ↓
User
  ↓
Tenant Membership
  ↓
Authorization
  ↓
Resource Ownership
  ↓
Access
```

Tenant context must propagate through:

- PostgreSQL;
- Qdrant;
- object storage;
- caches;
- events/jobs;
- logs/traces where appropriate.

Cross-tenant access must always be denied.

## 10. Document Ingestion

Documents progress through:

```text
UPLOADED
   ↓
QUEUED
   ↓
EXTRACTING
   ↓
PROCESSING
   ↓
EMBEDDING
   ↓
INDEXING
   ↓
READY

Any processing stage
   ↓
FAILED
```

Ingestion runs asynchronously so large documents do not block user-facing API requests.

## 11. Retrieval and RAG

V1 starts with **vector retrieval**.

```text
Question
   ↓
Query Embedding
   ↓
Vector Search
   ↓
Tenant / Domain / KB Filters
   ↓
Top-K Evidence
```

RAG then uses retrieved evidence:

```text
Question
   ↓
Search
   ↓
Evidence
   ↓
Context Selection
   ↓
LLM
   ↓
Grounded Answer
   ↓
Citations
```

Hybrid retrieval and reranking may be introduced later based on retrieval-quality evaluation.

## 12. Provider Abstractions

External infrastructure is accessed through application-owned interfaces.

```text
VectorStore
    ↑
QdrantAdapter

ObjectStore
    ↑
MinIOAdapter

LLMProvider
    ↑
OpenRouterAdapter

EmbeddingProvider
    ↑
EmbeddingAdapter

MessageBus
    ↑
NATSAdapter
```

This allows providers to be replaced without coupling core business logic to vendor-specific SDKs.

## 13. Local Development

Prerequisites:

- Go
- Docker
- Docker Compose
- OpenRouter/provider credentials

Local infrastructure:

```text
PostgreSQL
MinIO
Qdrant
NATS
```

Typical development flow:

```bash
docker compose up -d
go run ./cmd/api
go run ./cmd/worker
```

Environment-specific secrets must not be committed to source control.

## 14. Observability

Important operational context includes:

```text
request_id
trace_id
tenant_id
document_id
job_id
```

The system should expose:

- structured logs;
- API/search latency metrics;
- ingestion metrics;
- provider failures;
- AI token usage;
- AI cost information where available;
- retrieval metrics;
- distributed traces where appropriate.

## 15. Testing

The project includes or plans to include:

- unit tests;
- integration tests;
- API/event contract tests;
- end-to-end tests;
- tenant-isolation tests;
- retry/idempotency tests;
- failure tests;
- load tests;
- retrieval/RAG quality evaluations.

## 16. Roadmap

### V1

- Multi-tenant platform
- Legal domain
- Knowledge bases
- Document ingestion
- Semantic vector search
- RAG
- Citations
- Research saving
- Tenant/RBAC isolation
- Async ingestion
- Observability foundations

### Future

Potential future capabilities include:

- additional domains;
- hybrid retrieval;
- reranking;
- enterprise connectors;
- knowledge graphs;
- AI agents;
- advanced tenant isolation;
- SSO/SAML;
- additional AI/model providers.