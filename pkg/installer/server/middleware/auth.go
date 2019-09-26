/*
 * Copyright 2019 gosoon.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package middleware

import "net/http"

// Authenticate will create a authenticate middleware
// TODO(user): Modify this function to implement your logic.
func Authenticate(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
					// This is a example.
			        token := r.Header.Get("Authorization")
			        if len(token) != 0 {
			            bearerValue := strings.Split(token, " ")[1]
			            // token in config
			            if bearerValue == viper.GetString(config.token) {
			                next.ServeHTTP(w, r)
			            }
			        }
			        controller.Unauthorized(w, r, fmt.Sprintf("Authenticate failed,plz check your token."))
		*/
		next.ServeHTTP(w, r)
	}
}
