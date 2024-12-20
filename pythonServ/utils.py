import matplotlib.colors as mcolors
import numpy as np
import os
import requests
from PIL import Image


def create_custom_cmap(image_path):

    image = Image.open(image_path)
    image = image.resize((256, 1))

    data = np.array(image)

    if data.shape[2] == 4:
        data = data[:, :, :3]

    normalized_colors = data[0] / 255.0

    custom_cmap = mcolors.LinearSegmentedColormap.from_list("custom_map", normalized_colors)
    return custom_cmap

def choose_map(colormap: str):
    if colormap == "inferno":
        return create_custom_cmap("./src/inferno_map.png")
    else:
        return create_custom_cmap("./src/turbo_map.png")


def save_file(file, upload_folder, filename):
    file_path = os.path.join(upload_folder, filename)
    file.save(file_path)
    return file_path

def send_spectrogram(buffer, url):
    try:
        buffer.seek(0)
        headers = {"Content-Type": "application/octet-stream"}
        response = requests.post(url, headers=headers, data=buffer)
        response.raise_for_status()
        print("Спектрограмма построена и отправлена")
    except requests.exceptions.RequestException as e:
        raise RuntimeError(f"Error sending spectrogram: {e}")



