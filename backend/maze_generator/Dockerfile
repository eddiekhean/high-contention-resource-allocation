FROM python:3.11-slim

# Không tạo .pyc, log ra stdout ngay
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

WORKDIR /app

# Cài dependency trước để tận dụng cache
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY . .

EXPOSE 8000

# Chạy đúng 1 process uvicorn
CMD ["uvicorn", "maze.main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "1"]

