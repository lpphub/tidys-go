//go:build wireinject
// +build wireinject

//go:generate wire
package logic

import (
	"tidys-go/infra"
	"tidys-go/logic/auth"
	"tidys-go/logic/notes"
	notesRepo "tidys-go/logic/notes/repository"
	"tidys-go/logic/spaces"
	spacesRepo "tidys-go/logic/spaces/repository"
	"tidys-go/logic/tags"
	tagsRepo "tidys-go/logic/tags/repository"
	"tidys-go/logic/user"

	"github.com/google/wire"
)

type AppService struct {
	Auth  *auth.Service
	User  *user.Service
	Tag   *tags.Service
	Space *spaces.Service
	Note  *notes.Service
}

var svcSet = wire.NewSet(
	user.NewService,
	auth.NewService,
	tags.NewService,
	spaces.NewService,
	notes.NewService,
)

var providerSet = wire.NewSet(
	infra.ProvideDB,
	infra.ProvideRDB,
)

var repoSet = wire.NewSet(
	user.NewUserRepo,
	tagsRepo.NewTagRepo,
	tagsRepo.NewTagGroupRepo,
	spacesRepo.NewSpaceRepo,
	spacesRepo.NewMemberRepo,
	spacesRepo.NewInviteRepo,
	notesRepo.NewNoteRepo,
)

func initialize() *AppService {
	wire.Build(providerSet, repoSet, svcSet, wire.Struct(new(AppService), "*"))
	return nil
}
