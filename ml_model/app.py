from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import pandas as pd
import numpy as np
import pickle
import json
from model import FootballPredictionModel
import uvicorn
import sys

app = FastAPI(title="Football Prediction API")

# Global variables to hold model and metadata
model_wrapper = None
team_encoder = None
latest_stats = None
feature_cols = None

class MatchRequest(BaseModel):
    home_team: str
    away_team: str

class MatchResponse(BaseModel):
    home_win_prob: float
    draw_prob: float
    away_win_prob: float

@app.on_event("startup")
def load_assets():
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
            
        print("All assets loaded successfully.")
    except Exception as e:
        print(f"Error loading assets: {e}")
        # Not exiting so tests can run without models, but in prod you might exit
        # sys.exit(1)

@app.post("/predict", response_model=MatchResponse)
def predict_match(request: MatchRequest):
    if model_wrapper is None:
        raise HTTPException(status_code=500, detail="Model is not loaded. Train the model first.")
        
    home = request.home_team
    away = request.away_team
    
    # Check if teams exist
    known_teams = team_encoder.classes_
    if home not in known_teams:
        raise HTTPException(status_code=400, detail=f"Unknown home team: {home}")
    if away not in known_teams:
        raise HTTPException(status_code=400, detail=f"Unknown away team: {away}")
        
    # Encode teams
    home_encoded = team_encoder.transform([home])[0]
    away_encoded = team_encoder.transform([away])[0]
    
    # Get past stats
    home_past = latest_stats.get(home, [])
    away_past = latest_stats.get(away, [])
    
    stats_cols = ['GoalsScored', 'GoalsConceded', 'Shots', 'ShotsOnTarget', 'Corners', 'Fouls', 'YellowCards', 'RedCards']
    
    home_means = {f"Home_avg_{c}": (sum(x[c] for x in home_past)/len(home_past) if home_past else 0.0) for c in stats_cols}
    away_means = {f"Away_avg_{c}": (sum(x[c] for x in away_past)/len(away_past) if away_past else 0.0) for c in stats_cols}
    
    # Construct feature dictionary
    features = {
        'HomeTeam_encoded': home_encoded,
        'AwayTeam_encoded': away_encoded
    }
    features.update(home_means)
    features.update(away_means)
    
    # Create DataFrame with the exact columns order as during training
    df_features = pd.DataFrame([features])[feature_cols]
    
    # Predict (model returns probabilities for [Away(0), Draw(1), Home(2)])
    probs = model_wrapper.predict_proba(df_features)[0]
    
    return MatchResponse(
        away_win_prob=float(probs[0]),
        draw_prob=float(probs[1]),
        home_win_prob=float(probs[2])
    )

if __name__ == "__main__":
    uvicorn.run("app:app", host="0.0.0.0", port=8000, reload=True)
