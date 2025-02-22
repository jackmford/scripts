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

if response.status_code == 200:
    # The response comes as a URL-encoded string like:
    # "oauth_token=xxx&oauth_token_secret=xxx"
    token_data = dict(pair.split("=") for pair in response.text.split("&"))
    print("Access Token:", token_data.get("oauth_token"))
    print("Token Secret:", token_data.get("oauth_token_secret"))
else:
    print("Error:", response.status_code)
    print(response.text)
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
payload = {"folder_id": "archive"}
response = requests.post(url, auth=oauth, data=payload)
data = json.loads(response.text)

for data, i in enumerate(data):
    if 'title' in i and 'url' in i:
        print(data, i['title'], i['url'])


