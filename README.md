# Ministry of Education Portal — Distributed System (full stack)

A fault-tolerant **distributed system in Go** that serves Ethiopia's Grade-12 University Entrance Exam results at national scale. This is the complete monorepo: a geo-aware load balancer, **replicated** auth and backend services, a petition service, MySQL, and a ReactJS frontend.

![MoE portal](docs/screenshots/home.png)

## The problem

Each year 500,000+ students check their University Entrance Exam results in a short window. A single server buckles under that spike. This system distributes the load across replicated, coordinated services so results stay available under heavy traffic.

## Architecture

```
                  ┌─────────────────┐
   clients  ────▶ │  Load Balancer  │  geo-aware selection + etcd coordination
                  └────────┬────────┘
        ┌──────────────────┼──────────────────┐
   ┌────▼─────┐       ┌─────▼─────┐       ┌────▼─────┐
   │ backend  │       │   auth    │       │ petition │
   │ server_1 │       │ server_1  │       │   1 / 2  │
   │ server_2 │       │ server_2  │       └──────────┘
   └────┬─────┘       └───────────┘
        ▼                                   ┌──────────┐
     ┌───────┐                              │ frontend │  ReactJS
     │ MySQL │                              └──────────┘
     └───────┘
```

- **Load balancer** (`load_balancer/LoadBalancer.go`) — distributes requests across replicas with **geo-aware server selection** and **etcd**-based distributed locking/coordination (lease grants + compare-and-swap).
- **Auth servers** (`auth/server_1`, `auth/server_2`) — replicated; authenticate and authorize admins. Only authorized MoE staff can upload results. The backend calls auth over **RPC**.
- **Backend servers** (`backend/server_1`, `backend/server_2`) — replicated; serve student result lookups (open) and accept authorized uploads.
- **Petition service** (`backend/petition1`, `backend/petition2`) — handles student result petitions.
- **MySQL** — durable storage for results.
- **Frontend** (`frontend/`) — ReactJS UI for students.

## Tech stack

**Go** · **etcd** · **MySQL** · RPC · ReactJS · Docker

## Running locally

Each Go service is its own module (see `go.work`). Start etcd and MySQL, then bring up the load balancer, auth servers, and backend servers (see `start.sh`). Run the frontend with `cd frontend && npm install && npm run dev`.

> Team project (forked from the shared team repo). My work focused on the load balancer, etcd-based coordination, and the backend/auth RPC path.
