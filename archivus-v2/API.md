# Archivus v2 — API Reference

## Authentication

All endpoints except `/login` and `/files/download/{filepath}` require authentication via one of:

| Method | Header | Value |
|--------|--------|-------|
| JWT | `Authorization` | `Bearer <token>` |
| API Key | `X-API-Key` | `<api_key>` |

The auth middleware injects `userId` and `username` into downstream request headers automatically — you do not set these yourself.

---

## Response Format

All JSON responses follow a consistent envelope. Errors return HTTP 4xx/5xx with:
```json
{ "error": "message" }
```
Success responses return HTTP 200 with either:
```json
{ "message": "..." }
```
or a data payload directly.

---

## Endpoints

### Health

#### `GET /health`
No auth required.

**Response**
```json
{ "message": "OK" }
```

---

### Auth

#### `POST /login`
No auth required. Either `password` or `pin` must be provided (not both required).

**Request Body**
```json
{
  "username": "string",
  "password": "string",
  "pin": "string"
}
```

**Response**
```json
{
  "token": "string",
  "user_id": "uuid"
}
```

---

### Files

#### `GET /files/get/`
List files in a folder.

**Query Params**
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| `folder` | string | No | Folder path to list. Omit for root. |

**Response**
```json
{
  "files": [...],
  "size": 12345
}
```

---

#### `GET /files/get-signed-url/{filepath}`
Generate a time-limited signed URL for a file. Use this to get a download link to pass to clients.

**Path Params**
| Param | Description |
|-------|-------------|
| `filepath` | Full file path (e.g. `photos/image.jpg`) |

**Response**
```json
{
  "signed_url": "string"
}
```

---

#### `GET /files/download/{filepath}`
Download a file. **No auth required** — validated via signed URL query params.

**Path Params**
| Param | Description |
|-------|-------------|
| `filepath` | Full file path |

**Query Params**
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| `signature` | string | Yes | HMAC signature from signed URL |
| `expires_at` | int | Yes | Unix timestamp; request rejected after this time |
| `compressed` | string | No | `"true"` to serve compressed version |

**Response**: Binary file stream with `Content-Disposition: attachment`.

---

#### `POST /files/upload/`
Upload one or more files (max 500 MB total). Sent as `multipart/form-data`.

**Form Fields**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `folder` | string | No | Destination folder path |
| `file` | file | Yes | One or more files (repeat field for multiple) |

**Response**
```json
{ "message": "File uploaded successfully" }
```

---

#### `POST /files/move/`
Move a file to a new location.

**Request Body**
```json
{
  "filePath": "string",
  "dst": "string"
}
```

**Response**
```json
{ "message": "File moved successfully" }
```

---

#### `POST /files/delete/`
Delete a file.

**Request Body**
```json
{
  "filePath": "string"
}
```

**Response**
```json
{ "message": "File deleted successfully" }
```

---

### Folders

#### `POST /folder/add/`
Create a new folder. Requires the authenticated user to have `write_access`. If `user_dir_lock` is enabled, the user can only create folders under their own username directory.

**Request Body**
```json
{
  "Folder": "string"
}
```

**Response**
```json
{ "message": "Folder created successfully" }
```

**Errors**
- `403 Forbidden` — user does not have write access
- `400 Bad Request` — `Folder` field is empty

---

### Big Upload (Chunked)

For files too large for a single request. The flow is: **Initiate → Upload Chunks → Finalize**.

Chunks are verified with SHA-512 hashes. The server de-duplicates sessions by `wark` (a hash of file size + chunk hashes), so interrupted uploads can be resumed.

---

#### `POST /bigupload/initiate/`
Start or resume a chunked upload session.

**Request Body**
```json
{
  "file_name": "string",
  "file_size": 123456789,
  "chunk_hashes": ["sha512hex", "sha512hex", "..."],
  "mod_time": 1712345678
}
```

| Field | Type | Description |
|-------|------|-------------|
| `file_name` | string | Original filename |
| `file_size` | int64 | Total file size in bytes |
| `chunk_hashes` | string[] | SHA-512 hex hash of each chunk, in order |
| `mod_time` | int64 | File modification time (Unix timestamp) |

**Response**
```json
{
  "wark": "string",
  "needed_chunks": [0, 1, 3]
}
```

`needed_chunks` lists the 0-based indices of chunks that still need to be uploaded (already-uploaded chunks are excluded for resumable uploads).

---

#### `POST /bigupload/chunk/`
Upload a single chunk. Body is raw binary data.

**Headers**
| Header | Type | Description |
|--------|------|-------------|
| `X-Wark` | string | Session wark from initiate response |
| `X-Idx` | int | 0-based chunk index |
| `X-Chunk-Hash` | string | SHA-512 hex hash of this chunk |

**Body**: Raw binary chunk bytes.

**Response**: `200 OK` plain text.

**Errors**
- `400` — missing/invalid headers, hash mismatch, invalid index

---

#### `POST /bigupload/finalize/`
Finalize the upload after all chunks are sent. Moves the assembled file to its final destination.

**Request Body**
```json
{
  "wark": "string"
}
```

**Response**
```json
{
  "status": "ok",
  "path": "/absolute/path/to/file"
}
```

**Errors**
- `400` — not all chunks uploaded yet

---

### Tempora (Todos & Projects)

#### `GET /tempora/projects`
List all projects for the authenticated user.

**Response**
```json
[
  {
    "id": 1,
    "title": "string",
    "description": "string",
    "projectId": 1,
    "userId": "uuid"
  }
]
```

---

#### `POST /tempora/projects`
Create a project.

**Request Body**
```json
{
  "Title": "string",
  "Description": "string"
}
```

**Response**
```json
{ "message": "Project created successfully" }
```

---

#### `DELETE /tempora/projects`
Delete a project by ID.

**Request Body**: Raw numeric project ID (not wrapped in object).
```json
1
```

**Response**
```json
{ "message": "Project deleted successfully" }
```

---

#### `GET /tempora/todos`
List todos for the authenticated user.

**Query Params**
| Param | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | uint | No | Filter by project ID |

**Response**
```json
[
  {
    "id": 1,
    "title": "string",
    "description": "string",
    "status": 0,
    "priority": 1,
    "projectId": 1
  }
]
```

**Status values**: `0` = todo, `1` = in_progress, `2` = done

**Priority values**: `0` = low, `1` = medium, `2` = high

---

#### `POST /tempora/todos`
Create one or more todos.

**Request Body**: Array of todo objects.
```json
[
  {
    "title": "string",
    "description": "string",
    "status": 0,
    "priority": 1,
    "projectId": 1
  }
]
```

**Response**
```json
{ "message": "Todos created successfully" }
```

---

#### `POST /tempora/todos/update`
Update the status of one or more todos.

**Request Body**: Array of update objects.
```json
[
  { "Id": 1, "Status": 2 },
  { "Id": 2, "Status": 1 }
]
```

**Response**
```json
{ "message": "Todos marked as done successfully" }
```

---

#### `DELETE /tempora/todos/update`
Delete todos by ID.

**Request Body**: Array of uint IDs.
```json
[1, 2, 3]
```

**Response**
```json
{ "message": "Todos deleted successfully" }
```
