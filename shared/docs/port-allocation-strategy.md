# 🚀 USC Platform Port Allocation Strategy

## 📋 Overview

To avoid port conflicts and enable flexible deployment scenarios, USC Platform now uses **individual port allocation** for each of the 22 microservices.

## 🎯 Port Allocation Scheme

### Main Service Ports (8001-8022)
Each service gets a unique port based on its service ID:

| Service ID | Service Name | Main Port | Purpose |
|------------|--------------|-----------|---------|
| 01 | Gateway | 8001 | HTTP/REST API |
| 02 | Auth | 8002 | gRPC API |
| 03 | User | 8003 | gRPC API |
| 04 | Blockchain Core | 8004 | gRPC API |
| 05 | Wallet | 8005 | gRPC API |
| 06 | Security | 8006 | gRPC API |
| 07 | Caching | 8007 | gRPC API |
| 08 | Monitoring | 8008 | gRPC API |
| 09 | Social | 8009 | gRPC API |
| 10 | Bilateral Rewards | 8010 | gRPC API |
| 11 | Content Management | 8011 | gRPC API |
| 12 | Video | 8012 | gRPC API |
| 13 | AI | 8013 | gRPC API |
| 14 | Commerce | 8014 | gRPC API |
| 15 | Notification | 8015 | gRPC API |
| 16 | Search | 8016 | gRPC API |
| 17 | Analytics | 8017 | gRPC API |
| 18 | Moderation | 8018 | gRPC API |
| 19 | Recommendation | 8019 | gRPC API |
| 20 | Advertising | 8020 | gRPC API |
| 21 | Admin | 8021 | gRPC API |
| 22 | Kafka Messaging | 8022 | gRPC API |

### Metrics Ports (9001-9022)
Each service exposes Prometheus metrics on a unique port:

| Service | Metrics Port |
|---------|--------------|
| Gateway | 9001 |
| Auth | 9002 |
| User | 9003 |
| ... | ... |
| Kafka Messaging | 9022 |

### Special Purpose Ports

| Port | Service | Purpose |
|------|---------|---------|
| 4000 | Gateway | GraphQL API |
| 8090 | Social | WebSocket (real-time) |
| 8091 | Notification | WebSocket (real-time) |
| 30303 | Blockchain | P2P networking |
| 30301 | Blockchain | Node discovery |
| 7000 | Caching | Redis cluster |
| 1935 | Video | RTMP streaming |

## 💡 Benefits

### 1. **Zero Port Conflicts**
- Each service has dedicated ports
- Can deploy multiple services on same node
- No need for complex port mapping

### 2. **Clear Service Identification**
- Port number directly maps to service ID
- Easy debugging and monitoring
- Intuitive for developers

### 3. **Kubernetes Flexibility**
- Supports NodePort services
- Easy service mesh configuration
- Simplified ingress routing

### 4. **Development Friendly**
- Run multiple services locally
- No port collision during testing
- Clear separation of concerns

## 🔧 Implementation

### Go Constants
```go
// In shared/constants/ports.go
const (
    PortGateway = 8001
    PortAuth    = 8002
    PortUser    = 8003
    // ... etc
)

// Get service port
port := GetServicePort("service-01-gateway") // Returns 8001
```

### Kubernetes Deployment
```yaml
# Example for Auth service
ports:
- name: grpc
  containerPort: 8002  # Unique port
  protocol: TCP
- name: metrics
  containerPort: 9002  # Unique metrics port
  protocol: TCP
```

### Service Discovery
```yaml
# Service definition
spec:
  type: ClusterIP
  ports:
  - name: grpc
    port: 8002
    targetPort: 8002
    protocol: TCP
    appProtocol: h2c
```

## 🚀 Migration Guide

### 1. Update Constants
```bash
# Constants are already updated in shared/constants/ports.go
go build ./shared/constants  # Verify no errors
```

### 2. Update K8s Deployments
```bash
# Run the automated update script
chmod +x shared/scripts/update-k8s-ports.sh
./shared/scripts/update-k8s-ports.sh
```

### 3. Update Service Code
```go
// Old way
port := 8080  // Hardcoded

// New way  
import "github.com/usc-platform/shared/constants"
port := constants.GetServicePort(constants.ServiceAuth)  // 8002
```

### 4. Update Health Checks
```yaml
# Update readiness probes to use new ports
readinessProbe:
  grpc:
    port: 8002  # Service-specific port
```

### 5. Update Monitoring
```yaml
# Update Prometheus scrape configs
annotations:
  prometheus.io/port: "9002"  # Service-specific metrics port
```

## 🔍 Validation

### Test Port Allocation
```bash
cd shared/constants
go test -v  # Validates no port conflicts
```

### Verify K8s Updates  
```bash
cd shared/examples/k8s-clean/06-services/2-services
grep -r "containerPort:" . | sort  # Check port distribution
```

### Port Allocation Report
```go
// Generate port allocation table
fmt.Println(constants.PrintPortAllocation())
```

## 🎯 Best Practices

### 1. **Use Constants**
Always import and use `shared/constants` instead of hardcoding ports.

### 2. **Environment Variables**
```dockerfile
ENV SERVICE_PORT=8002
ENV METRICS_PORT=9002
```

### 3. **Health Checks**
Configure health checks to use service-specific ports.

### 4. **Service Mesh**
Update Istio/service mesh configurations with new ports.

### 5. **Monitoring**
Update Prometheus, Grafana dashboards with new metrics ports.

## 🔄 Rollback Strategy

If issues arise, restore from backups:
```bash
cd shared/examples/k8s-clean/06-services/2-services
for f in *.backup; do mv "$f" "${f%.backup}"; done
```

## 📊 Port Range Summary

| Range | Purpose | Count |
|-------|---------|-------|
| 8001-8022 | Main service APIs | 22 ports |
| 9001-9022 | Metrics endpoints | 22 ports |
| 4000 | GraphQL Gateway | 1 port |
| 8090-8091 | WebSocket services | 2 ports |
| 30301, 30303 | Blockchain P2P | 2 ports |
| 7000 | Redis clustering | 1 port |
| 1935 | RTMP streaming | 1 port |

**Total: 51 unique ports allocated** ✅

## 🎉 Benefits Achieved

✅ **Zero port conflicts** - Each service has unique ports  
✅ **Kubernetes ready** - Support for NodePort and service mesh  
✅ **Developer friendly** - Clear port-to-service mapping  
✅ **Production ready** - Scalable and maintainable  
✅ **Monitoring ready** - Unique metrics ports for each service  
✅ **Future proof** - Easy to add new services
