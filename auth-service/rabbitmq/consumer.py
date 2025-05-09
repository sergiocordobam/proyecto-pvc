import pika
import json
import os
import time
import requests

def handle_register_citizen(data):
    print("Handling register citizen:", data, flush=True)

    resp_json = {
        "full_name": data["full_name"],
        "document_id": str(data["document_id"]),
        "document_type": "CC",
        "address": "Carrera 1 # 2 - 3",
        "phone": "3003003030",
        "email": data["email"],
        "password": data["password"],
        "terms_accepted": data["terms_accepted"]
    }
    
    try:
        response = requests.post("http://auth-service:8000/auth/register", json=resp_json)
        print("Forwarded data, status:", response.status_code, flush=True)
        print("Response body:", response.text, flush=True)
    except Exception as e:
        print("Failed to insert user:", str(e), flush=True)

def handle_delete_citizen(data):
    print("Handling citizen delete:", data, flush=True)

    resp_json = {
        "document_id": data["citizenId"]
    }

    try:
        response = requests.post("http://auth-service:8000/auth/delete_user", json=resp_json)
        print("Forwarded data, status:", response.status_code, flush=True)
        print("Response body:", response.text, flush=True)
    except Exception as e:
        print("Failed to delete user:", str(e), flush=True)    

def start_rabbitmq_consumer():
    def callback_register(ch, method, properties, body):
        data = json.loads(body)
        handle_register_citizen(data)

    def callback_delete(ch, method, properties, body):
        data = json.loads(body)
        handle_delete_citizen(data)

    rabbitmq_user = os.getenv('RABBITMQ_DEFAULT_USER', 'guest')
    rabbitmq_password = os.getenv('RABBITMQ_DEFAULT_PASS', 'guest')
    rabbitmq_host = os.getenv('RABBITMQ_HOST', 'rabbitmq')

    credentials = pika.PlainCredentials(rabbitmq_user, rabbitmq_password)

    retries = 5
    for _ in range(retries):
        try:
            print("Waiting 15s to ensure RabbitMQ is ready...")
            time.sleep(15)

            print(f"Attempting to connect to RabbitMQ at {rabbitmq_host}...", flush=True)
            print(f"Attempting to connect to RabbitMQ at {rabbitmq_password}...", flush=True)
            print(f"Attempting to connect to RabbitMQ at {rabbitmq_user}...", flush=True)

            connection = pika.BlockingConnection(pika.ConnectionParameters(host="rabbitmq", credentials=credentials, heartbeat=600, blocked_connection_timeout=300))
            print("conn", flush=True)

            channel = connection.channel()
            print("channel", flush=True)

            channel.queue_declare(queue='register_citizen_queue', durable=True)
            channel.queue_declare(queue='delete_citizen_queue', durable=True)
            print("queue declare", flush=True)

            channel.basic_qos(prefetch_count=1)
            print("prefetch", flush=True)

            channel.basic_consume(queue='register_citizen_queue', on_message_callback=callback_register, auto_ack=True)
            channel.basic_consume(queue='delete_citizen_queue', on_message_callback=callback_delete, auto_ack=True)
            print("channel consume", flush=True)

            print(" [*] Waiting for messages. To exit press CTRL+C", flush=True)
            channel.start_consuming()

            break
        except Exception as e:
            print(f"Connection failed: {e.__class__.__name__} - {e}. Retrying...", flush=True)
            time.sleep(5)

    else:
        print("Failed to connect to RabbitMQ after several attempts. Exiting.", flush=True)
