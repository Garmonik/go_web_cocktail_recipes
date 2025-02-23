document.addEventListener("DOMContentLoaded", function () {
    console.log("create_post.js загружен!"); // Проверяем, загружается ли файл

    const form = document.getElementById("create-post-form");
    const imageInput = document.getElementById("image");
    const message = document.getElementById("message");
    const uploadBox = document.getElementById("upload-box");
    const imagePreview = document.getElementById("image-preview");
    const removeImage = document.getElementById("remove-image");
    const loadingIndicator = document.getElementById("loading-indicator");

    const allowedExtensions = [".jpg", ".jpeg", ".png", ".gif", ".webp"];

    imageInput.addEventListener("change", function () {
        const file = imageInput.files[0];
        if (file) {
            const fileExt = file.name.substring(file.name.lastIndexOf(".")).toLowerCase();
            if (!allowedExtensions.includes(fileExt)) {
                message.textContent = "Недопустимый формат файла. Разрешены: jpg, jpeg, png, gif, webp.";
                imageInput.value = "";
                return;
            }

            message.textContent = "";
            const reader = new FileReader();
            reader.onload = function (e) {
                imagePreview.src = e.target.result;
                imagePreview.style.display = "block";
                uploadBox.style.display = "none";
                removeImage.style.display = "block";
            };
            reader.readAsDataURL(file);
        } else {
            resetImage();
        }
    });

    removeImage.addEventListener("click", function () {
        resetImage();
    });

    function resetImage() {
        imageInput.value = "";
        imagePreview.src = "#";
        imagePreview.style.display = "none";
        uploadBox.style.display = "flex";
        removeImage.style.display = "none";
    }

    form.addEventListener("submit", async function (event) {
        event.preventDefault();
        console.log("Форма отправляется!");

        const formData = new FormData(form);
        loadingIndicator.style.display = "block";
        message.textContent = "";

        try {
            const response = await fetch("/api/recipes/", {
                method: "POST",
                body: formData
            });

            const text = await response.text(); 
            console.log("Ответ сервера:", text);

            if (!response.ok) {
                throw new Error(`Ошибка: ${response.status} ${response.statusText}`);
            }

            message.textContent = "Пост успешно создан!";
            form.reset();
            resetImage();

            setTimeout(() => {
                window.location.href = "/my_user/";
            }, 1000);
        } catch (error) {
            message.textContent = error.message;
            console.error("Ошибка:", error);
        } finally {
            loadingIndicator.style.display = "none";
        }
    });
});
