const countElement = document.getElementById('count');
const buttons = document.getElementsByClassName("button");

const buttonClickFunction = function () {
    let products = JSON.parse(localStorage.getItem("products"));
    if (products == null) {
        products = [];
    }

    let productNames = JSON.parse(localStorage.getItem("productNames"));
    if (productNames == null) {
        productNames = [];
    }

    if (this.classList.contains('success')) {
        products = products.filter(e => e !== this.getAttribute("data-id"));
        localStorage.setItem("products", JSON.stringify(products));

        productNames = productNames.filter(e => e !== this.getAttribute("data-name")+"####"+this.getAttribute("data-id"));
        localStorage.setItem("productNames", JSON.stringify(productNames));

        this.classList.remove('success');
    } else {
        products.push(this.getAttribute("data-id"));
        localStorage.setItem("products", JSON.stringify(products));

        productNames.push(this.getAttribute("data-name")+"####"+this.getAttribute("data-id"));
        localStorage.setItem("productNames", JSON.stringify(productNames));
        this.classList.add('success');
    }
    updateCount(products.length);
};

for (let i = 0; i < buttons.length; i++) {
    buttons[i].addEventListener('click', buttonClickFunction, false);
}

// Select All The Products From The Local Storage
window.onload = function () {
    let products = JSON.parse(localStorage.getItem("products"));
    if (products == null) {
        products = [];
    }
    const productButtons = document.querySelectorAll('[data-id]');
    for (let i = 0; i < productButtons.length; ++i) {
        if (products.includes(productButtons[i].getAttribute("data-id"))) {
            productButtons[i].classList.add('success');
        }
    }
    updateCount(products.length);
    appendSelectedProductsList();
    
    if (products.length === 0) {
        document.getElementById("superText").innerText = "Please feel free to contact us if you have any questions or comments about our services.";
        document.getElementById("subText").innerText = "Please provide the following information so that we can route your request to the appropriate person and thus respond to you faster.";
    }
};

function updateCount(number) {
    if (number > 0 && number != '' && number != null) {
        countElement.innerHTML = number.toString();
    } else {
        countElement.innerHTML = "0";
    }
}

// Add Products Name In View
function appendSelectedProductsList() {
    let productsList = document.getElementById("event-meta");
    let productNames = JSON.parse(localStorage.getItem("productNames"));

    for (let i = 0; i < productNames.length; ++i) {
        productsList.innerHTML += "<li><strong></strong>"+productNames[i].split("####")['0']+"</li>";
    }
}

/*
 * Handle Enquiry
 */
document.getElementById("enquirySubmitBtn").addEventListener('click', function (e) {
    e.preventDefault();
    let number = document.getElementById("mobile_no");
    if (number.value == null || number.value == "" || number.value.length !== 10) {
        alert("Please Enter 10 Digit Mobile Number");
        return;
    }

    let form = document.getElementById('enquiryForm');

    // Setup our serialized data
    let formData = new FormData();
    formData.append("id", localStorage.getItem("products"));


    // Loop through each field in the form
    for (let i = 0; i < form.elements.length; i++) {

        let field = form.elements[i];

        // Don't serialize fields without a name, submits, buttons, file and reset inputs, and disabled fields
        if (!field.name || field.disabled || field.type === 'file' || field.type === 'reset' || field.type === 'submit' || field.type === 'button') continue;

        // If a multi-select, get all selections
        if (field.type === 'select-multiple') {
            for (let n = 0; n < field.options.length; n++) {
                if (!field.options[n].selected) continue;
                formData.append(field.name, field.options[n].value);
            }
        }

        // Convert field data to a query string
        else if ((field.type !== 'checkbox' && field.type !== 'radio') || field.checked) {
            formData.append(field.name, field.value);
        }
    }

    console.log(formData);

    const XHR = new XMLHttpRequest();

    // Define what happens on successful data submission
    XHR.addEventListener('load', function (event) {
        
        let jsonResponse = JSON.parse(event.target.response);
        document.getElementById("superText").innerText = jsonResponse.message.superText;
        document.getElementById("subText").innerText = jsonResponse.message.subText;

        document.getElementById("enquiryForm").style.display = "none";

        let otpForm = document.createElement("form");
        otpForm.setAttribute("name", "otpForm");
        otpForm.setAttribute("id", "otpForm");
        otpForm.setAttribute("class", "otpForm");
        otpForm.setAttribute("autocomplete", "off");

        let linebreak = document.createElement("br");

        let otpInput = document.createElement("input");
        otpInput.setAttribute("name", "otp");
        otpInput.setAttribute("name", "otp");
        otpInput.setAttribute("id", "otp");
        otpInput.setAttribute("class", "h-full-width h-remove-bottom");
        otpInput.setAttribute("placeholder", "Please Enter OTP for " + jsonResponse.number);
        otpInput.setAttribute("type", "email");

        let otpInputDiv = document.createElement("div");
        otpInputDiv.setAttribute("class", "form-field");
        otpInputDiv.appendChild(otpInput);

        let customerId = document.createElement("input");
        customerId.style.display = "none";
        customerId.setAttribute("name", "customerId");
        customerId.setAttribute("id", "customerId");
        customerId.setAttribute("value", jsonResponse.customerId);

        let otpSubmitButton = document.createElement("input");
        otpSubmitButton.setAttribute("name", "otpSubmitBtn");
        otpSubmitButton.setAttribute("id", "otpSubmitBtn");
        otpSubmitButton.setAttribute("value", "Submit OTP");
        otpSubmitButton.setAttribute("class", "btn btn--primary btn--large h-full-width");
        otpSubmitButton.setAttribute("type", "submit");

        let otpResendButton = document.createElement("input");
        otpResendButton.setAttribute("onclick", "resendOTP(" + jsonResponse.number + ")");
        otpResendButton.setAttribute("class", "btn btn--primary btn--large h-full-width");
        otpResendButton.setAttribute("type", "submit");
        otpResendButton.setAttribute("id", "otpResendBtn");
        otpResendButton.setAttribute("value", "Resend OTP");

        document.getElementById("formDiv").appendChild(otpForm);

        otpForm.append(otpInputDiv,
            document.createElement("br"),
            customerId,
            otpSubmitButton,
            document.createElement("br"),
            otpResendButton
        );

        otpForm.scrollIntoView(false);
    });

    // Define what happens in case of error
    XHR.addEventListener(' error', function (event) {
        alert("Please Check Internet Connection And Try Again.");
        // alert( 'Oops! Something went wrong.' );
    });

    try {
        XHR.onreadystatechange = function() {
            console.log(XHR);
            if (XHR.readyState == 4 && XHR.status == 0) {
                alert("Please Check Internet Connection and Try Again");
            }
        };

        // Set up our request
        XHR.open('POST', '/enquiry', false);

        // Send our FormData object; HTTP headers are set automatically
        XHR.send(formData);
    } catch (e) {
        alert("Please Check Internet Connection And Try Again.");
    }

});

document.addEventListener('click',function(e){
    if (e.target && e.target.id === "otpSubmitBtn") {
        e.preventDefault();

        let otp = document.getElementById("otp");

        let otpValue = otp.value;
        if (otpValue === "" || otpValue == null) {
            alert("Please Enter Valid OTP");
            return;
        }

        let customerId = document.getElementById("customerId");
        let customerIdValue = customerId.value;

        // Setup our serialized data
        let formData = new FormData();
        formData.append("otp", otpValue);
        formData.append("customerId", customerIdValue);

        const XHR = new XMLHttpRequest();

        // Define what happens on successful data submission
        XHR.addEventListener('load', function (event) {

            let jsonResponse = JSON.parse(event.target.response);

            if (jsonResponse.success !== true) {
                alert("OTP Did Not Match.\nPlease Try Again.");
            } else {

                document.getElementById("otpForm").style.display = "none";

                document.getElementById("superText").innerText = jsonResponse.message.superText;
                document.getElementById("subText").innerText = jsonResponse.message.subText;

                try {
                    document.getElementById("productsForm").remove();
                } catch(e) {

                }

                let productsForm = document.createElement('form');
                productsForm.setAttribute("id", "productsForm");

                formData.append("customerId", jsonResponse.customerId);

                productsForm.append(customerId);
                productsForm.method = "post";
                productsForm.action = "/products";
                document.body.appendChild(productsForm);

                if (jsonResponse.showProducts === true) {
                    let showProductsButton = document.createElement("input");
                    showProductsButton.setAttribute("name", "showProducts");
                    showProductsButton.setAttribute("id", "showProducts");
                    showProductsButton.setAttribute("value", "SHOW WITH PRICE");
                    showProductsButton.setAttribute("class", "btn btn--primary btn--large h-full-width");
                    showProductsButton.setAttribute("type", "submit");

                    document.getElementById("formDiv").appendChild(showProductsButton)
                }


                // alert( 'Yeah! Data sent and response loaded.' );
            }
        });

        // Define what happens in case of error
        XHR.addEventListener(' error', function (event) {
            alert("Please Check Internet Connection And Try Again.")
        });

        try {
            XHR.onreadystatechange = function() {
                console.log(XHR);
                if (XHR.readyState == 4 && XHR.status == 0) {
                    alert("Please Check Internet Connection and Try Again");
                }
            };

            // Set up our request
            XHR.open('POST', '/verifyOTP', false);

            // Send our FormData object; HTTP headers are set automatically
            XHR.send(formData);
        } catch (e) {
            alert("Please Check Internet Connection And Try Again.");
        }
    }
});

document.addEventListener('click',function(e) {
    if (e.target && e.target.id === "showProducts") {
        let form = document.getElementById("productsForm");
        form.submit();
    }
});