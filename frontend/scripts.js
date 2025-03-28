document.addEventListener("DOMContentLoaded", function() {
    const hideTextForm = document.getElementById("hide-text-form");
    const hideFileForm = document.getElementById("hide-file-form");
    const extractContentForm = document.getElementById("extract-content-form");
    const analyzeMetadataForm = document.getElementById("analyze-metadata-form");
    const apiResponseContent = document.getElementById("api-response-content");

    hideTextForm.addEventListener("submit", function(event) {
        event.preventDefault();
        const formData = new FormData(hideTextForm);
        fetch("/api/hide", {
            method: "POST",
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            displayApiResponse(data);
        })
        .catch(error => {
            displayApiError(error);
        });
    });

    hideFileForm.addEventListener("submit", function(event) {
        event.preventDefault();
        const formData = new FormData(hideFileForm);
        fetch("/api/hideFile", {
            method: "POST",
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            displayApiResponse(data);
        })
        .catch(error => {
            displayApiError(error);
        });
    });

    extractContentForm.addEventListener("submit", function(event) {
        event.preventDefault();
        const formData = new FormData(extractContentForm);
        fetch("/api/extract", {
            method: "POST",
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            displayApiResponse(data);
        })
        .catch(error => {
            displayApiError(error);
        });
    });

    analyzeMetadataForm.addEventListener("submit", function(event) {
        event.preventDefault();
        const formData = new FormData(analyzeMetadataForm);
        fetch("/api/metadata", {
            method: "POST",
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            displayApiResponse(data);
        })
        .catch(error => {
            displayApiError(error);
        });
    });

    function displayApiResponse(data) {
        apiResponseContent.textContent = JSON.stringify(data, null, 2);
    }

    function displayApiError(error) {
        apiResponseContent.textContent = `Error: ${error.message}`;
    }
});
