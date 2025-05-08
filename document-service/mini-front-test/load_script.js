// script.js

const uploadForm = document.getElementById('uploadForm');
const userIdInput = document.getElementById('userId');
const fileInput = document.getElementById('fileInput');
const selectedFilesContainer = document.getElementById('selectedFilesContainer'); // Nuevo contenedor
const uploadStatusDiv = document.getElementById('uploadStatus');

// Almacenar los objetos File seleccionados para acceder a ellos después
let selectedFiles = [];

// --- Evento cuando se seleccionan archivos ---
fileInput.addEventListener('change', function() {
    selectedFiles = Array.from(fileInput.files); // Convertir FileList a Array
    displaySelectedFiles(selectedFiles); // Mostrar los campos dinámicamente
});

// --- Función para mostrar los archivos seleccionados y añadir campos de tipo ---
function displaySelectedFiles(files) {
    selectedFilesContainer.innerHTML = ''; // Limpiar contenedor anterior

    if (files.length === 0) {
        return; // No hay archivos seleccionados
    }

    const fileListTitle = document.createElement('h5');
    fileListTitle.textContent = 'Archivos Seleccionados:';
    selectedFilesContainer.appendChild(fileListTitle);

    files.forEach((file, index) => {
        const fileItemDiv = document.createElement('div');
        fileItemDiv.classList.add('file-info-item'); // Clase para estilizar

        // Nombre del archivo
        const fileNamePara = document.createElement('p');
        fileNamePara.innerHTML = `<strong>Archivo ${index + 1}:</strong> ${file.name}`;
        fileItemDiv.appendChild(fileNamePara);

        // Campo para el tipo de documento
        const typeLabel = document.createElement('label');
        typeLabel.setAttribute('for', `documentType_${index}`);
        typeLabel.classList.add('form-label');
        typeLabel.textContent = 'Tipo de Documento:';
        fileItemDiv.appendChild(typeLabel);

        const typeInput = document.createElement('input');
        typeInput.setAttribute('type', 'text');
        typeInput.setAttribute('id', `documentType_${index}`);
        typeInput.setAttribute('name', `documentType_${index}`); // Nombre para identificarlo
        typeInput.classList.add('form-control');
        typeInput.setAttribute('required', true); // Hacer el campo obligatorio
        typeInput.setAttribute('placeholder', 'Ej: Cédula, Recibo, Factura'); // Placeholder útil

        // Añadir un atributo de datos para asociar el input al índice del archivo original
        typeInput.dataset.fileIndex = index;

        fileItemDiv.appendChild(typeInput);

        selectedFilesContainer.appendChild(fileItemDiv);
    });
}


// --- Evento al enviar el formulario ---
uploadForm.addEventListener('submit', async function(event) {
    event.preventDefault(); // Evita que el formulario se envíe de la manera tradicional

    const userId = userIdInput.value.trim();

    // Validaciones básicas
    if (!userId) {
        alert('Por favor, ingresa un ID de usuario.');
        return;
    }
    if (selectedFiles.length === 0) {
        alert('Por favor, selecciona al menos un archivo.');
        return;
    }

    // Validar que todos los campos de tipo de documento estén llenos
    const documentTypeInputs = selectedFilesContainer.querySelectorAll('input[type="text"]');
    let allTypesFilled = true;
    documentTypeInputs.forEach(input => {
        if (input.value.trim() === '') {
            allTypesFilled = false;
            input.classList.add('is-invalid'); // Marcar el campo inválido con Bootstrap
        } else {
            input.classList.remove('is-invalid');
        }
    });

    if (!allTypesFilled) {
        alert('Por favor, especifica el tipo de documento para cada archivo.');
        return;
    }


    uploadStatusDiv.innerHTML = '<p>Preparando subida...</p>';

    // 1. Preparar el cuerpo de la solicitud para tu API backend
    const filesInfo = [];
    documentTypeInputs.forEach(input => {
        const index = parseInt(input.dataset.fileIndex); // Obtener el índice del archivo original
        const file = selectedFiles[index]; // Obtener el objeto File original

        filesInfo.push({
            fileName: file.name, // Nombre del archivo local
            contentType: file.type || 'application/octet-stream', // Tipo MIME
            size: file.size, // Tamaño del archivo
            documentType: input.value.trim(), // <-- ¡Obtener el tipo de documento ingresado!
            // Si 'documentType' es parte de 'importantMetadata' en tu backend,
            // tendrías que estructurarlo así: importantMetadata: { documentType: input.value.trim() }
            // Por ahora, lo enviamos como un campo separado 'documentType' en FileUploadInfo
        });
    });


    const requestBody = {
        userId: parseInt(userId),
        files: filesInfo,
    };

    // 2. Llamar a tu API backend para obtener las URLs firmadas
    const backendApiUrl = 'http://localhost:8001/files/upload'; // <-- ¡REEMPLAZA con la URL REAL de tu endpoint API!
    // const backendApiUrl = 'http://localhost:8080/files/request-upload-urls'; // Ejemplo local

    try {
        uploadStatusDiv.innerHTML += '<p>Solicitando URLs firmadas a la API...</p>';
        console.log(JSON.stringify(requestBody));
        const response = await fetch(backendApiUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody),
        });

        // Manejar la respuesta de la API
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Error de API: ${response.status} - ${errorText}`);
        }

        const responseData = await response.json();
        const documents = responseData.documents;
        const signedUrls = documents.map(document => document.signedUrl);

        if (!signedUrls || signedUrls.length === 0) {
            uploadStatusDiv.innerHTML += '<p class="text-warning">La API no devolvió URLs firmadas válidas.</p>';
            return;
        }
        // Opcional: Verificar que el número de URLs firmadas coincida con el número de archivos solicitados
        if (signedUrls.length !== selectedFiles.length) {
            uploadStatusDiv.innerHTML += `<p class="text-warning">Advertencia: El número de URLs firmadas (${signedUrls.length}) no coincide con el número de archivos seleccionados (${selectedFiles.length}).</p>`;
            // Decide si quieres continuar o detenerte aquí
        }


        uploadStatusDiv.innerHTML += `<p class="text-success">URLs firmadas recibidas (${signedUrls.length}). Iniciando subida directa...</p>`;

        // 3. Usar las URLs firmadas para subir cada archivo directamente a Cloud Storage
        // Es crucial que el orden de las URLs firmadas coincida con el orden de los archivos en selectedFiles
        // o que la respuesta de la API incluya alguna forma de mapear URL a archivo (ej: por FileName)
        const uploadPromises = documents.map(async (documentInfo, index) => {
            // Asumimos que el orden de urlInfo en la respuesta coincide con el orden de selectedFile;
            console.log("acá!!!: ", documentInfo)
            if (documentInfo.signedURL == undefined) {
                console.log("mal!!!: ", documentInfo.signedUrl)
            }
            const file = selectedFiles[index];
            if (!file) {
                uploadStatusDiv.innerHTML += `<p class="text-danger">Error: No se encontró el archivo local para la URL firmada de ${urlInfo.FileName}.</p>`;
                return; // Saltar este archivo si no se encuentra localmente
            }

            uploadStatusDiv.innerHTML += `<p>Subiendo "${file.name}" a Cloud Storage...</p>`;

            try {
                const uploadResponse = await fetch(documentInfo.signedUrl, { // Usar SignedUrl del objeto urlInfo
                    method: 'PUT', // *** DEBE ser PUT para subir ***
                    headers: {
                        'Content-Type': documentInfo.ContentType,
                    },
                    body: file, // El objeto File seleccionado por el usuario
                });

                if (uploadResponse.ok) {
                    uploadStatusDiv.innerHTML += `<p class="text-success">"${file.name}" subido exitosamente.</p>`;
                } else {
                    const errorText = await uploadResponse.text();
                    uploadStatusDiv.innerHTML += `<p class="text-danger">Error al subir "${file.name}": ${uploadResponse.status} - ${errorText}</p>`;
                }
            } catch (uploadError) {
                uploadStatusDiv.innerHTML += `<p class="text-danger">Error de red al subir "${file.name}": ${uploadError.message}</p>`;
            }
        });

        // Esperar a que todas las subidas se completen (o fallen)
        await Promise.all(uploadPromises);

        uploadStatusDiv.innerHTML += '<p class="text-info">Proceso de subida completado.</p>';

    } catch (error) {
        logError('Error general en el proceso de subida:', error);
        uploadStatusDiv.innerHTML = `<p class="text-danger">Ocurrió un error: ${error.message}</p>`;
    }
});

// Función de ayuda para loguear errores en el cliente (opcional)
function logError(message, error) {
    console.error(message, error);
    // Puedes enviar este error a un servicio de loggin de frontend si tienes uno
}
