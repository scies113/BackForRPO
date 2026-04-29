"""
Модуль модели предсказания результатов футбольных матчей.

Обёртка над XGBClassifier для удобного обучения, предсказания
и сериализации модели градиентного бустинга.
"""

import pickle

import xgboost as xgb


class FootballPredictionModel:
    """Модель предсказания исходов футбольных матчей на основе XGBoost.

    Классифицирует матчи на три класса: победа гостей (0),
    ничья (1), победа хозяев (2). Возвращает вероятности
    каждого исхода.
    """

    def __init__(self):
        self.model = xgb.XGBClassifier(
            objective='multi:softprob',
            n_estimators=100,
            learning_rate=0.1,
            max_depth=5,
            random_state=42
        )

    def train(self, X_train, y_train):
        """Обучает модель на переданных данных."""
        self.model.fit(X_train, y_train)

    def predict_proba(self, X):
        """Возвращает вероятности каждого класса для входных данных."""
        return self.model.predict_proba(X)

    def save(self, filepath='xgboost_model.pkl'):
        """Сохраняет обученную модель в файл."""
        with open(filepath, 'wb') as f:
            pickle.dump(self.model, f)

    def load(self, filepath='xgboost_model.pkl'):
        """Загружает модель из файла."""
        with open(filepath, 'rb') as f:
            self.model = pickle.load(f)
