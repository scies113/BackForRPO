# ==========================================
# ML Prediction Service — Dockerfile
# ==========================================
FROM python:3.11-slim

# Отключаем создание .pyc файлов (уменьшает размер) и буферизацию
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app

# Копируем зависимости
COPY requirements.txt .

# Используем кэш BuildKit и самое полное зеркало (Tsinghua) для ускорения загрузки
# Очищаем временные файлы для уменьшения итогового размера образа
RUN --mount=type=cache,target=/root/.cache/pip \
    pip install --upgrade pip && \
    pip install -i https://pypi.tuna.tsinghua.edu.cn/simple -r requirements.txt && \
    find /usr/local -name "*.pyc" -delete && \
    find /usr/local -type d -name "__pycache__" -delete && \
    rm -rf /tmp/*

# Копируем всё остальное
COPY . .

EXPOSE 8000

CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]
