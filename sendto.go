package sendto

import (
  "context"
  "embed"
  "errors"
  "flag"
  "fmt"
  "html/template"
  "io/fs"
  "log"
  "net"
  "net/http"
  "strings"

  "tailscale.com/client/tailscale"
  "tailscale.com/hostinfo"
  "tailscale.com/ipn"
  "tailscale.com/tsnet"

  "github.com/frankh/sendto/pkg/api"
  "github.com/frankh/sendto/pkg/db"
)

const defaultHostname = "sendto"

var (
  verbose    = flag.Bool("verbose", false, "be verbose")
  controlURL = flag.String("control-url", ipn.DefaultControlURL, "the URL base of the control plane (i.e. coordination server)")
  sqlitefile = flag.String("sqlitedb", "", "path of SQLite database to store files and messages")
  dev        = flag.String("dev-listen", "", "if non-empty, listen on this addr and run in dev mode; auto-set sqlitedb if empty and don't use tsnet")
  snapshot   = flag.String("snapshot", "", "file path of snapshot file")
  hostname   = flag.String("hostname", defaultHostname, "service name")
)

//go:embed frontend/public
var embeddedFS embed.FS
var staticFS fs.FS

var localClient *tailscale.LocalClient
var database db.DB

func init() {
  var err error
  staticFS, err = fs.Sub(embeddedFS, "frontend/public")
  if err != nil {
    panic(err)
  }
  homeTmpl = template.Must(template.ParseFS(staticFS, "index.html"))
  receiveTmpl = template.Must(template.ParseFS(staticFS, "receive.html"))
}

func Run() error {
  flag.Parse()

  hostinfo.SetApp("sendto")

  database = db.NewSqliteDB("sqlite.db")
  sendHandler := api.SendHandler{database}

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
  http.HandleFunc("/api/send", withUser(sendHandler.Serve))
  http.HandleFunc("/r/", withUser(serveReceive))
  http.HandleFunc("/", serveSendto)

  if *dev != "" {
    // override default hostname for dev mode
    if *hostname == defaultHostname {
      if h, p, err := net.SplitHostPort(*dev); err == nil {
        if h == "" {
          h = "localhost"
        }
        *hostname = fmt.Sprintf("%s:%s", h, p)
      }
    }

    log.Printf("Running in dev mode on %s ...", *dev)
    log.Fatal(http.ListenAndServeTLS(*dev, "localhost.crt", "localhost.key", nil))
  }

  if *hostname == "" {
    return errors.New("--hostname, if specified, cannot be empty")
  }

  srv := &tsnet.Server{
    ControlURL: *controlURL,
    Hostname:   *hostname,
    Logf: func(format string, args ...any) {
      // Show the log line with the interactive tailscale login link even when verbose is off
      if strings.Contains(format, "To start this tsnet server") {
        log.Printf(format, args...)
      }
    },
  }
  if *verbose {
    srv.Logf = log.Printf
  }
  if err := srv.Start(); err != nil {
    return err
  }
  localClient, _ = srv.LocalClient()

  l80, err := srv.Listen("tcp", ":80")
  if err != nil {
    return err
  }

  log.Printf("Serving http://%s/ ...", *hostname)
  if err := http.Serve(l80, nil); err != nil {
    return err
  }
  return nil
}

var (
  homeTmpl    *template.Template
  receiveTmpl *template.Template
)

func serveSendto(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "HTTP Method Unsupported", http.StatusBadRequest)
    return
  }

  to := strings.TrimPrefix(r.URL.Path, "/")
  if strings.Contains(to, "/") {
    http.Error(w, "Not found", http.StatusNotFound)
    return
  }
  to = findUser(to)

  homeTmpl.Execute(w, struct{ To string }{To: to})
}

func findUser(to string) string {
  st, err := localClient.Status(context.Background())
  if err != nil {
    return ""
  }
  for _, user := range st.User {
    if strings.Split(user.LoginName, "@")[0] == to {
      return user.LoginName
    }
  }

  return ""
}

func serveReceive(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "HTTP Method Unsupported", http.StatusBadRequest)
    return
  }
  id := strings.TrimPrefix(r.URL.Path, "/r/")
  message, err := database.Load(id)
  if err != nil {
    http.Error(w, "Not Found", http.StatusNotFound)
    return
  }

  login, err := currentUser(r)
  if err != nil || message.To != login {
    http.Error(w, "Wrong user", http.StatusNotFound)
    return
  }
  receiveTmpl.Execute(w, message)
}

// acceptHTML returns whether the request can accept a text/html response.
func acceptHTML(r *http.Request) bool {
  return strings.Contains(strings.ToLower(r.Header.Get("Accept")), "text/html")
}

func devMode() bool { return *dev != "" }

func currentUser(r *http.Request) (string, error) {
  login := ""
  if devMode() {
    login = "foo@example.com"
  } else {
    res, err := localClient.WhoIs(r.Context(), r.RemoteAddr)
    if err != nil {
      return "", err
    }
    login = res.UserProfile.LoginName
  }
  return login, nil

}

func withUser(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    login, err := currentUser(r)
    if err != nil {
      http.Error(w, "Bad login information", http.StatusBadRequest)
      return
    }
    next(w, r.WithContext(context.WithValue(r.Context(), "login", login)))
  })
}
