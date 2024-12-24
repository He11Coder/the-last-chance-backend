import numpy as np
import base64
from io import BytesIO
import json
from PIL import Image, ImageFile
import opennsfw2 as n2
#from flask import Flask, jsonify, request
import pika

import sys
sys.path.append('..')
from env_loader import NSFWValidatorParameters

params = NSFWValidatorParameters()

#app = Flask(__name__)

# The model is a `keras_core.Model` object.
model = n2.make_open_nsfw_model(weights_path=params.WEIGHTS_FILE_PATH)

connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost', port=8001))
channel = connection.channel()

channel.queue_declare(queue='nsfw_validation_queue', durable=True)


def nsfw_validator(pil_image: ImageFile):
    image = n2.preprocess_image(pil_image, n2.Preprocessing.SIMPLE)
    inputs = np.expand_dims(image, axis=0)  # Add batch axis (for single image).

    predictions = model.predict(inputs, verbose=2)
    sfw_probability, nsfw_probability = predictions[0]

    return (True, sfw_probability) if sfw_probability > nsfw_probability else (False, nsfw_probability)


def on_request(ch, method, props, body):
    try:
        data = json.loads(body)

        base64_image = data.get("image")
        if not base64_image:
            err_message = json.dumps({"error": "image data is required", "code": 400}).encode('utf-8')
            ch.basic_publish(exchange='',
                             routing_key=props.reply_to,
                             properties=pika.BasicProperties(correlation_id = props.correlation_id,
                                                             delivery_mode=pika.DeliveryMode.Persistent),
                             body=err_message)
            return

        image_bytes = base64.b64decode(base64_image)
        pil_image = Image.open(BytesIO(image_bytes))
        result = nsfw_validator(pil_image)
        res_to_return = json.dumps({"is_safe": result[0], "confidence": float(result[1]), "code": 200}).encode('utf-8')
        ch.basic_publish(exchange='',
                         routing_key=props.reply_to,
                         properties=pika.BasicProperties(correlation_id = props.correlation_id,
                                                         delivery_mode=pika.DeliveryMode.Persistent),
                         body=res_to_return)

    except Exception as e:
        err_message = json.dumps({"error": str(e), "code": 500}).encode('utf-8')
        ch.basic_publish(exchange='',
                         routing_key=props.reply_to,
                         properties=pika.BasicProperties(correlation_id = props.correlation_id,
                                                         delivery_mode=pika.DeliveryMode.Persistent),
                         body=err_message)

    finally:
        ch.basic_ack(delivery_tag=method.delivery_tag)
    

channel.basic_qos(prefetch_count=1)
channel.basic_consume(queue='nsfw_validation_queue', on_message_callback=on_request)

#channel.start_consuming()

'''@app.route("/validate_nsfw", methods=["POST"])
def validate_nsfw():
    """
    Handles POST requests to /validate_nsfw.
    """
    try:
        body_data = request.get_json()

        base64_image = body_data.get("image")
        if not base64_image:
            return jsonify({"error": "image data is required"}), 400
        
        image_bytes = base64.b64decode(base64_image)
        pil_image = Image.open(BytesIO(image_bytes))
        result = nsfw_validator(pil_image)

        return jsonify({"is_safe": result[0], "confidence": float(result[1])}), 200
    except Exception as e:
        error_response = {"error": "An error occurred", "details": str(e)}
        return jsonify(error_response), 500  # Return error as JSON with status 500 (Server Error)'''


if __name__ == "__main__":
    print("Awaiting RPC requests")
    channel.start_consuming()