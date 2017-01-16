package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"sync"
)

const (
	ISC_VULNERABILITY_MATRIX = "https://kb.isc.org/article/AA-00913/0/BIND-9-Security-Vulnerability-Matrix.html"
)

//status means the query running status.
type status int32

//those status will act like a enum type
const (
	QUERY_NOTSTART status = iota
	QUERY_READY
	QUERY_RUNNING
	QUERY_STOPPED
	QUERY_ERROR
)

var (
	ErrUnequalResults     = errors.New("Unequal results and urls return")
	ErrOutOfSelectedRange = errors.New("Selected range out of 3")
)

type BindValunerabilityDetail struct {
	CVE             string  `json:cve`
	DocumentVersion float32 `json:document_version`
	PostingDate     string  `json:posting_date`
	VersionAffected string  `json:"version_affected"`
	Severity        string  `json:"severity"`
	Expoitable      string  `json:"expoitable"`
	Description     string  `json:"description"`
	Impact          string  `json:"impact"`
	ActiveExpoits   string  `json:"active_expoits"`
	Solution        string  `json:"solution"`
}
type BindValunerability struct {
	ID               int                      `json:id`
	CVE              string                   `json:cve`
	ShortDescription string                   `json:short_descriotion`
	Url              string                   `json:url`
	Status           status                   `json:status`
	Detail           BindValunerabilityDetail `json:detail`
}

func (bv *BindValunerability) String() string {
	return fmt.Sprintf("[ID:%d] CVE=%s ShortDescription=%s Url=%s", bv.ID, bv.CVE, bv.ShortDescription, bv.Url)
}

type BindValunerabilityList struct {
	sync.Mutex
	Bvlist  []BindValunerability
	Status  status
	Fetched int
	Url     string
}

func NewBindValunerabilityList() *BindValunerabilityList {
	return &BindValunerabilityList{
		Status:  QUERY_NOTSTART,
		Fetched: 0,
		Url:     ISC_VULNERABILITY_MATRIX,
		Bvlist:  []BindValunerability{},
	}
}

func (self *BindValunerabilityList) Start() error {
	doc, err := goquery.NewDocument(self.Url)
	if err != nil {
		return err
	}
	doc.Find(".txt table").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			//position = 0
			s.Find("tr").Each(func(j int, sub *goquery.Selection) {
				//跳过第一个标题
				if j == 0 {
					return
				}
				bv := NewBindValunerabiltiy()

				sub.Children().Each(func(k int, subsub *goquery.Selection) {
					text := strings.Trim(strings.TrimSpace(subsub.Text()), "\n")
					switch k {
					case 0:
						id, _ := strconv.Atoi(text)
						bv.ID = id
					case 1:
						bv.CVE = text
					case 2:
						bv.ShortDescription = text
						link, exist := subsub.Find("a").Attr("href")
						if exist == false {
							bv.Url = ""
							bv.Status = QUERY_ERROR
						} else {
							bv.Url = link
							bv.Status = QUERY_READY
						}
					}
				})
				self.Bvlist = append(self.Bvlist, *bv)
			})
		}
	})

	for _, item := range self.Bvlist {
		println(item.String())
	}
	//log.Println(result)
	self.Status = QUERY_READY
	return nil
}

func (self *BindValunerabilityList) ShowFetchedNumber() int {
	self.Lock()
	defer self.Unlock()
	return self.Fetched
}
func (self *BindValunerabilityList) AddFetched() {
	self.Lock()
	defer self.Unlock()
	self.Fetched++
}

func NewBindValunerabiltiy() *BindValunerability {
	return &BindValunerability{
		Status: QUERY_NOTSTART,
		Detail: BindValunerabilityDetail{},
	}
}

func (bv *BindValunerability) GetDetail() error {
	return nil
}

func main() {
	bvlist := NewBindValunerabilityList()
	log.Printf("Will start send url request : %s", ISC_VULNERABILITY_MATRIX)
	err := bvlist.Start()
	if err != nil {
		log.Panicf("Request Error : %s", err)
		return
	}
	if bvlist.Status == QUERY_READY {
		log.Println("Will start concurrent query each item")
	}
	log.Println("Current Query Number is", bvlist.ShowFetchedNumber())
}
