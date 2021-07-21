package public

import (
	"github.com/iamrz1/ab-auth/model"
	"github.com/iamrz1/ab-auth/model/response"
	"github.com/iamrz1/ab-auth/utils"
	"net/http"
)

// listBDArea godoc
// @Summary Fetch BD area presets (division, district, sub-district)
// @Description Get a list of BD areas under selected parent (slug value). No parent returns list of divisions. Division as parent will return districts and so on)
// @Tags Common
// @Accept  json
// @Produce  json
// @Param parent query string false "Default value: empty-string"
// @Param page query integer false "Default value: 1"
// @Param limit query integer false "Default value: 10"
// @Success 200 {object} response.BDLocationListSuccessRes
// @Failure 500 {object} response.EmptyListErrorRes
// @Router /api/v1/public/bd-area [get]
func (pr *publicRouter) listBDArea(w http.ResponseWriter, r *http.Request) {
	parent := r.URL.Query().Get("parent")
	page, limit := utils.GetPageLimit(r)

	req := model.BDLocationReq{Parent: parent, Page: page, Limit: limit}

	res, count, err := pr.Services.CustomerService.GetBDLocation(r.Context(), &req)
	if err != nil {
		utils.HandleListError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Success. Very nice!", res, response.GetListMeta(page, limit, count), true)
}
