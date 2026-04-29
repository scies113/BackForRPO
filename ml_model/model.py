import xgboost as xgb
import pickle

class FootballPredictionModel:
    def __init__(self):
        # Using multi:softprob for multi-class probability prediction
        self.model = xgb.XGBClassifier(
            objective='multi:softprob',
            num_class=3,
            n_estimators=100,
            learning_rate=0.1,
            max_depth=5,
            random_state=42
        )
        
    def train(self, X_train, y_train):
        self.model.fit(X_train, y_train)
        
    def predict_proba(self, X):
        return self.model.predict_proba(X)
        
    def save(self, filepath='xgboost_model.pkl'):
        with open(filepath, 'wb') as f:
            pickle.dump(self.model, f)
            
    def load(self, filepath='xgboost_model.pkl'):
        with open(filepath, 'rb') as f:
            self.model = pickle.load(f)
