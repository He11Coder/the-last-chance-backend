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
            raise RuntimeError("one or more of the .env parameters haven't been specified")
        else:
            self.ADVISER_PORT = int(self.ADVISER_PORT)
            self.MODEL_TEMPERATURE = float(self.MODEL_TEMPERATURE)
