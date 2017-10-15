

createreactappmanifest integrate Create React App with Go web application




* [Type Config](#type-config)
* [Type Manifest](#type-manifest)
  * [New](#manifest-new)
  * [Get UR L](#manifest-get-ur-l)
  * [Mount](#manifest-mount)






## Type: Config
``` go
type Config struct {
    // the PUBLIC_URL environment variable when do yarn build in react app
    PublicURL   string
    ManifestDir string

    // sometime you want to excludes the files in react app public folder.
    MountExcludeForPublic string
    IsDev                 bool
    DevBundleURL          string
}
```









## Type: Manifest
``` go
type Manifest struct {
    // contains filtered or unexported fields
}
```






### Manifest: New
``` go
func New(cfg *Config) (m *Manifest, err error)
```

Mount expose assets with ServeMux, and GetURL get correct assets path for you.
```go
	mux := http.DefaultServeMux
	
	m, _ := manifest.New(&manifest.Config{
	    ManifestDir: "./example/build",
	    PublicURL:   "/cms",
	})
	
	renderTemplate := func() string {
	    buf := bytes.NewBuffer(nil)
	    fmt.Fprintf(buf, `<link href="%s" rel="stylesheet">`, m.GetURL("main.css"))
	    fmt.Fprintf(buf, `<script type="text/javascript" src="%s"></script>`, m.GetURL("main.js"))
	    return buf.String()
	}
	
	m.Mount(mux)
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	    fmt.Fprintf(w, renderTemplate())
	})
```



### Manifest: Get UR L
``` go
func (m *Manifest) GetURL(name string) (url string)
```
GetURL get dynamic compiled url like `/static/css/main.c17080f1.css` with simple name like `main.css` to be used in views




### Manifest: Mount
``` go
func (m *Manifest) Mount(mux *http.ServeMux)
```
Mount automatically mounts Create React App build directory into Go ServeMux





