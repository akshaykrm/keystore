# Keystore Platform

> A multi-tenant keystore platform built with Go, React and PostgreSQL.

## Overview

This project is a long-term effort to build a modern keystore platform capable of supporting businesses ranging from individual merchants to larger workspaces.

The platform will provide the infrastructure for managing products, customers, orders, users, permissions, integrations, and other commerce-related functionality. Customer-facing storefronts will be separate applications that consume the platform's API.

The focus of this project is on building a solid domain model and maintainable architecture before optimising for features or scale.

---

## Goals

- Build an API-first keystore platform.
- Support multi-tenancy through workspaces.
- Allow users to belong to multiple workspaces.
- Provide a flexible role and permission system.
- Design with extensibility in mind for third-party integrations.
- Maintain a clean, modular architecture that can evolve over time.

---

## Planned Technology Stack

### Backend

- Go
- PostgreSQL
- REST API
- JWT Authentication

### Frontend

- React
- TypeScript
- Tailwind CSS

### Infrastructure

- Docker (after MVP)
- Redis (when needed)
- Object Storage (S3-compatible)

---

## Core Principles

- Domain-driven design over framework-driven development.
- API-first architecture.
- Single source of truth for both code and documentation.
- Simple solutions before complex abstractions.
- Build for maintainability rather than short-term speed.
- Add complexity only when justified.

---

## High-Level Roadmap

### Phase 1 – Foundation

- Users
- Authentication
- Workspaces
- Memberships
- Roles & Permissions

### Phase 2 – Catalogue

- Categories
- Products
- Inventory
- Media

### Phase 3 – Customers

- Customers
- Addresses

### Phase 4 – Commerce

- Orders
- Payments
- Shipping

### Phase 5 – Platform

- Integrations
- Notifications
- Webhooks
- Analytics

---

## Repository Documentation

The `docs/` directory contains the project documentation.

- `ROADMAP.md` — High-level roadmap and milestones.
- `CURRENT.md` — Current focus and next steps.
- `domain/` — Business concepts and domain modelling.

---

## Project Philosophy

This project is intentionally developed without fixed deadlines.

The priority is to produce a well-designed platform that can evolve over time. Progress is measured by improving the architecture, domain model, and overall quality of the system rather than by implementing features as quickly as possible.

---

## Status

**Current Phase:** Phase 1
