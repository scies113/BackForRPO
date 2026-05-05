"""
Скрипт обучения модели.

Загружает данные, обучает XGBoost-классификатор и выводит
метрики качества (Accuracy, Log Loss, Classification Report).
Сохраняет обученную модель в файл xgboost_model.pkl.
"""

import numpy as np
from sklearn.metrics import accuracy_score, classification_report, log_loss

from dataset import load_and_preprocess_data
from model import FootballPredictionModel


def main():
    """Основной pipeline: загрузка данных → обучение → оценка → сохранение."""
    print("Загрузка и предобработка данных...")
    X_train, X_test, y_train, y_test, feature_cols = load_and_preprocess_data(
        'epl_final.csv', window_size=5
    )

    print("Обучение модели...")
    model = FootballPredictionModel()
    model.train(X_train, y_train)

    print("Оценка модели...")
    y_pred_proba = model.predict_proba(X_test)
    y_pred = np.argmax(y_pred_proba, axis=1)

    acc = accuracy_score(y_test, y_pred)
    loss = log_loss(y_test, y_pred_proba)

    print(f"Accuracy: {acc:.4f}")
    print(f"Log Loss: {loss:.4f}")
    print("\nClassification Report:")
    # Целевая переменная: A (гости) = 0, D (ничья) = 1, H (хозяева) = 2
    target_names = ['Away', 'Draw', 'Home']
    print(classification_report(y_test, y_pred, target_names=target_names))

    print("Сохранение модели...")
    model.save('xgboost_model.pkl')

    print("Обучение завершено! Модель сохранена в 'xgboost_model.pkl'")


if __name__ == "__main__":
    main()
