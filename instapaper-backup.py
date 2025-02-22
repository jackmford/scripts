import json
import os
from pprint import pprint

import requests
from requests_oauthlib import OAuth1

# Replace these values with your own credentials
CONSUMER_KEY = os.environ["INSTAPAPER_KEY"]
CONSUMER_SECRET = os.environ["INSTAPAPER_SECRET"]
USERNAME = os.environ["INSTAPAPER_USERNAME"]
PASSWORD = os.environ["INSTAPAPER_PASSWORD"]

if not all([CONSUMER_KEY, CONSUMER_SECRET, USERNAME, PASSWORD]):
    raise ValueError("Missing required environment variables for Instapaper API authentication.")

url = "https://www.instapaper.com/api/1/oauth/access_token"

# OAuth1 automatically adds the required OAuth 1.0 parameters
oauth = OAuth1(
    client_key=CONSUMER_KEY, client_secret=CONSUMER_SECRET, signature_method="HMAC-SHA1"
)

# Data payload for x_auth mode
payload = {
    "x_auth_username": USERNAME,
    "x_auth_password": PASSWORD,
    "x_auth_mode": "client_auth",
}

response = requests.post(url, auth=oauth, data=payload)
try:
    response.raise_for_status()
    token_data = dict(pair.split("=") for pair in response.text.split("&"))
except requests.exceptions.HTTPError as e:
    print(f"HTTP Error: {e.response.status_code} - {e.response.text}")
    exit(1)
except Exception as e:
    print(f"Error processing response: {e}")
    exit(1)

access_token = token_data.get("oauth_token")
access_token_secret = token_data.get("oauth_token_secret")

if not access_token or not access_token_secret:
    print("Error: Failed to retrieve access token.")
    exit(1)

url = "https://www.instapaper.com/api/1/bookmarks/list"

oauth = OAuth1(
    client_key=CONSUMER_KEY,
    client_secret=CONSUMER_SECRET,
    resource_owner_key=token_data.get("oauth_token"),
    resource_owner_secret=token_data.get("oauth_token_secret"),
    signature_method="HMAC-SHA1",
)

# Data payload for x_auth mode
payload = {"folder_id": "archive", "limit": 500}
response = requests.post(url, auth=oauth, data=payload)
data = json.loads(response.text)

for index, item in enumerate(data):
    if "title" in item and "url" in item:
        print(index, item["title"], item["url"])

