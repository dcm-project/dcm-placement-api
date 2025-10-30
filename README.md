# VM Placement Service

VM Placement Service API for optimizing virtual machine placement across infrastructure.

## How to Run the Project

### Prerequisites
- Go 1.23+
- Podman
- Cluster with KubeVirt - Find more information [here](https://kubevirt.io/quickstart_kind/)

### Steps
1. ** Login to openshift/k8s with CNV and create namespaces **
   ```bash
   oc login ...
   ```

2. **Run the application:**
   ```bash
   cat > deploy/podman/.env << 'EOF'
   # Database Configuration
   DATABASE_NAME=placement
   DATABASE_USER=admin
   DATABASE_PASSWORD=adminpass
   DATABASE_PORT=5432
   DATABASE_HOST=placement-db
   DATABASE_MASTER_USER=admin
   DATABASE_MASTER_PASSWORD=adminpass

   # Kubernetes configuration file path (used by k8s-service-provider)
   # Default: ${HOME}/.kube/config
   KUBECONFIG=${HOME}/.kube/config
   EOF
   make compose-up
   ```

3. **Create app:**
   ```bash
   curl -v -X POST -H "Content-type: application/json" --data '{"name": "myvm", "service": "webserver", "tier": 1}'  http://localhost:8080/applications
   ```

4. **Check VMs:**
   ```bash
   oc get vm -n us-east-1
   oc get vm -n us-east-2
   ```
