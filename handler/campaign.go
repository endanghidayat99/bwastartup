package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error get Campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of Campaigns", http.StatusOK, "successs", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}
