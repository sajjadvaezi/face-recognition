document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("registerForm");
    const responseMessage = document.getElementById("responseMessage");

    form.addEventListener("submit", async (event) => {
        event.preventDefault(); // Prevent form submission

        const formData = {
            name: document.getElementById("name").value,
            user_number: document.getElementById("user_number").value,
            role: document.getElementById("role").value,
        };

        try {
            const response = await fetch("https://192.168.164.125:8090/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(formData),
            });


            const data = await response.json();

            if (response.ok) {
                responseMessage.textContent = `Success! User ID: ${data.id}`;
                responseMessage.style.color = "green";
            } else {
                responseMessage.textContent = `Error: ${data.error || "Failed to register"}`;
                responseMessage.style.color = "red";
            }
        } catch (error) {
            responseMessage.textContent = `Error: ${error.message}`;
            responseMessage.style.color = "red";
        }
    });
});
