window.addEventListener('load', function () { // show hide barcode onload
    console.log("start");
    document.getElementById("preloader").style.display = "none";
    document.getElementsByTagName("BODY")[0].style.overflow = "unset";
});

document.addEventListener('DOMContentLoaded', function () {

    // adding hamburger menu for tablets and below
    document.getElementsByClassName("hamburger-menu")[0].addEventListener("click", function () {
        document.getElementsByTagName("nav")[0].classList.toggle("nav_active");
        document.getElementsByClassName("hamburger-menu")[0].classList.toggle("selected");
    });

    // adding svgs
    document.querySelectorAll(".icon")
        .forEach(item => {
            try {
                let svg_icon = '';
                let icon_name = item.classList[1].substring(5)
                switch (icon_name) {
                    case "phone": {
                        svg_icon = `<svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="phone-alt" class="svg-inline--fa fa-phone-alt fa-w-16" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M497.39 361.8l-112-48a24 24 0 0 0-28 6.9l-49.6 60.6A370.66 370.66 0 0 1 130.6 204.11l60.6-49.6a23.94 23.94 0 0 0 6.9-28l-48-112A24.16 24.16 0 0 0 122.6.61l-104 24A24 24 0 0 0 0 48c0 256.5 207.9 464 464 464a24 24 0 0 0 23.4-18.6l24-104a24.29 24.29 0 0 0-14.01-27.6z"></path></svg>`;
                        break;
                    }
                    case "mail": {
                        svg_icon = `<svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="envelope" class="svg-inline--fa fa-envelope fa-w-16" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path fill="currentColor" d="M502.3 190.8c3.9-3.1 9.7-.2 9.7 4.7V400c0 26.5-21.5 48-48 48H48c-26.5 0-48-21.5-48-48V195.6c0-5 5.7-7.8 9.7-4.7 22.4 17.4 52.1 39.5 154.1 113.6 21.1 15.4 56.7 47.8 92.2 47.6 35.7.3 72-32.8 92.3-47.6 102-74.1 131.6-96.3 154-113.7zM256 320c23.2.4 56.6-29.2 73.4-41.4 132.7-96.3 142.8-104.7 173.4-128.7 5.8-4.5 9.2-11.5 9.2-18.9v-19c0-26.5-21.5-48-48-48H48C21.5 64 0 85.5 0 112v19c0 7.4 3.4 14.3 9.2 18.9 30.6 23.9 40.7 32.4 173.4 128.7 16.8 12.2 50.2 41.8 73.4 41.4z"></path></svg>`
                        break;
                    }
                    case "website": {
                        svg_icon = `<svg aria-hidden="true" focusable="false" data-prefix="fas" data-icon="globe-asia" class="svg-inline--fa fa-globe-asia fa-w-16" role="img" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 496 512"><path fill="currentColor" d="M248 8C111.03 8 0 119.03 0 256s111.03 248 248 248 248-111.03 248-248S384.97 8 248 8zm-11.34 240.23c-2.89 4.82-8.1 7.77-13.72 7.77h-.31c-4.24 0-8.31 1.69-11.31 4.69l-5.66 5.66c-3.12 3.12-3.12 8.19 0 11.31l5.66 5.66c3 3 4.69 7.07 4.69 11.31V304c0 8.84-7.16 16-16 16h-6.11c-6.06 0-11.6-3.42-14.31-8.85l-22.62-45.23c-2.44-4.88-8.95-5.94-12.81-2.08l-19.47 19.46c-3 3-7.07 4.69-11.31 4.69H50.81C49.12 277.55 48 266.92 48 256c0-110.28 89.72-200 200-200 21.51 0 42.2 3.51 61.63 9.82l-50.16 38.53c-5.11 3.41-4.63 11.06.86 13.81l10.83 5.41c5.42 2.71 8.84 8.25 8.84 14.31V216c0 4.42-3.58 8-8 8h-3.06c-3.03 0-5.8-1.71-7.15-4.42-1.56-3.12-5.96-3.29-7.76-.3l-17.37 28.95zM408 358.43c0 4.24-1.69 8.31-4.69 11.31l-9.57 9.57c-3 3-7.07 4.69-11.31 4.69h-15.16c-4.24 0-8.31-1.69-11.31-4.69l-13.01-13.01a26.767 26.767 0 0 0-25.42-7.04l-21.27 5.32c-1.27.32-2.57.48-3.88.48h-10.34c-4.24 0-8.31-1.69-11.31-4.69l-11.91-11.91a8.008 8.008 0 0 1-2.34-5.66v-10.2c0-3.27 1.99-6.21 5.03-7.43l39.34-15.74c1.98-.79 3.86-1.82 5.59-3.05l23.71-16.89a7.978 7.978 0 0 1 4.64-1.48h12.09c3.23 0 6.15 1.94 7.39 4.93l5.35 12.85a4 4 0 0 0 3.69 2.46h3.8c1.78 0 3.35-1.18 3.84-2.88l4.2-14.47c.5-1.71 2.06-2.88 3.84-2.88h6.06c2.21 0 4 1.79 4 4v12.93c0 2.12.84 4.16 2.34 5.66l11.91 11.91c3 3 4.69 7.07 4.69 11.31v24.6z"></path></svg>`;
                        break;
                    }
                }
                item.innerHTML = svg_icon;
            } catch(e) {
                console.log(e);
            }

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

    /*const content = document.getElementsByTagName("main")[0];
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
    */

});

function hide_all_dropdowns() {
    document.querySelectorAll('.dropbtn')
        .forEach(reset_item => {
            reset_item.parentElement.children[1].classList.remove("display");
        });
}


window.addEventListener('load', function () {
    let gliderElement = document.querySelector('.glider');
    if(gliderElement == null) {
        return;
    }

    new Glider(document.querySelector('.glider'), {
        // Mobile-first defaults
        slidesToShow: 1,
        slidesToScroll: 1,
        scrollLock: true,
        dots: '.dots',
        arrows: {
            prev: '.glider-prev',
            next: '.glider-next'
        },
        responsive: [
            {
                // screens greater than >= 775px
                breakpoint: 775,
                settings: {
                    // Set to `auto` and provide item width to adjust to viewport
                    slidesToShow: 2,
                    slidesToScroll: 1,
                    itemWidth: 150,
                    duration: 0.25
                }
            },{
                // screens greater than >= 1024px
                breakpoint: 1024,
                settings: {
                    slidesToShow: 3,
                    slidesToScroll: 1,
                    itemWidth: 150,
                    duration: 0.25
                }
            }
        ]
    });
});

document.addEventListener("DOMContentLoaded", function(event) {
    document.querySelectorAll('img').forEach(function(img){
        img.onerror = function(){this.style.display='none';};
    })
});
