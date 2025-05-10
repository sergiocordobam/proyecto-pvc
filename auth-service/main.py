from flask import Flask, request, jsonify
from flask_cors import CORS, cross_origin
from user_register.register import UserRegister
from user_delete.delete import UserDelete
from user_login.login import UserLogin
from user_exists.exists import UserExists
from user_reset_password.reset_password import UserResetPassword
from google.cloud import pubsub_v1
import json
import os

os.environ["GOOGLE_APPLICATION_CREDENTIALS"] = "./serviceAccountKey2.json"

app = Flask(__name__)
CORS(app)

PROJECT_ID = "zeta-matrix-458323-p1"
TOPIC_ID = "auth-api-topic"

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

@app.route("/reset_password", methods=["POST"])
@cross_origin()
def reset_password():
    data = request.get_json()
    email = data["email"]
    document_type = data["document_type"]
    phone = data["phone"]
    address = data["address"]
    password = data["password"]

    user_reset_password_service = UserResetPassword()
    reset_password = user_reset_password_service.update_user_info(email, address, phone, document_type, password)

    return jsonify({
        "message": reset_password
    }), 200

@app.route("/publish_notifications", methods=["POST"])
@cross_origin()
def publish_notifications():
    publisher = pubsub_v1.PublisherClient()
    topic_path = publisher.topic_path(PROJECT_ID, TOPIC_ID)

    try:
        data = request.get_json()
        message_json = json.dumps(data)
        future = publisher.publish(topic_path, message_json.encode("utf-8"))
        message_id = future.result()

        return jsonify({"message": "Message published", "message_id": message_id}), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    from rabbitmq.consumer import start_all_consumers
    start_all_consumers()
    app.run(debug=True, host="0.0.0.0", port=5000)
