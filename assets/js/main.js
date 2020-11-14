document.addEventListener('DOMContentLoaded', function () {
    setTimeout(function () {
        document.getElementById("preloader").style.display = "none";
        document.getElementsByTagName("BODY")[0].style.overflow = "unset";
    }, 2000);


    // adding hamburger menu for tablets and below
    document.getElementsByClassName("hamburger-menu")[0].addEventListener("click", function() {
        document.getElementsByTagName("nav")[0].classList.toggle("nav_active");
    });

});
