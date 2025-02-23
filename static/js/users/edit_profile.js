document.addEventListener("DOMContentLoaded", async function () {
    console.log("edit_profile.js загружен!");

    const form = document.getElementById("edit-profile-form");
    const usernameField = document.getElementById("username");
    const emailField = document.getElementById("email");
    const bioField = document.getElementById("bio");
    const avatarImg = document.getElementById("image-preview");
    const avatarInput = document.getElementById("avatar-input");
    const uploadAvatarBtn = document.getElementById("upload-avatar");
    const removeAvatarBtn = document.getElementById("remove-image");
    const message = document.getElementById("message");

    let selectedAvatarFile = null;

    async function loadUserData() {
        try {
            const response = await fetch("/api/my_user/");
            if (!response.ok) throw new Error(`Ошибка загрузки данных: ${response.status}`);
            const data = await response.json();

            usernameField.value = data.username;
            emailField.value = data.email;
            bioField.value = data.bio || "";

            if (data.avatar) {
                avatarImg.src = `data:image/png;base64,${data.avatar}`;
            } else {
                avatarImg.src = "/static/images/base_avatar/default-avatar.png";
            }
        } catch (error) {
            console.error("Ошибка при загрузке данных:", error);
            message.textContent = "Ошибка загрузки данных!";
        }
    }

    uploadAvatarBtn.addEventListener("click", () => avatarInput.click());

    avatarInput.addEventListener("change", function () {
        if (avatarInput.files.length > 0) {
            const file = avatarInput.files[0];

            const allowedTypes = ["image/jpeg", "image/png", "image/webp"];
            if (!allowedTypes.includes(file.type)) {
                message.textContent = "Неверный формат файла! Разрешены: JPG, PNG, WEBP.";
                return;
            }

            selectedAvatarFile = file;
            const reader = new FileReader();
            reader.onload = function (e) {
                avatarImg.src = e.target.result;
            };
            reader.readAsDataURL(file);
        }
    });

    removeAvatarBtn.addEventListener("click", function () {
        avatarImg.src = "/static/images/base_avatar/default-avatar.png";
        selectedAvatarFile = null;
    });

    form.addEventListener("submit", async function (event) {
        event.preventDefault();

        const bioText = bioField.value.trim();

        const formData = new FormData();
        formData.append("bio", bioText);

        if (selectedAvatarFile) {
            formData.append("avatar", selectedAvatarFile);
        }

        try {
            const response = await fetch("/api/my_user/", {
                method: "PATCH",
                body: formData
            });

            if (!response.ok) {
                throw new Error(`Ошибка сохранения: ${response.status}`);
            }

            message.textContent = "Данные успешно обновлены!";
            message.style.color = "green";

            setTimeout(() => {
                window.location.href = "/my_user/";
            }, 1000);
        } catch (error) {
            console.error("Ошибка сохранения:", error);
            message.textContent = "Ошибка при сохранении данных!";
            message.style.color = "red";
        }
    });

    loadUserData();
});
