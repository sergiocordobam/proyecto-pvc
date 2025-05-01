import re
import firebase_admin
from firebase_admin import credentials, auth, firestore
import requests

cred = credentials.Certificate('serviceAccountKey.json')
firebase_admin.initialize_app(cred)
db = firestore.client()

class User:
    def __init__(self):
        self.db = db

    def is_valid_password(self, password):
        pattern = r'^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&\#]{8,}$'
        return re.match(pattern, password)

    def register(self, data):
        required_fields = ['full_name', 'document_id', 'document_type', 'address', 'phone', 'email', 'password', 'terms_accepted']
        if not all(field in data for field in required_fields):
            return {'error': 'Missing required fields'}, 400

        if not data['terms_accepted']:
            return {'error': 'Terms and conditions must be accepted'}, 400

        if not self.is_valid_password(data['password']):
            return {'error': 'Password does not meet security requirements'}, 400

        users_ref = self.db.collection('users')
        query = users_ref.where('document_id', '==', data['document_id']).stream()
        if any(query):
            return {'error': 'User with this document ID already exists'}, 400

        try:
            user_record = auth.create_user(
                email=data['email'],
                password=data['password'],
                display_name=data['full_name']
            )

            users_ref.document(user_record.uid).set({
                'full_name': data['full_name'],
                'document_id': data['document_id'],
                'document_type': data['document_type'],
                'address': data['address'],
                'phone': data['phone'],
                'email': data['email'],
            })

            # govcarpeta_response = requests.post(
            #     'https://govcarpeta-apis-4905ff3c005b.herokuapp.com/apis/registerCitizen',
            #     json={
            #         'id': data['document_id'],
            #         'name': data['full_name'],
            #         'address': data['address'],
            #         'email': data['email'],
            #         'operatorId': data['operator_id'],
            #         'operatorName': data['operator_name'],
            #     }
            # )

            # if govcarpeta_response.status_code != 200:
            #     return {'error': 'Failed to register with GovCarpeta'}, 500

            return {'message': 'User registered successfully'}, 201

        except Exception as e:
            return {'error': str(e)}, 500
