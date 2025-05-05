from firebase_admin import firestore
from firebase.firebase_initialization import initialize_firebase

initialize_firebase()
db = firestore.client()

class UserExists:
    def __init__(self):
        self.db = db

    def user_exists(self, document_id):
        users_ref = self.db.collection('users')
        query = users_ref.where('document_id', '==', document_id).stream()
        return any(query)