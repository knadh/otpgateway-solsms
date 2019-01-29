package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/knadh/otpgateway"
)

const (
	providerID  = "solsms"
	channelName = "SMS"
	maxOTPlen   = 6
	apiURL      = "http://alerts.solutionsinfini.co/apiv2/"
)

var reNum = regexp.MustCompile(`\+?([0-9]){8,15}`)

// sms is the default representation of the sms interface.
type sms struct {
	cfg *cfg
	h   *http.Client
}

type cfg struct {
	RootURL      string `json:"RootURL"`
	APIKey       string `json:"APIKey"`
	Sender       string `json:"Sender"`
	Timeout      int    `json:"Timeout"`
	MaxIdleConns int    `json:"MaxIdleConns"`
}

// New returns an instance of the SMS package. cfg is configuration
// represented as a JSON string. Supported options are.
// {
// 	RootURL: "", // Optional root URL of the API,
// 	APIKey: "", // API Key,
// 	Sender: "", // Sender name
// 	Timeout: 5 // Optional HTTP timeout in seconds
// }
func New(jsonCfg []byte) (otpgateway.Provider, error) {
	var c *cfg
	if err := json.Unmarshal(jsonCfg, &c); err != nil {
		return nil, err
	}
	if c.RootURL == "" {
		c.RootURL = apiURL
	}
	if c.APIKey == "" || c.Sender == "" {
		return nil, errors.New("invalid APIKey or Sender")
	}

	// Initialize the HTTP client.
	t := 5
	if c.Timeout != 0 {
		t = c.Timeout
	}
	h := &http.Client{
		Timeout: time.Duration(t) * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   1,
			ResponseHeaderTimeout: time.Second * time.Duration(t),
		},
	}

	return &sms{
		cfg: c,
		h:   h}, nil
}

// ID returns the Provider's ID.
func (s *sms) ID() string {
	return providerID
}

// ChannelName returns the Provider's name.
func (s *sms) ChannelName() string {
	return channelName
}

// Description returns help text for the SMS verification Provider.
func (s *sms) Description() string {
	return fmt.Sprintf(`
		We've sent a %d digit code in an SMS to your phone.
		Enter it here to verify your phone number.`, maxOTPlen)
}

// ValidateAddress "validates" a phone number.
func (s *sms) ValidateAddress(to string) error {
	if !reNum.MatchString(to) {
		return errors.New("invalid phone number")
	}
	return nil
}

// Push pushes out an SMS.
func (s *sms) Push(to, subject string, body []byte) error {
	p := url.Values{}
	p.Set("api", "http")
	p.Set("workingkey", s.cfg.APIKey)
	p.Set("sender", s.cfg.Sender)
	p.Set("to", to)
	p.Set("message", string(body))

	// Make the request.
	resp, err := s.h.PostForm(s.cfg.RootURL, p)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response.
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if !bytes.Contains(b, []byte("responsecode 200")) {
		return errors.New(string(b))
	}
	return nil
}

// MaxOTPLen returns the maximum allowed length of the OTP value.
func (s *sms) MaxOTPLen() int {
	return maxOTPlen
}

// MaxBodyLen returns the max permitted body size.
func (s *sms) MaxBodyLen() int {
	return 140
}
