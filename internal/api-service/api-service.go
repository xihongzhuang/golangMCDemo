package api_service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type APIServiceInstance struct {
	port     int
	router   *mux.Router
	server   *http.Server
	appStore *InMemoryAppStore
}

func NewAPIServiceInstance(port int) *APIServiceInstance {
	return &APIServiceInstance{
		port:     port,
		appStore: NewInMemoryAppStore(),
	}
}

//Shutdown gracefully shuts down the server without interrupting any active connections. Shutdown works by
//first closing all open listeners, then closing all idle connections, and then waiting indefinitely for
//connections to return to idle and then shut down. If the provided context expires before the shutdown is
//complete, Shutdown returns the context’s error, otherwise it returns any error returned from closing the
//Server’s underlying Listener(s)

func (ap *APIServiceInstance) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := ap.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

func ValidateRequest(req *http.Request) bool {
	isYaml := false
	if v, ok := req.Header["content-type"]; ok {
		for _, s := range v {
			if strings.Contains(s, "application/x-yml") || strings.Contains(s, "application/yaml") {
				isYaml = true
			}
		}
	} else {
		//default to yaml payload
		isYaml = true
	}
	return isYaml
}

func (ap *APIServiceInstance) StartService() error {
	if ap.router != nil {
		//started
		return nil
	}
	ap.router = mux.NewRouter()
	ap.router.Use(func(next http.Handler) http.Handler {
		timeoutHandler := func(w http.ResponseWriter, req *http.Request) {
			ctx, cancel := context.WithTimeout(req.Context(), time.Duration(120)*time.Second)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()
			log.Println(" request Method", req.Method)
			req = req.WithContext(ctx)
			next.ServeHTTP(w, req)
		}
		return http.HandlerFunc(timeoutHandler)
	})
	apiRoute := ap.router.PathPrefix("/api").Subrouter()

	apiRoute.HandleFunc("/appmetadata", func(w http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		log.Println("request all", req.Method, ",queryStr:", q)
		r, err := ap.appStore.GetAll2Yaml(q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(r)
	}).Methods(http.MethodGet)

	apiRoute.HandleFunc("/appmetadata/{metadataId}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, ok := vars["metadataId"]
		if !ok {
			http.Error(w, "metadataId not provided", http.StatusBadRequest)
			return
		}
		log.Println(" request ", id)
		if req.Method == http.MethodGet {
			r, err := ap.appStore.Get(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(r)
		} else {
			err := ap.appStore.Delete(id)
			if err != nil && errors.Is(err, ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Write([]byte(fmt.Sprintf("entry %s was deleted", id)))
		}
	}).Methods(http.MethodGet, http.MethodDelete)

	//Update

	apiRoute.HandleFunc("/appmetadata/{metadataId}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, ok := vars["metadataId"]
		if !ok {
			http.Error(w, "metadataId not provided", http.StatusBadRequest)
			return
		}
		log.Println(req.Method, " request ", id)
		isYaml := ValidateRequest(req)
		//	handling of json is out of scope in this demo
		if !isYaml {
			http.Error(w, "only yaml content is supported", http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var r []byte
		if req.Method == http.MethodPut {
			r, err = ap.appStore.Update(id, body)
		} else {
			r, err = ap.appStore.Patch(id, body)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(r)
	}).Methods(http.MethodPut, http.MethodPatch)

	apiRoute.HandleFunc("/appmetadata", func(w http.ResponseWriter, req *http.Request) {
		isYaml := ValidateRequest(req)
		//	handling of json is out of scope in this demo
		if !isYaml {
			http.Error(w, "only yaml content is supported", http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		r, err := ap.appStore.Create(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(r)
	}).Methods(http.MethodPost)
	//Search By Company Name

	ap.server = &http.Server{
		//Addr:    fmt.Sprintf("127.0.0.1:%d", ap.port),
		//to allow this service to be accessible from external of docker container
		Addr:    fmt.Sprintf("0.0.0.0:%d", ap.port),
		Handler: ap.router,
	}
	go func(ap *APIServiceInstance) {
		fmt.Println("Server Started, listening on", ap.server.Addr)
		if err := ap.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}(ap)
	return nil
}
