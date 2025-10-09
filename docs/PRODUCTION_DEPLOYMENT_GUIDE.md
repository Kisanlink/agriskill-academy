# AgriJobs Production Deployment Guide

## Quick Start

### 1. Environment Setup

Create a `.env` file from the template:

```bash
cp .env.example .env
```

### 2. Configure Required Variables

Edit `.env` and set these **REQUIRED** variables:

```bash
# Database Credentials
POSTGRES_USER=agrijobs
POSTGRES_PASS=your_secure_password_here
DB_NAME=asa_db

# JWT Secret (minimum 32 characters)
JWT_SECRET=your-super-secret-jwt-key-min-32-characters-long

# AWS S3 Configuration
AWS_REGION=us-east-1
AWS_S3_BUCKET=agrijobs-files
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key

# Redis Password
REDIS_PASSWORD=your_redis_password

# Base URL (no port, just the domain)
ASA_BASE_URL=http://localhost

# Admin Seeding (DISABLE IN PRODUCTION)
ASA_RUN_SEED=false
DEFAULT_ADMIN_EMAIL=admin@agrijobs.com
DEFAULT_ADMIN_PASSWORD=change_this_password
DEFAULT_ADMIN_NAME=System Administrator
DEFAULT_ADMIN_USERNAME=admin
```

### 3. Deploy with Docker Compose

```bash
# Start all services
docker-compose -f docker-compose.production.yml up -d

# View logs
docker-compose -f docker-compose.production.yml logs -f

# Stop all services
docker-compose -f docker-compose.production.yml down
```

### 4. Verify Deployment

Check if all services are healthy:

```bash
docker-compose -f docker-compose.production.yml ps
```

Access the API:
- Health Check: http://localhost:8080/health
- API Documentation: http://localhost:8080/swagger/index.html

## What Was Fixed

### Environment Variable Consistency

The following typos were corrected throughout the codebase:

| ❌ Old (Incorrect) | ✅ New (Correct) |
|-------------------|------------------|
| `POSTGRESS_USER` | `POSTGRES_USER` |
| `POSTGRESS_PASS` | `POSTGRES_PASS` |

### Files Updated:
1. **config/config.go** - Fixed environment variable reading
2. **docker-compose.production.yml** - Fixed container environment variables
3. **cloudformation/infrastructure.yml** - Fixed ECS task definition
4. **Dockerfile** - Updated health check to use ASA_BASE_URL

### Production-Ready Features:
- ✅ Dynamic port configuration via `SERVER_PORT`
- ✅ Configurable base URL via `ASA_BASE_URL`
- ✅ Health checks using environment variables
- ✅ Multi-stage Docker builds for optimization
- ✅ Non-root user for security
- ✅ Proper AWS S3 integration (no MinIO)
- ✅ Redis for job queue and caching
- ✅ PostgreSQL 16 with health checks

## AWS CloudFormation Deployment

For production AWS deployment, see the CloudFormation template:

```bash
# Deploy to AWS ECS Fargate
aws cloudformation create-stack \
  --stack-name agrijobs-production \
  --template-body file://cloudformation/infrastructure.yml \
  --parameters \
    ParameterKey=DBUsername,ParameterValue=agrijobs \
    ParameterKey=DBPassword,ParameterValue=YOUR_SECURE_PASSWORD \
    ParameterKey=JWTSecret,ParameterValue=YOUR_JWT_SECRET_MIN_32_CHARS \
    ParameterKey=ContainerImage,ParameterValue=YOUR_ECR_IMAGE_URI \
  --capabilities CAPABILITY_NAMED_IAM
```

## Environment Variables Reference

### Database (PostgreSQL)
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `POSTGRES_USER` - Database username
- `POSTGRES_PASS` - Database password
- `DB_NAME` - Database name
- `DB_SSLMODE` - SSL mode (disable for local, require for production)

### Server
- `SERVER_PORT` - Application port (default: 8080)
- `GIN_MODE` - Gin mode (debug/release)
- `ASA_BASE_URL` - Base URL without port

### AWS S3
- `AWS_REGION` - AWS region
- `AWS_S3_BUCKET` - S3 bucket name
- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret key

### Redis
- `REDIS_ADDR` - Redis address (host:port)
- `REDIS_PORT` - Redis port (default: 6379)
- `REDIS_PASSWORD` - Redis password
- `REDIS_DB` - Redis database number (default: 0)

### Authentication
- `JWT_SECRET` - JWT signing secret (min 32 characters)

### Security
- `RATE_LIMIT_REQUESTS` - Max requests per window (default: 100)
- `RATE_LIMIT_WINDOW` - Rate limit window (default: 1m)
- `MAX_REQUEST_SIZE` - Max request size in bytes (default: 10MB)
- `MAX_FILE_SIZE` - Max file upload size in bytes (default: 5MB)
- `ENABLE_CORS` - Enable CORS (true/false)
- `CORS_ALLOWED_ORIGINS` - Comma-separated allowed origins

### Admin Seeding
- `ASA_RUN_SEED` - Enable/disable seeding (true/false)
- `DEFAULT_ADMIN_EMAIL` - Default admin email
- `DEFAULT_ADMIN_PASSWORD` - Default admin password
- `DEFAULT_ADMIN_NAME` - Default admin name
- `DEFAULT_ADMIN_USERNAME` - Default admin username

### Logging
- `LOG_LEVEL` - Log level (debug/info/warn/error)
- `LOG_FORMAT` - Log format (json/console)
- `LOG_DEVELOPMENT` - Development mode logging (true/false)
- `LOG_OUTPUT_PATH` - Log output path (stdout/file path)

## Troubleshooting

### Database Connection Errors

If you see "failed to connect to database" errors:

1. Verify environment variables are set correctly:
```bash
docker-compose -f docker-compose.production.yml exec app env | grep POSTGRES
```

2. Check PostgreSQL logs:
```bash
docker-compose -f docker-compose.production.yml logs postgres
```

3. Verify PostgreSQL is healthy:
```bash
docker-compose -f docker-compose.production.yml ps postgres
```

### Health Check Failures

If health checks are failing:

1. Check the ASA_BASE_URL is set correctly
2. Verify the app container can reach the health endpoint:
```bash
docker-compose -f docker-compose.production.yml exec app curl -f ${ASA_BASE_URL}/health
```

### Missing Environment Variables

If you see warnings about missing variables:

1. Ensure `.env` file exists in the project root
2. Verify all required variables are set (see section 2 above)
3. Rebuild containers after updating `.env`:
```bash
docker-compose -f docker-compose.production.yml down
docker-compose -f docker-compose.production.yml up -d --build
```

## Security Checklist

Before deploying to production:

- [ ] Change all default passwords
- [ ] Set ASA_RUN_SEED to false
- [ ] Use strong JWT_SECRET (min 32 characters)
- [ ] Enable SSL/TLS for database (DB_SSLMODE=require)
- [ ] Configure proper CORS_ALLOWED_ORIGINS
- [ ] Review and adjust rate limits
- [ ] Enable AWS S3 bucket encryption
- [ ] Rotate AWS credentials regularly
- [ ] Monitor CloudWatch logs and alarms

## Support

For issues or questions:
- GitHub Issues: https://github.com/Kisanlink/agriskill-academy/issues
- Documentation: See DEPLOYMENT.md for detailed AWS instructions

## License

Copyright © 2025 KisanLink
