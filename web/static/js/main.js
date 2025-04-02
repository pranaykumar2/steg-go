window.stegGoInitialized = false;

document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM fully loaded, initializing application...');
    if (!window.stegGoInitialized) {
        initializeApp();
    }
});

function initializeApp() {
    console.log('Initializing StegGo application...');

    window.stegGoInitialized = true;

    if (typeof initThemeToggle === 'function') {
        initThemeToggle();
    }

    if (typeof initTabs === 'function') {
        initTabs();
    }

    if (typeof initFormHandlers === 'function') {
        initFormHandlers();
    }

    if (typeof initFileUploads === 'function') {
        console.log('Initializing file uploads from main.js');
        initFileUploads();
    }

    if (typeof initCloseButtons === 'function') {
        initCloseButtons();
    }

    // Initialize modals
    initModals();

    document.querySelectorAll('.btn-copy').forEach(button => {
        button.addEventListener('click', function() {
            const targetId = this.getAttribute('data-copy');
            if (targetId) {
                copyToClipboard(targetId);
            }
        });
    });

    if (typeof FileHandlers !== 'undefined' && FileHandlers.init && !window.fileHandlersInitialized) {
        window.fileHandlersInitialized = true;
        console.log('Initializing FileHandlers from main.js');
        try {
            FileHandlers.init();
        } catch (error) {
            console.error('Error initializing file handlers:', error);
        }
    }

    checkApiHealth();

    console.log('Application initialization complete.');
}

async function checkApiHealth() {
    try {
        console.log('Checking API health...');
        const response = await fetch('/api/health');
        if (!response.ok) {
            throw new Error('API health check failed');
        }
        const data = await response.json();
        console.log('API health check successful:', data);

        // Show welcome message if API is healthy
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(
                `StegGo Connected`
            );
        }
    } catch (error) {
        console.error('API health check failed:', error);
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast(
                'The StegGo API is currently unavailable. Please try again later.',
                'error',
                'API Error'
            );
        }
    }
}

/**
 * Show a toast notification (fallback if UIManager is not available)
 * @param {string} type - Type of toast ('success', 'error', 'warning', 'info')
 * @param {string} title - Toast title
 * @param {string} message - Toast message
 * @param {number} duration - Duration in milliseconds
 */
function showToast(type, title, message, duration = 5000) {
    // If UIManager exists, use that instead
    if (typeof UIManager !== 'undefined' && UIManager.showToast) {
        UIManager.showToast(message, type, title, duration);
        return;
    }

    const container = document.getElementById('toast-container');
    if (!container) {
        console.warn('Toast container not found:', message);
        return;
    }

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;

    let icon;
    switch (type) {
        case 'success': icon = '<i class="bx bx-check-circle"></i>'; break;
        case 'error': icon = '<i class="bx bx-error-circle"></i>'; break;
        case 'warning': icon = '<i class="bx bx-error"></i>'; break;
        case 'info': icon = '<i class="bx bx-info-circle"></i>'; break;
        default: icon = '<i class="bx bx-bell"></i>';
    }

    toast.innerHTML = `
    ${icon}
    <div class="toast-content">
      <h4>${title}</h4>
      <p>${message}</p>
    </div>
    <button class="toast-close"><i class="bx bx-x"></i></button>
  `;

    container.appendChild(toast);

    const closeButton = toast.querySelector('.toast-close');
    closeButton.addEventListener('click', () => {
        removeToast(toast);
    });

    setTimeout(() => {
        removeToast(toast);
    }, duration);
}

/**
 * Removes a toast notification with animation
 * @param {HTMLElement} toast - The toast element to remove
 */
function removeToast(toast) {
    toast.classList.add('removing');

    // Remove after animation completes
    setTimeout(() => {
        if (toast.parentNode) {
            toast.parentNode.removeChild(toast);
        }
    }, 300);
}


function initModals() {
    console.log('Initializing modals...');

    const privacyLinks = document.querySelectorAll('.privacy-link');
    const privacyModal = document.getElementById('privacy-modal');

    const termsLinks = document.querySelectorAll('.terms-link');
    const termsModal = document.getElementById('terms-modal');

    if ((!privacyModal && privacyLinks.length > 0) || (!termsModal && termsLinks.length > 0)) {
        console.warn('Some modal elements not found');
    }

    function openModal(modal) {
        if (!modal) return;
        modal.classList.add('active');
        document.body.classList.add('modal-open');
    }

    function closeModal(modal) {
        if (!modal) return;
        modal.classList.remove('active');
        document.body.classList.remove('modal-open');
    }

    function closeAllModals() {
        if (privacyModal) privacyModal.classList.remove('active');
        if (termsModal) termsModal.classList.remove('active');
        document.body.classList.remove('modal-open');
    }

    privacyLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            openModal(privacyModal);
        });
    });

    termsLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            openModal(termsModal);
        });
    });
    
    document.querySelectorAll('.modal .close-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            closeAllModals();
        });
    });

    document.querySelectorAll('.modal .modal-backdrop').forEach(backdrop => {
        backdrop.addEventListener('click', function() {
            closeAllModals();
        });
    });

    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape' && document.body.classList.contains('modal-open')) {
            closeAllModals();
        }
    });

    console.log('Modal initialization complete');
}

/**
 * Copy text to clipboard
 * @param {string} elementId - ID of element containing text to copy
 */
function copyToClipboard(elementId) {
    const element = document.getElementById(elementId);
    if (!element) return;

    let text;
    if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
        text = element.value;
    } else {
        text = element.textContent;
    }

    if (!text) return;

    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text)
            .then(() => {
                if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                    UIManager.showToast('Text copied to clipboard.', 'success', 'Copied!', 2000);
                }
            })
            .catch(err => {
                console.error('Error copying text: ', err);
                fallbackCopyToClipboard(text);
            });
    } else {
        fallbackCopyToClipboard(text);
    }
}

/**
 * Fallback method for copying text to clipboard
 * @param {string} text - Text to copy
 */
function fallbackCopyToClipboard(text) {
    const textArea = document.createElement('textarea');
    textArea.value = text;

    textArea.style.position = 'fixed';
    textArea.style.left = '-999999px';
    textArea.style.top = '-999999px';
    document.body.appendChild(textArea);

    textArea.focus();
    textArea.select();

    try {
        const successful = document.execCommand('copy');
        if (successful) {
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('Text copied to clipboard.', 'success', 'Copied!', 2000);
            }
        } else {
            if (typeof UIManager !== 'undefined' && UIManager.showToast) {
                UIManager.showToast('Failed to copy text.', 'error', 'Error', 3000);
            }
        }
    } catch (err) {
        console.error('Fallback: Oops, unable to copy', err);
        if (typeof UIManager !== 'undefined' && UIManager.showToast) {
            UIManager.showToast('Failed to copy text.', 'error', 'Error', 3000);
        }
    }

    document.body.removeChild(textArea);
}
