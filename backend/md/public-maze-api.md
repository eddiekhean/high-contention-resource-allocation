# Public User – Maze API Flow

## Tổng quan vai trò

- **Frontend**
  - Chọn maze có sẵn
  - Upload maze tạm thời
  - Chạy thuật toán giải maze

- **Backend API**
  - Gatekeeper cho public user
  - Rate limit
  - Quản lý session
  - Validate input

- **S3**
  - Chỉ chứa maze seed
  - Read-only
  - Không public trực tiếp

- **Session Store**
  - Lưu maze upload tạm thời của public user
  - Có TTL

---

## 1. Load danh sách maze có sẵn (S3 – Read Only)

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant BE as Backend API
    participant S3 as AWS S3

    FE->>BE: GET /public/mazes
    BE->>BE: Rate limit + anon_id
    BE->>S3: List maze metadata
    BE->>S3: Generate presigned GET URLs
    BE-->>FE: maze list + signedUrls

    FE->>S3: GET maze image
    FE->>FE: Render maze
```

## 2. Public user upload maze (Session Only)

### Constraints

- Maze upload của public user **không được lưu vĩnh viễn**
- Không ghi dữ liệu vào database
- Không upload lên S3
- Dữ liệu chỉ tồn tại trong session với TTL giới hạn

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant BE as Backend API
    participant SESSION as Session Store

    FE->>FE: Select maze image
    FE->>BE: POST /public/maze/session (multipart)

    BE->>BE: Validate size / type
    BE->>BE: Optional resize / normalize
    BE->>SESSION: Store image (ttl=15m)
    BE-->>FE: sessionMazeId
```

## 3. Lấy lại maze trong session (Preview / Rerun)

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant BE as Backend API
    participant SESSION as Session Store

    FE->>BE: GET /public/maze/session/:id
    BE->>SESSION: Load maze image
    BE-->>FE: image stream
```

## 4. Match / Nhận diện maze (dHash – Seed Only)

Public user chỉ được match với maze seed, không match với dữ liệu upload của user khác.

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant BE as Backend API
    participant DB as Database

    FE->>FE: Compute dHash
    FE->>BE: POST /public/maze/match (dhash)

    BE->>BE: Rate limit
    BE->>DB: Query seed maze hashes
    BE->>BE: Hamming distance

    alt Found
        BE-->>FE: matched=true + mazeId
    else Not found
        BE-->>FE: matched=false
    end
```

## 5. Chạy thuật toán giải maze (Stateless)

Backend không lưu trạng thái giải bài toán của public user.

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant BE as Backend API
    participant SESSION as Session Store

    FE->>BE: POST /public/maze/solve
    Note right of FE: mazeSource = seedId | sessionMazeId\nalgorithm = BFS / DFS / A*

    BE->>BE: Validate input
    BE->>SESSION: Load maze if sessionMazeId
    BE->>BE: Solve maze
    BE-->>FE: path + steps
```

## 6. High-level Flowchart (Public User)

```mermaid
flowchart TD
    A[Public User] --> B[GET /public/mazes]
    B --> C[Render seed maze]

    A --> D[Upload maze]
    D --> E[Store in session]

    C --> F[Run algorithm]
    E --> F

    F --> G[Render path]
    E --> H[Session expired]
```

## 7. API Scope Tóm Tắt

| Endpoint | Ghi DB | S3 | Session |
| :--- | :--- | :--- | :--- |
| `GET /public/mazes` | No | GET | No |
| `POST /public/maze/session` | No | No | Yes |
| `GET /public/maze/session/:id` | No | No | Yes |
| `POST /public/maze/match` | No | No | No |
| `POST /public/maze/solve` | No | No | Optional |

## Nguyên tắc thiết kế

1. **Quyền hạn**: Public user không có quyền ghi dữ liệu lâu dài.
2. **Bảo mật S3**: Không cấp presigned PUT cho public user.
3. **Lưu trữ tạm thời**: Mọi dữ liệu upload chỉ tồn tại trong session.
4. **Kiểm soát lưu lượng**: Rate limit bắt buộc cho tất cả public endpoints.
