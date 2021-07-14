window.onload = function(){

  initMicroform();
};

// Asks the backend for the microform JWT key
function initMicroform() {
    console.log("getMicroformJWTKey");

    var settings = {
        "async": true,
        "crossDomain": true,
        "url": "http://localhost:5000/microformJWTKey",
        "method": "GET",
        "headers": {
          "Accept": "*/*",
          "Cache-Control": "no-cache",
          "cache-control": "no-cache"
        }
      };
      
    $.ajax(settings).done(function (response, status) {

        console.log("Status: "+ status);

        console.log("response: "+ response);

        document.getElementById('microformContext').setAttribute('value',response);

        // custom styles that will be applied to each field we create using Microform
        var myStyles = {  
            'input': {    
              'font-size': '14px',    
              'font-family': 'helvetica, tahoma, calibri, sans-serif',    
              'color': '#555'  
            },  
            ':focus': { 'color': 'blue' },  
            ':disabled': { 'cursor': 'not-allowed' },  
            'valid': { 'color': '#3c763d' },  
            'invalid': { 'color': '#a94442' }
          };

        console.log("Loading Microform...");
        loadMicroform(response);
        console.log("Microform loaded");


    }).fail(function (failResponse) {
        console.log("Fail to get microform context JWT.");
        console.log(failResponse);
    });
}

function loadMicroform(contextJWT) {
  // JWK is set up on the server side route for /
  var form = document.querySelector('#my-sample-form');
  var tokenizeButton = document.querySelector('#tokenize-button');
  var flexResponse = document.querySelector('#flexresponse');
  var expMonth = document.querySelector('#expMonth');
  var expYear = document.querySelector('#expYear');
  var errorsOutput = document.querySelector('#errors-output');

  console.log("Invoking the Flex SDK.");
  var flex = new Flex(contextJWT);

  console.log("Initiating the microform object.");
   // custom styles that will be applied to each field we create using Microform
  var myStyles = {  
    'input': {    
      'font-size': '14px',    
      'font-family': 'helvetica, tahoma, calibri, sans-serif',    
      'color': '#555'  
    },  
    ':focus': { 'color': 'blue' },  
    ':disabled': { 'cursor': 'not-allowed' },  
    'valid': { 'color': '#3c763d' },  
    'invalid': { 'color': '#a94442' }
  };
  var microform = flex.microform({ styles: myStyles });

  console.log("Creating and attaching the microform fields to the HTML.");
  var number = microform.createField('number', { placeholder: 'Número do cartão' });
  var securityCode = microform.createField('securityCode', { placeholder: '•••' });
  number.load('#number-container');
  securityCode.load('#securityCode-container');

  console.log("Configuring tokenize button.");
  tokenizeButton.addEventListener('click', function() {  
    var options = {    
      expirationMonth: expMonth.value,  
      expirationYear: expYear.value 
    };
    microform.createToken(options, function (err, token) {
      if (err) {
        // handle error
        console.error(err);
        errorsOutput.textContent = err.message;
      } else {
        // At this point you may pass the token back to your server as you wish.
        // In this example we append a hidden input to the form and submit it.      
        console.log(JSON.stringify(token));
        flexResponse.value = JSON.stringify(token);
        form.submit();
      }
    });
  });
}