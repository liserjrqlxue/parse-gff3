package parseGff3

import (
	"bufio"
	"compress/gzip"
	"github.com/liserjrqlxue/simple-util"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// regexp
var (
	isGz      = regexp.MustCompile(`\.gz$`)
	isComment = regexp.MustCompile(`^#`)
)

type GFF3 struct {
	Seqid      string
	Source     string
	Type       string
	Start      uint64
	End        uint64
	Score      string
	Strand     string `+:"positive strand (relative to the landmark)",-:"minus strand",.:"not stranded",?:"unknown"`
	Phase      string
	Attributes map[string]string
}

func File2GFF3array(fileName string) (gff3Array []GFF3) {
	file, err := os.Open(fileName)
	simple_util.CheckErr(err)
	defer simple_util.DeferClose(file)

	var scanner *bufio.Scanner
	if isGz.MatchString(fileName) {
		gr, err := gzip.NewReader(file)
		simple_util.CheckErr(err)
		defer simple_util.DeferClose(gr)

		scanner = bufio.NewScanner(gr)
	} else {
		scanner = bufio.NewScanner(file)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if isComment.MatchString(line) {
			continue
		}
		array := strings.Split(line, "\t")
		if len(array) != 9 {
			log.Fatalf("GFF3 line not have 9 column:\n[%s]\n", line)
		}
		var item = new(GFF3)
		item.Seqid = array[0]
		item.Source = array[1]
		item.Type = array[2]
		start, err := strconv.Atoi(array[3])
		simple_util.CheckErr(err)
		item.Start = uint64(start)
		end, err := strconv.Atoi(array[4])
		simple_util.CheckErr(err)
		item.End = uint64(end)
		item.Score = array[5]
		item.Strand = array[6]
		item.Phase = array[7]
		attributes := strings.Split(array[8], ";")
		var attributeMap = make(map[string]string)
		for _, kv := range attributes {
			kvs := strings.SplitN(kv, "=", 2)
			if len(kvs) != 2 {
				log.Fatalf(
					"GFF3 item's attributes no have tag=value format\n\t[%s]\n\t[%s]\n",
					item.Attributes, kv,
				)
			}
			attributeMap[kvs[0]] = kvs[1]
		}
		item.Attributes = attributeMap
		gff3Array = append(gff3Array, *item)

	}
	return
}
