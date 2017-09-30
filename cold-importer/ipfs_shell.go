package main

import (
	"log"
	"net/http"
	"time"
)

type ipfsShell struct {
	url     string
	httpcli *http.Client
}

func ipfsInitStart(ipfsrpc string) *ipfsShell {

	c := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	for ipfsIsUp() {
		log.Println("ipfs daemon not running, waiting a second...")
		time.Sleep(1000 * time.Millisecond)
	}

	return &ipfsShell{
		url:     ipfsrpc,
		httpcli: c,
	}
}

func ipfsIsUp() bool {
	// PLACEHOLDER
	return false
	// PLACEHOLDER
}

/*
func NewRequest(ctx context.Context, url, command string, args ...string) *Request {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	opts := map[string]string{
		"encoding":        "json",
		"stream-channels": "true",
	}
	return &Request{
		ApiBase: url + "/api/v0",
		Command: command,
		Args:    args,
		Opts:    opts,
		Headers: make(map[string]string),
	}
}

func (r *Request) Send(c *http.Client) (*Response, error) {
	url := r.getURL()

	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		return nil, err
	}

	if fr, ok := r.Body.(*files.MultiFileReader); ok {
		req.Header.Set("Content-Type", "multipart/form-data; boundary="+fr.Boundary())
		req.Header.Set("Content-Disposition", "form-data: name=\"files\"")
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	parts := strings.Split(contentType, ";")
	contentType = parts[0]

	nresp := new(Response)

	nresp.Output = resp.Body
	if resp.StatusCode >= http.StatusBadRequest {
		e := &Error{
			Command: r.Command,
		}
		switch {
		case resp.StatusCode == http.StatusNotFound:
			e.Message = "command not found"
		case contentType == "text/plain":
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ipfs-shell: warning! response read error: %s\n", err)
			}
			e.Message = string(out)
		case contentType == "application/json":
			if err = json.NewDecoder(resp.Body).Decode(e); err != nil {
				fmt.Fprintf(os.Stderr, "ipfs-shell: warning! response unmarshall error: %s\n", err)
			}
		default:
			fmt.Fprintf(os.Stderr, "ipfs-shell: warning! unhandled response encoding: %s", contentType)
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ipfs-shell: response read error: %s\n", err)
			}
			e.Message = fmt.Sprintf("unknown ipfs-shell error encoding: %q - %q", contentType, out)
		}
		nresp.Error = e
		nresp.Output = nil

		// drain body and close
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	}

	return nresp, nil
}

func () DagPutPin(data interface{}, ienc, kind string) (string, error) {
	req := api.NewRequest(context.Background(), "localhost:5001", "dag/put")
	req.Opts = map[string]string{
		"input-enc": ienc,
		"format":    kind,
	}

	var r io.Reader
	switch data := data.(type) {
	case string:
		r = strings.NewReader(data)
	case []byte:
		r = bytes.NewReader(data)
	case io.Reader:
		r = data
	default:
		return "", fmt.Errorf("cannot current handle putting values of type %T", data)
	}
	rc := ioutil.NopCloser(r)
	fr := files.NewReaderFile("", "", rc, nil)
	slf := files.NewSliceFile("", "", []files.File{fr})
	fileReader := files.NewMultiFileReader(slf, true)
	req.Body = fileReader

	resp, err := req.Send(s.httpcli)
	if err != nil {
		return "", err
	}
	defer resp.Close()

	if resp.Error != nil {
		return "", resp.Error
	}

	var out struct {
		Cid struct {
			Target string `json:"/"`
		}
	}
	err = json.NewDecoder(resp.Output).Decode(&out)
	if err != nil {
		return "", err
	}

	return out.Cid.Target, nil
}
*/
