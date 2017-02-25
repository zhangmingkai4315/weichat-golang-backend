package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"log"
	"encoding/json"
	"github.com/zhangmingkai4315/weichat-golang-backend/utils"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	ISC_URL        = "https://kb.isc.org/article/AA-00913/0/BIND-9-Security-Vulnerability-Matrix.html"
	SAVE_FILE_PATH = "./data/lastest.file"
)

var (
	ErrUnequalResults     = errors.New("Unequal results and urls return")
	ErrOutOfSelectedRange = errors.New("Selected range out of 3")
)

type BindValunerabilityDetail struct {
	CVE             string `json:cve`
	DocumentVersion string `json:document_version`
	PostingDate     string `json:posting_date`
	VersionAffected string `json:"version_affected"`
	Severity        string `json:"severity"`
	Expoitable      string `json:"expoitable"`
	Description     string `json:"description"`
	Impact          string `json:"impact"`
}
type BindValunerability struct {
	ID               int                      `json:id`
	CVE              string                   `json:cve`
	ShortDescription string                   `json:short_descriotion`
	Url              string                   `json:url`
	QueryStatus      utils.Status             `json:status`
	Detail           BindValunerabilityDetail `json:detail`
}

func NewBindValunerabiltiy() *BindValunerability {
	return &BindValunerability{
		QueryStatus: utils.QUERY_INIT,
		Detail:      BindValunerabilityDetail{},
	}
}
func NewBindValunerabilityDetail(cve string) *BindValunerabilityDetail {
	return &BindValunerabilityDetail{CVE: cve}
}

func splitAndLast(s string) string {
	temp := strings.Split(s, ":")
	if len(temp) > 1 {
		return strings.TrimSpace(temp[1])
	} else {
		return "None"
	}
}
func (self *BindValunerabilityDetail) GetDetailInfo(cve, url string) error {
	self.CVE = cve
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	self.DocumentVersion = splitAndLast(doc.Find(".txt .field-field-document-version").Last().Text())
	self.PostingDate = splitAndLast(doc.Find(".txt .field-field-date").Last().Text())
	self.Impact = splitAndLast(doc.Find(".txt .field-field-project").Last().Text())
	self.VersionAffected = splitAndLast(doc.Find(".txt .field-field-versions-affected").Last().Text())
	self.Severity = splitAndLast(doc.Find(".txt .field-field-severity").Last().Text())
	self.Expoitable = splitAndLast(doc.Find(".txt .field-field-exploitable").Last().Text())

	lastPos := doc.Find(".txt .field-field-exploitable")
	self.Description = lastPos.Next().Next().Next().Text()
	return nil
}

func (bv *BindValunerability) String() string {
	return fmt.Sprintf("[ID:%d] CVE=%s ShortDescription=%s Url=%s", bv.ID, bv.CVE, bv.ShortDescription, bv.Url)
}

type BindValunerabilityList struct {
	sync.Mutex
	Bvlist      []*BindValunerability
	QueryStatus utils.Status
	Fetched     int32
	Url         string
}

func NewBindValunerabilityList() *BindValunerabilityList {
	return &BindValunerabilityList{
		QueryStatus: utils.QUERY_INIT,
		Fetched:     0,
		Url:         ISC_URL,
		Bvlist:      []*BindValunerability{},
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
							bv.QueryStatus = utils.QUERY_ERROR
						} else {
							bv.Url = link
							bv.QueryStatus = utils.QUERY_READY
						}
					}
				})
				self.Bvlist = append(self.Bvlist, bv)
			})
		}
	})

	self.QueryStatus = utils.QUERY_READY
	return nil
}
func (self *BindValunerabilityList) ShowList() {
	for _, item := range self.Bvlist {
		println(item.String())
	}
}
func (self *BindValunerabilityList) GetLastId() int {
	var lastid int
	for _, item := range self.Bvlist {
		if lastid < item.ID {
			lastid = item.ID
		}
	}
	return lastid
}

func (self *BindValunerabilityList) ShowFetchedNumber() int32 {
	self.Lock()
	defer self.Unlock()
	return self.Fetched
}

//AddFetched 当有一个已经结束的时候，在Fetched上面加1
func (self *BindValunerabilityList) AddFetched() {
	atomic.AddInt32(&self.Fetched, 1)
}

//SaveListToFile 保存查询到的数据到文件列表
func SaveListToFile(file string, bvl *BindValunerabilityList) error {
	fp, err := os.Create(file)
	if err != nil {
		log.Panicf("Unable to create %v. Err: %v.", file, err)
		return err
	}
	defer fp.Close()
	j, err := json.MarshalIndent(bvl, "", " ")
	if err != nil {
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

func ReadLastListFromFile(fileName string) (*BindValunerabilityList, error) {
	bvl := NewBindValunerabilityList()
	rawFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawFile, bvl)
	if err != nil {
		return nil, err
	}
	return bvl, nil
}

func main() {
	bvlist := NewBindValunerabilityList()
	log.Printf("Will start send url request : %s", ISC_URL)
	err := bvlist.Start()
	if err != nil {
		log.Panicf("Request Error : %s", err)
		return
	}
	//与之前下载的json文件进行比对,查看是否需要下载最新
	lastList, err := ReadLastListFromFile(SAVE_FILE_PATH)

	if err == nil {
		lastMaxId := lastList.GetLastId()
		currentMaxId := bvlist.GetLastId()
		// 比对最大记录id值
		if lastMaxId == currentMaxId {
			bvlist.QueryStatus = utils.QUERY_SKIP
			log.Printf("Skip query the detail, the current max id is same like before : %d,", lastMaxId)
			return
		} else if lastMaxId < currentMaxId {
			log.Printf("The current max id is larger than before : %d > %d try to call alarm methord", currentMaxId, lastMaxId)
			for _, b := range bvlist.Bvlist {
				if b.ID == currentMaxId {
					err := b.Detail.GetDetailInfo(b.CVE, b.Url)
					if err != nil {
						b.QueryStatus = utils.QUERY_ERROR
					} else {
						bvlist.AddFetched()
					}
				}
				lastList.Bvlist = append(lastList.Bvlist, b)
				SaveListToFile(SAVE_FILE_PATH, lastList)
			}
			return
		}
	}
	//如果无法获得之前的文件数据，初始化保存数据到文件。
	var wg sync.WaitGroup
	requestChan := make(chan struct{}, utils.MAXCONCURRENT)
	if bvlist.QueryStatus == utils.QUERY_READY {
		log.Println("Will start concurrent query each item")
		wg.Add(len(bvlist.Bvlist))
		for _, b := range bvlist.Bvlist {
			go func(b *BindValunerability) {
				defer wg.Done()
				requestChan <- struct{}{}
				err := b.Detail.GetDetailInfo(b.CVE, b.Url)
				if err != nil {
					b.QueryStatus = utils.QUERY_ERROR
				} else {
					bvlist.AddFetched()
				}
				<-requestChan
			}(b)
		}
		wg.Wait()
		SaveListToFile(SAVE_FILE_PATH, bvlist)
		log.Println("Current Query Number is", bvlist.ShowFetchedNumber())
	}

}
