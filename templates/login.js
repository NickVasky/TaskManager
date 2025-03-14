document.getElementById("signin-form").addEventListener("submit", async function(event) {
    event.preventDefault();

    if (!this.checkValidity()) {
        event.stopPropagation();
        this.classList.add('was-validated');
        return;
    }

    const errorDiv = document.getElementById("error-message");
    errorDiv.textContent = "";  // Clear previous error (if any)

    const formData = new FormData();

    // Add a text field
    formData.append("username", document.getElementById("username").value);
    formData.append("password", document.getElementById("password").value);

    const response = await fetch("/api/auth/login", {
        method: "POST",
        body: formData,
    });

    const resultMessage = await response.text();

    if (response.ok) {
        window.location.href = "/dashboard";
    } else {
        errorDiv.textContent = resultMessage
    }
});
