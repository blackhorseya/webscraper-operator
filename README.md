# WebScraper Operator

## Overview

The **WebScraper Operator** is a Kubernetes custom operator designed to manage and orchestrate large-scale web scraping
jobs using `CronJob` resources.
It simplifies the scheduling, monitoring, and management of web scraping tasks while
ensuring scalability and reliability in a Kubernetes-native environment.

## Features

- **Dynamic CronJob Management**: Automatically creates and manages `CronJob` resources for web scraping tasks based on
  custom resource definitions.
- **Retry Mechanism**: Handles failures with configurable retry policies.
- **Concurrency Control**: Prevents overlapping job executions with concurrency policies.
- **Resource Optimization**: Allows fine-grained control over resource requests and limits for each job.
- **Failure Notifications**: Sends notifications to a webhook when jobs fail or require manual intervention.
- **Scalability**: Supports thousands of scraping jobs with efficient scheduling and resource management.
