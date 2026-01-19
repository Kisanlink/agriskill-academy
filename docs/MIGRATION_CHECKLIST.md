# Migration Checklist: AWS to Render

**Date Started:** _______________
**Completed By:** _______________
**Status:** ⬜ Not Started | ⬜ In Progress | ⬜ Complete

---

## Pre-Migration (Before Starting)

### Account Setup
- [ ] Create Render.com account
- [ ] Create Neon.tech account
- [ ] Create Cloudflare account (for R2)
- [ ] Verify email addresses for all accounts

### Backup Current System
- [ ] Export full RDS database backup
  - Method: ⬜ CloudFormation | ⬜ Manual pg_dump
  - Backup location: _______________
  - Backup size: _______________ MB/GB
  - Backup date: _______________
- [ ] Download all S3 files (or verify they're safe)
  - Total files: _______________
  - Total size: _______________ GB
  - Backup location: _______________
- [ ] Document all current environment variables
  - Saved to: _______________

### Tool Installation
- [ ] Install PostgreSQL client tools (psql, pg_dump, pg_restore)
  - Version: _______________
- [ ] Install AWS CLI
  - Version: _______________
  - Configured: ⬜ Yes
- [ ] Install Cloudflare Wrangler (optional)

### Communication
- [ ] Schedule maintenance window
  - Date: _______________
  - Time: _______________
  - Duration: _______________ hours
- [ ] Notify users/stakeholders
- [ ] Prepare maintenance page (optional)

---

## Phase 1: New Infrastructure Setup

### Neon Database
- [ ] Create Neon project
  - Project name: _______________
  - Region: _______________
  - PostgreSQL version: _______________
- [ ] Save connection details
  - Host: _______________
  - Port: 5432
  - Database name: _______________
  - Username: _______________
  - Password: ⬜ Saved securely
- [ ] Test connection from local machine
  ```bash
  psql "postgresql://user:pass@host/db?sslmode=require"
  ```

### Cloudflare R2
- [ ] Create R2 bucket
  - Bucket name: _______________
  - Region: Automatic
- [ ] Create API token
  - Access Key ID: ⬜ Saved securely
  - Secret Access Key: ⬜ Saved securely
  - Account ID: _______________
- [ ] Note R2 endpoint
  - Endpoint: https://_______________.r2.cloudflarestorage.com
- [ ] (Optional) Setup custom domain
  - Domain: _______________
  - SSL: ⬜ Active

### Render Redis (Optional)
- [ ] Create Redis instance
  - Plan: ⬜ Free | ⬜ Starter | ⬜ Skip (optional)
  - Region: _______________
  - Connection string: ⬜ Saved securely

### Render Web Service
- [ ] Connect GitHub repository
  - Repository: _______________
  - Branch: _______________
- [ ] Configure service (DO NOT DEPLOY YET)
  - Service name: _______________
  - Region: _______________
  - Instance type: ⬜ Free | ⬜ Starter | ⬜ Standard
  - Environment variables: ⬜ Added (from template)

---

## Phase 2: Database Migration

### Export from RDS

**Method A: CloudFormation (Recommended)**
- [ ] Review CloudFormation template
- [ ] Update parameters
  - VPC ID: _______________
  - Subnet ID: _______________
  - RDS Security Group: _______________
  - RDS Endpoint: _______________
- [ ] Deploy CloudFormation stack
  ```bash
  aws cloudformation create-stack --stack-name agrijobs-rds-migration --template-body file://cloudformation/rds-migration-bastion.yml ...
  ```
- [ ] Monitor stack creation
  - Status: ⬜ CREATE_IN_PROGRESS | ⬜ CREATE_COMPLETE
- [ ] Wait for export to complete (~15-30 min)
- [ ] Download backup from S3
  - S3 URI: _______________
  - Downloaded to: _______________

**Method B: Manual Export**
- [ ] Connect to RDS (via VPN or bastion)
- [ ] Run pg_dump
  ```bash
  pg_dump -h <rds-host> -U agrijobs -d asa_db -F c -Z 9 -f asa_db_backup.dump
  ```
- [ ] Verify backup file created
  - File size: _______________ MB
  - Location: _______________

### Verify Export
- [ ] Check backup file integrity
  ```bash
  pg_restore --list asa_db_backup.dump | head -20
  ```
- [ ] Record table counts from RDS
  - Users: _______________
  - Job Posts: _______________
  - Applications: _______________
  - (Add other critical tables)

### Import to Neon
- [ ] Restore to Neon database
  ```bash
  pg_restore -h <neon-host> -U <user> -d neondb -v --no-owner --no-acl asa_db_backup.dump
  ```
- [ ] Monitor restore progress
  - Duration: _______________ minutes
  - Errors: ⬜ None | ⬜ See logs

### Verify Import
- [ ] Connect to Neon database
- [ ] Compare table counts
  - Users: RDS ___ | Neon ___ | ⬜ Match
  - Job Posts: RDS ___ | Neon ___ | ⬜ Match
  - Applications: RDS ___ | Neon ___ | ⬜ Match
- [ ] Spot check critical data
  - [ ] Recent users exist
  - [ ] Recent job posts exist
  - [ ] Recent applications exist
- [ ] Test queries
  - [ ] SELECT queries work
  - [ ] JOIN queries work
  - [ ] Indexes present

---

## Phase 3: File Storage Migration

### S3 to Cloudflare R2

**Method A: AWS CLI Sync**
- [ ] Configure R2 endpoint in AWS CLI
- [ ] Sync files from S3 to R2
  ```bash
  aws s3 sync s3://source-bucket/ s3://r2-bucket/ --endpoint-url https://...r2.cloudflarestorage.com
  ```
- [ ] Monitor sync progress
  - Files transferred: _______________
  - Total size: _______________ GB
  - Duration: _______________ minutes

**Method B: Console Upload**
- [ ] Download from S3
  ```bash
  aws s3 sync s3://source-bucket/ ./s3-backup/
  ```
- [ ] Upload to R2 via console
  - Files uploaded: _______________

### Verify Migration
- [ ] Compare file counts
  - S3 files: _______________
  - R2 files: _______________
  - ⬜ Counts match
- [ ] Compare total sizes
  - S3 size: _______________ GB
  - R2 size: _______________ GB
  - ⬜ Sizes match
- [ ] Spot check files
  - [ ] Sample resume accessible
  - [ ] Sample certificate accessible
  - [ ] Sample logo accessible

---

## Phase 4: Deploy on Render

### Environment Configuration
- [ ] Review environment variables template (docs/render.env.template)
- [ ] Add all variables to Render
  - [ ] Database credentials (Neon)
  - [ ] R2 credentials (Cloudflare)
  - [ ] Redis connection (if applicable)
  - [ ] JWT secret
  - [ ] Base URL
  - [ ] CORS origins
  - [ ] Firebase config (if used)
  - [ ] Admin seeding (set to false)
- [ ] Mark secrets as "Secret" in Render UI
- [ ] Double-check all values

### Deploy Application
- [ ] Click "Deploy" in Render dashboard
- [ ] Monitor deployment logs
  - Build status: ⬜ Building | ⬜ Success | ⬜ Failed
  - Deploy status: ⬜ Deploying | ⬜ Live | ⬜ Failed
- [ ] Wait for health check to pass
  - Health endpoint: ⬜ Passing
  - Database health: ⬜ Passing

### Initial Testing
- [ ] Test health endpoint
  ```bash
  curl https://your-service.onrender.com/health
  ```
- [ ] Test database health
  ```bash
  curl https://your-service.onrender.com/health/db
  ```
- [ ] Test authentication
  - [ ] User login works
  - [ ] JWT token generated
- [ ] Test file operations
  - [ ] Upload test file
  - [ ] Verify in R2 bucket
  - [ ] Download test file
- [ ] Test critical flows
  - [ ] Job posting creation
  - [ ] Application submission
  - [ ] User registration

---

## Phase 5: Go Live

### DNS Update (If using custom domain)
- [ ] Add custom domain in Render
  - Domain: _______________
- [ ] Create CNAME record
  - Type: CNAME
  - Name: api (or your subdomain)
  - Value: your-service.onrender.com
  - TTL: 300
- [ ] Wait for SSL certificate (automatic)
  - Status: ⬜ Pending | ⬜ Active
- [ ] Test custom domain
  ```bash
  curl https://api.agriskillacademy.com/health
  ```

### Frontend Update
- [ ] Update frontend API endpoint
  - Old: https://aws-alb-dns-name
  - New: https://your-service.onrender.com
  - Or: https://api.agriskillacademy.com
- [ ] Update CORS origins if needed
- [ ] Deploy frontend changes
- [ ] Test end-to-end flow

### Production Monitoring
- [ ] Monitor Render logs for errors
  - [ ] First 10 minutes: ⬜ No errors
  - [ ] First hour: ⬜ No errors
- [ ] Monitor Neon metrics
  - [ ] Database connections: _______________
  - [ ] Query performance: ⬜ Normal
- [ ] Monitor R2 analytics
  - [ ] File uploads: ⬜ Working
  - [ ] File downloads: ⬜ Working
- [ ] Monitor application metrics
  - [ ] Response times: ⬜ Acceptable
  - [ ] Error rates: ⬜ Low

---

## Phase 6: Cleanup (After 7-14 days)

### Verify Stability
- [ ] No critical issues for 7 days
- [ ] Performance acceptable
- [ ] All features working
- [ ] User feedback positive

### AWS Resource Cleanup
- [ ] Scale down ECS service to 0
  ```bash
  aws ecs update-service --cluster <cluster> --service <service> --desired-count 0
  ```
- [ ] Delete CloudFormation stack
  ```bash
  aws cloudformation delete-stack --stack-name <your-stack-name>
  ```
- [ ] Verify RDS snapshot created
  - Snapshot ID: _______________
  - Snapshot date: _______________
- [ ] Keep S3 files for 30 days (backup)
- [ ] After 30 days: Delete S3 bucket

### Delete Migration Resources
- [ ] Delete CloudFormation migration stack (if used)
  ```bash
  aws cloudformation delete-stack --stack-name agrijobs-rds-migration
  ```
- [ ] Delete local backup files (after verification)

---

## Monitoring & Alerts Setup

### Render Alerts
- [ ] Configure service alerts
  - [ ] Service crash notifications
  - [ ] High CPU/memory alerts
  - [ ] Deployment failure alerts
- [ ] Email: _______________

### Neon Monitoring
- [ ] Review consumption metrics
- [ ] Set up usage alerts (if available)

### Uptime Monitoring
- [ ] Setup external monitoring (UptimeRobot, etc.)
  - [ ] Monitor /health endpoint
  - [ ] Check interval: 5 minutes
  - [ ] Alert email: _______________

---

## Rollback Plan (If Needed)

### Emergency Rollback Steps
- [ ] Revert DNS to AWS ALB
  ```
  CNAME api → <aws-alb-dns-name>
  ```
- [ ] Revert frontend API URL
- [ ] Redeploy frontend
- [ ] Scale up ECS service
  ```bash
  aws ecs update-service --cluster <cluster> --service <service> --desired-count 1
  ```
- [ ] Verify AWS services
  - [ ] RDS accessible
  - [ ] S3 accessible
  - [ ] Redis accessible
- [ ] Monitor for stability

**Rollback Time Estimate:** 5-10 minutes

---

## Notes & Issues

### Issues Encountered
| Date | Issue | Resolution | Status |
|------|-------|------------|--------|
|      |       |            |        |
|      |       |            |        |
|      |       |            |        |

### Important Dates
- Migration started: _______________
- Database migrated: _______________
- Files migrated: _______________
- Went live: _______________
- AWS cleanup: _______________

### Contacts
- Technical lead: _______________
- DevOps support: _______________
- Emergency contact: _______________

---

## Success Criteria

- [ ] Zero data loss (all tables, rows match)
- [ ] All files migrated successfully
- [ ] All critical features working
- [ ] Performance acceptable or better
- [ ] Cost reduced by >60%
- [ ] No critical issues for 7 days
- [ ] User feedback positive
- [ ] AWS resources cleaned up

---

**Migration Completed:** ⬜ Yes | ⬜ No

**Sign-off:**
- Name: _______________
- Date: _______________
- Signature: _______________
