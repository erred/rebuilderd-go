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
</tr>
</thead>
</tbody>`

	trow = `
<tr class="%s">
<td><strong>%s</strong></td>
<td>%s</td>
<td>%s</td>
<td>%s</td>
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

func pkgs2page(packages []rebuilderd.PkgRelease, queue []rebuilderd.QueueItem) string {
	reponames := []string{"core", "extra", "community"}
	repos := make(map[string]*Repo, len(reponames))
	for _, n := range reponames {
		repos[n] = &Repo{
			pkgs: make(map[pkg]status),
			s:    make(map[string]int),
		}
	}
	for _, p := range packages {

		repos[p.Suite].pkgs[newPkg(p)] = status{status: p.Status, arch: p.Architecture}
		repos[p.Suite].s[p.Status]++
	}
	for _, qi := range queue {
		p := qi.Package
		p.Status = "QUEUED"
		if stat, ok := repos[p.Suite].pkgs[newPkg(p)]; ok {
			if stat.status == "UNKWN" {
				stat.status = p.Status
				repos[p.Suite].s["UNKWN"]--
			} else {
				stat.queued = true
			}
			repos[p.Suite].pkgs[newPkg(p)] = stat
		} else {
			repos[p.Suite].pkgs[newPkg(p)] = status{status: p.Status, arch: p.Architecture}
		}
		repos[p.Suite].s[p.Status]++
	}
	for _, n := range reponames {
		repo := repos[n]
		if len(repo.pkgs) == 0 {
			continue
		}

		repo.sum = fmt.Sprintf(shortlog, n, n, repo.perc(), repo.s["GOOD"], repo.s["BAD"], repo.s["QUEUED"], repo.s["UNKWN"])

		pkgs := make([]pkg, 0, len(repo.pkgs))
		for k := range repo.pkgs {
			pkgs = append(pkgs, k)
		}
		sort.Slice(pkgs, func(i, j int) bool {
			if pkgs[i].name == pkgs[j].name {
				return pkgs[i].vers < pkgs[j].vers
			}
			return pkgs[i].name < pkgs[j].name
		})
		buf := strings.Builder{}
		buf.WriteString(fmt.Sprintf(subsection, n, n))
		for _, k := range pkgs {
			p := repo.pkgs[k]
			buf.WriteString(fmt.Sprintf(trow, p.status, p.status, k.name, k.vers, p.arch))
			if p.queued {
				buf.WriteString(fmt.Sprintf(trow, "QUEUED", "QUEUED", k.name, k.vers, p.arch))
			}
		}
		buf.WriteString(ttail)
		repo.tab = buf.String()
	}

	var main strings.Builder
	main.WriteString(header)
	for _, n := range reponames {
		if len(repos[n].pkgs) == 0 {
			continue
		}
		main.WriteString(repos[n].sum)
	}
	for _, n := range reponames {
		if len(repos[n].pkgs) == 0 {
			continue
		}
		main.WriteString(repos[n].tab)
	}
	return main.String()
}

type pkg struct {
	name string
	vers string
}

func newPkg(p rebuilderd.PkgRelease) pkg {
	return pkg{
		name: p.Name,
		vers: p.Version,
	}
}

type status struct {
	arch   string
	status string
	queued bool
}

type Repo struct {
	b    strings.Builder
	pkgs map[pkg]status
	s    map[string]int
	sum  string
	tab  string
}

func (r Repo) perc() int {
	m := make(map[string]struct{})
	for k := range r.pkgs {
		m[k.name] = struct{}{}
	}
	return 100 * r.s["GOOD"] / len(m)
}
