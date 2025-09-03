# Info:
to run the test use:
````
go test ./...
````
to test the new endpoint use:
````
curl --location 'http://localhost:8080/beer/getFiltered?includeIpa=true&year=2000&hasFood=wolf&abvSortOrder=asc'
````
cache and mock api rate limits parameters can be adjusted in the config file.
# Interview Go — Candidate Task

Welcome! This repo is a minimal skeleton of an HTTP service in Go (Echo) that you will extend in ~60–90 minutes.
The goal of the exercise is to see how you structure code, validate input, expose REST endpoints etc...

> **Important:** Keep the solution **simple**. You don’t need a database, messaging, or external tools. Use what is already here.
> You can write your code and docs in **English**

---

## 1) Requirements

- Go **1.21+**
- Git
- IntelliJ IDEA / GoLand / VS Code
- cURL or Postman/Insomnia
---

## 2) How to run

```bash
go run .
# or from IntelliJ/GoLand: Run ▶️ on main.go
```

Default config is in `./config.yaml`.  
You can change the port in the file or override with environment variables (if mapped in code).

---

## 3) What’s already here

- Bootstrapping in `main.go` (flag `-config`, start server, OS signals, graceful shutdown).
- Endpoint: `GET /health` → `200 {"status":"ok"}`
- Endpoint: `GET /beer/getAll` → `get a list with 50 beers`

Health check:

```bash
curl -i http://localhost:8080/health
# 200 {"status":"ok"}
```

```bash
curl -s http://localhost:8080/beer/getAll | jq '.[0]'
# 200 {"status":"ok"}
```

---

## 4) Your assignment

Implement a **Beer Endpoint with a few clear requirements**

### Story

**Return all IPA beers brewed after year 2015, that go with chicken as food pairing, sorted in DESC order by alcohol content (ABV).**

You will:
1.	Design a middleware that configures default filtering/sorting parameters (and allows overriding them).
2.	Implement the endpoint that uses those parameters to filter/sort the beer list.
3.	Add in-memory caching of request → response. 
4.	Handle “API rate limit” errors from the downstream client (simulated).
5.	Write unit and integration tests.