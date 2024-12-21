// Ensure DOM is loaded
document.addEventListener("DOMContentLoaded", () => {
    const classNameInput = document.getElementById("class_name");
    const fetchButton = document.getElementById("fetch");
    const userListContainer = document.getElementById("user-list-container");
    const errorMessage = document.getElementById("error-message");

    // Fetch users who attended the class
    async function fetchAttendance() {
        const className = classNameInput.value.trim();

        if (!className) {
            errorMessage.textContent = "Class name is required!";
            userListContainer.innerHTML = ""; // Clear previous results
            return;
        }

        try {
            const response = await fetch(`http://localhost:8090/class/${className}`, {
                method: "GET",
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();

            if (result.users && result.users.length > 0) {
                displayUsers(result.users);
                errorMessage.textContent = ""; // Clear any error messages
            } else {
                userListContainer.innerHTML = "<p>No users attended this class.</p>";
                errorMessage.textContent = ""; // Clear any error messages
            }
        } catch (error) {
            errorMessage.textContent = "Error fetching attendance data.";
            console.error("Error fetching attendance:", error);
        }
    }

    // Display the list of users
    function displayUsers(users) {
        userListContainer.innerHTML = ""; // Clear previous results

        const table = document.createElement("table");
        table.style.borderCollapse = "collapse";
        table.style.width = "100%";
        table.style.marginTop = "10px";

        const headerRow = document.createElement("tr");
        headerRow.innerHTML = `
            <th style="border: 1px solid #ddd; padding: 8px;">User Number</th>
            <th style="border: 1px solid #ddd; padding: 8px;">Name</th>
        `;
        table.appendChild(headerRow);

        users.forEach((user) => {
            const row = document.createElement("tr");
            row.innerHTML = `
                <td style="border: 1px solid #ddd; padding: 8px;">${user.user_number}</td>
                <td style="border: 1px solid #ddd; padding: 8px;">${user.name}</td>
            `;
            table.appendChild(row);
        });

        userListContainer.appendChild(table);
    }

    // Event listener for fetch button
    fetchButton.addEventListener("click", fetchAttendance);
});
