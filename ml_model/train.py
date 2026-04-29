from dataset import load_and_preprocess_data
from model import FootballPredictionModel
from sklearn.metrics import accuracy_score, classification_report, log_loss
import numpy as np
import json

def main():
    print("Loading and preprocessing data...")
    X_train, X_test, y_train, y_test, feature_cols = load_and_preprocess_data('epl_final.csv', window_size=5)
    
    print("Training model...")
    model = FootballPredictionModel()
    model.train(X_train, y_train)
    
    print("Evaluating model...")
    y_pred_proba = model.predict_proba(X_test)
    y_pred = np.argmax(y_pred_proba, axis=1)
    
    acc = accuracy_score(y_test, y_pred)
    loss = log_loss(y_test, y_pred_proba)
    
    print(f"Accuracy: {acc:.4f}")
    print(f"Log Loss: {loss:.4f}")
    print("\nClassification Report:")
    # Mapping target Variable (A=0, D=1, H=2)
    target_names = ['Away', 'Draw', 'Home']
    print(classification_report(y_test, y_pred, target_names=target_names))
    
    print("Saving model...")
    model.save('xgboost_model.pkl')
    
    # Save feature names for reference during inference
    with open('feature_cols.json', 'w') as f:
        json.dump(feature_cols, f)
        
    print("Training complete! Model saved to 'xgboost_model.pkl'")

if __name__ == "__main__":
    main()
