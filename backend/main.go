package main

import (
	"net/http"

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

func main() {
	router := gin.Default()

	// Allow all origins (DEVELOPEMENT ONLY!!! REMOVE IN PRODUCTION)
	router.Use(cors.Default())

	router.GET("/get_house_prices", func(c *gin.Context) {
		housePrice := utils.ToFloat64(c.DefaultQuery("house_price", "300000"))
		agencyFee := utils.ToFloat64(c.DefaultQuery("agency_fee", "0"))

		isFirstHouse := utils.ToBool(c.DefaultQuery("is_first_house", "false"))
		isSoldByBuilder := utils.ToBool(c.DefaultQuery("is_sold_by_builder", "false"))

		mortgageAmount := utils.ToFloat64(c.DefaultQuery("mortgage_amount", "0"))
		mortgageDuration := utils.ToFloat64(c.DefaultQuery("mortgage_duration", "0"))
		mortgageTAEG := utils.ToFloat64(c.DefaultQuery("mortgage_TAEG", "0"))

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
		tassaArchivioMortgage = min(mortgageAmount*0.000146, 139.4)
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
		if agencyFee > 0 {
			agencyAmount = agencyFee * housePrice * 1.22
		}

		var registrazionePreliminare float64 = 200 + 16
		var tassaIpotecaria float64 = 90
		var tassaIpotecariaMortgage float64 = 35
		var impostaDiBollo float64 = 230

		var certificatoStatoLibero float64 = 16

		responseData := []CostItem{
			{
				Name:     "Property Price",
				Category: "House price",
				Value:    housePrice,
				Estimate: false,
				Info:     "House buy price",
			},

			{
				Name:     "Monthly Mortgage Installment",
				Category: "Installment",
				Value:    mortgageInstalment,
				Estimate: false,
				Info:     "Mortgage intrests",
			},
			{
				Name:     "Agency",
				Category: "Agency",
				Value:    agencyAmount,
				Estimate: false,
				Info:     "Agency fee",
			},
			{
				Name:     "VAT",
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

		c.JSON(http.StatusOK, gin.H{
			"data": responseData,
		})
	})
	router.Run(":8080")

}
