# Image Processing Service

Image Processing Service is a minimal image upload microservice designed for flexible embedding into third-party systems. It supports:

- Multipart image upload via REST
- Binary image upload via REST
- gRPC interface for image upload and retrieval
- WebP compression and thumbnail generation
- Storage in MinIO (S3-compatible)

---

## 🚀 Features

- Image compression to WebP
- Multiple thumbnail sizes
- Asynchronous processing
- Presigned URL generation
- Minimal external dependencies
- Clean Architecture (DDD + Ports & Adapters)
- Transport via REST or gRPC
- Optional API Key authentication for REST or gRPC (configurable in `.env`)
- Async processing

---

## 🛠 Tech Stack

- **Go (Golang)** — backend API
- **Gin** — HTTP web framework
- **gRPC** — for image upload and retrieval
- **GORM** — ORM for PostgreSQL
- **Redis** — for fast cache and session TTL
- **PostgreSQL** — for persistent game sessions
- **Docker & Docker Compose** — containerized development
- **Makefile** — for devtools and automation

---

## 🔗 API Endpoints

REST Base URL
👉 http://localhost:8000

GRPC Base URL
👉 localhost:50051

Ports are configurable in `.env`

---

## 📘 REST API

### Upload image (multipart/form-data)

**POST** `/image/upload`

**Headers:**

```
Content-Type: multipart/form-data
X-API-Key: <optional-token-from-env>
```

**Form Data:**

```
file=<image>
```

**Response:**

```json
{
  "id": "1a2b3c",
  "original_key": "uploads/originals/1a2b3c",
  "status": "pending"
}
```

### Upload image (binary)

**POST** `/image/upload-binary`

**Headers:**

```
Content-Type: image/jpeg
X-API-Key: <optional-token-from-env>
```

**Body:** binary content of the image

**Response:**

```json
{
  "id": "1a2b3c",
  "original_key": "uploads/originals/1a2b3c",
  "status": "pending"
}
```

### Get image info

**GET** `/image/{id}`

**Response:**

```json
{
  "id": "1a2b3c",
  "original_key": "uploads/originals/1a2b3c",
  "compressed_key": "uploads/compressed/1a2b3c.webp",
  "status": "ready",
  "error_message": null,
  "thumbnails": [
    {
      "size": "150x150",
      "key": "uploads/thumbnails/1a2b3c_150x150.webp",
      "type": "small"
    },
    ...
  ]
}
```

---

## 🔧 gRPC API

**Proto File:** `proto/image.proto`

```protobuf
service ImageService {
  rpc UploadImage(UploadImageRequest) returns (ImageResponse);
  rpc GetImage(GetImageRequest) returns (ImageResponse);
}

message UploadImageRequest {
  bytes data = 1;
}

message GetImageRequest {
  string id = 1;
}

message ThumbnailShort {
  string size = 1;
  string key = 2;
  string type = 3;
}

message ImageResponse {
  string id = 1;
  string original_key = 2;
  optional string compressed_key = 3;
  string status = 4;
  optional string error_message = 5;
  repeated ThumbnailShort thumbnails = 6;
}
```

---

## ⚡️ Getting Started

```shell
# Clone the repo
git clone https://github.com/m1n64/image-resizing-service
````
```shell
cd image-resizing-service
```

```shell
# Copy environment variables
cp .env.example .env
````

```shell
# Start environment
make up
```

---

## 📋 Makefile Commands

| Command                                 | Description                                  |
|-----------------------------------------|----------------------------------------------|
| `make up`                               | Start the DEV environment                    |
| `make prod`                             | Start the PROD environment                   |
| `make stop`                             | Stop all containers                          |
| `make down`                             | Remove all containers and volumes            |
| `make restart`                          | Restart all containers                       |
| `make restart-container CONTAINER=name` | Restart a specific container                 |
| `make stop-container CONTAINER=name`    | Stop a specific container                    |
| `make bash`                             | Open a bash shell inside the app container   |
| `make logs name`                        | View logs of a specific container            |
| `make app-logs`                         | View last logs from the app container        |
| `make psql`                             | Open psql shell with credentials from `.env` |
| `make redis`                            | Open redis-cli inside the Redis container    |
| `make seed`                             | Run seeders to populate countries & cities   |

---

## 🔧 gRPC build

```shell
protoc --go_out=. --go-grpc_out=. <.proto>
```

---

## 🗒️ MinIO S3 nginx config template

```nginx
server {
    listen 443 ssl;
    server_name localhost;

    #ssl_certificate /etc/nginx/ssl/cdn/server.crt;
    #ssl_certificate_key /etc/nginx/ssl/cdn/server.key;

    ssl_protocols TLSv1.2;
    ssl_prefer_server_ciphers on;

    access_log /var/log/nginx/cdn_access.log;
    error_log /var/log/nginx/cdn_error.log;

    location / {
        proxy_pass http://minio:9000;
        proxy_set_header Host minio:9000; #IMPORTANT FOR MinIO S3
        proxy_set_header X-Host-Override $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 🧙‍♂️ Author

Made with ❤️ by the **[Kirill Sakharov](https://github.com/m1n64) ([LinkedIn](https://www.linkedin.com/in/kirill-sakharov-862072227/))**
