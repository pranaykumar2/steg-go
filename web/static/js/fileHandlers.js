const FileHandlers = (function() {
    const initializedElements = new Set();

    const uploadConfig = {
        'hide-text-file': {
            uploadAreaSelector: '#hide-text-upload-area, .hide-text-upload, .upload-area[data-for="hide-text-file"]',
            previewContainerId: 'hide-text-preview-container',
            capacityId: 'capacity-indicator',
            acceptedTypes: ['image/jpeg', 'image/png'],
            maxSize: 50 * 1024 * 1024
        },
        'hide-file-image': {
            uploadAreaSelector: '#hide-file-upload-area, .hide-file-upload, .upload-area[data-for="hide-file-image"]',
            previewContainerId: 'hide-file-preview-container',
            acceptedTypes: ['image/jpeg', 'image/png'],
            maxSize: 50 * 1024 * 1024
        },
        'hide-file-file': {
            uploadAreaSelector: '#hide-file-file-area, .hide-file-file, .upload-area[data-for="hide-file-file"]',
            previewContainerId: 'hide-file-file-preview-container',
            isFileOnly: true,
            acceptedTypes: ['*/*'],
            maxSize: 50 * 1024 * 1024
        },
        'extract-file': {
            uploadAreaSelector: '#extract-upload-area, .extract-upload, .upload-area[data-for="extract-file"]',
            previewContainerId: 'extract-preview-container',
            acceptedTypes: ['image/jpeg', 'image/png'],
            maxSize: 50 * 1024 * 1024
        },
        'analyze-file': {
            uploadAreaSelector: '#analyze-upload-area, .analyze-upload, .upload-area[data-for="analyze-file"]',
            previewContainerId: 'analyze-preview-container',
            acceptedTypes: ['image/jpeg', 'image/png'],
            maxSize: 50 * 1024 * 1024
        }
    };

    /**
     * Format file size to human-readable string
     * @param {number} bytes - File size in bytes
     * @returns {string} - Formatted file size
     */
    const formatFileSize = function(bytes) {
        if (bytes === 0) return '0 Bytes';

        const units = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(1024));
        return parseFloat((bytes / Math.pow(1024, i)).toFixed(2)) + ' ' + units[i];
    };

    /**
     * Get file icon based on file type or extension
     * @param {File} file - File object
     * @returns {string} - BoxIcon class name
     */
    const getFileIcon = function(file) {
        const type = file.type;
        const extension = file.name.split('.').pop().toLowerCase();

        if (type.includes('image')) {
            return 'bx-image';
        } else if (type.includes('pdf') || extension === 'pdf') {
            return 'bx-file-pdf';
        } else if (type.includes('word') || extension === 'doc' || extension === 'docx') {
            return 'bx-file-doc';
        } else if (type.includes('excel') || extension === 'xls' || extension === 'xlsx') {
            return 'bx-spreadsheet';
        } else if (type.includes('video')) {
            return 'bx-video';
        } else if (type.includes('audio')) {
            return 'bx-music';
        } else if (type.includes('zip') || type.includes('compressed') ||
            ['zip', 'rar', '7z', 'tar', 'gz'].includes(extension)) {
            return 'bx-archive';
        } else if (type.includes('text') || extension === 'txt') {
            return 'bx-file-txt';
        } else {
            return 'bx-file-blank';
        }
    };

    /**
     * Validate file (type and size)
     * @param {File} file - File to validate
     * @param {Object} config - Configuration for validation
     * @returns {boolean} - Whether the file is valid
     */
    const validateFile = function(file, config) {
        if (file.size > config.maxSize) {
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast(
                    `The file exceeds the maximum size of ${formatFileSize(config.maxSize)}.`,
                    'error',
                    'File Too Large'
                );
            } else {
                console.error(`File too large: ${formatFileSize(file.size)} exceeds maximum ${formatFileSize(config.maxSize)}`);
                alert(`File too large: ${formatFileSize(file.size)} exceeds maximum ${formatFileSize(config.maxSize)}`);
            }
            return false;
        }

        if (config.acceptedTypes[0] !== '*/*') {
            if (!config.acceptedTypes.includes(file.type)) {
                if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                    UIManager.showToast(
                        `File type ${file.type} not supported. Please upload ${config.acceptedTypes.join(' or ')}.`,
                        'error',
                        'Invalid File Type'
                    );
                } else {
                    console.error(`Invalid file type: ${file.type}`);
                    alert(`Invalid file type: ${file.type}. Please upload ${config.acceptedTypes.join(' or ')}.`);
                }
                return false;
            }
        }

        return true;
    };

    /**
     * Calculate image capacity - uses client-side calculation if API fails
     * @param {File} file - Image file
     * @param {string} capacityId - ID of capacity element
     */
    const calculateImageCapacity = function(file, capacityId) {
        console.log('Calculating capacity for image', file.name);

        const capacityIndicator = document.getElementById(capacityId);
        const capacityFill = document.getElementById('capacity-fill');
        const capacityText = document.getElementById('capacity-text');

        if (!capacityIndicator || !capacityFill || !capacityText) {
            console.warn('Capacity indicator elements not found');
            return;
        }

        capacityIndicator.style.display = 'block';
        capacityFill.style.width = '0%';
        capacityText.textContent = 'Calculating...';

        // Try API capacity calculation first
        if (typeof StegGoAPI !== 'undefined' && StegGoAPI.calculateCapacity) {
            calculateClientSideCapacity(file, function(clientCapacity) {
                updateCapacityUI(clientCapacity);
                StegGoAPI.calculateCapacity(file)
                    .then(function(apiCapacity) {
                        console.log('API capacity calculation succeeded:', apiCapacity);
                        if (apiCapacity && apiCapacity.maxBytes) {
                            updateCapacityUI(apiCapacity);
                        }
                    })
                    .catch(function(error) {
                        console.log('Using client-side capacity calculation due to API error:', error);
                    });
            });
        } else {
            calculateClientSideCapacity(file, updateCapacityUI);
        }

        setupCapacityListeners(file, capacityId);
    };

    /**
     * Calculate capacity purely on client side (no API)
     * @param {File} file - Image file
     * @param {Function} callback - Callback with capacity result
     */
    const calculateClientSideCapacity = function(file, callback) {
        const reader = new FileReader();

        reader.onload = function(e) {
            const img = new Image();

            img.onload = function() {
                const width = this.width;
                const height = this.height;
                const pixelCount = width * height;
                const capacityBits = pixelCount * 3;
                const capacityBytes = Math.floor(capacityBits / 8);
                const usableCapacityBytes = Math.floor(capacityBytes * 0.85);

                const capacity = {
                    maxBytes: usableCapacityBytes,
                    maxBytesFormatted: formatFileSize(usableCapacityBytes),
                    percentageUsed: 0,
                    lsbBits: 3,
                    recommendation: `This image can store approximately ${formatFileSize(usableCapacityBytes)} of data safely.`
                };
                if (callback) callback(capacity);
            };

            img.onerror = function() {
                console.error('Error loading image for capacity calculation');
                const fallbackCapacity = {
                    maxBytes: 10240,
                    maxBytesFormatted: '10 KB',
                    percentageUsed: 0,
                    lsbBits: 3,
                    recommendation: 'Unable to calculate exact capacity. Using a conservative estimate.'
                };

                if (callback) callback(fallbackCapacity);
            };

            img.src = e.target.result;
        };

        reader.onerror = function() {
            console.error('Error reading file for capacity calculation');
            const fallbackCapacity = {
                maxBytes: 10240,
                maxBytesFormatted: '10 KB',
                percentageUsed: 0,
                lsbBits: 3,
                recommendation: 'Unable to calculate capacity. Using a conservative estimate.'
            };

            if (callback) callback(fallbackCapacity);
        };

        reader.readAsDataURL(file);
    };

    /**
     * Update the capacity UI with new capacity data
     * @param {Object} capacity - Capacity data
     */
    const updateCapacityUI = function(capacity) {
        const capacityFill = document.getElementById('capacity-fill');
        const capacityText = document.getElementById('capacity-text');

        if (!capacityFill || !capacityText) return;

        let currentSize = 0;
        const messageInput = document.getElementById('hide-text-message');
        const fileInput = document.getElementById('hide-file-file');

        if (messageInput && messageInput.value) {
            currentSize = new Blob([messageInput.value]).size;
        } else if (fileInput && fileInput.files && fileInput.files[0]) {
            currentSize = fileInput.files[0].size;
        }

        const percentage = capacity.maxBytes > 0 ?
            Math.min(100, Math.round((currentSize / capacity.maxBytes) * 100)) : 0;

        capacityFill.style.width = `${percentage}%`;
        capacityText.textContent = `${formatFileSize(currentSize)} / ${capacity.maxBytesFormatted} (${percentage}%)`;

        if (percentage > 90) {
            capacityFill.style.backgroundColor = '#ef4444'; // Red
        } else if (percentage > 70) {
            capacityFill.style.backgroundColor = '#f59e0b'; // Amber
        } else {
            capacityFill.style.backgroundColor = ''; // Default gradient
        }
    };

    /**
     * Set up listeners for capacity updates when message or file changes
     * @param {File} imageFile - The cover image file
     * @param {string} capacityId - ID of capacity element
     */
    const setupCapacityListeners = function(imageFile, capacityId) {

        let baseCapacity = null;

        calculateClientSideCapacity(imageFile, function(capacity) {
            baseCapacity = capacity;

            const messageInput = document.getElementById('hide-text-message');
            if (messageInput) {
                const newMessageInput = messageInput.cloneNode(true);
                messageInput.parentNode.replaceChild(newMessageInput, messageInput);

                newMessageInput.addEventListener('input', function() {
                    if (!baseCapacity) return;

                    const messageSize = new Blob([this.value]).size;
                    const percentage = Math.min(100, Math.round((messageSize / baseCapacity.maxBytes) * 100));

                    const capacityFill = document.getElementById('capacity-fill');
                    const capacityText = document.getElementById('capacity-text');

                    if (capacityFill) capacityFill.style.width = `${percentage}%`;
                    if (capacityText) capacityText.textContent = `${formatFileSize(messageSize)} / ${baseCapacity.maxBytesFormatted} (${percentage}%)`;

                    if (capacityFill) {
                        if (percentage > 90) {
                            capacityFill.style.backgroundColor = '#ef4444';
                        } else if (percentage > 70) {
                            capacityFill.style.backgroundColor = '#f59e0b';
                        } else {
                            capacityFill.style.backgroundColor = '';
                        }
                    }
                });
            }

            const fileInput = document.getElementById('hide-file-file');
            if (fileInput) {
                const newFileInput = fileInput.cloneNode(true);
                fileInput.parentNode.replaceChild(newFileInput, fileInput);

                newFileInput.addEventListener('change', function() {
                    if (!baseCapacity || !this.files || !this.files[0]) return;

                    const fileSize = this.files[0].size;
                    const percentage = Math.min(100, Math.round((fileSize / baseCapacity.maxBytes) * 100));

                    const capacityFill = document.getElementById('capacity-fill');
                    const capacityText = document.getElementById('capacity-text');

                    if (capacityFill) capacityFill.style.width = `${percentage}%`;
                    if (capacityText) capacityText.textContent = `${formatFileSize(fileSize)} / ${baseCapacity.maxBytesFormatted} (${percentage}%)`;

                    if (capacityFill) {
                        if (percentage > 90) {
                            capacityFill.style.backgroundColor = '#ef4444'; // Red
                        } else if (percentage > 70) {
                            capacityFill.style.backgroundColor = '#f59e0b'; // Amber
                        } else {
                            capacityFill.style.backgroundColor = ''; // Default gradient
                        }
                    }
                });
            }
        });
    };

    /**
     * Create image preview
     * @param {File} file - Image file
     * @param {string} previewContainerId - ID of preview container
     * @param {string} inputId - ID of file input
     */
    const createImagePreview = function(file, previewContainerId, inputId) {
        console.log('Creating image preview for', file.name, 'in', previewContainerId);

        const previewContainer = document.getElementById(previewContainerId);
        if (!previewContainer) {
            console.error('Preview container not found:', previewContainerId);
            return;
        }

        const reader = new FileReader();

        reader.onload = function(e) {
            previewContainer.innerHTML = `
                <img src="${e.target.result}" alt="Preview" class="preview-image">
                <div class="preview-info">
                    <div class="file-info">
                        <div class="file-name">${file.name}</div>
                        <div class="file-size">${formatFileSize(file.size)}</div>
                    </div>
                    <button type="button" class="remove-file" aria-label="Remove file">
                        <i class='bx bx-trash'></i>
                    </button>
                </div>
            `;
            previewContainer.style.display = 'block';

            // Add event listener to remove button
            const removeBtn = previewContainer.querySelector('.remove-file');
            if (removeBtn) {
                removeBtn.addEventListener('click', function() {
                    const fileInput = document.getElementById(inputId);
                    if (fileInput) {
                        fileInput.value = '';
                    }
                    previewContainer.innerHTML = '';
                    previewContainer.style.display = 'none';
                    const config = uploadConfig[inputId];
                    if (config) {
                        const uploadArea = document.querySelector(config.uploadAreaSelector);
                        if (uploadArea) {
                            uploadArea.style.display = 'flex';
                        }
                    }
                    const capacityIndicator = document.getElementById('capacity-indicator');
                    if (capacityIndicator) {
                        capacityIndicator.style.display = 'none';
                    }
                    updateFormStatus(document.getElementById(inputId).closest('form'));
                });
            }
            const config = uploadConfig[inputId];
            if (config && config.capacityId) {
                calculateImageCapacity(file, config.capacityId);
            }

            updateFormStatus(document.getElementById(inputId).closest('form'));
        };

        reader.onerror = function() {
            console.error('Error reading file for preview');
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('Error creating preview for ' + file.name, 'error');
            }
        };

        reader.readAsDataURL(file);
    };

    /**
     * Create file preview
     * @param {File} file - File object
     * @param {string} previewContainerId - ID of preview container
     * @param {string} inputId - ID of file input
     */
    const createFilePreview = function(file, previewContainerId, inputId) {
        console.log('Creating file preview for', file.name, 'in', previewContainerId);

        const previewContainer = document.getElementById(previewContainerId);
        if (!previewContainer) {
            console.error('Preview container not found:', previewContainerId);
            return;
        }

        previewContainer.innerHTML = `
            <div class="file-preview">
                <div class="file-icon">
                    <i class='bx ${getFileIcon(file)}'></i>
                </div>
                <div class="preview-info">
                    <div class="file-info">
                        <div class="file-name">${file.name}</div>
                        <div class="file-size">${formatFileSize(file.size)}</div>
                    </div>
                    <button type="button" class="remove-file" aria-label="Remove file">
                        <i class='bx bx-trash'></i>
                    </button>
                </div>
            </div>
        `;

        previewContainer.style.display = 'block';

        const removeBtn = previewContainer.querySelector('.remove-file');
        if (removeBtn) {
            removeBtn.addEventListener('click', function() {
                // Reset file input
                const fileInput = document.getElementById(inputId);
                if (fileInput) {
                    fileInput.value = '';
                }

                previewContainer.innerHTML = '';
                previewContainer.style.display = 'none';
                const config = uploadConfig[inputId];
                if (config) {
                    const uploadArea = document.querySelector(config.uploadAreaSelector);
                    if (uploadArea) {
                        uploadArea.style.display = 'flex';
                    }
                }
                updateFormStatus(document.getElementById(inputId).closest('form'));
            });
        }
        updateFormStatus(document.getElementById(inputId).closest('form'));
    };

    /**
     * Handle file selection
     * @param {File} file - Selected file
     * @param {string} inputId - ID of file input
     */
    const handleFileSelection = function(file, inputId) {
        console.log('Handling file selection:', file.name, 'for input', inputId);

        if (!file || !inputId) {
            console.error('Invalid file or input ID');
            return;
        }

        const config = uploadConfig[inputId];
        if (!config) {
            console.error('No configuration found for input:', inputId);
            return;
        }
        if (!validateFile(file, config)) {
            return;
        }

        const uploadArea = document.querySelector(config.uploadAreaSelector);
        if (uploadArea) {
            uploadArea.style.display = 'none';
        } else {
            console.warn('Upload area not found for selector:', config.uploadAreaSelector);
        }

        if (file.type.startsWith('image') && !config.isFileOnly) {
            createImagePreview(file, config.previewContainerId, inputId);
        } else {
            createFilePreview(file, config.previewContainerId, inputId);
        }
    };

    /**
     * Update form status (enable/disable submit button)
     * @param {HTMLFormElement} form - Form to update
     */
    const updateFormStatus = function(form) {
        if (!form) return;

        const submitBtn = form.querySelector('button[type="submit"]');
        if (!submitBtn) return;

        switch (form.id) {
            case 'hide-text-form':
                const textImage = document.getElementById('hide-text-file').files[0];
                const message = document.getElementById('hide-text-message').value.trim();
                submitBtn.disabled = !(textImage && message);
                break;

            case 'hide-file-form':
                const fileImage = document.getElementById('hide-file-image').files[0];
                const fileToHide = document.getElementById('hide-file-file').files[0];
                submitBtn.disabled = !(fileImage && fileToHide);
                break;

            case 'extract-form':
                const extractImage = document.getElementById('extract-file').files[0];
                const key = document.getElementById('extract-key').value.trim();
                submitBtn.disabled = !(extractImage && key);
                break;

            case 'analyze-form':
                const analyzeImage = document.getElementById('analyze-file').files[0];
                submitBtn.disabled = !analyzeImage;
                break;
        }
    };

    /**
     * Set up drag and drop functionality
     * @param {HTMLElement} uploadArea - Upload area element
     * @param {HTMLElement} fileInput - File input element
     */
    const setupDragAndDrop = function(uploadArea, fileInput) {
        if (!uploadArea || !fileInput) return;

        if (initializedElements.has(uploadArea)) return;
        initializedElements.add(uploadArea);
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            uploadArea.addEventListener(eventName, function(e) {
                e.preventDefault();
                e.stopPropagation();
            }, false);
        });
        uploadArea.addEventListener('dragenter', function() {
            this.classList.add('dragover');
        }, false);

        uploadArea.addEventListener('dragover', function() {
            this.classList.add('dragover');
        }, false);

        uploadArea.addEventListener('dragleave', function() {
            this.classList.remove('dragover');
        }, false);

        uploadArea.addEventListener('drop', function() {
            this.classList.remove('dragover');
        }, false);

        uploadArea.addEventListener('drop', function(e) {
            if (e.dataTransfer.files && e.dataTransfer.files[0]) {
                const dt = new DataTransfer();
                dt.items.add(e.dataTransfer.files[0]);
                fileInput.files = dt.files;

                const event = new Event('change');
                fileInput.dispatchEvent(event);
            }
        }, false);
    };

    /**
     * Initialize file handlers
     */
    const init = function() {
        if (window.fileHandlersInitialized) {
            console.log('FileHandlers already initialized');
            return;
        }

        console.log('Initializing FileHandlers...');
        window.fileHandlersInitialized = true;
        Object.keys(uploadConfig).forEach(inputId => {
            const fileInput = document.getElementById(inputId);
            if (!fileInput) {
                console.warn('File input not found:', inputId);
                return;
            }

            console.log('Setting up file input:', inputId);
            const newInput = fileInput.cloneNode(true);
            fileInput.parentNode.replaceChild(newInput, fileInput);
            newInput.addEventListener('change', function(e) {
                if (this.files && this.files[0]) {
                    handleFileSelection(this.files[0], this.id);
                }
            });
            const config = uploadConfig[inputId];
            const uploadAreas = document.querySelectorAll(config.uploadAreaSelector);

            if (uploadAreas.length === 0) {
                console.warn('Upload area not found for:', inputId);
            }

            uploadAreas.forEach(uploadArea => {
                const newArea = uploadArea.cloneNode(true);
                uploadArea.parentNode.replaceChild(newArea, uploadArea);
                newArea.addEventListener('click', function(e) {
                    if (e.target.tagName === 'BUTTON' || e.target.closest('button')) {
                        return;
                    }
                    newInput.click();
                });

                setupDragAndDrop(newArea, newInput);
            });
        });

        document.querySelectorAll('form').forEach(form => {
            form.querySelectorAll('input[type="text"], textarea').forEach(input => {
                const newInput = input.cloneNode(true);
                input.parentNode.replaceChild(newInput, input);
                newInput.addEventListener('input', function() {
                    updateFormStatus(form);
                });
            });

            updateFormStatus(form);
        });

        console.log('FileHandlers initialization complete');
    };

    return {
        init,
        formatFileSize,
        getFileIcon,
        handleFileSelection
    };
})();

document.addEventListener('DOMContentLoaded', function() {
    setTimeout(function() {
        if (!window.fileHandlersInitialized && FileHandlers && typeof FileHandlers.init === 'function') {
            console.log('Initializing FileHandlers from DOMContentLoaded');
            FileHandlers.init();
        }
    }, 100);
});
