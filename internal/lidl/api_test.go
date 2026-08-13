package lidl

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mateuszkozakiewicz/lidl-coupons/internal/config"
)

func newTestClient(t *testing.T, serverURL string) *Client {
	t.Helper()
	return New(&config.LidlConfig{
		ApiURL: serverURL + "/",
		Token:  "test-token",
	})
}

func TestRequest_NoToken(t *testing.T) {
	c := New(&config.LidlConfig{ApiURL: "https://example.invalid/"})
	_, err := c.request("promotionslist", "GET")
	if err == nil {
		t.Fatal("expected error when token is not set")
	}
}

func TestActivatePromotion_Expired(t *testing.T) {
	c := New(&config.LidlConfig{ApiURL: "https://example.invalid/"})
	p := Promotion{
		ID:                "1",
		StartValidityDate: "2000-01-01T00:00:00Z",
		EndValidityDate:   "2000-01-02T00:00:00Z",
	}
	if err := c.activatePromotion(p); err != ErrPromotionExpired {
		t.Errorf("got %v, want %v", err, ErrPromotionExpired)
	}
}

func TestActivatePromotion_NotYetValid(t *testing.T) {
	c := New(&config.LidlConfig{ApiURL: "https://example.invalid/"})
	future := time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z")
	p := Promotion{
		ID:                "1",
		StartValidityDate: future,
		EndValidityDate:   future,
	}
	if err := c.activatePromotion(p); err != ErrPromotionNotYetValid {
		t.Errorf("got %v, want %v", err, ErrPromotionNotYetValid)
	}
}

func TestActivatePromotion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if got := r.URL.Path; got != "/promotions/42/activation" {
			t.Errorf("unexpected path: %s", got)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL)
	past := time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05Z")
	future := time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z")
	p := Promotion{ID: "42", StartValidityDate: past, EndValidityDate: future}

	if err := c.activatePromotion(p); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestActivatePromotion_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL)
	past := time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05Z")
	future := time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z")
	p := Promotion{ID: "42", StartValidityDate: past, EndValidityDate: future}

	if err := c.activatePromotion(p); err == nil {
		t.Error("expected error for non-202 status code")
	}
}

func TestGetPromotions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Cookie"); got != "authToken=test-token" {
			t.Errorf("unexpected Cookie header: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"sections": [
				{"promotions": [{"id": "1", "title": "Promo 1"}]},
				{"promotions": [{"id": "2", "title": "Promo 2"}]}
			]
		}`))
	}))
	defer server.Close()

	c := newTestClient(t, server.URL)
	promotions, err := c.getPromotions()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(promotions) != 2 {
		t.Fatalf("got %d promotions, want 2", len(promotions))
	}
	if promotions[0].ID != "1" || promotions[1].ID != "2" {
		t.Errorf("unexpected promotion IDs: %+v", promotions)
	}
}

func TestGetPromotions_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL)
	if _, err := c.getPromotions(); err == nil {
		t.Error("expected error for non-200 status code")
	}
}
