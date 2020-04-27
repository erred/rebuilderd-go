package status

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"go.seankhliao.com/rebuilderd-go"
	"go.seankhliao.com/usvc"
	"go.seankhliao.com/webstyle"
)

const (
	thead = `<table>
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
<td>%s</td>
<td>%s</td>
<td>%s</td>
<td>%s</td>
<td><a href="%s">link</a></td>
</tr>
`
	ttail = `
</tbody>
</table>`

	style = `
header{display: none;}
main{margin-top: 0;}
.BAD{color: #BF616A;}
.GOOD{color: #A3BE8C;}
.UNKWN{color: #999;}
`

	header = `
<h3><em>Archlinux</em> Reproducible Builds</h3>
<p>
<a href="https://github.com/kpcyrd/rebuilderd">rebuilderd</a>
run by
<em><a href="https://seankhliao.com">seankhliao</a></em>
</p>

<p><em><a href="#core">Core</a></em>:
%d%% reproducible with
%d <span class="GOOD">good</span> /
%d <span class="BAD">bad</span> /
%d <span class="UNKWN">unknown</span></p>

<p><em><a href="#community">Community</a></em>:
%d%% reproducible with
%d <span class="GOOD">good</span> /
%d <span class="BAD">bad</span> /
%d <span class="UNKWN">unknown</span></p>

<h4 id="core"><em>Core</em></h4>`
)

type Server struct {
	rc   *rebuilderd.Client
	t    *template.Template
	page map[string]string

	Svc *usvc.ServerSecure
}

func NewServer(args []string) *Server {
	var certFile, keyFile, endpoint, gaid string
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.StringVar(&certFile, "cert", "/etc/letsencrypt/live/sne.seankhliao.com/fullchain.pem", "fullchain certificate file")
	fs.StringVar(&keyFile, "key", "/etc/letsencrypt/live/sne.seankhliao.com/privkey.pem", "private key file")
	fs.StringVar(&endpoint, "endpoint", "http://145.100.104.117:8910", "rebuilderd api endpoint")
	fs.StringVar(&gaid, "gaid", "UA-114337586-1", "google analytics id")
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
	s.Svc.Mux.Handle("/", s)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pkgs, err := s.rc.PkgsList(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	page := make(map[string]string, len(s.page))
	for k, v := range s.page {
		page[k] = v
	}

	var core, comm strings.Builder
	var coret, coreg, coreb, coreu, commt, commg, commb, commu int
	core.WriteString(thead)
	comm.WriteString(thead)

	for _, p := range pkgs {
		s := fmt.Sprintf(trow, p.Status, p.Status, p.Name, p.Architecture, p.Version, p.URL)
		switch p.Suite {
		case "core":
			core.WriteString(s)
			coret++
			switch p.Status {
			case "GOOD":
				coreg++
			case "BAD":
				coreb++
			case "UNKWN":
				coreu++
			}
		case "community":
			comm.WriteString(s)
			commt++
			switch p.Status {
			case "GOOD":
				commg++
			case "BAD":
				commb++
			case "UNKWN":
				commu++
			}
		}
	}

	core.WriteString(ttail)
	comm.WriteString(ttail)

	var main strings.Builder
	main.WriteString(fmt.Sprintf(header, (100 * coreg / coret), coreg, coreb, coreu, (100 * commg / commt), commg, commb, commu))
	main.WriteString(core.String())
	main.WriteString(`<h4 id="community"><em>Community</em></h4>`)
	main.WriteString(comm.String())
	page["Main"] = main.String()

	s.t.ExecuteTemplate(w, "LayoutGohtml", page)
}
