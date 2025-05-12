from firebase_admin import auth, firestore
from firebase_admin.exceptions import FirebaseError
from firebase.firebase_initialization import initialize_firebase

initialize_firebase()
db = firestore.client()

class UserDelete:
    def __init__(self):
        self.db = db

    def delete_user_by_document_id(self, document_id):
        try:
            users_ref = self.db.collection('users')
            query = users_ref.where('document_id', '==', document_id).stream()
            print("query", flush=True)

            user_found = None
            for user in query:
                user_found = user
                print("for", flush=True)

            if not user_found:
                return {'error': 'User not found with this document ID'}, 404

            user_uid = user_found.id

            auth.delete_user(user_uid)
            print("user_delete", flush=True)

            self.db.collection('users').document(user_uid).delete()

            return {'message': 'User deleted successfully'}, 200

        except FirebaseError as e:
            return {'error': f'Error deleting user: {str(e)}'}, 500