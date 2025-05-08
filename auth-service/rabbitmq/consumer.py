import pika
import json
import os
import threading

def handle_register_citizen(data):
    print("Handling register citizen:", data)
    # TODO: Add logic to insert user into DB or perform other actions

# def handle_delete_documents(data):
#     print("Handling delete documents:", data)
#     # TODO: Add logic to delete documents from DB

def start_rabbitmq_consumer():
    def callback_register(ch, method, properties, body):
        data = json.loads(body)
        handle_register_citizen(data)

    # def callback_delete(ch, method, properties, body):
    #     data = json.loads(body)
    #     handle_delete_documents(data)

    rabbitmq_user = os.getenv('RABBITMQ_USER', 'guest')
    rabbitmq_password = os.getenv('RABBITMQ_PASSWORD', 'guest')
    rabbitmq_host = os.getenv('RABBITMQ_HOST', 'localhost')

    credentials = pika.PlainCredentials(rabbitmq_user, rabbitmq_password)
    connection = pika.BlockingConnection(pika.ConnectionParameters(host=rabbitmq_host, credentials=credentials))
    channel = connection.channel()

    channel.queue_declare(queue='register_citizen_queue', durable=True)
    # channel.queue_declare(queue='delete_documents_queue', durable=True)

    channel.basic_consume(queue='register_citizen_queue', on_message_callback=callback_register, auto_ack=True)
    # channel.basic_consume(queue='delete_documents_queue', on_message_callback=callback_delete, auto_ack=True)

    print(" [*] Waiting for messages. To exit press CTRL+C")
    channel.start_consuming()

def run_consumer_in_background():
    thread = threading.Thread(target=start_rabbitmq_consumer, daemon=True)
    thread.start()
