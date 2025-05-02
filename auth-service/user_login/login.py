from firebase_admin import auth

class UserLogin:
    def __init__(self):
        pass

    def verify_token(self, id_token):
        if not id_token:
            return {"error": "ID token is required"}, 400

        try:
            decoded_token = auth.verify_id_token(id_token)
            uid = decoded_token["uid"]
            return {"message": "Token is valid", "uid": uid}, 200
        except Exception as e:
            return {"error": "Invalid token", "details": str(e)}, 401
