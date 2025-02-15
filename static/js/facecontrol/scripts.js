window.onload = function() {
    let images = [];
    if (window.location.pathname.includes("login")) {
        images = [
            "/static/images/facecontrol/base_login1.png",
            "/static/images/facecontrol/base_login2.png",
            "/static/images/facecontrol/base_login3.png",
            "/static/images/facecontrol/base_login4.png"
        ];
    } else if (window.location.pathname.includes("register")) {
        images = [
            "/static/images/facecontrol/register1.png",
            "/static/images/facecontrol/register2.png",
            "/static/images/facecontrol/register3.png",
            "/static/images/facecontrol/register4.png"
        ];
    }

    let currentIndex = 0;
    const bgContainer = document.getElementById("bg-container");

    let bg1 = document.createElement("div");
    let bg2 = document.createElement("div");

    bg1.className = "bg-image active";
    bg2.className = "bg-image";

    bg1.style.backgroundImage = `url(${images[currentIndex]})`;
    bgContainer.appendChild(bg1);
    bgContainer.appendChild(bg2);

    function changeBackground() {
        currentIndex = (currentIndex + 1) % images.length;

        let activeBg = document.querySelector(".bg-image.active");
        let inactiveBg = document.querySelector(".bg-image:not(.active)");

        inactiveBg.style.backgroundImage = `url(${images[currentIndex]}?v=${Date.now()})`;
        inactiveBg.classList.add("active");

        setTimeout(() => {
            activeBg.classList.remove("active");
        }, 1500);
    }

    setInterval(changeBackground, 5000);
};
