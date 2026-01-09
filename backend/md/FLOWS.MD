# Image Similarity & Upload Flow (dHash + S3)

## Luồng 1: Match ảnh (KHÔNG upload)

Quy trình kiểm tra xem ảnh đã tồn tại trong hệ thống chưa dựa trên mã băm dHash.

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant BE as Backend API
    participant DB as Database
    participant S3 as AWS S3

    FE->>FE: Compute dHash (resize 9x8, grayscale)
    FE->>BE: POST /images/match (dhash)

    BE->>DB: Query candidate hashes
    BE->>BE: Compute Hamming distance

    alt Found (distance <= threshold)
        BE->>S3: Generate presigned GET URL
        BE-->>FE: matched=true + metadata + signedUrl
        FE->>S3: GET image
        FE->>FE: Render image
    else Not found
        BE-->>FE: matched=false
    end
```

## Luồng 2: Upload ảnh mới

Quy trình tải ảnh lên S3 và lưu thông tin vào Database nếu không tìm thấy ảnh trùng lặp.

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant BE as Backend API
    participant S3 as AWS S3
    participant DB as Database

    FE->>BE: POST /images/upload-url
    BE->>S3: Generate presigned PUT URL
    BE-->>FE: uploadUrl + s3Key

    FE->>S3: PUT image (direct upload)
    FE->>BE: POST /images (dhash, s3Key, metadata)

    BE->>S3: HeadObject (verify object exists)
    BE->>DB: Insert image record
    BE-->>FE: Upload success
```

## Tổng quan quy trình (Flowchart)

```mermaid
flowchart TD
    A[User selects image] --> B[Compute dHash]
    B --> C[POST /images/match]

    C -->|matched| D[Request signed GET URL]
    D --> E[Render image]

    C -->|not matched| F[POST /images/upload-url]
    F --> G[PUT image to S3]
    G --> H[POST /images (save DB)]
    H --> I[Render uploaded image]
```
