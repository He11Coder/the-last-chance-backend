from yandex_cloud_ml_sdk import YCloudML
from flask import Flask, jsonify, request

import sys
sys.path.append('..')
from env_loader import AdviserParameters

#Initialize Flask parameters and Yandex Cloud connection
params = AdviserParameters()

app = Flask(__name__)

sdk = YCloudML(folder_id=params.YACLOUD_FOLDER_ID, auth=params.YAGPT_API_KEY)
model = sdk.models.completions('yandexgpt')
model = model.configure(temperature=params.MODEL_TEMPERATURE)


def get_advice(animal: str, prompt: str):
    """Get advice for 'animal' from YandexGPT by 'prompt'"""

    full_prompt = f"Дай мне несколько рекомендаций по уходу за домашним животным ({animal}) исходя из его описания: " + prompt

    result = model.run(full_prompt)
    for alternative in result:
        text_result = alternative.text

    return {"animal": animal, "prompt": prompt, "advice": text_result}


@app.route("/get_advice", methods=["GET"])
def handle_get_advice():
    """
    Handles GET requests to /get_advice.
    """
    try:
        animal = request.args.get("animal", default="", type=str)
        prompt = request.args.get("prompt", default="", type=str)
        if (not prompt) or prompt == "":
            return jsonify({"error": "missing required query parameter 'prompt'"}), 400
        if (not animal) or animal == "":
            return jsonify({"error": "missing required query parameter 'animal'"}), 400

        advice = get_advice(animal, prompt)
        return jsonify(advice), 200  # Return advice as JSON with status 200 (OK)
    except Exception as e:
        error_response = {"error": "An error occurred", "details": str(e)}
        return jsonify(error_response), 500  # Return error as JSON with status 500 (Server Error)


if __name__ == "__main__":
    app.run(host=params.ADVISER_HOST, port=params.ADVISER_PORT)