# ADR - Decision to share database servers across instances
Author: Axel Ramone
Creation Date: 2025-10-18

## Problem
The goal of a microservice architecture is to have decoupled, independently deployable applications that ideally have no shared dependencies. This includes database servers and instances.

However, for this project, the cost of deploying both each service along with a separate database instance will far surpass the budget allocated toward it considering its current use case to compute commute times for only 4 users.

## Analysis
A cost estimate for a subset of the services available can be found [here]() for a gcp deployment. The largest cost by far is a postgresql server instance running 24/7. Even with minimal memory and cpu configurations, the cost far outpaces the service's expected usage, costing 2 to 3 times as much to store the data than to run the service itself.

Many of the services' data models can be expressed in either NoSQL or SQL as each service is responsible for a very small domain with few relationships. However, some services, such as the roommate service, do have relationships not only within the service but with other services as well.

The data stored is expected to be small as well and this project will likely not be maintained after a year. Following a true microservice architecture, although a great learning experience, will cost more than it's worth. Dependency management will not come into picture in terms of upgrading a database to a new version.

## Solution
Instead of deploying independent database servers, services will instead share one database instance if the service's data models require normalization and is easier to model with SQL. If the service is better fit for a NoSQL solution, it can deploy its own database instance either through Cloud Firestore or deploy a Cloud Storage Bucket.

Services using the shared PostgreSQL server instance must have their own schemas to prevent conflicts. This will still allow for some benefit to be gained from a microservice architecture as services are only responsible for and can only query from their own schemas.