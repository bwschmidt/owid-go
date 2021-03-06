/* ****************************************************************************
 * Copyright 2020 51 Degrees Mobile Experts Limited (51degrees.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 * use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 * ***************************************************************************/

package owid

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	s, err := getServices()
	if err != nil {
		t.Fatal(err)
	}

	data := url.Values{}
	data.Set("name", testOrgName)

	req, err := http.NewRequest(
		"POST",
		"/owid/api/v1/register",
		strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Host = testDomain
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	setCommon(req)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerRegister(s))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}
}
func TestCreatorHandler(t *testing.T) {
	s, err := getServices()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/owid/api/v1/creator", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = testDomain
	setCommon(req)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlerCreator(s))

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := testCreatorJSON
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got '%v' want '%v'",
			rr.Body.String(), expected)
	}

}

func setCommon(r *http.Request) {
	q := r.URL.Query()

	// set the access key
	q.Set("accessKey", "key1")

	r.URL.RawQuery = q.Encode()
}

func getServices() (*Services, error) {
	c := NewConfig("appsettings.Test.json")
	a := NewAccessSimple([]string{"key1", "key2"})
	s, err := newTestStore()
	if err != nil {
		return nil, err
	}
	return NewServices(c, s, a), nil
}
