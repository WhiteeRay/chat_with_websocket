document.addEventListener("DOMContentLoaded", () => {
    const storedUsername = localStorage.getItem("username");

    if (storedUsername) {
        const nameInput = document.querySelector("#name");
        nameInput.value = storedUsername;


        const event = new Event("submit", { bubbles: true });
        document.querySelector("#usernameForm")?.dispatchEvent(event);
    }
});
