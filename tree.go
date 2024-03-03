package main

//. C-code voor omzetten van dot naar svg

/*
#cgo LDFLAGS: -lgvc -lcgraph
#include <graphviz/gvc.h>
#include <graphviz/cgraph.h>
#include <stdlib.h>

char *makeGraph(char *data) {
	Agraph_t *G;
	char *s;
	unsigned int n;
	GVC_t *gvc;

	s = NULL;
	gvc = gvContext();
	G = agmemread(data);
	free(data);
	if (G == NULL) {
		gvFreeContext(gvc);
		return s;
	}
	gvLayout(gvc, G, "dot");
	gvRenderData(gvc, G, "svg", &s, &n);
	gvFreeLayout(gvc, G);
	agclose(G);
	gvFreeContext(gvc);

	return s;
}
*/
import "C"

//. Imports

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/rug-compling/alpinods"
)

//. Structs

type TreeContext struct {
	//	marks    map[string]bool
	refs   map[string]bool
	mnodes map[int]bool
	graph  bytes.Buffer // definitie dot-bestand
	start  int
	words  []string
	// ud1      map[string]bool
	// ud2      map[string]bool
	SkipThis map[int]bool
	fp       io.Writer
}

//. Variables

var treeMu sync.Mutex

//. Functies

func tree(data []byte, fp io.Writer, filename string) {
	ctx := &TreeContext{
		//		marks:    make(map[string]bool), // node met vette rand en edges van en naar de node, inclusief coindex
		refs:   make(map[string]bool),
		mnodes: make(map[int]bool), // gekleurde nodes in boom
		words:  make([]string, 0),
		// ud1:      make(map[string]bool),
		// ud2:      make(map[string]bool),
		SkipThis: make(map[int]bool),
		fp:       fp,
	}

	if ids, ok := IDs[filename]; ok {
		for _, id := range ids {
			ctx.mnodes[id] = true
		}
	}

	/*
		if *optU != "" {
			for _, m := range strings.Split(*optU, ",") {
				ctx.ud1[m] = true
			}
		}

		if *optE != "" {
			for _, m := range strings.Split(*optE, ",") {
				ctx.ud2[m] = true
			}
		}
	*/

	var alpino alpinods.AlpinoDS
	x(xml.Unmarshal(data, &alpino))

	title := html.EscapeString(alpino.Sentence.Sentence)
	ctx.words = strings.Fields(title)

	fmt.Fprintf(fp, `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<script type="text/javascript">
var tooltip=function(){
    var id = 'tt';
    var top = 3;
    var left = 3;
    var maxw = 500;
    var speed = 10;
    var timer = 20;
    var endalpha = 95;
    var alpha = 0;
    var tt,t,c,b,h;
    var ie = document.all ? true : false;
    return{
	show:function(v,w,below){
	    if (v.search("&lt;table") == 0) {
		v = v.replace(/&lt;/g, "<")
		    .replace(/&gt;/g, ">")
		    .replace(/&quot;/g, "\"")
		    .replace(/&apos;/g, "'")
		    .replace(/&amp;/g, "&");
	    }
	    if(tt == null){
		tt = document.createElement('div');
		tt.setAttribute('id',id);
		t = document.createElement('div');
		t.setAttribute('id',id + 'top');
		c = document.createElement('div');
		c.setAttribute('id',id + 'cont');
		b = document.createElement('div');
		b.setAttribute('id',id + 'bot');
		tt.appendChild(t);
		tt.appendChild(c);
		tt.appendChild(b);
		document.body.appendChild(tt);
		tt.style.opacity = 0;
		tt.style.filter = 'alpha(opacity=0)';
	    }
	    document.onmousemove = below ? this.pos2 : this.pos;
	    tt.style.display = 'block';
	    c.innerHTML = v;
	    tt.style.width = w ? w + 'px' : 'auto';
	    if(!w && ie){
		t.style.display = 'none';
		b.style.display = 'none';
		tt.style.width = tt.offsetWidth;
		t.style.display = 'block';
		b.style.display = 'block';
	    }
	    if(tt.offsetWidth > maxw){tt.style.width = maxw + 'px'}
	    h = parseInt(tt.offsetHeight) + top;
	    clearInterval(tt.timer);
	    tt.timer = setInterval(function(){tooltip.fade(1)},timer);
	},
	pos:function(e){
	    var u = ie ? event.clientY + document.documentElement.scrollTop : e.pageY;
	    var l = ie ? event.clientX + document.documentElement.scrollLeft : e.pageX;
	    var w = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth;
	    var o = window.pageXOffset || document.documentElement.scrollLeft || document.body.scrollLeft || 0;
	    var scroll = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop || 0;
	    var top = u - h;
	    if (top < scroll + 10) {
		top = scroll + 10;
	    }
	    tt.style.top = top + 'px';
	    if (w > maxw && l + maxw > w + o) {
		tt.style.right = (w - l - left + 10) + 'px';
		tt.style.left = 'auto';
	    } else {
		tt.style.left = (l + left + 10) + 'px';
		tt.style.right = 'auto';
	    }
	},
	pos2:function(e){
	    var u = ie ? event.clientY + document.documentElement.scrollTop : e.pageY;
	    var l = ie ? event.clientX + document.documentElement.scrollLeft : e.pageX;
	    var w = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth;
	    var o = window.pageXOffset || document.documentElement.scrollLeft || document.body.scrollLeft || 0;
	    var scroll = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop || 0;
	    var top = u + 24;
	    if (top < scroll + 10) {
		top = scroll + 10;
	    }
	    tt.style.top = top + 'px';
	    if (w > maxw && l + maxw > w + o) {
		tt.style.right = (w - l - left + 10) + 'px';
		tt.style.left = 'auto';
	    } else {
		tt.style.left = (l + left + 10) + 'px';
		tt.style.right = 'auto';
	    }
	},
	fade:function(d){
	    var a = alpha;
	    if((a != endalpha && d == 1) || (a != 0 && d == -1)){
		var i = speed;
		if(endalpha - a < speed && d == 1){
		    i = endalpha - a;
		}else if(alpha < speed && d == -1){
		    i = a;
		}
		alpha = a + (i * d);
		tt.style.opacity = alpha * .01;
		tt.style.filter = 'alpha(opacity=' + alpha + ')';
	    }else{
		clearInterval(tt.timer);
		if(d == -1){tt.style.display = 'none'}
	    }
	},
	hide:function(){
	    clearInterval(tt.timer);
	    tt.timer = setInterval(function(){tooltip.fade(-1)},timer);
	}
    };
}();
</script>
<style type="text/css">
// body {font:14px/1.5 Verdana, Arial, Helvetica, sans-serif; }

#tt {position:absolute; display:block; }
#tttop {display:block; height:5px; overflow:hidden}
#ttcont {display:block; padding:.5em 1em; background:#666; color:#FFF; overflow:visible;}
#ttbot {display:block; height:5px; overflow:hidden}

table.attr tr {
    padding: 0px;
    margin: 0px;
}

table.attr tr td {
    margin: 0px;
    padding: 0px;
    vertical-align: top;
}

table.attr tr td.lbl {
    padding-right: 1em;
}
  div.break {
    margin-top: 1em;
    padding-top: 1em;
    border-top: 1px solid grey;
  }
  div.warning {
    margin: 2em 20%%;
    padding: 1em;
    text-align: center;
    border: 4px solid red;
    background-color: #ffe0e0;
  }
  div.unidep {
    overflow-x: auto;
  }
  .udcontrol {
    margin-bottom: 200px;
  }
  .udcontrol input,
  .udcontrol label {
    cursor: pointer;
  }
  .udcontrol label:hover {
    color: #0000e0;
    text-decoration: underline;
  }
</style>
</head>
<body>
<em>%s</em>
<p>
<div style="overflow-x:auto">
`, title, strings.Join(ctx.words, " "))

	if alpino.Metadata != nil && len(alpino.Metadata.Meta) > 0 {
		for _, m := range alpino.Metadata.Meta {
			var v string
			/*
				if m.Type == "date" {
					t, err := time.Parse("2006-01-02", m.Value)
					if err != nil {
						v = err.Error()
					} else {
						v = ranges.PrintDate(t, false)
					}
				} else if m.Type == "datetime" {
					t, err := time.Parse("2006-01-02 15:04", m.Value)
					if err != nil {
						v = err.Error()
					} else {
						v = ranges.PrintDate(t, true)
					}
				} else {
			*/
			v = m.Value
			//}
			fmt.Fprintf(fp, "%s: %s<br>\n", html.EscapeString(m.Name), html.EscapeString(v))
		}
		fmt.Fprintln(fp, "<p>")
	}

	if alpino.Parser != nil && (alpino.Parser.Cats != "" || alpino.Parser.Skips != "") {
		fmt.Fprintf(fp, "cats: %s<br>\nskips: %s\n<p>\n", html.EscapeString(alpino.Parser.Cats), html.EscapeString(alpino.Parser.Skips))
	}

	// BEGIN: definitie van dot-bestand aanmaken.

	ctx.graph.WriteString(`strict graph gr {

    ranksep=".25 equally"
    nodesep=.05
    ordering=out

    node [shape=plaintext, height=0, width=0, fontsize=12, fontname="Helvetica"];

`)

	// Registreer alle markeringen van nodes met een verwijzing.
	// set_refs(fp, alpino.Node)

	// Nodes
	print_nodes(ctx, alpino.Node)

	// Terminals
	ctx.graph.WriteString("\n    node [fontname=\"Helvetica-Oblique\", shape=box, color=\"#d3d3d3\", style=filled];\n\n")
	ctx.start = 0
	terms := print_terms(ctx, alpino.Node)
	sames := strings.Split(strings.Join(terms, " "), "|")
	for _, same := range sames {
		same = strings.TrimSpace(same)
		if same != "" {
			ctx.graph.WriteString("\n    {rank=same; " + same + " }\n")
		}
	}

	// Edges
	ctx.graph.WriteString("\n    edge [sametail=true, color=\"#d3d3d3\"];\n\n")
	print_edges(ctx, alpino.Node)

	ctx.graph.WriteString("}\n")

	// EINDE: definitie van dot-bestand aanmaken.

	// Dot omzetten naar svg.
	// De C-string wordt in de C-code ge'free'd.
	// Zeldzame crash toen er zo te zien twee bomen tegelijk getekend werden.
	// Is graphviz soms niet thread-safe?
	treeMu.Lock()
	s := C.makeGraph(C.CString(ctx.graph.String()))
	svg := C.GoString(s)
	C.free(unsafe.Pointer(s))
	treeMu.Unlock()

	// BEGIN: svg nabewerken en printen

	// XML-declaratie en DOCtype overslaan
	if i := strings.Index(svg, "<svg"); i < 0 {
		x(fmt.Errorf("BUG"))
	} else {
		svg = svg[i:]
	}

	a := ""
	for _, line := range strings.SplitAfter(svg, "\n") {
		// alles wat begint met <title> weghalen
		i := strings.Index(line, "<title")
		if i >= 0 {
			line = line[:i] + "\n"
		}

		// <a xlink> omzetten in tooltip
		i = strings.Index(line, "<a xlink")
		if i >= 0 {
			s := line[i:]
			line = line[:i] + "\n"
			i = strings.Index(s, "\"")
			s = s[i+1:]
			i = strings.LastIndex(s, "\"")
			a = strings.TrimSpace(s[:i])

		}
		if strings.HasPrefix(line, "<text ") && a != "" {
			line = "<text onmouseover=\"tooltip.show('" + html.EscapeString(a) + "')\" onmouseout=\"tooltip.hide()\"" + line[5:]
		}
		if strings.HasPrefix(line, "</a>") {
			line = ""
			a = ""
		}

		fmt.Fprint(fp, line)
	}
	// EIND: svg nabewerken en printen

	fmt.Fprint(fp, "</div>\n<p>\n")

	conllu2svg(fp, 1, &alpino, ctx, data)

	fmt.Fprint(fp, "\n</body>\n</html>\n")
}

//. Genereren van dot

func print_nodes(ctx *TreeContext, node *alpinods.Node) {
	idx := ""
	style := ""

	// Als dit een node met index is, dan in vierkant zetten.
	// Als de node gemarkeerd is, dan in zwart, anders in lichtgrijs.
	// Index als nummer in label zetten.

	if node.Index > 0 {
		idx = fmt.Sprintf("\\n%v", node.Index)
		style += ", color=\"#d3d3d3\""
		if ctx.mnodes[node.ID] {
			style += ", style=filled, fillcolor=\"#ffa07a\""
		} else {
			style += ", shape=box"
		}
	} // else

	if ctx.mnodes[node.ID] {
		style += ", color=\"#ffa07a\", style=filled"
	}

	// attributen
	var tooltip bytes.Buffer
	tooltip.WriteString("<table class=\"attr\">")
	for _, attr := range node.Attrs() {
		tooltip.WriteString(fmt.Sprintf("<tr><td class=\"lbl\">%s:<td>%s", html.EscapeString(attr.Name), html.EscapeString(attr.Value)))
	}
	tooltip.WriteString("</table>")

	lbl := dotquote(node.Rel) + idx

	// als dit geen lege index-node is, dan attributen toevoegen
	if !(node.Index > 0 && (node.Node == nil || len(node.Node) == 0) && node.Word == "") {
		if node.Cat != "" && node.Cat != node.Rel {
			lbl += "\\n" + dotquote(node.Cat)
		} else if node.Pt != "" && node.Pt != node.Rel {
			lbl += "\\n" + dotquote(node.Pt)
		}
	}

	ctx.graph.WriteString(fmt.Sprintf("    n%v [label=\"%v\"%s, tooltip=\"%s\"];\n", node.ID, lbl, style, dotquote2(tooltip.String())))
	for _, d := range node.Node {
		print_nodes(ctx, d)
	}
}

// Geeft een lijst terminals terug die op hetzelfde niveau moeten komen te staan,
// met "|" ingevoegd voor onderbrekingen in niveaus.
func print_terms(ctx *TreeContext, node *alpinods.Node) []string {
	terms := make([]string, 0)

	if node.Node == nil || len(node.Node) == 0 {
		if node.Word != "" {
			// Een terminal
			idx := ""
			col := ""
			if node.Begin != ctx.start {
				// Onderbeking
				terms = append(terms, "|")
				// Onzichtbare node invoegen om te scheiden van node die links staat
				ctx.graph.WriteString(fmt.Sprintf("    e%v [label=\" \", tooltip=\" \", style=invis];\n", node.ID))
				terms = append(terms, fmt.Sprintf("e%v", node.ID))
				ctx.SkipThis[node.ID] = true
			}
			ctx.start = node.End
			terms = append(terms, fmt.Sprintf("t%v", node.ID))
			if node.Lemma == "" {
				ctx.graph.WriteString(fmt.Sprintf("    t%v [label=\"%s%s\", tooltip=\"%s\"%s];\n",
					node.ID, idx, dotquote(node.Word), dotquote2(node.Postag), col))
			} else {
				ctx.graph.WriteString(fmt.Sprintf("    t%v [label=\"%s%s\", tooltip=\"%s:%s\"%s];\n",
					node.ID, idx, dotquote(node.Word), dotquote2(node.Lemma), dotquote(node.Postag), col))
			}
		} else {
			// Een lege node met index
		}
	} else {
		for _, d := range node.Node {
			t := print_terms(ctx, d)
			terms = append(terms, t...)
		}
	}
	return terms
}

func print_edges(ctx *TreeContext, node *alpinods.Node) {
	if node.Node == nil || len(node.Node) == 0 {
		if ctx.SkipThis[node.ID] {
			// Extra: Onzichtbare edge naar extra onzichtbare terminal
			ctx.graph.WriteString(fmt.Sprintf("    n%v -- e%v [style=invis];\n", node.ID, node.ID))
		}

		// geen edge voor lege indexen
		if node.Index == 0 || node.Word != "" {
			// Gewone edge naar terminal
			ctx.graph.WriteString(fmt.Sprintf("    n%v -- t%v;\n", node.ID, node.ID))
		}
	} else {
		// Edges naar dochters
		for _, d := range node.Node {
			// Gewone edge naar dochter
			ctx.graph.WriteString(fmt.Sprintf("    n%v -- n%v;\n", node.ID, d.ID))
		}
		for _, d := range node.Node {
			print_edges(ctx, d)
		}
	}
}

func dotquote(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	return s
}

func dotquote2(s string) string {
	s = strings.Replace(s, "\\", "\\\\\\\\", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	return s
}

// Zet lijst van indexen (string met komma's) om in map[int]bool
func indexes(s string) map[int]bool {
	m := make(map[int]bool)
	for _, i := range strings.Split(s, ",") {
		j, err := strconv.Atoi(i)
		if err == nil {
			m[j] = true
		}
	}
	return m
}
