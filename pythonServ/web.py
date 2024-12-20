from flask import Flask, request, jsonify
from utils import save_file
from specGenerate import flac_to_spectrogram
from cfg import UPLOAD_FOLDER
import re

app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

@app.route('/upload', methods=['POST'])
def upload_file():
    if 'file' not in request.files:
        return jsonify({'error': 'No file part'}), 400
    
    file = request.files['file']
    if file.filename == '':
        return jsonify({'error': 'No selected file'}), 400
    colormapurl = request.args.get('colormap')
    fileidurl = request.args.get('fileid')
    filename = re.sub(r'[^\w\sа-яА-Я.-]', '', file.filename, flags=re.UNICODE).strip()
    print("Получен аудиофайл:", filename)
    file_path = save_file(file, app.config['UPLOAD_FOLDER'], filename)
    try:
        flac_to_spectrogram(file_path, colormapurl, fileidurl, filename)
        return  jsonify({'message': 'file converted and send successfully'}), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500
