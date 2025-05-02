from flask import Flask, request, jsonify
from flask_cors import CORS, cross_origin
from user_register.register import UserRegister
from user_delete.delete import UserDelete

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

if __name__ == "__main__":
    app.run(debug=True)
