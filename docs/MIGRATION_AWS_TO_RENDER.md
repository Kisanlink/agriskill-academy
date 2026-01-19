# Migration Guide: AWS to Render.com

**Date:** January 8, 2026
**Status:** Ready for execution
**Estimated Total Time:** 4-6 hours
**Estimated Downtime:** 15-30 minutes

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Pre-Migration Checklist](#pre-migration-checklist)
3. [Migration Components](#migration-components)
4. [Step-by-Step Migration](#step-by-step-migration)
5. [Rollback Plan](#rollback-plan)
6. [Cost Comparison](#cost-comparison)
7. [Troubleshooting](#troubleshooting)

---

## Overview

### Current AWS Stack

- **Compute:** ECS Fargate (256 CPU, 512 MB RAM)
- **Database:** RDS PostgreSQL 17.4 (db.t3.micro)
- **Cache:** ElastiCache Redis 7.1 (cache.t3.micro)
- **Storage:** S3 bucket for files
- **Networking:** VPC, NAT Gateway, ALB
- **Monitoring:** CloudWatch

### Target Render Stack

- **Compute:** Render Web Service (Docker)
- **Database:** Neon Serverless PostgreSQL
- **Cache:** Render Redis or Upstash Redis
- **Storage:** Cloudflare R2 (S3-compatible)
- **Networking:** Render built-in load balancing
- **Monitoring:** Render dashboard + logs

### Why Migrate?

- **Cost Reduction:** Eliminate NAT Gateway ($32/month), reduce compute costs
- **Simplicity:** Less infrastructure management
- **Developer Experience:** Easier deployments, better logs
- **Scalability:** Serverless database, auto-scaling web service

---

## Pre-Migration Checklist

### ✅ Before You Start

- [ ] **Backup Current Data**
  - [ ] Export RDS database
  - [ ] Download all S3 files
  - [ ] Document current environment variables

- [ ] **Create Accounts**
  - [ ] Render.com account
  - [ ] Neon.tech account
  - [ ] Cloudflare account (for R2)

- [ ] **Prepare Environment**
  - [ ] Install PostgreSQL client tools (`pg_dump`, `pg_restore`)
  - [ ] Install AWS CLI
  - [ ] Install Cloudflare Wrangler CLI (optional)

- [ ] **Communication**
  - [ ] Notify users of planned maintenance window
  - [ ] Set up maintenance page (optional)

---

## Migration Components

### 1. Database: RDS PostgreSQL → Neon

**Compatibility:** ✅ 100% compatible (both PostgreSQL)

**Migration Method:** `pg_dump` + `pg_restore`

**Data to Migrate:**
- All tables (users, job_posts, applications, etc.)
- Indexes and constraints
- No triggers or stored procedures (Go app handles logic)

### 2. File Storage: AWS S3 → Cloudflare R2

**Compatibility:** ✅ S3-compatible API

**Migration Method:** AWS CLI sync or console-based transfer

**Files to Migrate:**
- Resumes (PDFs)
- Certificates (images/PDFs)
- Employer logos (images)
- Student profile photos (images)

**Estimated Total Size:** Check with `aws s3 ls s3://your-bucket --recursive --summarize`

### 3. Cache: ElastiCache Redis → Render Redis

**Compatibility:** ✅ Redis protocol compatible

**Migration Method:** Fresh start (no persistent data needed)

**Note:** Redis is used for background job queue - safe to start fresh

### 4. Application: ECS Container → Render Web Service

**Method:** Git-based deployment from GitHub

**Requirements:**
- Dockerfile (already exists)
- Environment variables configuration

---

## Step-by-Step Migration

### Phase 1: Setup New Infrastructure (No Downtime)

#### Step 1.1: Create Neon Database

1. **Sign up at neon.tech**
2. **Create new project:**
   - Name: `agrijobs-production`
   - Region: Choose closest to your users (e.g., `us-east-1`)
   - PostgreSQL version: `17` (or latest)

3. **Get connection details:**
   - Navigate to Dashboard → Connection Details
   - Copy **Connection String** (will look like):
     ```
     postgresql://username:password@ep-xxxx.us-east-1.aws.neon.tech/neondb?sslmode=require
     ```

4. **Save credentials:**
   ```
   DB_HOST=ep-xxxx.us-east-1.aws.neon.tech
   DB_PORT=5432
   POSTGRES_USER=username
   POSTGRES_PASS=password
   DB_NAME=neondb
   DB_SSLMODE=require
   ```

#### Step 1.2: Create Cloudflare R2 Bucket

1. **Log in to Cloudflare Dashboard**
2. **Navigate to R2 → Create Bucket**
   - Bucket name: `agrijobs-files` (or your preferred name)
   - Location: Automatic

3. **Create API Token:**
   - R2 → Manage R2 API Tokens → Create API Token
   - Permissions: Object Read & Write
   - **Save credentials:**
     ```
     AWS_ACCESS_KEY_ID=<R2_ACCESS_KEY_ID>
     AWS_SECRET_ACCESS_KEY=<R2_SECRET_ACCESS_KEY>
     AWS_S3_ENDPOINT=https://<account_id>.r2.cloudflarestorage.com
     AWS_S3_BUCKET=agrijobs-files
     AWS_REGION=auto
     AWS_S3_FORCE_PATH_STYLE=false
     AWS_S3_DISABLE_SSL=false
     ```

4. **Optional - Setup Custom Domain:**
   - R2 → Your Bucket → Settings → Public Access
   - Connect custom domain: `files.agriskillacademy.com`

#### Step 1.3: Setup Render Redis (Optional)

**Option A: Render Managed Redis** (Recommended)
1. Render Dashboard → New → Redis
2. Name: `agrijobs-redis`
3. Plan: Free (25 MB) or Starter ($10/month, 256 MB)
4. Region: Same as web service
5. Copy connection details:
   ```
   REDIS_ADDR=redis://red-xxx.redis.render.com:6379
   ```

**Option B: Upstash Redis** (Serverless alternative)
1. Sign up at upstash.com
2. Create database → Choose region
3. Copy connection string

**Option C: Make Redis Optional**
- If your code already supports Redis as optional, you can skip this

#### Step 1.4: Create Render Web Service

1. **Connect GitHub Repository**
   - Render Dashboard → New → Web Service
   - Connect your GitHub account
   - Select repository: `Kisanlink/agrijobs` (or your fork)
   - Branch: `main` or `production`

2. **Configure Service:**
   - Name: `agrijobs-backend`
   - Region: Same as Neon database
   - Branch: `main`
   - Runtime: Docker
   - Dockerfile Path: `./Dockerfile`
   - Instance Type: `Starter` ($7/month) or `Standard` ($25/month)

3. **Do NOT deploy yet** - we need to configure environment variables first

---

### Phase 2: Migrate Database (15-30 min downtime)

#### Step 2.1: Export from AWS RDS

**Method A: Using CloudFormation (Automated)**

1. **Deploy migration stack:**
   ```bash
   aws cloudformation create-stack \
     --stack-name agrijobs-rds-migration \
     --template-body file://cloudformation/rds-migration-bastion.yml \
     --parameters \
       ParameterKey=VpcId,ParameterValue=vpc-xxxxx \
       ParameterKey=PublicSubnetId,ParameterValue=subnet-xxxxx \
       ParameterKey=RDSSecurityGroupId,ParameterValue=sg-xxxxx \
       ParameterKey=RDSEndpoint,ParameterValue=your-rds.xxx.rds.amazonaws.com \
       ParameterKey=DBPassword,ParameterValue=your-password \
       ParameterKey=NotificationEmail,ParameterValue=your-email@example.com \
     --capabilities CAPABILITY_IAM
   ```

2. **Monitor progress:**
   ```bash
   aws cloudformation describe-stacks \
     --stack-name agrijobs-rds-migration \
     --query 'Stacks[0].StackStatus'
   ```

3. **Download backup from S3:**
   ```bash
   aws s3 cp s3://agrijobs-rds-migration-xxxxx/asa_db_backup.dump ./
   ```

**Method B: Manual Export**

1. **Connect to RDS from local machine** (requires VPN or public access):
   ```bash
   pg_dump -h your-rds.xxx.rds.amazonaws.com \
     -U agrijobs \
     -d asa_db \
     -F c \
     -Z 9 \
     -f asa_db_backup.dump
   ```
   - Enter password when prompted
   - This will create a compressed custom-format dump

2. **Verify backup:**
   ```bash
   pg_restore --list asa_db_backup.dump | head -20
   ```

#### Step 2.2: Import to Neon

1. **Restore database:**
   ```bash
   pg_restore -h ep-xxxx.us-east-1.aws.neon.tech \
     -U username \
     -d neondb \
     -v \
     --no-owner \
     --no-acl \
     asa_db_backup.dump
   ```
   - Enter Neon password when prompted
   - `-v` shows verbose progress
   - `--no-owner` and `--no-acl` skip ownership info

2. **Verify data:**
   ```bash
   psql "postgresql://username:password@ep-xxxx.us-east-1.aws.neon.tech/neondb?sslmode=require"

   -- Check table counts
   SELECT
     schemaname,
     tablename,
     (SELECT COUNT(*) FROM quote_ident(schemaname) || '.' || quote_ident(tablename)) as row_count
   FROM pg_tables
   WHERE schemaname = 'public';

   -- Verify specific critical tables
   SELECT COUNT(*) FROM users;
   SELECT COUNT(*) FROM job_posts;
   SELECT COUNT(*) FROM applications;
   ```

3. **Compare with RDS counts** to ensure all data migrated

---

### Phase 3: Migrate Files to Cloudflare R2

#### Method A: AWS CLI Sync (Recommended for large datasets)

1. **Install and configure AWS CLI** (if not already done)

2. **Configure Cloudflare R2 endpoint:**
   ```bash
   aws configure set aws_access_key_id <R2_ACCESS_KEY_ID>
   aws configure set aws_secret_access_key <R2_SECRET_ACCESS_KEY>
   aws configure set region auto
   ```

3. **Sync files from S3 to R2:**
   ```bash
   aws s3 sync \
     s3://production-agrijobs-files-<account-id>/ \
     s3://agrijobs-files/ \
     --endpoint-url https://<account_id>.r2.cloudflarestorage.com \
     --source-region us-east-1
   ```

4. **Verify transfer:**
   ```bash
   # Check file count in R2
   aws s3 ls s3://agrijobs-files/ \
     --recursive \
     --endpoint-url https://<account_id>.r2.cloudflarestorage.com \
     --summarize

   # Compare with S3
   aws s3 ls s3://production-agrijobs-files-<account-id>/ \
     --recursive \
     --summarize
   ```

#### Method B: Console-Based Transfer (Small datasets)

1. **Download from S3:**
   ```bash
   aws s3 sync s3://production-agrijobs-files-<account-id>/ ./s3-backup/
   ```

2. **Upload to R2 using web console:**
   - Cloudflare Dashboard → R2 → agrijobs-files → Upload
   - Drag and drop folders

**Estimated Time:**
- < 1 GB: 5-15 minutes
- 1-10 GB: 30-60 minutes
- > 10 GB: Use CLI method

---

### Phase 4: Configure and Deploy on Render

#### Step 4.1: Set Environment Variables

1. **In Render Dashboard** → Your Web Service → Environment
2. **Add all variables** (see `docs/render.env.template`):

**Database:**
```
DB_HOST=ep-xxxx.us-east-1.aws.neon.tech
DB_PORT=5432
POSTGRES_USER=<neon-username>
POSTGRES_PASS=<neon-password>
DB_NAME=neondb
DB_SSLMODE=require
```

**Cloudflare R2:**
```
AWS_REGION=auto
AWS_S3_BUCKET=agrijobs-files
AWS_ACCESS_KEY_ID=<R2_ACCESS_KEY_ID>
AWS_SECRET_ACCESS_KEY=<R2_SECRET_ACCESS_KEY>
AWS_S3_ENDPOINT=https://<account_id>.r2.cloudflarestorage.com
AWS_S3_FORCE_PATH_STYLE=false
AWS_S3_DISABLE_SSL=false
```

**Redis:**
```
REDIS_ADDR=<redis-connection-string>
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=<if-required>
```

**Application:**
```
SERVER_PORT=8080
GIN_MODE=release
JWT_SECRET=<your-jwt-secret>
ASA_BASE_URL=https://agrijobs-backend.onrender.com
LOG_LEVEL=info
LOG_FORMAT=json
ENABLE_CORS=true
CORS_ALLOWED_ORIGINS=https://agriskillacademy.com
ASA_RUN_SEED=false
APP_ENV=production
```

**Firebase (if used):**
```
FIREBASE_PROJECT_ID=<your-project-id>
FIREBASE_CREDENTIALS_JSON=<service-account-json>
FIREBASE_WEB_API_KEY=<web-api-key>
```

#### Step 4.2: Deploy Application

1. **Click "Deploy"** in Render dashboard
2. **Monitor deployment logs**
3. **Wait for health check to pass** (checks `/health` endpoint)

#### Step 4.3: Verify Deployment

1. **Check health endpoint:**
   ```bash
   curl https://agrijobs-backend.onrender.com/health
   ```

2. **Check database connection:**
   ```bash
   curl https://agrijobs-backend.onrender.com/health/db
   ```

3. **Test file upload/download:**
   - Upload a test file through your application
   - Verify it appears in Cloudflare R2
   - Test download

4. **Test critical flows:**
   - User login
   - Job posting
   - Application submission
   - File uploads

---

### Phase 5: Update DNS and Go Live

#### Step 5.1: Update Frontend Configuration

1. **Update frontend API endpoint:**
   - Change `API_BASE_URL` from AWS ALB to Render URL
   - Example: `https://agrijobs-backend.onrender.com`

2. **Deploy frontend changes**

#### Step 5.2: DNS Configuration (If using custom domain)

1. **Add custom domain in Render:**
   - Web Service → Settings → Custom Domains
   - Add: `api.agriskillacademy.com`

2. **Update DNS records:**
   - Type: `CNAME`
   - Name: `api`
   - Value: `agrijobs-backend.onrender.com`
   - TTL: `300` (5 minutes for faster propagation)

3. **Wait for SSL certificate** (automatic via Let's Encrypt)

#### Step 5.3: Monitor Production

1. **Monitor Render logs:**
   - Dashboard → Logs (real-time)

2. **Check for errors:**
   ```bash
   # Watch logs
   render logs --service agrijobs-backend --tail
   ```

3. **Monitor database performance:**
   - Neon Dashboard → Metrics
   - Check query performance

4. **Monitor R2 usage:**
   - Cloudflare → R2 → Analytics

---

## Rollback Plan

### If Migration Fails

**Before Deleting AWS Resources:**
1. Keep AWS infrastructure running in parallel for 7-14 days
2. Monitor Render deployment for issues
3. Keep RDS snapshots for 30 days

**Emergency Rollback Steps:**

1. **Revert DNS:**
   ```
   CNAME api → <aws-alb-dns-name>
   ```

2. **Revert frontend API URL:**
   - Change back to AWS ALB endpoint
   - Redeploy frontend

3. **Verify AWS services:**
   - RDS still running
   - S3 still accessible
   - Redis still connected

**Recovery Time:** 5-10 minutes

---

## Cost Comparison

### AWS (Current - Monthly)

| Service | Configuration | Cost |
|---------|---------------|------|
| ECS Fargate | 0.25 vCPU, 0.5 GB | ~$13 |
| RDS PostgreSQL | db.t3.micro | ~$15 |
| ElastiCache Redis | cache.t3.micro | ~$13 |
| NAT Gateway | Data transfer | ~$32 |
| S3 Storage | 10 GB + requests | ~$1 |
| ALB | Load balancer | ~$16 |
| Data Transfer | 100 GB outbound | ~$9 |
| **Total** | | **~$99/month** |

### Render + Neon + R2 (New - Monthly)

| Service | Configuration | Cost |
|---------|---------------|------|
| Render Web Service | Starter (512 MB RAM) | $7 |
| Neon PostgreSQL | Free tier (0.5 GB) | $0 |
| Neon PostgreSQL | Serverless compute | ~$10 |
| Render Redis | Starter (256 MB) | $10 |
| Cloudflare R2 | 10 GB storage | $0.15 |
| Cloudflare R2 | 1M requests | ~$1 |
| **Total** | | **~$28/month** |

**Savings: ~$71/month (72% reduction)**

### Cost Optimization Tips

1. **Neon Free Tier:**
   - 0.5 GB storage free
   - 3 GB-hours compute free
   - Autosuspends when inactive

2. **Render Free Resources:**
   - Can start with free tier for testing
   - Redis free tier available (25 MB)

3. **Cloudflare R2:**
   - No egress fees (huge saving vs S3)
   - First 10 GB storage free
   - Very cheap request pricing

---

## Troubleshooting

### Database Connection Issues

**Problem:** Can't connect to Neon database

**Solutions:**
- ✅ Verify `DB_SSLMODE=require` is set
- ✅ Check Neon database is not suspended (auto-suspend after 5 min inactive)
- ✅ Verify connection string format
- ✅ Check IP allowlist in Neon (if configured)

### File Upload Failures

**Problem:** Files not uploading to R2

**Solutions:**
- ✅ Verify R2 API token has write permissions
- ✅ Check `AWS_S3_ENDPOINT` is correctly set
- ✅ Verify `AWS_REGION=auto` for R2
- ✅ Check CORS configuration in R2 bucket settings

### Redis Connection Issues

**Problem:** Redis connection failures

**Solutions:**
- ✅ Verify Redis is not required for basic functionality
- ✅ Check `REDIS_ADDR` format (should include redis:// protocol)
- ✅ Verify Redis password if authentication enabled
- ✅ Consider making Redis optional in code

### Performance Issues

**Problem:** Slow database queries on Neon

**Solutions:**
- ✅ Neon may auto-suspend - first query wakes it up (cold start)
- ✅ Enable connection pooling in your app
- ✅ Upgrade to paid plan for dedicated compute
- ✅ Check query performance in Neon dashboard

---

## Post-Migration Tasks

### Cleanup AWS Resources (After 7-14 days)

1. **Delete ECS Service:**
   ```bash
   aws ecs update-service --cluster <cluster> --service <service> --desired-count 0
   aws ecs delete-service --cluster <cluster> --service <service> --force
   ```

2. **Delete CloudFormation Stack:**
   ```bash
   aws cloudformation delete-stack --stack-name <your-stack-name>
   ```

3. **Verify Deletion:**
   - RDS snapshot created ✅
   - S3 files backed up ✅
   - Monitor for 30 days before final deletion

### Monitoring Setup

1. **Render Alerts:**
   - Dashboard → Service → Alerts
   - Configure email notifications for:
     - Service crashes
     - High CPU/memory usage
     - Deployment failures

2. **Neon Monitoring:**
   - Dashboard → Metrics
   - Monitor query performance
   - Set up consumption alerts

3. **Uptime Monitoring:**
   - Use UptimeRobot or similar
   - Monitor `/health` endpoint every 5 minutes

---

## Support Resources

- **Render Documentation:** https://render.com/docs
- **Neon Documentation:** https://neon.tech/docs
- **Cloudflare R2 Documentation:** https://developers.cloudflare.com/r2/
- **PostgreSQL Migration Guide:** https://www.postgresql.org/docs/current/backup-dump.html

---

**Migration Checklist:** See `docs/MIGRATION_CHECKLIST.md`
**Environment Template:** See `docs/render.env.template`
**CloudFormation Template:** See `cloudformation/rds-migration-bastion.yml`
