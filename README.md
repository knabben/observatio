# Observātiō 

[![Build](https://github.com/knabben/observatio/actions/workflows/build.yml/badge.svg)](https://github.com/knabben/observatio/actions/workflows/build.yml)

<p align="center">
<img src="front/public/logo.png" alt="logo" width="300"/>
</p>

The project focuses on monitoring Kubernetes clusters managed ty [ClusterAPI](https://cluster-api.sigs.k8s.io/), 
providing tools and solutions to enhance visibility and efficiency. By collecting and consolidating data from diverse sources, 
it offers comprehensive insights into cluster performance and health. Equipped with advanced dashboards and real-time visualization, 
the project enables users to swiftly identify and address issues, improving operational reliability and reducing downtime. 
This solution empowers organizations to maintain optimal cluster functionality, streamline troubleshooting efforts, 
and ensure robust management of their cloud-native environments.

## Development

### Prerequisites

- Go 1.23.1
- Node.js and pnpm
- Linux and Make

### Backend Setup

1. Install backend dependencies:

   ```bash
   cd webserver
   go mod tidy
   ```

2. Build the backend webserver job:
   ```bash
   make run-backend what=serve
   ```

3. Running unit tests
   ```bash
   make run-tests-backend
   ```

The backend server will start and listen for WebSocket connections. By default, it runs on port 8080.

### Frontend Setup

1. Install frontend dependencies:
   ```bash
   cd front
   pnpm install
   ```

2. Run the development server:
   make run-frontend
   ```bash
   ```

3. Run tests for the frontend:
   ```bash
   make run-tests-frontend
 

The frontend development server will start and be available at http://localhost:3000.
