import discord
import asyncio
import requests
import json
import threading
import time

f = open("tokens.txt")
lines = f.readlines()
token = lines[0].rstrip()
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
time.sleep(10)
bot_state = "TEST"