# 🔒 USC Platform Shared Library - Security Guide

## 🚨 **CRITICAL SECURITY REQUIREMENTS**

### **1. Environment Variables Setup**

**NEVER** use hardcoded secrets in your code or configuration files. Always use environment variables:

```bash
# Copy the example file
cp examples/k8s/env.example .env

# Edit with your actual values
nano .env
```

### **2. Required Environment Variables**

The following environment variables **MUST** be set for production:

#### **Authentication**
```bash
JWT_SECRET=your-super-secure-jwt-secret-key-must-be-at-least-32-characters-long-for-production
```

#### **Database Passwords**
```bash
POSTGRES_PASSWORD=your-secure-postgres-password
REDIS_PASSWORD=your-secure-redis-password
CLICKHOUSE_PASSWORD=your-secure-clickhouse-password
QUICKWIT_PASSWORD=your-secure-quickwit-password
INFLUXDB_TOKEN=your-secure-influxdb-token
```

#### **Kafka Security**
```bash
KAFKA_SASL_USERNAME=your-kafka-username
KAFKA_SASL_PASSWORD=your-secure-kafka-password
```

### **3. Security Best Practices**

#### **Password Requirements**
- **Minimum 16 characters** for database passwords
- **Minimum 32 characters** for JWT secrets
- Use **cryptographically secure random strings**
- **Unique passwords** for each service and environment

#### **JWT Secret Generation**
```bash
# Generate a secure JWT secret (64 characters)
openssl rand -base64 48

# Or using Go
go run -c 'import "crypto/rand"; import "encoding/base64"; b := make([]byte, 32); rand.Read(b); print(base64.StdEncoding.EncodeToString(b))'
```

#### **Database Password Generation**
```bash
# Generate secure passwords
openssl rand -base64 32
```

### **4. Environment-Specific Configuration**

#### **Development**
```bash
ENVIRONMENT=development
POSTGRES_PASSWORD=dev-password-123
JWT_SECRET=dev-jwt-secret-key-32-chars-minimum
```

#### **Staging**
```bash
ENVIRONMENT=staging
POSTGRES_PASSWORD=staging-secure-password-456
JWT_SECRET=staging-jwt-secret-key-32-chars-minimum
```

#### **Production**
```bash
ENVIRONMENT=production
POSTGRES_PASSWORD=prod-ultra-secure-password-789
JWT_SECRET=prod-jwt-secret-key-32-chars-minimum
```

### **5. Secret Management**

#### **Option 1: Environment Variables**
```bash
# Set in your shell
export JWT_SECRET="your-secret-here"
export POSTGRES_PASSWORD="your-password-here"

# Or use .env file (never commit this!)
echo "JWT_SECRET=your-secret-here" >> .env
```

#### **Option 2: Docker Secrets**
```yaml
# docker-compose.yml
services:
  your-service:
    environment:
      - JWT_SECRET_FILE=/run/secrets/jwt_secret
    secrets:
      - jwt_secret

secrets:
  jwt_secret:
    file: ./secrets/jwt_secret.txt
```

#### **Option 3: Kubernetes Secrets**
```yaml
# k8s-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: usc-secrets
type: Opaque
data:
  jwt-secret: <base64-encoded-secret>
  postgres-password: <base64-encoded-password>
```

#### **Option 4: HashiCorp Vault**
```bash
# Install Vault CLI
vault kv put secret/usc-platform/jwt-secret value="your-secret-here"
vault kv put secret/usc-platform/postgres-password value="your-password-here"
```

### **6. Configuration Validation**

The shared library automatically validates security requirements:

```go
// JWT Secret validation
if len(cfg.Auth.JWTSecret) < 32 {
    return fmt.Errorf("jwt secret must be at least 32 characters")
}

// Empty secret validation
if cfg.Auth.JWTSecret == "" {
    return fmt.Errorf("jwt secret cannot be empty")
}
```

### **7. SSL/TLS Configuration**

#### **Production Database Connections**
```yaml
database:
  sslmode: "require"
  sslcert: "/etc/ssl/certs/postgres.crt"
  sslkey: "/etc/ssl/certs/postgres.key"
  sslrootcert: "/etc/ssl/certs/postgres.crt"
```

#### **Kafka SSL Configuration**
```yaml
kafka:
  security_protocol: "SSL"
  ssl_ca_file: "/etc/ssl/certs/kafka-ca.crt"
  ssl_cert_file: "/etc/ssl/certs/kafka-client.crt"
  ssl_key_file: "/etc/ssl/certs/kafka-client.key"
```

### **8. Security Checklist**

Before deploying to production, ensure:

- [ ] All hardcoded secrets removed from code
- [ ] Environment variables properly set
- [ ] JWT secret is at least 32 characters
- [ ] Database passwords are strong and unique
- [ ] SSL/TLS enabled for all connections
- [ ] Secrets are not committed to version control
- [ ] Production secrets are different from development
- [ ] Secret rotation plan is in place
- [ ] Access to secrets is properly restricted
- [ ] Monitoring and alerting for security events

### **9. Common Security Mistakes**

#### **❌ DON'T DO THIS**
```go
// Hardcoded secret
cfg.Auth.JWTSecret = "my-secret-key"

// Weak default password
cfg.Database.Password = "password"

// Same secret for all environments
cfg.Auth.JWTSecret = "same-secret-everywhere"
```

#### **✅ DO THIS INSTEAD**
```go
// Use environment variable
cfg.Auth.JWTSecret = os.Getenv("JWT_SECRET")

// No default password
cfg.Database.Password = os.Getenv("POSTGRES_PASSWORD")

// Environment-specific secrets
cfg.Auth.JWTSecret = os.Getenv("JWT_SECRET")
```

### **10. Emergency Response**

If secrets are compromised:

1. **Immediately rotate** all affected secrets
2. **Revoke** any issued JWT tokens
3. **Update** all environment configurations
4. **Restart** all services
5. **Audit** access logs for suspicious activity
6. **Notify** security team and stakeholders

### **11. Security Monitoring**

Monitor for:
- Failed authentication attempts
- Unusual access patterns
- Configuration changes
- Secret access attempts
- SSL/TLS certificate expiration

---

## 🛡️ **SECURITY CONTACT**

For security issues or questions:
- **Email**: security@usc-platform.com
- **Documentation**: [Security Wiki](https://wiki.usc-platform.com/security)
- **Incident Response**: [Security Playbook](https://wiki.usc-platform.com/security/incident-response)

---

**Remember: Security is everyone's responsibility. When in doubt, ask the security team!**
