# MinIO Sample Application with Chainguard Containers

This project demonstrates how to use MinIO with Chainguard's secure, minimal container images. It includes a Python sample application that performs common MinIO operations, along with deployment configurations for both Docker Compose and Kubernetes.

## Overview

This sample application showcases:
- MinIO server running on Chainguard's minimal container image
- Python client application demonstrating bucket and object operations
- Docker Compose setup for local development
- Kubernetes manifests for production deployment
- Multi-stage Docker builds using Chainguard Python images

## Prerequisites

### For Docker Compose
- Docker and Docker Compose
- Access to `cgr.dev/chainguard/minio` and `cgr.dev/chainguard/python` images
  - **OR** use the Docker Hub alternatives (see [Using Docker Hub Images](#using-docker-hub-images))

### For Kubernetes
- Kubernetes cluster 
- kubectl CLI tool
- Access to Chainguard container registry
  - **OR** use the Docker Hub alternatives (see [Using Docker Hub Images](#using-docker-hub-images))

## Project Structure

```
.
├── README.md
├── .gitignore
├── docker-compose.yml              # Local development (Chainguard images)
├── docker-compose.dockerhub.yml    # Local development (Docker Hub alternative)
├── app/
│   ├── Dockerfile                 # Python app (Chainguard images)
│   ├── Dockerfile.dockerhub       # Python app (Docker Hub alternative)
│   ├── main.py                    # Sample MinIO client application
│   ├── requirements.txt           # Python dependencies
│   └── .dockerignore
├── k8s/                            # Kubernetes manifests (Chainguard images)
│   ├── namespace.yaml
│   ├── secret.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── test-job.yaml
│   └── kustomization.yaml
└── k8s-dockerhub/                  # Kubernetes manifests (Docker Hub alternative)
    ├── namespace.yaml
    ├── secret.yaml
    ├── deployment.yaml
    ├── service.yaml
    ├── test-job.yaml
    └── kustomization.yaml
```

## Quick Start with Docker Compose

1. **Start the services:**
   ```bash
   docker-compose up -d
   ```

2. **View the logs:**
   ```bash
   docker-compose logs -f
   ```

3. **Access MinIO:**
   - API: http://localhost:9000
   - Console: http://localhost:9001
   - Credentials: `minioadmin` / `minioadmin123`

4. **Run the sample application:**
   ```bash
   docker-compose up app
   ```

5. **Stop the services:**
   ```bash
   docker-compose down
   ```

## Kubernetes Deployment

### Deploy with kubectl

1. **Apply all manifests:**
   ```bash
   kubectl apply -f k8s/namespace.yaml
   kubectl apply -f k8s/secret.yaml
   kubectl apply -f k8s/deployment.yaml
   kubectl apply -f k8s/service.yaml
   ```

2. **Or use Kustomize:**
   ```bash
   kubectl apply -k k8s/
   ```

3. **Check deployment status:**
   ```bash
   kubectl get pods -n minio-demo
   kubectl get svc -n minio-demo
   ```

4. **Port-forward to access MinIO:**
   ```bash
   # API access
   kubectl port-forward -n minio-demo svc/minio 9000:9000

   # Console access
   kubectl port-forward -n minio-demo svc/minio 9001:9001
   ```

5. **Run the test job (optional):**
   ```bash
   kubectl apply -f k8s/test-job.yaml
   kubectl logs -n minio-demo job/minio-test
   ```

## Sample Application Features

The Python application (`app/main.py`) demonstrates:

1. **Connecting to MinIO** - Establishes secure connection with retry logic
2. **Creating buckets** - Creates a sample bucket if it doesn't exist
3. **Uploading objects** - Uploads text and JSON data
4. **Listing objects** - Displays all objects in a bucket
5. **Downloading objects** - Retrieves and displays object contents
6. **Deleting objects** - Removes objects from storage
7. **Listing buckets** - Shows all available buckets

### Running the Sample App Locally

1. **Install dependencies:**
   ```bash
   cd app
   pip install -r requirements.txt
   ```

2. **Set environment variables:**
   ```bash
   export MINIO_ENDPOINT=localhost:9000
   export MINIO_ACCESS_KEY=minioadmin
   export MINIO_SECRET_KEY=minioadmin123
   export MINIO_SECURE=false
   ```

3. **Run the application:**
   ```bash
   python main.py
   ```

## Configuration

### Environment Variables

#### MinIO Server
- `MINIO_ROOT_USER` - Admin username (default: minioadmin)
- `MINIO_ROOT_PASSWORD` - Admin password (default: minioadmin123)

#### Python Application
- `MINIO_ENDPOINT` - MinIO server endpoint (default: localhost:9000)
- `MINIO_ACCESS_KEY` - Access key for authentication
- `MINIO_SECRET_KEY` - Secret key for authentication
- `MINIO_SECURE` - Use TLS (true/false, default: false)

## Extending the MinIO Image

This demo uses the Chainguard MinIO image directly without building a custom Dockerfile. The image is designed to be used as-is:

```yaml
# In docker-compose.yml
minio:
  image: cgr.dev/chainguard/minio:latest  # Pull and run directly
```

### When You Need Customization

If you need to add custom configurations, certificates, or policies, create a Dockerfile:

```dockerfile
FROM cgr.dev/chainguard/minio:latest

# Example: Add custom CA certificates
COPY custom-ca-cert.crt /etc/ssl/certs/

# Example: Add MinIO policy files
COPY policies/ /policies/
```

Then update docker-compose.yml to build instead of pull:

```yaml
minio:
  build:
    context: .
    dockerfile: Dockerfile
  # ... rest of config
```

**For most use cases**, use the image directly as shown in this demo. The minimal approach reduces attack surface and simplifies deployment.

## Security Considerations

This sample uses default credentials for demonstration purposes. For production:

1. **Use strong credentials** - Generate secure random passwords
2. **Use Kubernetes Secrets** - Store credentials securely
3. **Enable TLS** - Configure HTTPS for MinIO
4. **Network policies** - Restrict access to MinIO pods
5. **Use persistent volumes** - Configure proper storage for production
6. **Update the secret.yaml** - Use proper secret management tools like:
   - Sealed Secrets
   - External Secrets Operator
   - HashiCorp Vault
   - Cloud provider secret managers (AWS Secrets Manager, GCP Secret Manager, Azure Key Vault)

## Why Chainguard Containers?

Chainguard containers offer several advantages:

- **Minimal attack surface** - Only essential components included
- **No CVEs** - Regularly updated to minimize vulnerabilities
- **Non-root by default** - Enhanced security posture
- **Small image size** - Faster pulls and deployments
- **SBOM included** - Full software bill of materials
- **Verifiable** - Signed with Sigstore/Cosign

## Using Docker Hub Images

This project includes alternative configurations using official images from Docker Hub. As of October 23, 2025, these images are no longer being updated. These alternatives are included as a way to compare vulnerabilities across images. 

### Docker Compose with Docker Hub

Use the alternative compose file:

```bash
# Start services
docker-compose -f docker-compose.dockerhub.yml up -d

# View logs
docker-compose -f docker-compose.dockerhub.yml logs -f

# Run the sample app
docker-compose -f docker-compose.dockerhub.yml up app

# Stop services
docker-compose -f docker-compose.dockerhub.yml down
```

**Images used:**
- `minio/minio:latest` - Official MinIO server
- `python:3.13-slim` - Official Python base image
- `minio/mc:latest` - Official MinIO Client (for testing)

### Kubernetes with Docker Hub

Use the alternative Kubernetes manifests:

```bash
# Deploy using kustomize
kubectl apply -k k8s-dockerhub/

# Or apply individually
kubectl apply -f k8s-dockerhub/namespace.yaml
kubectl apply -f k8s-dockerhub/secret.yaml
kubectl apply -f k8s-dockerhub/deployment.yaml
kubectl apply -f k8s-dockerhub/service.yaml

# Run the test job
kubectl apply -f k8s-dockerhub/test-job.yaml
```

### Building with Docker Hub Images

To build the Python app using Docker Hub images:

```bash
cd app
docker build -f Dockerfile.dockerhub -t minio-app:dockerhub .
```

### Differences from Chainguard Images

| Feature | Chainguard Images | Docker Hub Images |
|---------|------------------|-------------------|
| **Registry Access** | Free & public | Free & public |
| **Image Size** | Smaller (minimal) | Larger (more packages) |
| **Security** | No CVEs, regular security updates | No longer updated |
| **User** | Non-root by default | May require configuration |
| **SBOM** | Included | Available separately |
| **Updates** | Daily automated builds | Regular releases |

Both options provide fully functional MinIO deployments.

## Troubleshooting

### MinIO server won't start
```bash
# Check logs
docker-compose logs minio

# Or in Kubernetes
kubectl logs -n minio-demo -l app=minio
```

### Application can't connect to MinIO
```bash
# Verify MinIO is running
curl http://localhost:9000/minio/health/live

# Check network connectivity in Docker
docker-compose exec app ping minio

# Check service in Kubernetes
kubectl get svc -n minio-demo
```

### Permission denied errors
Chainguard containers run as non-root. Ensure volume permissions are correct:

```bash
# For Docker volumes, this is handled automatically
# For Kubernetes PVs, you may need to set proper ownership
```

## Additional Resources

- [MinIO Documentation](https://min.io/docs/minio/linux/index.html)
- [Chainguard Images](https://www.chainguard.dev/docs/chainguard/containers)
- [MinIO Python SDK](https://min.io/docs/minio/linux/developers/python/minio-py.html)
- [Chainguard Console](https://console.chainguard.dev/)

## License

This sample application is provided as-is for demonstration purposes.
