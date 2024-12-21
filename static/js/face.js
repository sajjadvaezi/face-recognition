document.addEventListener("DOMContentLoaded", () => {
    const videoElement = document.getElementById("video");
    const captureButton = document.getElementById("captureButton");
    const userNumberInput = document.getElementById("userNumber");
    const responseMessage = document.getElementById("responseMessage");
    const canvas = document.createElement("canvas");

    // Access the user's camera
    async function startCamera() {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ video: true });
            videoElement.srcObject = stream;
            videoElement.play();
        } catch (error) {
            console.error("Error accessing the camera:", error);
            responseMessage.textContent = "Error accessing the camera.";
            responseMessage.style.color = "red";
        }
    }

    // Capture the image and convert it to base64
    function captureImage() {
        canvas.width = videoElement.videoWidth;
        canvas.height = videoElement.videoHeight;
        const context = canvas.getContext("2d");
        context.drawImage(videoElement, 0, 0, canvas.width, canvas.height);
        return canvas.toDataURL("image/jpeg").split(",")[1]; // Extract base64 data
    }

    // Handle the form submission
    captureButton.addEventListener("click", async () => {
        const userNumber = userNumberInput.value.trim();
        if (!userNumber) {
            responseMessage.textContent = "User number is required.";
            responseMessage.style.color = "red";
            return;
        }

        const base64Image = captureImage();

        try {
            const response = await fetch("http://localhost:8090/face", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    image: base64Image,
                    user_number: userNumber,
                }),
            });

            const data = await response.json();

            if (response.ok) {
                responseMessage.textContent = "Face data added successfully!";
                responseMessage.style.color = "green";
            } else {
                responseMessage.textContent = `Error: ${data.error || "Failed to add face data"}`;
                responseMessage.style.color = "red";
            }
        } catch (error) {
            responseMessage.textContent = `Error: ${error.message}`;
            responseMessage.style.color = "red";
        }
    });

    // Start the camera on page load
    startCamera();
});
