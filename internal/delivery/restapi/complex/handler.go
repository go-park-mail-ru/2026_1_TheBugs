package complex

import (
	"net/http"
	"strconv"

	complexGRPC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/response"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UtilityCompanyHandler struct {
	grpcClient complexGRPC.ComplexServiceClient
}

func NewUtilityCompanyHandler(client *grpc.ClientConn) *UtilityCompanyHandler {
	return &UtilityCompanyHandler{grpcClient: complexGRPC.NewComplexServiceClient(client)}
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
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		utils.EasyJSONResponse(w, http.StatusBadRequest, response.ErrorResponse{Error: "invalid alias"})
		return
	}
	company, err := h.grpcClient.GetComplex(r.Context(), &complexGRPC.GetComplexRequest{Alias: alias})
	if err != nil {
		log.Errorf("h.uc.GetUtilityCompany: %v", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.EasyJSONResponse(w, http.StatusOK, response.MapGRPCComplexToDTO(company))

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
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)
	res, err := h.grpcClient.GetAllDevelopers(r.Context(), &emptypb.Empty{})
	if err != nil {
		log.Printf("h.uc.GetAllDevelopers: %v", err)
		utils.HandelError(w, err)
		return
	}
	utils.EasyJSONResponse(w, http.StatusOK, response.MapGRPCDevelopersToDTO(res))

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
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

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

	res, err := h.grpcClient.GetComplexesByDeveloperID(r.Context(), &complexGRPC.GetComplexesByDeveloperIDRequest{DeveloperID: int32(developerID)})
	if err != nil {
		log.Printf("h.uc.GetAllByDeveloperID: %v", err)
		utils.HandelError(w, err)
		return
	}

	utils.EasyJSONResponse(w, http.StatusOK, response.MapGRPCComplexesToDTO(res.Complexes))

}
