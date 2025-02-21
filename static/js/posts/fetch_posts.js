document.addEventListener("DOMContentLoaded", async function () {
    const contentContainer = document.querySelector("#posts-container");
    const newPostsBtn = document.getElementById("new-posts-btn");
    const popularPostsBtn = document.getElementById("popular-posts-btn");
    const postsHeader = document.querySelector(".posts-header");

    let limit = 10;
    let offset = 0;
    let orderBy = "created"; // По умолчанию "Новые"

    async function fetchPosts() {
        try {
            const response = await fetch(`/api/recipes/?limit=${limit}&offset=${offset}&order_by=${orderBy}`);
            if (!response.ok) {
                throw new Error("Ошибка при загрузке постов");
            }
            const posts = await response.json();
            renderPosts(posts);
        } catch (error) {
            console.error("Ошибка:", error);
        }
    }

    function renderPosts(posts) {
        if (offset === 0) {
            contentContainer.innerHTML = ""; // Очищаем контейнер при новой загрузке
        }

        posts.forEach(post => {
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            postElement.innerHTML = `
                <div class="post-header">
                    <img src="data:image/png;base64,${post.author.avatar}" alt="Avatar" class="avatar">
                    <span class="username">${post.author.username}</span>
                </div>
                <div class="post-content">
                    <h3>${post.name}</h3>
                    <p>${post.description}</p>
                    <img src="data:image/png;base64,${post.image}" alt="Post Image" class="post-image">
                </div>
                <div class="post-footer">
                    <button class="like-button" data-post-id="${post.id}">
                        ${post.like ? "❤️" : "🤍"}
                    </button>
                </div>
            `;
            contentContainer.appendChild(postElement);
        });

        // Показываем заголовок с кнопками, если он вдруг исчезал
        postsHeader.style.opacity = "1";
    }

    // Обработчики кнопок
    newPostsBtn.addEventListener("click", () => {
        orderBy = "created";
        offset = 0;
        newPostsBtn.classList.add("active");
        popularPostsBtn.classList.remove("active");
        fetchPosts();
    });

    popularPostsBtn.addEventListener("click", () => {
        orderBy = "popular";
        offset = 0;
        popularPostsBtn.classList.add("active");
        newPostsBtn.classList.remove("active");
        fetchPosts();
    });

    // Функция для подгрузки постов при прокрутке
    window.addEventListener("scroll", () => {
        if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 100) {
            offset += limit;
            fetchPosts();
        }
    });

    await fetchPosts();
});
