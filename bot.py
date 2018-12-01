from python_anticaptcha import AnticaptchaClient, NoCaptchaTaskProxylessTask, AnticaptchaException
import discord
import asyncio
import requests
import json
import threading
import time
import random

current = False
bot_state_queue = asyncio.Queue()
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
loop = asyncio.new_event_loop()
client = discord.Client()

class Bot(threading.Thread):
    def __init__(self, name, token, bot_state_queue, client, loop):
        threading.Thread.__init__(self)
        self.name = name
        self.token = token
        self.bot_state_queue = bot_state_queue
        self.client = client
    def run(self):
        asyncio.set_event_loop(loop)
        @client.event
        async def on_ready():
            print("Connected!")
            channel = self.client.get_channel("433195705486540800")
            while True:
                update_state = await self.bot_state_queue.get()
                bot_state = update_state[0]
                state_type = update_state[1]
                if state_type == 0:
                    try:
                        await self.client.send_message(channel,bot_state)
                    except:
                        pass
                if state_type == 1:
                    try:
                        await self.client.change_presence(game=bot_state)
                    except:
                        pass
                state_type = False

        self.client.run(self.token)


bot_thread = Bot("discord_bot", token, bot_state_queue, client, loop)
bot_thread.start()

def send(msg):
    bot_state_queue.put_nowait([msg,0])

def change_game(game):
    bot_state_queue.put_nowait([game,1])

def grab_oauth(user,password):
    send("Requesting authentication token from rooster teeth...")
    s = requests.Session()
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
    send("Requesting account info...")
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
    send("Requesting to begin a FIRST trial...")
    r = requests.post(url,headers=headers,data=data)
    r = requests.get(url,headers=headers)
    sub_uuid = json.loads(r.text[1:-1])["uuid"]
    url = "https://business-service.roosterteeth.com/api/v1/recurly_service/subscriptions/"+sub_uuid+"/cancel"
    send("Requesting to cancel FIRST trial membership subscription...")
    requests.delete(url, headers=headers)

def generate_rt_account(api_key):
    send("Generating a RT account...")
    site_key = '6LeZAyAUAAAAAKXhHLkm7QSka-pPFSRLgL7fjS_g'
    url = 'https://roosterteeth.com/signup/'
    random.seed=(time.ctime())
    send("Generating username and password...")
    email = ''.join(random.choice("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890") for _ in range(15))+"@how2trianglemuygud.com"
    password = ''.join(random.choice("ABCDEFGHIJKLMNOPQRSTUVWXYZ") for _ in range(10))
    try:
        send("Waiting for the ReCaptcha solution (This may take several minutes)")
        client_anticaptcha = AnticaptchaClient(api_key)
        task = NoCaptchaTaskProxylessTask(url, site_key)
        job = client_anticaptcha.createTask(task)
        job.join()
        recaptcha_response = job.get_solution_response()
    except AnticaptchaException as e:
        if e.error_code == 'ERROR_ZERO_BALANCE':
            send("FATAL ERROR: <!@360457422181105666>, your Anti-Captcha balance is empty, topup and restart the bot!")
            recaptcha_response = "NO BAL"
        else:
            recaptcha_response = "UNIDENTIFIED ERROR"
            send("FATAL ERROR: <!@360457422181105666>, Anti-Captcha was unable to solve the ReCaptcha!")
    if recaptcha_response == "NO BAL" or recaptcha_response == "UNIDENTIFIED ERROR":
        return recaptcha_response
    send("Got ReCaptcha solution!")
    data = dict()
    user = dict()
    user["email"] = email
    user["password"] = password
    data["recaptcha_response"] = recaptcha_response
    data["user"] = user
    data = json.dumps(data)
    send("Creating account...")
    r = requests.post("https://business-service.roosterteeth.com/api/v1/users", data=data)
    if "error" in r:
        return False
    else:
        send("The response didn't contain any errors, assuming we are good.")
        return email,password

def parse_latest_video(token):
    send("Parsing the video please wait...")
    r = requests.get("https://svod-be.roosterteeth.com/api/v1/seasons/rwby-volume-6/episodes?order=des&per_page=1")
    json_data = json.loads(r.text)
    headers = {"authorization":token}
    send("Getting the magic numbers...")
    r = requests.get("https://svod-be.roosterteeth.com/api/v1/episodes/" + json_data["data"][0]["uuid"] + "/videos/", headers=headers)
    video = json.loads(r.text)
    end = video["data"][0]["attributes"]["url"][len("https://rtv3-video.roosterteeth.com/store/"):]
    pos = video["data"][0]["attributes"]["url"][len("https://rtv3-video.roosterteeth.com/store/"):].find("/ts/")
    magic_nums = end[:pos]
    return magic_nums

def update_json_info(magic):
    send("Compiling results and storing to the server...")
    magic_long, magic_short = magic.split("-")
    r = requests.get("https://svod-be.roosterteeth.com/api/v1/seasons/rwby-volume-6/episodes?order=des&per_page=1")
    video = json.loads(r.text)
    epnum = video["data"][0]["attributes"]["number"]
    title = video["data"][0]["attributes"]["title"]
    site_info = dict()
    site_info["epnum"] = epnum
    site_info["title"] = title
    site_info["magic_short"] = magic_short
    site_info["magic_long"] = magic_long
    site_info = json.dumps(site_info)
    f = open("rwby_info.json", "w+")
    lines = list()
    lines.append(site_info)
    f.writelines(lines)
    f.close()
    
def update():
    send("Updating local stream data....")
    login = generate_rt_account(anticaptchaapikey)
    if not type(login) == tuple:
        send("Something went wrong...")
        return
    user, password = login
    token = grab_oauth(user, password)
    activate_first(token,card_info)
    magic = parse_latest_video(token)
    update_json_info(magic)

def setup():
    global current
    change_game(discord.Game(name="Bot init in progress..."))
    send("Initalizing...")
    send("Comparing stored info against the server...")
    r = requests.get("https://svod-be.roosterteeth.com/api/v1/seasons/rwby-volume-6/episodes?order=des&per_page=1")
    json_data = json.loads(r.text)
    current = json_data["data"][0]["uuid"]
    f = open("rwby_info.json")
    lines = f.readlines()
    f.close()
    info = json.loads(lines[0].rstrip())
    if not info["epnum"] == json_data["data"][0]["attributes"]["number"]:
        send("Local data is out of data, grabbing latest episode and updating H2TMG.")
        update()
        send("Done! Init finished!")
    else:
        send("Everything looks good! Init finished!")
    f = open("rwby_info.json")
    lines = f.readlines()
    f.close()
    info = json.loads(lines[0].rstrip())
    change_game(discord.Game(name="Ep " + str(info["epnum"]) + " - " + info["title"]))

def wait():
    while True:
        time_current = time.localtime()
        if time_current.tm_wday == 5 and time_current.tm_hour > 10:
            r = requests.get("https://svod-be.roosterteeth.com/api/v1/seasons/rwby-volume-6/episodes?order=des&per_page=1")
            r = json.loads(r.text)
            last_date = r["data"][0]["attributes"]["sponsor_golive_at"]
            year,month,day = last_date.split("-")
            day = day.split("T")[0]
            if not (int(year) == time_current.tm_year and int(month) == time_current.tm_mon and int(day) == time_current.tm_mday):
                break
        time.sleep(30)

def check_loop():
    global current
    change_game(discord.Game(name="Waiting for new episode..."))
    while True:
        r = requests.get("https://svod-be.roosterteeth.com/api/v1/seasons/rwby-volume-6/episodes?order=des&per_page=1")
        json_data = json.loads(r.text)
        if current != json_data["data"][0]["uuid"]:
            change_game(discord.Game(name="Updating..."))
            send("New Episode detected!")
            send("The title is: '" + json_data["data"][0]["attributes"]["title"] + "'")
            update()
            send("Done!")
            send("@everyone , RWBY Volume 6 Episode "+ str(json_data["data"][0]["attributes"]["number"])+": "+json_data["data"][0]["attributes"]["title"]+" is now avalible at: https://how2trianglemuygud.com/rwbyvol6/")
            current = json_data["data"][0]["uuid"]
            change_game(discord.Game(name="Ep " + str(json_data["data"][0]["attributes"]["number"]) + " - " + json_data["data"][0]["attributes"]["title"]))
            break
        time.sleep(30)

setup()
while True:
    wait()
    check_loop()