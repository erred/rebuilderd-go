package status

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"text/template"

	"github.com/NYTimes/gziphandler"
	"go.seankhliao.com/rebuilderd-go"
	"go.seankhliao.com/usvc"
	"go.seankhliao.com/webstyle"
)

const (
	style = `
header{display: none;}
main{margin-top: 0;}
td {font-weight:700}
.BAD{color: #bf616a;}
.GOOD{color: #a3be8c;}
.UNKWN{color: #999;}
.QUEUED{color: #ebcb8b;}
`

	header = `
<h3><em>Arch Linux</em> Reproducible Builds</h3>
<p>
<a href="https://github.com/kpcyrd/rebuilderd">rebuilderd</a>
run by
<em><a href="https://seankhliao.com">seankhliao</a></em>,
page source: <a href="https://github.com/seankhliao/rebuilderd-go">github</a>
</p>`

	shortlog = `
<p><em><a href="#%s">%s</a></em>:
%d%% reproducible with
%d <span class="GOOD">good</span> /
%d <span class="BAD">bad</span> /
%d <span class="QUEUED">queued</span> /
%d <span class="UNKWN">unknown</span></p>
`

	subsection = `<h4 id="%s"><em>%s</em></h4>
<table>
<thead>
<tr>
<th>Status</th>
<th>Package</th>
<th>Version</th>
<th>Architecture</th>
<th>url</th>
</tr>
</thead>
</tbody>`

	trow = `
<tr class="%s">
<td><strong>%s</strong></td>
<td>%s</td>
<td>%s</td>
<td>%s</td>
<td><a href="%s">link</a></td>
</tr>
`
	ttail = `
</tbody>
</table>`
)

type Server struct {
	rc   *rebuilderd.Client
	t    *template.Template
	page map[string]string

	Svc *usvc.ServerSecure
}

func NewServer(args []string) *Server {
	var certFile, keyFile, endpoint, gaid, acf string
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.StringVar(&certFile, "cert", "/etc/letsencrypt/live/sne.seankhliao.com/fullchain.pem", "fullchain certificate file")
	fs.StringVar(&keyFile, "key", "/etc/letsencrypt/live/sne.seankhliao.com/privkey.pem", "private key file")
	fs.StringVar(&endpoint, "endpoint", "http://145.100.104.117:8484", "rebuilderd api endpoint")
	fs.StringVar(&gaid, "gaid", "UA-114337586-6", "google analytics id")
	fs.StringVar(&acf, "authcookiefile", "/var/lib/rebuilderd/auth-cookie", "file containing auth cookie")
	c := usvc.NewConfig(fs)
	fs.Parse(args[1:])

	svc, err := usvc.NewServerSecure(c, certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	rc, err := rebuilderd.NewClient(endpoint, nil)
	if err != nil {
		svc.Log.Fatal().Err(err).Msg("rebuilderd client")
	}

	b, err := ioutil.ReadFile(acf)
	if err != nil {
		svc.Log.Error().Str("file", acf).Err(err).Msg("read auth cookie file")
	}
	rc.AuthCookie = string(bytes.TrimSpace(b))

	s := &Server{
		rc: rc,
		t:  webstyle.Template,
		page: map[string]string{
			"Title":           "arch rebuilder | seankhliao",
			"Description":     "rebuilderd status for arch linux by seankhliao",
			"URLCanonical":    "https://rebuilder.seankhliao.com/",
			"Style":           style,
			"URLLogger":       "https://statslogger.seankhliao.com/form",
			"GoogleAnalytics": gaid,
		},
		Svc: svc,
	}
	s.Svc.Mux.HandleFunc("/favicon.ico", favicon)
	s.Svc.Mux.Handle("/", gziphandler.GzipHandler(s))
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=300")
	pkgs, err := s.rc.PkgsList(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	queue, err := s.rc.QueueList(rebuilderd.ListQueue{})
	if err != nil {
		s.Svc.Log.Error().Err(err).Msg("get queue")
	}

	page := make(map[string]string, len(s.page))
	for k, v := range s.page {
		page[k] = v
	}
	page["Main"] = pkgs2page(pkgs, queue.Queue)
	s.t.ExecuteTemplate(w, "LayoutGohtml", page)
}

func pkgs2page(pkgs []rebuilderd.PkgRelease, queue []rebuilderd.QueueItem) string {
	m := make(map[pkgver]rebuilderd.PkgRelease, len(pkgs)+len(queue))
	for _, p := range pkgs {
		m[newPkgver(p)] = p
	}
	for _, qi := range queue {
		qi.Package.Status = "QUEUED"
		m[newPkgver(qi.Package)] = qi.Package
	}

	packages := make([]rebuilderd.PkgRelease, len(m))
	for _, v := range m {
		packages = append(packages, v)
	}
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	reponames := []string{"Core", "Extra", "Community"}
	repos := make([]Repo, len(reponames))

	for _, p := range packages {
		s := fmt.Sprintf(trow, p.Status, p.Status, p.Name, p.Version, p.Architecture, p.URL)
		var i = -1
		switch p.Suite {
		case "core":
			i = 0
		case "extra":
			i = 1
		case "community":
			i = 2
		default:
			continue
		}
		repos[i].b.WriteString(s)
		repos[i].total++
		switch p.Status {
		case "GOOD":
			repos[i].good++
		case "BAD":
			repos[i].bad++
		case "UNKWN":
			repos[i].unknown++
		case "QUEUED":
			repos[i].queued++
		}
	}

	var main strings.Builder
	main.WriteString(header)
	for i, n := range reponames {
		if repos[i].empty() {
			continue
		}
		main.WriteString(fmt.Sprintf(shortlog, n, n, repos[i].perc(), repos[i].good, repos[i].bad, repos[i].queued, repos[i].unknown))
	}
	for i, n := range reponames {
		if repos[i].empty() {
			continue
		}
		repos[i].b.WriteString(ttail)
		main.WriteString(fmt.Sprintf(subsection, n, n))
		main.WriteString(repos[i].b.String())
	}
	return main.String()
}

type pkgver struct {
	distro string
	suite  string
	name   string
}

func newPkgver(p rebuilderd.PkgRelease) pkgver {
	return pkgver{
		distro: p.Distro,
		suite:  p.Suite,
		name:   p.Name,
	}
}

type Repo struct {
	b       strings.Builder
	total   int
	good    int
	bad     int
	unknown int
	queued  int
}

func (r Repo) empty() bool {
	return r.total == 0
}

func (r Repo) perc() int {
	if r.empty() {
		return 0
	}
	return 100 * r.good / (r.total)
}
