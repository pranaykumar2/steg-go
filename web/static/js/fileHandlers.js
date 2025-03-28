/**
 * StegGo - Image Steganography Tool
 * File Handling Module
 * Author: pranaykumar2
 * Version: 1.0.0
 * Date: 2025-03-28 15:51:28
 * User: pranaykumar2
 */

class FileHandler {
    constructor() {
        // Store uploaded files
        this.files = {};

        // Image capacity data for text hiding
        this.capacityData = {
            hideText: {
                capacity: 0,
                used: 0
            }
        };

        // Maximum file size (50MB)
        this.maxFileSize = 50 * 1024 * 1024;

        // Define section mapping based on element IDs
        this.sectionMapping = {
            // Upload areas to sections
            'hide-text-upload-area': 'hide-text',
            'hide-file-upload-area': 'hide-file',
            'hide-file-file-area': 'hide-file', // For file to hide
            'extract-upload-area': 'extract',
            'analyze-upload-area': 'analyze',

            // File inputs to sections
            'hide-text-file': 'hide-text',
            'hide-file-image': 'hide-file',
            'hide-file-file': 'hide-file',
            'extract-file': 'extract',
            'analyze-file': 'analyze'
        };

        // Define type mapping for file inputs
        this.typeMapping = {
            'hide-file-file': 'file', // This is for the file to hide, not the cover image
            'hide-text-file': 'image',
            'hide-file-image': 'image',
            'extract-file': 'image',
            'analyze-file': 'image'
        };
    }

    /**
     * Initialize file handler
     */
    initialize() {
        console.log('Initializing FileHandler...');
        this.setupUploadAreas();
        this.setupFileInputs();
        console.log('FileHandler initialization complete');
    }

    /**
     * Set up file upload areas
     */
    setupUploadAreas() {
        const uploadAreas = document.querySelectorAll('.upload-area');
        console.log(`Found ${uploadAreas.length} upload areas`);

        uploadAreas.forEach(area => {
            // Get section from ID
            const id = area.id;
            const section = this.sectionMapping[id];

            if (!section) {
                console.warn(`Upload area with unknown ID: ${id}`);
                return;
            }

            // Determine file type (image or file)
            let type = 'image';
            if (id === 'hide-file-file-area') {
                type = 'file';
            }

            console.log(`Setting up upload area for section: ${section}, type: ${type}`);

            // Setup drag and drop
            area.addEventListener('dragover', e => {
                e.preventDefault();
                e.stopPropagation();
                area.classList.add('dragover');
            });

            area.addEventListener('dragleave', e => {
                e.preventDefault();
                e.stopPropagation();
                area.classList.remove('dragover');
            });

            area.addEventListener('drop', e => {
                e.preventDefault();
                e.stopPropagation();
                area.classList.remove('dragover');

                if (e.dataTransfer.files.length > 0) {
                    const file = e.dataTransfer.files[0];
                    console.log(`File dropped: ${file.name} for ${section}`);
                    this.processFile(file, section, type);
                }
            });

            // Setup click to select file
            area.addEventListener('click', () => {
                // Find corresponding file input
                let inputId;
                if (section === 'hide-file' && type === 'file') {
                    inputId = 'hide-file-file';
                } else if (section === 'hide-file') {
                    inputId = 'hide-file-image';
                } else {
                    inputId = `${section}-file`;
                }

                const fileInput = document.getElementById(inputId);
                if (fileInput) {
                    console.log(`Clicking file input: ${inputId}`);
                    fileInput.click();
                } else {
                    console.warn(`File input not found: ${inputId}`);
                }
            });
        });
    }

    /**
     * Set up file input elements
     */
    setupFileInputs() {
        const fileInputs = document.querySelectorAll('input[type="file"]');
        console.log(`Found ${fileInputs.length} file inputs`);

        fileInputs.forEach(input => {
            const id = input.id;
            const section = this.sectionMapping[id];

            if (!section) {
                console.warn(`File input with unknown ID: ${id}`);
                return;
            }

            // Get the type (image or file)
            const type = this.typeMapping[id] || 'image';

            console.log(`Setting up file input ${id} for section: ${section}, type: ${type}`);

            input.addEventListener('change', e => {
                if (e.target.files.length > 0) {
                    const file = e.target.files[0];
                    console.log(`File selected: ${file.name} for ${section}`);
                    this.processFile(file, section, type);
                }
            });
        });
    }

    /**
     * Process uploaded file
     * @param {File} file - Uploaded file
     * @param {string} section - Section ID
     * @param {string} type - File type (image or file)
     */
    processFile(file, section, type) {
        console.log(`Processing file: ${file.name} for section: ${section}, type: ${type}`);

        // Check file size
        if (file.size > this.maxFileSize) {
            showToast('error', 'Error', `File too large. Maximum size is ${this.maxFileSize / (1024 * 1024)}MB.`);
            return;
        }

        // Process based on type
        if (type === 'image') {
            this.handleImageUpload(file, section);
        } else {
            this.handleFileUpload(file, section);
        }
    }

    /**
     * Handle image file upload
     * @param {File} file - Uploaded image
     * @param {string} section - Section ID
     */
    handleImageUpload(file, section) {
        console.log(`Handling image upload: ${file.name} for section: ${section}`);

        // Check if it's an image
        if (!file.type.startsWith('image/')) {
            showToast('error', 'Error', 'Please select a valid image file.');
            return;
        }

        // Store the file
        this.files[section] = file;

        // Create preview
        const reader = new FileReader();

        reader.onload = e => {
            // Find preview container based on section
            const previewContainerId = `${section}-preview-container`;
            let previewContainer = document.getElementById(previewContainerId);

            if (!previewContainer) {
                console.log(`Creating preview container: ${previewContainerId}`);

                // Get upload area to replace
                const uploadAreaId = section === 'hide-file' ? 'hide-file-upload-area' : `${section}-upload-area`;
                const uploadArea = document.getElementById(uploadAreaId);

                if (!uploadArea) {
                    console.error(`Upload area not found: ${uploadAreaId}`);
                    return;
                }

                // Create preview container
                previewContainer = document.createElement('div');
                previewContainer.id = previewContainerId;
                previewContainer.className = 'preview-container';

                // Replace upload area with preview container
                uploadArea.parentNode.replaceChild(previewContainer, uploadArea);
            }

            // Clear previous preview
            previewContainer.innerHTML = '';

            // Create image preview
            const img = document.createElement('img');
            img.src = e.target.result;
            img.alt = 'Preview';
            img.className = 'preview-image';

            // Create info div
            const info = document.createElement('div');
            info.className = 'preview-info';
            info.innerHTML = `
        <span class="file-name">${file.name}</span>
        <span class="file-size">${stegApi.formatFileSize(file.size)}</span>
        <button class="remove-file" data-section="${section}">
          <i class="fas fa-times"></i>
        </button>
      `;

            // Add to container
            previewContainer.appendChild(img);
            previewContainer.appendChild(info);

            // Setup remove button
            const removeBtn = previewContainer.querySelector('.remove-file');
            if (removeBtn) {
                removeBtn.addEventListener('click', () => {
                    this.removeFile(section, 'image');
                });
            }

            // Calculate capacity if in hide-text section
            if (section === 'hide-text') {
                this.calculateImageCapacity(img, section);
            }

            // Dispatch event for validation
            document.dispatchEvent(new CustomEvent('file:uploaded', {
                detail: {
                    type: 'image',
                    section: section,
                    file: file
                }
            }));

            console.log(`Image preview created for ${file.name} in ${section}`);
        };

        reader.onerror = () => {
            console.error(`Error reading file: ${file.name}`);
            showToast('error', 'Error', 'Failed to read the image file.');
        };

        reader.readAsDataURL(file);
    }

    /**
     * Handle regular file upload (for hiding)
     * @param {File} file - Uploaded file
     * @param {string} section - Section ID
     */
    handleFileUpload(file, section) {
        console.log(`Handling file upload: ${file.name} for section: ${section}`);

        // Store the file with special key
        this.files[`${section}-file`] = file;

        // Find or create preview container
        const previewContainerId = `${section}-file-preview-container`;
        let previewContainer = document.getElementById(previewContainerId);

        if (!previewContainer) {
            console.log(`Creating file preview container: ${previewContainerId}`);

            // Get upload area to replace
            const uploadArea = document.getElementById('hide-file-file-area');

            if (!uploadArea) {
                console.error('File upload area not found: hide-file-file-area');
                return;
            }

            // Create preview container
            previewContainer = document.createElement('div');
            previewContainer.id = previewContainerId;
            previewContainer.className = 'file-preview-container';

            // Replace upload area with preview container
            uploadArea.parentNode.replaceChild(previewContainer, uploadArea);
        }

        // Clear container
        previewContainer.innerHTML = '';

        // Get file icon based on extension
        const fileExt = file.name.split('.').pop().toLowerCase();
        const iconClass = StegGoAPI.getFileIconClass(file.name);

        // Create file preview
        const filePreview = document.createElement('div');
        filePreview.className = 'file-preview';
        filePreview.innerHTML = `
      <div class="file-icon">
        <i class="fas ${iconClass}"></i>
      </div>
      <div class="file-details">
        <span class="file-name">${file.name}</span>
        <span class="file-size">${stegApi.formatFileSize(file.size)}</span>
      </div>
      <button class="remove-file" data-section="${section}">
        <i class="fas fa-times"></i>
      </button>
    `;

        // Add to container
        previewContainer.appendChild(filePreview);

        // Setup remove button
        const removeBtn = previewContainer.querySelector('.remove-file');
        if (removeBtn) {
            removeBtn.addEventListener('click', () => {
                this.removeFile(section, 'file');
            });
        }

        // Dispatch event for validation
        document.dispatchEvent(new CustomEvent('file:uploaded', {
            detail: {
                type: 'file',
                section: section,
                file: file
            }
        }));

        console.log(`File preview created for ${file.name} in ${section}`);
    }

    /**
     * Remove a file
     * @param {string} section - Section ID
     * @param {string} type - File type (image or file)
     */
    removeFile(section, type) {
        console.log(`Removing ${type} from ${section}`);

        if (type === 'image') {
            // Remove from storage
            delete this.files[section];

            // Get preview container
            const previewContainer = document.getElementById(`${section}-preview-container`);
            if (!previewContainer) return;

            // Restore upload area
            const uploadAreaId = section === 'hide-file' ? 'hide-file-upload-area' : `${section}-upload-area`;
            let uploadArea = document.createElement('div');
            uploadArea.id = uploadAreaId;
            uploadArea.className = 'upload-area';

            // Add appropriate content
            uploadArea.innerHTML = `
        <div class="upload-icon">
          <i class="fas fa-cloud-upload-alt"></i>
        </div>
        <p>Drag & drop an image here, or click to select</p>
        <p class="upload-hint">PNG, JPG up to 50MB</p>
      `;

            // Replace preview with upload area
            previewContainer.parentNode.replaceChild(uploadArea, previewContainer);

            // Re-setup the upload area
            this.setupUploadAreas();

            // Reset capacity if in hide-text section
            if (section === 'hide-text') {
                this.capacityData.hideText = {
                    capacity: 0,
                    used: 0
                };
                this.updateCapacityIndicator('hide-text');
            }
        } else if (type === 'file') {
            // Remove from storage
            delete this.files[`${section}-file`];

            // Get preview container
            const filePreviewContainer = document.getElementById(`${section}-file-preview-container`);
            if (!filePreviewContainer) return;

            // Restore file upload area
            const fileArea = document.createElement('div');
            fileArea.id = 'hide-file-file-area';
            fileArea.className = 'upload-area file-upload-area';

            // Add content
            fileArea.innerHTML = `
        <div class="upload-icon">
          <i class="fas fa-file-upload"></i>
        </div>
        <p>Drag & drop a file here, or click to select</p>
        <p class="upload-hint">Any file type up to 50MB</p>
      `;

            // Replace preview with upload area
            filePreviewContainer.parentNode.replaceChild(fileArea, filePreviewContainer);

            // Re-setup the upload area
            this.setupUploadAreas();
        }

        // Dispatch event for validation
        document.dispatchEvent(new CustomEvent('file:removed', {
            detail: {
                type: type,
                section: section
            }
        }));
    }

    /**
     * Calculate image capacity for steganography
     * @param {HTMLImageElement} img - Image element
     * @param {string} section - Section ID
     */
    calculateImageCapacity(img, section) {
        if (section !== 'hide-text') return;

        img.onload = () => {
            // Simple capacity estimation: 3 bits per pixel (1 bit per color channel)
            const width = img.naturalWidth;
            const height = img.naturalHeight;
            const pixelCount = width * height;

            // Calculate capacity in bytes (3 bits per pixel / 8 bits per byte)
            const capacityBytes = Math.floor((pixelCount * 3) / 8);

            console.log(`Image dimensions: ${width}x${height}, capacity: ${capacityBytes} bytes`);

            // Store capacity
            this.capacityData.hideText.capacity = capacityBytes;

            // Update capacity indicator
            this.updateCapacityIndicator(section);

            // Set up message field listener
            const messageField = document.getElementById('hide-text-message');
            if (messageField) {
                // Remove any existing listeners
                const newField = messageField.cloneNode(true);
                messageField.parentNode.replaceChild(newField, messageField);

                // Add new listener
                newField.addEventListener('input', () => {
                    this.capacityData.hideText.used = newField.value.length;
                    this.updateCapacityIndicator(section);
                });
            }
        };
    }

    /**
     * Update capacity indicator UI
     * @param {string} section - Section ID
     */
    updateCapacityIndicator(section) {
        if (section !== 'hide-text') return;

        const indicator = document.getElementById('capacity-indicator');
        if (!indicator) return;

        const barElement = document.getElementById('capacity-bar');
        const textElement = document.getElementById('capacity-text');

        if (!barElement || !textElement) return;

        const { capacity, used } = this.capacityData.hideText;

        // Hide if no capacity
        if (capacity === 0) {
            indicator.style.display = 'none';
            return;
        }

        // Show indicator
        indicator.style.display = 'block';

        // Calculate percentage
        const percentage = Math.min(100, Math.floor((used / capacity) * 100));

        // Update bar
        barElement.style.width = `${percentage}%`;
        textElement.textContent = `${used} / ${capacity} bytes (${percentage}%)`;

        // Set color based on usage
        if (percentage < 70) {
            barElement.className = 'capacity-bar-fill green';
        } else if (percentage < 90) {
            barElement.className = 'capacity-bar-fill orange';
        } else {
            barElement.className = 'capacity-bar-fill red';
        }
    }

    /**
     * Get a stored file
     * @param {string} section - Section ID
     * @param {string} type - File type (image or file)
     * @returns {File|null} The stored file
     */
    getFile(section, type = 'image') {
        const key = type === 'image' ? section : `${section}-file`;
        return this.files[key] || null;
    }
}

// Create instance
const fileHandler = new FileHandler();
