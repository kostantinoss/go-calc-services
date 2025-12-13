package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
)

type ApiGateway struct {
	server           http.Server
	forwardingClient http.Client
	router           Router
	logger           log.Logger
	config           Config
}

func (gw *ApiGateway) forwardRequest(server string, w http.ResponseWriter, req *http.Request) {
	targetUrl := server + req.RequestURI

	// Copy request body
	requestBodyCopy, err := io.ReadAll(req.Body)
	if err != nil {
		gw.logger.Printf("Error reading request body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// Create new request
	forwardRequest, err := http.NewRequest(req.Method, targetUrl, bytes.NewReader(requestBodyCopy))
	if err != nil {
		gw.logger.Printf("Error creating forwarding request: %v", err)
		http.Error(w, "Error creating request to forward", http.StatusInternalServerError)
		return
	}

	// Copy original headers
	maps.Copy(forwardRequest.Header, req.Header)

	// Forward new request to backend servers
	gw.logger.Println("Forwarding request to: ", targetUrl)
	resp, err := gw.forwardingClient.Do(forwardRequest)
	if err != nil {
		gw.logger.Printf("Error forwarding request to %s: %v", targetUrl, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy response headers, body and status code to the client
	gw.logger.Printf("Received response from %s: %s", targetUrl, resp.Status)
	maps.Copy(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		gw.logger.Printf("Error copying response body: %v", err)
	}
}

func (gw *ApiGateway) mainHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: check and sanitize url before routing
	urlPath := req.URL.Path
	handler, exists := gw.router.matchRoute(urlPath)
	if !exists {
		gw.logger.Printf("Route not found: %s", urlPath)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	handler(w, req)
}

func (gw *ApiGateway) Init() error {
	log.Println("Initializing gateway...")

	gw.logger = *log.New(
		os.Stdout,
		"[GATEWAY]",
		log.LstdFlags|log.Lshortfile,
	)

	conf_ := Config{}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		gw.logger.Println("CONFIG_PATH not set, using default: config.yaml")
		configPath = "/app/config.yaml"
	}
	if err := conf_.LoadFromFile(configPath); err != nil {
		return err
	}
	gw.config = conf_

	port := gw.config.ApiGateway.GatewayServer.Port
	if port == "" {
		port = "8080"
		gw.logger.Println("PORT not set, using default: 8080")
	}

	gw.server = http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: http.HandlerFunc(gw.mainHandler),
	}

	gw.forwardingClient = http.Client{}

	gw.router = Router{
		logger: &gw.logger,
	}

	// Initialize the router
	gw.router.Init()
	gw.logger.Println("Router initialized")

	// Add routes to the router
	for _, route := range gw.config.ApiGateway.Routes {
		serverUrl, err := gw.config.GetTargetServerUrl(route.Server)
		if err != nil {
			gw.logger.Printf("Error: %v", err)
			return err
		}

		gw.router.AddRoute(route.Path, func(w http.ResponseWriter, req *http.Request) {
			gw.forwardRequest(serverUrl, w, req)
		})
		gw.logger.Printf("Route added: %s -> %s", route.Path, serverUrl)
	}

	gw.logger.Println("Routes have been registered")

	return nil
}

func (gw *ApiGateway) Start() {
	gw.logger.Printf("Starting server")
	gw.logger.Printf("Listening on: %s ...", gw.server.Addr)

	if err := gw.server.ListenAndServe(); err != nil {
		gw.logger.Printf("Server error: %v", err)
	}

}
