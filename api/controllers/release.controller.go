package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/welschmorgan/go-release-manager/release"
)

type ReleaseEndpoint struct {
	*RestEndpoint
}

func NewReleaseEndpoint() *ReleaseEndpoint {
	ret := &ReleaseEndpoint{}
	ret.RestEndpoint = NewRestEndpoint(map[string]http.HandlerFunc{
		http.MethodPost:   ret.Create,
		http.MethodGet:    ret.List,
		http.MethodPut:    ret.Update,
		http.MethodDelete: ret.Delete,
	})
	return ret
}

func (e *ReleaseEndpoint) Create(w http.ResponseWriter, r *http.Request) {
	data := []byte{}
	params := map[string]interface{}{}
	if _, err := r.Body.Read(data); err != nil {
		e.Error(w, http.StatusBadRequest, "NO_BODY", "Missing 'type' field (values: major, minor, build, revision, preRelease, buildMetaTag)\n"+err.Error())
	} else if err = json.Unmarshal(data, &params); err != nil {
		e.Error(w, http.StatusBadRequest, "INVALID_BODY", "Couldn't deserialize body into map, "+err.Error())
	} else if mode, ok := params["type"]; !ok {
		e.Error(w, http.StatusBadRequest, "MISSING_FIELD", "Missing 'type' field (values: major, minor, build, revision, preRelease, buildMetaTag)\n"+err.Error())
	} else if err = release.DoRelease(mode.(string)); err != nil {
		e.Error(w, http.StatusBadRequest, "", "Failed to create release, "+err.Error())
	}
}

func (e *ReleaseEndpoint) Read(w http.ResponseWriter, r *http.Request)   {}
func (e *ReleaseEndpoint) Update(w http.ResponseWriter, r *http.Request) {}
func (e *ReleaseEndpoint) Delete(w http.ResponseWriter, r *http.Request) {}
func (e *ReleaseEndpoint) List(w http.ResponseWriter, r *http.Request) {
	if releases, err := release.LoadReleaseTrain(); err != nil {
		e.Error(w, 500, "", err.Error())
	} else {
		e.Json(w, 200, releases)
	}
}
