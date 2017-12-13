package main
// http://hackmd.diverse-team.fr/s/HkPqyy8-G#
import (
    "fmt"
    "io/ioutil"
    "time"
    "net/http"
    "encoding/json"
    "github.com/apaxa-go/eval"
)


type JWT struct {
    Sub string `json:"sub"`
    Iat int64 `json:"iat"`
    Exp int64 `json:"exp"`
}

type MSG struct {
    Login string `json:"login,omitempty"`
    Pass string `json:"pass,omitempty"`
    Token *JWT `json:"token,omitempty"`
    Expression string `json:"expression,omitempty"`
    Result interface{} `json:"result,omitempty"`
}

// High security level
var login, pwd string = "admin", "admin"

// POST JSON with login and pass
// Return JSON with JWT as token
func loginHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var req MSG
    err := json.Unmarshal(body, &req)
    if(req.Login == login && req.Pass == pwd){
      now := time.Now().Unix()
      token := JWT{req.Login, now, now+3600}
      res, _ := json.Marshal(MSG{Token:&token})
      w.Write(res)
    } else if(err!=nil){
      w.WriteHeader(http.StatusBadRequest)
      fmt.Fprintf(w, `{"fault":%s}`, err)
    } else {
      w.WriteHeader(http.StatusUnauthorized)
      fmt.Fprintf(w, `{"fault":"error invalid login and password"}`)
    }
}

// POST JSON with valid token and expression
// Return JSON with result
func computeHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var req MSG
    err := json.Unmarshal(body, &req)
    if(req.Token.Sub == "admin" && req.Token.Exp >= time.Now().Unix() && req.Token.Iat <= time.Now().Unix()){
      expr,err:=eval.ParseString(req.Expression,"")
      r,err:=expr.EvalToInterface(nil)
      if(err==nil){
        res, _ := json.Marshal(MSG{Result:r})
        //w.WriteHeader(http.StatusOK)
        w.Write(res)
        //fmt.Fprintf(w, `%s`, res)
      } else if(err!=nil){
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, `{"fault":%s}`, err)
      }
    } else if(err!=nil){
      w.WriteHeader(http.StatusBadRequest)
      fmt.Fprintf(w, `{"fault":%s}`, err)
    } else {
      w.WriteHeader(http.StatusUnauthorized)
      fmt.Fprintf(w, `{"fault":"invalid token"}`)
    }

}

func main() {
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/compute", computeHandler)
    http.ListenAndServe(":8080", nil)
}
