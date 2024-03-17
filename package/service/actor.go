package service

import (
	filmoteka "vk_restAPI"
	"vk_restAPI/package/repository"
)

type ActorService struct {
	repo repository.Actors
}

type ActorsWithMoviesService struct {
	repo repository.ActorsWithMovies
}

func NewActorService(repo repository.Actors) *ActorService {
	return &ActorService{repo: repo}
}

func NewActorsWithMoviesService(repo repository.ActorsWithMovies) *ActorsWithMoviesService {
	return &ActorsWithMoviesService{repo: repo}
}

func (a *ActorService) CreateActor(actor filmoteka.Actors) (int, error) {
	return a.repo.CreateActor(actor)
}

func (a *ActorsWithMoviesService) GetActors() ([]filmoteka.ActorsWithMovies, error) {
	return a.repo.GetActors()
}

func (a *ActorsWithMoviesService) GetActorById(actorId int) (filmoteka.ActorsWithMovies, error) {
	return a.repo.GetActorById(actorId)
}

func (a *ActorService) UpdateActor(actorId int, input filmoteka.UpdateActors) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return a.repo.UpdateActor(actorId, input)
}

func (a *ActorService) DeleteActor(actorId int) error {
	return a.repo.DeleteActor(actorId)
}
