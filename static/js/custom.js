const buttons = document.getElementsByClassName("button");

const buttonClickFunction = function () {
    if (this.classList.contains('success')) {
        this.classList.remove('success');
    } else {
        this.classList.add('success');
    }
};

for (let i = 0; i < buttons.length; i++) {
    buttons[i].addEventListener('click', buttonClickFunction, false);
}