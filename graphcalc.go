package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"time"

	navi "github.com/camen6ert/osmNavigator/navigator"

	"github.com/camen6ert/goXmlParser"
)

//Node is bla
//<node id="6145747792" lat="49.5245371" lon="10.9062853" version="1" timestamp="2018-12-18T17:33:24Z" changeset="0"/>
type Node struct {
	lat float64
	lon float64
}

//Way is bla
type Way struct {
	id    int64
	nodes []int64
	tags  map[string]string
}

//Tag is bla
type Tag struct {
	k string
	v string
}

type graphnode struct {
	e []edge
}

type edge struct {
	w1 int
	w2 int
}

func main() {

	f := flag.String("file", "", "path of osm file")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	file, err := os.Open(*f)
	if err != nil {
		fmt.Printf("File not found %s", err)
		return
	}

	defer file.Close()

	stat, _ := file.Stat()
	fmt.Printf("File size is %d\n", stat.Size())

	start := time.Now()
	r := io.Reader(file)

	nodes := make(map[int64]Node)
	//ways := make(map[int64]Way)
	edges := make(map[int64][]Way)
	way := Way{}
	way.tags = make(map[string]string)

	parentIsNode := true

	goXmlParser.Parse(r,
		func(t goXmlParser.Tag) {
			if t.Name == "way" {
				parentIsNode = false
				id, _ := strconv.ParseInt(t.Attributes["id"], 10, 64)
				way = Way{id: id}

				if id%100 == 0 {
					fmt.Println(id)
				}
			} else if t.Name == "tag" {
				if parentIsNode {
					return
				}

				if way.tags == nil {
					way.tags = make(map[string]string)
				}
				way.tags[t.Attributes["k"]] = t.Attributes["v"]
			} else if t.Name == "nd" {
				id, _ := strconv.ParseInt(t.Attributes["ref"], 10, 64)
				way.nodes = append(way.nodes, id)
				edges[id] = append(edges[id], way)
			} else if t.Name == "node" {
				parentIsNode = true
				id, _ := strconv.ParseInt(t.Attributes["id"], 10, 64)
				lat, _ := strconv.ParseFloat(t.Attributes["lat"], 32)
				lng, _ := strconv.ParseFloat(t.Attributes["lat"], 32)
				n := Node{lat: lat, lon: lng}
				nodes[id] = n
			}
		},
		func(t goXmlParser.Tag) {

			if t.Name == "way" {
				way = Way{}
				way.tags = make(map[string]string)
			}

		})

	processWays(&edges, &nodes)

	elapsed := time.Since(start)
	fmt.Printf("Binomial took %s", elapsed)

}

func processWays(e *map[int64][]Way, nodesPtr *map[int64]Node) {

	fmt.Println("Process ways")

	n := navi.NewNavigator()
	nodes := *nodesPtr

	for buf, edge := range *e {

		fmt.Println(buf)

		for _, w := range edge {

			street := navi.Street{}
			street.ID = w.id
			if s, ok := w.tags["addr:street"]; ok {
				street.Name = s
			}
			if co, ok := w.tags["addr:city"]; ok {
				street.City = co
			}
			if c, ok := w.tags["addr:country"]; ok {
				street.Country = c
			}

			for _, p := range w.nodes {
				pos := navi.Pos{}
				pos.Lat = nodes[p].lat
				pos.Lon = nodes[p].lon
				street.Pos = append(street.Pos, &pos)
			}

			for _, c := range edge {
				if c.id != w.id {
					street.Con = append(street.Con, c.id)
				}
			}

			n.AddStreet(street)

			// if p, ok := w.tags["addr:postcode"]; ok {
			// }
			// if h, ok := way.tags["addr:housenumber"]; ok {
			// }

		}
	}
}
