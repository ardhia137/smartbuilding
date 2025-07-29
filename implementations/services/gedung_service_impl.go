package services

import (
	"errors"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"

	"gorm.io/gorm"
)

type GedungServiceImpl struct {
	gedungRepo repositories.GedungRepository
	torentRepo repositories.TorentRepository
}

func NewGedungService(gedungRepo repositories.GedungRepository, torentRepo repositories.TorentRepository) services.GedungService {
	return &GedungServiceImpl{gedungRepo: gedungRepo, torentRepo: torentRepo}
}

func (s *GedungServiceImpl) CreateGedung(request entities.CreateGedungRequest) (*entities.GedungResponseCreate, error) {
	var createdGedung entities.Gedung
	var createdTorent []entities.Torent

	db := s.gedungRepo.WithTransaction()

	err := db.Transaction(func(tx *gorm.DB) error {
		gedung := entities.Gedung{
			NamaGedung:   request.NamaGedung,
			HaosURL:      request.HaosURL,
			HaosToken:    request.HaosToken,
			Scheduler:    request.Scheduler,
			HargaListrik: request.HargaListrik,
			JenisListrik: request.JenisListrik,
		}
		if err := tx.Create(&gedung).Error; err != nil {
			return err
		}
		createdGedung = gedung

		var torentList []entities.Torent
		for _, toren := range request.DataToren {
			torentList = append(torentList, entities.Torent{
				MonitoringName: toren.MonitoringName,
				KapasitasToren: toren.KapasitasToren,
				IDGedung:       gedung.ID,
			})
		}
		if err := tx.Create(&torentList).Error; err != nil {
			return err
		}
		createdTorent = torentList
		return nil
	})

	if err != nil {
		return nil, err
	}

	response := entities.GedungResponseCreate{
		ID:           createdGedung.ID,
		NamaGedung:   createdGedung.NamaGedung,
		HaosURL:      createdGedung.HaosURL,
		HaosToken:    createdGedung.HaosToken,
		Scheduler:    createdGedung.Scheduler,
		HargaListrik: createdGedung.HargaListrik,
		JenisListrik: createdGedung.JenisListrik,
		DataToren:    createdTorent,
	}

	return &response, nil
}

func (s *GedungServiceImpl) GetAllCornJobs() ([]entities.GedungResponse, error) {
	gedungList, err := s.gedungRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []entities.GedungResponse
	for _, gedung := range gedungList {
		response = append(response, entities.GedungResponse{
			ID:           gedung.ID,
			NamaGedung:   gedung.NamaGedung,
			HaosURL:      gedung.HaosURL,
			HaosToken:    gedung.HaosToken,
			Scheduler:    gedung.Scheduler,
			HargaListrik: gedung.HargaListrik,
			JenisListrik: gedung.JenisListrik,
		})
	}

	return response, nil
}

func (s *GedungServiceImpl) GetAllGedung(role string, userID uint) ([]entities.GedungResponse, error) {
	var gedungList []entities.Gedung
	var err error

	if role == "admin" {
		gedungList, err = s.gedungRepo.FindAll()
	} else {
		gedungList, err = s.gedungRepo.FindByUserId(userID)
	}

	if err != nil {
		return nil, err
	}

	// Jika tidak ada data, return error unauthorized
	if len(gedungList) == 0 {
		return nil, errors.New("no data")
	}

	// Ambil status monitoring dari memory
	monitoringStatusMap := utils.GetMonitoringStatus()

	var response []entities.GedungResponse
	for _, gedung := range gedungList {
		// Cari status monitoring berdasarkan nama gedung
		var monitoringStatus []map[string]string
		if status, exists := monitoringStatusMap[gedung.NamaGedung]; exists {
			monitoringStatus = status
		} else {
			// Default status jika tidak ditemukan
			monitoringStatus = []map[string]string{
				{"monitoring air": "unknown"},
				{"monitoring listrik": "unknown"},
			}
		}

		response = append(response, entities.GedungResponse{
			ID:               gedung.ID,
			NamaGedung:       gedung.NamaGedung,
			HaosURL:          gedung.HaosURL,
			HaosToken:        gedung.HaosToken,
			Scheduler:        gedung.Scheduler,
			HargaListrik:     gedung.HargaListrik,
			JenisListrik:     gedung.JenisListrik,
			MonitoringStatus: monitoringStatus,
		})
	}

	return response, nil
}

func (s *GedungServiceImpl) GetGedungByID(id int) (*entities.GedungResponse, error) {
	gedung, err := s.gedungRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Ambil status monitoring dari memory
	monitoringStatusMap := utils.GetMonitoringStatus()

	// Cari status monitoring berdasarkan nama gedung
	var monitoringStatus []map[string]string
	if status, exists := monitoringStatusMap[gedung.NamaGedung]; exists {
		monitoringStatus = status
	} else {
		// Default status jika tidak ditemukan
		monitoringStatus = []map[string]string{
			{"monitoring air": "unknown"},
			{"monitoring listrik": "unknown"},
		}
	}

	response := entities.GedungResponse{
		ID:               gedung.ID,
		NamaGedung:       gedung.NamaGedung,
		HaosURL:          gedung.HaosURL,
		HaosToken:        gedung.HaosToken,
		Scheduler:        gedung.Scheduler,
		HargaListrik:     gedung.HargaListrik,
		JenisListrik:     gedung.JenisListrik,
		MonitoringStatus: monitoringStatus,
	}

	return &response, nil
}

func (s *GedungServiceImpl) UpdateGedung(id int, request entities.CreateGedungRequest) (*entities.GedungResponse, error) {
	gedung, err := s.gedungRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	gedung.NamaGedung = request.NamaGedung
	gedung.HaosURL = request.HaosURL
	gedung.HaosToken = request.HaosToken
	gedung.Scheduler = request.Scheduler
	gedung.HargaListrik = request.HargaListrik
	gedung.JenisListrik = request.JenisListrik

	updatedGedung, err := s.gedungRepo.Update(gedung)
	if err != nil {
		return nil, err
	}

	response := entities.GedungResponse{
		ID:           updatedGedung.ID,
		NamaGedung:   updatedGedung.NamaGedung,
		HaosURL:      updatedGedung.HaosURL,
		HaosToken:    updatedGedung.HaosToken,
		Scheduler:    updatedGedung.Scheduler,
		HargaListrik: updatedGedung.HargaListrik,
		JenisListrik: updatedGedung.JenisListrik,
	}

	return &response, nil
}

func (s *GedungServiceImpl) DeleteGedung(id int) error {
	return s.gedungRepo.Delete(id)
}
