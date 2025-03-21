document.getElementById("signin-form").addEventListener("submit", async function(event) {
    event.preventDefault();

    if (!this.checkValidity()) {
        event.stopPropagation();
        this.classList.add('was-validated');
        return;
    }

    const errorDiv = document.getElementById("message");
    errorDiv.textContent = "";  // Clear previous error (if any)
    errorDiv.classList.remove(["text-danger", "text-success"])

    const formData = new FormData();

    // Add a text field
    formData.append("username", document.getElementById("username").value);
    formData.append("password", document.getElementById("password").value);

    const response = await fetch("/api/users/login", {
        method: "POST",
        body: formData,
    });

    const resultMessage = await response.text();

    if (response.ok) {
        window.location.href = "/dashboard";
    } else {
        errorDiv.classList.add("text-danger")
        errorDiv.textContent = resultMessage
    }
});


document.getElementById("signup-form").addEventListener("submit", async function(event) {
    event.preventDefault();

    if (!this.checkValidity()) {
        event.stopPropagation();
        this.classList.add('was-validated');
        return;
    }

    const errorDiv = document.getElementById("message");
    errorDiv.textContent = "";  // Clear previous error (if any)
    errorDiv.classList.remove(["text-danger", "text-success"])

    const formData = new FormData();

    // Add a text field
    formData.append("username", document.getElementById("username").value);
    formData.append("password", document.getElementById("password").value);

    const response = await fetch("/api/users/register", {
        method: "POST",
        body: formData,
    });

    const resultMessage = await response.text();

    if (response.ok) {
        //window.location.href = "/dashboard";
        errorDiv.textContent = resultMessage
        errorDiv.classList.add("text-success")
        
    } else {
        errorDiv.classList.add("text-danger")
        errorDiv.textContent = resultMessage
    }
});
