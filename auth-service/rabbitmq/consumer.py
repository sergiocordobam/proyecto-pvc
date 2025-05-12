import pika
import json
import os
import time
import requests
import threading

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
        response = requests.post("http://auth-service:5000/register", json=resp_json)
        print("Forwarded data, status:", response.status_code, flush=True)
        print("Response body:", response.text, flush=True)
    except Exception as e:
        print("Failed to insert user:", str(e), flush=True)

def handle_delete_citizen(data):
    print("Handling citizen delete:", data, flush=True)
    resp_json = {
        "document_id": str(data["citizenId"])
    }
    try:
        response = requests.post("http://auth-service:5000/delete_user", json=resp_json)
        print("Forwarded data, status:", response.status_code, flush=True)
        print("Response body:", response.text, flush=True)
    except Exception as e:
        print("Failed to delete user:", str(e), flush=True)

def start_consumer(queue_name, callback):
    rabbitmq_user = os.getenv('RABBITMQ_DEFAULT_USER', 'guest')
    rabbitmq_password = os.getenv('RABBITMQ_DEFAULT_PASS', 'guest')
    rabbitmq_host = os.getenv('RABBITMQ_HOST', 'rabbitmq')

    credentials = pika.PlainCredentials(rabbitmq_user, rabbitmq_password)

    for _ in range(5):
        try:
            print(f"[{queue_name}] Waiting for RabbitMQ...", flush=True)
            time.sleep(30)

            connection = pika.BlockingConnection(pika.ConnectionParameters(
                host=rabbitmq_host,
                credentials=credentials,
                heartbeat=600,
                blocked_connection_timeout=300
            ))
            channel = connection.channel()
            channel.queue_declare(queue=queue_name, durable=True)
            channel.basic_qos(prefetch_count=1)
            channel.basic_consume(queue=queue_name, on_message_callback=callback, auto_ack=True)

            print(f"[{queue_name}] Listening for messages...", flush=True)
            channel.start_consuming()
            break

        except Exception as e:
            print(f"[{queue_name}] Connection failed: {e}. Retrying...", flush=True)
            time.sleep(5)
    else:
        print(f"[{queue_name}] Failed to connect after several attempts.", flush=True)

def start_all_consumers():
    threading.Thread(
        target=start_consumer,
        args=('register_citizen_queue', lambda ch, method, props, body: handle_register_citizen(json.loads(body))),
        daemon=True
    ).start()

    threading.Thread(
        target=start_consumer,
        args=('delete_citizen_queue', lambda ch, method, props, body: handle_delete_citizen(json.loads(body))),
        daemon=True
    ).start()
