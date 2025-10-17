# Firebase Email Service Setup Guide

## Overview

AgriJobs uses **Firebase ONLY for sending emails** (verification and password reset). All authentication, password storage, and user management remain 100% local.

## Architecture

```
User Signup → Local DB (bcrypt password) → Firebase Email Service → Verification Email
User clicks link → Verify token in Local DB → Update email_verified = true

User Forgot Password → Generate token in Local DB → Firebase Email Service → Reset Email
User clicks link → Validate token in Local DB → Update password in Local DB
```

**Key Points:**
- ✅ Passwords stored ONLY in local database (bcrypt hashed)
- ✅ Login validates against local database (no Firebase call)
- ✅ Firebase only sends emails (verification & reset links)
- ✅ Tokens generated and managed locally
- ✅ No user data sync with Firebase needed

---

## Firebase Project Setup

### Step 1: Create Firebase Project

1. Go to [Firebase Console](https://console.firebase.google.com)
2. Click **Add project**
3. Enter project name: `agrijobs-production` (or your preferred name)
4. Disable Google Analytics (optional for email service)
5. Click **Create project**

### Step 2: Enable Authentication

1. In Firebase Console, go to **Authentication** → **Get started**
2. Click **Sign-in method** tab
3. Click **Email/Password**
4. **Enable** the Email/Password provider
5. Click **Save**

**Note:** We only need Firebase Auth enabled for email sending API access. Users are NOT authenticated via Firebase.

### Step 3: Configure Email Templates

1. Go to **Authentication** → **Templates**
2. Customize the following templates:

#### Email Address Verification Template
```
Subject: Verify your email for AgriJobs

Body:
Hello {{userName}},

Thank you for registering with AgriJobs!

Please verify your email address by clicking the link below:
{{verificationLink}}

This link will expire in 24 hours.

If you didn't create an account, please ignore this email.

Best regards,
The AgriJobs Team
```

#### Password Reset Template
```
Subject: Reset your AgriJobs password

Body:
Hello,

We received a request to reset your password for your AgriJobs account.

Click the link below to reset your password:
{{resetLink}}

This link will expire in 1 hour.

If you didn't request a password reset, please ignore this email.

Best regards,
The AgriJobs Team
```

### Step 4: Generate Service Account Credentials

1. Go to **Project Settings** (gear icon) → **Service accounts**
2. Click **Generate new private key**
3. Click **Generate key**
4. Save the downloaded JSON file as `serviceAccountKey.json`
5. **Keep this file secure** - it contains sensitive credentials

---

## Local Development Setup

### 1. Copy Service Account Key

```bash
# Place the downloaded key in your project root
cp ~/Downloads/serviceAccountKey.json ./serviceAccountKey.json

# Add to .gitignore (IMPORTANT!)
echo "serviceAccountKey.json" >> .gitignore
```

### 2. Update .env File

```bash
# Copy example file
cp .env.example .env

# Edit .env and add Firebase configuration
FIREBASE_PROJECT_ID=agrijobs-production
FIREBASE_CREDENTIALS_PATH=./serviceAccountKey.json
FRONTEND_URL=http://localhost:3000
```

### 3. Install Go Dependencies

```bash
go get firebase.google.com/go/v4
go get google.golang.org/api/option
```

### 4. Test Email Sending

Start your application and test signup:

```bash
# Start the server
go run cmd/server/main.go

# In another terminal, test signup
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "user_name": "testuser",
    "email": "lisise3548@datoinf.com",
    "password": "SecurePass123!",
    "confirm_password": "SecurePass123!",
    "role": "student",
    "phone_number": "1234567890"
  }'
```

Check the email inbox for `test@example.com` - you should receive a verification email.

---

## Production Deployment (AWS ECS)

### 1. Base64 Encode Service Account JSON

```bash
# Encode the JSON file
base64 -w 0 serviceAccountKey.json > firebase-credentials-base64.txt

# Copy the output
cat firebase-credentials-base64.txt
```

### 2. Store in AWS Secrets Manager

```bash
# Create secret in AWS Secrets Manager
aws secretsmanager create-secret \
  --name agrijobs/production/firebase-credentials \
  --description "Firebase service account credentials for email sending" \
  --secret-string "$(cat firebase-credentials-base64.txt)" \
  --region us-east-1
```

### 3. Update CloudFormation Template

Add to `cloudformation/infrastructure.yml`:

```yaml
# Add to Parameters section
FirebaseProjectID:
  Type: String
  Description: Firebase project ID
  Default: agrijobs-production

FrontendURL:
  Type: String
  Description: Frontend URL for email links
  Default: https://agrijobs.com

# Add to Resources section
FirebaseCredentialsSecret:
  Type: AWS::SecretsManager::Secret
  Properties:
    Name: !Sub "${EnvironmentName}-firebase-credentials"
    Description: Firebase service account credentials (base64)
    SecretString: !Ref FirebaseCredentialsBase64

# Add to ECS Task Definition environment variables
- Name: FIREBASE_PROJECT_ID
  Value: !Ref FirebaseProjectID
- Name: FIREBASE_CREDENTIALS_JSON
  ValueFrom: !Ref FirebaseCredentialsSecret
- Name: FRONTEND_URL
  Value: !Ref FrontendURL
```

### 4. Deploy

```bash
# Deploy CloudFormation stack with Firebase parameters
aws cloudformation create-stack \
  --stack-name agrijobs-production \
  --template-body file://cloudformation/infrastructure.yml \
  --parameters \
    ParameterKey=FirebaseProjectID,ParameterValue=agrijobs-production \
    ParameterKey=FrontendURL,ParameterValue=https://agrijobs.com \
    # ... other parameters
  --capabilities CAPABILITY_NAMED_IAM \
  --region us-east-1
```

---

## Docker Compose Setup

### 1. Update docker-compose.production.yml

```yaml
services:
  backend:
    environment:
      # Firebase Configuration
      FIREBASE_PROJECT_ID: ${FIREBASE_PROJECT_ID}
      FIREBASE_CREDENTIALS_JSON: ${FIREBASE_CREDENTIALS_JSON}
      FRONTEND_URL: ${FRONTEND_URL:-http://localhost:3000}
```

### 2. Update .env

```bash
FIREBASE_PROJECT_ID=agrijobs-production
FIREBASE_CREDENTIALS_JSON=<base64_encoded_json_here>
FRONTEND_URL=http://localhost:3000
```

### 3. Run

```bash
docker-compose -f docker-compose.production.yml up -d
```

---

## API Endpoints

### 1. User Signup (with Email Verification)

```bash
POST /api/auth/signup
Content-Type: application/json

{
  "name": "John Doe",
  "user_name": "johndoe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "confirm_password": "SecurePass123!",
  "role": "student",
  "phone_number": "1234567890"
}

Response:
{
  "success": true,
  "message": "User registered successfully. Please check your email to verify your account.",
  "user": {
    "id": "USER_abc123",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "student",
    "email_verified": false
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. Verify Email

```bash
GET /api/auth/verify-email?token=<verification_token>

Response:
{
  "success": true,
  "message": "Email verified successfully! You can now login."
}
```

### 3. Forgot Password

```bash
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}

Response:
{
  "success": true,
  "message": "Password reset link sent to your email"
}
```

### 4. Reset Password

```bash
POST /api/auth/reset-password
Content-Type: application/json

{
  "token": "<reset_token>",
  "new_password": "NewSecurePass123!"
}

Response:
{
  "success": true,
  "message": "Password reset successful"
}
```

### 5. Login (100% Local - No Firebase)

```bash
POST /api/auth/login
Content-Type: application/json

{
  "user_name": "johndoe",
  "password": "SecurePass123!"
}

Response:
{
  "success": true,
  "message": "Login successful",
  "user": {
    "id": "USER_abc123",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "student"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## Troubleshooting

### Issue 1: Firebase initialization fails

**Error:**
```
Failed to initialize Firebase email service: credentials not provided
```

**Solution:**
- Check that `FIREBASE_CREDENTIALS_PATH` or `FIREBASE_CREDENTIALS_JSON` is set
- Verify the service account JSON file exists and is valid
- Ensure the JSON is properly base64 encoded for container deployments

### Issue 2: Email not received

**Possible causes:**
1. **Firebase Auth not enabled**: Go to Firebase Console → Authentication → Enable Email/Password
2. **Invalid email address**: Check that the email is valid and not a temporary/disposable email
3. **Spam folder**: Check user's spam/junk folder
4. **Rate limiting**: Firebase has rate limits on emails sent

**Debug:**
```bash
# Check application logs
docker logs agrijobs-backend | grep "Firebase"
docker logs agrijobs-backend | grep "verification email"
```

### Issue 3: Token invalid or expired

**Error:**
```
{
  "success": false,
  "message": "invalid or expired verification token"
}
```

**Causes:**
- Token was already used
- Token doesn't exist in database
- User clicked wrong link

**Solution:**
- User needs to request a new verification email (resend feature)
- Check database for `verification_token` field

### Issue 4: Firebase permission denied

**Error:**
```
Permission denied (Service: Firebase Auth)
```

**Solution:**
- Verify the service account has correct permissions
- Regenerate service account key from Firebase Console
- Ensure Firebase Authentication is enabled

---

## Security Best Practices

1. ✅ **Never commit service account credentials** to git
   ```bash
   echo "serviceAccountKey.json" >> .gitignore
   echo "firebase-credentials-base64.txt" >> .gitignore
   ```

2. ✅ **Rotate credentials periodically**
   - Regenerate service account key every 90 days
   - Update in Secrets Manager immediately

3. ✅ **Use environment-specific projects**
   - Development: `agrijobs-dev`
   - Staging: `agrijobs-staging`
   - Production: `agrijobs-production`

4. ✅ **Monitor Firebase usage**
   - Set up billing alerts in Firebase Console
   - Monitor daily email sending limits

5. ✅ **Restrict Firebase API access**
   - Only enable required APIs (Authentication)
   - Disable unused features to reduce attack surface

---

## Cost Estimation

Firebase Authentication pricing:
- **Free tier**: 50 emails/day
- **Blaze (Pay-as-you-go)**: $0.00 for authentication
- **Email sending**: No additional cost for verification/reset emails

**Estimated cost for AgriJobs:**
- Development: FREE (< 50 emails/day)
- Production (1000 signups/day): FREE (email sending is free)

**Note:** You only pay for Firebase if you exceed free tier limits on other services (Firestore, Storage, etc.). Since we only use Authentication for email sending, costs should be $0/month.

---

## Support

For Firebase-related issues:
- Firebase Documentation: https://firebase.google.com/docs/auth
- Firebase Support: https://firebase.google.com/support

For AgriJobs-specific issues:
- GitHub Issues: https://github.com/Kisanlink/agriskill-academy/issues
- Technical Architecture: See `TECHNICAL_ARCHITECTURE.md`

---

## License

Copyright © 2025 KisanLink
