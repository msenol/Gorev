# Gorev MCP Configuration Guide

This guide explains how to configure Gorev MCP server for different deployment scenarios.

## Connection Modes

Gorev MCP supports three connection modes:

### 1. Local Installation (npm/npx)

**Best for**: Individual developers, quick testing

#### Setup:
```bash
# Install globally
npm install -g @mehmetsenol/gorev-mcp-server

# Or use with npx (no installation)
npx @mehmetsenol/gorev-mcp-server daemon --detach
```

#### Configuration (`.mcp.json`):
```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "mcp-proxy"],
      "env": {
        "GOREV_MCP_CONNECTION_MODE": "local",
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

### 2. Docker Container

**Best for**: Team deployments, CI/CD, consistent environments

#### Setup:
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

#### Configuration (`.mcp.json`):
```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "mcp-proxy"],
      "env": {
        "GOREV_MCP_CONNECTION_MODE": "docker",
        "GOREV_DOCKER_COMPOSE_FILE": "./docker-compose.yml",
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

### 3. Remote Server

**Best for**: Centralized deployments, multiple clients

#### Server Setup:
```bash
# On remote server
npm install -g @mehmetsenol/gorev-mcp-server
gorev daemon --detach --port 5082
```

#### Configuration (`.mcp.json`):
```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "mcp-proxy"],
      "env": {
        "GOREV_MCP_CONNECTION_MODE": "remote",
        "GOREV_API_HOST": "your-server.com",
        "GOREV_API_PORT": "5082",
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

## Environment Variables

| Variable | Description | Default | Modes |
|----------|-------------|---------|-------|
| `GOREV_MCP_CONNECTION_MODE` | Connection mode: `auto`, `local`, `docker`, `remote` | `auto` | All |
| `GOREV_API_HOST` | API server host | `localhost` | remote |
| `GOREV_API_PORT` | API server port | `5082` | docker, remote |
| `GOREV_SERVER_PATH` | Path to gorev binary | `gorev` (from PATH) | local |
| `GOREV_DOCKER_COMPOSE_FILE` | Path to docker-compose.yml | `./docker-compose.yml` | docker |
| `GOREV_LANG` | UI language (`tr` or `en`) | `tr` | All |
| `GOREV_DB_PATH` | Database path (for daemon) | Auto-detect | All |

## Troubleshooting

### Daemon not starting

Check daemon status:
```bash
curl http://localhost:5082/api/health
```

If not running, manually start:
```bash
gorev daemon --detach --debug
```

### Docker issues

Check docker logs:
```bash
docker-compose logs gorev
```

### Remote connection issues

Verify connectivity:
```bash
curl http://your-server.com:5082/api/health
```

### MCP proxy logs

Enable debug mode:
```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "mcp-proxy", "--debug"]
    }
  }
}
```

## Best Practices

1. **Development**: Use `local` or `npx` mode
2. **Team**: Use `docker` mode with shared compose file
3. **Production**: Use `remote` mode with proper server
4. **CI/CD**: Use `docker` mode in pipelines
5. **Always**: Set `GOREV_LANG` to your preferred language

## Examples

### Multiple MCP Clients

Same `.mcp.json` works for:
- VS Code MCP Extension
- Claude Desktop
- Kilocode MCP Client
- Windsurf
- Cursor

### Per-Project Configuration

Place `.mcp.json` in project root, Gorev will auto-detect workspace.

### Environment-Specific

Use `.mcp.json` for local, environment variables for CI/CD:
```bash
export GOREV_MCP_CONNECTION_MODE=docker
export GOREV_DOCKER_COMPOSE_FILE=./docker-compose.prod.yml
npx @mehmetsenol/gorev-mcp-server mcp-proxy
```

