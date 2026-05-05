# ==========================================
# ML Prediction Service — Dockerfile
# ==========================================
FROM python:3.11-slim

WORKDIR /app

# Копируем зависимости и устанавливаем
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Копируем исходный код и артефакты модели
COPY app.py model.py dataset.py train.py ./
COPY xgboost_model.pkl team_encoder.pkl latest_team_stats.pkl feature_cols.json ./
COPY epl_final.csv ./

EXPOSE 8000

CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]
