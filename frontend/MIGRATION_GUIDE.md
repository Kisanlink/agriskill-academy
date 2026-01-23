# Migration Guide: AWS to Render (Neon DB + Cloudflare R2)

This guide details how to migrate your **Private AWS RDS (PostgreSQL)** to **Neon DB** and your **AWS S3** data to **Cloudflare R2** without changing your application code.

## 🛑 Important Constraints

1.  **Private RDS:** Your database is in a **Private Subnet**. You cannot connect to it directly from your local computer or the internet. You must use a "Jump Box" (Bastion Host) or an SSH Tunnel to reach it.
2.  **S3 Permissions:** You will need an AWS IAM User with `AmazonS3FullAccess` (or at least Read permissions) to export files.

---

## 🛠 Prerequisites

Ensure you have these installed on your local machine (Windows):

1.  **AWS CLI**: [Install Guide](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
2.  **PostgreSQL Tools** (`pg_dump`, `psql`): [Download Installer](https://www.postgresql.org/download/windows/)
    *   *Note: Add the `bin` folder (e.g., `C:\Program Files\PostgreSQL\16\bin`) to your System PATH.*
3.  **Rclone** (for file migration): [Download](https://rclone.org/downloads/)
    *   *Or use Cloudflare's web-based importer (easiest).*

---

## Part 1: Database Migration (Private RDS → Neon DB)

Since your RDS is private, we will create a temporary "bridge" (EC2 Bastion Host) to access it.

### Step 1: Create a Temporary Bastion Host
1.  Log in to the **AWS Console** > **EC2**.
2.  Click **Launch Instance**.
3.  **Name:** `Migration-Bastion`.
4.  **AMI:** Amazon Linux 2023 (Free tier eligible).
5.  **Instance Type:** `t2.micro` or `t3.micro`.
6.  **Key Pair:** Create a new key pair (e.g., `migration-key.pem`) and **download it**.
7.  **Network Settings (Crucial)**:
    *   **VPC:** Select the **AgriJobs VPC** (the one where your RDS lives).
    *   **Subnet:** Select a **PUBLIC Subnet** (e.g., `agrijobs-public-1`).
    *   **Auto-assign Public IP:** **Enable**.
    *   **Security Group:** Create a new one named `Bastion-SG`. Allow **SSH (Port 22)** from **My IP**.
8.  Launch the instance.
9.  Copy the **Public IPv4 address** of this new instance.

### Step 2: Allow Bastion to Access RDS
1.  Go to **RDS Console** > **Databases** > Click your DB instance.
2.  Find the **VPC security groups** (under Connectivity & security). Click the active security group (e.g., `rds-sg`).
3.  Go to the **Inbound rules** tab > **Edit inbound rules**.
4.  **Add Rule**:
    *   **Type:** PostgreSQL (5432)
    *   **Source:** Select "Custom" and start typing `Bastion-SG` (select the security group you created in Step 1).
5.  **Save rules**.

### Step 3: Open an SSH Tunnel
On your local machine (PowerShell or Command Prompt):

1.  Move your downloaded key (`migration-key.pem`) to a safe folder (e.g., `C:\Users\Karthikeya\.ssh\`).
2.  Run this command to create a tunnel. This forwards your local port `5433` to the remote RDS port `5432`.

```powershell
# Syntax: ssh -i <path-to-key> -L <local-port>:<rds-endpoint>:5432 ec2-user@<bastion-public-ip> -N

ssh -i "C:\Users\Karthikeya\.ssh\migration-key.pem" -L 5433:agrijobs-db.xxxxxx.us-east-1.rds.amazonaws.com:5432 ec2-user@54.123.45.67
```

*   *Replace the RDS endpoint and Bastion IP with your actual values.*
*   *The command will appear to "hang" or sit silently. This is normal. **Keep this window open.** Result: accessing `localhost:5433` now sends traffic to your AWS RDS.*

### Step 4: Dump AWS Data
Open a **new** terminal window.

1.  Run `pg_dump` pointing to your local tunnel.

```powershell
# -h localhost -p 5433 targets the tunnel
# -U <db_username> is your AWS DB username (e.g., agrijobs)
# -d <db_name> is your AWS DB name (e.g., asa_db)

pg_dump -h localhost -p 5433 -U agrijobs -d asa_db --no-owner --no-acl --clean --if-exists -f dump.sql
```
*   You will be asked for a password. Enter your **AWS RDS Password**.

### Step 5: Restore to Neon DB
1.  Get your **Neon Connection String** from the Neon Dashboard (select "Pooled connection" if available).
    *   Format: `postgres://user:pass@ep-xyz.region.neon.tech/dbname?sslmode=require`
2.  Run `psql` to import the data.

```powershell
# Replace the URL below with your Neon connection string
psql "postgres://user:pass@ep-xyz.region.neon.tech/dbname?sslmode=require" -f dump.sql
```

3.  **Verification**:
    *   Run: `psql "postgres://..." -c "\dt"` to list tables and ensure they exist.

---

## Part 2: File Storage Migration (S3 → R2)

### Option A: Cloudflare "Super Slurper" (Recommended)
This runs entirely on Cloudflare's servers. No local bandwidth used.

1.  Log in to **Cloudflare Dashboard** > **R2**.
2.  Go to **Data Migration** (on the sidebar).
3.  Click **Migrate files**.
4.  **Source Bucket (AWS)**:
    *   Bucket Name: `production-agrijobs-files-xxxxxxxx`
    *   Region: `us-east-1` (or your specific region).
    *   **Credentials:** You need an AWS Access Key ID and Secret Access Key. (Create an IAM User with `AmazonS3ReadOnlyAccess` if you don't have one).
5.  **Destination Bucket (R2)**:
    *   Select or create your new R2 bucket (e.g., `agrijobs-files`).
6.  Start Migration.

### Option B: Using Rclone (CLI)
If you prefer manual control or Super Slurper fails.

1.  Run `rclone config` and follow the prompts:
    *   **New Remote `aws`:** Type `s3`, Provider `AWS`, paste AWS Keys, Region `us-east-1`.
    *   **New Remote `r2`:** Type `s3`, Provider `Cloudflare`, paste R2 Access Key/Secret, Endpoint (from R2 dashboard).
2.  Sync:
    ```powershell
    rclone sync aws:production-agrijobs-files-12345 r2:agrijobs-files --progress --transfers=32
    ```

---

## Part 3: Cleanup

Once verification is done:
1.  **Terminate the Bastion Host** in AWS EC2 (Stop/Terminate instance) to stop paying for it.
2.  **Remove the Inbound Rule** from your RDS Security Group (delete the rule allowing access from Bastion-SG).

```