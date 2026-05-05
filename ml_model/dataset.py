"""
Модуль предобработки данных.

Загружает CSV-файл с матчами EPL, вычисляет скользящие средние
показатели команд за последние N матчей и формирует обучающую
и тестовую выборки.
"""

import json
import pickle

import numpy as np
import pandas as pd
from sklearn.preprocessing import LabelEncoder


def load_and_preprocess_data(csv_path='epl_final.csv', window_size=5):
    """Загружает данные из CSV и подготавливает признаки для обучения.

    Для каждого матча вычисляются средние показатели домашней и гостевой
    команд за последние ``window_size`` игр (голы, удары, угловые и т.д.).
    Это предотвращает утечку данных, так как используется только информация,
    доступная до начала матча.

    Args:
        csv_path: Путь к CSV-файлу с матчами (разделитель — точка с запятой).
        window_size: Количество предыдущих матчей для расчёта средних.

    Returns:
        Кортеж (X_train, X_test, y_train, y_test, feature_cols).
    """
    # Загрузка данных
    df = pd.read_csv(csv_path, sep=';')

    # Приведение числовых столбцов к нужному типу
    numeric_cols = [
        'FullTimeHomeGoals', 'FullTimeAwayGoals', 'HomeShots', 'AwayShots',
        'HomeShotsOnTarget', 'AwayShotsOnTarget', 'HomeCorners', 'AwayCorners',
        'HomeFouls', 'AwayFouls', 'HomeYellowCards', 'AwayYellowCards',
        'HomeRedCards', 'AwayRedCards'
    ]
    for col in numeric_cols:
        df[col] = pd.to_numeric(df[col], errors='coerce').fillna(0)

    # Сортировка по дате (формат DD.MM.YYYY)
    df['MatchDate'] = pd.to_datetime(df['MatchDate'], format='%d.%m.%Y', errors='coerce')
    df = df.dropna(subset=['MatchDate']).sort_values('MatchDate').reset_index(drop=True)

    # Расчёт скользящих средних показателей
    stats_cols = [
        'GoalsScored', 'GoalsConceded', 'Shots', 'ShotsOnTarget',
        'Corners', 'Fouls', 'YellowCards', 'RedCards'
    ]
    team_stats = {team: [] for team in set(df['HomeTeam']).union(set(df['AwayTeam']))}

    new_features = []

    for idx, row in df.iterrows():
        home = row['HomeTeam']
        away = row['AwayTeam']

        # Получение статистики за последние window_size матчей
        home_past = team_stats[home][-window_size:] if len(team_stats[home]) > 0 else []
        away_past = team_stats[away][-window_size:] if len(team_stats[away]) > 0 else []

        home_means = {
            f"Home_avg_{c}": (sum(x[c] for x in home_past) / len(home_past) if home_past else 0.0)
            for c in stats_cols
        }
        away_means = {
            f"Away_avg_{c}": (sum(x[c] for x in away_past) / len(away_past) if away_past else 0.0)
            for c in stats_cols
        }

        # Добавление текущего матча в историю команд
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

    # Кодирование названий команд в числа
    le = LabelEncoder()
    all_teams = np.unique(df[['HomeTeam', 'AwayTeam']].values)
    le.fit(all_teams)

    df['HomeTeam_encoded'] = le.transform(df['HomeTeam'])
    df['AwayTeam_encoded'] = le.transform(df['AwayTeam'])

    # Сохранение энкодера для использования при инференсе
    with open('team_encoder.pkl', 'wb') as f:
        pickle.dump(le, f)

    # Целевая переменная: A (гости) = 0, D (ничья) = 1, H (хозяева) = 2
    target_mapping = {'A': 0, 'D': 1, 'H': 2}
    df['Target'] = df['FullTimeResult'].map(target_mapping)
    df = df.dropna(subset=['Target'])

    feature_cols = ['HomeTeam_encoded', 'AwayTeam_encoded'] + list(features_df.columns)

    X = df[feature_cols]
    y = df['Target']

    # Сохранение списка признаков для инференса
    with open('feature_cols.json', 'w') as f:
        json.dump(feature_cols, f)

    # Сохранение последней статистики команд для инференса
    latest_team_stats = {team: stats[-window_size:] for team, stats in team_stats.items()}
    with open('latest_team_stats.pkl', 'wb') as f:
        pickle.dump(latest_team_stats, f)

    # Временной сплит: 80% на обучение, 20% на тест
    split_index = int(len(X) * 0.8)
    X_train, X_test = X.iloc[:split_index], X.iloc[split_index:]
    y_train, y_test = y.iloc[:split_index], y.iloc[split_index:]

    return X_train, X_test, y_train, y_test, feature_cols


if __name__ == "__main__":
    X_train, X_test, y_train, y_test, cols = load_and_preprocess_data('epl_final.csv')
    print(f"Данные загружены. Train: {X_train.shape}, Test: {X_test.shape}")
    print("Признаки:", cols)
