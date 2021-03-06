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
	"fmt"
	"net/http"
)

// Services references all the information needed for every method.
type Services struct {
	config Configuration // Configuration used by the server.
	store  Store         // Instance of storage service for node data
	access Access        // Instance of access service
}

// NewServices a set of services to use with Shared Web State. These provide
// defaults via the configuration parameter, and access to persistent storage
// via the store parameter.
func NewServices(
	config Configuration,
	store Store,
	access Access) *Services {
	var s Services
	s.config = config
	s.store = store
	s.access = access
	return &s
}

// Config returns the configuration service.
func (s *Services) Config() *Configuration { return &s.config }

// GetCreator returns the store service
func (s *Services) GetCreator(host string) (*Creator, error) {
	return s.store.GetCreator(host)
}

// Returns true if the request is allowed to access the handler, otherwise false.
// If false is returned then no further action is needed as the method will have
// responded to the request already.
func (s *Services) getAccessAllowed(w http.ResponseWriter, r *http.Request) bool {
	err := r.ParseForm()
	if err != nil {
		returnAPIError(s, w, err, http.StatusInternalServerError)
		return false
	}
	v, err := s.access.GetAllowed(r.FormValue("accesskey"))
	if v == false || err != nil {
		returnAPIError(
			s,
			w,
			fmt.Errorf("Access denied"),
			http.StatusNetworkAuthenticationRequired)
		return false
	}
	return true
}
