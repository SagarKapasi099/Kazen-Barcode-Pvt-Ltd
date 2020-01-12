const buttons = document.getElementsByClassName("button");

const buttonClickFunction = function () {
    let products = JSON.parse(localStorage.getItem("products"));
    if (products == null) {
        products = [];
    }
    if (this.classList.contains('success')) {
        products = products.filter(e => e !== this.getAttribute("data-id"));
        localStorage.setItem("products", JSON.stringify(products));
        this.classList.remove('success');
    } else {
        products.push(this.getAttribute("data-id"));
        localStorage.setItem("products", JSON.stringify(products));
        this.classList.add('success');
    }
};

for (let i = 0; i < buttons.length; i++) {
    buttons[i].addEventListener('click', buttonClickFunction, false);
}

// Select All The Products From The Local Storage
window.onload = function () {
    let products = JSON.parse(localStorage.getItem("products"));
    if (products != null) {
        const productButtons = document.querySelectorAll('[data-id]');
        for (let i = 0; i < productButtons.length; ++i) {
            if (products.includes(productButtons[i].getAttribute("data-id"))) {
                productButtons[i].classList.add('success');
            }
        }
    }
};
