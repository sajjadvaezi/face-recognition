// Ensure DOM is loaded
document.addEventListener("DOMContentLoaded", () => {
    const video = document.getElementById("video");
    const captureButton = document.getElementById("capture");
    const submitButton = document.getElementById("submit");
    const classNameInput = document.getElementById("class_name");
    const errorMessage = document.getElementById("error-message");
    const successMessage = document.getElementById("success-message");
    let base64Image = "";

    // Initialize camera
    async function startCamera() {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({
                video: true,
            });
            video.srcObject = stream;
        } catch (error) {
            errorMessage.textContent = "Camera access denied or unavailable.";
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
        successMessage.textContent = "Image captured successfully!";
        errorMessage.textContent = ""; // Clear error message
    }

    // Submit attendance
    async function submitAttendance() {
        const className = classNameInput.value.trim();

        if (!base64Image || !className) {
            errorMessage.textContent = "Both image and class name are required!";
            return;
        }

        try {
            const response = await fetch("http://localhost:8090/class/attend", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    image: base64Image,
                    class_name: className,
                }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();

            if (result.status === "success") {
                successMessage.textContent = "Attendance marked successfully!";
                errorMessage.textContent = ""; // Clear error message
            } else if (
                result.error &&
                result.error.includes("could not attend the class")
            ) {
                errorMessage.textContent = "User already attended the class.";
            } else {
                errorMessage.textContent = "Failed to mark attendance.";
                console.error("Error:", result.error);
            }
        } catch (error) {
            errorMessage.textContent = "Error submitting attendance.";
            console.error("Error submitting attendance:", error);
        }
    }

    // Event listeners
    captureButton.addEventListener("click", captureImage);
    submitButton.addEventListener("click", submitAttendance);

    // Start the camera on page load
    startCamera();
});
