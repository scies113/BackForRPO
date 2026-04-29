import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder
import pickle
import os

def load_and_preprocess_data(csv_path='epl_final.csv', window_size=5):
    # Load data
    df = pd.read_csv(csv_path, sep=';')
    
    # Fill NaN values for numerical stats
    numeric_cols = [
        'FullTimeHomeGoals', 'FullTimeAwayGoals', 'HomeShots', 'AwayShots',
        'HomeShotsOnTarget', 'AwayShotsOnTarget', 'HomeCorners', 'AwayCorners',
        'HomeFouls', 'AwayFouls', 'HomeYellowCards', 'AwayYellowCards',
        'HomeRedCards', 'AwayRedCards'
    ]
    for col in numeric_cols:
        df[col] = pd.to_numeric(df[col], errors='coerce').fillna(0)
        
    # Sort by date just to be sure (MatchDate format is DD.MM.YYYY)
    df['MatchDate'] = pd.to_datetime(df['MatchDate'], format='%d.%m.%Y', errors='coerce')
    df = df.dropna(subset=['MatchDate']).sort_values('MatchDate').reset_index(drop=True)
    
    # Calculate rolling averages
    stats_cols = ['GoalsScored', 'GoalsConceded', 'Shots', 'ShotsOnTarget', 'Corners', 'Fouls', 'YellowCards', 'RedCards']
    team_stats = {team: [] for team in set(df['HomeTeam']).union(set(df['AwayTeam']))}
    
    new_features = []
    
    for idx, row in df.iterrows():
        home = row['HomeTeam']
        away = row['AwayTeam']
        
        # Get current rolling stats for home and away
        home_past = team_stats[home][-window_size:] if len(team_stats[home]) > 0 else []
        away_past = team_stats[away][-window_size:] if len(team_stats[away]) > 0 else []
        
        home_means = {f"Home_avg_{c}": (sum(x[c] for x in home_past)/len(home_past) if home_past else 0.0) for c in stats_cols}
        away_means = {f"Away_avg_{c}": (sum(x[c] for x in away_past)/len(away_past) if away_past else 0.0) for c in stats_cols}
        
        # Append current match stats to team histories
        team_stats[home].append({
            'GoalsScored': row['FullTimeHomeGoals'], 'GoalsConceded': row['FullTimeAwayGoals'],
            'Shots': row['HomeShots'], 'ShotsOnTarget': row['HomeShotsOnTarget'],
            'Corners': row['HomeCorners'], 'Fouls': row['HomeFouls'],
            'YellowCards': row['HomeYellowCards'], 'RedCards': row['HomeRedCards']
        })
        team_stats[away].append({
            'GoalsScored': row['FullTimeAwayGoals'], 'GoalsConceded': row['FullTimeHomeGoals'],
            'Shots': row['AwayShots'], 'ShotsOnTarget': row['AwayShotsOnTarget'],
            'Corners': row['AwayCorners'], 'Fouls': row['AwayFouls'],
            'YellowCards': row['AwayYellowCards'], 'RedCards': row['AwayRedCards']
        })
        
        features = {**home_means, **away_means}
        new_features.append(features)
        
    features_df = pd.DataFrame(new_features, index=df.index)
    df = pd.concat([df, features_df], axis=1)
    
    # Encode Teams
    le = LabelEncoder()
    # Fit on all unique teams
    all_teams = np.unique(df[['HomeTeam', 'AwayTeam']].values)
    le.fit(all_teams)
    
    df['HomeTeam_encoded'] = le.transform(df['HomeTeam'])
    df['AwayTeam_encoded'] = le.transform(df['AwayTeam'])
    
    # Save encoder for inference
    with open('team_encoder.pkl', 'wb') as f:
        pickle.dump(le, f)
        
    # Map Target Variable (A=0, D=1, H=2)
    # This ordering makes sense: Away win, Draw, Home win
    target_mapping = {'A': 0, 'D': 1, 'H': 2}
    df['Target'] = df['FullTimeResult'].map(target_mapping)
    df = df.dropna(subset=['Target']) # Drop any rows where result is missing
    
    # Drop rows that don't have enough history (e.g. first few matches of the dataset where rolling avg is 0)
    # Actually, we can keep them, tree models handle zeros fine, but it might be noise.
    
    feature_cols = ['HomeTeam_encoded', 'AwayTeam_encoded'] + list(features_df.columns)
    
    X = df[feature_cols]
    y = df['Target']
    
    import json
    # Save feature names for reference
    with open('feature_cols.json', 'w') as f:
        json.dump(feature_cols, f)
        
    # Save the latest stats for each team to be used during inference
    latest_team_stats = {team: stats[-window_size:] for team, stats in team_stats.items()}
    with open('latest_team_stats.pkl', 'wb') as f:
        pickle.dump(latest_team_stats, f)
    
    # Temporal train/test split (last 20% for testing)
    split_index = int(len(X) * 0.8)
    X_train, X_test = X.iloc[:split_index], X.iloc[split_index:]
    y_train, y_test = y.iloc[:split_index], y.iloc[split_index:]
    
    return X_train, X_test, y_train, y_test, feature_cols

if __name__ == "__main__":
    X_train, X_test, y_train, y_test, cols = load_and_preprocess_data('epl_final.csv')
    print(f"Data loaded. Train shape: {X_train.shape}, Test shape: {X_test.shape}")
    print("Features:", cols)
