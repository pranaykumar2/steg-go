const StegGoAPI = (function() {
    const API_BASE_URL = '/api';

    const ENDPOINTS = {
        HIDE_TEXT: '/hide',
        HIDE_FILE: '/hideFile',
        EXTRACT: '/extract',
        ANALYZE: '/metadata',
        HEALTH: '/health'
    };

    /**
     * Generic function to make API requests
     * @param {string} endpoint - API endpoint
     * @param {string} method - HTTP method (GET, POST, etc.)
     * @param {FormData|Object} data - Data to send to the API
     * @returns {Promise} - Promise resolving to the API response
     */
    async function makeRequest(endpoint, method = 'GET', data = null) {
        const url = API_BASE_URL + endpoint;

        const options = {
            method: method,
            headers: {},
            credentials: 'same-origin'
        };

        if (data) {
            if (data instanceof FormData) {
                options.body = data;
            } else {
                options.headers['Content-Type'] = 'application/json';
                options.body = JSON.stringify(data);
            }
        }

        try {
            if (typeof UIManager !== 'undefined' && UIManager.showLoading) {
                UIManager.showLoading();
            }

            const response = await fetch(url, options);

            if (!response.ok) {
                let errorMessage = 'An unknown error occurred';
                try {
                    const errorData = await response.json();
                    errorMessage = errorData.error || errorData.message || errorMessage;
                } catch (e) {
                    // If unable to parse JSON, use status text
                    errorMessage = response.statusText || errorMessage;
                }

                throw new Error(errorMessage);
            }

            const responseData = await response.json();

            if (responseData && responseData.success === false) {
                throw new Error(responseData.error || 'Operation failed');
            }

            return responseData;
        } catch (error) {
            console.error('API Request Error:', error);
            throw error;
        } finally {
            // Hide loading indicator if UIManager is available
            if (typeof UIManager !== 'undefined' && UIManager.hideLoading) {
                UIManager.hideLoading();
            }
        }
    }

    /**
     * Check API health
     * @returns {Promise} - Promise resolving to the health status
     */
    async function getHealth() {
        return makeRequest(ENDPOINTS.HEALTH, 'GET');
    }

    /**
     * Hide text in an image
     * @param {FormData} formData - Form data with image and message
     * @returns {Promise} - Promise resolving to the API response
     */
    async function hideText(formData) {
        const apiFormData = new FormData();
        const message = formData.get('message');
        if (message) {
            apiFormData.append('Message', message);
        }

        const image = formData.get('image');
        if (image) {
            apiFormData.append('image', image);
        }

        return makeRequest(ENDPOINTS.HIDE_TEXT, 'POST', apiFormData);
    }

    /**
     * Hide file in an image
     * @param {FormData} formData - Form data with image and file to hide
     * @returns {Promise} - Promise resolving to the API response
     */
    async function hideFile(formData) {
        const apiFormData = new FormData();

        const image = formData.get('image');
        if (image) {
            apiFormData.append('image', image);
        }

        const file = formData.get('file');
        if (file) {
            apiFormData.append('file', file);
        }

        return makeRequest(ENDPOINTS.HIDE_FILE, 'POST', apiFormData);
    }

    /**
     * Extract hidden content from an image
     * @param {FormData} formData - Form data with image and decryption key
     * @returns {Promise} - Promise resolving to the API response
     */
    async function extract(formData) {
        const apiFormData = new FormData();

        const key = formData.get('key');
        if (key) {
            apiFormData.append('Key', key);
        }

        const image = formData.get('image');
        if (image) {
            apiFormData.append('image', image);
        }

        return makeRequest(ENDPOINTS.EXTRACT, 'POST', apiFormData);
    }

    /**
     * Analyze an image for steganographic properties
     * @param {FormData} formData - Form data with image to analyze
     * @returns {Promise} - Promise resolving to the API response
     */
    async function analyze(formData) {
        return makeRequest(ENDPOINTS.ANALYZE, 'POST', formData);
    }

    /**
     * Calculate estimated capacity for image steganography
     * @param {File} imageFile - The image file to calculate capacity for
     * @returns {Promise} - Promise resolving to the capacity details
     */
    async function calculateCapacity(imageFile) {
        console.log('API: Calculating capacity for image:', imageFile.name);
        return new Promise((resolve) => {
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

                    resolve({
                        maxBytes: usableCapacityBytes,
                        maxBytesFormatted: formatFileSize(usableCapacityBytes),
                        percentageUsed: 0,
                        lsbBits: 3,
                        recommendation: `This image can store approximately ${formatFileSize(usableCapacityBytes)} of data safely.`
                    });
                };

                img.onerror = function() {
                    console.error('Error loading image for capacity calculation');
                    // Provide a fallback capacity
                    resolve({
                        maxBytes: 10240, // 10KB fallback
                        maxBytesFormatted: '10 KB',
                        percentageUsed: 0,
                        lsbBits: 3,
                        recommendation: 'Unable to calculate exact capacity. Using a conservative estimate.'
                    });
                };

                img.src = e.target.result;
            };

            reader.onerror = function() {
                console.error('Error reading file for capacity calculation');
                // Provide a fallback capacity
                resolve({
                    maxBytes: 10240, // 10KB fallback
                    maxBytesFormatted: '10 KB',
                    percentageUsed: 0,
                    lsbBits: 3,
                    recommendation: 'Unable to calculate capacity. Using a conservative estimate.'
                });
            };

            reader.readAsDataURL(imageFile);
        });
    }

    /**
     * Process the response from hide text operation
     * @param {Object} response - API response
     * @returns {Object} - Processed response
     */
    function processHideTextResponse(response) {
        console.log('Processing hide text response:', response);

        if (!response) {
            throw new Error('Invalid response from server');
        }

        const data = response.data || response;

        return {
            success: response.success !== false, // Default to true if not explicitly false
            message: response.message || data.message || 'Message hidden successfully',
            key: data.Key || data.key,
            image: {
                url: data.OutputFileURL || data.outputFileURL || '',
                download: data.OutputFileURL || data.outputFileURL || '',
                type: 'image/png', // Assume PNG output
                size: null // Not provided by API
            }
        };
    }

    /**
     * Process the response from hide file operation
     * @param {Object} response - API response
     * @returns {Object} - Processed response
     */
    function processHideFileResponse(response) {
        console.log('Processing hide file response:', response);

        if (!response) {
            throw new Error('Invalid response from server');
        }

        // Check if response has nested data property
        const data = response.data || response;

        return {
            success: response.success !== false, // Default to true if not explicitly false
            message: response.message || data.message || 'File hidden successfully',
            key: data.Key || data.key,
            image: {
                url: data.OutputFileURL || data.outputFileURL || '',
                download: data.OutputFileURL || data.outputFileURL || '',
                type: 'image/png', // Assume PNG output
                size: null // Not provided by API
            },
            hiddenFile: data.FileDetails ? {
                name: data.FileDetails.OriginalName || data.FileDetails.originalName || '',
                size: data.FileDetails.FileSize || data.FileDetails.fileSize || 0,
                type: data.FileDetails.FileType || data.FileDetails.fileType || ''
            } : null
        };
    }


    /**
     * Process the response from extract operation
     * @param {Object} response - API response
     * @returns {Object} - Processed response
     */
    function processExtractResponse(response) {
        console.log('Processing extract response:', response);

        if (!response) {
            throw new Error('Invalid response from server');
        }

        if (response.success === false) {
            throw new Error(response.error || 'Extraction failed');
        }

        const data = response.data || response;
        console.log('Extract data:', data);

        if (data.IsFile || data.isFile) {
            // Get filename from response or use default
            const fileName = data.FileName || data.fileName || 'extracted-file';

            const fileURL = data.FileURL || data.fileURL || '';
            console.log('Using exact file URL from server:', fileURL);

            let contentType = data.ContentType || data.contentType || 'application/octet-stream';

            return {
                success: true,
                type: 'file',
                file: {
                    name: fileName,
                    type: contentType,
                    size: data.FileSize || data.fileSize || 0,
                    url: fileURL,
                    download: fileURL,
                    icon: getFileIconClass(contentType, fileName.split('.').pop())
                }
            };
        }
        else {
            return {
                success: true,
                type: 'text',
                content: data.Message || data.message || ''
            };
        }
    }

    /**
     * Process the response from analyze operation
     * @param {Object} response - API response
     * @returns {Object} - Processed response
     */
    function processAnalyzeResponse(response) {
        console.log('Processing analyze response:', response);

        if (!response) {
            throw new Error('Invalid response from server');
        }

        const data = response.data || response;

        const privacyRisks = [];
        if (data.PrivacyRisks && Array.isArray(data.PrivacyRisks)) {
            data.PrivacyRisks.forEach(risk => {
                privacyRisks.push({
                    level: 'warning',
                    message: risk
                });
            });
        } else if (data.privacyRisks && Array.isArray(data.privacyRisks)) {
            data.privacyRisks.forEach(risk => {
                privacyRisks.push({
                    level: 'warning',
                    message: risk
                });
            });
        }

        if (privacyRisks.length === 0) {
            privacyRisks.push({
                level: 'info',
                message: 'No privacy risks detected in this image.'
            });
        }

        const imageInfo = {
            width: data.ImageWidth || data.imageWidth || 0,
            height: data.ImageHeight || data.imageHeight || 0,
            type: data.MimeType || data.mimeType || '',
            size: formatFileSize(data.FileSize || data.fileSize || 0),
            name: data.Filename || data.filename || ''
        };

        let capacityBytes = 0;
        if (data.SteganoCapacity && data.SteganoCapacity.Bytes) {
            capacityBytes = data.SteganoCapacity.Bytes;
        } else if (data.steganoCapacity && data.steganoCapacity.bytes) {
            capacityBytes = data.steganoCapacity.bytes;
        } else if (imageInfo.width && imageInfo.height) {
            // Calculate theoretical capacity if not provided
            const pixelCount = imageInfo.width * imageInfo.height;
            capacityBytes = Math.floor((pixelCount * 3) / 8) * 0.85;
        }

        const capacity = {
            maxBytes: capacityBytes,
            maxBytesFormatted: formatFileSize(capacityBytes),
            percentageUsed: 0,
            lsbBits: 3,
            recommendation: capacityBytes > 0
                ? `This image can store approximately ${formatFileSize(capacityBytes)} of data.`
                : 'Unable to determine capacity.'
        };

        // Normalize EXIF data
        const exifData = data.Properties || data.properties || {};

        return {
            success: true,
            imageInfo: imageInfo,
            capacity: capacity,
            privacyRisks: privacyRisks,
            exifData: exifData
        };
    }

    /**
     * Get appropriate icon class based on file type
     * @param {string} mimeType - MIME type of the file
     * @param {string} extension - File extension
     * @returns {string} - BoxIcon class name
     */
    function getFileIconClass(mimeType, extension) {
        let iconClass = 'bx-file';

        if (!mimeType && !extension) {
            return iconClass;
        }

        if (mimeType) {
            if (mimeType.startsWith('image/')) {
                iconClass = 'bx-image';
            } else if (mimeType.startsWith('video/')) {
                iconClass = 'bx-video';
            } else if (mimeType.startsWith('audio/')) {
                iconClass = 'bx-music';
            } else if (mimeType.startsWith('text/')) {
                iconClass = 'bx-file-txt';
            } else if (mimeType.includes('pdf')) {
                iconClass = 'bx-file-pdf';
            } else if (mimeType.includes('zip') || mimeType.includes('compressed')) {
                iconClass = 'bx-archive';
            } else if (mimeType.includes('word')) {
                iconClass = 'bx-file-doc';
            } else if (mimeType.includes('excel') || mimeType.includes('spreadsheet')) {
                iconClass = 'bx-spreadsheet';
            } else if (mimeType.includes('presentation') || mimeType.includes('powerpoint')) {
                iconClass = 'bx-slideshow';
            }
        }

        if (iconClass === 'bx-file' && extension) {
            switch(extension.toLowerCase()) {
                case 'pdf':
                    iconClass = 'bx-file-pdf';
                    break;
                case 'doc':
                case 'docx':
                    iconClass = 'bx-file-doc';
                    break;
                case 'xls':
                case 'xlsx':
                    iconClass = 'bx-spreadsheet';
                    break;
                case 'ppt':
                case 'pptx':
                    iconClass = 'bx-slideshow';
                    break;
                case 'jpg':
                case 'jpeg':
                case 'png':
                case 'gif':
                case 'webp':
                case 'bmp':
                    iconClass = 'bx-image';
                    break;
                case 'mp3':
                case 'wav':
                case 'ogg':
                case 'flac':
                    iconClass = 'bx-music';
                    break;
                case 'mp4':
                case 'avi':
                case 'mov':
                case 'wmv':
                    iconClass = 'bx-video';
                    break;
                case 'zip':
                case 'rar':
                case '7z':
                case 'tar':
                case 'gz':
                    iconClass = 'bx-archive';
                    break;
                case 'txt':
                case 'md':
                case 'html':
                case 'css':
                case 'js':
                    iconClass = 'bx-file-txt';
                    break;
            }
        }

        return iconClass;
    }

    /**
     * Format file size in human-readable format
     * @param {number} bytes - File size in bytes
     * @returns {string} - Formatted file size
     */
    function formatFileSize(bytes) {
        if (bytes === 0 || bytes === undefined || bytes === null) return '0 Bytes';

        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));

        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    /**
     * Get filename from a URL
     * @param {string} url - URL containing a filename
     * @returns {string} - Extracted filename
     */
    function getFileName(url) {
        if (!url) return '';

        const urlWithoutParams = url.split('?')[0];

        const parts = urlWithoutParams.split('/');
        return parts[parts.length - 1] || '';
    }

    return {
        getHealth,
        hideText,
        hideFile,
        extract,
        analyze,
        calculateCapacity,

        processHideTextResponse,
        processHideFileResponse,
        processExtractResponse,
        processAnalyzeResponse,

        getFileIconClass,
        formatFileSize,
        getFileName
    };
})();

const UIManager = (function() {
    // Keep track of active toasts to prevent duplicates
    const activeToasts = new Set();

    // Get DOM elements
    const loadingOverlay = document.getElementById('loading-overlay');
    const toastContainer = document.getElementById('toast-container');

    function showLoading() {
        if (loadingOverlay) {
            loadingOverlay.style.display = 'flex';
        }
    }

    function hideLoading() {
        if (loadingOverlay) {
            setTimeout(() => {
                loadingOverlay.style.display = 'none';
            }, 300);
        }
    }

    /**
     * Show a toast notification
     * @param {string} message - Message to display
     * @param {string} type - Notification type (success, error, warning, info)
     * @param {string} title - Notification title
     * @param {number} duration - Duration in milliseconds
     */
    function showToast(message, type = 'info', title = '', duration = 5000) {
        if (!toastContainer) {
            console.error('Toast container not found!');
            return;
        }

        const toastKey = `${type}-${title}-${message}`;
        if (activeToasts.has(toastKey)) {
            console.log('Preventing duplicate toast:', toastKey);
            return;
        }

        activeToasts.add(toastKey);

        if (!title) {
            switch (type) {
                case 'success': title = 'Success'; break;
                case 'error': title = 'Error'; break;
                case 'warning': title = 'Warning'; break;
                case 'info': title = 'Information'; break;
            }
        }

        let icon = 'bx-info-circle';
        switch (type) {
            case 'success': icon = 'bx-check-circle'; break;
            case 'error': icon = 'bx-error-circle'; break;
            case 'warning': icon = 'bx-error'; break;
        }

        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        toast.innerHTML = `
            <div class="toast-icon">
                <i class='bx ${icon}'></i>
            </div>
            <div class="toast-content">
                <div class="toast-title">${title}</div>
                <p class="toast-message">${message}</p>
            </div>
            <button class="toast-close" aria-label="Close notification">
                <i class='bx bx-x'></i>
            </button>
        `;

        toastContainer.appendChild(toast);
        setTimeout(() => {
            toast.classList.add('showing');
        }, 10);

        const closeButton = toast.querySelector('.toast-close');
        closeButton.addEventListener('click', () => {
            removeToast(toast, toastKey);
        });

        setTimeout(() => {
            removeToast(toast, toastKey);
        }, duration);
    }

    /**
     * Remove a toast notification
     * @param {HTMLElement} toast - Toast element to remove
     * @param {string} toastKey - Key to remove from active toasts
     */
    function removeToast(toast, toastKey) {
        toast.classList.add('removing');
        toast.classList.remove('showing');
        activeToasts.delete(toastKey);

        // Remove from DOM after animation completes
        setTimeout(() => {
            if (toast.parentNode) {
                toast.parentNode.removeChild(toast);
            }
        }, 300);
    }
    return {
        showLoading,
        hideLoading,
        showToast
    };
})();
