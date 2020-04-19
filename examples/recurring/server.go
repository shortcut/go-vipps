package main

import (
	"github.com/gorilla/mux"
	"github.com/torfjor/go-vipps"
	"github.com/torfjor/go-vipps/auth"
	"github.com/torfjor/go-vipps/recurring"
	"github.com/unrolled/render"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	rnd := render.New(render.Options{
		Layout:        "base",
		IsDevelopment: true,
		Funcs: []template.FuncMap{
			{
				"divide": func(a, b int) int {
					return a / b
				},
				"formatDate": func(t *time.Time) string {
					return t.Local().String()
				},
				"statusClass": func(status recurring.ChargeStatus) string {
					switch status {
					case recurring.ChargeStatusPending,
						recurring.ChargeStatusDue,
						recurring.ChargeStatusProcessing:
						return "warning"
					case recurring.ChargeStatusFailed,
						recurring.ChargeStatusCancelled:
						return "danger"
					default:
						return "success"
					}
				},
			},
		},
	})

	env := vipps.EnvironmentTesting
	authClient := auth.NewClient(env, vipps.Credentials{
		APISubscriptionKey: os.Getenv("API_KEY"),
		ClientID:           os.Getenv("CLIENT_ID"),
		ClientSecret:       os.Getenv("CLIENT_SECRET"),
	})

	recurringClient := recurring.NewClient(vipps.ClientConfig{
		Logger:      log.New(os.Stdout, "", log.LstdFlags),
		HTTPClient:  authClient,
		Environment: env,
	})

	m := mux.NewRouter()

	srv := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: m,
	}

	h := handler{
		r:     rnd,
		vipps: recurringClient,
	}

	m.Handle("/", http.RedirectHandler("/agreements", http.StatusFound))
	m.HandleFunc("/agreements", h.listAgreements).Methods("GET")
	m.HandleFunc("/agreements/{id}", h.getAgreement).Methods("GET")

	log.Fatal(srv.ListenAndServe())
}
