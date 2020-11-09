package dns

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/common/utils/retry"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//DNS obj
type DNS struct {
	DNSApiURL      string `json:"dna_api_url"`
	Prefix         string `json:"prefix"`
	Suffix         string `json:"suffix"`
	PublicPrefix   string `json:"public_prefix"`
	InternalPrefix string `json:"internal_prefix"`

	AppCode    string `json:"app_code"`
	ERP        string `json:"erp"`
	TimeStamp  string `json:"time_stamp"`
	BusinessID string `json:"business_id"`
	Sign       string `json:"sign"`

	Ticket string `json:"ticket"`

	NetworkType int    `json:"network_type"`
	Primary     string `json:"primary"`
	AppEnv      string `json:"app_env"`
	ProjectName string `json:"project_name"`
	ManageERP   string `json:"manage_erp"`
	Remark      string `json:"remark"`
	AuthERP     string `json:"aurh_erp"`
}

//DomainReserve ...
type DomainReserve struct {
	Domain      string `json:"domain"`
	Network     int    `json:"network"`
	Primary     string `json:"primary"`
	AppEnv      string `json:"app_env"`
	ProjectName string `json:"project_name"`
	ManageErp   string `json:"manage_erp"`
	Remark      string `json:"remark"`
	Auth        string `json:"auth"`
}

//DomainReserveResult ...
type DomainCommonResult struct {
	AppCode   string `json:"appCode"`
	ResStatus int    `json:"resStatus"`
	ResMsg    string `json:"resMsg"`
}

//DomainBind ...
type DomainBind struct {
	Domain string             `json:"domain"`
	Data   []ViewsRecordsData `json:"data"`
}

type ViewsRecordsData struct {
	Views   []string      `json:"views"`
	Records []RecordsData `json:"records"`
}

type RecordsData struct {
	Type    string   `json:"type"`
	Records []string `json:"records"`
}

type Records struct {
	Records string `json:"records"`
	Weight  int    `json:"weight"`
}

type DomainBindData struct {
	Type    string   `json:"type"`
	Records []string `json:"records"`
}

//DomainCheckResultData ...
type DomainCheckResultData struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

//DomainCheckResult ...
type DomainCheckResult struct {
	DomainCommonResult
	Data DomainCheckResultData `json:"data"`
}

//DomainResult obj
type DomainResult struct {
	Status bool   `json:"status"`
	Domain string `json:"domain"`
}

const (
	AllNetwork             = 0
	InternalUserNetwork    = 1
	InternalManagerNetwork = 2
	PublicNetwork          = 3
)

var networkMap = map[int]int{
	AllNetwork:             0,
	InternalUserNetwork:    1,
	InternalManagerNetwork: 1,
	PublicNetwork:          2,
}

//DNSConf obj
type DNSConf struct {
	DNSApiURL       string `toml:"dns_api_url"`
	DNSDomainSuffix string `toml:"dns_domain_suffix"`
	DNSDomainPrefix string `toml:"dns_domain_prefix"`
	PublicPrefix    string `toml:"public_prefix"`
	InternalPrefix  string `toml:"internal_prefix"`
	AppCode         string `toml:"app_code"`
	ERP             string `toml:"erp"`
	BusinessID      string `toml:"business_id"`

	Primary     string `toml:"primary"`
	AppEnv      string `toml:"app_env"`
	ProjectName string `toml:"project_name"`
	ManageERP   string `toml:"manage_erp"`
	Remark      string `toml:"remark"`
	AuthERP     string `toml:"auth_erp"`
}

//NewCreateDNS api
func NewDNSClient(dnsConf *DNSConf, networkType int) *DNS {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	reqTime := time.Now().In(loc)
	timeStamp := strconv.FormatInt(reqTime.UTC().Unix(), 10)
	StartTime := fmt.Sprintf("%02d%02d%d%02d%02d", reqTime.Hour(), reqTime.Minute(), reqTime.Year(), reqTime.Month(), reqTime.Day())
	sumStr := dnsConf.ERP + "#" + dnsConf.BusinessID + "NP" + StartTime
	sumData := []byte(sumStr)
	signMD5 := md5.Sum(sumData)
	sign := fmt.Sprintf("%x", signMD5)
	dns := &DNS{}
	dns.DNSApiURL = dnsConf.DNSApiURL
	dns.Prefix = dnsConf.DNSDomainPrefix
	dns.Suffix = dnsConf.DNSDomainSuffix
	dns.PublicPrefix = dnsConf.PublicPrefix
	dns.InternalPrefix = dnsConf.InternalPrefix
	dns.AppCode = dnsConf.AppCode
	dns.ERP = dnsConf.ERP
	dns.BusinessID = dnsConf.BusinessID
	dns.TimeStamp = timeStamp
	dns.Sign = sign
	dns.Ticket = ""

	dns.NetworkType = networkType
	dns.Primary = dnsConf.Primary
	dns.AppEnv = dnsConf.AppEnv
	dns.ProjectName = dnsConf.ProjectName
	dns.ManageERP = dnsConf.ManageERP
	dns.Remark = dnsConf.Remark
	dns.AuthERP = dnsConf.AuthERP
	return dns
}

//CheckDomain is
func (dns *DNS) CheckDomain(requestID, domain string) error {

	var req *http.Request
	dnsCheckURL := "http://" + dns.DNSApiURL + "/V1/Dns/domainCheck"
	req, err := http.NewRequest("GET", dnsCheckURL, nil)
	req.URL.RawQuery = "domain=" + domain

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("appCode", dns.AppCode)
	req.Header.Add("erp", dns.ERP)
	req.Header.Add("timestamp", dns.TimeStamp)
	req.Header.Add("sign", dns.Sign)

	golog.Infof(requestID, "dns url %s check... params %+v , reqHeader %+v", dnsCheckURL, domain, req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		golog.Errorf(requestID, "Error send request,error:%s", err.Error())
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		golog.Errorf(requestID, "Error receive response,error:%s", err.Error())
		return err
	}

	var result DomainCheckResult
	golog.Infof(requestID, "check dns ..................result body:%s", string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		golog.Errorf("0", "Error json.Unmarshal,error:%s", err.Error())
		return err
	}

	//if result.ResStatus != 200 || result.Data.Status != -1 {
	//	err = errors.New(result.ResMsg + "," + result.Data.Msg)
	//	return err
	//}

	//response error
	if result.ResStatus != 200 {
		return errors.New(result.ResMsg + "," + result.Data.Msg)
	}
	//domain Unavailable
	if result.Data.Status != -1 {
		return errors.New(result.ResMsg + "," + result.Data.Msg + " domain is exist")
	}

	golog.Infof(requestID, "check dns ..................result:%v", result)

	return nil
}

//PingDomain is
func (dns *DNS) PingDomain(requestID, domain string, ip string) error {
	golog.Infof(requestID, "ping domain %s", domain)
	cmd := exec.Command("ping", "-c", "3", domain)
	out, err := cmd.Output()
	if err != nil {
		errMsg := fmt.Sprintf("domain %s not avaiable, error: %", domain, err)
		err = errors.New(errMsg)
		golog.Warnf(requestID, err.Error())
		return err
	}

	if ip != "" && !strings.Contains(string(out), ip) {
		errMsg := fmt.Sprintf("domain %s is avaiable, but not ip %s", domain, ip)
		err = errors.New(errMsg)
		golog.Warn(requestID, err.Error())
		return err
	}
	return nil
}

//Binding api
func (dns *DNS) DomainReserve(requestID string, domain string) error {

	newDomain := &DomainReserve{Domain: domain, Network: dns.NetworkType, Primary: dns.Primary, ProjectName: dns.ProjectName, ManageErp: dns.ManageERP, AppEnv: dns.AppEnv}
	bs, err := json.Marshal(newDomain)
	if err != nil {
		golog.Errorf(requestID, "Error json.Marshal error:%s", err.Error())
		return err
	}
	golog.Infof(requestID, "DNSReserve dns para :%+v", newDomain)

	var req *http.Request
	dnsReserveURL := "http://" + dns.DNSApiURL + "/V1/Dns/reserve"
	req, err = http.NewRequest("POST", dnsReserveURL, bytes.NewReader(bs))

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("appCode", dns.AppCode)
	req.Header.Add("erp", dns.ERP)
	req.Header.Add("timestamp", dns.TimeStamp)
	req.Header.Add("sign", dns.Sign)

	if err != nil {
		golog.Errorf(requestID, "Error newrequest, error:%s", err.Error())
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		golog.Errorf(requestID, "Error send request,error:%s", err.Error())
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		golog.Errorf(requestID, "Error receive response,error:%s", err.Error())
		return err
	}

	var result DomainCommonResult
	golog.Infof(requestID, "reserve dns ..................result body:%s", string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		golog.Errorf(requestID, "Error json.Unmarshal error:%s", err.Error())
		return err
	}
	if result.ResStatus != 200 {
		return errors.New(result.ResMsg)
	}
	golog.Infof(requestID, "DNSReserve dns ..................result:%v", result)
	return nil

}

func (dns *DNS) DomainBind(requestID string, domain string, v4Ip, v6Ip string) error {
	//data := []DomainBindData{{Type: "A", Records: ipv4}}
	data := make([]RecordsData, 0)
	if v4Ip != "" {
		data = append(data, RecordsData{Type: "A", Records: []string{v4Ip}})
	}
	if v6Ip != "" {
		data = append(data, RecordsData{Type: "AAAA", Records: []string{v6Ip}})
	}
	domainBind := &DomainBind{Domain: domain, Data: []ViewsRecordsData{{Records: data}}}

	golog.Infof(requestID, "dns domain params domain %+v ,data %+v", domain, data)
	bs, err := json.Marshal(domainBind)
	if err != nil {
		golog.Errorf(requestID, "Error json.Marshal error:%s", err.Error())
		return err
	}

	var req *http.Request

	dnsBindURL := "http://" + dns.DNSApiURL + "/V2/Dns/resolution"
	req, err = http.NewRequest("POST", dnsBindURL, bytes.NewReader(bs))
	if err != nil {
		golog.Errorf(requestID, "Error new request, error:%s", err.Error())
		return err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("appCode", dns.AppCode)
	req.Header.Add("erp", dns.ERP)
	req.Header.Add("timestamp", dns.TimeStamp)
	req.Header.Add("sign", dns.Sign)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		golog.Errorf(requestID, "Error send request,error:%s", err.Error())
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		golog.Errorf(requestID, "Error receive response,error:%s", err.Error())
		return err
	}

	var result DomainCommonResult
	golog.Infof(requestID, "Bind dns ..................result body:%s", string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		golog.Errorf(requestID, "Error json.Unmarshal error:%s", err.Error())
		return err
	}
	if result.ResStatus != 200 {
		return errors.New(result.ResMsg)
	}
	golog.Infof(requestID, "Bind dns ..................result:%v", result)
	return nil
}

//DomainDelete ...
func (dns *DNS) DomainDelete(requestID string, domain string) error {

	var req *http.Request

	dnsDeleteURL := "http://" + dns.DNSApiURL + "/V1/Dns/delDomain"
	req, err := http.NewRequest("GET", dnsDeleteURL, nil)
	if err != nil {
		golog.Errorf(requestID, "Error newrequest, error:%s", err.Error())
		return err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("appCode", dns.AppCode)
	req.Header.Add("erp", dns.ERP)
	req.Header.Add("timestamp", dns.TimeStamp)
	req.Header.Add("sign", dns.Sign)

	req.URL.RawQuery = "domain=" + domain

	golog.Infof(requestID, "dns url %s deleDomain... params %+v", domain, *req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		golog.Errorf(requestID, "Error send request,error:%s", err.Error())
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		golog.Errorf(requestID, "Error receive response,error:%s", err.Error())
		return err
	}

	var result DomainCommonResult
	golog.Infof(requestID, "DelDomain dns ..................result body:%s", string(body))
	err = json.Unmarshal(body, &result)
	if err != nil {
		golog.Errorf(requestID, "Error json.Unmarshal error:%s", err.Error())
		return err
	}

	if result.ResStatus != 200 {
		return errors.New(result.ResMsg)
	}

	golog.Infof(requestID, "DelDomain dns ..................result:%v", result)
	return nil

}

func (dns *DNS) DomainGenerateAndReserve(requestID string, serviceCode string, retryCountLow int, retryIntervalLow time.Duration) (string, error) {
	//Configuration:
	//dns.InternalPrefix=sqlserver-dms-north1-
	//dns.PublicPrefix=sqlserver-internet-north1-
	//dns.Prefix=sqlserver-north1-
	//dns.Suffix=.rds.jdcloud.com

	//serviceCode := pg, mysql, mariadb, procona, and so on

	//如果产品线信息已经固化到 xxxPrefix配置里面，serviceCode传 "" 即可；
	//如果产品线信息是动态生成的，xxxPrefix配置里面不要设置产品信息；

	var domain string
	uuid := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	switch dns.NetworkType {
	case InternalManagerNetwork:
		domain = serviceCode + dns.InternalPrefix + uuid[0:16] + dns.Suffix //example: sqlserver-dms-north1-ae0d9b210de94d.rds.jdcloud.com
		break
	case PublicNetwork:
		domain = serviceCode + dns.PublicPrefix + uuid[0:16] + dns.Suffix //example: sqlserver-internet-north1-ae0d9b210de94d.rds.jdcloud.com
		break
	default:
		domain = serviceCode + dns.Prefix + uuid[0:16] + dns.Suffix //example: sqlserver-north1-ae0d9b210de94d.rds.jdcloud.com
	}
	op := func() error {
		return dns.CheckDomain(requestID, domain)
	}
	err := retry.Do(op, retry.Timeout(0), retry.MaxTries(retryCountLow), retry.Sleep(retryIntervalLow))
	if err != nil {
		golog.Info(requestID, "failed to check domain %s；　%s", domain, err.Error())
		return "", err
	}
	dns.NetworkType = networkMap[dns.NetworkType]
	op = func() error {
		return dns.DomainReserve(requestID, domain)
	}
	err = retry.Do(op, retry.Timeout(0), retry.MaxTries(retryCountLow), retry.Sleep(retryIntervalLow))
	if err != nil {
		golog.Info(requestID, "failed to reserve domain %s；　%s", domain, err.Error())
		return "", err
	}
	return domain, nil
}
