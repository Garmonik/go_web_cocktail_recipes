document.addEventListener("DOMContentLoaded", async function () {
    const contentContainer = document.querySelector("main.content");

    async function fetchPosts(limit = 10, offset = 0, orderBy = "created") {
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
        contentContainer.innerHTML = "";
        posts.forEach(post => {
            const postElement = document.createElement("div");
            postElement.classList.add("post");
            postElement.innerHTML = `
                <div class="post-header">
                    <img  src="data:image/png;base64,${post.author.avatar}" alt="Avatar" class="avatar">
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
    }

    await fetchPosts();
});
