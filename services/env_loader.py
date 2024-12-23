import os
from dotenv import load_dotenv

load_dotenv()

class AdviserParameters:
    """Loads environment variables for GPT Adviser Service"""

    def __init__(self):
        self.YACLOUD_FOLDER_ID = os.getenv("YACLOUD_FOLDER_ID")
        self.YAGPT_API_KEY = os.getenv("YAGPT_API_KEY")
        self.ADVISER_HOST = os.getenv("ADVISER_HOST")
        self.ADVISER_PORT = os.getenv("ADVISER_PORT")
        self.MODEL_TEMPERATURE = os.getenv("MODEL_TEMPERATURE")

        if self.YACLOUD_FOLDER_ID == None or self.YAGPT_API_KEY == None or self.ADVISER_HOST == None or self.ADVISER_PORT == None or self.MODEL_TEMPERATURE == None:
            raise RuntimeError("one or more of the .env parameters for Adviser Service haven't been specified")
        else:
            self.ADVISER_PORT = int(self.ADVISER_PORT)
            self.MODEL_TEMPERATURE = float(self.MODEL_TEMPERATURE)


class NSFWValidatorParameters:
    """Loads environment variables for NSFW Validator Service"""

    def __init__(self):
        self.WEIGHTS_FILE_PATH = os.getenv("WEIGHTS_FILE_PATH")
        self.NSFW_VALIDATOR_HOST = os.getenv("NSFW_VALIDATOR_HOST")
        self.NSFW_VALIDATOR_PORT = os.getenv("NSFW_VALIDATOR_PORT")

        if self.WEIGHTS_FILE_PATH == None or self.NSFW_VALIDATOR_HOST == None or self.NSFW_VALIDATOR_PORT == None:
            raise RuntimeError("one or more of the .env parameters for NSFW Validator Service haven't been specified")
        else:
            self.NSFW_VALIDATOR_PORT = int(self.NSFW_VALIDATOR_PORT)