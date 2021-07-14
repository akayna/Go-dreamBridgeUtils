package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rafaelcunha/Go-CyberSource/cybersource3ds"
	"github.com/rafaelcunha/Go-CyberSource/cybersourcecommons"
	"github.com/rafaelcunha/Go-CyberSource/cybersourceflex"
	"github.com/rafaelcunha/Go-CyberSource/cybersourcegateway"
	"github.com/rafaelcunha/Go-CyberSource/cybersourcejwt"
	"github.com/rafaelcunha/Go-CyberSource/cybersourcetms"
	"github.com/rafaelcunha/Go-dreambridgeUtils/jsonfile"
	"github.com/rafaelcunha/Go-dreambridgeUtils/timeutils"
)

type middleware func(http.Handler) http.Handler
type middlewares []middleware

func (mws middlewares) apply(hdlr http.Handler) http.Handler {
	if len(mws) == 0 {
		return hdlr
	}
	return mws[1:].apply(mws[0](hdlr))
}

func (c *controller) shutdown(ctx context.Context, server *http.Server) context.Context {
	ctx, done := context.WithCancel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer done()

		<-quit
		signal.Stop(quit)
		close(quit)

		atomic.StoreInt64(&c.healthy, 0)
		server.ErrorLog.Printf("Server is shutting down...\n")

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			server.ErrorLog.Fatalf("Could not gracefully shutdown the server: %s\n", err)
		}
	}()

	return ctx
}

type controller struct {
	logger        *log.Logger
	nextRequestID func() string
	healthy       int64
}

var credentials cybersourcecommons.Credentials

func main() {
	// Carrega credenciais
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_safra.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_rafaelcunha.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_ezpay.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_willian.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_enext.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_asaas.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_mrpagamentos.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_adiq.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_adiqbr.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_tycoon.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_giropagamentos.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_iterative.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_redeceler.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_zoop.json", &credentials)
	//err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_latam.json", &credentials)

	err := jsonfile.ReadJSONFile2("../../Testes/arquivos/", "credenciais_brazil3ds_prod.json", &credentials)

	if err != nil {
		log.Println("Erro ao ler credenciais.")
		log.Println("Erro: ", err)

		return
	}

	listenAddr := ":5000"
	if len(os.Args) == 2 {
		listenAddr = os.Args[1]
	}

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Printf("Server is starting...")

	c := &controller{logger: logger, nextRequestID: func() string { return strconv.FormatInt(time.Now().UnixNano(), 36) }}
	router := http.NewServeMux()
	//router.HandleFunc("/", c.index)
	router.HandleFunc("/healthz", c.healthz)

	router.HandleFunc("/getKey", getKey)

	router.HandleFunc("/3DSJWT", getJWT)

	router.HandleFunc("/3DSEnrollment", doEnrollment)

	router.HandleFunc("/3DSValidate", doValidate)

	router.HandleFunc("/vcoReceiveData", vcoReceiveData)

	router.HandleFunc("/authorize", doAuthorize)

	router.HandleFunc("/microformJWTKey", getMicroformJWTKey)

	directory := flag.String("d", "./", "the directory of static file to host")
	flag.Parse()

	router.Handle("/", http.StripPrefix(strings.TrimRight("/", "/"), http.FileServer(http.Dir(*directory))))

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      (middlewares{c.tracing, c.logging}).apply(router),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	ctx := c.shutdown(context.Background(), server)

	logger.Printf("Server is ready to handle requests at %q\n", listenAddr)
	atomic.StoreInt64(&c.healthy, time.Now().UnixNano())

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %q: %s\n", listenAddr, err)
	}
	<-ctx.Done()
	logger.Printf("Server stopped\n")
}

func (c *controller) index(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Hello, World!\n")
}

func (c *controller) healthz(w http.ResponseWriter, req *http.Request) {
	if h := atomic.LoadInt64(&c.healthy); h == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		fmt.Fprintf(w, "uptime: %s\n", time.Since(time.Unix(0, h)))
	}
}

func (c *controller) logging(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func(start time.Time) {
			requestID := w.Header().Get("X-Request-Id")
			if requestID == "" {
				requestID = "unknown"
			}
			c.logger.Println(requestID, req.Method, req.URL.Path, req.RemoteAddr, req.UserAgent(), time.Since(start))
		}(time.Now())
		hdlr.ServeHTTP(w, req)
	})
}

func (c *controller) tracing(hdlr http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := req.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = c.nextRequestID()
		}
		w.Header().Set("X-Request-Id", requestID)
		hdlr.ServeHTTP(w, req)
	})
}

func getMicroformJWTKey(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Generating Micrform JWT:")

	targetorigin := ""

	if req.TLS == nil {
		targetorigin = "http://"
	} else {
		targetorigin = "https://"
	}

	targetorigin = targetorigin + req.Host

	response, responseMessage, err := cybersourceflex.GenerateMicroformKey(&credentials.CyberSourceCredential, targetorigin)
	if err != nil {
		log.Println("main - getMicroformJWTKey - Erro ao gerar JWT para o microform.")
		return
	}

	fmt.Println("response:", responseMessage)

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println("main - getMicroformJWTKey - Error converting struct to json string.")
		return
	}

	log.Println("Objeto de resposta: ")
	log.Println(string(jsonData))

	w.Write([]byte(*response.KeyID))
}

func getKey(w http.ResponseWriter, req *http.Request) {
	generatedKey, msg, err := cybersourceflex.GenerateKey(&credentials.CyberSourceCredential, nil)

	if err != nil {
		log.Println("main - Error generating key.")
		log.Println(err)
		return
	}

	fmt.Println(msg)
	fmt.Printf("Key: %+v\n", generatedKey)
	fmt.Printf("KeyID: %+v\n", generatedKey.KeyID)
	fmt.Printf("Der: %+v\n", generatedKey.Der)
	fmt.Printf("JWK: %+v\n", generatedKey.Jwk)

	w.Write([]byte(*generatedKey.KeyID))
}

type teste struct {
	OrderDetails struct {
		OrderNumber string `json:"OrderNumber"`
	} `json:"OrderDetails"`
}

type getJWTResponse struct {
	JWT         string `json:"jwt"`
	OrderNumber string `json:"orderNumber"`
}

func getJWT(w http.ResponseWriter, req *http.Request) {
	orderNumber := "order_" + timeutils.GetTimeID()

	fmt.Println("OrderNumber: " + orderNumber)

	var payload teste

	payload.OrderDetails.OrderNumber = orderNumber

	jwtClaims := cybersourcejwt.Claims{
		CardinalCredentials: &credentials.CardinalCredential,
		ReferenceID:         payload.OrderDetails.OrderNumber,
		Payload:             payload,
	}

	jwtString := jwtClaims.GetJWTString()

	fmt.Print("JWT Gerado: ")
	fmt.Println(jwtString)

	response := getJWTResponse{
		JWT:         jwtString,
		OrderNumber: orderNumber,
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		errorString := "getJWT - Erro converting struct to JSON - " + err.Error()
		log.Printf("getJWT - Erro converting struct to JSON - %s", err)
		w.Write([]byte(errorString))
		return
	}

	w.Write(responseJSON)
}

func doEnrollment(w http.ResponseWriter, req *http.Request) {
	log.Println("doEnrollment")

	var enrollmentData cybersource3ds.EnrollmentRequestData

	err := json.NewDecoder(req.Body).Decode(&enrollmentData)
	if err != nil {
		log.Printf("doEnrollment - Erro converting Post Body to JSON - %s", err)
		return
	}

	//log.Printf("Header: %+v\n", req.Header)
	//log.Println("Host: " + req.Host)

	if enrollmentData.PaymentInformation.Customer != nil {
		log.Println("Teste de Token retrieval")

		PaymentInstrumentResponse, returnMsg, err := cybersourcetms.RetrievePaymentInstrument(&credentials.CyberSourceCredential, *enrollmentData.PaymentInformation.Customer.CustomerID)

		log.Println("Token retrieve teste message: " + returnMsg)

		if err != nil {
			log.Printf("doEnrollment - Erro testing token retrieve - %s\n", err)
		} else {
			log.Printf("Token retrieve response: %+v\n", PaymentInstrumentResponse)
		}
	}

	var mcc = "5399"
	var messageCategory = "01"
	var productCode = "01"
	var transactionMode = cybersource3ds.TransactionModeECOMMERCE
	var acsWindowSize = "05"
	//var requestorID = "CARDCYBS_5b16ebc085282c2b20313e7b"
	//var requestorName = "Braspag"

	enrollmentData.ConsumerAuthenticationInformation.MCC = &mcc
	//enrollmentData.ConsumerAuthenticationInformation.RequestorId = &requestorID
	//enrollmentData.ConsumerAuthenticationInformation.RequestorName = &requestorName
	enrollmentData.ConsumerAuthenticationInformation.MessageCategory = &messageCategory
	enrollmentData.ConsumerAuthenticationInformation.ProductCode = &productCode
	enrollmentData.ConsumerAuthenticationInformation.TransactionMode = &transactionMode
	enrollmentData.ConsumerAuthenticationInformation.AcsWindowSize = &acsWindowSize

	// Acrescenta dados de Merchant
	//var merchantName = "Brazil Test"
	//var merchantURL = "https://merchantrul.com"

	//MerchantInformationData := new(cybersource3ds.MerchantInformation)
	//MerchantDescriptorData := new(cybersource3ds.MerchantDescriptor)

	//MerchantInformationData.MerchantName = &merchantName
	//MerchantDescriptorData.Name = &merchantName
	//MerchantDescriptorData.URL = &merchantURL

	//MerchantInformationData.MerchantDescriptor = MerchantDescriptorData
	//enrollmentData.MerchantInformation = MerchantInformationData

	log.Println("Paymenta data:")
	log.Printf("%+v\n", enrollmentData)

	enrollmentResponse, returnString, err := cybersource3ds.EnrollmentRequest(&credentials.CyberSourceCredential, &enrollmentData)

	if err != nil || enrollmentResponse == nil {
		log.Println("doEnrollment - Error during enrollment request: " + returnString)
		http.Error(w, returnString, http.StatusBadRequest)
		return
	}

	enrollmentResponseJSON, err := json.Marshal(enrollmentResponse)

	if err != nil {
		errorString := "doEnrollment - Erro converting struct to JSON - " + err.Error()
		log.Printf("doEnrollment - Erro converting struct to JSON - %s", err)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	log.Println("Enrollment response:")
	log.Println(string(enrollmentResponseJSON))

	w.Write(enrollmentResponseJSON)
}

func doValidate(w http.ResponseWriter, req *http.Request) {
	log.Println("doValidate")

	var jwt string

	err := json.NewDecoder(req.Body).Decode(&jwt)
	if err != nil {
		log.Println("doValidate - Error converting json to struct: ", err)
		http.Error(w, "doValidate - Error converting json to struct.", http.StatusBadRequest)
		return
	}

	claims := cybersourcejwt.ValidateReadJWT(jwt, credentials.CardinalCredential.APIKeyID)

	if claims == nil {
		log.Printf("doValidate - JWT error.")
		http.Error(w, "doValidate - JWT error.", http.StatusBadRequest)
		return
	}

	payload, ok := claims["Payload"].(map[string]interface{})
	if !ok {
		log.Printf("doValidate - JWT error.")
		http.Error(w, "doValidate - JWT error.", http.StatusBadRequest)
		return
	}

	payment, ok := payload["Payment"].(map[string]interface{})

	if !ok {
		log.Printf("doValidate - JWT error.")
		http.Error(w, "doValidate - JWT error.", http.StatusBadRequest)
		return
	}

	authenticationTransactionID := payment["ProcessorTransactionId"].(string)

	validationRequestData := cybersource3ds.ValidationRequestData{
		ConsumerAuthenticationInformation: &cybersourcecommons.ConsumerAuthenticationInformation{
			AuthenticationTransactionID: &authenticationTransactionID,
		},
	}

	validationResponse, returnString, err := cybersource3ds.ValidationtRequest(&credentials.CyberSourceCredential, &validationRequestData)

	if err != nil || validationResponse == nil {
		log.Println("doValidate - Error during validation request: " + returnString)
		http.Error(w, returnString, http.StatusBadRequest)
		return
	}

	validationResponseJSON, err := json.Marshal(validationResponse)

	if err != nil {
		errorString := "doValidate - Erro converting struct to JSON - " + err.Error()
		log.Printf("doValidate - Erro converting struct to JSON - %s", err)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	w.Write(validationResponseJSON)
}

func vcoReceiveData(w http.ResponseWriter, req *http.Request) {
	log.Println("vcoReceiveData")

	log.Println("VCO Request:")
	log.Println(printHTTPRequest(req))
}

func doAuthorize(w http.ResponseWriter, req *http.Request) {
	log.Println("doAuthorize")

	var paymentData cybersourcegateway.Payment

	err := json.NewDecoder(req.Body).Decode(&paymentData)
	if err != nil {
		log.Printf("doAuthorize - Erro converting Post Body to JSON - %s", err)
		return
	}

	paymentDataJSON, err := json.Marshal(paymentData)

	if err != nil {
		errorString := "doAuthorize - Erro converting struct to JSON - " + err.Error()
		log.Printf("doAuthorize - Erro converting struct to JSON - %s", err)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	log.Println("Paymenta data:")
	log.Println(string(paymentDataJSON))

	auhorizationResponse, responseString, err := cybersourcegateway.ProcessPayment(&credentials.CyberSourceCredential, &paymentData)

	if err != nil {
		log.Println("doAuthorize - Error authorizing payment.")
		log.Println("Response String: " + responseString)
		log.Println("Error: ", err)
		return
	}

	log.Println("Authorization Response String: " + responseString)

	jsonData, err := json.Marshal(auhorizationResponse)
	if err != nil {
		log.Println("doAuthorize - Error converting struct to json string.")
		return
	}

	log.Println("Authorization response:")
	log.Println(string(jsonData))
}

func printHTTPRequest(req *http.Request) string {
	// Create return string
	var request []string

	// Add the request string
	url := fmt.Sprintf("%v %v %v", req.Method, req.URL, req.Proto)
	request = append(request, url)

	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", req.Host))

	// Loop through headers
	for name, headers := range req.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if req.Method == "POST" {
		req.ParseForm()
		//request = append(request, "\n")
		request = append(request, "Parameters:\n")
		request = append(request, req.Form.Encode())
	}

	// Return the request as a string
	return strings.Join(request, "\n")
}

// main_test.go
var (
	_ http.Handler = http.HandlerFunc((&controller{}).index)
	_ http.Handler = http.HandlerFunc((&controller{}).healthz)
	_ middleware   = (&controller{}).logging
	_ middleware   = (&controller{}).tracing
)
