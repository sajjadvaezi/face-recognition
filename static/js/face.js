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
            showResponseMessage("Error accessing the camera. Please allow access or check your device.", "red");
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

    // Display the response message with color
    function showResponseMessage(message, color) {
        responseMessage.textContent = message;
        responseMessage.style.backgroundColor = color === "red" ? "#FEE2E2" : "#D1FAE5";
        responseMessage.style.color = color === "red" ? "#B91C1C" : "#065F46";
        responseMessage.classList.remove("hidden");
    }

    // Clear the response message
    function clearResponseMessage() {
        responseMessage.classList.add("hidden");
        responseMessage.textContent = "";
    }

    // Handle the form submission
    captureButton.addEventListener("click", async () => {
        clearResponseMessage(); // Clear any previous messages
        const userNumber = userNumberInput.value.trim();
        if (!userNumber) {
            showResponseMessage("User number is required.", "red");
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
                showResponseMessage("Face data added successfully!", "green");
            } else {
                const errorMessage =
                    data.error.includes("UNIQUE constraint failed")
                        ? "Face data already exists for this user."
                        : data.error || "Failed to add face data.";
                showResponseMessage(`Error: ${errorMessage}`, "red");
            }
        } catch (error) {
            console.error("Error submitting face data:", error);
            showResponseMessage("Error connecting to the server. Please try again.", "red");
        }
    });

    // Start the camera on page load
    startCamera();
});
