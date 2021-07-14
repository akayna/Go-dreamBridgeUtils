// Asks the backend for one KeyID to generate one cardToken
function getCardTokenKey() {
    console.log("getCardTokenKey");

    var settings = {
        "async": true,
        "crossDomain": true,
        "url": "http://localhost:5000/getKey",
        "method": "GET",
        "headers": {
          "Accept": "*/*",
          "Cache-Control": "no-cache",
          //"Host": "localhost:5000",
          //"accept-encoding": "gzip, deflate",
          //"Connection": "keep-alive",
          "cache-control": "no-cache"
        }
      };
      
    $.ajax(settings).done(function (response, status) {

      console.log("Status: "+ status);

      console.log("keyID: "+ response);

      setKeyID(response);

      generateCardToken();

      return false;
    }).fail(failResponse);

    return false;
}

// Generates one cardtoken using the card form and the KeyID generated into the backend
function generateCardToken() {
    console.log("generateCardToken");

    var cardTokenData = {
        keyId: document.getElementById('keyID').value, 
        cardInfo: {
            cardNumber: document.getElementById("ccnum").value,
            cardExpirationMonth: document.getElementById("expmonth").value,
            cardExpirationYear: document.getElementById("expyear").value,
            cardType: document.getElementById("type").value
        }
    };

    var payload = JSON.stringify(cardTokenData);
    console.log("Payload: " + payload);

    var settings = {
        "async": true,
        "crossDomain": true,
        "url": "https://apitest.cybersource.com/flex/v1/tokens",
        "method": "POST",
        "headers": {
          "Content-Type": "application/json",
          "Accept": "*/*"
        },
        "processData": false,
        "data": payload
    };
      
      $.ajax(settings).done(function (response, status) {

        console.log("Status: "+ status);

        console.log("Response: ");
        console.log(response);

        setCardToken(response.token);

        return false;
      }).fail(failResponse);

    return false;
}

// Treats the one fail response.
function failResponse(response, status) {
  console.log("Status: "+ status);

  console.log("Response object:");

  console.log(response);

}