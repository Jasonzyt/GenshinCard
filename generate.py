# -*- coding: utf-8 -*-
import argparse
import hashlib
import math
import os
import random
import time

import requests

from selenium import webdriver
from PIL import Image, ImageDraw

MIYOUSHE_API = 0
HOYOLAB_API = 1
DS_SALT = "xV8v4Qu54lUKrEYFZkJhB8cuOh9Asafs"

REGION_CN_OFFICIAL = "cn_gf01"
REGION_CN_BILIBILI = "cn_qd01"
REGION_OS_NA = "os_usa"
REGION_OS_EU = "os_euro"
REGION_OS_AS = "os_asia"
REGION_OS_SAR = "os_cht"


APIS = [
    {
        "PlayerIndexUrl": "https://api-takumi-record.mihoyo.com/game_record/app/genshin/api/index",
        "PlayerCharacterUrl": "https://api-takumi-record.mihoyo.com/game_record/app/genshin/api/character",
        "PlayerSprialAbyssUrl": "https://api-takumi-record.mihoyo.com/game_record/app/genshin/api/spiralAbyss",
    },
    {
        "PlayerIndexUrl": "https://bbs-api-os.hoyolab.com/game_record/genshin/api/index",
        "PlayerCharacterUrl": "",  # Not supported
        "PlayerSprialAbyssUrl": "https://bbs-api-os.hoyolab.com/game_record/genshin/api/spiralAbyss",
    },
]

SVG_TEMPLATE = """
<svg width="720" height="420" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <style>
        @font-face {{
            font-family: 'HYWenHei-85W';
            src: url('./assets/font/subset-HYWenHei-HEW.woff2') format('woff2'), url('./assets/font/subset-HYWenHei-HEW.woff') format('woff');
            font-weight: 900;
            font-style: normal;
            font-display: swap;
        }}
        @keyframes fadein {{
            to {{
                opacity: 1;
            }}
        }}
        svg {{
            background: rgba(0, 0, 0, 0);
            box-shadow: 10px 16px 8px 0 #bbb;
            border-radius: 40px;
        }}
        * {{
            font-family: 'HYWenHei-85W', sans-serif;
            font-weight: 900;
        }}
        th {{
            font-size: 30px;
        }}
        td {{
            font-size: 48px;
        }}
        tr, td, th {{
            padding: 0;
        }}
    </style>
    <image id="background" xlink:href="{background}" width="100%" height="100%" preserveAspectRatio="xMaxYMid slice"/>
    <text x="60px" y="90px" width="600px" fill="white">
        <tspan font-size="48px">{name}
            <tspan font-size="32px" dx="30px">Lv.<tspan font-size="40px" dx="4px">{level}</tspan>
            </tspan>
        </tspan>
    </text>
    <text x="60px" y="136px" width="600px" fill="white">
        <tspan font-size="32px">UID: {uid}</tspan>
    </text>
    <foreignObject x="40px" y="152px" width="640px" height="200px">
        <table xmlns="http://www.w3.org/1999/xhtml" style="width: 640px; color: white; line-height: 3em; text-align: center;">
            <tr><th>活跃天数</th><th>深境螺旋</th><th>获得角色</th><th>达成成就</th></tr>
            <tr><td>{active_days}</td><td>{spiral_abyss}</td><td>{characters}</td><td>{achievements}</td></tr>
            <tr><th>开启宝箱</th><th>供奉神瞳</th><th>激活锚点</th><th>洞天仙力</th></tr>
            <tr><td>{chests}</td><td>{culus}</td><td>{waypoints}</td><td>{comfort_num}</td></tr>
        </table>
    </foreignObject>
</svg>
"""

global_config = {"cookie": ""}


def get_region_by_uid(uid):
    if len(uid) < 2:
        return None
    c = uid[:1]
    if c == "1" or c == "2":
        return REGION_CN_OFFICIAL
    elif c == "5":
        return REGION_CN_BILIBILI
    elif c == "6":
        return REGION_OS_NA
    elif c == "7":
        return REGION_OS_EU
    elif c == "8":
        return REGION_OS_AS
    elif c == "9":
        return REGION_OS_SAR
    return None


def get_apis_by_region(region):
    if len(region) < 2:
        return None
    if region[:2] == "cn":
        return APIS[MIYOUSHE_API]
    elif region[:2] == "os":
        return APIS[HOYOLAB_API]
    return None


def get_random_background():
    backgrounds = []
    for file in os.listdir("./assets/img"):
        if file.endswith(".jpg") or file.endswith(".png"):
            backgrounds.append(file)
    if len(backgrounds) == 0:
        raise Exception("No backgrounds found")
    return "assets/img/" + random.choice(backgrounds)


def render_svg(fn):
    driver = webdriver.Chrome()
    driver.get(f"file:///{os.getcwd()}/{fn}")
    driver.set_window_size(750, 570)
    time.sleep(2)
    png = driver.get_screenshot_as_png()
    driver.quit()
    with open(fn.replace(".svg", ".png"), "wb") as f:
        f.write(png)
    process_png(fn.replace(".svg", ".png"))


def is_in_square(x, y, sx, sy, l):
    """
    (x, y) is the point to check
    (sx, sy) is the top left corn
    l is the length of the square
    """
    return x >= sx and x <= sx + l and y >= sy and y <= sy + l


def is_in_rect(x, y, sx, sy, w, h):
    """
    (x, y) is the point to check
    (sx, sy) is the top left corn
    w is the width of the rect
    h is the height of the rect
    """
    return x >= sx and x <= sx + w and y >= sy and y <= sy + h


def is_in_circle(x, y, cx, cy, radius):
    """
    (x, y) is the point to check
    (cx, cy) is the center of the circle
    radius is the radius of the circle
    """
    return (x - cx) ** 2 + (y - cy) ** 2 < radius**2


def is_in_rounded_rect(x, y):
    WIDTH = 720
    HEIGHT = 420
    RADIUS = 40
    CIRCLE_CENTERS = [
        (0 + RADIUS, 0 + RADIUS),
        (WIDTH - RADIUS, 0 + RADIUS),
        (0 + RADIUS, HEIGHT - RADIUS),
        (WIDTH - RADIUS, HEIGHT - RADIUS),
    ]
    SQUARE_CORNERS = [
        (0, 0),
        (WIDTH - RADIUS, 0),
        (0, HEIGHT - RADIUS),
        (WIDTH - RADIUS, HEIGHT - RADIUS),
    ]

    if x >= WIDTH or y >= HEIGHT:
        return False

    if x >= RADIUS and x <= WIDTH - RADIUS and y >= RADIUS and y <= HEIGHT - RADIUS:
        return True

    outside_squares = True
    for i in range(0, 4):
        corner = SQUARE_CORNERS[i]
        center = CIRCLE_CENTERS[i]
        if is_in_square(x, y, corner[0], corner[1], RADIUS):
            if is_in_circle(x, y, center[0], center[1], RADIUS):
                return True
            outside_squares = False
    if outside_squares:
        return True

    return False


def process_png(fn):
    img = Image.open(fn)
    img = img.convert("RGBA")
    pixels = img.load()
    for x in range(img.size[0]):
        for y in range(img.size[1]):
            if not is_in_rounded_rect(x, y):
                r, g, b = pixels[x, y][:3]
                # convert white to transparent
                if r > 240 and g > 240 and b > 240:
                    pixels[x, y] = (255, 255, 255, 0)
                # convert grey to transparent gradient
                elif r == g and g == b:
                    pixels[x, y] = (r, g, b, 255 - r)
    img.save("out.png")
    img.close()
    os.remove(fn)


def generate_ds(uid, region):
    query = f"role_id={uid}&server={region}"
    # Current time in seconds
    t = int(time.time())
    r = random.randint(100_000, 200_000)
    if r == 100_000:
        r = 642367
    # Generate DS
    text = f"salt={DS_SALT}&t={t}&r={r}&b=&q={query}"
    sign = hashlib.md5(text.encode("utf-8")).hexdigest()
    return f"{t},{r},{sign}"


def http_get(url, ds):
    headers = {
        "Referer": "https://webstatic.mihoyo.com/",
        "User-Agent": "Mozilla/5.0 (Linux; Android 13; M2101K9C Build/TKQ1.220829.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/108.0.5359.128 Mobile Safari/537.36 miHoYoBBS/2.44.1",
        "X-Requested-With": "com.mihoyo.hyperion",
        "DS": ds,
        "Origin": "https://api-takumi-record.mihoyo.com",
        "Host": "api-takumi-record.mihoyo.com",
        "x-rpc-app_version": "2.44.1",
        "x-rpc-client_type": "5",
        "Cookie": global_config["cookie"],
    }
    return requests.get(url, headers=headers, timeout=10)


def get_player_index(uid, region):
    apis = get_apis_by_region(region)
    if apis is None:
        raise Exception("Failed to get apis")
    ds = generate_ds(uid, region)
    url = apis["PlayerIndexUrl"] + f"?role_id={uid}&server={region}"
    resp = http_get(url, ds)
    if resp.status_code != 200:
        resp.raise_for_status()
    # print(resp.text)
    return resp.json()


def generate_image(uid, region):
    data = get_player_index(uid, region)
    if data is None:
        raise Exception("Failed to get player index")
    if data["retcode"] != 0:
        raise Exception(f"Failed to get player index: {data['message']}")
    data = data["data"]
    if data is None:
        raise Exception("Bad player index data")
    role_data = data["role"]
    if role_data is None:
        raise Exception("Bad player index data")
    stats_data = data["stats"]
    if stats_data is None:
        raise Exception("Bad player index data")
    home_list = data["homes"]
    if home_list is None or type(home_list) != list:
        raise Exception("Bad player index data")
    name = role_data["nickname"]
    level = role_data["level"]
    """
    ActiveDayNumber      int    `json:"active_day_number"`
	AchievementNumber    int    `json:"achievement_number"`
	AnemoculusNumber     int    `json:"anemoculus_number"`
	GeoculusNumber       int    `json:"geoculus_number"`
	AvatarNumber         int    `json:"avatar_number"`
	WayPointNumber       int    `json:"way_point_number"`
	DomainNumber         int    `json:"domain_number"`
	SpiralAbyss          string `json:"spiral_abyss"`
	PreciousChestNumber  int    `json:"precious_chest_number"`
	LuxuriousChestNumber int    `json:"luxurious_chest_number"`
	ExquisiteChestNumber int    `json:"exquisite_chest_number"`
	CommonChestNumber    int    `json:"common_chest_number"`
	ElectroculusNumber   int    `json:"electroculus_number"`
	MagicChestNumber     int    `json:"magic_chest_number"`
	DendroculusNumber    int    `json:"dendroculus_number"`
    """
    active_days = stats_data["active_day_number"]
    spiral_abyss = stats_data["spiral_abyss"]
    achievements = stats_data["achievement_number"]
    characters = stats_data["avatar_number"]
    chests = (
        stats_data["precious_chest_number"]
        + stats_data["luxurious_chest_number"]
        + stats_data["exquisite_chest_number"]
        + stats_data["common_chest_number"]
    )
    culus = (
        stats_data["anemoculus_number"]
        + stats_data["geoculus_number"]
        + stats_data["electroculus_number"]
        + stats_data["dendroculus_number"]
    )
    waypoints = stats_data["way_point_number"]
    comfort_num = 0
    for home in home_list:
        comfort_num = max(comfort_num, home["comfort_num"])
    svg = SVG_TEMPLATE.format(
        background=get_random_background(),
        name=name,
        level=level,
        uid=uid,
        active_days=active_days,
        spiral_abyss=spiral_abyss,
        achievements=achievements,
        characters=characters,
        chests=chests,
        culus=culus,
        waypoints=waypoints,
        comfort_num=comfort_num,
    )
    # print(svg)
    with open("temp.svg", "w", encoding="utf-8") as f:
        f.write(svg)
    render_svg("temp.svg")
    os.remove("temp.svg")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--uid", help="User ID", required=True)
    parser.add_argument("--cookie", help="Cookie", required=True)
    args = parser.parse_args()
    global_config["cookie"] = args.cookie
    region = get_region_by_uid(args.uid)
    if region is None:
        raise Exception("Failed to get region")
    generate_image(args.uid, region)
