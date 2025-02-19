package main
import "net/http"

// setup cors middleware to  be open for telex
func CorsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// recover application panic error... and return 500 server response
// func recoverpanic(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		defer func() {

// 			if err := recover(); err != nil {
// 				w.Header().Set("Connection", "close")
// 				// send a 500 Internal server error to the user
// 				app.serverErrorResponse(w, err)
// 			}
// 		}()
// 		next.ServeHTTP(w, r)
// 	})
// }
