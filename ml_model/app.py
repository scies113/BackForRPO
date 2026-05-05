"""
FastAPI-сервер для предсказания результатов футбольных матчей.

Загружает обученную XGBoost-модель и предоставляет HTTP-эндпоинт
/predict для получения вероятностей исходов матча (победа хозяев,
ничья, победа гостей).
"""

import json
import pickle
from contextlib import asynccontextmanager

import pandas as pd
import uvicorn
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from model import FootballPredictionModel

# Глобальные переменные для хранения модели и метаданных
model_wrapper = None
team_encoder = None
latest_stats = None
feature_cols = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Загружает модель и артефакты при старте сервера."""
    global model_wrapper, team_encoder, latest_stats, feature_cols
    try:
        model_wrapper = FootballPredictionModel()
        model_wrapper.load('xgboost_model.pkl')

        with open('team_encoder.pkl', 'rb') as f:
            team_encoder = pickle.load(f)

        with open('latest_team_stats.pkl', 'rb') as f:
            latest_stats = pickle.load(f)

        with open('feature_cols.json', 'r') as f:
            feature_cols = json.load(f)

        print("Все артефакты успешно загружены.")
    except Exception as e:
        print(f"Ошибка при загрузке артефактов: {e}")

    yield


app = FastAPI(title="Football Prediction API", lifespan=lifespan)


@app.get("/health")
def health():
    """Проверка работоспособности сервиса."""
    return {"status": "ok", "model_loaded": model_wrapper is not None}


@app.get("/teams")
def list_teams():
    """Возвращает список всех команд, известных модели."""
    if team_encoder is None:
        return {"teams": []}
    return {"teams": sorted(team_encoder.classes_.tolist())}


class MatchRequest(BaseModel):
    """Входные данные для предсказания: названия домашней и гостевой команд."""
    home_team: str
    away_team: str


class MatchResponse(BaseModel):
    """Результат предсказания: вероятности трёх исходов матча."""
    home_win_prob: float
    draw_prob: float
    away_win_prob: float


@app.post("/predict", response_model=MatchResponse)
def predict_match(request: MatchRequest):
    """Предсказывает вероятности исходов матча по названиям команд."""
    if model_wrapper is None:
        raise HTTPException(
            status_code=500,
            detail="Модель не загружена. Сначала обучите модель командой: python train.py"
        )

    home = request.home_team
    away = request.away_team

    # Проверка существования команд
    known_teams = team_encoder.classes_
    if home not in known_teams:
        raise HTTPException(status_code=400, detail=f"Неизвестная домашняя команда: {home}")
    if away not in known_teams:
        raise HTTPException(status_code=400, detail=f"Неизвестная гостевая команда: {away}")

    # Кодирование команд
    home_encoded = team_encoder.transform([home])[0]
    away_encoded = team_encoder.transform([away])[0]

    # Получение статистики за последние матчи
    home_past = latest_stats.get(home, [])
    away_past = latest_stats.get(away, [])

    stats_cols = [
        'GoalsScored', 'GoalsConceded', 'Shots', 'ShotsOnTarget',
        'Corners', 'Fouls', 'YellowCards', 'RedCards'
    ]

    home_means = {
        f"Home_avg_{c}": (sum(x[c] for x in home_past) / len(home_past) if home_past else 0.0)
        for c in stats_cols
    }
    away_means = {
        f"Away_avg_{c}": (sum(x[c] for x in away_past) / len(away_past) if away_past else 0.0)
        for c in stats_cols
    }

    # Формирование вектора признаков
    features = {
        'HomeTeam_encoded': home_encoded,
        'AwayTeam_encoded': away_encoded
    }
    features.update(home_means)
    features.update(away_means)

    # Создание DataFrame с правильным порядком столбцов
    df_features = pd.DataFrame([features])[feature_cols]

    # Предсказание: модель возвращает вероятности [Away(0), Draw(1), Home(2)]
    probs = model_wrapper.predict_proba(df_features)[0]

    return MatchResponse(
        away_win_prob=float(probs[0]),
        draw_prob=float(probs[1]),
        home_win_prob=float(probs[2])
    )


if __name__ == "__main__":
    uvicorn.run("app:app", host="0.0.0.0", port=8000, reload=True)
