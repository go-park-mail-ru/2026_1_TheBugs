package complex

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/utils"
	usacase "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/complex"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ComplexServiceServer struct {
	complex.UnimplementedComplexServiceServer
	uc *usacase.UtilityCompanyUseCase
}

func NewComplexServiceServer(uc *usacase.UtilityCompanyUseCase) *ComplexServiceServer {
	return &ComplexServiceServer{
		uc: uc,
	}
}
func (c *ComplexServiceServer) GetComplex(
	ctx context.Context,
	req *complex.GetComplexRequest,
) (*complex.GetComplexResponse, error) {

	company, err := c.uc.GetUtilityCompany(ctx, req.Alias)
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return utils.MapComplex(company), nil
}

func (c *ComplexServiceServer) GetAllDevelopers(
	ctx context.Context,
	in *emptypb.Empty,
) (*complex.GetAllDevelopersResponse, error) {

	devs, err := c.uc.GetAllDevelopers(ctx)
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &complex.GetAllDevelopersResponse{
		Length:     int32(len(devs)),
		Developers: utils.MapDevelopers(devs),
	}, nil
}

func (c *ComplexServiceServer) GetComplexesByDeveloperID(
	ctx context.Context,
	req *complex.GetComplexesByDeveloperIDRequest,
) (*complex.GetComplexesByDeveloperIDResponse, error) {

	if req.DeveloperID == 0 {
		return nil, status.Error(codes.InvalidArgument, "developer_id is required")
	}

	comps, err := c.uc.GetAllByDeveloperID(ctx, int(req.DeveloperID))
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &complex.GetComplexesByDeveloperIDResponse{
		Length:    int32(len(comps)),
		Complexes: utils.MapComplexes(comps),
	}, nil
}
