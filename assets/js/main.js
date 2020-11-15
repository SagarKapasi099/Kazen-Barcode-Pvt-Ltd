document.addEventListener('DOMContentLoaded', function () {
    setTimeout(function () {
        document.getElementById("preloader").style.display = "none";
        document.getElementsByTagName("BODY")[0].style.overflow = "unset";
    }, 2000);


    // adding hamburger menu for tablets and below
    document.getElementsByClassName("hamburger-menu")[0].addEventListener("click", function () {
        document.getElementsByTagName("nav")[0].classList.toggle("nav_active");
    });

    // toggle dropdown hide
    document.querySelectorAll('.dropbtn')
        .forEach(item => {
            item.addEventListener('click', function () {

                let class_list = this.parentElement.children[1].classList;
                if (class_list.contains("display")) {
                    this.parentElement.children[1].classList.remove("display");
                } else {
                    hide_all_dropdowns();
                    this.parentElement.children[1].classList.add("display");
                }
            })
        });

    const content = document.getElementsByTagName("main")[0];
    document.addEventListener("scroll", (e) => {

        let scrolled = document.scrollingElement.scrollTop;
        let position = content.offsetTop;
        let navbar = document.getElementsByTagName("header")[0];
        if (scrolled > position + 100) {
            navbar.classList.add("blur");
        } else {
            navbar.classList.remove("blur");
        }
    });

});

function hide_all_dropdowns() {
    document.querySelectorAll('.dropbtn')
        .forEach(reset_item => {
            reset_item.parentElement.children[1].classList.remove("display");
        });
}
