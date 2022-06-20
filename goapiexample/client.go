package goapiexample

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

type HttpClient struct {
	baseURL *url.URL
	client  *http.Client
	signer  *v4.Signer
	cfg     aws.Config
	now     func() time.Time
}

const emptyHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

func (c HttpClient) buildSignedRequest(ctx context.Context, method string, url url.URL, body interface{}) (r *http.Request, err error) {
	// If there is no payload, default to an empty buffer and the empty hash.
	hash := emptyHash
	payloadBuffer := new(bytes.Buffer)
	if body != nil {
		// Hasher to compute body hash.
		bodyHasher := sha256.New()

		// JSON encode and calculate the hash concurrently.
		err = json.NewEncoder(io.MultiWriter(bodyHasher, payloadBuffer)).Encode(body)
		if err != nil {
			err = fmt.Errorf("goexample-api client: failed to encode request: %w", err)
			return
		}
		hash = hex.EncodeToString(bodyHasher.Sum(nil))
	}

	// Create the request.
	r, err = http.NewRequestWithContext(ctx, method, url.String(), payloadBuffer)
	if err != nil {
		err = fmt.Errorf("goexample-api client: failed to create request: %w", err)
		return
	}

	// Get signing credentials.
	creds, err := c.cfg.Credentials.Retrieve(ctx)
	if err != nil {
		err = fmt.Errorf("goexample-api client: failed to retrieve creds: %w", err)
		return
	}

	// Sign the request.
	err = c.signer.SignHTTP(ctx, creds, r, hash, "execute-api", "eu-west-2", c.now())
	if err != nil {
		err = fmt.Errorf("goexample-api client: failed to sign request: %w", err)
		return
	}
	return
}

func NewClient(baseUrl string) (c HttpClient, err error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return
	}
	c.baseURL = u
	c.client = &http.Client{
		Timeout: 15 * time.Second,
	}
	c.cfg, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		return
	}
	c.signer = v4.NewSigner()
	c.now = time.Now
	return
}

func (c HttpClient) MathsAdd(ctx context.Context, a, b int) (result int, ok bool, err error) {
	resource := *c.baseURL
	resource.Path = path.Join(resource.Path, "maths", "add")
	req, err := c.buildSignedRequest(ctx, http.MethodGet, resource, nil)
	if err != nil {
		return
	}
	res, err := c.client.Do(req)
	if err != nil {
		return
	}
	if statusOK := res.StatusCode >= 200 && res.StatusCode < 300; !statusOK {
		body, _ := ioutil.ReadAll(res.Body)
		err = fmt.Errorf("api failed to respond with a 2xx status code, got: %d. Body: %s", res.StatusCode, string(body))
		return
	}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return
	}
	ok = true
	return
}
