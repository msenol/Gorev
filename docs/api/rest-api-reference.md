# REST API Reference

**Last Updated:** 2 December 2025
**API Version:** v1
**Base URL:** `http://localhost:5082/api/v1`
**Framework:** Fiber (Go)

## 📋 Overview

The Gorev REST API provides a complete HTTP interface for task management operations. It's designed for the embedded Web UI but can also be used by external applications.

## 🔐 Authentication

Currently, the API does not require authentication as it's designed for local development use only.

**CORS Policy:**

- Allowed Origins: `http://localhost:5000-5003` (development only)
- Allowed Methods: `GET, POST, PUT, DELETE, OPTIONS`
- Credentials: Not allowed

## 📊 Response Format

All API responses follow a consistent JSON structure:

### Success Response

```json
{
  "success": true,
  "data": { /* response data */ },
  "total": 42,           // Optional: for list endpoints
  "message": "Success"   // Optional: for mutations
}
```

### Error Response

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

## 🛠️ Endpoints

### Health Check

#### GET `/api/v1/health`

Check if the API server is running.

**Response:**

```json
{
  "status": "ok",
  "time": 1727654321
}
```

---

### Tasks

#### GET `/api/v1/tasks`

List tasks with optional filtering.

**Query Parameters:**

| Parameter | Type | Required | Description | Default |
|-----------|------|----------|-------------|---------|
| `durum` | string | ❌ | Filter by status: `beklemede`, `devam_ediyor`, `tamamlandi` | All |
| `tum_projeler` | boolean | ❌ | Include tasks from all projects | `false` |
| `sirala` | string | ❌ | Sort by: `son_tarih_asc`, `son_tarih_desc` | - |
| `filtre` | string | ❌ | Time filter: `acil` (within 7 days), `gecmis` (overdue) | - |
| `etiket` | string | ❌ | Filter by tag name | - |
| `limit` | number | ❌ | Maximum tasks to return | `50` |
| `offset` | number | ❌ | Number of tasks to skip | `0` |

**Example Request:**

```bash
curl "http://localhost:5082/api/v1/tasks?durum=devam_ediyor&limit=10"
```

**Example Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "baslik": "JWT Authentication",
      "aciklama": "Token based auth system",
      "durum": "devam_ediyor",
      "oncelik": "yuksek",
      "proje_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "parent_id": null,
      "olusturma_zamani": "2025-09-15T14:30:00Z",
      "guncelleme_zamani": "2025-09-16T10:45:00Z",
      "son_tarih": null,
      "etiketler": ["security", "backend"],
      "alt_gorevler": [],
      "bagimliliklar": []
    }
  ],
  "total": 1
}
```

#### GET `/api/v1/tasks/:id`

Get detailed information about a specific task.

**Path Parameters:**

- `id` (string, required): Task UUID

**Example Request:**

```bash
curl "http://localhost:5082/api/v1/tasks/550e8400-e29b-41d4-a716-446655440000"
```

**Example Response:**

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "baslik": "JWT Authentication",
    "aciklama": "## Implementation\n\n- Spring Security\n- JWT tokens\n- Refresh mechanism",
    "durum": "devam_ediyor",
    "oncelik": "yuksek",
    "proje_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "proje_isim": "E-ticaret Sitesi",
    "parent_id": null,
    "olusturma_zamani": "2025-09-15T14:30:00Z",
    "guncelleme_zamani": "2025-09-16T10:45:00Z",
    "son_tarih": "2025-10-01T00:00:00Z",
    "etiketler": ["security", "backend", "urgent"],
    "alt_gorevler": [
      {
        "id": "child-task-id",
        "baslik": "Setup JWT library",
        "durum": "tamamlandi"
      }
    ],
    "bagimliliklar": [
      {
        "kaynak_id": "dependency-task-id",
        "hedef_id": "550e8400-e29b-41d4-a716-446655440000",
        "baglanti_tipi": "onceki"
      }
    ]
  }
}
```

#### POST `/api/v1/tasks`

**⚠️ DEPRECATED**: Direct task creation is deprecated. Use `/api/v1/tasks/from-template` instead.

**Response:**

```json
{
  "success": false,
  "error": "Direct task creation is deprecated. Use /tasks/from-template endpoint with a template."
}
```

#### POST `/api/v1/tasks/from-template`

Create a new task from a template.

**Request Body:**

```json
{
  "template_id": "39f28dbd-10f3-454c-8b35-52ae6b7ea391",
  "degerler": {
    "baslik": "Login button not working",
    "aciklama": "Users cannot click login button",
    "modul": "auth",
    "ortam": "production",
    "adimlar": "1. Go to login page\n2. Enter credentials\n3. Click login",
    "beklenen": "User should be redirected to dashboard",
    "mevcut": "Nothing happens, button is unresponsive",
    "oncelik": "yuksek",
    "etiketler": "bug,urgent,auth"
  }
}
```

**Example Response:**

```json
{
  "success": true,
  "data": {
    "id": "d7f4e8b9-2a1c-4f5e-9d3b-8c1a2e3f4d5b",
    "baslik": "🐛 [auth] Login button not working",
    "durum": "beklemede",
    "oncelik": "yuksek"
  },
  "message": "Task created from template successfully"
}
```

#### PUT `/api/v1/tasks/:id`

Update an existing task.

**Path Parameters:**

- `id` (string, required): Task UUID

**Request Body:**

```json
{
  "durum": "tamamlandi"
}
```

**Supported Fields:**

- `durum`: Status (beklemede, devam_ediyor, tamamlandi)
- `oncelik`: Priority (dusuk, orta, yuksek)
- `baslik`: Title (string)
- `aciklama`: Description (markdown string)

**Example Response:**

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "baslik": "JWT Authentication",
    "durum": "tamamlandi"
  },
  "message": "Task updated successfully"
}
```

#### DELETE `/api/v1/tasks/:id`

Delete a task permanently.

**Path Parameters:**

- `id` (string, required): Task UUID

**Example Response:**

```json
{
  "success": true,
  "message": "Task deleted successfully"
}
```

---

### Projects

#### GET `/api/v1/projects`

List all projects with task counts.

**Example Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "isim": "E-ticaret Sitesi",
      "tanim": "Online satış platformu geliştirme projesi",
      "olusturma_zamani": "2025-09-01T10:00:00Z",
      "gorev_sayisi": 12
    },
    {
      "id": "6ba7b814-9dad-11d1-80b4-00c04fd430c8",
      "isim": "Mobil Uygulama v2.0",
      "tanim": "React Native ile cross-platform mobil uygulama",
      "olusturma_zamani": "2025-09-10T14:30:00Z",
      "gorev_sayisi": 8
    }
  ],
  "total": 2
}
```

#### GET `/api/v1/projects/:id`

Get detailed information about a specific project.

**Path Parameters:**

- `id` (string, required): Project UUID

**Example Response:**

```json
{
  "success": true,
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "isim": "E-ticaret Sitesi",
    "tanim": "Online satış platformu geliştirme projesi. Kullanıcı yönetimi, ürün katalogu, sepet ve ödeme sistemlerini içerir.",
    "olusturma_zamani": "2025-09-01T10:00:00Z",
    "gorev_sayisi": 12
  }
}
```

#### POST `/api/v1/projects`

Create a new project.

**Request Body:**

```json
{
  "isim": "Mobile App v3.0",
  "tanim": "Next generation mobile application with AI features"
}
```

**Example Response:**

```json
{
  "success": true,
  "data": {
    "id": "new-project-uuid",
    "isim": "Mobile App v3.0",
    "tanim": "Next generation mobile application with AI features",
    "olusturma_zamani": "2025-09-30T00:00:00Z",
    "gorev_sayisi": 0
  },
  "message": "Project created successfully"
}
```

#### GET `/api/v1/projects/:id/tasks`

Get all tasks for a specific project.

**Path Parameters:**

- `id` (string, required): Project UUID

**Query Parameters:**

- `limit` (number, optional): Max tasks to return (default: 50)
- `offset` (number, optional): Number of tasks to skip (default: 0)

**Example Request:**

```bash
curl "http://localhost:5082/api/v1/projects/6ba7b810-9dad-11d1-80b4-00c04fd430c8/tasks?limit=10"
```

**Example Response:**

```json
{
  "success": true,
  "data": [
    { /* task object */ }
  ],
  "total": 12
}
```

#### PUT `/api/v1/projects/:id/activate`

Set a project as the active project.

**Path Parameters:**

- `id` (string, required): Project UUID

**Example Response:**

```json
{
  "success": true,
  "data": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "isim": "E-ticaret Sitesi",
    "tanim": "Online satış platformu"
  },
  "message": "Project activated successfully"
}
```

---

### Templates

#### GET `/api/v1/templates`

List all available task templates.

**Query Parameters:**

- `kategori` (string, optional): Filter by category

**Example Request:**

```bash
curl "http://localhost:5082/api/v1/templates?kategori=Teknik"
```

**Example Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "39f28dbd-10f3-454c-8b35-52ae6b7ea391",
      "isim": "Bug Raporu",
      "aciklama": "Yazılım hatası bildirimi için detaylı template",
      "kategori": "Teknik",
      "baslik_sablonu": "🐛 [{{modul}}] {{baslik}}",
      "alanlar": [
        {
          "key": "baslik",
          "label": "Başlık",
          "type": "text",
          "required": true
        },
        {
          "key": "ortam",
          "label": "Ortam",
          "type": "select",
          "required": true,
          "options": ["development", "staging", "production"]
        }
      ]
    }
  ],
  "total": 1
}
```

---

### System

#### GET `/api/v1/summary`

Get system-wide summary statistics.

**Example Response:**

```json
{
  "success": true,
  "data": {
    "message": "Summary endpoint - to be implemented"
  }
}
```

---

### Subtask Management

#### GET `/api/v1/tasks/:id/subtasks`

Get all subtasks for a parent task.

**Path Parameters:**

- `id` (string, required): Parent task UUID

**Example Request:**

```bash
curl "http://localhost:5082/api/v1/tasks/550e8400-e29b-41d4-a716-446655440000/subtasks"
```

**Example Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "subtask-uuid-1",
      "baslik": "Setup JWT library",
      "aciklama": "Install and configure JWT dependencies",
      "durum": "tamamlandi",
      "oncelik": "orta",
      "parent_id": "550e8400-e29b-41d4-a716-446655440000",
      "olusturma_zamani": "2025-09-15T15:00:00Z"
    },
    {
      "id": "subtask-uuid-2",
      "baslik": "Implement token generation",
      "aciklama": "Create JWT token generation logic",
      "durum": "devam_ediyor",
      "oncelik": "yuksek",
      "parent_id": "550e8400-e29b-41d4-a716-446655440000",
      "olusturma_zamani": "2025-09-15T15:30:00Z"
    }
  ],
  "total": 2
}
```

#### POST `/api/v1/tasks/:id/subtasks`

Create a subtask under a parent task.

**Path Parameters:**

- `id` (string, required): Parent task UUID

**Request Body:**

```json
{
  "baslik": "Implement API endpoint",
  "aciklama": "Create REST endpoint for subtask creation",
  "oncelik": "yuksek",
  "son_tarih": "2025-10-15",
  "etiketler": "api,backend"
}
```

**Example Response:**

```json
{
  "success": true,
  "data": {
    "id": "subtask-uuid",
    "baslik": "Implement API endpoint",
    "durum": "beklemede",
    "parent_id": "parent-task-uuid"
  },
  "message": "Subtask created successfully"
}
```

#### PUT `/api/v1/tasks/:id/parent`

Change the parent of a task.

**Path Parameters:**

- `id` (string, required): Task UUID to move

**Request Body:**

```json
{
  "new_parent_id": "new-parent-uuid"
}
```

**Note:** Leave `new_parent_id` empty string to move task to root level.

**Example Response:**

```json
{
  "success": true,
  "data": {
    "id": "task-uuid",
    "parent_id": "new-parent-uuid"
  },
  "message": "Parent changed successfully"
}
```

#### GET `/api/v1/tasks/:id/hierarchy`

Get the full hierarchy of a task (parent chain and all subtasks).

**Path Parameters:**

- `id` (string, required): Task UUID

**Example Response:**

```json
{
  "success": true,
  "data": {
    "gorev": {
      "id": "task-uuid",
      "baslik": "Main Task"
    },
    "alt_gorevler": [
      {
        "id": "subtask-1",
        "baslik": "Subtask 1",
        "alt_gorevler": []
      }
    ],
    "toplam_alt_gorev": 5,
    "tamamlanan_alt_gorev": 2
  }
}
```

---

### Dependency Management

#### POST `/api/v1/tasks/:id/dependencies`

Add a dependency between tasks (task `:id` depends on `kaynak_id`).

**Path Parameters:**

- `id` (string, required): Target task UUID (dependent task)

**Request Body:**

```json
{
  "kaynak_id": "prerequisite-task-uuid",
  "baglanti_tipi": "onceki"
}
```

**Example Response:**

```json
{
  "success": true,
  "message": "Dependency added successfully"
}
```

#### DELETE `/api/v1/tasks/:id/dependencies/:dep_id`

Remove a dependency between tasks.

**Path Parameters:**

- `id` (string, required): Target task UUID
- `dep_id` (string, required): Source task UUID (dependency to remove)

**Status:** Not yet implemented (returns 501)

---

### Active Project Management

#### GET `/api/v1/active-project`

Get the currently active project.

**Example Response (with active project):**

```json
{
  "success": true,
  "data": {
    "id": "project-uuid",
    "isim": "E-commerce Site",
    "tanim": "Online shopping platform"
  }
}
```

**Example Response (no active project):**

```json
{
  "success": true,
  "data": null,
  "message": "No active project set"
}
```

#### DELETE `/api/v1/active-project`

Remove the active project setting.

**Example Response:**

```json
{
  "success": true,
  "message": "Active project removed successfully"
}
```

---

### Language Management

#### GET `/api/v1/language`

Get the current MCP server language setting.

**Example Response:**

```json
{
  "success": true,
  "language": "tr"
}
```

#### POST `/api/v1/language`

Change the MCP server language.

**Request Body:**

```json
{
  "language": "en"
}
```

**Supported Languages:**

- `tr`: Turkish
- `en`: English

**Example Response:**

```json
{
  "success": true,
  "language": "en",
  "message": "Language changed to en"
}
```

**Notes:**

- Language change affects all subsequent MCP tool responses
- Web UI automatically syncs language with this endpoint
- Language preference is stored server-side (not persisted across restarts)

---

## 🚨 Error Codes

| HTTP Status | Description | Example |
|-------------|-------------|---------|
| 200 | Success | Request completed successfully |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid request body or parameters |
| 404 | Not Found | Resource not found (task, project, etc.) |
| 500 | Internal Server Error | Server-side error occurred |

### Common Error Responses

**400 Bad Request:**

```json
{
  "success": false,
  "error": "Task ID is required"
}
```

**404 Not Found:**

```json
{
  "success": false,
  "error": "failed to get task with ID abc123: record not found"
}
```

**500 Internal Server Error:**

```json
{
  "success": false,
  "error": "failed to list tasks: database connection error"
}
```

## 📊 Rate Limiting

Currently, there is no rate limiting implemented. The API is designed for local development use.

**Future Considerations:**

- Implement token bucket algorithm for production use
- Default: 100 requests per minute per IP
- Burst: 20 requests

## 🔄 Pagination

List endpoints support pagination via `limit` and `offset` parameters.

**Example:**

```bash
# Get first 10 tasks
GET /api/v1/tasks?limit=10&offset=0

# Get next 10 tasks
GET /api/v1/tasks?limit=10&offset=10
```

**Best Practices:**

- Default `limit`: 50
- Maximum `limit`: 100
- Use `total` field in response to calculate total pages

## 🧪 Testing the API

### Using cURL

```bash
# List all tasks
curl http://localhost:5082/api/v1/tasks

# Get specific task
curl http://localhost:5082/api/v1/tasks/task-uuid

# Create task from template
curl -X POST http://localhost:5082/api/v1/tasks/from-template \
  -H "Content-Type: application/json" \
  -d '{
    "template_id": "template-uuid",
    "degerler": {
      "baslik": "Test Task",
      "aciklama": "Test description"
    }
  }'

# Update task status
curl -X PUT http://localhost:5082/api/v1/tasks/task-uuid \
  -H "Content-Type: application/json" \
  -d '{"durum": "tamamlandi"}'

# Delete task
curl -X DELETE http://localhost:5082/api/v1/tasks/task-uuid

# Change language
curl -X POST http://localhost:5082/api/v1/language \
  -H "Content-Type: application/json" \
  -d '{"language": "en"}'
```

### Using HTTPie

```bash
# List tasks
http GET localhost:5082/api/v1/tasks

# Create task
http POST localhost:5082/api/v1/tasks/from-template \
  template_id="template-uuid" \
  degerler:='{"baslik": "Test"}'

# Update task
http PUT localhost:5082/api/v1/tasks/task-uuid \
  durum="tamamlandi"
```

## 📚 SDK Support

### JavaScript/TypeScript

The Web UI includes a React Query client. See `gorev-web/src/api/client.ts` for implementation.

### Go (Future)

Native Go client SDK planned for v0.17.0.

## 🔮 Future Enhancements

- [ ] WebSocket support for real-time updates
- [ ] Bulk operations endpoint (`POST /api/v1/tasks/bulk`)
- [ ] Advanced filtering with query language
- [ ] File attachment support
- [ ] Audit log endpoint
- [ ] OpenAPI/Swagger specification
- [ ] GraphQL endpoint (alternative to REST)

## 🤝 Contributing

When adding new endpoints:

1. Follow RESTful conventions
2. Use consistent response format
3. Add comprehensive error handling
4. Document all parameters and responses
5. Update this reference document
6. Add integration tests

---

**Questions?** Check the [main README](../../README.md) or open an issue on GitHub.
