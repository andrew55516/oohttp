package http

import (
	"context"
	"net/http"
)

const connCloserFuncKey = "connCloserFuncKey"

// StdlibTransport is an adapter for integrating net/http dependend code.
// It looks like an http.RoundTripper but uses this fork internally.
type StdlibTransport struct {
	*Transport
}

// RoundTrip implements the http.RoundTripper interface.
func (txp *StdlibTransport) RoundTrip(stdReq *http.Request) (*http.Response, error) {
	req := &Request{
		Method:           stdReq.Method,
		URL:              stdReq.URL,
		Proto:            stdReq.Proto,
		ProtoMajor:       stdReq.ProtoMajor,
		ProtoMinor:       stdReq.ProtoMinor,
		Header:           Header(stdReq.Header),
		Body:             stdReq.Body,
		GetBody:          stdReq.GetBody,
		ContentLength:    stdReq.ContentLength,
		TransferEncoding: stdReq.TransferEncoding,
		Close:            stdReq.Close,
		Host:             stdReq.Host,
		Form:             stdReq.Form,
		PostForm:         stdReq.PostForm,
		MultipartForm:    stdReq.MultipartForm,
		Trailer:          Header(stdReq.Trailer),
		RemoteAddr:       stdReq.RemoteAddr,
		RequestURI:       stdReq.RequestURI,
		TLS:              stdReq.TLS,
		Cancel:           stdReq.Cancel,
		Response:         nil, // cannot assign this field
		ctx:              stdReq.Context(),
	}
	resp, err := txp.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	stdResp := &http.Response{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMinor:       resp.ProtoMinor,
		ProtoMajor:       resp.ProtoMajor,
		Header:           http.Header(resp.Header),
		Body:             resp.Body,
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Close:            resp.Close,
		Uncompressed:     resp.Uncompressed,
		Trailer:          http.Header(resp.Trailer),
		Request:          stdReq.WithContext(context.WithValue(stdReq.Context(), connCloserFuncKey, resp.ConnCloserFunc)),
		TLS:              resp.TLS,
	}
	return stdResp, nil
}

func CloseUsedConnection(resp *http.Response, err error) {
	if resp != nil && resp.Request.Context().Value(connCloserFuncKey) != nil {
		resp.Request.Context().Value(connCloserFuncKey).(func(err error))(err)
	}
}
