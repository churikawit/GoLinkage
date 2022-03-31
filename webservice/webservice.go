package webservice

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	ami "github.com/churikawit/goami"
	scapi "github.com/churikawit/goscapi"
)

var (
	status_ok string = ""
)

type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type Person struct {
	Name       string    `form:"name"`
	Address    string    `form:"address"`
	Birthday   time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
	CreateTime time.Time `form:"createTime" time_format:"unixNano"`
	UnixTime   time.Time `form:"unixTime" time_format:"unix"`
}

type CardData struct {
	Cid        string
	Pid        string
	FullName   string
	Title      string
	FirstName  string
	MiddleName string
	LastName   string

	FullName_En   string
	Title_En      string
	FirstName_En  string
	MiddleName_En string
	LastName_En   string

	BirthDate string
	Gender    string

	RequestID     string
	BP1NO         string
	IssueLocation string
	IssuePersonID string
	IssueDate     string
	ExpireDate    string
	CardType      string
	Address       string
}

func filterIP(c *gin.Context) bool {
	fmt.Printf("Client IP: %v\n", c.ClientIP())
	return true

	if c.ClientIP() != "127.0.0.1" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return false
	}
	return true
}

func handleReadCard(c *gin.Context) {
	// Check IP Address
	if !filterIP(c) {
		return
	}

	smartcard, err := scapi.ReadCardData()
	if err != nil {
		output := gin.H{
			"status": string(err.Error()),
		}
		c.JSON(http.StatusOK, output)
		return
	}

	cid, err := scapi.GetCID()
	if err != nil {
		cid = ""
	}

	cdata := CardData{
		Cid:        cid,
		Pid:        smartcard.GetPID(),
		FullName:   smartcard.GetFullName(),
		Title:      smartcard.GetTitle(),
		FirstName:  smartcard.GetFirstName(),
		MiddleName: smartcard.GetMiddleName(),
		LastName:   smartcard.GetLastName(),

		FullName_En:   smartcard.GetFullName_En(),
		Title_En:      smartcard.GetTitle_En(),
		FirstName_En:  smartcard.GetFirstName_En(),
		MiddleName_En: smartcard.GetMiddleName_En(),
		LastName_En:   smartcard.GetLastName_En(),
		BirthDate:     smartcard.GetBirthDate(),
		Gender:        smartcard.GetGender(),
		RequestID:     smartcard.GetRequestID(),
		BP1NO:         smartcard.GetBP1NO(),
		IssueLocation: smartcard.GetIssueLocation(),
		IssuePersonID: smartcard.GetIssuePersonID(),
		IssueDate:     smartcard.GetIssueDate(),
		ExpireDate:    smartcard.GetExpireDate(),
		CardType:      smartcard.GetCardType(),
		Address:       smartcard.GetAddress(),
	}

	output := gin.H{
		"status": status_ok,
		"data":   cdata,
	}
	c.JSON(http.StatusOK, output)
}

// handleGetLinkageToken
// input: {"officecode":"00244"}
func handleGetLinkageToken(c *gin.Context) {
	// Check IP Address
	if !filterIP(c) {
		return
	}

	json := struct {
		OfficeCode string `form:"officecode" json:"officecode" xml:"officecode"  binding:"required"`
	}{}

	// bind JSON data to struct
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}

	officecode := json.OfficeCode
	if len(officecode) != 5 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "officecode length is invalid"})
		return
	}

	// ---------------------- Get Card Data ------------------------------
	smartcard, err := scapi.ReadCardData()
	if err != nil {
		output := gin.H{
			"status": string(err.Error()),
		}
		c.JSON(http.StatusOK, output)
		return
	}
	cid, err := scapi.GetCID()
	if err != nil {
		output := gin.H{
			"status": string(err.Error()),
		}
		c.JSON(http.StatusOK, output)
		return
	}

	// ----------------------- Get Random --------------------------------
	pid := smartcard.GetPID()
	random, err := ami.AMI_REQUEST_9080(pid, cid, officecode)
	if err != nil {
		output := gin.H{
			"status": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, output)
		return
	}
	fmt.Printf("random: %v\n", random)

	// ----------------------- Verify Pin --------------------------------
	envelope, err := scapi.GetEnvelopeByVerifyPin(random)
	if err != nil {
		output := gin.H{
			"status": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, output)
		return
	}
	envelope = strings.Trim(envelope, " ")
	envelope = strings.Trim(envelope, "\u0000")
	envelope = fmt.Sprintf("%v:%s", len(envelope), envelope)
	fmt.Printf("envelop: %v\n", envelope)

	// ----------------------- Get Token --------------------------------
	token, err := ami.AMI_REQUEST_9081(pid, cid, random, envelope)
	if err != nil {
		output := gin.H{
			"status": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, output)
		return
	}
	fmt.Printf("token: %v\n", token)

	// -------------------------------------------------------------------

	output := gin.H{
		"status": status_ok,
		"random": random,
		"token":  token,
	}
	c.JSON(http.StatusOK, output)
}

// handleGetLinkageToken
// input: {"token":"xxxxxxxxxx", "pid":"pid for inquiry" }
func handleInquireIdData(c *gin.Context) {
	// Check IP Address
	if !filterIP(c) {
		return
	}

	json := struct {
		Token string `form:"token" json:"token" xml:"token"  binding:"required"`
		Pid   string `form:"pid" json:"pid" xml:"pid"  binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}
	token := json.Token
	pid := json.Pid
	if len(pid) != 13 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "pid length is invalid"})
		return
	}
	// ----------------------- Send 5000 --------------------------------
	officecode := "00023"
	versioncode := "01"
	servicecode := "001"
	data, err := ami.AMI_REQUEST_5000(token, officecode, versioncode, servicecode, pid)
	if err != nil {
		output := gin.H{
			"status": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	data = strings.Trim(data, "\u0000")

	fmt.Printf("data: %v\n", data)
	// -------------------------------------------------------------------
	d, err := ami.BindIdData(data)
	if err != nil {
		output := gin.H{
			"status": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	output := gin.H{
		"status": status_ok,
		"data":   d,
	}
	c.JSON(http.StatusOK, output)
}

// handleInquireHomeData
// input: {"token":"xxxxxxxxxx", "pid":"pid for inquiry" }
func handleInquireHomeData(c *gin.Context) {
	// Check IP Address
	if !filterIP(c) {
		return
	}

	json := struct {
		Token string `form:"token" json:"token" xml:"token"  binding:"required"`
		Pid   string `form:"pid" json:"pid" xml:"pid"  binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
		return
	}
	token := json.Token
	pid := json.Pid
	if len(pid) != 13 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "pid length is invalid"})
		return
	}
	// ----------------------- Send 5000 --------------------------------
	officecode := "00023"
	versioncode := "01"
	servicecode := "027"
	data, err := ami.AMI_REQUEST_5000(token, officecode, versioncode, servicecode, pid)
	if err != nil {
		output := gin.H{
			"status": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	data = strings.Trim(data, "\u0000")
	// -------------------------------------------------------------------
	d, err := ami.BindHomeData(data)
	if err != nil {
		output := gin.H{
			"status": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, output)
		return
	}

	output := gin.H{
		"status": status_ok,
		"data":   d,
	}
	c.JSON(http.StatusOK, output)
}

func Run() {
	// gin.SetMode(gin.ReleaseMode)

	// Default With the Logger and Recovery middleware already attached
	var r *gin.Engine = gin.Default()

	r.GET("/ReadCard", handleReadCard)

	r.POST("/GetLinkageToken", handleGetLinkageToken)

	r.POST("/InquireIdData", handleInquireIdData)

	r.POST("/InquireHomeData", handleInquireHomeData)

	r.Run("0.0.0.0:8080")
}
