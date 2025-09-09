# VM Placement Service

VM Placement Service API for optimizing virtual machine placement across infrastructure.

## How to Run the Project

### Prerequisites
- Go 1.23+
- Podman and podman-compose

### Steps

1. **Start the database:**
   ```bash
   make deploy-db
   ```

2. **Run the application:**
   ```bash
   make run
   ```

3. **Test the health endpoint:**
   ```bash
   curl -v http://localhost:8080/health
   ```
