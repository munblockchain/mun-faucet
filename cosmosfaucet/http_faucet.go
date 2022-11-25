package cosmosfaucet

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"faucet/cmd/config"

	"github.com/ignite/cli/ignite/pkg/xhttp"
)

type TransferRequest struct {
	// AccountAddress to request for coins.
	AccountAddress string `json:"address"`
}

func NewTransferRequest(accountAddress string, coins []string) TransferRequest {
	return TransferRequest{
		AccountAddress: accountAddress,
	}
}

type TransferResponse struct {
	Error string `json:"error,omitempty"`
}

func (f Faucet) faucetHandler(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	cookie_captcha, err := r.Cookie("response")

	if err != nil {
		fmt.Println("Invalid request")
		responseError(w, http.StatusBadRequest, err)
		return
	}

	referer := r.Referer()
	origin := r.Header.Get("Origin") + "/"

	if referer != origin || referer != f.configuration.Referer {
		fmt.Println("Invalid request")
		responseError(w, http.StatusBadRequest, err)
		return
	}

	err = f.captcha.Verify(cookie_captcha.Value)
	if err != nil {
		fmt.Println("Invalid reCaptha")
		responseError(w, http.StatusBadRequest, err)
		return
	}

	// decode request into req.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	// try performing the transfer
	if errCode, err := f.Transfer(r.Context(), req.AccountAddress); err != nil {
		if err == context.Canceled {
			return
		}
		responseError(w, (int)(errCode), err)
	} else {
		responseSuccess(w)
	}
}

func (f Faucet) validateReCAPTCHA(recaptchaResponse string) (bool, error) {
	// Check this URL verification details from Google
	// https://developers.google.com/recaptcha/docs/verify
	req, err := http.PostForm(f.configuration.ReCAPTCHA_VerifyURL, url.Values{
		"secret":   {f.configuration.ReCAPTCHA_ServerKey},
		"response": {recaptchaResponse},
	})
	if err != nil { // Handle error from HTTP POST to Google reCAPTCHA verify server
		return false, err
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body) // Read the response from Google
	if err != nil {
		return false, err
	}

	var googleResponse config.GoogleRecaptchaResponse
	err = json.Unmarshal(body, &googleResponse) // Parse the JSON response from Google
	if err != nil {
		return false, err
	}
	return true, nil
}

// FaucetInfoResponse is the faucet info payload.
type FaucetInfoResponse struct {
	// IsAFaucet indicates that this is a faucet endpoint.
	// useful for auto discoveries.
	IsAFaucet bool `json:"is_a_faucet"`

	// ChainID is chain id of the chain that faucet is running for.
	ChainID string `json:"chain_id"`
}

func (f Faucet) faucetInfoHandler(w http.ResponseWriter, r *http.Request) {
	xhttp.ResponseJSON(w, http.StatusOK, FaucetInfoResponse{
		IsAFaucet: true,
		ChainID:   f.chainID,
	})
}

func responseSuccess(w http.ResponseWriter) {
	xhttp.ResponseJSON(w, http.StatusOK, TransferResponse{})
}

func responseError(w http.ResponseWriter, code int, err error) {
	xhttp.ResponseJSON(w, code, TransferResponse{
		Error: err.Error(),
	})
}
