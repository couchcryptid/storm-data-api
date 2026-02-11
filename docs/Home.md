# Storm Data API

A Go service that consumes transformed storm weather reports from a Kafka topic, persists them to PostgreSQL, and serves them through a GraphQL API. Part of the storm data pipeline.

## Pages

- [[Architecture]] -- Project structure, layer responsibilities, database schema, and capacity
- [[Configuration]] -- Environment variables
- [[Deployment]] -- Docker Compose setup and Docker image
- [[Development]] -- Build, test, lint, CI, and project conventions
- [[Data Model]] -- Kafka message shape, event types, field mapping
- [[API Reference]] -- GraphQL types, queries, filter options
