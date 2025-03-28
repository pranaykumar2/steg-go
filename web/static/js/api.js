/**
 * StegGo - Image Steganography Tool
 * API Client
 * Author: pranaykumar2
 * Version: 1.0.0
 * Date: 2025-03-28 15:32:49 UTC
 * User: pranaykumar2
 */

class StegGoAPI {
    constructor(baseUrl = '/api') {
        this.baseUrl = baseUrl;
        this.lastRequest = null;
        this.pendingRequests = 0;
    }

    /**
     * Get API health status
     * @returns {Promise<Object>} Health status information
     */
    async getHealth() {
        return this.request('/health');
    }

    /**
     * Hide text in an image
     * @param {File} image - Cover image file
     * @param {string} message - Secret message to hide
     * @returns {Promise<Object>} Result with key and output file
     */
    async hideText(image, message) {
        if (!image) {
            throw new Error('No image selected.');
        }

        if (!message || message.trim() === '') {
            throw new Error('Message cannot be empty.');
        }

        const formData = new FormData();
        formData.append('image', image);
        formData.append('Message', message);

        console.log(`Hiding text in ${image.name}. Message length: ${message.length} characters`);

        return this.request('/hide', 'POST', formData);
    }

    /**
     * Hide a file in an image
     * @param {File} image - Cover image file
     * @param {File} file - File to hide
     * @returns {Promise<Object>} Result with key and output file
     */
    async hideFile(image, file) {
        if (!image) {
            throw new Error('No cover image selected.');
        }

        if (!file) {
            throw new Error('No file selected to hide.');
        }

        const formData = new FormData();
        formData.append('image', image);
        formData.append('file', file);

        console.log(`Hiding file ${file.name} (${this.formatFileSize(file.size)}) in ${image.name}`);

        return this.request('/hideFile', 'POST', formData);
    }

    /**
     * Extract hidden content from an image
     * @param {File} image - Stego image file
     * @param {string} key - Decryption key
     * @returns {Promise<Object>} Extracted content
     */
    async extract(image, key) {
        if (!image) {
            throw new Error('No image selected.');
        }

        if (!key || key.trim() === '') {
            throw new Error('Encryption key is required.');
        }

        // Validate key format (64 hex characters)
        if (!/^[0-9a-fA-F]{64}$/.test(key)) {
            throw new Error('Invalid key format. Key must be 64 hexadecimal characters.');
        }

        const formData = new FormData();
        formData.append('image', image);
        formData.append('Key', key); // Capital K to match Go struct field name

        console.log(`Extracting from ${image.name} using key: ${key.substring(0, 8)}...`);

        return this.request('/extract', 'POST', formData);
    }

    /**
     * Analyze image metadata
     * @param {File} image - Image file to analyze
     * @returns {Promise<Object>} Image analysis results
     */
    async analyzeMetadata(image) {
        if (!image) {
            throw new Error('No image selected.');
        }

        const formData = new FormData();
        formData.append('image', image);

        console.log(`Analyzing metadata for ${image.name}`);

        return this.request('/metadata', 'POST', formData);
    }

    /**
     * Make an API request with error handling
     * @param {string} endpoint - API endpoint
     * @param {string} method - HTTP method
     * @param {FormData|Object|null} data - Request data
     * @returns {Promise<Object>} Response data
     */
    async request(endpoint, method = 'GET', data = null) {
        // Increment pending requests counter
        this.pendingRequests++;

        // Dispatch event to notify about the request starting
        if (this.pendingRequests === 1) {
            document.dispatchEvent(new CustomEvent('api:loading:start'));
        }

        try {
            const url = `${this.baseUrl}${endpoint}`;
            const options = {
                method,
                headers: {}
            };

            if (data) {
                if (data instanceof FormData) {
                    options.body = data;
                } else {
                    options.headers['Content-Type'] = 'application/json';
                    options.body = JSON.stringify(data);
                }
            }

            // Cache the request
            this.lastRequest = { url, options };

            // Make the request
            const response = await fetch(url, options);
            const responseData = await response.json();

            // Handle API error (non-200 response)
            if (!response.ok) {
                throw new Error(responseData.error || `API error (${response.status}): ${response.statusText}`);
            }

            // Check API-specific error format
            if (responseData.success === false && responseData.error) {
                throw new Error(responseData.error);
            }

            // Return the data
            return responseData;
        } catch (error) {
            // Log the error
            console.error('API request failed:', error);

            // Rethrow to let the caller handle it
            throw error;
        } finally {
            // Decrement pending requests counter
            this.pendingRequests--;

            // Notify about all requests being complete
            if (this.pendingRequests === 0) {
                document.dispatchEvent(new CustomEvent('api:loading:end'));
            }
        }
    }

    /**
     * Get a full URL to a file
     * @param {string} path - File path
     * @returns {string} Full URL
     */
    getFileUrl(path) {
        // Handle already absolute URLs
        if (path.startsWith('http') || path.startsWith('/')) {
            return path;
        }

        return `${this.baseUrl}/files/${path}`;
    }

    /**
     * Format file size into human-readable string
     * @param {number} bytes - Size in bytes
     * @param {number} decimals - Decimal places
     * @returns {string} Formatted size
     */
    formatFileSize(bytes, decimals = 2) {
        if (bytes === 0) return '0 Bytes';

        const k = 1024;
        const dm = decimals < 0 ? 0 : decimals;
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];

        const i = Math.floor(Math.log(bytes) / Math.log(k));

        return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    }

    /**
     * Get the file name from a URL path
     * @param {string} url - URL or path
     * @returns {string} File name
     */
    static getFileName(url) {
        return url.split('/').pop();
    }

    /**
     * Get appropriate file icon class based on file type
     * @param {string} filename - File name
     * @returns {string} Font Awesome icon class
     */
    static getFileIconClass(filename) {
        const ext = filename.split('.').pop().toLowerCase();

        const iconMap = {
            // Documents
            'pdf': 'fa-file-pdf',
            'doc': 'fa-file-word',
            'docx': 'fa-file-word',
            'xls': 'fa-file-excel',
            'xlsx': 'fa-file-excel',
            'ppt': 'fa-file-powerpoint',
            'pptx': 'fa-file-powerpoint',
            'txt': 'fa-file-alt',

            // Images
            'jpg': 'fa-file-image',
            'jpeg': 'fa-file-image',
            'png': 'fa-file-image',
            'gif': 'fa-file-image',
            'svg': 'fa-file-image',

            // Archives
            'zip': 'fa-file-archive',
            'rar': 'fa-file-archive',
            'tar': 'fa-file-archive',
            'gz': 'fa-file-archive',

            // Code
            'html': 'fa-file-code',
            'css': 'fa-file-code',
            'js': 'fa-file-code',
            'json': 'fa-file-code',
            'xml': 'fa-file-code',

            // Audio
            'mp3': 'fa-file-audio',
            'wav': 'fa-file-audio',
            'ogg': 'fa-file-audio',

            // Video
            'mp4': 'fa-file-video',
            'avi': 'fa-file-video',
            'mov': 'fa-file-video',
            'wmv': 'fa-file-video',
        };

        return iconMap[ext] || 'fa-file';
    }
}

// Create a global instance
const stegApi = new StegGoAPI();

// Add loading event handlers to show/hide loading overlay
document.addEventListener('api:loading:start', () => {
    const loadingOverlay = document.getElementById('loading-overlay');
    if (loadingOverlay) loadingOverlay.style.display = 'flex';
});

document.addEventListener('api:loading:end', () => {
    const loadingOverlay = document.getElementById('loading-overlay');
    if (loadingOverlay) loadingOverlay.style.display = 'none';
});
