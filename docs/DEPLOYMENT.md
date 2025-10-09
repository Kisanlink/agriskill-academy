# AWS Production Deployment Guide - AgriJobs Backend

Complete guide for deploying the AgriJobs backend to AWS using Docker and CloudFormation.

## 📋 Prerequisites

### Required Tools
- **AWS CLI** (v2.x or later) - [Install Guide](https://aws.amazon.com/cli/)
- **Docker** (20.10+) - [Install Guide](https://docs.docker.com/get-docker/)
- **AWS Account** with appropriate permissions
- **Git** for version control

### AWS Permissions Required
Your AWS IAM user/role needs permissions for:
- CloudFormation (full access)
- ECS (full access)
- RDS (full access)
- ElastiCache (full access)
- S3 (full access)
- EC2 (VPC, Security Groups, Load Balancers)
- IAM (role creation)
- Secrets Manager
- CloudWatch (Logs and Alarms)
- ECR (Elastic Container Registry)

---

## 🚀 Quick Start Deployment

### Step 1: Configure AWS CLI

```bash
# Configure AWS credentials
aws configure

# Enter your credentials:
# AWS Access Key ID: YOUR_ACCESS_KEY
# AWS Secret Access Key: YOUR_SECRET_KEY
# Default region name: us-east-1
# Default output format: json
```

### Step 2: Set Environment Variables

Create a `.env.production` file:

```bash
# Environment
export AWS_REGION=us-east-1
export ENVIRONMENT_NAME=production

# Database
export DB_PASSWORD=YourStrongPassword123!
export DB_USERNAME=agrijobs
export DB_NAME=asa_db

# Security
export JWT_SECRET=your-super-secret-jwt-key-min-32-characters

# Application
export CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
```

Load the environment:
```bash
source .env.production
```

### Step 3: Build and Push Docker Image to ECR

```bash
# 1. Create ECR repository
aws ecr create-repository \
    --repository-name production-agrijobs-backend \
    --region us-east-1

# 2. Get ECR login
aws ecr get-login-password --region us-east-1 | \
    docker login --username AWS --password-stdin \
    YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com

# 3. Build Docker image
docker build -t agrijobs-backend:latest .

# 4. Tag image
docker tag agrijobs-backend:latest \
    YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/production-agrijobs-backend:latest

# 5. Push to ECR
docker push YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/production-agrijobs-backend:latest
```

### Step 4: Deploy Infrastructure with CloudFormation

```bash
# Deploy the CloudFormation stack
aws cloudformation create-stack \
    --stack-name production-agrijobs-infrastructure \
    --template-body file://cloudformation/infrastructure.yml \
    --parameters \
        ParameterKey=EnvironmentName,ParameterValue=production \
        ParameterKey=DBPassword,ParameterValue="${DB_PASSWORD}" \
        ParameterKey=JWTSecret,ParameterValue="${JWT_SECRET}" \
        ParameterKey=ContainerImage,ParameterValue="YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/production-agrijobs-backend:latest" \
        ParameterKey=DesiredCount,ParameterValue=2 \
        ParameterKey=ContainerCpu,ParameterValue=512 \
        ParameterKey=ContainerMemory,ParameterValue=1024 \
    --capabilities CAPABILITY_NAMED_IAM \
    --region us-east-1

# Monitor stack creation
aws cloudformation wait stack-create-complete \
    --stack-name production-agrijobs-infrastructure \
    --region us-east-1

# Check status
aws cloudformation describe-stacks \
    --stack-name production-agrijobs-infrastructure \
    --region us-east-1 \
    --query 'Stacks[0].StackStatus'
```

### Step 5: Get Application URL

```bash
# Get the Load Balancer URL
aws cloudformation describe-stacks \
    --stack-name production-agrijobs-infrastructure \
    --region us-east-1 \
    --query 'Stacks[0].Outputs[?OutputKey==`ALBEndpoint`].OutputValue' \
    --output text
```

Your application will be available at: `http://YOUR-ALB-DNS-NAME`

---

## 📦 What Gets Deployed

### Infrastructure Components

1. **Networking**
   - VPC with public and private subnets across 2 AZs
   - Internet Gateway and NAT Gateways
   - Route tables and security groups

2. **Database Layer**
   - **RDS PostgreSQL 16.1** (Multi-AZ in production)
   - Automated backups (7 days retention)
   - Encryption at rest
   - Private subnet deployment

3. **Cache Layer**
   - **ElastiCache Redis 7.1** (Multi-node in production)
   - Encryption at rest and in transit
   - Automatic failover

4. **Storage**
   - **S3 Bucket** for file storage
   - Versioning enabled
   - Encryption at rest
   - CORS configured

5. **Compute**
   - **ECS Fargate** cluster
   - Auto-scaling (2-10 tasks)
   - Blue/green deployments
   - Health checks

6. **Load Balancing**
   - **Application Load Balancer**
   - Health checks
   - HTTP/HTTPS support

7. **Monitoring**
   - **CloudWatch Logs** (14 days retention)
   - **CloudWatch Alarms** for CPU and memory
   - Container Insights

8. **Security**
   - **Secrets Manager** for sensitive data
   - IAM roles with least privilege
   - Security groups with restricted access

---

## 🔧 Configuration

### Environment-Specific Parameters

#### Production
```yaml
EnvironmentName: production
DBInstanceClass: db.t3.small
RedisNodeType: cache.t3.small
DesiredCount: 2
ContainerCpu: 512
ContainerMemory: 1024
MultiAZ: true
BackupRetentionPeriod: 7
```

#### Staging
```yaml
EnvironmentName: staging
DBInstanceClass: db.t3.micro
RedisNodeType: cache.t3.micro
DesiredCount: 1
ContainerCpu: 256
ContainerMemory: 512
MultiAZ: false
BackupRetentionPeriod: 3
```

---

## 🔄 Updating the Application

### Update Docker Image

```bash
# 1. Build new image
docker build -t agrijobs-backend:latest .

# 2. Tag with version
docker tag agrijobs-backend:latest \
    YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/production-agrijobs-backend:v1.0.1

# 3. Push to ECR
docker push YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/production-agrijobs-backend:v1.0.1

# 4. Update ECS service (triggers rolling update)
aws ecs update-service \
    --cluster production-agrijobs-cluster \
    --service production-agrijobs-service \
    --force-new-deployment \
    --region us-east-1
```

### Update Infrastructure

```bash
# Update CloudFormation stack
aws cloudformation update-stack \
    --stack-name production-agrijobs-infrastructure \
    --template-body file://cloudformation/infrastructure.yml \
    --parameters \
        ParameterKey=EnvironmentName,UsePreviousValue=true \
        ParameterKey=DBPassword,UsePreviousValue=true \
        ParameterKey=JWTSecret,UsePreviousValue=true \
        ParameterKey=ContainerImage,ParameterValue="YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/production-agrijobs-backend:v1.0.1" \
        ParameterKey=DesiredCount,ParameterValue=3 \
    --capabilities CAPABILITY_NAMED_IAM \
    --region us-east-1
```

---

## 📊 Monitoring and Logs

### View Application Logs

```bash
# Get log group
aws logs tail /ecs/production-agrijobs \
    --follow \
    --region us-east-1

# Filter for errors
aws logs filter-log-events \
    --log-group-name /ecs/production-agrijobs \
    --filter-pattern "ERROR" \
    --region us-east-1
```

### Check ECS Service Health

```bash
# Get service details
aws ecs describe-services \
    --cluster production-agrijobs-cluster \
    --services production-agrijobs-service \
    --region us-east-1

# List running tasks
aws ecs list-tasks \
    --cluster production-agrijobs-cluster \
    --service production-agrijobs-service \
    --region us-east-1
```

### CloudWatch Alarms

Access CloudWatch console to view:
- High CPU usage alerts (>80%)
- High memory usage alerts (>85%)
- RDS performance metrics
- Redis cache metrics

---

## 🔐 Security Best Practices

### Secrets Management

```bash
# Update database password
aws secretsmanager update-secret \
    --secret-id production/agrijobs/db-password \
    --secret-string "NewStrongPassword123!" \
    --region us-east-1

# Update JWT secret
aws secretsmanager update-secret \
    --secret-id production/agrijobs/jwt-secret \
    --secret-string "new-jwt-secret-min-32-characters" \
    --region us-east-1
```

### Database Access

```bash
# Create bastion host for database access (if needed)
# Or use AWS Systems Manager Session Manager

# Connect to RDS (from bastion/VPN)
psql -h RDS_ENDPOINT -U agrijobs -d asa_db
```

### S3 Bucket Access

```bash
# List files in S3 bucket
aws s3 ls s3://production-agrijobs-files-YOUR_ACCOUNT_ID/ \
    --region us-east-1

# Download a file
aws s3 cp s3://production-agrijobs-files-YOUR_ACCOUNT_ID/path/to/file.pdf ./
```

---

## 💰 Cost Optimization

### Estimated Monthly Costs (us-east-1)

**Production (2 tasks running 24/7):**
- ECS Fargate (2 x 0.5vCPU, 1GB): ~$30
- RDS PostgreSQL (db.t3.small, Multi-AZ): ~$60
- ElastiCache Redis (cache.t3.micro, 2 nodes): ~$30
- NAT Gateway (2 AZs): ~$65
- Application Load Balancer: ~$22
- S3 Storage (100GB): ~$2.50
- Data Transfer: ~$10
- **Total: ~$220/month**

**Staging (1 task, smaller instances):**
- **Total: ~$80-100/month**

### Cost Reduction Tips

1. **Use Fargate Spot** (already configured, up to 70% savings)
2. **Single NAT Gateway** for non-production
3. **Schedule** for dev environments (stop at night)
4. **Reserved Instances** for RDS in production
5. **S3 Lifecycle Policies** (already configured)

---

## 🐛 Troubleshooting

### Common Issues

#### 1. Service Not Starting
```bash
# Check ECS service events
aws ecs describe-services \
    --cluster production-agrijobs-cluster \
    --services production-agrijobs-service \
    --region us-east-1 \
    --query 'services[0].events[0:10]'

# Check task stopped reasons
aws ecs describe-tasks \
    --cluster production-agrijobs-cluster \
    --tasks TASK_ID \
    --region us-east-1
```

#### 2. Database Connection Issues
- Verify security group allows ECS tasks
- Check RDS endpoint in Secrets Manager
- Verify DB credentials

#### 3. S3 Access Issues
- Check IAM task role permissions
- Verify bucket policy
- Check bucket name in environment variables

#### 4. High Costs
- Review NAT Gateway usage
- Check data transfer costs
- Optimize ECS task sizing

---

## 🗑️ Cleanup / Deletion

### Delete Stack (Warning: This deletes everything!)

```bash
# Delete CloudFormation stack
aws cloudformation delete-stack \
    --stack-name production-agrijobs-infrastructure \
    --region us-east-1

# Monitor deletion
aws cloudformation wait stack-delete-complete \
    --stack-name production-agrijobs-infrastructure \
    --region us-east-1
```

### Delete ECR Repository

```bash
# Delete ECR repository and all images
aws ecr delete-repository \
    --repository-name production-agrijobs-backend \
    --force \
    --region us-east-1
```

---

## 📚 Additional Resources

- [AWS ECS Best Practices](https://docs.aws.amazon.com/AmazonECS/latest/bestpracticesguide/)
- [RDS PostgreSQL Guide](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html)
- [ElastiCache Redis Guide](https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/)
- [CloudFormation Documentation](https://docs.aws.amazon.com/cloudformation/)

---

## 🆘 Support

For issues or questions:
1. Check CloudWatch Logs
2. Review ECS service events
3. Check security group configurations
4. Verify environment variables

## 📝 Notes

- **First Deployment**: Takes ~20-30 minutes
- **Updates**: ECS rolling updates take ~5-10 minutes
- **Backups**: RDS automated daily backups enabled
- **Monitoring**: CloudWatch Container Insights enabled
- **Security**: All data encrypted at rest and in transit
