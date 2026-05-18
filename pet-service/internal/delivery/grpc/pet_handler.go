package grpcdelivery

import (
	"context"
	"time"

	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-service/internal/entity"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/pet-service/internal/usecase"
	"github.com/ZhantoreGaineden/FinalProjectAdvancedProgramming2/proto/gen/petpb"
)

type PetHandler struct {
	petpb.UnimplementedPetServiceServer
	usecase *usecase.PetUsecase
}

func NewPetHandler(usecase *usecase.PetUsecase) *PetHandler {
	return &PetHandler{usecase: usecase}
}

func (h *PetHandler) CreatePet(ctx context.Context, req *petpb.CreatePetRequest) (*petpb.PetResponse, error) {
	pet := entity.Pet{
		Name:     req.GetName(),
		Category: req.GetCategory(),
		Breed:    req.GetBreed(),
		Age:      req.GetAge(),
		Price:    req.GetPrice(),
		Status:   req.GetStatus(),
	}

	created, err := h.usecase.CreatePet(ctx, pet)
	if err != nil {
		return nil, err
	}

	return &petpb.PetResponse{Pet: toProtoPet(created)}, nil
}

func (h *PetHandler) GetPet(ctx context.Context, req *petpb.GetPetRequest) (*petpb.PetResponse, error) {
	pet, err := h.usecase.GetPet(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &petpb.PetResponse{Pet: toProtoPet(pet)}, nil
}

func (h *PetHandler) ListPets(ctx context.Context, req *petpb.ListPetsRequest) (*petpb.ListPetsResponse, error) {
	pets, err := h.usecase.ListPets(ctx, req.GetCategory(), req.GetStatus())
	if err != nil {
		return nil, err
	}

	response := &petpb.ListPetsResponse{
		Pets: make([]*petpb.Pet, 0, len(pets)),
	}

	for _, pet := range pets {
		response.Pets = append(response.Pets, toProtoPet(pet))
	}

	return response, nil
}

func (h *PetHandler) UpdatePet(ctx context.Context, req *petpb.UpdatePetRequest) (*petpb.PetResponse, error) {
	pet := entity.Pet{
		ID:       req.GetId(),
		Name:     req.GetName(),
		Category: req.GetCategory(),
		Breed:    req.GetBreed(),
		Age:      req.GetAge(),
		Price:    req.GetPrice(),
		Status:   req.GetStatus(),
	}

	updated, err := h.usecase.UpdatePet(ctx, pet)
	if err != nil {
		return nil, err
	}

	return &petpb.PetResponse{Pet: toProtoPet(updated)}, nil
}

func (h *PetHandler) DeletePet(ctx context.Context, req *petpb.DeletePetRequest) (*petpb.DeletePetResponse, error) {
	err := h.usecase.DeletePet(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &petpb.DeletePetResponse{
		Success: true,
		Message: "pet deleted successfully",
	}, nil
}

func toProtoPet(pet entity.Pet) *petpb.Pet {
	createdAt := ""
	if !pet.CreatedAt.IsZero() {
		createdAt = pet.CreatedAt.Format(time.RFC3339)
	}

	return &petpb.Pet{
		Id:        pet.ID,
		Name:      pet.Name,
		Category:  pet.Category,
		Breed:     pet.Breed,
		Age:       pet.Age,
		Price:     pet.Price,
		Status:    pet.Status,
		CreatedAt: createdAt,
	}
}
