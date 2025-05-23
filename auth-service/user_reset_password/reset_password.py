from firebase_admin import auth
from firebase_admin import firestore
from firebase.firebase_initialization import initialize_firebase
import requests

initialize_firebase()
db = firestore.client()

class UserResetPassword:
    def __init__(self):
        self.db = db

    def update_user_info(self, email, address=None, phone=None, document_type=None, password=None):
        users_ref = self.db.collection('users')
        query = users_ref.where("email", "==", email).limit(1).get()
        if not query:
            return "User not found"

        user_doc = query[0]
        update_data = {}

        if address is not None:
            update_data["address"] = address
        if phone is not None:
            update_data["phone"] = phone
        if document_type is not None:
            update_data["document_type"] = document_type

        if password is not None:
            try:
                user_record = auth.get_user_by_email(email)
                auth.update_user(user_record.uid, password=password)
                print("Firebase Auth password updated.", flush=True)
                return "Firebase Auth password updated"
            except Exception as e:
                print("Failed to update password in Firebase Auth:", str(e), flush=True)
                return f"Failed to update Firebase password: {str(e)}"

        if update_data:
            self.db.collection('users').document(user_doc.id).update(update_data)
            print("Firestore user data updated:", update_data, flush=True)

            send_email = requests.post("http://auth-service:5000/publish_notifications",
                                       json={
                                            "event": "register",
                                            "user": 1234567890,
                                            "name": "Reset Password",
                                            "user_email": email,
                                            "extra_data": {
                                                "title": "Actualiza tu contraseña por favor",
                                                "body": "http://localhost:8000"
                                            }
                                        })
            return f"Firestore user data updated: {update_data}"
        else:
            print("No Firestore fields to update.", flush=True)
            return "No Firestore fields to update"
