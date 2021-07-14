// This function executes with a cardinal init (setup) finish
function setupCompleteCallback(setupCompleteData) {
  console.log("setupCompletCallback");

  // Enable Authenticate button
  enablePaymentInputs();

  getBIN(document.getElementById("ccnum"));
}

// Configure Cadinal´s songbird.js
function configureSongbird() {
  console.log("configureSongbird");

  Cardinal.configure({
    timeout: "6000", //The time in milliseconds to wait before a request to Centinel API is considered a timeout
    extendedtimeout: "4000", // extendedTimeout is only used in the event of the first request timing out. This configuration allows the merchant to set the timeout (in milliseconds) for subsequent retry attempts to complete a specific request to CentinelAPI.
                            // (This configuration would be useful when the merchant wants to set higher timeout values on requests).
    maxRequestRetries: "1",
    logging: {  //The level of logging to the browser console. Enable this feature to help debug and implement Songbird.
      level: "verbose"
              //Possible Values:
              //off (default) - No logging to console enabled. This is the setting to use for production systems.
              //on - Similar to info level logging, this value will provide some information about whats occurring during a transaction. This is recommended setting for merchants implementing Songbird
              //verbose - All logs are output to console. This method can be thought of as debug level logging and will be very loud when implementing Songbird, but is the level needed when getting support from the Cardinal team.
  },
    payment:{ //An object to describe how you want the user interactions to behave.
      displayLoading: true, //A flag to enable / disable a loading screen while requests are being made to Centinel API services. This can provide feedback to the end user that processing is taking place and they should not try to reload the page, or navigate away.
      //displayExitButton: true // Enables or disable the exit icon on the modal
  }});

  // Configurando funções de callback

  Cardinal.on("payments.setupComplete", function(setupCompleteData) {
    console.log("payments.setupComplete event");
    console.log(setupCompleteData);

    setupCompleteCallback(setupCompleteData);
  });

  Cardinal.on("payments.validated", function (data, jwt) {
    console.log("payments.validated event");
    console.log("Data:");
    console.log(data);

    if (jwt == null) {
      console.log("JWT is null.");
    } else {
      console.log("JWT:");
      console.log(jwt);

      validateChallenge(jwt);
    }
  });
}

// Initiate the authentication process in the backend
function initiateAuthentication() {
  console.log("Initiating Authentication.");
  disablePaymentInputs();

  var authenticationData = getAuthenticationData();

  console.log("Data sent to erollment:");
  console.log(JSON.stringify(authenticationData));
    
  var settings = {
    "async": true,
    "crossDomain": true,
    "url": "http://localhost:5000/3DSEnrollment",
    "method": "POST",
    "headers": {
      "Content-Type": "application/json",
      "Accept": "*/*",
      "Cache-Control": "no-cache",
      "cache-control": "no-cache"
    },
    "processData": false,
    "data": JSON.stringify(authenticationData),
  };
    
  $.ajax(settings).done(function (response, status) {

    console.log("Status: "+ status);

    console.log("Enrollment Response: "+ response);

    treatEnrollmentResponse(response);

  }).fail(function (failResponse) {
    console.log("Problem during enrollment request.");
    console.log(failResponse);
    window.alert("Problem during enrollment request.");   
  });
}

function initiateChallenge(acsURL, PAReq, processorTransactionId) {
  Cardinal.continue('cca', {
      "AcsUrl": acsURL,
      "Payload": PAReq,
    },
    {
      "OrderDetails":{
        "TransactionId": processorTransactionId
      }
  });
}

function treatEnrollmentResponse(enrollmentResponse) {
  var objEnrollment = JSON.parse(enrollmentResponse);

  saveAutenticationData(objEnrollment);

  console.log("3DS protocol version: "+ objEnrollment.consumerAuthenticationInformation.specificationVersion);

  switch(objEnrollment.consumerAuthenticationInformation.veresEnrolled) {
    case "Y":
        switch(objEnrollment.consumerAuthenticationInformation.paresStatus) {
          case "Y":
            console.log("Successful silent authentication.");
            window.alert("Successful silent authentication.");
            break;
          case "N":
            console.log("Payer cannot be authenticated - Unsuccessful.");
            window.alert("Payer cannot be authenticated - Unsuccessful.");
            break;
          case "A":
              console.log("Stand-in silent authentication.");
              window.alert("Stand-in silent authentication.");
            break;
          case "U":
            console.log("Payer cannot be authenticated - Unavailable.");
            window.alert("Payer cannot be authenticated - Unavailable.");
            break;
          case "R":
            console.log("Payer cannot be authenticated - Rejected.");
            window.alert("Payer cannot be authenticated - Rejected.");
            break;
          default: //case "C":
            console.log("Step-up.");

            if (confirm("Accomplish challenge?")) {
              Cardinal.continue('cca',
              {
                "AcsUrl": objEnrollment.consumerAuthenticationInformation.acsUrl,
                "Payload": objEnrollment.consumerAuthenticationInformation.pareq,
              },
              {
              "OrderDetails":{
                "TransactionId": objEnrollment.consumerAuthenticationInformation.authenticationTransactionId
              }
              });
            } else {
              console.log("Step-up canceled.");
            }
            break;
        }
      break;
    case "U":
      console.log("Payer cannot be authenticated - not Available.");
      window.alert("Payer cannot be authenticated - not Available.");
      break;
    case "B":
      console.log("Bypassed Authentication.");
      window.alert("Bypassed Authentication.");
      break;
    case "R":
      console.log("Payer cannot be authenticated - Rejected.");
      window.alert("Payer cannot be authenticated - Rejected.");
      break;
    default:
      console.log("Authentication failure. Unknown response received.");
      break;
  }
}
 
// Asks the backend for the jwt for the transaction
function getJWT() {
    console.log("getCardTokenKey");

    var settings = {
        "async": true,
        "crossDomain": true,
        "url": "http://localhost:5000/3DSJWT",
        "method": "GET",
        "headers": {
          "Accept": "*/*",
          "Cache-Control": "no-cache",
          "cache-control": "no-cache"
        }
      };
      
    $.ajax(settings).done(function (response, status) {

      console.log("Status: "+ status);

      console.log("Response: "+ response);

      var responseObj = JSON.parse(response);

      console.log("ResponseObj: ");

      console.log(responseObj);

      document.getElementById('jwt').setAttribute('value',responseObj.jwt);
      document.getElementById('orderNumber').setAttribute('value',responseObj.orderNumber);

      console.log("Initiating SongBird. Cancelled");

      Cardinal.setup("init", {
        jwt: responseObj.jwt,
      });

    }).fail(failResponse);

}

// Executes the Cardinal´s bin detection
function executeBinDetection(bin) {
  console.log("Initiating BIN detection.");

  //Cardinal.trigger("accountNumber.update", bin);
  Cardinal.trigger("bin.process", bin);
}

// Verifies the BIN input field ensuring only numbers
var previousBin = "";

function getBIN(e) {
  e.value=e.value.replace(/[^\d]/,'');

  inputValue = e.value;

  if (inputValue.length >= 6) {

    bin = inputValue.slice(0,6);

    console.log("BIN: " + bin);

    if (bin != previousBin) {
      previousBin = bin;

      executeBinDetection(bin);
    }
  }
}

function validateChallenge(jwt) {

  var settings = {
    "async": true,
    "crossDomain": true,
    "url": "http://localhost:5000/3DSValidate",
    "method": "POST",
    "headers": {
      "Content-Type": "text/html; charset=utf-8",
      "Accept": "*/*",
      "Cache-Control": "no-cache",
      "cache-control": "no-cache"
    },
    "processData": false,
    "data": JSON.stringify(jwt),
  };
    
  $.ajax(settings).done(function (response, status) {

    console.log("Status: "+ status);

    console.log("Validation Response: "+ response);

    treatValidationResponse(response);


  }).fail(function (failResponse) {
    console.log("Problem during velidation request.");
    window.alert("Problem during velidation request.");   
  });
}

function treatValidationResponse(validationResponse) {
  var objValidation = JSON.parse(validationResponse);

  saveAutenticationData(objValidation);

  console.log("3DS protocol version: "+ objValidation.consumerAuthenticationInformation.specificationVersion);

  switch(objValidation.consumerAuthenticationInformation.paresStatus) {
    case "Y":
      console.log("Successful authentication.");
      window.alert("Successful authentication.");
      break;
    case "N":
      console.log("Payer could be authenticated - Unsuccessful.");
      window.alert("Payer could be authenticated - Unsuccessful.");
      break;
    case "U":
      console.log("Payer could be authenticated - Unavailable.");
      window.alert("Payer could be authenticated - Unavailable.");
      break;
    case "B":
      console.log("Payer could be authenticated - MerchantBypass.");
      window.alert("Payer could be authenticated - MerchantBypass.");
      break;
    default:
        console.log("Unexpected response. paresStatus: " + objValidation.consumerAuthenticationInformation.paresStatus);
        window.alert("Unexpected response. paresStatus: " + objValidation.consumerAuthenticationInformation.paresStatus);
      break;
  }
}

function saveAutenticationData(objAutenticationData) {

  if (objAutenticationData == null) {
    return;
  }

  if(objAutenticationData.hasOwnProperty('consumerAuthenticationInformation'))
  {
    if(objAutenticationData.consumerAuthenticationInformation.hasOwnProperty('eci') && objAutenticationData.consumerAuthenticationInformation.eci != "")
    {
      saveAuthenticationInfo(objAutenticationData.consumerAuthenticationInformation.eci, 
        objAutenticationData.consumerAuthenticationInformation.cavv,
        objAutenticationData.consumerAuthenticationInformation.xid);
    } else if(objAutenticationData.consumerAuthenticationInformation.hasOwnProperty('ucafCollectionIndicator') && objAutenticationData.consumerAuthenticationInformation.ucafCollectionIndicator != "") {
      saveAuthenticationInfo(objAutenticationData.consumerAuthenticationInformation.ucafCollectionIndicator, 
        objAutenticationData.consumerAuthenticationInformation.ucafAuthenticationData,
        objAutenticationData.consumerAuthenticationInformation.xid);
    }
  }
}

function mountDeviceFingerprint() {
  console.log("Getting 3DS devicefingerprint information.");

  httpBrowserColorDepth = screen.colorDepth;
  document.getElementById('httpBrowserColorDepth').setAttribute('value',httpBrowserColorDepth);
  console.log("httpBrowserColorDepth: " + httpBrowserColorDepth);

  httpBrowserJavaEnabled = "";
  if (navigator.javaEnabled() == true)
  {
    httpBrowserJavaEnabled= "Y";
  }
  else
  {
    httpBrowserJavaEnabled= "N";
  }
  document.getElementById('httpBrowserJavaEnabled').setAttribute('value',httpBrowserJavaEnabled);
  console.log("httpBrowserJavaEnabled: " + httpBrowserJavaEnabled);

  httpBrowserJavaScriptEnabled = "Y";
  document.getElementById('httpBrowserJavaScriptEnabled').setAttribute('value',httpBrowserJavaScriptEnabled);
  console.log("httpBrowserJavaScriptEnabled: " + httpBrowserJavaScriptEnabled);

  httpBrowserLanguage = navigator.language || navigator.userLanguage;
  document.getElementById('httpBrowserLanguage').setAttribute('value',httpBrowserLanguage);
  console.log("httpBrowserLanguage: " + httpBrowserLanguage);

  httpBrowserScreenHeight = window.innerHeight;
  document.getElementById('httpBrowserScreenHeight').setAttribute('value',httpBrowserScreenHeight);
  console.log("httpBrowserScreenHeight: " + httpBrowserScreenHeight);

  httpBrowserScreenWidth = window.innerWidth;
  document.getElementById('httpBrowserScreenWidth').setAttribute('value',httpBrowserScreenWidth);
  console.log("httpBrowserScreenWidth: " + httpBrowserScreenWidth);

  var date = new Date();
  httpBrowserTimeDifference = date.getTimezoneOffset();
  document.getElementById('httpBrowserTimeDifference').setAttribute('value',httpBrowserTimeDifference);
  console.log("httpBrowserTimeDifference: " + httpBrowserTimeDifference);

  // Fazer no backend
  //httpAcceptBrowserValue
  //IP

}