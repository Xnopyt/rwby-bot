from python_anticaptcha import AnticaptchaClient, NoCaptchaTaskProxylessTask, AnticaptchaException
import discord
import asyncio
import requests
import json
import threading
import time
import random

global msg_queue
msg_queue = list()
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
    channels = [client.get_channel("445190274902261770")]
    print("Connected to discord!")
    bot_state_cache = "NOT SET YET!"
    while True:
        if not bot_state == bot_state_cache:
            await client.change_presence(game=discord.Game(name=bot_state))
            bot_state_cache = bot_state
        if len(msg_queue) > 0:
            for chan in channels:
                await client.send_message(chan, msg_queue[0])
            del msg_queue[0]
        asyncio.sleep(1)

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
    msg_queue.append("Requesting authentication token from rooster teeth...")
    s = requests.Session()
    r=s.get("https://roosterteeth.com/login/")
    client_id = "4338d2b4bdc8db1239360f28e72f0d9ddb1fd01e7a38fbb07b4b1f4ba4564cc5"
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
    msg_queue.append("Requesting account info...")
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
    msg_queue.append("Requesting to begin a FIRST trial...")
    r = requests.post(url,headers=headers,data=data)
    r = requests.get(url,headers=headers)
    sub_uuid = json.loads(r.text[1:-1])["uuid"]
    url = "https://business-service.roosterteeth.com/api/v1/recurly_service/subscriptions/"+sub_uuid+"/cancel"
    msg_queue.append("Requesting to cancle FIRST trial membership subscription...")
    requests.delete(url, headers=headers)

def generate_rt_account(api_key):
    msg_queue.append("Generating a RT account...")
    site_key = '6LeZAyAUAAAAAKXhHLkm7QSka-pPFSRLgL7fjS_g'
    url = 'https://roosterteeth.com/signup/'
    random.seed=(time.ctime())
    msg_queue.append("Generating username and password...")
    email = ''.join(random.choice("ABCDEFGHIJKLMNOPQRSTUVWXYZ") for _ in range(10))+"@how2trianglemuygud.com"
    password = ''.join(random.choice("ABCDEFGHIJKLMNOPQRSTUVWXYZ") for _ in range(10))
    try:
        msg_queue.append("Creating a task for Anti-Captcha...")
        client_anticaptcha = AnticaptchaClient(api_key)
        task = NoCaptchaTaskProxylessTask(url, site_key)
        job = client_anticaptcha.createTask(task)
        job.join()
        msg_queue.append("Waiting for the ReCaptcha solution (This may take several minuets)")
        recaptcha_response = job.get_solution_response()
    except AnticaptchaException as e:
        if e.error_code == 'ERROR_ZERO_BALANCE':
            msg_queue.append("FATAL ERROR: <!@360457422181105666>, your Anti-Captcha balance is empty, topup and restart the bot!")
            recaptcha_response = "NO BAL"
        else:
            recaptcha_response = "UNIDENTIFIED ERROR"
            msg_queue.append("FATAL ERROR: <!@360457422181105666>, Anti-Captcha was unable to solve the ReCaptcha!")
    if recaptcha_response == "NO BAL" or recaptcha_response == "UNIDENTIFIED ERROR":
        return recaptcha_response
    msg_queue.append("Got ReCaptcha solution!")
    data = dict()
    user = dict()
    user["email"] = email
    user["password"] = password
    data["recaptcha_response"] = recaptcha_response
    data["user"] = user
    data = json.dumps(data)
    msg_queue.append("Creating account...")
    r = requests.post("https://business-service.roosterteeth.com/api/v1/users", data=data)
    if "error" in r:
        return False
    else:
        msg_queue.append("The response didn't contain any errors, assuming we are good.")
        return email,password

def parse_latest_video(token):
    msg_queue.append("Parsing the video please wait...")
    r = requests.get("https://svod-be.roosterteeth.com/api/v1/seasons/rwby-volume-6/episodes?order=des&per_page=1")
    json_data = json.loads(r.text)
    headers = {"authorization":token}
    msg_queue.append("Getting the magic numbers...")
    r = requests.get("https://svod-be.roosterteeth.com/api/v1/episodes/" + json_data["data"][0]["uuid"] + "/videos/", headers=headers)
    video = json.loads(r.text)
    end = video["data"][0]["attributes"]["url"][len("https://rtv3-video.roosterteeth.com/store/"):]
    pos = video["data"][0]["attributes"]["url"][len("https://rtv3-video.roosterteeth.com/store/"):].find("/ts/")
    magic_nums = end[:pos]
    return magic_nums

def update_json_info(magic):
    msg_queue.append("Compiling results and storing to the server...")
    magic_short, magic_long = magic
    r = requests.get("https://svod-be.roosterteeth.com/api/v1/seasons/rwby-volume-6/episodes?order=des&per_page=1")
    video = json.loads(r.text)
    epnum = video["data"][0]["attributes"]["season_number"]
    title = video["data"][0]["attributes"]["title"]
    site_info = dict()
    site_info["epnum"] = epnum
    site_info["title"] = title
    site_info["magic_short"] = magic_short
    site_info["magic_long"] = magic_long
    site_info = json.dumps(site_info)
    f = open("rwby_info.json", "w+")
    lines[0] = site_info
    f.writelines(lines)
    f.close()