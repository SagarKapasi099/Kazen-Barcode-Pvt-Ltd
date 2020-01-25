$(document).ready(function () {
    $('input').keydown(function () {
        $(this).parent().removeClass("has-danger");
    });
});

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
        debugger;

        if (jsonResponse.error == null) {
            alert("Username and Password does not match records\nPlease Try Again.");
        } else if (jsonResponse.token == null || jsonResponse.token.length === 0){
            alert("Please Check Internet Connection And Try Again.");
        } else {
            localStorage.setItem("jwt", jsonResponse.token);
            window.location = "administrator";

        }
    });

    // Define what happens in case of error
    XHR.addEventListener(' error', function (event) {
        alert("Please Check Internet Connection And Try Again.");
    });

    try {
        XHR.onreadystatechange = function () {
            console.log(XHR);
            if (XHR.readyState == 4 && XHR.status == 0) {
                alert("Please Check Internet Connection and Try Again");
            }
        };

        // Set up our request
        XHR.open('POST', '/getToken', false);

        // Send our FormData object; HTTP headers are set automatically
        XHR.send(formData);
    } catch (e) {
        alert("Please Check Internet Connection And Try Again.");
    }
});