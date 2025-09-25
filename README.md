# VM Placement Service

VM Placement Service API for optimizing virtual machine placement across infrastructure.

## How to Run the Project

### Prerequisites
- Go 1.23+
- Podman
- Cluster with KubeVirt - Find more information [here](https://kubevirt.io/quickstart_kind/)

### Steps
0. ** Login to openshift/k8s with CNV and create namespaces **
   ```bash
   oc login ...
   oc create ns us-east-1
   oc create ns us-east-2
   oc create ns us-west-1
   oc create ns us-west-2
   ```

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

5. **Create app:**
   ```bash
   curl -v -X POST -H "Content-type: application/json" --data '{"name": "myvm", "service": "webserver", "tier": "1"}'  http://localhost:8080/applications
   ```

6. **Check VMs:**
   ```bash
   oc get vm -n us-east-1
   oc get vm -n us-east-2
   ```
