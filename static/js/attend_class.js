document.addEventListener("DOMContentLoaded", () => {
    const video = document.getElementById("video");
    const captureButton = document.getElementById("capture");
    const submitButton = document.getElementById("submit");
    const classNameInput = document.getElementById("class_name");
    const errorMessage = document.getElementById("error-message");
    const successMessage = document.getElementById("success-message");
    let base64Image = "";

    // Helper function to show messages
    function showMessage(message, type) {
        if (type === "error") {
            errorMessage.textContent = message;
            errorMessage.classList.remove("hidden");
            successMessage.classList.add("hidden");
        } else {
            successMessage.textContent = message;
            successMessage.classList.remove("hidden");
            errorMessage.classList.add("hidden");
        }
    }

    // Initialize camera
    async function startCamera() {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ video: true });
            video.srcObject = stream;
        } catch (error) {
            showMessage("Camera access denied or unavailable.", "error");
            console.error("Error starting camera:", error);
        }
    }

    // Capture image from video
    function captureImage() {
        const canvas = document.createElement("canvas");
        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;
        const context = canvas.getContext("2d");
        context.drawImage(video, 0, 0, canvas.width, canvas.height);
        base64Image = canvas.toDataURL("image/jpeg").split(",")[1]; // Convert to Base64
        showMessage("Image captured successfully!", "success");
    }

    // Submit attendance
    async function submitAttendance() {
        const className = classNameInput.value.trim();

        if (!base64Image || !className) {
            showMessage("Both image and class name are required!", "error");
            return;
        }

        try {
            const response = await fetch("https://192.168.164.125:8090/class/attend", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ image: base64Image, class_name: className }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();

            if (result.status === "success") {
                showMessage("Attendance marked successfully!", "success");
                classNameInput.value = ""; // Clear input
                base64Image = ""; // Reset image
            } else if (result.error && result.error.includes("could not attend the class")) {
                showMessage("User already attended the class.", "error");
            } else {
                showMessage("Failed to mark attendance.", "error");
                console.error("Error:", result.error);
            }
        } catch (error) {
            showMessage("Error submitting attendance.", "error");
            console.error("Error submitting attendance:", error);
        }
    }

    // Event listeners
    captureButton.addEventListener("click", captureImage);
    submitButton.addEventListener("click", submitAttendance);

    // Start the camera on page load
    startCamera();
});
