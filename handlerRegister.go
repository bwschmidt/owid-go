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

import "net/http"

// HandlerRegister - Handler for the registering of a domain.
func HandlerRegister(s *Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var d Register
		d.Services = s
		d.Domain = r.Host
		d.Name = ""

		// Check that the domain has not already been registered.
		n, err := s.store.GetCreator(r.Host)
		if err != nil {
			returnServerError(s, w, err)
			return
		}
		if n != nil {
			return
		}

		// Get any values from the form.
		err = r.ParseForm()
		if err != nil {
			returnServerError(s, w, err)
			return
		}
		d.DisplayErrors = len(r.Form) > 0

		// Get the network information.
		d.Name = r.FormValue("name")
		if len(d.Name) <= 5 {
			d.NameError = "Name must be longer than 5 characters"
		} else if len(d.Name) > 20 {
			d.NameError = "Name can not be longer than 20 characters"
		}

		// If the form data is valid then store the new node.
		if d.NameError == "" {
			storeCreator(s, &d)
		}

		// Return the HTML page.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		err = registerTemplate.Execute(w, &d)
		if err != nil {
			returnServerError(s, w, err)
		}
	}
}

func storeCreator(s *Services, d *Register) {

	// Create the new node ready to have it's secret added and stored.
	cry, err := NewCrypto()
	c := newCreator(
		d.Domain,
		cry.privateKeyToPemString(),
		cry.publicKeyToPemString(),
		d.Name)
	if err != nil {
		d.Error = err.Error()
		return
	}

	// Store the node and it successful mark the registration process as
	// complete.
	err = s.store.setCreator(c)
	if err != nil {
		d.Error = err.Error()
	} else {
		d.ReadOnly = true
	}
}
