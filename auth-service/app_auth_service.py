from flask import Flask, request, jsonify
from user.user import User

app = Flask(__name__)

user_service = User()

@app.route("/register", methods=["POST"])
def register():
    data = request.get_json()
    response, status = user_service.register(data)

    return jsonify(response), status

if __name__ == "__main__":
    app.run(debug=True)
