import re
from firebase_admin import auth, firestore
import requests
from firebase.firebase_initialization import initialize_firebase

initialize_firebase()
db = firestore.client()

class UserRegister:
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

            # user_record = auth.create_user(
            #     email=data['email'],
            #     password=data['password'],
            #     display_name=data['full_name']
            # )

            # users_ref.document(user_record.uid).set({
            #     'full_name': data['full_name'],
            #     'document_id': data['document_id'],
            #     'document_type': data['document_type'],
            #     'address': data['address'],
            #     'phone': data['phone'],
            #     'email': data['email'],
            # })

            operator_info = requests.get("http://interoperator-service:3000/comunication/operators/self")
            resp = operator_info.json()
            print("resp", resp, flush=True)
            # operator_id = operator_info["_id"]
            # operator_name = operator_info["operatorName"]
            # print("op id", operator_id, flush=True)
            # print("op name", operator_name, flush=True)

            # govcarpeta_response = requests.post(
            #     'https://govcarpeta-apis-4905ff3c005b.herokuapp.com/apis/registerCitizen',
            #     json={
            #         'id': data['document_id'],
            #         'name': data['full_name'],
            #         'address': data['address'],
            #         'email': data['email'],
            #         'operatorId': data['operator_id'],
            #         'operatorName': "carpeta_PVC",
            #     }
            # )

            # if govcarpeta_response.status_code == 500:
            #     return {'error': 'GovCarpeta error'}, 500
            # elif govcarpeta_response.status_code == 501:
            #     return {'error': 'El ciudadano ya se encuentra registrado'}

            return {'message': 'User registered successfully'}, 201

        except Exception as e:
            return {'error': str(e)}, 500
