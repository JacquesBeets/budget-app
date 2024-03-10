package controllers

import (
	"budget-app/internal/models"
	"budget-app/internal/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func returnOptions(xAxisData []string, seriesData []opts.LineData) map[string]interface{} {
	// Manually create the chart options
	options := map[string]interface{}{
		"grid": map[string]interface{}{
			"left":         "0%",
			"right":        "0%",
			"bottom":       "0%",
			"top":          "0%",
			"containLabel": false,
		},
		"color": []string{"#FFFFFF"},
		"tooltip": map[string]interface{}{
			"show":      true,
			"formatter": "R {c}",
		},
		"legend": map[string]interface{}{
			"show": false,
		},
		"xAxis": []map[string]interface{}{
			{
				"data":  xAxisData,
				"scale": true,
				"splitLine": map[string]interface{}{
					"show": false,
					"lineStyle": map[string]interface{}{
						"color":   "white",
						"opacity": 0.1,
					},
				},
			},
		},
		"yAxis": []map[string]interface{}{
			{
				"scale": true,
				"splitLine": map[string]interface{}{
					"show": true,
					"lineStyle": map[string]interface{}{
						"color":   "white",
						"opacity": 0.1,
					},
				},
			},
		},
		"series": []map[string]interface{}{
			{
				"name": "Category A",
				"type": "bar",
				"data": seriesData,
			},
		},
	}
	return options
}

func (ge *GinEngine) RenderLineChart(c *gin.Context) {
	var chartLabels []string
	var lineData []opts.LineData

	var portfolioHistory models.CryptoPortfolioTotalHistories
	_, err := portfolioHistory.FetchAll(ge.db())
	if err != nil {
		ge.ReturnErrorPage(c, err)
		return
	}

	for _, history := range portfolioHistory {
		chartLabels = append(chartLabels, fmt.Sprintf("%d", history.ID))
		lineData = append(lineData, opts.LineData{
			Value: utils.FormatNumberToTwoDecimalPlaces(history.TotalValueZar),
		})
	}

	ge.Router.LoadHTMLFiles(LineChart)

	options := returnOptions(chartLabels, lineData)

	// Convert the chart options to a JSON string
	optionsJSON, err := json.Marshal(options)
	if err != nil {
		log.Fatalf("Error marshaling chart options: %v", err)
	}

	// Pass the options to your template
	c.HTML(http.StatusOK, "portfolio_line.html", gin.H{
		"line": template.JS(optionsJSON),
	})
}
