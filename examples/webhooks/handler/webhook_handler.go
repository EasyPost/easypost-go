package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

type Handler struct {
	Username string
	Password string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if h.Username != "" {
		user, pass, ok := req.BasicAuth()
		if !ok || user != h.Username || pass != h.Password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	var event easypost.Event
	if err := json.NewDecoder(req.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("failed decoding request:", err)
		return
	}

	switch event.Description {
	case "batch.created", "batch.updated":
		h.HandleBatchEvent(&event)
	case "insurance.cancelled", "insurance.purchased":
		h.HandleInsuranceEvent(&event)
	case "payment.completed", "payment.created", "payment.failed":
		h.HandlePaymentEvent(&event)
	case "refund.successful":
		h.HandleRefundEvent(&event)
	case "report.available", "report.failed", "report.new":
		h.HandleReportEvent(&event)
	case "scan_form.created", "scan_form.updated":
		h.HandleScanFormEvent(&event)
	case "tracker.created", "tracker.updated":
		h.HandleTrackerEvent(&event)
	default:
		log.Println("unrecognized event type:", event.Description)
	}
}

func (h *Handler) HandleBatchEvent(event *easypost.Event) {
	batch, ok := event.Result.(*easypost.Batch)
	if !ok {
		log.Printf("unexpected result type for batch event: %T\n", event.Result)
		return
	}
	verb := strings.TrimPrefix(event.Description, "batch.")
	log.Printf(
		"batch %s %s with %d shipments and status "+
			"(postage_purchased: %d, postage_purchase_failed: %d, "+
			"queued_for_purchase: %d, creation_failed: %d)\n",
		batch.ID, verb, batch.NumShipments, batch.Status.PostagePurchased,
		batch.Status.PostagePurchaseFailed, batch.Status.QueuedForPurchase,
		batch.Status.CreationFailed,
	)
}

func (h *Handler) HandleInsuranceEvent(event *easypost.Event) {
	insurance, ok := event.Result.(*easypost.Insurance)
	if !ok {
		log.Printf(
			"unexpected result type for insurance event: %T\n", event.Result,
		)
		return
	}
	verb := strings.TrimPrefix(event.Description, "insurance.")
	log.Printf(
		"insurance %s %s of %s for shipment %s\n",
		insurance.ID, verb, insurance.Amount, insurance.ShipmentID,
	)
}

func (h *Handler) HandlePaymentEvent(event *easypost.Event) {
	payment, ok := event.Result.(*easypost.PaymentLog)
	if !ok {
		log.Printf("unexpected result type for payment event: %T\n", event.Result)
		return
	}
	verb := strings.TrimPrefix(event.Description, "payment.")
	log.Printf(
		"payment %s %s with amount %s and status %d\n",
		payment.ID, verb, payment.Amount, payment.Status,
	)
}

func (h *Handler) HandleRefundEvent(event *easypost.Event) {
	refund, ok := event.Result.(*easypost.Refund)
	if !ok {
		log.Printf("unexpected result type for refund event: %T\n", event.Result)
		return
	}
	verb := strings.TrimPrefix(event.Description, "refund.")
	log.Printf(
		"refund %s %s with status %s for shipment %s\n",
		refund.ID, verb, refund.Status, refund.ShipmentID,
	)
}

func (h *Handler) HandleReportEvent(event *easypost.Event) {
	report, ok := event.Result.(*easypost.Report)
	if !ok {
		log.Printf("unexpected result type for report event: %T\n", event.Result)
		return
	}
	verb := strings.TrimPrefix(event.Description, "report.")
	log.Printf(
		"report %s %s with status %s and URL %s\n",
		report.ID, verb, report.Status, report.URL,
	)
}

func (h *Handler) HandleScanFormEvent(event *easypost.Event) {
	scanForm, ok := event.Result.(*easypost.ScanForm)
	if !ok {
		log.Printf("unexpected result type for batch event: %T\n", event.Result)
		return
	}
	verb := strings.TrimPrefix(event.Description, "scan_form.")
	log.Printf(
		"scan form %s %s with status %s and tracking codes %s\n",
		scanForm.ID, verb, scanForm.Status,
		strings.Join(scanForm.TrackingCodes, ", "),
	)
}

func (h *Handler) HandleTrackerEvent(event *easypost.Event) {
	tracker, ok := event.Result.(*easypost.Tracker)
	if !ok {
		log.Printf("unexpected result type for tracker event: %T\n", event.Result)
		return
	}
	verb := strings.TrimPrefix(event.Description, "tracker.")
	log.Printf(
		"tracker %s %s with status %s and tracking code %s\n",
		tracker.ID, verb, tracker.Status, tracker.TrackingCode,
	)
}

func main() {
	var addr, user, pass, path string
	flag.StringVar(&addr, "addr", ":8080", "Local HTTP listener address")
	flag.StringVar(&user, "user", "", "HTTP user name required in requests")
	flag.StringVar(&pass, "pass", "", "HTTP user password required in requests")
	flag.StringVar(&path, "path", "/easypost/events", "HTTP webhook handler URI path")
	flag.Parse()
	handler := &Handler{Username: user, Password: pass}
	mux := http.NewServeMux()
	mux.Handle(path, handler)
	if err := http.ListenAndServe(addr, mux); err != http.ErrServerClosed {
		fmt.Println("error:", err)
	}
}
