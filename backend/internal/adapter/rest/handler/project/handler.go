package project

import (
	inproject "github.com/Unknowns24/akritas/backend/internal/core/ports/in/project"
)

type Handler struct {
	projects         inproject.Create
	getter           inproject.Get
	lister           inproject.List
	updater          inproject.Update
	monitoringGetter inproject.GetMonitoring
	monitoringPutter inproject.PutMonitoring
	paginationSecret string
}

func NewHandler(
	projects inproject.Create,
	getter inproject.Get,
	lister inproject.List,
	updater inproject.Update,
	monitoringGetter inproject.GetMonitoring,
	monitoringPutter inproject.PutMonitoring,
	paginationSecret string,
) *Handler {
	return &Handler{
		projects: projects, getter: getter, lister: lister, updater: updater,
		monitoringGetter: monitoringGetter, monitoringPutter: monitoringPutter,
		paginationSecret: paginationSecret,
	}
}
