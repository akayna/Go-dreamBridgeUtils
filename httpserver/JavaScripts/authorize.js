
// Gather all the information and send the authorization request to the backend
function sendAuthorizationRequest() {
    console.log("Initiating Authorization.");
    disablePaymentInputs();

    var paymentData = getPaymentData();

    console.log("Data sent to authorization:");
    console.log(JSON.stringify(paymentData));

    var settings = {
        "async": true,
        "crossDomain": true,
        "url": "http://localhost:5000/authorize",
        "method": "POST",
        "headers": {
            "Content-Type": "application/json",
            "Accept": "*/*",
            "Cache-Control": "no-cache",
            "cache-control": "no-cache"
        },
        "processData": false,
        "data": JSON.stringify(paymentData),
    };
    $.ajax(settings).done(function (response, status) {

        console.log("Status: " + status);

        console.log("Authorization Response: " + response);

        //treatEnrollmentResponse(response);
    }).fail(function (failResponse) {
        console.log("Problem during authoriation request.");
        console.log(failResponse);
        window.alert("Problem during authoriation request.");
    });
}
