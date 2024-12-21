document.addEventListener("DOMContentLoaded", () => {
    const addClassForm = document.getElementById("addClassForm");
    const classNameInput = document.getElementById("className");
    const userNumberInput = document.getElementById("userNumber");
    const responseMessage = document.getElementById("responseMessage");

    addClassForm.addEventListener("submit", async (event) => {
        event.preventDefault();

        const className = classNameInput.value.trim();
        const userNumber = userNumberInput.value.trim();

        if (!className || !userNumber) {
            responseMessage.textContent = "Both fields are required.";
            responseMessage.style.color = "red";
            return;
        }

        try {
            const response = await fetch("http://localhost:8090/add/class", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    class_name: className,
                    user_number: userNumber,
                }),
            });

            const data = await response.json();

            if (response.ok) {
                responseMessage.textContent = "Class added successfully!";
                responseMessage.style.color = "green";
            } else {
                // Check for specific error indicating class name already exists
                if (data.error && data.error.includes("UNIQUE constraint failed")) {
                    responseMessage.textContent = "Error: Class name already exists.";
                } else if (data.error && data.error.includes("no user found ")) {
                    responseMessage.textContent = "Error: no teacher found with given user number.";
                }
                else {
                    responseMessage.textContent = `Error: ${data.error || "Failed to add class"}`;
                }
                responseMessage.style.color = "red";
            }
        } catch (error) {
            responseMessage.textContent = `Error: ${error.message}`;
            responseMessage.style.color = "red";
        }
    });
});
