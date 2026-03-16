package complex

import (
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/response"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/complex"
)

type UtilityCompanyHandler struct {
	uc *complex.UtilityCompanyUseCase
}

func NewUtilityCompanyHandler(uc *complex.UtilityCompanyUseCase) *UtilityCompanyHandler {
	return &UtilityCompanyHandler{uc: uc}
}

// @Summary Get utility complex by alias
// @Description Returns utility complex details along with its photos by the given alias
// @Tags utility_complexes
// @Produce json
// @Param alias path string true "Utility Complex Alias"
// @Success 200 {object} dto.UtilityCompanyDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /utility_companies/alias/{alias} [get]
func (h *UtilityCompanyHandler) GetUtilityCompany(w http.ResponseWriter, r *http.Request) {
	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, response.ErrorResponse{Error: "invalid alias"})
		return
	}
	UtilityCompany, err := h.uc.GetUtilityCompany(r.Context(), alias)
	if err != nil {
		log.Printf("h.uc.GetUtilityCompany: %v", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, UtilityCompany)

}
