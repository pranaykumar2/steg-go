/**
 * StegGo - Image Steganography Tool
 * Main Application Script
 * Author: pranaykumar2
 * Version: 1.0.0
 * Date: 2025-03-28 12:43:20 UTC
 */

// Wait for DOM to be fully loaded
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM fully loaded, initializing application...');
    initializeApp();
});

/**
 * Initialize the application
 */
function initializeApp() {
    console.log('Initializing StegGo application...');

    // Initialize UI components
    uiController.initialize();

    // Initialize file handling
    fileHandler.initialize();

    // Initialize form validation
    formValidator.initialize();

    // Check API health
    checkApiHealth();

    console.log('Application initialization complete.');
}

/**
 * Check API health
 */
async function checkApiHealth() {
    try {
        console.log('Checking API health...');
        const response = await stegApi.getHealth();
        console.log('API health check successful:', response);
    } catch (error) {
        console.error('API health check failed:', error);
        showToast('error', 'API Error', 'The StegGo API is currently unavailable. Please try again later.');
    }
}

/**
 * Shows a toast notification
 * @param {string} type - Type of toast ('success', 'error', 'warning', 'info')
 * @param {string} title - Toast title
 * @param {string} message - Toast message
 * @param {number} duration - Duration in milliseconds
 */
function showToast(type, title, message, duration = 5000) {
    const container = document.getElementById('toast-container');

    // Create toast element
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;

    // Icon based on type
    let icon;
    switch (type) {
        case 'success': icon = '<i class="fas fa-check-circle"></i>'; break;
        case 'error': icon = '<i class="fas fa-exclamation-circle"></i>'; break;
        case 'warning': icon = '<i class="fas fa-exclamation-triangle"></i>'; break;
        case 'info': icon = '<i class="fas fa-info-circle"></i>'; break;
        default: icon = '<i class="fas fa-bell"></i>';
    }

    // Set toast content
    toast.innerHTML = `
    ${icon}
    <div class="toast-content">
      <h4>${title}</h4>
      <p>${message}</p>
    </div>
    <button class="toast-close"><i class="fas fa-times"></i></button>
  `;

    // Add to container
    container.appendChild(toast);

    // Add close functionality
    const closeButton = toast.querySelector('.toast-close');
    closeButton.addEventListener('click', () => {
        removeToast(toast);
    });

    // Auto remove after duration
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
