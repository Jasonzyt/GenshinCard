import os
import pathlib

files = os.listdir(".")
out = ""

for file in files:
    name = pathlib.Path(file).name
    ext = pathlib.Path(file).suffix
    stem = pathlib.Path(file).stem
    if ext == ".png" or ext == ".jpg":
        parts = stem.split("_p")
        if len(parts) == 2:  # is pixiv
            out += f"[{name}](./{name}) | https://pixiv.net/artworks/{parts[0]}"
            if parts[1] != "0":
                out += f"#{parts[1]}"
            out += "\n"
        else:
            url = stem.replace("-", "/")
            out += f"[{name}](./{name}) | https://{url}\n"

print(out)
