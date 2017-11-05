/*
    This is where your secret config info will go.
*/
package config

import (
    "net/http"
)



var (
  Secret = &SecretConfig{
    URL: "http://example.com",
    Data: "Hello!",
  }
  CatSecret *SecretConfig
)

/*
  //init will always run on bootup for you noobies
  func init() {
      //calculate config variables at start if you're feeling crazy
      //maybe make a map string -> string, then a header
      h := http.Header{}
      for k, v := range header {
          h.Add(k, v)
      }

      //and why not a cookie
      cookie := &http.Cookie {
          Name: "cool",
          Value: "guy",
      }

      //set new config with request, if you don't do this then fetch won't work
      Secret = &SecretConfig{
          URL: "dingo",
          Data: "dango",
          Header: h,
          Cookie: cookie,
      }
  }
*/

type SecretConfig struct {
    URL      string
    Header   http.Header
    Cookie   *http.Cookie
    Data     string
}
