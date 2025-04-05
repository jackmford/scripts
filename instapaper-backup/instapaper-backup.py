import json
import os
import requests
from requests_oauthlib import OAuth1

def get_env_var(var_name):
    """Retrieve environment variable and raise an error if missing."""
    value = os.getenv(var_name)
    if not value:
        raise ValueError(f"Missing environment variable: {var_name}")
    return value

def get_access_token():
    """Fetch Instapaper OAuth access token."""
    url = "https://www.instapaper.com/api/1/oauth/access_token"
    oauth = OAuth1(
        client_key=CONSUMER_KEY,
        client_secret=CONSUMER_SECRET,
        signature_method="HMAC-SHA1",
    )
    payload = {
        "x_auth_username": USERNAME,
        "x_auth_password": PASSWORD,
        "x_auth_mode": "client_auth",
    }
    
    try:
        response = requests.post(url, auth=oauth, data=payload)
        response.raise_for_status()
        token_data = dict(pair.split("=") for pair in response.text.split("&"))
        return token_data.get("oauth_token"), token_data.get("oauth_token_secret")
    except requests.exceptions.RequestException as e:
        print(f"Failed to obtain access token: {e}")
        exit(1)

def get_bookmarks(access_token, access_token_secret):
    """Retrieve bookmarks from Instapaper."""
    url = "https://www.instapaper.com/api/1/bookmarks/list"
    oauth = OAuth1(
        client_key=CONSUMER_KEY,
        client_secret=CONSUMER_SECRET,
        resource_owner_key=access_token,
        resource_owner_secret=access_token_secret,
        signature_method="HMAC-SHA1",
    )
    
    params = {"folder_id": "archive", "limit": 500}

    try:
        response = requests.get(url, auth=oauth, params=params)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        print(f"Request failed: {e}")
        exit(1)
    except json.JSONDecodeError:
        print("Error: Unable to parse JSON response from Instapaper.")
        exit(1)

def main():
    global CONSUMER_KEY, CONSUMER_SECRET, USERNAME, PASSWORD
    CONSUMER_KEY = get_env_var("INSTAPAPER_KEY")
    CONSUMER_SECRET = get_env_var("INSTAPAPER_SECRET")
    USERNAME = get_env_var("INSTAPAPER_USERNAME")
    PASSWORD = get_env_var("INSTAPAPER_PASSWORD")

    access_token, access_token_secret = get_access_token()
    bookmarks = get_bookmarks(access_token, access_token_secret)

    print(bookmarks)
    for index, item in enumerate(bookmarks):
        if "title" in item and "url" in item:
            print(index, item["title"], item["url"])

if __name__ == "__main__":
    main()
