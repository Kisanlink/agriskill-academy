# Docker Port Mapping Guide

## Understanding Docker Networking

### Port Mapping Format
```yaml
ports:
  - "HOST_PORT:CONTAINER_PORT"
```

- **HOST_PORT**: Port on your computer (localhost)
- **CONTAINER_PORT**: Port inside the container

## Current Configuration

### PostgreSQL
```yaml
ports:
  - "5433:5432"
```

**What this means:**
- ✅ **Inside Docker network**: App connects to `postgres:5432`
- ✅ **From your computer**: Use `localhost:5433` (for pgAdmin, DBeaver, etc.)
- ✅ PostgreSQL runs on its default port 5432 inside the container
- ✅ Port 5433 on your host maps to container's 5432

**Important:** You cannot change PostgreSQL's internal port easily. It always runs on 5432 inside the container.

### Redis
```yaml
ports:
  - "6379:6379"
```

**What this means:**
- ✅ **Inside Docker network**: App connects to `redis:6379`
- ✅ **From your computer**: Use `localhost:6379` (for Redis CLI, Redis Desktop Manager, etc.)

### Application Server
```yaml
ports:
  - "8080:8080"
```

**What this means:**
- ✅ **From your computer**: Access API at `http://localhost:8080`
- ✅ App listens on port 8080 inside the container

## Fixed Configuration

### What Was Wrong
```yaml
# ❌ WRONG - This tried to expose container port 5433, but PostgreSQL runs on 5432
ports:
  - "${DB_PORT:-5432}:5433"

# ❌ WRONG - App tried to connect to postgres:5433
environment:
  DB_PORT: ${DB_PORT:-5432}  # If you set DB_PORT=5433 in .env
```

### What Is Now Correct
```yaml
# ✅ CORRECT - Maps host 5433 to container 5432 (where PostgreSQL actually runs)
ports:
  - "5433:5432"

# ✅ CORRECT - App always connects to container's internal port
environment:
  DB_PORT: 5432  # Fixed value for Docker network communication
```

## How to Connect

### From Your Application (Inside Docker)
```yaml
DB_HOST: postgres
DB_PORT: 5432  # Always 5432 (container internal port)
```

### From Your Computer (Outside Docker)
```bash
# PostgreSQL
psql -h localhost -p 5433 -U agrijobs -d asa_db

# Redis
redis-cli -h localhost -p 6379

# API
curl http://localhost:8080/health
```

### From Tools (Outside Docker)
**pgAdmin / DBeaver / TablePlus:**
- Host: `localhost`
- Port: `5433` ← Use this, not 5432!
- Username: Your POSTGRES_USER
- Password: Your POSTGRES_PASS
- Database: Your DB_NAME

**Redis Desktop Manager:**
- Host: `localhost`
- Port: `6379`
- Password: Your REDIS_PASSWORD

## Important Rules

### ✅ DO:
1. Use hardcoded internal ports for container-to-container communication
2. Map host ports if you need external access
3. Remember: container services talk to each other using service names and internal ports

### ❌ DON'T:
1. Try to change PostgreSQL's or Redis's internal container port
2. Use ${DB_PORT} or ${REDIS_PORT} for Docker container environment variables
3. Expect containers to connect using host ports

## Troubleshooting

### Connection Refused Error
```
dial tcp 172.18.0.2:5433: connect: connection refused
```

**Problem:** App is trying to connect to wrong port inside Docker network

**Solution:**
- App should use `DB_PORT: 5432` (container internal)
- NOT `DB_PORT: 5433` (host port)

### Can't Connect from pgAdmin
**Problem:** Using port 5432 from your computer

**Solution:** Use port 5433 from host:
```
psql -h localhost -p 5433 -U agrijobs -d asa_db
```

### Port Already in Use
```
Error: bind: address already in use
```

**Problem:** Another service is using the host port

**Solution:** Change the HOST port (left side):
```yaml
# If port 5433 is busy, use a different host port
ports:
  - "5434:5432"  # Now access from localhost:5434
```

## Quick Reference Table

| Service | Internal Port | Host Port | Container-to-Container | Host-to-Container |
|---------|--------------|-----------|----------------------|-------------------|
| PostgreSQL | 5432 | 5433 | `postgres:5432` | `localhost:5433` |
| Redis | 6379 | 6379 | `redis:6379` | `localhost:6379` |
| App | 8080 | 8080 | `app:8080` | `localhost:8080` |

## Environment Variables in .env

**You do NOT need to set these in .env for Docker Compose:**
- ❌ `DB_PORT` - Hardcoded in docker-compose.yml as 5432
- ❌ `REDIS_PORT` - Hardcoded in docker-compose.yml as 6379

**You DO need to set these:**
- ✅ `POSTGRES_USER`
- ✅ `POSTGRES_PASS`
- ✅ `DB_NAME`
- ✅ `REDIS_PASSWORD`
- ✅ `JWT_SECRET`
- ✅ All other application configuration

## Testing Your Setup

```bash
# 1. Start services
docker-compose -f docker-compose.production.yml up -d

# 2. Check if all services are healthy
docker-compose -f docker-compose.production.yml ps

# 3. Test PostgreSQL from host
docker-compose -f docker-compose.production.yml exec postgres psql -U agrijobs -d asa_db -c "SELECT version();"

# 4. Test app can connect to PostgreSQL
docker-compose -f docker-compose.production.yml logs app | grep "DB_PORT"

# 5. Test API
curl http://localhost:8080/health
```

## Summary

**The key insight:** Containers have TWO sets of ports:
1. **Internal ports** (for container-to-container communication) - Fixed by the container image
2. **External/Host ports** (for accessing from your computer) - You can customize these

Your app inside Docker MUST use internal ports. You can map any host port to access from outside.
