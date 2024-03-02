package controllers

import (
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// generate random data for line chart
func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func (ge *GinEngine) RenderLineChart(c *gin.Context) {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	// line.SetGlobalOptions(
	// 	charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
	// 	charts.WithTitleOpts(opts.Title{
	// 		Title:    "Line example in Westeros theme",
	// 		Subtitle: "Line chart rendered by the http server this time",
	// 	}))

	// Put data into instance
	line.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Category A", generateLineItems()).
		AddSeries("Category B", generateLineItems()).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(c.Writer)
}
