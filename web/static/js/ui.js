/**
 * StegGo - Image Steganography Tool
 * UI Controller
 * Author: pranaykumar2
 * Version: 1.0.0
 * Date: 2025-03-28 12:43:20 UTC
 */

class UiController {
    constructor() {
        // Track current tab
        this.activeTab = 'hide-text';

        // Track theme
        this.isDarkMode = localStorage.getItem('theme') === 'dark';
    }

    /**
     * Initialize UI controller
     */
    initialize() {
        // Initialize tab switching
        this.initializeTabs();

        // Initialize copy buttons
        this.initializeCopyButtons();

        // Initialize result close buttons
        this.initializeCloseButtons();

        // Initialize theme toggle
        this.initializeThemeToggle();

        // Initialize tooltips
        this.initializeTooltips();
    }

    /**
     * Initialize tab switching
     */
    initializeTabs() {
        const tabButtons = document.querySelectorAll('.tab-button');
        const panels = document.querySelectorAll('.panel');

        tabButtons.forEach(button => {
            button.addEventListener('click', () => {
                const tabId = button.getAttribute('data-tab');
                this.switchTab(tabId);
            });
        });
    }

    /**
     * Switch active tab
     * @param {string} tabId - ID of tab to activate
     */
    switchTab(tabId) {
        // Update active tab tracking
        this.activeTab = tabId;

        // Update tab buttons
        const tabButtons = document.querySelectorAll('.tab-button');
        tabButtons.forEach(button => {
            const isActive = button.getAttribute('data-tab') === tabId;
            button.classList.toggle('active', isActive);
        });

        // Update panels
        const panels = document.querySelectorAll('.panel');
        panels.forEach(panel => {
            const isPanelActive = panel.id === `${tabId}-panel`;
            panel.classList.toggle('active', isPanelActive);
        });
    }

    /**
     * Initialize copy buttons
     */
    initializeCopyButtons() {
        const copyButtons = document.querySelectorAll('.btn-copy');

        copyButtons.forEach(button => {
            button.addEventListener('click', () => {
                const targetId = button.getAttribute('data-copy');
                this.copyToClipboard(targetId);
            });
        });
    }

    /**
     * Initialize close buttons for results sections
     */
    initializeCloseButtons() {
        const closeButtons = document.querySelectorAll('.close-results');

        closeButtons.forEach(button => {
            button.addEventListener('click', () => {
                const resultsSection = button.closest('.results-section');
                if (resultsSection) {
                    resultsSection.style.display = 'none';
                }
            });
        });
    }

    /**
     * Initialize theme toggle
     */
    initializeThemeToggle() {
        const themeToggle = document.querySelector('.theme-toggle');
        if (!themeToggle) return;

        // Set initial theme
        document.documentElement.setAttribute('data-theme', this.isDarkMode ? 'dark' : 'light');

        // Toggle theme on click
        themeToggle.addEventListener('click', () => {
            this.isDarkMode = !this.isDarkMode;
            document.documentElement.setAttribute('data-theme', this.isDarkMode ? 'dark' : 'light');
            localStorage.setItem('theme', this.isDarkMode ? 'dark' : 'light');
        });
    }

    /**
     * Initialize tooltips
     */
    initializeTooltips() {
        // Elements with data-tooltip attribute already have CSS-based tooltips
        // This is just a placeholder for any JavaScript-based tooltip enhancements
    }

    /**
     * Copy text to clipboard
     * @param {string} elementId - ID of element containing text to copy
     */
    copyToClipboard(elementId) {
        const element = document.getElementById(elementId);
        if (!element) return;

        // Get text to copy
        let text;
        if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
            text = element.value;
        } else {
            text = element.textContent;
        }

        // Use navigator.clipboard API if available
        if (navigator.clipboard) {
            navigator.clipboard.writeText(text)
                .then(() => {
                    showToast('success', 'Copied!', 'Text copied to clipboard');
                })
                .catch(err => {
                    console.error('Could not copy text: ', err);
                    this.fallbackCopy(text);
                });
        } else {
            this.fallbackCopy(text);
        }
    }

    /**
     * Fallback method for copying to clipboard
     * @param {string} text - Text to copy
     */
    fallbackCopy(text) {
        // Create temporary textarea
        const textarea = document.createElement('textarea');
        textarea.value = text;

        // Make it invisible but part of the document
        textarea.style.position = 'absolute';
        textarea.style.left = '-9999px';
        document.body.appendChild(textarea);

        // Select and copy
        textarea.select();

        try {
            document.execCommand('copy');
            showToast('success', 'Copied!', 'Text copied to clipboard');
        } catch (err) {
            console.error('Fallback copy failed:', err);
            showToast('error', 'Copy Failed', 'Could not copy to clipboard');
        }

        // Clean up
        document.body.removeChild(textarea);
    }

    /**
     * Show a specific results panel
     * @param {string} panelId - ID of results panel to show
     */
    showResults(panelId) {
        const panel = document.getElementById(panelId);
        if (panel) {
            panel.style.display = 'block';
            panel.scrollIntoView({ behavior: 'smooth' });
        }
    }

    /**
     * Hide a specific results panel
     * @param {string} panelId - ID of results panel to hide
     */
    hideResults(panelId) {
        const panel = document.getElementById(panelId);
        if (panel) {
            panel.style.display = 'none';
        }
    }
}

// Create instance
const uiController = new UiController();
