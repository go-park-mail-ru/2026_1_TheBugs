package poster

import (
	"context"
	"io"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/utils/mapper"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PosterServiceServer struct {
	poster.UnimplementedPosterServiceServer
	uc delivery.PostersUseCase
}

func NewPosterServiceServer(uc delivery.PostersUseCase) *PosterServiceServer {
	return &PosterServiceServer{
		uc: uc,
	}
}

func int32PtrToIntPtr(v *int32) *int {
	if v == nil {
		return nil
	}
	res := int(*v)
	return &res
}

func (s *PosterServiceServer) SearchPosters(ctx context.Context, req *poster.SearchPostersRequest) (*poster.SearchPostersResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "SearchPosters")

	if req.Limit <= 0 {
		log.Error("invalid field: limit")
		return nil, status.Error(codes.InvalidArgument, "limit must be positive")
	}

	if req.Offset < 0 {
		log.Error("invalid field: offset")
		return nil, status.Error(codes.InvalidArgument, "offset must be non-negative")
	}

	filters := dto.PostersFiltersDTO{
		Limit:            int(req.Limit),
		Offset:           int(req.Offset),
		SearchQuery:      req.SearchQuery,
		UtilityCompany:   req.UtilityCompany,
		Category:         req.Category,
		MaxPrice:         int32PtrToIntPtr(req.MaxPrice),
		MinPrice:         int32PtrToIntPtr(req.MinPrice),
		RoomCount:        int32PtrToIntPtr(req.RoomCount),
		MaxSquare:        int32PtrToIntPtr(req.MaxSquare),
		MinSquare:        int32PtrToIntPtr(req.MinSquare),
		Facilities:       req.Facilities,
		MaxFlatFloor:     int32PtrToIntPtr(req.MaxFlatFloor),
		MinFlatFloor:     int32PtrToIntPtr(req.MinFlatFloor),
		IsNotFirstFloor:  req.IsNotFirstFloor,
		IsNotLastFloor:   req.IsNotLastFloor,
		MaxBuildingFloor: int32PtrToIntPtr(req.MaxBuildingFloor),
		MinBuildingFloor: int32PtrToIntPtr(req.MinBuildingFloor),
	}

	response, err := s.uc.SearchPostersUseCase(ctx, filters)
	if err != nil {
		log.Errorf("s.uc.SearchPostersUseCase: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.SearchPostersResponse{
		Len:     int32(response.Len),
		Posters: mapper.PosterCardsDTOToProto(response.Posters),
	}, nil
}

func (s *PosterServiceServer) GetPosterByAlias(ctx context.Context, req *poster.GetPosterByAliasRequest) (*poster.GetPosterByAliasResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "GetPosterByAlias")

	if req.PosterAlias == "" {
		log.Error("missing required field: poster_alias")
		return nil, status.Error(codes.InvalidArgument, "poster_alias required")
	}

	var userID *int
	if req.UserId != nil {
		id := int(*req.UserId)
		userID = &id
	}

	posterDTO, err := s.uc.GetPosterByAliasUseCase(ctx, req.PosterAlias, userID)
	if err != nil {
		log.Errorf("s.uc.GetPosterByAliasUseCase: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.GetPosterByAliasResponse{
		Poster: mapper.PosterDTOToProto(posterDTO),
	}, nil
}

func (s *PosterServiceServer) GetPostersByUserID(ctx context.Context, req *poster.GetPostersByUserIDRequest) (*poster.GetPostersByUserIDResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "GetPostersByUserID")

	if req.UserId <= 0 {
		log.Error("invalid field: user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

	posters, err := s.uc.GetPosterByUserID(ctx, int(req.UserId))
	if err != nil {
		log.Errorf("s.uc.GetPosterByUserID: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.GetPostersByUserIDResponse{
		Posters: mapper.MyPostersDTOToProto(posters),
	}, nil
}

func (s *PosterServiceServer) AddViewPoster(ctx context.Context, req *poster.AddViewPosterRequest) (*poster.AddViewPosterResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "AddViewPoster")

	if req.Alias == "" {
		log.Error("missing required field: alias")
		return nil, status.Error(codes.InvalidArgument, "alias required")
	}

	if req.UserId <= 0 {
		log.Error("invalid field: user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

	err := s.uc.AddViewPoster(ctx, req.Alias, int(req.UserId))
	if err != nil {
		log.Errorf("s.uc.AddViewPoster: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.AddViewPosterResponse{}, nil
}

func (s *PosterServiceServer) GetViewsPoster(ctx context.Context, req *poster.GetViewsPosterRequest) (*poster.GetViewsPosterResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "GetViewsPoster")

	if req.Alias == "" {
		log.Error("missing required field: alias")
		return nil, status.Error(codes.InvalidArgument, "alias required")
	}

	views, err := s.uc.GetViewsPoster(ctx, req.Alias)
	if err != nil {
		log.Errorf("s.uc.GetViewsPoster: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.GetViewsPosterResponse{
		Views: int32(views),
	}, nil
}

func (s *PosterServiceServer) AddFavoritePoster(ctx context.Context, req *poster.AddFavoritePosterRequest) (*poster.AddFavoritePosterResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "AddFavoritePoster")

	if req.Alias == "" {
		log.Error("missing required field: alias")
		return nil, status.Error(codes.InvalidArgument, "alias required")
	}

	if req.UserId <= 0 {
		log.Error("invalid field: user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

	err := s.uc.AddFavoritePoster(ctx, req.Alias, int(req.UserId))
	if err != nil {
		log.Errorf("s.uc.AddFavoritePoster: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.AddFavoritePosterResponse{}, nil
}

func (s *PosterServiceServer) GetFavoritePosters(ctx context.Context, req *poster.GetFavoritePostersRequest) (*poster.GetFavoritePostersResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "GetFavoritePosters")

	if req.UserId <= 0 {
		log.Error("invalid field: user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

	response, err := s.uc.GetFavoritesPoster(ctx, int(req.UserId))
	if err != nil {
		log.Errorf("s.uc.GetFavoritesPoster: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.GetFavoritePostersResponse{
		Len:     int32(response.Len),
		Posters: mapper.PosterCardsDTOToProto(response.Posters),
	}, nil
}

func (s *PosterServiceServer) DeleteFavoritePoster(ctx context.Context, req *poster.DeleteFavoritePosterRequest) (*poster.DeleteFavoritePosterResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "DeleteFavoritePoster")

	if req.Alias == "" {
		log.Error("missing required field: alias")
		return nil, status.Error(codes.InvalidArgument, "alias required")
	}

	if req.UserId <= 0 {
		log.Error("invalid field: user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

	err := s.uc.DeleteFavoritePoster(ctx, req.Alias, int(req.UserId))
	if err != nil {
		log.Errorf("s.uc.DeleteFavoritePoster: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.DeleteFavoritePosterResponse{}, nil
}

func (s *PosterServiceServer) GetFavoritesCountPoster(ctx context.Context, req *poster.GetFavoritesCountPosterRequest) (*poster.GetFavoritesCountPosterResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "GetFavoritesCountPoster")

	if req.PosterAlias == "" {
		log.Error("missing required field: poster_alias")
		return nil, status.Error(codes.InvalidArgument, "poster_alias required")
	}

	var userID *int
	if req.UserId != nil {
		id := int(*req.UserId)
		userID = &id
	}

	count, isFavorite, err := s.uc.GetFavoritesCountPoster(ctx, req.PosterAlias, userID)
	if err != nil {
		log.Errorf("s.uc.GetFavoritesCountPoster: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.GetFavoritesCountPosterResponse{
		Count:      int32(count),
		IsFavorite: isFavorite,
	}, nil
}

func (s *PosterServiceServer) GetPostersByCoords(ctx context.Context, req *poster.GetPostersByCoordsRequest) (*poster.GetPostersByCoordsResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "GetPostersByCoords")

	if req.Bounds == nil {
		log.Error("missing required field: bounds")
		return nil, status.Error(codes.InvalidArgument, "bounds required")
	}

	if req.Bounds.Bbox == nil {
		log.Error("missing required field: bbox")
		return nil, status.Error(codes.InvalidArgument, "bbox required")
	}

	if req.Bounds.Bbox.SouthWest == nil || req.Bounds.Bbox.NorthEast == nil {
		log.Error("missing required fields: south_west and north_east")
		return nil, status.Error(codes.InvalidArgument, "south_west and north_east required")
	}

	if req.Filters == nil {
		log.Error("missing required field: filters")
		return nil, status.Error(codes.InvalidArgument, "filters required")
	}

	if req.Filters.Limit <= 0 {
		log.Error("invalid field: limit")
		return nil, status.Error(codes.InvalidArgument, "limit must be positive")
	}

	if req.Filters.Offset < 0 {
		log.Error("invalid field: offset")
		return nil, status.Error(codes.InvalidArgument, "offset must be non-negative")
	}

	bounds := dto.MapBounds{
		BBox: dto.BBox{
			SouthWest: dto.GeographyDTO{
				Lat: req.Bounds.Bbox.SouthWest.Lat,
				Lon: req.Bounds.Bbox.SouthWest.Lon,
			},
			NorthEast: dto.GeographyDTO{
				Lat: req.Bounds.Bbox.NorthEast.Lat,
				Lon: req.Bounds.Bbox.NorthEast.Lon,
			},
		},
		Zoom: int(req.Bounds.Zoom),
	}

	filters := dto.PostersFiltersDTO{
		Limit:            int(req.Filters.Limit),
		Offset:           int(req.Filters.Offset),
		SearchQuery:      req.Filters.SearchQuery,
		UtilityCompany:   req.Filters.UtilityCompany,
		Category:         req.Filters.Category,
		MaxPrice:         int32PtrToIntPtr(req.Filters.MaxPrice),
		MinPrice:         int32PtrToIntPtr(req.Filters.MinPrice),
		RoomCount:        int32PtrToIntPtr(req.Filters.RoomCount),
		MaxSquare:        int32PtrToIntPtr(req.Filters.MaxSquare),
		MinSquare:        int32PtrToIntPtr(req.Filters.MinSquare),
		Facilities:       req.Filters.Facilities,
		MaxFlatFloor:     int32PtrToIntPtr(req.Filters.MaxFlatFloor),
		MinFlatFloor:     int32PtrToIntPtr(req.Filters.MinFlatFloor),
		IsNotFirstFloor:  req.Filters.IsNotFirstFloor,
		IsNotLastFloor:   req.Filters.IsNotLastFloor,
		MaxBuildingFloor: int32PtrToIntPtr(req.Filters.MaxBuildingFloor),
		MinBuildingFloor: int32PtrToIntPtr(req.Filters.MinBuildingFloor),
	}

	response, err := s.uc.GetPostersByCoords(ctx, bounds, filters)
	if err != nil {
		log.Errorf("s.uc.GetPostersByCoords: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.GetPostersByCoordsResponse{
		Features: mapper.GeoJSONFeaturesDTOToProto(response.Posters),
		Len:      int32(response.Len),
	}, nil
}

func (s *PosterServiceServer) GetPostersByRadius(ctx context.Context, req *poster.GetPostersByRadiusRequest) (*poster.GetPostersByRadiusResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "GetPostersByRadius")

	if req.Point == nil {
		log.Error("missing required field: point")
		return nil, status.Error(codes.InvalidArgument, "point required")
	}

	posters, err := s.uc.GetPostersByRadius(ctx, dto.GeographyDTO{
		Lat: req.Point.Lat,
		Lon: req.Point.Lon,
	})
	if err != nil {
		log.Errorf("s.uc.GetPostersByRadius: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.GetPostersByRadiusResponse{
		Posters: mapper.MyPostersDTOToProto(posters),
	}, nil
}

func (s *PosterServiceServer) GenerateDescription(ctx context.Context, req *poster.GenerateDescriptionRequest) (*poster.GenerateDescriptionResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "GenerateDescription")

	if req.Category == "" {
		log.Error("missing required field: category")
		return nil, status.Error(codes.InvalidArgument, "category required")
	}

	if req.Area <= 0 {
		log.Error("invalid field: area")
		return nil, status.Error(codes.InvalidArgument, "area must be positive")
	}

	if req.FlatCategory == "" {
		log.Error("missing required field: flat_category")
		return nil, status.Error(codes.InvalidArgument, "flat_category required")
	}

	description, err := s.uc.GenerateDescription(ctx, dto.GenerateDescriptionDTO{
		Category:     req.Category,
		Area:         req.Area,
		FlatCategory: req.FlatCategory,
		City:         req.City,
		Features:     req.Features,
	})
	if err != nil {
		log.Errorf("s.uc.GenerateDescription: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.GenerateDescriptionResponse{
		Description: description,
	}, nil
}

func (s *PosterServiceServer) GetPriceHistoryPoster(ctx context.Context, req *poster.GetPriceHistoryPosterRequest) (*poster.GetPriceHistoryPosterResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "GetPriceHistoryPoster")

	if req.PosterAlias == "" {
		log.Error("missing required field: poster_alias")
		return nil, status.Error(codes.InvalidArgument, "poster_alias required")
	}

	history, err := s.uc.GetPriceHistoryPoster(ctx, req.PosterAlias)
	if err != nil {
		log.Errorf("s.uc.GetPriceHistoryPoster: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.GetPriceHistoryPosterResponse{
		History: mapper.PriceHistoryDTOsToProto(history),
	}, nil
}

func (s *PosterServiceServer) CreateFlatPoster(stream poster.PosterService_CreateFlatPosterServer) error {
	log := ctxLogger.GetLogger(stream.Context()).WithField("method", "CreateFlatPoster")

	var meta *poster.FlatPosterMeta

	files := make(map[int][]byte)
	photosMeta := make(map[int]*poster.FlatPosterPhotoMeta)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv: %v", err)
			return err
		}

		switch payload := req.Payload.(type) {

		case *poster.CreateFlatPosterRequest_PosterMeta:
			meta = payload.PosterMeta

		case *poster.CreateFlatPosterRequest_PhotoMeta:
			photoMeta := payload.PhotoMeta
			photosMeta[int(photoMeta.Order)] = photoMeta

		case *poster.CreateFlatPosterRequest_PhotoChunk:
			chunk := payload.PhotoChunk

			files[int(chunk.Order)] = append(
				files[int(chunk.Order)],
				chunk.Data...,
			)
		}
	}

	if meta == nil {
		log.Error("missing poster meta")
		return status.Error(codes.InvalidArgument, "poster meta is required")
	}

	dtoReq, err := mapper.FlatPosterProtoToDTO(meta, photosMeta, files)
	if err != nil {
		log.Errorf("mapCreateFlatPosterToDTO: %v", err)
		return status.Error(codes.Internal, err.Error())
	}

	result, err := s.uc.CreateFlatPoster(stream.Context(), dtoReq)
	if err != nil {
		log.Errorf("s.uc.CreateFlatPoster: %v", err)
		return utils.TranslateDomainsError(err)
	}

	return stream.SendAndClose(&poster.CreateFlatPosterResponse{
		Id:    int64(result.ID),
		Alias: result.Alias,
	})
}

func (s *PosterServiceServer) UpdateFlatPoster(stream poster.PosterService_UpdateFlatPosterServer) error {
	log := ctxLogger.GetLogger(stream.Context()).WithField("method", "UpdateFlatPoster")

	var meta *poster.UpdateFlatPosterMeta

	files := make(map[int][]byte)
	photosMeta := make(map[int]*poster.FlatPosterPhotoMeta)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("stream.Recv: %v", err)
			return err
		}

		switch payload := req.Payload.(type) {
		case *poster.UpdateFlatPosterRequest_PosterMeta:
			meta = payload.PosterMeta

		case *poster.UpdateFlatPosterRequest_PhotoMeta:
			photoMeta := payload.PhotoMeta
			photosMeta[int(photoMeta.Order)] = photoMeta

		case *poster.UpdateFlatPosterRequest_PhotoChunk:
			chunk := payload.PhotoChunk

			files[int(chunk.Order)] = append(
				files[int(chunk.Order)],
				chunk.Data...,
			)
		}
	}

	if meta == nil {
		log.Error("missing poster meta")
		return status.Error(codes.InvalidArgument, "poster meta is required")
	}

	if meta.Alias == "" {
		log.Error("missing alias")
		return status.Error(codes.InvalidArgument, "alias is required")
	}

	if meta.Poster == nil {
		log.Error("missing poster data")
		return status.Error(codes.InvalidArgument, "poster data is required")
	}

	dtoReq, err := mapper.FlatPosterProtoToDTO(meta.Poster, photosMeta, files)
	if err != nil {
		log.Errorf("mapFlatPosterMetaToDTO: %v", err)
		return err
	}

	result, err := s.uc.UpdateFlatPoster(stream.Context(), meta.Alias, dtoReq)
	if err != nil {
		log.Errorf("s.uc.UpdateFlatPoster: %v", err)
		return utils.TranslateDomainsError(err)
	}

	return stream.SendAndClose(&poster.UpdateFlatPosterResponse{
		Id:    int64(result.ID),
		Alias: result.Alias,
	})
}

func (s *PosterServiceServer) DeleteFlatPoster(ctx context.Context, req *poster.DeleteFlatPosterRequest) (*poster.DeleteFlatPosterResponse, error) {
	log := ctxLogger.GetLogger(ctx).WithField("method", "DeleteFlatPoster")

	if req.Alias == "" {
		log.Error("missing required field: alias")
		return nil, status.Error(codes.InvalidArgument, "alias required")
	}

	if req.UserId <= 0 {
		log.Error("invalid field: user_id")
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

	deletedPoster, err := s.uc.DeleteFlatPoster(ctx, req.Alias, int(req.UserId))
	if err != nil {
		log.Errorf("s.uc.DeleteFlatPoster: %s", err)
		return nil, utils.TranslateDomainsError(err)
	}

	return &poster.DeleteFlatPosterResponse{
		Id:    int64(deletedPoster.ID),
		Alias: deletedPoster.Alias,
	}, nil
}
