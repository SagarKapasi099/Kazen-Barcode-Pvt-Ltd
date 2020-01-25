$(document).ready(function () {
    $('input').keydown(function () {
        $(this).parent().removeClass("has-danger");
    });
});

function displayDanger() {
    let usernameInput = $('#username');
    usernameInput.parent().addClass("has-danger");
    usernameInput.parent().removeClass("has-success");

    let passwordInput = $('#password');
    passwordInput.parent().addClass("has-danger");
    passwordInput.parent().removeClass("has-success");
}

function displaySuccess() {
    let usernameInput = $('#username');
    usernameInput.parent().removeClass("has-danger");
    usernameInput.parent().addClass("has-success");

    let passwordInput = $('#password');
    passwordInput.parent().removeClass("has-danger");
    passwordInput.parent().addClass("has-success");
}

$('#signIn').click(function (e) {
    e.preventDefault();
    let usernameInput = $('#username');

    let username = usernameInput.val();
    if (username === "" || username == null || username.length === 0) {
        usernameInput.parent().addClass("has-danger");
        return;
    }

    let passwordInput = $('#password');
    let password = passwordInput.val();
    if (password === "" || password == null || password.length === 0) {
        passwordInput.parent().addClass("has-danger");
        return;
    }

    let formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);


    const XHR = new XMLHttpRequest();

    // Define what happens on successful data submission
    XHR.addEventListener('load', function (event) {
        let jsonResponse = JSON.parse(event.target.response);

        if (jsonResponse.token == null || jsonResponse.token.length === 0){
            displayDanger();
        } else {
            displaySuccess();
            window.location = "administrator";

        }
    });

    // Define what happens in case of error
    XHR.addEventListener(' error', function (event) {
        displayDanger();
    });

    try {
        XHR.onreadystatechange = function () {
            console.log(XHR);
            if (XHR.readyState == 4 && XHR.status == 0) {
                alert("Please Check Internet Connection and Try Again");
            }
        };

        // Set up our request
        XHR.open('POST', '/getToken', true);

        // Send our FormData object; HTTP headers are set automatically
        XHR.send(formData);
    } catch (e) {
        if (XHR.status === 401) {
            displayDanger();
        } else {
            alert("Please Check Internet Connection And Try Again.");
        }
    }
});