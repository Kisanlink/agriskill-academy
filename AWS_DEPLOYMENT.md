# AWS ECS Deployment Guide

## Overview

This guide covers deploying the AgriJobs backend to AWS using ECS Fargate, with proper IAM role-based S3 access.

## Architecture

```
Internet → ALB → ECS Fargate Tasks → RDS PostgreSQL
                       ↓
                   ElastiCache Redis
                       ↓
                   S3 Bucket (via IAM Task Role)
```

## Prerequisites

1. AWS CLI configured with appropriate credentials
2. Docker image built and pushed to ECR
3. AWS account with permissions to create CloudFormation stacks

## S3 Access - IAM Task Role (Recommended)

### ✅ How It Works

The CloudFormation template creates an **ECS Task Role** with S3 permissions. When your application runs in ECS:

1. **IAM Task Role** is attached to the ECS task (lines 631-650 in infrastructure.yml)
2. **S3 Bucket Policy** trusts the Task Role (lines 512-529 in infrastructure.yml)
3. **AWS SDK** automatically uses the Task Role credentials (no explicit credentials needed)
4. **Application code** detects empty credentials and uses IAM role (cmd/server/main.go:147-150)

### Required Environment Variables

**In ECS (CloudFormation):**
```yaml
AWS_REGION: us-east-1
AWS_S3_BUCKET: production-agrijobs-files-123456789
# AWS_ACCESS_KEY_ID: NOT SET (uses IAM role)
# AWS_SECRET_ACCESS_KEY: NOT SET (uses IAM role)
```

**Application automatically:**
- Detects credentials are empty
- Uses AWS SDK's default credential chain
- SDK finds IAM Task Role credentials
- Accesses S3 with proper permissions

### Security Benefits

✅ No hardcoded credentials
✅ Automatic credential rotation by AWS
✅ Fine-grained permissions per task
✅ CloudTrail audit logging
✅ No credential leakage risk

## Deployment Steps

### 1. Build and Push Docker Image

```bash
# Set your AWS account ID and region
AWS_ACCOUNT_ID=123456789
AWS_REGION=us-east-1
ECR_REPO=agrijobs

# Login to ECR
aws ecr get-login-password --region $AWS_REGION | \
  docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Create ECR repository (if not exists)
aws ecr create-repository \
  --repository-name $ECR_REPO \
  --region $AWS_REGION

# Build Docker image
docker build -t $ECR_REPO:latest .

# Tag for ECR
docker tag $ECR_REPO:latest \
  $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO:latest

# Push to ECR
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO:latest
```

### 2. Deploy CloudFormation Stack

```bash
aws cloudformation create-stack \
  --stack-name agrijobs-production \
  --template-body file://cloudformation/infrastructure.yml \
  --parameters \
    ParameterKey=EnvironmentName,ParameterValue=production \
    ParameterKey=ContainerImage,ParameterValue=$AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPO:latest \
    ParameterKey=DBUsername,ParameterValue=agrijobs \
    ParameterKey=DBPassword,ParameterValue=YOUR_SECURE_DB_PASSWORD \
    ParameterKey=JWTSecret,ParameterValue=YOUR_JWT_SECRET_MIN_32_CHARS \
    ParameterKey=DefaultAdminEmail,ParameterValue=admin@your-domain.com \
    ParameterKey=DefaultAdminPassword,ParameterValue=YOUR_SECURE_ADMIN_PASSWORD \
    ParameterKey=DefaultAdminName,ParameterValue="System Administrator" \
    ParameterKey=DefaultAdminUsername,ParameterValue=admin \
  --capabilities CAPABILITY_NAMED_IAM \
  --region $AWS_REGION

# Monitor stack creation
aws cloudformation wait stack-create-complete \
  --stack-name agrijobs-production \
  --region $AWS_REGION
```

**Important Notes:**
- **Admin Account**: The application automatically creates an admin account on first startup using the provided credentials
- **S3 Access**: Do NOT provide `AwsAccessKeyId` or `AwsSecretAccessKey` parameters - IAM Task Role is used automatically
- **Password Security**: Change the admin password immediately after first login!
- **Seeding Safety**: Admin creation only happens if no admin exists (safe for restarts/redeployments)

### 3. Get Stack Outputs

```bash
aws cloudformation describe-stacks \
  --stack-name agrijobs-production \
  --region $AWS_REGION \
  --query 'Stacks[0].Outputs' \
  --output table
```

Key outputs:
- **ALBEndpoint**: Your application URL
- **S3BucketName**: S3 bucket for file storage
- **RDSEndpoint**: Database endpoint
- **RedisEndpoint**: Cache endpoint

### 4. Verify Deployment

```bash
# Get ALB endpoint
ALB_URL=$(aws cloudformation describe-stacks \
  --stack-name agrijobs-production \
  --region $AWS_REGION \
  --query 'Stacks[0].Outputs[?OutputKey==`ALBEndpoint`].OutputValue' \
  --output text)

# Test health endpoint
curl $ALB_URL/health

# Check API documentation
open $ALB_URL/swagger/index.html
```

### 5. Monitor Application Logs

```bash
# Get log group name
aws logs describe-log-groups \
  --log-group-name-prefix /ecs/production-agrijobs \
  --region $AWS_REGION

# Tail logs
aws logs tail /ecs/production-agrijobs \
  --follow \
  --region $AWS_REGION
```

## Admin Account Seeding

The application includes an automatic admin account creation feature (seeding) that runs on first startup.

### How It Works

1. **On Application Startup:**
   - Checks if `ASA_RUN_SEED=true` environment variable is set
   - Queries database for existing admin accounts (role = 'asa_admin')
   - If no admin exists, creates one using provided credentials

2. **Safety Features:**
   - ✅ **Idempotent**: Safe to restart/redeploy (won't create duplicate admins)
   - ✅ **Transaction-based**: Prevents race conditions in multi-task deployments
   - ✅ **Checks by role AND email/username**: Prevents conflicts

3. **Required Environment Variables:**
   ```yaml
   ASA_RUN_SEED: 'true'
   DEFAULT_ADMIN_EMAIL: admin@your-domain.com
   DEFAULT_ADMIN_PASSWORD: (stored in Secrets Manager)
   DEFAULT_ADMIN_NAME: System Administrator
   DEFAULT_ADMIN_USERNAME: admin
   ```

4. **First Login:**
   - Use the credentials provided in CloudFormation parameters
   - **⚠️ CRITICAL**: Change the password immediately after first login
   - The seeding logs will show the admin account details (check CloudWatch Logs)

### Seeding Logs

Check the application logs to confirm admin creation:

```bash
aws logs tail /ecs/production-agrijobs --follow | grep -A 5 "admin account created"
```

You should see:
```
✅ Default admin account created successfully!
   Username: admin
   Email: admin@your-domain.com
   Password: ********
   ID: <uuid>
   ⚠️  IMPORTANT: Change the default password after first login!
```

### Disabling Seeding

If you want to create admin accounts manually, set `ASA_RUN_SEED=false` in the task definition.

**Note:** Seeding is recommended for initial deployment. After the first admin is created, additional admins should be created through the admin panel or API.

## Local Development (MinIO)

For local development, use MinIO as S3-compatible storage:

### 1. Start services with Docker Compose

```bash
# Create .env file
cp .env.example .env

# Edit .env - set MinIO credentials
AWS_ACCESS_KEY_ID=minio
AWS_SECRET_ACCESS_KEY=minio123
AWS_S3_ENDPOINT=http://minio:9000
AWS_S3_FORCE_PATH_STYLE=true
AWS_S3_DISABLE_SSL=true

# Start all services
docker-compose -f docker-compose.production.yml up -d
```

### 2. Access MinIO Console

```
URL: http://localhost:9001
Username: minio
Password: minio123
```

### 3. Create S3 Bucket

In MinIO console, create bucket named: `agrijobs-files`

## Troubleshooting

### S3 Access Denied Errors

**Symptom:**
```
ERROR: Access Denied (Service: Amazon S3; Status Code: 403)
```

**Solutions:**

1. **Verify IAM Task Role has S3 permissions:**
```bash
aws iam get-role-policy \
  --role-name production-agrijobs-task-role \
  --policy-name S3FullAccess
```

2. **Check bucket policy trusts task role:**
```bash
aws s3api get-bucket-policy \
  --bucket production-agrijobs-files-123456789
```

3. **Verify task is using correct role:**
```bash
aws ecs describe-tasks \
  --cluster production-agrijobs-cluster \
  --tasks <task-arn>
```

### Application Can't Connect to S3

**Check application logs:**
```bash
aws logs tail /ecs/production-agrijobs --follow
```

**Look for:**
```
Using S3 storage with IAM role (AWS ECS): bucket=production-agrijobs-files-123456789
```

If you see:
```
Using S3 storage with explicit credentials
```

**Problem:** Credentials were set in environment variables (shouldn't be set in ECS)

**Solution:** Remove AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY from task definition

## Cost Estimate

**Monthly costs (us-east-1, production environment):**

- ECS Fargate (2 tasks, 0.5 vCPU, 1GB RAM): ~$30
- ALB: ~$20
- RDS PostgreSQL (db.t3.micro): ~$15
- ElastiCache Redis (cache.t3.micro): ~$15
- S3 (100GB storage, 1TB transfer): ~$30
- NAT Gateway (2 AZs): ~$90
- CloudWatch Logs: ~$10
- Data Transfer: ~$20

**Total: ~$230/month**

**Cost optimization tips:**
- Use FARGATE_SPOT for non-critical tasks (60% cheaper)
- Enable S3 Intelligent-Tiering
- Use Single-AZ RDS for staging
- Reduce NAT Gateway count in staging

## Security Best Practices

1. ✅ **Use IAM Task Roles** (not hardcoded credentials)
2. ✅ **Enable encryption at rest** (RDS, S3, Redis)
3. ✅ **Use Secrets Manager** for sensitive data
4. ✅ **Enable CloudTrail** for audit logging
5. ✅ **Use private subnets** for tasks
6. ✅ **Restrict security groups** to minimum required access
7. ✅ **Enable container insights** for monitoring
8. ✅ **Use Application Load Balancer** with health checks
9. ✅ **Enable VPC Flow Logs** for network monitoring
10. ✅ **Implement least privilege** IAM policies

## Updating the Application

### Rolling Update

```bash
# Build and push new image
docker build -t agrijobs:v2 .
docker tag agrijobs:v2 $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/agrijobs:v2
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/agrijobs:v2

# Update service
aws ecs update-service \
  --cluster production-agrijobs-cluster \
  --service production-agrijobs-service \
  --force-new-deployment
```

### Blue/Green Deployment

Use AWS CodeDeploy for zero-downtime deployments:

```bash
aws deploy create-deployment \
  --application-name agrijobs \
  --deployment-group-name production \
  --deployment-config-name CodeDeployDefault.ECSAllAtOnce \
  --description "Deploy version 2.0"
```

## Monitoring

### CloudWatch Metrics

Key metrics to monitor:
- ECS CPU/Memory utilization
- ALB request count and latency
- RDS connections and read/write IOPS
- S3 request count and errors

### CloudWatch Alarms

The CloudFormation template creates alarms for:
- High CPU usage (>80%)
- High memory usage (>85%)

### Application Logs

View logs in CloudWatch Logs:
```bash
aws logs tail /ecs/production-agrijobs \
  --since 1h \
  --format short
```

## Backup and Disaster Recovery

### RDS Automated Backups

- Retention: 7 days (production), 3 days (staging)
- Backup window: 03:00-04:00 UTC
- Maintenance window: Sunday 04:00-05:00 UTC

### S3 Versioning

- Enabled by default
- Lifecycle policy: Delete old versions after 90 days

### Manual Backup

```bash
# Create RDS snapshot
aws rds create-db-snapshot \
  --db-snapshot-identifier agrijobs-manual-backup-$(date +%Y%m%d) \
  --db-instance-identifier production-agrijobs-db
```

## Scaling

### Auto Scaling Configuration

The template includes auto-scaling based on:
- CPU utilization (target: 70%)
- Memory utilization (target: 80%)

Min tasks: 2
Max tasks: 10

### Manual Scaling

```bash
# Scale to 5 tasks
aws ecs update-service \
  --cluster production-agrijobs-cluster \
  --service production-agrijobs-service \
  --desired-count 5
```

## Support

For issues or questions:
- GitHub Issues: https://github.com/Kisanlink/agriskill-academy/issues
- Documentation: See README.md

## License

Copyright © 2025 KisanLink
