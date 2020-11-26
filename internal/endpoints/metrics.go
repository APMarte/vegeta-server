package endpoints

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
	reqStsCode                                                *prometheus.GaugeVec
	resSuccessRatio                                           *prometheus.GaugeVec
	histogram                                                 *prometheus.HistogramVec

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

		var metricId string

		attackInfo := e.GetIdList(c)
		jsonReports := e.GetAllReports(c)

		for _, elem := range attackInfo {
			for _, element := range jsonReports {
				if elem.ID != element.ID {
					continue
				}

				metricId = element.ID
				p.reqCnt.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(float64(element.Requests))
				p.reqRt.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(float64(element.Rate))
				p.reqDur.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(milliseconds(element.Duration + element.Wait))
				p.reqAttck.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(milliseconds(element.Duration))
				p.reqWait.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(milliseconds(element.Wait))
				p.reqLatMean.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(milliseconds(element.Latencies.Mean))
				p.reqLat50th.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(milliseconds(element.Latencies.P50th))
				p.reqLat95th.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(milliseconds(element.Latencies.P95th))
				p.reqLat99th.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(milliseconds(element.Latencies.P99th))
				p.reqLatMax.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(milliseconds(element.Latencies.Max))
				p.resSuccessRatio.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration).Add(element.Success)

				for key, mapElem := range element.StatusCodes {
					p.reqStsCode.WithLabelValues(metricId, strconv.Itoa(elem.Params.Rate), elem.Params.Duration, key).Add(float64(mapElem))
				}
			}
		}

		//add histogram metrics to prometheus
		jsonHistogramReport := e.GetHistogram(metricId, c)

		for idx, value := range strings.Split(jsonHistogramReport, "\n") {

			if idx == 0 || value == "" {
				continue
			}

			ms := strings.Fields(value)[0]

			tofloat, err := strconv.ParseFloat(ms, 64)

			if err != nil {
				continue
			}

			p.histogram.WithLabelValues(metricId).Observe(float64(tofloat))
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
			p.reqStsCode = metric.(*prometheus.GaugeVec)
		case models.ResSuccessRatio:
			p.resSuccessRatio = metric.(*prometheus.GaugeVec)
		case models.Histogram:
			p.histogram = metric.(*prometheus.HistogramVec)
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

func milliseconds(d int) float64 {
	convertMs := time.Duration(d)

	msec, nsec := convertMs/time.Millisecond, convertMs%time.Millisecond
	return float64(msec) + float64(nsec)/1e6
}
