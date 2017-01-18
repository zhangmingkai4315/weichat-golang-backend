package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"log"
	"encoding/json"
	"os"
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
const MAXCONCURRENT = 5
var (
	ErrUnequalResults     = errors.New("Unequal results and urls return")
	ErrOutOfSelectedRange = errors.New("Selected range out of 3")
)

type BindValunerabilityDetail struct {
	CVE             string  `json:cve`
	DocumentVersion string `json:document_version`
	PostingDate     string  `json:posting_date`
	VersionAffected string  `json:"version_affected"`
	Severity        string  `json:"severity"`
	Expoitable      string  `json:"expoitable"`
	Description     string  `json:"description"`
	Impact          string  `json:"impact"`
}
type BindValunerability struct {
	ID               int                      `json:id`
	CVE              string                   `json:cve`
	ShortDescription string                   `json:short_descriotion`
	Url              string                   `json:url`
	Status           status                   `json:status`
	Detail           BindValunerabilityDetail `json:detail`
}
func NewBindValunerabiltiy() *BindValunerability {
	return &BindValunerability{
		Status: QUERY_NOTSTART,
		Detail: BindValunerabilityDetail{},
	}
}
func NewBindValunerabilityDetail(cve string)*BindValunerabilityDetail{
	return &BindValunerabilityDetail{CVE:cve}
}

func splitAndLast(s string)string{
	temp:=strings.Split(s,":")
	if len(temp)>1{
		return strings.TrimSpace(temp[1])
	}else{
		return "None"
	}
}
func(self *BindValunerabilityDetail)GetDetailInfo(cve,url string)error{
	self.CVE=cve
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	self.DocumentVersion=splitAndLast(doc.Find(".txt .field-field-document-version").Last().Text())
	self.PostingDate=splitAndLast(doc.Find(".txt .field-field-date").Last().Text())
	self.Impact=splitAndLast(doc.Find(".txt .field-field-project").Last().Text())
	self.VersionAffected=splitAndLast(doc.Find(".txt .field-field-versions-affected").Last().Text())
	self.Severity=splitAndLast(doc.Find(".txt .field-field-severity").Last().Text())
	self.Expoitable=splitAndLast(doc.Find(".txt .field-field-exploitable").Last().Text())

	lastPos:=doc.Find(".txt .field-field-exploitable")
	self.Description=lastPos.Next().Next().Next().Text()
	return nil
}



func (bv *BindValunerability) String() string {
	return fmt.Sprintf("[ID:%d] CVE=%s ShortDescription=%s Url=%s", bv.ID, bv.CVE, bv.ShortDescription, bv.Url)
}

type BindValunerabilityList struct {
	sync.Mutex
	Bvlist  []*BindValunerability
	Status  status
	Fetched int32
	Url     string
}

func NewBindValunerabilityList() *BindValunerabilityList {
	return &BindValunerabilityList{
		Status:  QUERY_NOTSTART,
		Fetched: 0,
		Url:     ISC_VULNERABILITY_MATRIX,
		Bvlist:  []*BindValunerability{},
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
				self.Bvlist = append(self.Bvlist, bv)
			})
		}
	})

	self.Status = QUERY_READY
	return nil
}
func (self *BindValunerabilityList) ShowList() {
	for _, item := range self.Bvlist {
		println(item.String())
	}
}
func (self *BindValunerabilityList) ShowFetchedNumber() int32 {
	self.Lock()
	defer self.Unlock()
	return self.Fetched
}
//AddFetched 当有一个已经结束的时候，在Fetched上面加1
func (self *BindValunerabilityList) AddFetched() {
	atomic.AddInt32(&self.Fetched,1)
}

func SaveListToFile(file string,bvl *BindValunerabilityList)error{
	fp, err := os.Create(file)
	if err != nil {
		log.Panicf("Unable to create %v. Err: %v.", file, err)
		return err
	}
	defer fp.Close()
	j,err:=json.MarshalIndent(bvl,""," ")
	if err!=nil{
		return err
	}
	_, werr := fp.Write(j)
	if werr != nil {
		return werr
	}
	return nil
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
	var wg sync.WaitGroup
	requestChan:=make(chan struct{},MAXCONCURRENT)
	if bvlist.Status == QUERY_READY {
		log.Println("Will start concurrent query each item")
		wg.Add(len(bvlist.Bvlist))
		for _,b:=range(bvlist.Bvlist){
			go func(b *BindValunerability){
				defer wg.Done()
				requestChan<-struct{}{}
				err:=b.Detail.GetDetailInfo(b.CVE,b.Url)
				if err!=nil{
					b.Status=QUERY_ERROR
				}else{
					bvlist.AddFetched()
				}
				<-requestChan
			}(b)
		}
		wg.Wait()
		SaveListToFile("lastest.file",bvlist)
	}
	log.Println("Current Query Number is", bvlist.ShowFetchedNumber())




}
