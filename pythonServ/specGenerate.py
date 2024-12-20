import librosa
import io
import os
import matplotlib.pyplot as plt
import numpy as np

from matplotlib.ticker import FuncFormatter
from utils import send_spectrogram
from cfg import TICKS, TICK_labels
from utils import choose_map

def flac_to_spectrogram(input_file: str, colormap: str, fileid: str, filename: str):
    try:
        audio_data, sample_rate = librosa.load(input_file, sr=None)

        stft = librosa.stft(audio_data)
        db_scale = librosa.amplitude_to_db(np.abs(stft), ref=np.max)

        plt.figure(figsize=(8, 4))
        plt.imshow(
            db_scale,
            aspect='auto',
            origin='lower',
            cmap= choose_map(colormap),
            vmin=-85,
            extent=(0, len(audio_data) / sample_rate, 0, 22050)
        )
        cbar=plt.colorbar(format='%+2.0f dB')
        cbar.ax.tick_params(labelsize=8)
        dbticks = np.array([-85, -80, -70, -60, -50, -40, -30, -20, -10, 0])
        cbar.set_ticks(dbticks)
        cbar.set_ticklabels([f"{int(tick)} dB" for tick in dbticks])

        plt.ylim(0, 22050)
        plt.title(filename, fontsize=12)
        plt.xlabel('Time (mm:ss)', fontsize=12)
        plt.ylabel('Frequency (Hz)', fontsize=12)
        plt.yticks(TICKS, TICK_labels)

        def time_formatter(x, pos):
            minutes = int(x // 60)
            seconds = int(x % 60)
            return f"{minutes}:{seconds:02d}"

        ax = plt.gca()
        ax.xaxis.set_major_formatter(FuncFormatter(time_formatter))

        max_time = len(audio_data) / sample_rate
        num_ticks = 5
        ticks = np.linspace(0, max_time, num_ticks)
        ticks = np.append(ticks, max_time)
        ax.set_xticks(ticks)

        plt.tight_layout()

        buffer = io.BytesIO()
        plt.savefig(buffer, format="png", dpi=300)
        plt.close()

        send_spectrogram(buffer, "http://localhost:8080/uploadSpec?fileid=" + fileid)
    except Exception as e:
        print(f"Произошла ошибка: {e}")
    finally:
            if os.path.exists(input_file):
                os.remove(input_file)
