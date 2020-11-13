from datetime import datetime
from subprocess import call 
import sys
import os
import time
import json


def get_current_image():
    with open("image_version", "r") as fr:
        image_version = fr.readline()
    return image_version
    

def is_new_image(prev_version):
    new_version = get_current_image()
    if new_version != prev_version:
        return True

    return False

def spinup(image: str):
    os.system("docker run --log-opt max-size=10m --log-opt max-file=5 --restart always -d " + image + ">> id.txt")
    
    with open("id.txt", "r") as fr:
        container: str = fr.readlines()[-1].replace("\n", "")
    os.system("rm id.txt")
    return container

def start_db(container: str):
    time.sleep(20)
    os.system(f"docker exec -it {container} ash  -c 'touch start'")
    print(f"Switched state to {container}")

def destroy(container: str):
    os.system(f"docker kill {container}")

def run(blue_container: str, green_container: str):
    update_time: bool = False
    current_image: str = get_current_image()
    while True: 
        if blue_container == None:
            blue_container = spinup(current_image)
            start_db(blue_container)

        if update_time:
            green_container = spinup(current_image)
            time.sleep(5)
            destroy(blue_container)
            start_db(green_container)
            blue_container = green_container
            green_container = None
            update_time = False

        if (datetime.now().hour % 4 == 0 and datetime.now().minute == 0 and datetime.now().second == 0):
            update_time = True
            print("Re-establishing connections")

        elif is_new_image(current_image):
            update_time = True
            current_image = get_current_image()
            print("Updating to new image: ", current_image)


if __name__ == '__main__':
    run(None, None)
