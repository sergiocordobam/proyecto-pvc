from flask import Flask, request, jsonify
from flask_cors import CORS, cross_origin
from user_register.register import UserRegister
from user_delete.delete import UserDelete
from user_login.login import UserLogin
from user_exists.exists import UserExists

app = Flask(__name__)
CORS(app)

@app.route("/register", methods=["POST"])
@cross_origin()
def register():
    data = request.get_json()
    user_service = UserRegister()

    response, status = user_service.register(data)

    return jsonify(response), status

@app.route("/delete_user", methods=["POST"])
@cross_origin()
def delete_user():
    data = request.json
    document_id = data.get('document_id')

    if not document_id:
        return jsonify({'error': 'Missing document_id'}), 400

    user_delete = UserDelete()
    result, status_code = user_delete.delete_user_by_document_id(document_id)

    return jsonify(result), status_code

@app.route("/verify_token", methods=["POST"])
def verify_token():
    data = request.get_json()
    id_token = data.get("idToken")

    login_handler = UserLogin()
    
    response, status = login_handler.verify_token(id_token)
    return jsonify(response), status

@app.route("/user_exists", methods=["POST"])
@cross_origin()
def user_exists():
    data = request.get_json()
    document_id = data.get('document_id')

    if not document_id:
        return jsonify({'error': 'Missing document_id'}), 400

    user_exists_service = UserExists()
    exists = user_exists_service.user_exists(document_id)

    return jsonify({'exists': exists}), 200

if __name__ == "__main__":
    app.run(debug=True)
