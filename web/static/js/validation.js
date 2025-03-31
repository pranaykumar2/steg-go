class FormValidator {
    constructor() {
        this.forms = {
            'hide-text': { valid: false },
            'hide-file': { valid: false },
            'extract': { valid: false },
            'analyze': { valid: false }
        };

        this.fields = {
            'hide-text-message': { valid: false, errorMessage: 'Please enter a message to hide.' },
            'extract-key': { valid: false, errorMessage: 'Please enter a valid 64-character encryption key.' }
        };

        this.files = {
            'hide-text': { valid: false, errorMessage: 'Please select an image file.' },
            'hide-file': { valid: false, errorMessage: 'Please select a cover image.' },
            'hide-file-file': { valid: false, errorMessage: 'Please select a file to hide.' },
            'extract': { valid: false, errorMessage: 'Please select an image file.' },
            'analyze': { valid: false, errorMessage: 'Please select an image file to analyze.' }
        };
    }

    initialize() {
        console.log('Initializing FormValidator...');

        this.setupResultContainers();

        this.setupEventListeners();

        this.validateAllForms();

        this.setupDirectInputMonitoring();

        console.log('FormValidator initialization complete');
    }

    setupResultContainers() {
        this.createHideFileResults();

        this.createExtractResults();

        this.createAnalyzeResults();
    }

    createHideFileResults() {
        let resultsSection = document.getElementById('hide-file-results');
        if (!resultsSection) {
            console.log('Creating hide-file-results section');
            const hideFileForm = document.getElementById('hide-file-form');
            if (!hideFileForm) {
                console.warn('Unable to create hide-file-results: hide-file-form not found');
                return;
            }

            resultsSection = document.createElement('div');
            resultsSection.id = 'hide-file-results';
            resultsSection.className = 'results-section card mt-4';
            resultsSection.style.display = 'none';

            resultsSection.innerHTML = `
        <div class="card-header">
          <h3>File Hidden Successfully</h3>
        </div>
        <div class="card-body">
          <div class="row">
            <div class="col-md-6">
              <div class="result-image-container">
                <img id="hide-file-result-img" alt="Steganographic image with hidden file" class="result-image">
              </div>
            </div>
            <div class="col-md-6">
              <div class="result-details">
                <div class="form-group">
                  <label for="hide-file-key">Decryption Key:</label>
                  <div class="input-group">
                    <input type="text" id="hide-file-key" class="form-control" readonly>
                    <div class="input-group-append">
                      <button class="btn btn-outline-secondary" type="button" onclick="copyToClipboard('hide-file-key')">
                        <i class="fas fa-copy"></i>
                      </button>
                    </div>
                  </div>
                  <small class="form-text text-muted">Save this key to extract the file later.</small>
                </div>
                <div id="hide-file-details" class="details-container"></div>
                <div class="mt-3">
                  <a id="hide-file-download" href="#" class="btn btn-primary" download="stego-image.png">
                    <i class="fas fa-download"></i> Download Image
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      `;

            hideFileForm.parentNode.insertBefore(resultsSection, hideFileForm.nextSibling);
        }
    }

    createExtractResults() {
        let resultsSection = document.getElementById('extract-results');
        if (!resultsSection) {
            console.log('Creating extract-results section');
            const extractForm = document.getElementById('extract-form');
            if (!extractForm) {
                console.warn('Unable to create extract-results: extract-form not found');
                return;
            }

            resultsSection = document.createElement('div');
            resultsSection.id = 'extract-results';
            resultsSection.className = 'results-section card mt-4';
            resultsSection.style.display = 'none';

            resultsSection.innerHTML = `
        <div class="card-header">
          <h3>Content Extracted Successfully</h3>
        </div>
        <div class="card-body">
          <div id="extract-text-results" style="display: none;">
            <h4>Extracted Message:</h4>
            <div class="extracted-text">
              <pre id="extract-text-content"></pre>
            </div>
            <button class="btn btn-primary mt-3" onclick="copyToClipboard('extract-text-content')">
              <i class="fas fa-copy"></i> Copy Text
            </button>
          </div>
          
          <div id="extract-file-results" style="display: none;">
            <h4>Extracted File:</h4>
            <div class="file-details">
              <div class="detail-item">
                <span class="detail-label">File Name:</span>
                <span id="extract-file-name" class="detail-value">-</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">File Size:</span>
                <span id="extract-file-size" class="detail-value">-</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">File Type:</span>
                <span id="extract-file-type" class="detail-value">-</span>
              </div>
              <div class="mt-3">
                <a id="extract-file-download" href="#" class="btn btn-primary" download>
                  <i class="fas fa-download"></i> Download File
                </a>
              </div>
            </div>
          </div>
        </div>
      `;
            extractForm.parentNode.insertBefore(resultsSection, extractForm.nextSibling);
        }
    }

    createAnalyzeResults() {
        let resultsSection = document.getElementById('analyze-results');

        if (!resultsSection) {
            console.log('Creating analyze-results section');

            const analyzeForm = document.getElementById('analyze-form');
            if (!analyzeForm) {
                console.warn('Unable to create analyze-results: analyze-form not found');
                return;
            }

            resultsSection = document.createElement('div');
            resultsSection.id = 'analyze-results';
            resultsSection.className = 'results-section card mt-4';
            resultsSection.style.display = 'none';

            resultsSection.innerHTML = `
        <div class="card-header">
          <h3>Image Analysis Results</h3>
        </div>
        <div class="card-body">
          <div class="row">
            <div class="col-md-6">
              <div class="analysis-section">
                <h4>Image Details</h4>
                <div id="analyze-image-details" class="details-container"></div>
              </div>
              
              <div class="analysis-section mt-4">
                <h4>Steganographic Capacity</h4>
                <div id="analyze-capacity-details" class="details-container"></div>
              </div>
            </div>
            
            <div class="col-md-6">
              <div id="analyze-privacy-risks" class="analysis-section"></div>
              
              <div id="analyze-exif-data" class="analysis-section mt-4"></div>
            </div>
          </div>
        </div>
      `;
            analyzeForm.parentNode.insertBefore(resultsSection, analyzeForm.nextSibling);
        }
    }

    setupDirectInputMonitoring() {
        setInterval(() => {
            const messageField = document.getElementById('hide-text-message');
            if (messageField) {
                const value = messageField.value.trim();
                const currentValidState = this.fields['hide-text-message'].valid;
                const newValidState = value.length > 0;
                if (currentValidState !== newValidState) {
                    console.log(`Message field validity changed: ${currentValidState} -> ${newValidState}`);
                    this.fields['hide-text-message'].valid = newValidState;
                    this.validateForm('hide-text');
                }
            }
        }, 500);
    }

    setupEventListeners() {
        document.addEventListener('file:uploaded', (e) => {
            console.log('File uploaded event:', e.detail);
            const { type, section, file } = e.detail;

            const mappedSection = section === 'hide' ? 'hide-text' : section;

            if (type === 'image') {
                if (this.files[mappedSection]) {
                    console.log(`Setting ${mappedSection} image as valid`);
                    this.files[mappedSection].valid = true;
                } else {
                    console.warn(`Unknown section for validation: ${mappedSection}`);
                }
            } else if (type === 'file') {
                const fileKey = `${mappedSection}-file`;
                if (this.files[fileKey]) {
                    console.log(`Setting ${fileKey} as valid`);
                    this.files[fileKey].valid = true;
                } else {
                    console.warn(`Unknown section for file validation: ${fileKey}`);
                }
            }

            this.validateForm(mappedSection);
        });

        document.addEventListener('file:removed', (e) => {
            console.log('File removed event:', e.detail);
            const { type, section } = e.detail;

            const mappedSection = section === 'hide' ? 'hide-text' : section;

            if (type === 'image') {
                if (this.files[mappedSection]) {
                    this.files[mappedSection].valid = false;
                }
            } else if (type === 'file') {
                const fileKey = `${mappedSection}-file`;
                if (this.files[fileKey]) {
                    this.files[fileKey].valid = false;
                }
            }

            this.validateForm(mappedSection);
        });

        const messageField = document.getElementById('hide-text-message');
        if (messageField) {
            console.log('Setting up message field listener');

            ['input', 'change', 'keyup', 'paste'].forEach(eventType => {
                messageField.addEventListener(eventType, () => {
                    console.log(`Message ${eventType} event:`, messageField.value);
                    this.validateField('hide-text-message');
                    this.validateForm('hide-text');
                });
            });

            setTimeout(() => {
                if (messageField.value.trim()) {
                    console.log('Found existing message text, validating');
                    this.validateField('hide-text-message');
                    this.validateForm('hide-text');
                }
            }, 100);
        } else {
            console.warn('Message field not found!');
        }

        const keyField = document.getElementById('extract-key');
        if (keyField) {
            ['input', 'change', 'keyup', 'paste'].forEach(eventType => {
                keyField.addEventListener(eventType, () => {
                    this.validateField('extract-key');
                    this.validateForm('extract');
                });
            });
        }

        this.setupFormSubmissions();

        this.addDebugTools();

        if (typeof window.copyToClipboard !== 'function') {
            window.copyToClipboard = function(elementId) {
                const element = document.getElementById(elementId);
                if (!element) return;

                let text;
                if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
                    text = element.value;
                } else {
                    text = element.textContent;
                }

                navigator.clipboard.writeText(text).then(() => {
                    showToast('success', 'Copied!', 'Text copied to clipboard.');
                }, (err) => {
                    console.error('Error copying text: ', err);
                    showToast('error', 'Error', 'Failed to copy text.');
                });
            };
        }
    }

    addDebugTools() {
        if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
            document.addEventListener('keydown', e => {
                if (e.altKey && e.key === 'd') {
                    console.log('Form States:', this.forms);
                    console.log('Field States:', this.fields);
                    console.log('File States:', this.files);
                    this.highlightFormsDebug();
                }

                if (e.altKey && e.key === 'f') {
                    console.log('Force enabling all submit buttons');
                    document.querySelectorAll('button[type="submit"]').forEach(btn => {
                        btn.disabled = false;
                        btn.classList.add('debug-force-enabled');
                    });
                }

                if (e.altKey && e.key === 'r') {
                    console.log('Force showing all result sections');
                    document.querySelectorAll('[id$="-results"]').forEach(section => {
                        section.style.display = 'block';
                        console.log(`Showing result section: ${section.id}`);
                    });
                }
            });
        }
    }

    highlightFormsDebug() {
        Object.entries(this.forms).forEach(([id, { valid }]) => {
            const form = document.getElementById(`${id}-form`);
            if (form) {
                form.style.outline = valid ? '2px solid green' : '2px solid red';
                setTimeout(() => {
                    form.style.outline = '';
                }, 3000);
            }
        });
    }

    setupFormSubmissions() {
        const hideTextForm = document.getElementById('hide-text-form');
        if (hideTextForm) {
            hideTextForm.addEventListener('submit', (e) => {
                e.preventDefault();
                console.log('Hide Text form submitted, validity:', this.forms['hide-text'].valid);

                if (this.forms['hide-text'].valid) {
                    this.submitHideTextForm();
                } else {
                    console.warn('Form validation failed');
                    this.showFormErrors('hide-text');
                }
            });
        } else {
            console.warn('Hide Text form not found!');
        }

        const hideFileForm = document.getElementById('hide-file-form');
        if (hideFileForm) {
            hideFileForm.addEventListener('submit', (e) => {
                e.preventDefault();
                console.log('Hide File form submitted, validity:', this.forms['hide-file'].valid);

                if (this.forms['hide-file'].valid) {
                    this.submitHideFileForm();
                } else {
                    this.showFormErrors('hide-file');
                }
            });
        }

        const extractForm = document.getElementById('extract-form');
        if (extractForm) {
            extractForm.addEventListener('submit', (e) => {
                e.preventDefault();
                console.log('Extract form submitted, validity:', this.forms['extract'].valid);

                if (this.forms['extract'].valid) {
                    this.submitExtractForm();
                } else {
                    this.showFormErrors('extract');
                }
            });
        }

        const analyzeForm = document.getElementById('analyze-form');
        if (analyzeForm) {
            analyzeForm.addEventListener('submit', (e) => {
                e.preventDefault();
                console.log('Analyze form submitted, validity:', this.forms['analyze'].valid);

                if (this.forms['analyze'].valid) {
                    this.submitAnalyzeForm();
                } else {
                    this.showFormErrors('analyze');
                }
            });
        }
    }

    /**
     * Show form validation errors
     * @param {string} formId - Form ID
     */
    showFormErrors(formId) {
        let errorMessages = [];

        switch (formId) {
            case 'hide-text':
                if (!this.files['hide-text'].valid) {
                    errorMessages.push('Please select an image file.');
                }
                if (!this.fields['hide-text-message'].valid) {
                    errorMessages.push('Please enter a message to hide.');
                }
                break;

            case 'hide-file':
                if (!this.files['hide-file'].valid) {
                    errorMessages.push('Please select a cover image.');
                }
                if (!this.files['hide-file-file'].valid) {
                    errorMessages.push('Please select a file to hide.');
                }
                break;

            case 'extract':
                if (!this.files['extract'].valid) {
                    errorMessages.push('Please select an image file.');
                }
                if (!this.fields['extract-key'].valid) {
                    errorMessages.push('Please enter a valid encryption key.');
                }
                break;

            case 'analyze':
                if (!this.files['analyze'].valid) {
                    errorMessages.push('Please select an image file to analyze.');
                }
                break;
        }

        if (errorMessages.length > 0) {
            showToast('error', 'Form Error', errorMessages.join('<br>'));
        }
    }

    /**
     * Validate all forms
     */
    validateAllForms() {
        Object.keys(this.forms).forEach(formId => {
            this.validateForm(formId);
        });
    }

    /**
     * Validate a form
     * @param {string} formId - Form ID
     */
    validateForm(formId) {
        console.log(`Validating form: ${formId}`);
        let isValid = false;

        switch (formId) {
            case 'hide-text':
                const messageField = document.getElementById('hide-text-message');
                const messageText = messageField ? messageField.value.trim() : '';
                const messageValid = messageText.length > 0;

                this.fields['hide-text-message'].valid = messageValid;

                const imageValid = this.files['hide-text'].valid;

                console.log(`Hide text validation - Message: ${messageValid} (text: "${messageText.substring(0, 20)}${messageText.length > 20 ? '...' : ''}"), Image: ${imageValid}`);
                isValid = messageValid && imageValid;
                break;

            case 'hide-file':
                isValid = this.files['hide-file'].valid && this.files['hide-file-file'].valid;
                break;

            case 'extract':
                const keyField = document.getElementById('extract-key');
                const keyValid = keyField && /^[0-9a-fA-F]{64}$/.test(keyField.value.trim());
                this.fields['extract-key'].valid = keyValid;

                const extractImageValid = this.files['extract'].valid;
                isValid = keyValid && extractImageValid;
                break;

            case 'analyze':
                isValid = this.files['analyze'].valid;
                break;
        }

        this.forms[formId].valid = isValid;
        console.log(`Form ${formId} is now ${isValid ? 'valid' : 'invalid'}`);

        this.updateSubmitButton(formId);
    }

    /**
     * Update submit button state
     * @param {string} formId - Form ID
     */
    updateSubmitButton(formId) {
        const form = document.getElementById(`${formId}-form`);
        if (!form) {
            console.warn(`Form element not found: ${formId}-form`);
            return;
        }

        const submitButton = form.querySelector('button[type="submit"]');
        if (!submitButton) {
            console.warn(`Submit button not found in form: ${formId}-form`);
            return;
        }

        const isValid = this.forms[formId].valid;
        submitButton.disabled = !isValid;

        console.log(`Updated submit button for ${formId}: disabled=${!isValid}`);

        if (isValid) {
            submitButton.classList.remove('disabled');
            submitButton.classList.add('enabled');
        } else {
            submitButton.classList.add('disabled');
            submitButton.classList.remove('enabled');
        }
    }

    /**
     * Validate a specific field
     * @param {string} fieldId - Field ID
     */
    validateField(fieldId) {
        const field = document.getElementById(fieldId);
        if (!field) {
            console.warn(`Field not found: ${fieldId}`);
            return;
        }

        let isValid = false;

        switch (fieldId) {
            case 'hide-text-message':
                isValid = field.value.trim().length > 0;
                console.log(`Validating message field: ${isValid ? 'valid' : 'invalid'} (length: ${field.value.trim().length})`);
                break;

            case 'extract-key':
                isValid = /^[0-9a-fA-F]{64}$/.test(field.value.trim());
                break;
        }

        this.fields[fieldId].valid = isValid;

        field.classList.toggle('invalid', !isValid);
        field.classList.toggle('valid', isValid);
    }

    async submitHideTextForm() {
        try {
            const image = fileHandler.getFile('hide-text');
            const messageElement = document.getElementById('hide-text-message');
            const message = messageElement ? messageElement.value.trim() : '';

            if (!image) {
                showToast('error', 'Error', 'Please select an image file.');
                return;
            }

            if (!message) {
                showToast('error', 'Error', 'Please enter a message to hide.');
                return;
            }

            const submitButton = document.querySelector('#hide-text-form button[type="submit"]');
            if (submitButton) {
                submitButton.disabled = true;
                submitButton.innerHTML = '<i class="fas fa-circle-notch fa-spin"></i> Processing...';
            }

            console.log('Calling hideText API with:', { image, messageLength: message.length });
            const response = await stegApi.hideText(image, message);
            console.log('Hide text API response:', response);

            this.handleHideTextSuccess(response);
        } catch (error) {
            console.error('Error hiding text:', error);
            showToast('error', 'Error', error.message || 'Failed to hide text. Please try again.');
        } finally {
            const submitButton = document.querySelector('#hide-text-form button[type="submit"]');
            if (submitButton) {
                submitButton.disabled = !this.forms['hide-text'].valid;
                submitButton.innerHTML = 'Hide Message';
            }
        }
    }

    /**
     * Handle successful text hiding
     * @param {Object} response - API response
     */
    handleHideTextSuccess(response) {
        console.log('Hide text success response:', response);

        if (!response.data) {
            showToast('warning', 'Warning', 'Unexpected server response format');
            console.error('Invalid response format:', response);
            return;
        }

        const data = response.data;
        console.log('Processing hide text response data:', data);

        const keyInput = document.getElementById('hide-text-key');
        const resultImage = document.getElementById('hide-text-result-img');
        const downloadLink = document.getElementById('hide-text-download');

        if (keyInput) {
            keyInput.value = data.key || '';
            console.log('Set key input value:', data.key);
        } else {
            console.warn('Key input element not found: hide-text-key');
        }

        if (resultImage) {
            resultImage.src = data.outputFileURL || '';
            resultImage.alt = 'Steganographic image with hidden message';
            console.log('Set result image src:', data.outputFileURL);
        } else {
            console.warn('Result image element not found: hide-text-result-img');
        }

        if (downloadLink) {
            downloadLink.href = data.outputFileURL || '#';
            downloadLink.download = data.outputFileURL ? StegGoAPI.getFileName(data.outputFileURL) : 'stego-image.png';
            console.log('Set download link href:', data.outputFileURL);
        } else {
            console.warn('Download link element not found: hide-text-download');
        }

        const resultsSection = document.getElementById('hide-text-results');
        if (resultsSection) {
            resultsSection.style.display = 'block';
            console.log('Show hide-text-results section');

            setTimeout(() => {
                resultsSection.scrollIntoView({ behavior: 'smooth' });
            }, 100);
        } else {
            console.warn('Results section element not found: hide-text-results');
        }

        showToast('success', 'Success', 'Your message has been successfully hidden in the image.');
    }

    async submitHideFileForm() {
        try {
            const image = fileHandler.getFile('hide-file');
            const file = fileHandler.getFile('hide-file', 'file');

            if (!image) {
                showToast('error', 'Error', 'Please select a cover image.');
                return;
            }

            if (!file) {
                showToast('error', 'Error', 'Please select a file to hide.');
                return;
            }

            const submitButton = document.querySelector('#hide-file-form button[type="submit"]');
            if (submitButton) {
                submitButton.disabled = true;
                submitButton.innerHTML = '<i class="fas fa-circle-notch fa-spin"></i> Processing...';
            }

            console.log('Calling hideFile API with:', { image, file });
            const response = await stegApi.hideFile(image, file);
            console.log('Hide file API response:', response);

            this.handleHideFileSuccess(response);
        } catch (error) {
            console.error('Error hiding file:', error);
            showToast('error', 'Error', error.message || 'Failed to hide file. Please try again.');
        } finally {
            const submitButton = document.querySelector('#hide-file-form button[type="submit"]');
            if (submitButton) {
                submitButton.disabled = !this.forms['hide-file'].valid;
                submitButton.innerHTML = 'Hide File';
            }
        }
    }

    /**
     * Handle successful file hiding
     * @param {Object} response - API response
     */
    handleHideFileSuccess(response) {
        console.log('Hide file success response:', response);

        if (!response.data) {
            showToast('warning', 'Warning', 'Unexpected server response format');
            console.error('Invalid response format:', response);
            return;
        }

        const data = response.data;
        console.log('Processing hide file response data:', data);

        const keyInput = document.getElementById('hide-file-key');
        const resultImage = document.getElementById('hide-file-result-img');
        const downloadLink = document.getElementById('hide-file-download');
        const fileDetailsElement = document.getElementById('hide-file-details');

        if (keyInput) {
            keyInput.value = data.key || '';
            console.log('Set file key input value:', data.key);
        } else {
            console.warn('Key input element not found: hide-file-key');
        }

        if (resultImage) {
            resultImage.src = data.outputFileURL || '';
            resultImage.alt = 'Steganographic image with hidden file';
            console.log('Set file result image src:', data.outputFileURL);
        } else {
            console.warn('Result image element not found: hide-file-result-img');
        }

        if (downloadLink) {
            downloadLink.href = data.outputFileURL || '#';
            downloadLink.download = data.outputFileURL ? StegGoAPI.getFileName(data.outputFileURL) : 'stego-image.png';
            console.log('Set file download link href:', data.outputFileURL);
        } else {
            console.warn('Download link element not found: hide-file-download');
        }

        if (fileDetailsElement) {
            if (data.fileDetails) {
                const details = data.fileDetails;
                fileDetailsElement.innerHTML = `
          <div class="detail-item">
            <span class="detail-label">File Name:</span>
            <span class="detail-value">${details.originalName || 'Unknown'}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">File Type:</span>
            <span class="detail-value">${details.fileType || 'Unknown'}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">File Size:</span>
            <span class="detail-value">${stegApi.formatFileSize(details.fileSize || 0)}</span>
          </div>
        `;
                console.log('Set file details:', details);
            } else {
                const hiddenFile = fileHandler.getFile('hide-file', 'file');
                fileDetailsElement.innerHTML = `
          <div class="detail-item">
            <span class="detail-label">File Name:</span>
            <span class="detail-value">${hiddenFile ? hiddenFile.name : 'Unknown'}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">File Size:</span>
            <span class="detail-value">${hiddenFile ? stegApi.formatFileSize(hiddenFile.size) : 'Unknown'}</span>
          </div>
        `;
                console.log('Set fallback file details from original file');
            }
        } else {
            console.warn('File details element not found: hide-file-details');
        }

        const resultsSection = document.getElementById('hide-file-results');
        if (resultsSection) {
            console.log('Showing hide-file-results section');
            resultsSection.style.display = 'block';

            setTimeout(() => {
                resultsSection.scrollIntoView({ behavior: 'smooth' });
            }, 100);
        } else {
            console.error('Results section element not found: hide-file-results');
        }

        showToast('success', 'Success', 'Your file has been successfully hidden in the image.');
    }

    async submitExtractForm() {
        try {
            const image = fileHandler.getFile('extract');
            const keyElement = document.getElementById('extract-key');
            const key = keyElement ? keyElement.value.trim() : '';

            if (!image) {
                showToast('error', 'Error', 'Please select an image to extract from.');
                return;
            }

            if (!key) {
                showToast('error', 'Error', 'Please enter a decryption key.');
                return;
            }

            if (!/^[0-9a-fA-F]{64}$/.test(key)) {
                showToast('error', 'Error', 'Invalid key format. Key must be 64 hexadecimal characters.');
                return;
            }

            const submitButton = document.querySelector('#extract-form button[type="submit"]');
            if (submitButton) {
                submitButton.disabled = true;
                submitButton.innerHTML = '<i class="fas fa-circle-notch fa-spin"></i> Processing...';
            }
            console.log('Calling extract API with:', { image, key });
            const response = await stegApi.extract(image, key);
            console.log('Extract API response:', response);

            this.handleExtractSuccess(response);
        } catch (error) {
            console.error('Error extracting content:', error);
            showToast('error', 'Error', error.message || 'Failed to extract content. Please check your key and try again.');
        } finally {
            const submitButton = document.querySelector('#extract-form button[type="submit"]');
            if (submitButton) {
                submitButton.disabled = !this.forms['extract'].valid;
                submitButton.innerHTML = 'Extract Content';
            }
        }
    }

    /**
     * Handle successful extraction
     * @param {Object} response - API response
     */
    handleExtractSuccess(response) {
        console.log('Extract success response:', response);

        if (!response.data) {
            showToast('warning', 'Warning', 'Unexpected server response format');
            console.error('Invalid response format:', response);
            return;
        }

        const data = response.data;
        console.log('Processing extract response data:', data);

        const resultsContainer = document.getElementById('extract-results');
        if (!resultsContainer) {
            console.error('Results container not found: extract-results');
            return;
        }

        const textResultsContainer = document.getElementById('extract-text-results');
        const fileResultsContainer = document.getElementById('extract-file-results');

        if (!data.isFile && data.message) {
            console.log('Extracted text content:', data.message);

            const textContent = document.getElementById('extract-text-content');
            if (textContent) {
                textContent.textContent = data.message;
                console.log('Set extract-text-content');
            } else {
                console.warn('Text content element not found: extract-text-content');
            }

            if (textResultsContainer) {
                textResultsContainer.style.display = 'block';
                console.log('Show extract-text-results section');
            }
            if (fileResultsContainer) fileResultsContainer.style.display = 'none';

        } else if (data.isFile) {
            console.log('Extracted file content');

            const fileNameElement = document.getElementById('extract-file-name');
            const fileSizeElement = document.getElementById('extract-file-size');
            const fileTypeElement = document.getElementById('extract-file-type');
            const fileDownloadLink = document.getElementById('extract-file-download');

            if (fileNameElement) {
                fileNameElement.textContent = data.fileName || 'Unknown';
                console.log('Set extract-file-name:', data.fileName);
            }

            if (fileSizeElement) {
                fileSizeElement.textContent = stegApi.formatFileSize(data.fileSize || 0);
                console.log('Set extract-file-size:', data.fileSize);
            }

            if (fileTypeElement) {
                fileTypeElement.textContent = data.fileType || 'Unknown';
                console.log('Set extract-file-type:', data.fileType);
            }

            if (fileDownloadLink && data.fileURL) {
                fileDownloadLink.href = data.fileURL;
                fileDownloadLink.download = data.fileName || 'extracted_file';
                console.log('Set extract-file-download href:', data.fileURL);
            }

            if (fileResultsContainer) {
                fileResultsContainer.style.display = 'block';
                console.log('Show extract-file-results section');
            }
            if (textResultsContainer) textResultsContainer.style.display = 'none';
        }

        resultsContainer.style.display = 'block';
        console.log('Show extract-results section');

        setTimeout(() => {
            resultsContainer.scrollIntoView({ behavior: 'smooth' });
        }, 100);

        showToast('success', 'Success', 'Content successfully extracted!');
    }

    async submitAnalyzeForm() {
        try {
            const image = fileHandler.getFile('analyze');
            if (!image) {
                showToast('error', 'Error', 'Please select an image to analyze.');
                return;
            }
            const submitButton = document.querySelector('#analyze-form button[type="submit"]');
            if (submitButton) {
                submitButton.disabled = true;
                submitButton.innerHTML = '<i class="fas fa-circle-notch fa-spin"></i> Processing...';
            }
            console.log('Calling analyzeMetadata API with:', { image });
            const response = await stegApi.analyzeMetadata(image);
            console.log('Analyze API response:', response);
            this.handleAnalyzeSuccess(response);
        } catch (error) {
            console.error('Error analyzing image:', error);
            showToast('error', 'Error', error.message || 'Failed to analyze image. Please try again.');
        } finally {
            const submitButton = document.querySelector('#analyze-form button[type="submit"]');
            if (submitButton) {
                submitButton.disabled = !this.forms['analyze'].valid;
                submitButton.innerHTML = 'Analyze Image';
            }
        }
    }

    /**
     * Handle successful analysis
     * @param {Object} response - API response
     */
    handleAnalyzeSuccess(response) {
        console.log('Analyze success response:', response);
        if (!response.data) {
            showToast('warning', 'Warning', 'Unexpected server response format');
            console.error('Invalid response format:', response);
            return;
        }

        const data = response.data;
        console.log('Processing analyze response data:', data);

        const resultsContainer = document.getElementById('analyze-results');
        if (!resultsContainer) {
            console.error('Results container not found: analyze-results');
            return;
        }

        const imageDetailsElement = document.getElementById('analyze-image-details');
        if (imageDetailsElement) {
            imageDetailsElement.innerHTML = `
        <div class="detail-item">
          <span class="detail-label">File Name:</span>
          <span class="detail-value">${data.filename || 'Unknown'}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">File Size:</span>
          <span class="detail-value">${stegApi.formatFileSize(data.fileSize || 0)}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">File Type:</span>
          <span class="detail-value">${data.fileType || 'Unknown'}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">MIME Type:</span>
          <span class="detail-value">${data.mimeType || 'Unknown'}</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">Dimensions:</span>
          <span class="detail-value">${data.imageWidth || 0} x ${data.imageHeight || 0} pixels</span>
        </div>
        <div class="detail-item">
          <span class="detail-label">Modified:</span>
          <span class="detail-value">${data.modTime ? new Date(data.modTime).toLocaleString() : 'Unknown'}</span>
        </div>
      `;
            console.log('Set analyze-image-details');
        } else {
            console.warn('Image details element not found: analyze-image-details');
        }

        const capacityElement = document.getElementById('analyze-capacity-details');
        if (capacityElement) {
            if (data.steganoCapacity) {
                const capacity = data.steganoCapacity;
                capacityElement.innerHTML = `
          <div class="detail-item">
            <span class="detail-label">Total Capacity:</span>
            <span class="detail-value">${stegApi.formatFileSize(capacity.bytes || 0)}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Text Capacity:</span>
            <span class="detail-value">${capacity.text?.characters || 0} characters (approx. ${capacity.text?.words || 0} words)</span>
          </div>
        `;
                console.log('Set analyze-capacity-details');
            } else {
                const image = fileHandler.getFile('analyze');
                if (image) {
                    const img = new Image();
                    img.onload = () => {
                        const width = img.width;
                        const height = img.height;
                        const pixelCount = width * height;
                        const byteCapacity = Math.floor((pixelCount * 3) / 8);
                        const charCapacity = Math.floor(byteCapacity * 0.8);
                        const wordCapacity = Math.floor(charCapacity / 6);

                        capacityElement.innerHTML = `
              <div class="detail-item">
                <span class="detail-label">Total Capacity:</span>
                <span class="detail-value">${stegApi.formatFileSize(byteCapacity)}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Text Capacity:</span>
                <span class="detail-value">${charCapacity} characters (approx. ${wordCapacity} words)</span>
              </div>
            `;
                        console.log('Set fallback analyze-capacity-details');
                    };
                    img.src = URL.createObjectURL(image);
                }
            }
        } else {
            console.warn('Capacity details element not found: analyze-capacity-details');
        }

        const risksElement = document.getElementById('analyze-privacy-risks');
        if (risksElement) {
            if (data.privacyRisks && data.privacyRisks.length > 0) {
                const risksList = data.privacyRisks.map(risk => `<li>${risk}</li>`).join('');
                risksElement.innerHTML = `
          <h4>Privacy Risks:</h4>
          <ul class="risk-list">
            ${risksList}
          </ul>
        `;
                console.log('Set analyze-privacy-risks with risks');
            } else {
                risksElement.innerHTML = `
          <h4>Privacy Risks:</h4>
          <p>No significant privacy risks detected.</p>
        `;
                console.log('Set analyze-privacy-risks with no risks');
            }
        } else {
            console.warn('Privacy risks element not found: analyze-privacy-risks');
        }

        const exifElement = document.getElementById('analyze-exif-data');
        if (exifElement) {
            if (data.hasEXIF && data.properties && Object.keys(data.properties).length > 0) {
                const exifRows = Object.entries(data.properties)
                    .map(([key, value]) => `
            <tr>
              <td>${key}</td>
              <td>${value}</td>
            </tr>
          `).join('');

                exifElement.innerHTML = `
          <h4>EXIF Metadata:</h4>
          <div class="table-container">
            <table class="exif-table">
              <thead>
                <tr>
                  <th>Property</th>
                  <th>Value</th>
                </tr>
              </thead>
              <tbody>
                ${exifRows}
              </tbody>
            </table>
          </div>
        `;
                console.log('Set analyze-exif-data with EXIF data');
            } else {
                exifElement.innerHTML = `
          <h4>EXIF Metadata:</h4>
          <p>No EXIF metadata found in this image.</p>
        `;
                console.log('Set analyze-exif-data with no EXIF');
            }
        } else {
            console.warn('EXIF data element not found: analyze-exif-data');
        }

        resultsContainer.style.display = 'block';
        console.log('Show analyze-results section');

        setTimeout(() => {
            resultsContainer.scrollIntoView({ behavior: 'smooth' });
        }, 100);

        showToast('success', 'Success', 'Image analysis completed.');
    }

    /**
     * Get a stored file
     * @param {string} section - Section ID
     * @param {string} type - File type
     * @returns {File|null} The stored file
     */
    getFile(section, type = 'image') {
        const mappedSection = section === 'hide' ? 'hide-text' : section;

        return fileHandler.getFile(mappedSection, type);
    }
}

const formValidator = new FormValidator();

document.addEventListener('DOMContentLoaded', () => {
    const style = document.createElement('style');
    style.textContent = `
    .invalid {
      border-color: red !important;
    }
    .valid {
      border-color: green !important;
    }
    .enabled {
      cursor: pointer;
      opacity: 1;
    }
    .disabled {
      cursor: not-allowed;
      opacity: 0.7;
    }
    .debug-force-enabled {
      border: 2px dashed green !important;
    }
    
    /* Results styling */
    .results-section {
      margin-top: 30px;
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      overflow: hidden;
    }
    
    .card-header {
      background-color: #f8f9fa;
      border-bottom: 1px solid rgba(0, 0, 0, 0.125);
      padding: 12px 20px;
    }
    
    .card-header h3 {
      margin: 0;
      font-size: 1.5rem;
      color: #333;
    }
    
    .card-body {
      padding: 20px;
    }
    
    .result-image-container {
      text-align: center;
      margin-bottom: 20px;
    }
    
    .result-image {
      max-width: 100%;
      border-radius: 4px;
      box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
    }
    
    .result-details {
      padding: 10px;
    }
    
    .details-container {
      background-color: #f8f9fa;
      border-radius: 4px;
      padding: 15px;
      margin-top: 15px;
    }
    
    .detail-item {
      margin-bottom: 10px;
      display: flex;
    }
    
    .detail-label {
      font-weight: bold;
      margin-right: 10px;
      min-width: 100px;
    }
    
    .detail-value {
      flex-grow: 1;
    }
    
    .extracted-text {
      background-color: #f8f9fa;
      border-radius: 4px;
      padding: 15px;
      margin: 15px 0;
      max-height: 300px;
      overflow-y: auto;
    }
    
    .extracted-text pre {
      white-space: pre-wrap;
      word-wrap: break-word;
      margin: 0;
    }
    
    .analysis-section {
      margin-bottom: 30px;
    }
    
    .analysis-section h4 {
      margin-bottom: 15px;
      color: #333;
    }
    
    .exif-table {
      width: 100%;
      border-collapse: collapse;
    }
    
    .exif-table th, .exif-table td {
      border: 1px solid #ddd;
      padding: 8px;
      text-align: left;
    }
    
    .exif-table th {
      background-color: #f2f2f2;
    }
    
    .table-container {
      max-height: 300px;
      overflow-y: auto;
    }
    
    .risk-list {
      margin: 0;
      padding-left: 20px;
    }
    
    .risk-list li {
      margin-bottom: 8px;
    }
  `;
    document.head.appendChild(style);

    setTimeout(() => {
        const messageField = document.getElementById('hide-text-message');
        if (messageField && messageField.value.trim().length > 0) {
            console.log("Emergency fix: Found pre-filled message field, validating");
            formValidator.validateField('hide-text-message');
            formValidator.validateForm('hide-text');
        }

        const keyField = document.getElementById('extract-key');
        if (keyField && keyField.value.trim().length > 0) {
            formValidator.validateField('extract-key');
            formValidator.validateForm('extract');
        }
    }, 500);
});
