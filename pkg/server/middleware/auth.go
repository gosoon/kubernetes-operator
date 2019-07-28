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

// Authenticate is check client source
//func Authenticate(next http.Handler) http.HandlerFunc {
//return func(w http.ResponseWriter, r *http.Request) {
//token := r.Header.Get("Authorization")
//// check token
//if len(token) != 0 {
//bearerValue := strings.Split(token, " ")[1]
//userInfo, err := Db.GetUserByToken(bearerValue)
//if err != nil {
//glog.Errorf("get user info by token with err [%v]", err)
//public.FailedResponse(w, r, fmt.Sprintf("Authenticate failed,plz check your token."),
//httpcode.StatusUnauthorized)
//} else {
//// set username and role in cookie
//expiration := time.Now().Add(24 * time.Hour)
//r.AddCookie(&http.Cookie{Name: "username", Value: userInfo.Name, Expires: expiration})
//r.AddCookie(&http.Cookie{Name: "role", Value: userInfo.Role, Expires: expiration})
//valueCtx := context.WithValue(r.Context(), "opUser", userInfo.Name)
//r = r.WithContext(valueCtx)
//next.ServeHTTP(w, r)
//}
//}
//}
//}
