document.addEventListener('DOMContentLoaded', function() {
    initTabs();
    initThemeToggle();
    initFormHandlers();
    initCloseButtons();
});

function initTabs() {
    const tabButtons = document.querySelectorAll('.tab-btn');
    if (!tabButtons.length) {
        console.warn('No tab buttons found');
        return;
    }

    tabButtons.forEach(button => {
        button.addEventListener('click', function() {
            const tabId = this.getAttribute('data-tab');
            if (!tabId) return;

            document.querySelectorAll('.tab-btn').forEach(btn => {
                btn.classList.remove('active');
                btn.setAttribute('aria-selected', 'false');
            });

            document.querySelectorAll('.panel').forEach(panel => {
                panel.classList.remove('active');
            });

            this.classList.add('active');
            this.setAttribute('aria-selected', 'true');
            const panel = document.getElementById(`${tabId}-panel`);
            if (panel) {
                panel.classList.add('active');
            } else {
                console.warn(`Panel not found: ${tabId}-panel`);
            }
        });
    });
}

const initializedUIComponents = new Set();

function initThemeToggle() {
    if (initializedUIComponents.has('themeToggle')) {
        console.log('Theme toggle already initialized');
        return;
    }

    const themeToggle = document.getElementById('theme-toggle');
    if (!themeToggle) {
        console.warn('Theme toggle button not found');
        return;
    }

    initializedUIComponents.add('themeToggle');
    console.log('Initializing theme toggle');

    const htmlElement = document.documentElement;
    const savedTheme = localStorage.getItem('theme');
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

    if (savedTheme) {
        htmlElement.setAttribute('data-theme', savedTheme);
    } else {
        const initialTheme = systemPrefersDark ? 'dark' : 'light';
        htmlElement.setAttribute('data-theme', initialTheme);
        localStorage.setItem('theme', initialTheme);
    }

    updateThemeToggleIcon(themeToggle, htmlElement.getAttribute('data-theme'));

    const newThemeToggle = themeToggle.cloneNode(true);
    themeToggle.parentNode.replaceChild(newThemeToggle, themeToggle);

    newThemeToggle.addEventListener('click', function() {
        const currentTheme = htmlElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';

        htmlElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);

        updateThemeToggleIcon(this, newTheme);

        console.log(`Switched to ${newTheme} theme`);
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(`Switched to ${newTheme} theme`, 'info', 'Theme Changed', 2000);
        }
    });
}

/**
 * Update theme toggle icon based on current theme
 * @param {HTMLElement} toggleButton - The theme toggle button
 * @param {string} theme - Current theme ('dark' or 'light')
 */
function updateThemeToggleIcon(toggleButton, theme) {
    if (!toggleButton) return;

    const iconElement = toggleButton.querySelector('i') || toggleButton;

    if (theme === 'dark') {
        iconElement.className = 'bx bx-sun';
        toggleButton.setAttribute('aria-label', 'Switch to light theme');
        toggleButton.setAttribute('title', 'Switch to light theme');
    } else {
        iconElement.className = 'bx bx-moon';
        toggleButton.setAttribute('aria-label', 'Switch to dark theme');
        toggleButton.setAttribute('title', 'Switch to dark theme');
    }
}

function initFormHandlers() {
    const hideTextForm = document.getElementById('hide-text-form');
    if (hideTextForm) {
        hideTextForm.addEventListener('submit', handleHideTextSubmit);
    } else {
        console.warn('Hide text form not found');
    }

    const hideFileForm = document.getElementById('hide-file-form');
    if (hideFileForm) {
        hideFileForm.addEventListener('submit', handleHideFileSubmit);
    } else {
        console.warn('Hide file form not found');
    }

    const extractForm = document.getElementById('extract-form');
    if (extractForm) {
        extractForm.addEventListener('submit', handleExtractSubmit);
    } else {
        console.warn('Extract form not found');
    }

    const analyzeForm = document.getElementById('analyze-form');
    if (analyzeForm) {
        analyzeForm.addEventListener('submit', handleAnalyzeSubmit);
    } else {
        console.warn('Analyze form not found');
    }
}

function initFileUploads() {
    function setupFileUpload(containerId, fileInputId, previewContainerId, acceptedTypes) {
        const fileInput = document.getElementById(fileInputId);
        if (!fileInput) {
            console.warn(`Missing elements for file upload setup: ${containerId}`);
            return;
        }

        const previewContainer = document.getElementById(previewContainerId);
        if (!previewContainer) {
            console.warn(`Preview container not found: ${previewContainerId}`);
            return;
        }

        const uploadArea =
            document.getElementById(`${containerId}-area`) ||
            document.querySelector(`.${containerId}-area`) ||
            document.getElementById(containerId) ||
            document.querySelector(`.${containerId}`) ||
            fileInput.closest('.upload-area');

        if (!uploadArea) {
            console.warn(`Upload area not found for: ${containerId}`);
            return;
        }

        console.log(`Found upload elements for ${containerId}`);

        uploadArea.addEventListener('click', function(e) {
            if (e.target.tagName === 'BUTTON' || e.target.closest('button')) {
                return;
            }
            fileInput.click();
        });

        fileInput.addEventListener('change', function() {
            if (this.files && this.files[0]) {
                handleFileSelection(this.files[0], previewContainer, acceptedTypes);
            }
        });

        uploadArea.addEventListener('dragover', function(e) {
            e.preventDefault();
            e.stopPropagation();
            this.classList.add('dragover');
        });

        uploadArea.addEventListener('dragleave', function(e) {
            e.preventDefault();
            e.stopPropagation();
            this.classList.remove('dragover');
        });

        uploadArea.addEventListener('drop', function(e) {
            e.preventDefault();
            e.stopPropagation();
            this.classList.remove('dragover');

            if (e.dataTransfer.files && e.dataTransfer.files[0]) {
                fileInput.files = e.dataTransfer.files;
                handleFileSelection(e.dataTransfer.files[0], previewContainer, acceptedTypes);
            }
        });
    }

    try {
        setupFileUpload('hide-text-upload', 'hide-text-file', 'hide-text-preview-container', ['image/jpeg', 'image/png']);
        setupFileUpload('hide-file-upload', 'hide-file-image', 'hide-file-preview-container', ['image/jpeg', 'image/png']);
        setupFileUpload('hide-file-file', 'hide-file-file', 'hide-file-file-preview-container', null);
        setupFileUpload('extract-upload', 'extract-file', 'extract-preview-container', ['image/jpeg', 'image/png']);
        setupFileUpload('analyze-upload', 'analyze-file', 'analyze-preview-container', ['image/jpeg', 'image/png']);
    } catch (e) {
        console.error('Error setting up file uploads:', e);
    }
}

/**
 * Handle file selection and preview
 * @param {File} file - Selected file
 * @param {HTMLElement} previewContainer - Container for preview
 * @param {Array} acceptedTypes - Array of accepted mime types
 */
function handleFileSelection(file, previewContainer, acceptedTypes) {
    if (!file || !previewContainer) return;
    if (acceptedTypes && !acceptedTypes.includes(file.type) && acceptedTypes[0] !== null) {
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(
                `File type "${file.type}" is not supported. Please select ${acceptedTypes.join(' or ')}.`,
                'error',
                'Invalid File'
            );
        } else {
            alert(`File type "${file.type}" is not supported. Please select ${acceptedTypes.join(' or ')}.`);
        }
        return;
    }

    previewContainer.style.display = 'block';

    previewContainer.innerHTML = '';

    if (file.type.startsWith('image/')) {
        createImagePreview(file, previewContainer);
    } else {
        createFilePreview(file, previewContainer);
    }

    if (file.type.startsWith('image/') &&
        (previewContainer.id === 'hide-text-preview-container' ||
            previewContainer.id === 'hide-file-preview-container')) {
        updateCapacity(file, previewContainer.id.split('-')[1]);
    }

    const form = previewContainer.closest('form');
    if (form) {
        const submitBtn = form.querySelector('button[type="submit"]');
        if (submitBtn) {
            submitBtn.disabled = false;
        }
    }
}

/**
 * Create image preview
 * @param {File} file - Image file
 * @param {HTMLElement} container - Container for preview
 */
function createImagePreview(file, container) {
    const reader = new FileReader();

    reader.onload = function(e) {
        container.innerHTML = `
            <img src="${e.target.result}" alt="Preview" class="preview-image">
            <div class="preview-info">
                <div class="file-info">
                    <div class="file-name">${file.name}</div>
                    <div class="file-size">${typeof StegGoAPI !== 'undefined' && StegGoAPI.formatFileSize ? StegGoAPI.formatFileSize(file.size) : formatFileSizeFallback(file.size)}</div>
                </div>
                <button type="button" class="remove-file" aria-label="Remove file">
                    <i class='bx bx-trash'></i>
                </button>
            </div>
        `;

        const removeBtn = container.querySelector('.remove-file');
        if (removeBtn) {
            removeBtn.addEventListener('click', function() {
                const form = container.closest('form');
                if (form) {
                    const fileInput = form.querySelector('input[type="file"]');
                    if (fileInput) {
                        fileInput.value = '';
                    }

                    const submitBtn = form.querySelector('button[type="submit"]');
                    if (submitBtn) {
                        submitBtn.disabled = true;
                    }
                }

                container.style.display = 'none';
                container.innerHTML = '';

                const operation = container.id.split('-')[1];
                const capacityMeter = document.getElementById(`capacity-indicator`);
                if (capacityMeter && (operation === 'text' || operation === 'file')) {
                    capacityMeter.style.display = 'none';
                }
            });
        }
    };

    reader.readAsDataURL(file);
}

/**
 * Create file preview for non-image files
 * @param {File} file - File to preview
 * @param {HTMLElement} container - Container for preview
 */
function createFilePreview(file, container) {
    const fileTypeIcon = typeof StegGoAPI !== 'undefined' && StegGoAPI.getFileIconClass ?
        StegGoAPI.getFileIconClass(file.type, file.name.split('.').pop()) :
        getFileIconClassFallback(file);

    container.innerHTML = `
        <div class="file-preview neu-pushed">
            <i class='bx ${fileTypeIcon}'></i>
            <div class="preview-info">
                <div class="file-info">
                    <div class="file-name">${file.name}</div>
                    <div class="file-size">${typeof StegGoAPI !== 'undefined' && StegGoAPI.formatFileSize ? StegGoAPI.formatFileSize(file.size) : formatFileSizeFallback(file.size)}</div>
                </div>
                <button type="button" class="remove-file" aria-label="Remove file">
                    <i class='bx bx-trash'></i>
                </button>
            </div>
        </div>
    `;

    const removeBtn = container.querySelector('.remove-file');
    if (removeBtn) {
        removeBtn.addEventListener('click', function() {
            const form = container.closest('form');
            if (form) {
                const fileInput = form.querySelector('input[type="file"]');
                if (fileInput) {
                    fileInput.value = '';
                }

                const submitBtn = form.querySelector('button[type="submit"]');
                if (submitBtn) {
                    submitBtn.disabled = true;
                }
            }

            container.style.display = 'none';
            container.innerHTML = '';
        });
    }
}

/**
 * Format file size (fallback function)
 * @param {number} bytes - File size in bytes
 * @returns {string} - Formatted file size
 */
function formatFileSizeFallback(bytes) {
    if (bytes === 0) return '0 Bytes';

    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

/**
 * Get file icon class (fallback function)
 * @param {File} file - File object
 * @returns {string} - BoxIcon class
 */
function getFileIconClassFallback(file) {
    const type = file.type;
    const extension = file.name.split('.').pop().toLowerCase();

    if (type.includes('image')) return 'bx-image';
    if (type.includes('pdf') || extension === 'pdf') return 'bx-file-pdf';
    if (type.includes('word') || extension === 'doc' || extension === 'docx') return 'bx-file-doc';
    if (type.includes('excel') || extension === 'xls' || extension === 'xlsx') return 'bx-spreadsheet';
    if (type.includes('video')) return 'bx-video';
    if (type.includes('audio')) return 'bx-music';
    if (type.includes('zip') || type.includes('compressed') ||
        ['zip', 'rar', '7z', 'tar', 'gz'].includes(extension)) return 'bx-archive';
    if (type.includes('text') || extension === 'txt') return 'bx-file-txt';

    return 'bx-file-blank';
}

/**
 * Update capacity meter
 * @param {File} imageFile - Image file
 * @param {string} operation - Operation type ('text' or 'file')
 */
function updateCapacity(imageFile, operation) {
    const capacityMeter = document.getElementById('capacity-indicator');
    if (!capacityMeter) return;

    capacityMeter.style.display = 'block';

    const capacityText = document.getElementById('capacity-text');
    const capacityFill = document.getElementById('capacity-fill');

    if (capacityText) {
        capacityText.textContent = 'Calculating...';
    }
    if (capacityFill) {
        capacityFill.style.width = '0%';
    }

    if (typeof StegGoAPI !== 'undefined' && StegGoAPI.calculateCapacity) {
        StegGoAPI.calculateCapacity(imageFile)
            .then(capacity => {
                if (!capacity) return;

                let usedBytes = 0;
                let percentageUsed = 0;

                if (operation === 'text') {
                    const messageTextarea = document.getElementById('hide-text-message');
                    if (messageTextarea) {
                        usedBytes = new Blob([messageTextarea.value]).size;
                        percentageUsed = (usedBytes / capacity.maxBytes) * 100;

                        messageTextarea.addEventListener('input', function() {
                            const newUsedBytes = new Blob([this.value]).size;
                            const newPercentage = (newUsedBytes / capacity.maxBytes) * 100;

                            if (capacityText) {
                                capacityText.textContent = `${StegGoAPI.formatFileSize(newUsedBytes)} / ${capacity.maxBytesFormatted} (${Math.min(100, newPercentage.toFixed(1))}%)`;
                            }
                            if (capacityFill) {
                                capacityFill.style.width = `${Math.min(100, newPercentage)}%`;
                            }

                            if (newPercentage > 90 && newPercentage <= 100) {
                                capacityFill.style.backgroundColor = '#f59e0b';
                            } else if (newPercentage > 100) {
                                capacityFill.style.backgroundColor = '#ef4444';
                                if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                                    UIManager.showToast('Message exceeds image capacity. Some data may be lost.', 'warning', 'Capacity Warning');
                                }
                            } else {
                                capacityFill.style.backgroundColor = '';
                            }
                        });
                    }
                } else if (operation === 'file') {
                    const fileInput = document.getElementById('hide-file-file');
                    if (fileInput && fileInput.files[0]) {
                        usedBytes = fileInput.files[0].size;
                        percentageUsed = (usedBytes / capacity.maxBytes) * 100;

                        if (percentageUsed > 100) {
                            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                                UIManager.showToast('File is too large for this image. Please choose a smaller file or a larger image.', 'error', 'Capacity Error');
                            }
                        } else if (percentageUsed > 90) {
                            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                                UIManager.showToast('File is approaching the image capacity limit.', 'warning', 'Capacity Warning');
                            }
                        }
                    }
                }

                if (capacityText) {
                    capacityText.textContent = `${StegGoAPI.formatFileSize(usedBytes)} / ${capacity.maxBytesFormatted} (${Math.min(100, percentageUsed.toFixed(1))}%)`;
                }
                if (capacityFill) {
                    capacityFill.style.width = `${Math.min(100, percentageUsed)}%`;

                    if (percentageUsed > 90 && percentageUsed <= 100) {
                        capacityFill.style.backgroundColor = '#f59e0b';
                    } else if (percentageUsed > 100) {
                        capacityFill.style.backgroundColor = '#ef4444';
                    }
                }
            })
            .catch(error => {
                console.error('Error calculating capacity:', error);
                if (capacityText) {
                    capacityText.textContent = 'Error calculating capacity';
                }
                if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                    UIManager.showToast('Failed to calculate image capacity.', 'error');
                }
            });
    } else {
        estimateCapacityFromImage(imageFile, operation, capacityText, capacityFill);
    }
}

/**
 * Estimate capacity from image dimensions (fallback when API not available)
 * @param {File} imageFile - Image file
 * @param {string} operation - Operation type
 * @param {HTMLElement} capacityText - Text element
 * @param {HTMLElement} capacityFill - Fill element
 */
function estimateCapacityFromImage(imageFile, operation, capacityText, capacityFill) {
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

            let usedBytes = 0;

            if (operation === 'text') {
                const messageTextarea = document.getElementById('hide-text-message');
                if (messageTextarea) {
                    usedBytes = new Blob([messageTextarea.value]).size;

                    messageTextarea.addEventListener('input', function() {
                        updateCapacityDisplay(new Blob([this.value]).size, usableCapacityBytes);
                    });
                }
            } else if (operation === 'file') {
                const fileInput = document.getElementById('hide-file-file');
                if (fileInput && fileInput.files[0]) {
                    usedBytes = fileInput.files[0].size;
                }
            }

            updateCapacityDisplay(usedBytes, usableCapacityBytes);

            function updateCapacityDisplay(usedBytes, maxBytes) {
                const percentageUsed = (usedBytes / maxBytes) * 100;
                if (capacityText) {
                    capacityText.textContent = `${formatFileSizeFallback(usedBytes)} / ${formatFileSizeFallback(maxBytes)} (${Math.min(100, percentageUsed.toFixed(1))}%)`;
                }
                if (capacityFill) {
                    capacityFill.style.width = `${Math.min(100, percentageUsed)}%`;

                    if (percentageUsed > 90 && percentageUsed <= 100) {
                        capacityFill.style.backgroundColor = '#f59e0b';
                    } else if (percentageUsed > 100) {
                        capacityFill.style.backgroundColor = '#ef4444';
                    } else {
                        capacityFill.style.backgroundColor = '';
                    }
                }
            }
        };
        img.src = e.target.result;
    };

    reader.readAsDataURL(imageFile);
}

function initCloseButtons() {
    document.querySelectorAll('.close-results').forEach(button => {
        button.addEventListener('click', function() {
            const resultsSection = this.closest('.results-section');
            if (resultsSection) {
                resultsSection.style.display = 'none';
            }
        });
    });
}

/**
 * Handle hide text form submission
 * @param {Event} e - Submit event
 */
async function handleHideTextSubmit(e) {
    e.preventDefault();

    try {
        const imageFile = document.getElementById('hide-text-file').files[0];
        const message = document.getElementById('hide-text-message').value;
        const passwordInput = document.getElementById('hide-text-password');
        const password = passwordInput ? passwordInput.value : '';

        if (!imageFile) {
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('Please select an image file.', 'error');
            }
            return;
        }

        if (!message || message.trim() === '') {
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('Please enter a message to hide.', 'error');
            }
            return;
        }

        const formData = new FormData();
        formData.append('image', imageFile);
        formData.append('message', message);
        if (password) {
            formData.append('password', password);
        }

        if (typeof UIManager !== 'undefined' && UIManager.showLoading) {
            UIManager.showLoading();
        }

        let response;
        if (typeof StegGoAPI !== 'undefined' && StegGoAPI.hideText) {
            response = await StegGoAPI.hideText(formData);

            const processedResponse = StegGoAPI.processHideTextResponse(response);

            displayHideTextResults(processedResponse);
        } else {
            console.error('StegGoAPI not available');
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('API not available. Please try again later.', 'error');
            }
        }
    } catch (error) {
        console.error('Error in hide text operation:', error);
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(error.message || 'An error occurred during the operation.', 'error');
        }
    } finally {
        if (typeof UIManager !== 'undefined' && UIManager.hideLoading) {
            UIManager.hideLoading();
        }
        if (passwordInput) passwordInput.value = ''; // Clear password field
    }
}

/**
 * Display hide text results
 * @param {Object} response - Processed API response
 */
function displayHideTextResults(response) {
    if (!response || !response.success) {
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(response.message || 'Operation failed. Please try again.', 'error');
        }
        return;
    }

    // Remove key display, add encryption status display
    const encryptionStatusDiv = document.getElementById('hide-text-encryption-status');
    if (encryptionStatusDiv && response.encryption) {
        encryptionStatusDiv.innerHTML = `
            <div class="detail-item">
                <span class="detail-item-label">Encryption:</span>
                <span class="detail-item-value ${response.encryption === 'enabled' ? 'text-success' : 'text-muted'}">
                    ${response.encryption.charAt(0).toUpperCase() + response.encryption.slice(1)}
                </span>
            </div>
        `;
    } else if (encryptionStatusDiv) {
        encryptionStatusDiv.innerHTML = ''; // Clear if no status
    }

    const resultImage = document.getElementById('hide-text-result-img');
    if (resultImage && response.image) {
        if (response.image.url) {
            const imageUrl = response.image.url.startsWith('http') ?
                response.image.url :
                (response.image.url.startsWith('/') ?
                    window.location.origin + response.image.url :
                    response.image.url);

            resultImage.src = imageUrl;
            console.log('Setting result image to:', imageUrl);
        } else {
            console.error('Missing image URL in response:', response);
        }
    }
    const downloadLink = document.getElementById('hide-text-download');
    if (downloadLink && response.image && response.image.download) {
        const downloadUrl = response.image.download.startsWith('http') ?
            response.image.download :
            (response.image.download.startsWith('/') ?
                window.location.origin + response.image.download :
                response.image.download);

        downloadLink.href = downloadUrl;
        downloadLink.download = getFileNameFromUrl(downloadUrl) || 'stego-image.png';
    }

    const resultsSection = document.getElementById('hide-text-results');
    if (resultsSection) {
        resultsSection.style.display = 'block';

        resultsSection.scrollIntoView({ behavior: 'smooth' });
    }

    if (typeof UIManager !== 'undefined' && UIManager.showToast) {
        UIManager.showToast('Message hidden successfully!', 'success');
    }
}

/**
 * Handle hide file form submission
 * @param {Event} e - Submit event
 */
async function handleHideFileSubmit(e) {
    e.preventDefault();

    try {
        const imageFile = document.getElementById('hide-file-image').files[0];
        const fileToHide = document.getElementById('hide-file-file').files[0];
        const passwordInput = document.getElementById('hide-file-password');
        const password = passwordInput ? passwordInput.value : '';

        if (!imageFile) {
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('Please select a cover image.', 'error');
            }
            return;
        }

        if (!fileToHide) {
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('Please select a file to hide.', 'error');
            }
            return;
        }

        const formData = new FormData();
        formData.append('image', imageFile);
        formData.append('file', fileToHide);
        if (password) {
            formData.append('password', password);
        }

        if (typeof UIManager !== 'undefined' && UIManager.showLoading) {
            UIManager.showLoading();
        }

        let response;
        if (typeof StegGoAPI !== 'undefined' && StegGoAPI.hideFile) {
            response = await StegGoAPI.hideFile(formData);

            const processedResponse = StegGoAPI.processHideFileResponse(response);

            displayHideFileResults(processedResponse);
        } else {
            console.error('StegGoAPI not available');
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('API not available. Please try again later.', 'error');
            }
        }
    } catch (error) {
        console.error('Error in hide file operation:', error);
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(error.message || 'An error occurred during the operation.', 'error');
        }
    } finally {
        if (typeof UIManager !== 'undefined' && UIManager.hideLoading) {
            UIManager.hideLoading();
        }
        if (passwordInput) passwordInput.value = ''; // Clear password field
    }
}

/**
 * Display hide file results
 * @param {Object} response - API response
 */
function displayHideFileResults(response) {
    if (!response || !response.success) {
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(response.message || 'Operation failed. Please try again.', 'error');
        }
        return;
    }

    // Remove key display, add encryption status display
    const encryptionStatusDiv = document.getElementById('hide-file-encryption-status');
     if (encryptionStatusDiv && response.encryption) {
        encryptionStatusDiv.innerHTML = `
            <div class="detail-item">
                <span class="detail-item-label">Encryption:</span>
                <span class="detail-item-value ${response.encryption === 'enabled' ? 'text-success' : 'text-muted'}">
                    ${response.encryption.charAt(0).toUpperCase() + response.encryption.slice(1)}
                </span>
            </div>
        `;
    } else if (encryptionStatusDiv) {
        encryptionStatusDiv.innerHTML = ''; // Clear if no status
    }
    
    const resultImage = document.getElementById('hide-file-result-img');
    if (resultImage && response.image && response.image.url) {
        const imageUrl = response.image.url.startsWith('http') ?
            response.image.url :
            (response.image.url.startsWith('/') ?
                window.location.origin + response.image.url :
                response.image.url);

        resultImage.src = imageUrl;
    }

    const downloadLink = document.getElementById('hide-file-download');
    if (downloadLink && response.image && response.image.download) {
        const downloadUrl = response.image.download.startsWith('http') ?
            response.image.download :
            (response.image.download.startsWith('/') ?
                window.location.origin + response.image.download :
                response.image.download);

        downloadLink.href = downloadUrl;
        downloadLink.download = getFileNameFromUrl(downloadUrl) || 'stego-image.png';
    }

    const detailsContainer = document.getElementById('hide-file-details');
    if (detailsContainer && response.hiddenFile) {
        const fileInfo = response.hiddenFile;
        const formatFileSize = typeof StegGoAPI !== 'undefined' && StegGoAPI.formatFileSize ?
            StegGoAPI.formatFileSize : formatFileSizeFallback;

        detailsContainer.innerHTML = `
            <div class="detail-item">
                <span class="detail-item-label">File Name:</span>
                <span class="detail-item-value">${fileInfo.name || 'Unknown'}</span>
            </div>
            <div class="detail-item">
                <span class="detail-item-label">File Size:</span>
                <span class="detail-item-value">${formatFileSize(fileInfo.size) || 'Unknown'}</span>
            </div>
            <div class="detail-item">
                <span class="detail-item-label">File Type:</span>
                <span class="detail-item-value">${fileInfo.type || 'Unknown'}</span>
            </div>
        `;
    } else if (detailsContainer) {
        const fileToHide = document.getElementById('hide-file-file').files[0];
        if (fileToHide) {
            const formatFileSize = typeof StegGoAPI !== 'undefined' && StegGoAPI.formatFileSize ?
                StegGoAPI.formatFileSize : formatFileSizeFallback;

            detailsContainer.innerHTML = `
                <div class="detail-item">
                    <span class="detail-item-label">File Name:</span>
                    <span class="detail-item-value">${fileToHide.name}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-item-label">File Size:</span>
                    <span class="detail-item-value">${formatFileSize(fileToHide.size)}</span>
                </div>
                <div class="detail-item">
                    <span class="detail-item-label">File Type:</span>
                    <span class="detail-item-value">${fileToHide.type || 'Unknown'}</span>
                </div>
            `;
        }
    }

    const resultsSection = document.getElementById('hide-file-results');
    if (resultsSection) {
        resultsSection.style.display = 'block';

        resultsSection.scrollIntoView({ behavior: 'smooth' });
    }

    if (typeof UIManager !== 'undefined' && UIManager.showToast) {
        UIManager.showToast('File hidden successfully!', 'success');
    }
}

/**
 * Handle extracted file download using XHR instead of fetch
 * @param {Event} e - Click event
 */
function handleExtractedFileDownload(e) {
    e.preventDefault();

    const downloadUrl = this.dataset.downloadUrl;
    const fileName = this.dataset.fileName || 'extracted-file';
    const fileType = this.dataset.fileType || 'application/octet-stream';

    if (!downloadUrl) {
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast('Download URL not available', 'error');
        }
        return;
    }

    console.log(`Attempting direct download: ${fileName} from ${downloadUrl}`);

    if (typeof UIManager !== 'undefined' && UIManager.showLoading) {
        UIManager.showLoading();
    }
    const xhr = new XMLHttpRequest();
    xhr.open('GET', downloadUrl, true);
    xhr.responseType = 'arraybuffer';

    xhr.onload = function() {
        if (this.status === 200) {
            const blob = new Blob([this.response], { type: fileType });

            downloadBlob(blob, fileName);

            if (typeof UIManager !== 'undefined') {
                UIManager.hideLoading();
                UIManager.showToast('File downloaded successfully', 'success');
            }
        } else {
            console.error(`Download failed with status ${this.status}`);

            tryDirectLink(downloadUrl, fileName);

            if (typeof UIManager !== 'undefined') {
                UIManager.hideLoading();
                UIManager.showToast('Download failed. Trying alternative method...', 'warning');
            }
        }
    };

    xhr.onerror = function(error) {
        console.error('XHR download error:', error);

        tryIframeDownload(downloadUrl);

        if (typeof UIManager !== 'undefined') {
            UIManager.hideLoading();
            UIManager.showToast('Error during download. Trying alternative method...', 'warning');
        }
    };

    xhr.send();
}

/**
 * Helper function to download a blob
 * @param {Blob} blob - Blob to download
 * @param {string} fileName - Filename to use
 */
function downloadBlob(blob, fileName) {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = url;
    a.download = fileName;
    document.body.appendChild(a);
    a.click();

    setTimeout(() => {
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
    }, 100);
}

/**
 * Try direct link download as fallback
 * @param {string} url - URL to download
 * @param {string} fileName - Filename to use
 */
function tryDirectLink(url, fileName) {
    console.log('Trying direct link download...');
    const downloadUrl = url.includes('?') ?
        `${url}&download=true&filename=${encodeURIComponent(fileName)}` :
        `${url}?download=true&filename=${encodeURIComponent(fileName)}`;

    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = downloadUrl;
    a.download = fileName;
    a.target = '_blank';
    document.body.appendChild(a);
    a.click();
    setTimeout(() => {
        document.body.removeChild(a);
    }, 100);
}

/**
 * Try iframe download as a last resort
 * @param {string} url - URL to download
 */
function tryIframeDownload(url) {
    console.log('Trying iframe download...');
    const iframe = document.createElement('iframe');
    iframe.style.display = 'none';
    iframe.src = url;
    document.body.appendChild(iframe);
    setTimeout(() => {
        document.body.removeChild(iframe);
    }, 1000);
}

/**
 * Handle extracted file download click with improved error handling
 * @param {Event} e - Click event
 */
function handleExtractedFileDownload(e) {
    e.preventDefault();

    const downloadUrl = this.dataset.downloadUrl;
    const fileName = this.dataset.fileName || 'extracted-file';
    const fileType = this.dataset.fileType || 'application/octet-stream';

    if (!downloadUrl) {
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast('Download URL not available', 'error');
        }
        return;
    }

    console.log(`Attempting to download file: ${fileName} (${fileType}) from ${downloadUrl}`);
    if (typeof UIManager !== 'undefined' && UIManager.showLoading) {
        UIManager.showLoading();
    }

    const directDownload = () => {
        try {
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = downloadUrl;
            a.download = fileName;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);

            if (typeof UIManager !== 'undefined') {
                UIManager.hideLoading();
                UIManager.showToast('Download started', 'success');
            }
        } catch (error) {
            console.error('Direct download failed:', error);
            if (typeof UIManager !== 'undefined') {
                UIManager.hideLoading();
                UIManager.showToast(`Download failed. Please try again.`, 'error');
            }
        }
    };

    fetch(downloadUrl, {
        method: 'GET',
        credentials: 'include',
        headers: {
            'Accept': '*/*'
        }
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Status: ${response.status}`);
            }
            return response.blob();
        })
        .then(blob => {
            const blobWithType = new Blob([blob], { type: fileType });
            const url = window.URL.createObjectURL(blobWithType);
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = url;
            a.download = fileName;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);

            if (typeof UIManager !== 'undefined') {
                UIManager.hideLoading();
                UIManager.showToast('File downloaded successfully', 'success');
            }
        })
        .catch(error => {
            console.error('Download error:', error);

            console.log('Trying alternative download method...');
            directDownload();
        });
}


/**
 * Handle extract form submission
 * @param {Event} e - Submit event
 */
async function handleExtractSubmit(e) {
    e.preventDefault();

    try {
        const imageFile = document.getElementById('extract-file').files[0];
        const passwordInput = document.getElementById('extract-password');
        const password = passwordInput ? passwordInput.value : '';
        const keyInput = document.getElementById('extract-key'); // Old key input

        if (!imageFile) {
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('Please select an image file.', 'error');
            }
            return;
        }

        // Password is now optional for extraction. API handles logic if it's needed but not provided.
        // The old key field is no longer used for API call.

        const formData = new FormData();
        formData.append('image', imageFile);
        if (password) {
            formData.append('password', password);
        }
        // Old: formData.append('key', key);

        if (typeof UIManager !== 'undefined' && UIManager.showLoading) {
            UIManager.showLoading();
        }

        if (typeof StegGoAPI !== 'undefined' && StegGoAPI.extract) {
            const response = await StegGoAPI.extract(formData);

            const processedResponse = StegGoAPI.processExtractResponse(response);

            displayExtractResults(processedResponse);
        } else {
            console.error('StegGoAPI not available');
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('API not available. Please try again later.', 'error');
            }
        }
    } catch (error) {
        console.error('Error in extract operation:', error);
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(error.message || 'An error occurred during the operation.', 'error');
        }
    } finally {
        if (typeof UIManager !== 'undefined' && UIManager.hideLoading) {
            UIManager.hideLoading();
        }
        if (passwordInput) passwordInput.value = ''; // Clear new password field
        if (keyInput) keyInput.value = ''; // Clear old key field if it still exists in HTML (it should be removed)
    }
}

/**
 * Display extract results with direct file URL handling
 * @param {Object} response - API response
 */
function displayExtractResults(response) {
    if (!response || !response.success) {
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast('Extraction failed. Please verify the image and key.', 'error');
        }
        return;
    }

    const textResults = document.getElementById('extract-text-results');
    const fileResults = document.getElementById('extract-file-results');

    if (textResults) textResults.style.display = 'none';
    if (fileResults) fileResults.style.display = 'none';

    if (response.type === 'text') {
        const textContent = document.getElementById('extract-text-content');
        if (textContent && textResults) {
            textContent.textContent = response.content || '';
            textResults.style.display = 'block';
        }
    } else if (response.type === 'file' && response.file) {
        const fileInfo = response.file;

        if (fileResults) {
            const nameEl = document.getElementById('extract-file-name');
            const sizeEl = document.getElementById('extract-file-size');
            const typeEl = document.getElementById('extract-file-type');
            const iconEl = document.getElementById('extract-file-icon');
            const downloadLink = document.getElementById('extract-file-download');

            const formatFileSize = typeof StegGoAPI !== 'undefined' && StegGoAPI.formatFileSize ?
                StegGoAPI.formatFileSize : formatFileSizeFallback;

            if (nameEl) nameEl.textContent = fileInfo.name || 'Unknown';
            if (sizeEl) sizeEl.textContent = formatFileSize(fileInfo.size) || 'Unknown';
            if (typeEl) typeEl.textContent = fileInfo.type || 'Unknown';

            if (iconEl) {
                iconEl.className = `bx ${fileInfo.icon || 'bx-file'}`;
            }

            if (downloadLink && fileInfo.url) {
                console.log('Setting up download with original URL:', fileInfo.url);

                downloadLink.href = fileInfo.url;

                downloadLink.target = '_blank';
                downloadLink.setAttribute('download', fileInfo.name || 'extracted-file');
                downloadLink.onclick = function(e) {
                    console.log(`Attempting to download using direct URL: ${this.href}`);

                    if (needsManualDownloadHandling()) {
                        e.preventDefault();
                        window.open(this.href, '_blank');
                    }
                };
            }

            fileResults.style.display = 'block';
        }
    }

    const resultsSection = document.getElementById('extract-results');
    if (resultsSection) {
        resultsSection.style.display = 'block';

        resultsSection.scrollIntoView({ behavior: 'smooth' });
    }

    if (typeof UIManager !== 'undefined' && UIManager.showToast) {
        UIManager.showToast('Content extracted successfully!', 'success');
    }
}

/**
 * Check if the browser needs manual download handling
 * @returns {boolean} - Whether manual download handling is needed
 */
function needsManualDownloadHandling() {
    const ua = navigator.userAgent;
    return /safari/i.test(ua) && !/chrome|chromium/i.test(ua) ||
        /mobile/i.test(ua) ||
        /ipad|iphone|ipod/i.test(ua);
}


function downloadExtractedFile(downloadButton) {
    const fileId = downloadButton.getAttribute('data-file-id');
    const fileName = downloadButton.getAttribute('data-file-name');

    if (!fileId) {
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast('File ID not available', 'error');
        }
        return;
    }

    console.log(`Downloading extracted file. ID: ${fileId}, Name: ${fileName}`);

    if (typeof UIManager !== 'undefined' && UIManager.showLoading) {
        UIManager.showLoading();
    }

    const downloadUrl = `/api/download?id=${encodeURIComponent(fileId)}`;

    const form = document.createElement('form');
    form.method = 'GET';
    form.action = downloadUrl;

    if (fileName) {
        const filenameInput = document.createElement('input');
        filenameInput.type = 'hidden';
        filenameInput.name = 'filename';
        filenameInput.value = fileName;
        form.appendChild(filenameInput);
    }

    document.body.appendChild(form);

    try {
        form.submit();
        console.log('File download form submitted');

        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast('Download started', 'success');
        }
    } catch (error) {
        console.error('Form submission error:', error);
        window.location.href = downloadUrl + (fileName ? `&filename=${encodeURIComponent(fileName)}` : '');

        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast('Using alternative download method', 'info');
        }
    }

    setTimeout(() => {
        if (form.parentNode) {
            document.body.removeChild(form);
        }

        if (typeof UIManager !== 'undefined' && UIManager.hideLoading) {
            UIManager.hideLoading();
        }
    }, 1000);
}

/**
 * Handle download link click to properly download the file
 * @param {Event} e - Click event
 */
function handleDownloadClick(e) {
    e.preventDefault();

    const downloadUrl = this.dataset.downloadUrl;
    const fileName = this.dataset.fileName || 'extracted-file';
    const fileType = this.dataset.fileType || 'application/octet-stream';

    if (!downloadUrl) {
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast('Download URL not available', 'error');
        }
        return;
    }

    console.log(`Attempting to download file: ${fileName} (${fileType}) from ${downloadUrl}`);

    if (typeof UIManager !== 'undefined' && UIManager.showLoading) {
        UIManager.showLoading();
    }

    fetch(downloadUrl)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.blob();
        })
        .then(blob => {
            const fileBlob = new Blob([blob], { type: fileType });

            const url = window.URL.createObjectURL(fileBlob);
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = url;
            a.download = fileName;
            document.body.appendChild(a);
            a.click();


            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);

            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('File download started', 'success');
            }
        })
        .catch(error => {
            console.error('Download error:', error);
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast(`Download failed: ${error.message}`, 'error');
            }
        })
        .finally(() => {
            if (typeof UIManager !== 'undefined' && UIManager.hideLoading) {
                UIManager.hideLoading();
            }
        });
}

/**
 * Handle analyze form submission
 * @param {Event} e - Submit event
 */
async function handleAnalyzeSubmit(e) {
    e.preventDefault();

    try {
        const imageFile = document.getElementById('analyze-file').files[0];

        if (!imageFile) {
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('Please select an image file.', 'error');
            }
            return;
        }

        const formData = new FormData();
        formData.append('image', imageFile);

        if (typeof UIManager !== 'undefined' && UIManager.showLoading) {
            UIManager.showLoading();
        }

        if (typeof StegGoAPI !== 'undefined' && StegGoAPI.analyze) {
            const response = await StegGoAPI.analyze(formData);
            const processedResponse = StegGoAPI.processAnalyzeResponse(response);

            displayAnalyzeResults(processedResponse);
        } else {
            console.error('StegGoAPI not available');
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('API not available. Please try again later.', 'error');
            }
        }
    } catch (error) {
        console.error('Error in analyze operation:', error);
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(error.message || 'An error occurred during the operation.', 'error');
        }
    } finally {
        if (typeof UIManager !== 'undefined' && UIManager.hideLoading) {
            UIManager.hideLoading();
        }
    }
}

/**
 * Display analyze results
 * @param {Object} response - API response
 */
function displayAnalyzeResults(response) {
    if (!response || !response.success) {
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast('Analysis failed. Please try again.', 'error');
        }
        return;
    }

    const formatFileSize = typeof StegGoAPI !== 'undefined' && StegGoAPI.formatFileSize ?
        StegGoAPI.formatFileSize : formatFileSizeFallback;

    const imageDetails = document.getElementById('analyze-image-details');
    if (imageDetails && response.imageInfo) {
        const info = response.imageInfo;

        imageDetails.innerHTML = `
            <div class="detail-item">
                <span class="detail-item-label">File Name:</span>
                <span class="detail-item-value">${info.name || 'Unknown'}</span>
            </div>
            <div class="detail-item">
                <span class="detail-item-label">Dimensions:</span>
                <span class="detail-item-value">${info.width || 0} × ${info.height || 0} px</span>
            </div>
            <div class="detail-item">
                <span class="detail-item-label">File Size:</span>
                <span class="detail-item-value">${info.size || formatFileSize(0)}</span>
            </div>
            <div class="detail-item">
                <span class="detail-item-label">Format:</span>
                <span class="detail-item-value">${info.type ? (info.type.split('/')[1] || '').toUpperCase() : 'Unknown'}</span>
            </div>
        `;
    }

    const capacityDetails = document.getElementById('analyze-capacity-details');
    if (capacityDetails && response.capacity) {
        const capacity = response.capacity;

        capacityDetails.innerHTML = `
            <div class="detail-item">
                <span class="detail-item-label">Max Capacity:</span>
                <span class="detail-item-value">${capacity.maxBytesFormatted || formatFileSize(capacity.maxBytes || 0)}</span>
            </div>
            <div class="detail-item">
                <span class="detail-item-label">LSB Bits:</span>
                <span class="detail-item-value">${capacity.lsbBits || 3}</span>
            </div>
            <div class="detail-item">
                <span class="detail-item-label">Used:</span>
                <span class="detail-item-value">${capacity.percentageUsed || 0}%</span>
            </div>
            <div class="detail-item">
                <span class="detail-item-label">Recommendation:</span>
                <span class="detail-item-value">${capacity.recommendation || 'Safe for steganography.'}</span>
            </div>
        `;
    }

    const privacyRisks = document.getElementById('analyze-privacy-risks');
    if (privacyRisks) {
        privacyRisks.innerHTML = '<h4>Privacy Risks</h4>';

        if (response.privacyRisks && response.privacyRisks.length > 0) {
            response.privacyRisks.forEach(risk => {
                const riskItem = document.createElement('div');
                riskItem.className = 'detail-item';
                riskItem.innerHTML = `
                    <span class="detail-item-label">
                        <i class='bx ${risk.level === 'info' ? 'bx-info-circle text-info' : 'bx-error-circle text-warning'}'></i>
                    </span>
                    <span class="detail-item-value">${risk.message || ''}</span>
                `;
                privacyRisks.appendChild(riskItem);
            });
        } else {
            privacyRisks.innerHTML += `
                <div class="detail-item">
                    <span class="detail-item-label">
                        <i class='bx bx-check-circle text-success'></i>
                    </span>
                    <span class="detail-item-value">No privacy risks detected.</span>
                </div>
            `;
        }
    }

    const exifContainer = document.getElementById('analyze-exif-container');
    if (exifContainer) {
        exifContainer.innerHTML = '<h4>EXIF Metadata</h4>';

        if (response.exifData && Object.keys(response.exifData).length > 0) {
            let tableHTML = `
                <div class="table-container">
                    <table class="exif-table">
                        <thead>
                            <tr>
                                <th>Property</th>
                                <th>Value</th>
                            </tr>
                        </thead>
                        <tbody>
            `;

            for (const [key, value] of Object.entries(response.exifData)) {
                tableHTML += `
                    <tr>
                        <td>${key}</td>
                        <td>${value}</td>
                    </tr>
                `;
            }

            tableHTML += `
                        </tbody>
                    </table>
                </div>
            `;

            exifContainer.innerHTML += tableHTML;
        } else {
            exifContainer.innerHTML += `
                <p class="text-muted">No EXIF data found in this image.</p>
            `;
        }
    }

    const resultsSection = document.getElementById('analyze-results');
    if (resultsSection) {
        resultsSection.style.display = 'block';

        resultsSection.scrollIntoView({ behavior: 'smooth' });
    }

    if (typeof UIManager !== 'undefined' && UIManager.showToast) {
        UIManager.showToast('Image analysis completed successfully!', 'success');
    }
}

/**
 * Get a filename from URL
 * @param {string} url - URL to extract filename from
 * @returns {string} - Extracted filename or default value
 */
function getFileNameFromUrl(url) {
    if (!url) return 'file';

    const urlWithoutParams = url.split('?')[0];

    const parts = urlWithoutParams.split('/');
    const filename = parts[parts.length - 1];

    return filename || 'file';
}

/**
 * Copy text to clipboard
 * @param {string} elementId - Element ID containing text to copy
 */
function copyToClipboard(elementId) {
    const element = document.getElementById(elementId);
    if (!element) return;

    let textToCopy = '';

    if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
        textToCopy = element.value;
    } else {
        textToCopy = element.textContent;
    }

    if (!textToCopy) return;

    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(textToCopy)
            .then(() => {
                if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                    UIManager.showToast('Copied to clipboard!', 'success', null, 2000);
                }
            })
            .catch(err => {
                console.error('Error copying to clipboard:', err);
                fallbackCopyToClipboard(textToCopy);
            });
    } else {
        fallbackCopyToClipboard(textToCopy);
    }
}

/**
 * Fallback method to copy text to clipboard
 * @param {string} text - Text to copy
 */
function fallbackCopyToClipboard(text) {
    const tempElement = document.createElement('textarea');
    tempElement.value = text;
    tempElement.style.position = 'absolute';
    tempElement.style.left = '-9999px';
    document.body.appendChild(tempElement);
    tempElement.select();
    document.execCommand('copy');

    document.body.removeChild(tempElement);

    if (typeof UIManager !== 'undefined' && UIManager.showToast) {
        UIManager.showToast('Copied to clipboard!', 'success', null, 2000);
    }
}

document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.btn-copy').forEach(button => {
        button.addEventListener('click', function() {
            const targetId = this.getAttribute('data-copy');
            if (targetId) {
                copyToClipboard(targetId);
            }
        });
    });
});
