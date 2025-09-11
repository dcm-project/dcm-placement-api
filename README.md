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

3. **Run the OPA server:**
   ```bash
   make opa
   ```

4. **Run the application:**
   ```bash
   make run
   ```

5. **Create a VM:**
   ```bash
   # Not allowed:
   curl -X POST -H "Content-type: application/json" --data '{"name": "myvm", "env": "PROD", "region": "us-east-2", "ram": 1, "os": "RHEL", "cpu": 2, "role": "public-facing", "tenantid": "PRCR-001"}'  http://localhost:8080/place/vm
   # Allowed
   curl -X POST -H "Content-type: application/json" --data '{"name": "myvm", "env": "PROD", "region": "us-east-1", "ram": 1, "os": "RHEL", "cpu": 2, "role": "public-facing", "tenantid": "PRCR-001"}'  http://localhost:8080/place/vm
   ```

5. **Get a VMs:**
   ```bash
   curl http://localhost:8080/declaredvms
   curl http://localhost:8080/requestedvms
   ```
