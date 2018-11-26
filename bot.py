from python_anticaptcha import AnticaptchaClient, NoCaptchaTaskProxylessTask, AnticaptchaException
import discord
import asyncio
import requests
import json
import threading
import time
import random

f = open("tokens.txt")
lines = f.readlines()
token = lines[0].rstrip()
card_info = dict()
card_info["first_name"] = lines[1].rstrip()
card_info["last_name"] = lines[2].rstrip()
card_info["postal_code"] = lines[3].rstrip()
card_info["number"] = lines[4].rstrip()
card_info["month"] = lines[5].rstrip()
card_info["year"] = lines[6].rstrip()
card_info["cvv"] = lines[7].rstrip()
card_info["version"] = "4.9.3"
card_info["key"] = "ewr1-2beFfL1PHAOpBH03tu5h6j"
anticaptchaapikey = lines[8].rstrip()
f.close()
client = discord.Client()
loop = asyncio.get_event_loop()
bot_state = "Bot init in progress, please wait..."

@client.event
async def on_ready():
    print("Connected to discord!")
    bot_state_cache = "NOT SET YET!"
    while True:
        if not bot_state == bot_state_cache:
            await client.change_presence(game=discord.Game(name=bot_state))
            bot_state_cache = bot_state
        asyncio.sleep(5)

class Bot(threading.Thread):
    def __init__(self, name):
        threading.Thread.__init__(self)
        self.name = name
    def run(self):
        loop.run_until_complete(client.login(token))
        loop.run_until_complete(client.connect())

bot_thread = Bot("discord_bot")
bot_thread.start()

def grab_oauth(user,password):
    s = requests.Session()
    r=s.get("https://roosterteeth.com/login/")
    client_id = r.text.find("REACT_APP_AUTH_CLIENT_ID: rtConfigSetup('")
    client_id = r.text[client_id+len("REACT_APP_AUTH_CLIENT_ID: rtConfigSetup('"):]
    client_id_end = client_id.find("'")
    client_id = client_id[:client_id_end]
    data = dict()
    data["client_id"] = client_id
    data["grant_type"] = "password"
    data["username"] = user
    data["password"] = password
    data["scope"] = "user public"
    data = json.dumps(data)
    r = s.post("https://auth.roosterteeth.com/oauth/token",data=data)
    response = json.loads(r.text)
    return "Bearer "+ response["access_token"]

def activate_first(token, card_info):
    headers = {"authorization":token}
    r = requests.get("https://business-service.roosterteeth.com/api/v1/me", headers=headers)
    uuid = json.loads(r.text)["id"]
    r = requests.post("https://api.recurly.com/js/v1/token", data=card_info)
    recurly_id = json.loads(r.text)["id"]
    url = "https://business-service.roosterteeth.com/api/v1/recurly_service/accounts/"+uuid+"/subscriptions"
    data = dict()
    subscription = dict()
    subscription["coupon_code"] = None
    subscription["first_name"] = card_info["first_name"]
    subscription["last_name"] = card_info["last_name"]
    subscription["plan_code"] = "1month"
    subscription["recurly_token"] = recurly_id
    data["subscription"] = subscription
    data = json.dumps(data)
    r = requests.post(url,headers=headers,data=data)
    r = requests.get(url,headers=headers)
    sub_uuid = json.loads(r.text[1:-1])["uuid"]
    url = "https://business-service.roosterteeth.com/api/v1/recurly_service/subscriptions/"+sub_uuid+"/cancel"
    requests.delete(url, headers=headers)

def generate_rt_account(api_key):
    site_key = '6LeZAyAUAAAAAKXhHLkm7QSka-pPFSRLgL7fjS_g'
    url = 'https://roosterteeth.com/signup/'
    random.seed=(time.ctime())
    email = ''.join(random.choice("ABCDEFGHIJKLMNOPQRSTUVWXYZ") for _ in range(10))+"@how2trianglemuygud.com"
    password = ''.join(random.choice("ABCDEFGHIJKLMNOPQRSTUVWXYZ") for _ in range(10))
    try:
        client_anticaptcha = AnticaptchaClient(api_key)
        task = NoCaptchaTaskProxylessTask(url, site_key)
        job = client_anticaptcha.createTask(task)
        job.join()
        recaptcha_response = job.get_solution_response()
    except AnticaptchaException as e:
        if e.error_code == 'ERROR_ZERO_BALANCE':
            recaptcha_response = "NO BAL"
        else:
            recaptcha_response = "UNIDENTIFIED ERROR"
    if recaptcha_response == "NO BAL" or recaptcha_response == "UNIDENTIFIED ERROR":
        return recaptcha_response
    data = dict()
    user = dict()
    user["email"] = email
    user["password"] = password
    data["recaptcha_response"] = recaptcha_response
    data["user"] = user
    data = json.dumps(data)
    r = requests.post("https://business-service.roosterteeth.com/api/v1/users", data=data)
    if "error" in r:
        return False
    else:
        return email,password