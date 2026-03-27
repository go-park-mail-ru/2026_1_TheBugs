package complex

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
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
// @Router /utility-companies/by-alias/{alias} [get]
func (h *UtilityCompanyHandler) GetUtilityCompany(w http.ResponseWriter, r *http.Request) {
	op := "UtilityCompanyHandler.GetUtilityCompany"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

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

// @Summary Get developers
// @Description Returns developers
// @Tags utility_complexes
// @Produce json
// @Success 200 {object} response.DevelopersResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /utility-companies/developers [get]
func (h *UtilityCompanyHandler) GetAllDevelopers(w http.ResponseWriter, r *http.Request) {
	op := "UtilityCompanyHandler.GetAllDevelopers"
	log := middleware.GetLogger(r.Context()).WithField("op", op)
	developers, err := h.uc.GetAllDevelopers(r.Context())
	if err != nil {
		log.Printf("h.uc.GetAllDevelopers: %v", err)
		utils.HandelError(w, err)
		return
	}

	res := response.DevelopersResponse{
		Len:        len(developers),
		Developers: developers,
	}

	utils.JSONResponse(w, http.StatusOK, res)

}

// @Summary Get utility complexies by developers
// @Description Returns utility complexies by developers
// @Tags utility_complexes
// @Produce json
// @Param developer_id query int true "Developer id"
// @Success 200 {object} response.CompaniesResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /utility-companies/ [get]
func (h *UtilityCompanyHandler) GetUtilityCompaniesByDeveloper(w http.ResponseWriter, r *http.Request) {
	op := "UtilityCompanyHandler.GetUtilityCompaniesByDeveloper"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	val := r.URL.Query().Get("developer_id")
	if val == "" {
		log.Errorf("val is empty")
		utils.WriteError(w, "query param 'developer_id' is requered", http.StatusBadRequest)
		return
	}

	developerID, err := strconv.Atoi(val)
	if err != nil {
		log.Errorf("Atoi: %s", err)
		utils.WriteError(w, "invalid value for query param 'developer_id'", http.StatusBadRequest)
		return
	}

	comps, err := h.uc.GetAllByDeveloperID(r.Context(), developerID)
	if err != nil {
		log.Printf("h.uc.GetAllByDeveloperID: %v", err)
		utils.HandelError(w, err)
		return
	}

	res := response.CompaniesResponse{
		Len:       len(comps),
		Companies: comps,
	}

	utils.JSONResponse(w, http.StatusOK, res)

}
