package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"vegeta-server/models"
	"vegeta-server/pkg/vegeta"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	reqCnt, reqRt                                             *prometheus.GaugeVec
	reqDur, reqAttck, reqWait                                 *prometheus.GaugeVec
	reqLatMean, reqLat50th, reqLat95th, reqLat99th, reqLatMax *prometheus.GaugeVec
	reqStsCode                                                *prometheus.SummaryVec
	resSuccessRatio                                           *prometheus.SummaryVec

	MetricsList []*models.Metric
}

func NewPrometheus(subsystem string) *Prometheus {

	var metricsList []*models.Metric

	for _, metric := range models.StandardMetrics {
		metricsList = append(metricsList, metric)
	}

	p := &Prometheus{
		MetricsList: metricsList,
	}

	p.registerMetrics(subsystem)

	return p
}

func (e *Endpoints) HandlerFunc(p *Prometheus) gin.HandlerFunc {
	return func(c *gin.Context) {

		attackInfo := e.GetIdList(c)
		jsonReports := e.GetAllReports(c)

		for _, elem := range attackInfo {
			for _, element := range jsonReports {
				if elem.ID != element.ID {
					continue
				}

				p.reqCnt.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(float64(element.Requests))
				p.reqRt.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(float64(element.Rate))
				p.reqDur.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(makeTimestamp(element.Duration + element.Wait))
				p.reqAttck.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(makeTimestamp(element.Duration))
				p.reqWait.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(makeTimestamp(element.Wait))
				p.reqLatMean.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(makeTimestamp(element.Latencies.Mean))
				p.reqLat50th.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(makeTimestamp(element.Latencies.P50th))
				p.reqLat95th.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(makeTimestamp(element.Latencies.P95th))
				p.reqLat99th.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(makeTimestamp(element.Latencies.P99th))
				p.reqLatMax.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(makeTimestamp(element.Latencies.Max))
				//p.resSuccessRatio.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(element.Success)
				//p.reqStsCode.WithLabelValues(element.ID, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Observe(element.StatusCodes)
			}

			//add histogram metrics to prometheus
			//TODO
			//jsonHistogramReport := e.GetHistogram(elem.ID, c)

		}

		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (p *Prometheus) registerMetrics(subsystem string) {

	for _, metricDef := range p.MetricsList {
		metric := models.NewMetric(metricDef, subsystem)
		if err := prometheus.Register(metric); err != nil {
			return
		}
		switch metricDef {
		case models.ReqCnt:
			p.reqCnt = metric.(*prometheus.GaugeVec)
		case models.ReqRt:
			p.reqRt = metric.(*prometheus.GaugeVec)
		case models.ReqDur:
			p.reqDur = metric.(*prometheus.GaugeVec)
		case models.ReqAttck:
			p.reqAttck = metric.(*prometheus.GaugeVec)
		case models.ReqWait:
			p.reqWait = metric.(*prometheus.GaugeVec)
		case models.ReqLatMean:
			p.reqLatMean = metric.(*prometheus.GaugeVec)
		case models.ReqLat50th:
			p.reqLat50th = metric.(*prometheus.GaugeVec)
		case models.ReqLat95th:
			p.reqLat95th = metric.(*prometheus.GaugeVec)
		case models.ReqLat99th:
			p.reqLat99th = metric.(*prometheus.GaugeVec)
		case models.ReqLatMax:
			p.reqLatMax = metric.(*prometheus.GaugeVec)
		case models.ReqStsCode:
			p.reqStsCode = metric.(*prometheus.SummaryVec)
		case models.ResSuccessRatio:
			p.resSuccessRatio = metric.(*prometheus.SummaryVec)
		}
		metricDef.MetricCollector = metric
	}
}

func (e *Endpoints) GetIdList(c *gin.Context) []*models.AttackBaseInfo {

	filterMap := make(models.FilterParams)
	filterMap["status"] = c.DefaultQuery("status", "completed")
	attackInfo := e.dispatcher.ListIds(
		filterMap,
	)

	return attackInfo
}

func (e *Endpoints) GetAllReports(c *gin.Context) []models.JSONReportResponse {
	resp := e.reporter.GetAll()
	jsonReports := make([]models.JSONReportResponse, 0)
	for _, report := range resp {
		var jsonReport models.JSONReportResponse
		err := json.Unmarshal(report, &jsonReport)
		if err != nil {
			ginErrInternalServerError(c, err)
			return nil
		}
		jsonReports = append(jsonReports, jsonReport)
	}
	return jsonReports
}

func (e *Endpoints) GetHistogram(id string, c *gin.Context) string {

	format := vegeta.NewFormat("hdrplot")

	resp, err := e.reporter.GetInFormat(id, format)
	if err != nil {

	}

	response := fmt.Sprintf("%s", resp)

	return response
}

func makeTimestamp(value int) float64 {
	return float64(value) / float64(time.Second)
}
