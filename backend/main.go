package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"quanto.mi.costi.com/utils"
)

type CostItem struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Value    float64 `json:"value"`
	Estimate bool    `json:"estimate"`
	Info     string  `json:"info"`
}

type CashVsMortgageItem struct {
	Percentage  int       `json:"percentage"`
	Values      []float64 `json:"values"`
	Installment float64   `json:"installment"`
}

type MortgageCompareItem struct {
	Duration    int       `json:"duration"`
	Values      []float64 `json:"values"`
	Valid       bool      `json:"valid"`
	Installment float64   `json:"installment"`
}

type Dataset struct {
	Label string      `json:"label"`
	Data  interface{} `json:"data"` // Can be []float64 or []string
}

type ChartData struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

type PriceVolumeResponse struct {
	MarketSize      string    `json:"market_size"`       // Changed to string
	CurrentVolumeMq string    `json:"current_volume_mq"` // Changed to string
	MqRange         string    `json:"mq_range"`
	CurrentMaxPrice float64   `json:"current_max_price"`
	CurrentMinPrice float64   `json:"current_min_price"`
	VolumeTrend     ChartData `json:"volume_trend"`
	PriceTrend      ChartData `json:"price_trend"`
}

type MunicipalitiesInfoResponse struct {
	Data [][]interface{} `json:"DATA"`
}

type MunicipalitiesListResponse struct {
	Data []string `json:"DATA"`
}

const pythonAPIBaseURL = "https://quanto-mi-costi-934184719806.europe-west8.run.app"

func callPythonAPI(endpoint string, params map[string]string) ([]byte, error) {
	// Build URL with query parameters
	baseURL := pythonAPIBaseURL + endpoint
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// Make HTTP request
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the error response body for debugging
		errorBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Python API error (status %d): %s\n", resp.StatusCode, string(errorBody))
		return nil, fmt.Errorf("python API returned status %d: %s", resp.StatusCode, string(errorBody))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func calculateAllCosts(housePrice float64,
	agencyFee float64,
	isFirstHouse bool,
	isSoldByBuilder bool,
	mortgageAmount float64,
	mortgageDuration int,
	mortgageTAEG float64) []CostItem {

	var istruttoriaCost float64 = 0
	var expertiseCost float64 = 0
	var mortgageInstalment float64 = 0

	if mortgageAmount > 0 {
		istruttoriaCost = mortgageAmount * 0.007
		mortgageInstalment = utils.CalculateMonthlyInstallment(mortgageAmount, mortgageDuration, mortgageTAEG)

		if mortgageAmount > 0 {
			expertiseCost = 300*1.03 + 50
		}
	}

	var tassaArchivio float64 = 27.1
	var tassaArchivioMortgage float64 = 27.1
	var visureIpotecarie float64 = 0
	var renditaCatastale float64 = 0

	// Derived calculations (automatically update when dependencies change)
	//estimate to calculate the tassa_archivio cost, min 27.1 max 139.4
	tassaArchivio = min(housePrice*0.0001614, 139.4)
	tassaArchivioMortgage = min(mortgageAmount*0.00017, 139.4)
	//estimate to calculate the visure_ipotecarie cost, can vary between notary
	visureIpotecarie = housePrice * 0.000489 * 1.22
	renditaCatastale = housePrice / 18.89 / 20.5

	var valoreCatastale float64 = 0
	//if sold by builder there isn't valoreCatastale
	if !isSoldByBuilder {
		if isFirstHouse {
			valoreCatastale = renditaCatastale * 1.05 * 110
		} else {
			valoreCatastale = renditaCatastale * 1.05 * 120
		}
	}

	var impostaDiRegistro float64 = 0

	if isSoldByBuilder {
		impostaDiRegistro = 200
	} else {
		if isFirstHouse {
			impostaDiRegistro = max((2*valoreCatastale)/100, 1000)
		} else {
			impostaDiRegistro = max((9*valoreCatastale)/100, 1000)
		}
	}

	var impostaIpotecaria float64 = 50
	var impostaCatastale float64 = 50

	if isSoldByBuilder {
		impostaIpotecaria = 200
		impostaCatastale = 200
	}

	var vat float64 = 0

	if isSoldByBuilder {
		if isFirstHouse {
			vat = 0.04 * housePrice
		} else {
			vat = 0.1 * housePrice
		}
	}

	//multiplied by a coefficent to estimete the cost voice
	var notaryAmount float64 = housePrice * 0.0042 * 1.22
	var notaryAmountMortage float64 = mortgageAmount * 0.0042 * 1.22
	var impostaDiBolloMotgage float64 = mortgageAmount * 0.0025

	if !isFirstHouse {
		impostaDiBolloMotgage = mortgageAmount * 0.02
	}

	var agencyAmount float64 = 0
	var agencyVat float64 = 0
	if agencyFee > 0 {
		agencyAmount = agencyFee * housePrice
		agencyVat = agencyAmount * 0.22
	}

	var registrazionePreliminare float64 = 200 + 16
	var tassaIpotecaria float64 = 90
	var tassaIpotecariaMortgage float64 = 35
	var impostaDiBollo float64 = 230

	var certificatoStatoLibero float64 = 16

	responseData := []CostItem{
		{
			Name:     "Costo casa",
			Category: "House price",
			Value:    housePrice,
			Estimate: false,
			Info:     "House buy price",
		},

		{
			Name:     "Rata mutuo mensile",
			Category: "Installment",
			Value:    mortgageInstalment,
			Estimate: false,
			Info:     "Mortgage intrests",
		},
		{
			Name:     "Provvigione",
			Category: "Agency",
			Value:    agencyAmount,
			Estimate: false,
			Info:     "Agency fee",
		},
		{
			Name:     "IVA",
			Category: "Agency",
			Value:    agencyVat,
			Estimate: false,
			Info:     "Agency fee",
		},
		{
			Name:     "IVA",
			Category: "Tax",
			Value:    vat,
			Estimate: false,
			Info:     "IVA dovuta al venditore, incide solamente se il venditore è passibile di IVA.",
		},
		{
			Name:     "Imposta di registro",
			Category: "Tax",
			Value:    impostaDiRegistro,
			Estimate: false,
			Info:     "Imposta di registro, tassa che varia in base al tipo di venditore.",
		},
		{
			Name:     "Imposta ipotecaria",
			Category: "Tax",
			Value:    impostaIpotecaria,
			Estimate: false,
			Info:     "Imposta ipotecaria, tassa fissa che varia in base al tipo di venditore.",
		},
		{
			Name:     "Imposta catastale",
			Category: "Tax",
			Value:    impostaCatastale,
			Estimate: false,
			Info:     "Imposta catastale, tassa fissa che varia in base al tipo di venditore.",
		},
		{
			Name:     "Registrazione preliminare",
			Category: "Tax",
			Value:    registrazionePreliminare,
			Estimate: false,
			Info:     "Registrazione preliminare, spesa da sostenere durante il preliminare",
		},
		{
			Name:     "Compenso notaio",
			Category: "Notary",
			Value:    notaryAmount,
			Estimate: true,
			Info:     "Compenso onorario notaio",
		},
		{
			Name:     "Tassa di archivio",
			Category: "Notary",
			Value:    tassaArchivio,
			Estimate: true,
			Info:     "Tassa di archivio, varia in relazione all'importo dell'immobile.",
		},
		{
			Name:     "Tassa ipotecaria",
			Category: "Notary",
			Value:    tassaIpotecaria,
			Estimate: false,
			Info:     "Tassa ipotecaria, tassa fissa durante il rogito.",
		},
		{
			Name:     "Visure ipotecarie",
			Category: "Notary",
			Value:    visureIpotecarie,
			Estimate: true,
			Info:     "Visure ipotecarie, spese variabile da sostenere per la consultazione delle visure ipotecarie da parte del notaio durante il rogito.",
		},
		{
			Name:     "Imposta di bollo",
			Category: "Notary",
			Value:    impostaDiBollo,
			Estimate: false,
			Info:     "Imposta di bollo, tassa fissa durante il rogito.",
		},
		{
			Name:     "Certificato stato libero",
			Category: "Notary",
			Value:    certificatoStatoLibero,
			Estimate: false,
			Info:     "Certificato stato libero, necessario per dare valenza legale all'atto di rogito.",
		},
		{
			Name:     "Compenso notaio (mutuo)",
			Category: "Notary",
			Value:    notaryAmountMortage,
			Estimate: true,
			Info:     "Compenso onorario notaio per l'atto riguardante il mutuo.",
		},
		{
			Name:     "Tassa ipotecaria (mutuo)",
			Category: "Notary",
			Value:    tassaIpotecariaMortgage,
			Estimate: false,
			Info:     "Tassa ipotecaria, tassa fissa durante il rogito per l'atto riguardante il mutuo.",
		},
		{
			Name:     "Tassa di archivio (mutuo)",
			Category: "Notary",
			Value:    tassaArchivioMortgage,
			Estimate: true,
			Info:     "Tassa di archivio, varia in relazione all'importo del mutuo.",
		},
		{
			Name:     "Imposta di bollo",
			Category: "Bank",
			Value:    impostaDiBolloMotgage,
			Estimate: false,
			Info:     "Imposta di bollo, tassa da pagare in relazione all'importo del mutuo richiesto.",
		},
		{
			Name:     "Spese di istruttoria",
			Category: "Bank",
			Value:    istruttoriaCost,
			Estimate: true,
			Info:     "Spese di istruttoria, da pagare alla banca per sostenere i costi di apertura del mutuo.",
		},
		{
			Name:     "Spese di perizia",
			Category: "Bank",
			Value:    expertiseCost,
			Estimate: true,
			Info:     "Spese di perizia, da pagare alla banca per sostenere i costi di perizia dell'immobile.",
		},
	}
	return responseData
}

func calculateCashVsMortgage(
	yearlySaving float64,
	mortgageDuration int,
	taeg float64,
	yearlyGrowthgRate float64,
	yearlySavingRate float64,
	morgagePercentages []float64,
	housePrice float64) []CashVsMortgageItem {

	installment := 0.0

	yearlyIncome := make([]float64, mortgageDuration)
	incomeValue := (yearlySaving * yearlySavingRate)
	for i := range yearlyIncome {
		yearlyIncome[i] = incomeValue
	}

	myData := make([]CashVsMortgageItem, len(morgagePercentages))

	for i, percentage := range morgagePercentages {
		installment = utils.CalculateYearlyInstallment((housePrice * morgagePercentages[i] / 100), mortgageDuration, taeg)
		myData[i].Installment = installment / 12
		myData[i].Percentage = int(percentage)
		myData[i].Values = utils.SimulateSavingsCashVsMortgage(
			yearlyIncome,
			installment,
			mortgageDuration,
			taeg,
			yearlyGrowthgRate,
			housePrice*morgagePercentages[i]/100,
			housePrice)

	}

	return myData
}

func calculateMortgageCompare(
	yearlySaving float64,
	taeg float64,
	yearlyGrowthgRate float64,
	yearlySavingRate float64,
	durations []int,
	mortgageAmount float64,
	housePrice float64) []MortgageCompareItem {

	yearlyIncome := make([]float64, slices.Max(durations))
	incomeValue := (yearlySaving * yearlySavingRate)
	for i := range yearlyIncome {
		yearlyIncome[i] = incomeValue
	}

	myData := make([]MortgageCompareItem, len(durations))

	for i, duration := range durations {
		myData[i].Duration = int(duration)
		myData[i].Installment = utils.CalculateYearlyInstallment(mortgageAmount, duration, taeg)
		myData[i].Values = utils.SimulateSavings(
			yearlyIncome,
			myData[i].Installment,
			duration,
			taeg,
			yearlyGrowthgRate,
			mortgageAmount,
			housePrice)
		myData[i].Valid = utils.IsValidMortgage(myData[i].Installment, yearlySaving)

		fmt.Println(myData[i].Valid)
	}

	return myData
}

func main() {
	router := gin.Default()
	utils.Init()

	utils.InizializeRateLimiter()
	router.Use(cors.Default())

	anyoneCanJoin := true

	router.GET("/get_house_costs", func(c *gin.Context) {
		housePrice := utils.ToFloat64(c.DefaultQuery("house_price", "300000"))
		agencyFee := utils.ToFloat64(c.DefaultQuery("agency_fee", "0"))

		isFirstHouse := utils.ToBool(c.DefaultQuery("is_first_house", "false"))
		isSoldByBuilder := utils.ToBool(c.DefaultQuery("is_sold_by_builder", "false"))

		mortgageAmount := utils.ToFloat64(c.DefaultQuery("mortgage_amount", "0"))
		mortgageDuration := utils.ToInt(c.DefaultQuery("mortgage_duration", "0"))
		mortgageTAEG := utils.ToFloat64(c.DefaultQuery("mortgage_TAEG", "0"))

		if utils.IsAllowed(c.ClientIP()) {
			c.JSON(http.StatusOK, gin.H{
				"data": calculateAllCosts(housePrice, agencyFee, isFirstHouse, isSoldByBuilder, mortgageAmount, mortgageDuration, mortgageTAEG),
			})
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		}

	})
	router.GET("/get_cash_vs_mortgage", func(c *gin.Context) {
		housePrice := utils.ToFloat64(c.DefaultQuery("house_price", "300000"))
		yearlySaving := utils.ToFloat64(c.DefaultQuery("yearly_saving", "20000"))
		yearlyGrowthRate := utils.ToFloat64(c.DefaultQuery("yearly_growth_rate", "0.03"))
		yearlySavingRate := utils.ToFloat64(c.DefaultQuery("yearly_saving_rate", "0.3"))
		mortgageDuration := utils.ToInt(c.DefaultQuery("mortgage_duration", "0"))
		mortgageTAEG := utils.ToFloat64(c.DefaultQuery("mortgage_TAEG", "0"))
		mortgagePercentages := utils.ToFloatArray(c.DefaultQuery("mortgage_percentage", "25,50,75,80"), []float64{25, 50, 75, 80})

		UID := c.DefaultQuery("UID", "null")
		err := utils.IsUserStillPro(utils.FirebaseClient, UID)
		if err != nil || anyoneCanJoin {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if utils.IsAllowed(c.ClientIP()) {
			c.JSON(http.StatusOK, gin.H{
				"data": calculateCashVsMortgage(yearlySaving, mortgageDuration, mortgageTAEG, yearlyGrowthRate, yearlySavingRate, mortgagePercentages, housePrice),
			})
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		}

	})

	router.GET("/get_mortgage_compare", func(c *gin.Context) {
		housePrice := utils.ToFloat64(c.DefaultQuery("house_price", "300000"))
		yearlySaving := utils.ToFloat64(c.DefaultQuery("yearly_saving", "20000"))
		yearlyGrowthRate := utils.ToFloat64(c.DefaultQuery("yearly_growth_rate", "0.03"))
		yearlySavingRate := utils.ToFloat64(c.DefaultQuery("yearly_saving_rate", "0.3"))
		mortgageAmount := utils.ToFloat64(c.DefaultQuery("mortgage_amount", "0"))
		mortgageTAEG := utils.ToFloat64(c.DefaultQuery("mortgage_TAEG", "0"))
		durations := utils.ToIntArray(c.DefaultQuery("durations", "10,20,30"), []int{10, 20, 30})

		UID := c.DefaultQuery("UID", "null")
		err := utils.IsUserStillPro(utils.FirebaseClient, UID)
		if err != nil || anyoneCanJoin {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if utils.IsAllowed(c.ClientIP()) {
			c.JSON(http.StatusOK, gin.H{
				"data": calculateMortgageCompare(yearlySaving, mortgageTAEG, yearlyGrowthRate, yearlySavingRate, durations, mortgageAmount, housePrice),
			})
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		}

	})

	router.GET("/get-price-volume-data", func(c *gin.Context) {

		UID := c.DefaultQuery("UID", "null")
		err := utils.IsUserStillPro(utils.FirebaseClient, UID)
		if err != nil || anyoneCanJoin {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !utils.IsAllowed(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		// Get query parameters
		com := c.Query("com")
		zone := c.Query("zone")
		typeParam := c.Query("type")
		state := c.Query("state")
		mq := c.Query("mq")

		// Validate required parameters
		if com == "" || zone == "" || typeParam == "" || state == "" || mq == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters: com, zone, type, state, mq"})
			return
		}

		// Prepare parameters for Python API
		params := map[string]string{
			"com":   com,
			"zone":  zone,
			"type":  typeParam,
			"state": state,
			"mq":    mq,
		}

		// Call Python API
		responseBody, err := callPythonAPI("/get-price-volume-data", params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to call python API: %v", err)})
			return
		}

		// Log the raw response for debugging
		fmt.Printf("Python API raw response: %s\n", string(responseBody))

		// Parse and return the response
		var priceVolumeData PriceVolumeResponse
		if err := json.Unmarshal(responseBody, &priceVolumeData); err != nil {
			fmt.Printf("JSON parsing error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to parse python API response: %v", err)})
			return
		}

		c.JSON(http.StatusOK, priceVolumeData)
	})

	// Get municipalities info
	router.GET("/get-municipalities-info", func(c *gin.Context) {

		UID := c.DefaultQuery("UID", "null")
		err := utils.IsUserStillPro(utils.FirebaseClient, UID)
		if err != nil || anyoneCanJoin {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !utils.IsAllowed(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		com := c.Query("com")
		if com == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameter: com"})
			return
		}

		params := map[string]string{
			"com": com,
		}

		responseBody, err := callPythonAPI("/get-municipalities-info", params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to call python API: %v", err)})
			return
		}

		var municipalitiesInfo MunicipalitiesInfoResponse
		if err := json.Unmarshal(responseBody, &municipalitiesInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse python API response"})
			return
		}

		c.JSON(http.StatusOK, municipalitiesInfo)
	})

	// Get municipalities list
	router.GET("/get-municipalities-list", func(c *gin.Context) {

		UID := c.DefaultQuery("UID", "null")
		err := utils.IsUserStillPro(utils.FirebaseClient, UID)
		if err != nil || anyoneCanJoin {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !utils.IsAllowed(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		responseBody, err := callPythonAPI("/get-municipalities-list", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to call python API: %v", err)})
			return
		}

		var municipalitiesList MunicipalitiesListResponse
		if err := json.Unmarshal(responseBody, &municipalitiesList); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse python API response"})
			return
		}

		c.JSON(http.StatusOK, municipalitiesList)
	})

	// Update PRO
	router.POST("/is-pro", func(c *gin.Context) {

		type UID struct {
			Uid string `json:"uid" binding:"required"`
		}
		type ErrorResponse struct {
			Error   string `json:"error"`
			Message string `json:"message,omitempty"`
		}

		var req UID
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_request",
				Message: err.Error(),
			})
			return
		}

		utils.IsUserStillPro(utils.FirebaseClient, req.Uid)

		c.JSON(http.StatusOK, nil)
	})

	api := router.Group("/api")
	{
		api.POST("/create-payment-intent", utils.CreatePaymentIntent)
		router.POST("/api/confirm-payment", utils.ConfirmPayment)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local dev
	}
	router.Run(":" + port)

}
