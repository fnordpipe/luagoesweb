package luahttp

import (
  "net/http"
  "strings"

  "github.com/gorilla/mux"
  "github.com/yuin/gopher-lua"
  "metagit.org/fnordpipe/gorillua/modules/logger"
)

type RouterInfo struct {
  Method string
  Context string
  Callback lua.LValue
}

var m = map[string]lua.LGFunction{
  "serve": serve,
}

func Loader(L *lua.LState) int {
  module := L.SetFuncs(L.NewTable(), m)
  L.Push(module)
  return 1
}

func handleRequest(L *lua.LState, ctx RouterInfo, w http.ResponseWriter, r *http.Request) {
  var _w = map[string]lua.LGFunction{
    "addHeader": func(L *lua.LState) int {
      key := L.CheckString(1)
      value := L.CheckString(2)
      w.Header().Add(key, value)
      return 0
    },
    "setCookie": func(L *lua.LState) int {
      name := L.CheckString(1)
      value := L.CheckString(2)
      path := L.CheckString(3)
      maxage := L.CheckNumber(4)
      httponly := L.CheckBool(5)

      cookie := http.Cookie{
        Name: name,
        Value: value,
        Path: path,
        MaxAge: int(lua.LVAsNumber(maxage)),
        HttpOnly: httponly,
      }

      http.SetCookie(w, &cookie)
      return 0
    },
    "setStatus": func(L *lua.LState) int {
      status := L.CheckNumber(1)
      w.WriteHeader(int(status))
      return 0
    },
    "write": func(L *lua.LState) int {
      content := L.CheckString(1)
      w.Write([]byte(content))
      return 0
    },
  }

  var _r = map[string]lua.LGFunction{
    "getCookie": func(L *lua.LState) int {
      name := L.CheckString(1)
      cookie, err := r.Cookie(name)
      if err != nil {
        L.Push(lua.LNil)
        L.Push(lua.LString(err.Error()))
        return 2
      }

      L.Push(lua.LString(cookie.Value))
      return 1
    },
    "getHeader": func(L *lua.LState) int {
      key := L.CheckString(1)
      header := r.Header.Get(key)
      L.Push(lua.LString(header))
      return 1
    },
    "parseForm": func(L *lua.LState) int {
      r.ParseForm()
      if len(r.Form) <= 0 {
        L.Push(lua.LNil)
        L.Push(lua.LString("no form data to parse"))
        return 2
      }

      t := L.CreateTable(0, len(r.Form))
      for k, v := range r.Form {
        t.RawSetH(lua.LString(k), lua.LString(strings.Join(v, "")))
      }
      L.Push(t)
      return 1
    },
  }

  mw := L.SetFuncs(L.NewTable(), _w)
  mr := L.SetFuncs(L.NewTable(), _r)
  L.Push(ctx.Callback)
  L.Push(mw)
  L.Push(mr)
  err := L.PCall(2, 0, nil)
  if err != nil {
    logger.Debug(err.Error())
  }

  logger.Stdout(ctx.Context)
  return
}

func serve(L *lua.LState) int {
  address := L.CheckString(1)
  lrouter := L.CheckTable(2)
  lstatic := L.CheckAny(3)
  router := mux.NewRouter()
  var routes []RouterInfo

  lrouter.ForEach(func(k, v lua.LValue) {
    var route RouterInfo
    switch lv := v.(type) {
      case *lua.LTable:
        lv.ForEach(func(k, v lua.LValue) {
          if k.String() == "method" { route.Method = v.String() }
          if k.String() == "context" { route.Context = v.String() }
          if k.String() == "callback" { route.Callback = v }
        })
        routes = append(routes, route)
    }
  })

  for k, _ := range routes {
    j := k
    router.HandleFunc(routes[j].Context, func(w http.ResponseWriter, r *http.Request) {
      handleRequest(L, routes[j], w, r)
    }).Methods(routes[j].Method)
  }

  switch lv := lstatic.(type) {
    case lua.LString:
      router.PathPrefix("/").Handler(
        http.StripPrefix("/static/", http.FileServer(http.Dir(lv.String()))))
  }

  http.ListenAndServe(address, router)

  return 0
}
