# Gorev Docker Guide

## Quick Start

### Basic Docker Run

```bash
docker run -d \
  --name gorev \
  -p 5082:5082 \
  -v $(pwd)/workspace:/workspace \
  -e GOREV_LANG=tr \
  msenol/gorev:latest \
  daemon --detach
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'
services:
  gorev:
    image: msenol/gorev:latest
    ports:
      - "5082:5082"
    volumes:
      - ./workspace:/workspace
    environment:
      - GOREV_LANG=tr
    command: daemon --detach
```

Start:
```bash
docker-compose up -d
```

### With Database Persistence

```yaml
version: '3.8'
services:
  gorev:
    image: msenol/gorev:latest
    ports:
      - "5082:5082"
    volumes:
      - ./workspace:/workspace
      - ./data:/data
    environment:
      - GOREV_LANG=tr
      - GOREV_DB_PATH=/data/gorev.db
    command: daemon --detach
```

## Advanced Usage

### Custom Configuration

```yaml
version: '3.8'
services:
  gorev:
    image: msenol/gorev:latest
    ports:
      - "5082:5082"
    volumes:
      - ./workspace:/workspace
      - ./config:/config
    environment:
      - GOREV_LANG=en
      - GOREV_DB_PATH=/config/gorev.db
      - GOREV_CONFIG_PATH=/config/gorev.toml
    command: daemon --detach --config /config/gorev.toml
```

### Multi-Workspace Setup

```yaml
version: '3.8'
services:
  gorev-project1:
    image: msenol/gorev:latest
    ports:
      - "5082:5082"
    volumes:
      - ./project1:/workspace
    environment:
      - GOREV_LANG=tr
      - GOREV_DB_PATH=/workspace/.gorev/gorev.db
    command: daemon --detach

  gorev-project2:
    image: msenol/gorev:latest
    ports:
      - "5083:5082"
    volumes:
      - ./project2:/workspace
    environment:
      - GOREV_LANG=tr
      - GOREV_DB_PATH=/workspace/.gorev/gorev.db
    command: daemon --detach
```

### Development Mode

```yaml
version: '3.8'
services:
  gorev:
    image: msenol/gorev:latest
    ports:
      - "5082:5082"
    volumes:
      - ./workspace:/workspace
      - ./logs:/logs
    environment:
      - GOREV_LANG=tr
      - GOREV_DEBUG=true
    command: daemon --detach --debug --log-dir /logs
```

## Environment Variables

- `GOREV_LANG`: UI language (`tr` or `en`)
- `GOREV_DB_PATH`: Database file path
- `GOREV_CONFIG_PATH`: Config file path
- `GOREV_DEBUG`: Enable debug logging (`true`/`false`)
- `GOREV_LOG_DIR`: Log directory path

## Volumes

- `/workspace`: Project directory
- `/data`: Database storage
- `/config`: Configuration files
- `/logs`: Log files

## Health Check

```bash
curl http://localhost:5082/api/health
```

Docker health check in compose:
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:5082/api/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 10s
```

## Troubleshooting

### Container not starting

Check logs:
```bash
docker logs gorev
docker-compose logs gorev
```

### Port already in use

Change port mapping:
```yaml
ports:
  - "5083:5082"  # Host:Container
```

### Permission issues

Fix volume permissions:
```bash
mkdir -p workspace data logs
chmod 755 workspace data logs
```

### Database locked

Remove lock file:
```bash
docker exec gorev rm -f /workspace/.gorev/gorev.db-shm
docker exec gorev rm -f /workspace/.gorev/gorev.db-wal
```

## Production Considerations

### Security

- Use reverse proxy (nginx/traefik)
- Enable authentication
- Use HTTPS
- Restrict network access

### Backup

Backup database volume:
```bash
docker run --rm \
  -v gorev_data:/data \
  -v $(pwd):/backup \
  busybox tar cvf /backup/gorev-backup.tar /data
```

### Monitoring

```yaml
version: '3.8'
services:
  gorev:
    image: msenol/gorev:latest
    ports:
      - "5082:5082"
    volumes:
      - ./workspace:/workspace
      - ./logs:/logs
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

