package handler

import (
	"net/http"

	"github.com/Jeff1507/go-to-save-api/internal/core/store/query"
	"github.com/Jeff1507/go-to-save-api/internal/infra/api/request"
	"github.com/gin-gonic/gin"
)

type StatsResponse struct {
	TodaySales           int64 `json:"todaySales"`
	ActiveBags           int64 `json:"activeBags"`
	CompletedOrdersMonth int64 `json:"completedOrdersMonth"`
}

type GetStatsHandler struct {
	getStats *query.GetStats
}

func NewGetStatsHandler(getStats *query.GetStats) *GetStatsHandler {
	return &GetStatsHandler{
		getStats: getStats,
	}
}

func (h *GetStatsHandler) Handle(c *gin.Context) {
	storeID, ok := request.StoreID(c.Request.Context())
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	stats, err := h.getStats.Execute(c.Request.Context(), storeID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": StatsResponse{
			TodaySales:           stats.TodaySales,
			ActiveBags:           stats.ActiveBags,
			CompletedOrdersMonth: stats.CompletedOrdersMonth,
		},
	})
}
