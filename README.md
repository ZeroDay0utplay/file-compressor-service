# File Compression Service

A lightweight service for general‑purpose file compression. Built with clean architecture and SOLID principles, it isolates CPU‑heavy work from the main app, supports controlled concurrency, low memory usage, and returns compressed files with no persistent storage.

---

## Features

* General‑purpose compression (files of any type)
* Multiple algorithms (e.g., zstd, gzip, brotli – configurable)
* Ephemeral file handling (no persistent storage)
* Designed for performance and low memory usage
* Controlled concurrency to protect CPU
* Clean architecture and SOLID principles
* Container‑friendly

---

## Architecture Overview

```
Client (Rails / API)
      |
      v
File Compression Service
      |
      v
Compressed File (binary response)
```

* Receives file as binary or multipart upload
* Streams to temp file
* Compresses with selected algorithm
* Reads compressed file
* Deletes temp files
* Returns binary response

---

## Tech Stack

* Language: Go (or your chosen runtime)
* Compression Engines: zstd, gzip, brotli (pluggable)
* Deployment: Docker (Fly.io / Render / VPS)
* CI/CD: GitLab CI

---

## API Contract

### Compress File

`POST /compress`

**Request**

* Content-Type: multipart/form-data
* Field: `file` (any file)
* Optional field: `algorithm` (zstd | gzip | brotli)

**Response**

* Content-Type: application/octet-stream
* Body: compressed file binary

**Errors**

* 400: invalid file or missing file
* 413: file too large
* 500: compression failure

---

## Compression Strategy

Default algorithm can be chosen based on file size:

| Size   | Algorithm               |
| ------ | ----------------------- |
| < 5MB  | gzip                    |
| 5–50MB | zstd                    |
| > 50MB | zstd (high compression) |

---

## Concurrency Model

* Requests are accepted concurrently
* Actual compression is limited by a worker pool
* Max workers should match available CPU cores

Example:

* 2 vCPU → max 1–2 concurrent compressions

---

## Memory Strategy

* Stream uploads to disk
* Avoid keeping large files in memory
* Temp files deleted immediately after request
* Recommended memory: 512MB–1GB depending on file size

---

## Local Development

### Requirements

* Go 1.22+ (or chosen runtime)
* Compression libs/tools
* Docker (optional)

### Run locally

```bash
go run main.go
```

---

## Docker

### Build

```bash
docker build -t file-compression-service .
```

### Run

```bash
docker run -p 3000:3000 file-compression-service
```

---

## Deployment

### Vercel (Primary Target)

> Note: Vercel is optimized for lightweight, short‑running workloads. For CPU‑heavy or very large files, apply strict limits.

**Guidelines**

* Prefer lightweight algorithms (gzip, brotli, low/medium zstd levels)
* Avoid very large files (recommend < 20–30MB)
* Keep execution time short
* Use streaming where possible
* No persistent storage — rely on ephemeral temp files

**Setup**

* Use Vercel Serverless Functions (Node or Go runtime if supported)
* Configure max memory/time within Vercel limits
* Enforce file size limits at the edge

**Hybrid Option**

* Small files → compress on

## CI/CD (GitLab)

* Tests on every push
* Builds Docker image
* Pushes to registry
* Auto‑deploy on main branch

---

## Security

* Limit file size
* Validate input
* Rate limit endpoint
* Optional auth token

---

## When to Use

* When compression is not your core feature
* When you need to protect your main app from CPU spikes
* When files can be large or numerous

---

## License

MIT or your preferred license
