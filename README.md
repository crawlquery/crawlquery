# System Overview

This document provides an overview of the different parts of the system, organized into two main components: the API and the Node packages. Each section and sub-section is described to give a clear understanding of the structure and functionality of the system.

## API Packages

The API packages handle the core functionalities, authentication, routing, and business logic of the application.

### ./api

- Root directory for the API packages.

### ./api/migration

- Handles database migrations.

### ./api/domain

- Contains domain models and interfaces used throughout the API.

### ./api/errorutil

- Provides utilities for error handling.

### ./api/shard

- Manages data sharding within the application.
  - **./api/shard/service**
    - Services related to shard management.
  - **./api/shard/repository**
    - Repositories for storing shard-related data.
    - **./api/shard/repository/mysql**
      - MySQL implementation of shard repositories.
    - **./api/shard/repository/mem**
      - In-memory implementation of shard repositories.

### ./api/auth

- Handles authentication and authorization.
  - **./api/auth/service**
    - Services related to authentication.
  - **./api/auth/handler**
    - Handlers for authentication-related endpoints.

### ./api/search

- Manages search functionalities within the application.
  - **./api/search/service**
    - Services related to search operations.
  - **./api/search/handler**
    - Handlers for search-related endpoints.

### ./api/cmd

- Command-line utilities and entry points for the API.
  - **./api/cmd/run**
    - Main entry point to run the API.

### ./api/middleware

- Middleware for request processing and handling.

### ./api/router

- Manages routing of API requests.

### ./api/account

- Handles account management functionalities.
  - **./api/account/service**
    - Services related to account operations.
  - **./api/account/handler**
    - Handlers for account-related endpoints.
  - **./api/account/repository**
    - Repositories for storing account-related data.
    - **./api/account/repository/mysql**
      - MySQL implementation of account repositories.
    - **./api/account/repository/mem**
      - In-memory implementation of account repositories.

### ./api/factory

- Factory pattern implementations for creating service instances.

### ./api/dto

- Data Transfer Objects for API requests and responses.

### ./api/node

- Manages interactions and data processing at the node level.
  - **./api/node/service**
    - Services related to node operations.
  - **./api/node/handler**
    - Handlers for node-related endpoints.
  - **./api/node/repository**
    - Repositories for storing node-related data.
    - **./api/node/repository/mysql**
      - MySQL implementation of node repositories.
    - **./api/node/repository/mem**
      - In-memory implementation of node repositories.

### ./api/link

- Manages link data within the application.
  - **./api/link/service**
    - Services related to link operations.
  - **./api/link/handler**
    - Handlers for link-related endpoints.
  - **./api/link/repository**
    - Repositories for storing link-related data.
    - **./api/link/repository/mysql**
      - MySQL implementation of link repositories.
    - **./api/link/repository/mem**
      - In-memory implementation of link repositories.

### ./api/page

- Manages page data within the application.
  - **./api/page/service**
    - Services related to page operations.
  - **./api/page/repository**
    - Repositories for storing page-related data.
    - **./api/page/repository/mysql**
      - MySQL implementation of page repositories.
    - **./api/page/repository/mem**
      - In-memory implementation of page repositories.

### ./api/crawl

- Manages crawling operations within the application.
  - **./api/crawl/job**
    - Handles crawl job management.
    - **./api/crawl/job/service**
      - Services related to crawl jobs.
    - **./api/crawl/job/handler**
      - Handlers for crawl job-related endpoints.
    - **./api/crawl/job/repository**
      - Repositories for storing crawl job-related data.
      - **./api/crawl/job/repository/mysql**
        - MySQL implementation of crawl job repositories.
      - **./api/crawl/job/repository/mem**
        - In-memory implementation of crawl job repositories.
  - **./api/crawl/restriction**
    - Manages crawl restrictions.
    - **./api/crawl/restriction/service**
      - Services related to crawl restrictions.
    - **./api/crawl/restriction/repository**
      - Repositories for storing crawl restriction-related data.
      - **./api/crawl/restriction/repository/mysql**
        - MySQL implementation of crawl restriction repositories.
      - **./api/crawl/restriction/repository/mem**
        - In-memory implementation of crawl restriction repositories.

## Node Packages

The Node packages handle data processing, indexing, crawling, and related operations at the node level.




### How to start a node
```
go run node/cmd/run/main.go -key 0dc709b5-7565-40e9-b795-d74f5752a6a0 -html /home/ross/cq/node1/html -pdb /home/ross/cq/node1/pdb
```

### ./node

- Root directory for the Node packages.

### ./node/phrase

- Handles phrase extraction and processing.

### ./node/stat

- Manages statistical data processing.
  - **./node/stat/service**
    - Services related to statistical operations.
  - **./node/stat/handler**
    - Handlers for statistical endpoints.

### ./node/signal

- Handles signals for the search algorithm.

### ./node/html

- Manages HTML content processing and storage.
  - **./node/html/service**
    - Services related to HTML operations.
  - **./node/html/repository**
    - Repositories for storing HTML-related data.
    - **./node/html/repository/disk**
      - Disk-based implementation of HTML repositories.
    - **./node/html/repository/mem**
      - In-memory implementation of HTML repositories.

### ./node/domain

- Contains domain models and interfaces used within the Node packages.

### ./node/index

- Handles indexing of data.
  - **./node/index/service**
    - Services related to indexing operations.
  - **./node/index/handler**
    - Handlers for indexing-related endpoints.

### ./node/peer

- Manages peer-to-peer interactions between nodes.
  - **./node/peer/service**
    - Services related to peer operations.

### ./node/filter

- Handles filtering of data.
  - **./node/filter/service**
    - Services related to filtering operations.

### ./node/cmd

- Command-line utilities and entry points for the Node.
  - **./node/cmd/run**
    - Main entry point to run the Node.

### ./node/token

- Handles tokenization of data.

### ./node/quality

- Manages quality assessment of data.

### ./node/router

- Manages routing of Node requests.

### ./node/dump

- Handles data dumping operations.
  - **./node/dump/service**
    - Services related to data dumping.
  - **./node/dump/handler**
    - Handlers for data dump-related endpoints.

### ./node/dto

- Data Transfer Objects for Node requests and responses.

### ./node/page

- Manages page data within the Node.
  - **./node/page/service**
    - Services related to page operations.
  - **./node/page/repository**
    - Repositories for storing page-related data.
    - **./node/page/repository/bolt**
      - BoltDB implementation of page repositories.
    - **./node/page/repository/mem**
      - In-memory implementation of page repositories.

### ./node/crawl

- Manages crawling operations within the Node.
  - **./node/crawl/service**
    - Services related to crawling operations.
  - **./node/crawl/handler**
    - Handlers for crawl-related endpoints.

### ./node/parse

- Handles parsing of data within the Node.

---

This document provides a high-level overview of the system's structure, detailing the various packages and their responsibilities. Each package is designed to handle specific functionalities, ensuring a modular and maintainable codebase.