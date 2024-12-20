// Обновляем текст метки кнопки выбора файла при выборе файла
document.getElementById('file-input').addEventListener('change', (event) => {
    const fileInput = event.target;
    const fileName = fileInput.files[0]?.name || "Выберите файл";
    const fileLabel = document.querySelector('.custom-file-label');
    fileLabel.textContent = fileName;
});

// Обработчик кнопки загрузки
document.getElementById('uploadButton').addEventListener('click', async () => {
    const fileInput = document.getElementById('file-input');
    const colormapSelect = document.getElementById('colormap-select');
    const loadingIndicator = document.getElementById('loading-indicator');
    const file = fileInput.files[0];

    if (!file) {
        document.getElementById('status').textContent = "Выберите файл для загрузки!";
        return;
    }

    loadingIndicator.style.display = 'block';

    const colormap = colormapSelect.value;
    if (!colormap) {
        document.getElementById('status').textContent = "Выберите цветовую палитру!";
        loadingIndicator.style.display = 'none'; //
        return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
        const response = await fetch(`http://localhost:8080/uploadFlac?colormap=${encodeURIComponent(colormap)}`, {
            method: 'POST',
            body: formData
        });

        if (response.ok) {
            const data = await response.json();
            const fileID = data.fileID;
            displaySpectrogram(fileID, colormap);
        } else {
            const errorText = await response.text();
            document.getElementById('status').textContent = "Ошибка загрузки файла: " + errorText;
        }
    } catch (error) {
        document.getElementById('status').textContent = "Ошибка: " + error.message;
    } finally {
        loadingIndicator.style.display = 'none';
    }
});

// Логика вывода спектрограммы
async function displaySpectrogram(fileID, colormap) {
    // Скрываем контейнер загрузки
    const uploadContainer = document.getElementById('upload-container');
    uploadContainer.style.display = 'none';

    // Отображаем контейнер со спектрограммой
    const spectrogramContainer = document.getElementById('spectrogram-container');
    spectrogramContainer.style.display = 'flex';

    // Устанавливаем изображение спектрограммы
    const img = document.getElementById('spectrogram-image');
    img.src = `http://localhost:8080/getSpec?fileid=${encodeURIComponent(fileID)}`;
}

// Обработчик для кнопки "Назад"
document.getElementById('backButton').addEventListener('click', () => {
    const uploadContainer = document.getElementById('upload-container');
    const spectrogramContainer = document.getElementById('spectrogram-container');
    const fileInput = document.getElementById('file-input');
    const fileLabel = document.querySelector('.custom-file-label');

    spectrogramContainer.style.display = 'none';
    uploadContainer.style.display = 'flex';

    fileInput.value = '';
    fileLabel.textContent = 'Выберите файл';

    document.getElementById('status').textContent = '';
});