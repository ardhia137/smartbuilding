package services

import (
	"errors"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"

	"gorm.io/gorm"
)

type SettingServiceImpl struct {
	haosRepo   repositories.SettingRepository
	torentRepo repositories.TorentRepository
}

func NewSettingService(haosRepo repositories.SettingRepository, torentRepo repositories.TorentRepository) services.SettingService {
	return &SettingServiceImpl{haosRepo: haosRepo, torentRepo: torentRepo}
}

func (s *SettingServiceImpl) CreateSetting(request entities.CreateSettingRequest) (*entities.SettingResponseCreate, error) {
	var createdSetting entities.Setting
	var createdTorent []entities.Torent

	db := s.haosRepo.WithTransaction()

	err := db.Transaction(func(tx *gorm.DB) error {
		haos := entities.Setting{
			NamaGedung:   request.NamaGedung,
			HaosURL:      request.HaosURL,
			HaosToken:    request.HaosToken,
			Scheduler:    request.Scheduler,
			HargaListrik: request.HargaListrik,
			JenisListrik: request.JenisListrik,
		}
		if err := tx.Create(&haos).Error; err != nil {
			return err
		}
		createdSetting = haos

		var torentList []entities.Torent
		for _, toren := range request.DataToren {
			torentList = append(torentList, entities.Torent{
				MonitoringName: toren.MonitoringName,
				KapasitasToren: toren.KapasitasToren,
				IDSetting:      haos.ID,
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

	response := entities.SettingResponseCreate{
		ID:           createdSetting.ID,
		NamaGedung:   createdSetting.NamaGedung,
		HaosURL:      createdSetting.HaosURL,
		HaosToken:    createdSetting.HaosToken,
		Scheduler:    createdSetting.Scheduler,
		HargaListrik: createdSetting.HargaListrik,
		JenisListrik: createdSetting.JenisListrik,
		DataToren:    createdTorent,
	}

	return &response, nil
}

func (s *SettingServiceImpl) GetAllCornJobs() ([]entities.SettingResponse, error) {
	haosList, err := s.haosRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []entities.SettingResponse
	for _, haos := range haosList {
		response = append(response, entities.SettingResponse{
			ID:           haos.ID,
			NamaGedung:   haos.NamaGedung,
			HaosURL:      haos.HaosURL,
			HaosToken:    haos.HaosToken,
			Scheduler:    haos.Scheduler,
			HargaListrik: haos.HargaListrik,
			JenisListrik: haos.JenisListrik,
		})
	}

	return response, nil
}

func (s *SettingServiceImpl) GetAllSetting(role string, userID uint) ([]entities.SettingResponse, error) {
	var haosList []entities.Setting
	var err error

	if role == "admin" {
		haosList, err = s.haosRepo.FindAll()
	} else {
		haosList, err = s.haosRepo.FindByUserId(userID)
	}

	if err != nil {
		return nil, err
	}

	// Jika tidak ada data, return error unauthorized
	if len(haosList) == 0 {
		return nil, errors.New("no data")
	}

	// Ambil status monitoring dari memory
	monitoringStatusMap := utils.GetMonitoringStatus()

	var response []entities.SettingResponse
	for _, haos := range haosList {
		// Cari status monitoring berdasarkan nama gedung
		var monitoringStatus []map[string]string
		if status, exists := monitoringStatusMap[haos.NamaGedung]; exists {
			monitoringStatus = status
		} else {
			// Default status jika tidak ditemukan
			monitoringStatus = []map[string]string{
				{"monitoring air": "unknown"},
				{"monitoring listrik": "unknown"},
			}
		}

		response = append(response, entities.SettingResponse{
			ID:               haos.ID,
			NamaGedung:       haos.NamaGedung,
			HaosURL:          haos.HaosURL,
			HaosToken:        haos.HaosToken,
			Scheduler:        haos.Scheduler,
			HargaListrik:     haos.HargaListrik,
			JenisListrik:     haos.JenisListrik,
			MonitoringStatus: monitoringStatus,
		})
	}

	return response, nil
}

func (s *SettingServiceImpl) GetSettingByID(id int) (*entities.SettingResponse, error) {
	haos, err := s.haosRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Ambil status monitoring dari memory
	monitoringStatusMap := utils.GetMonitoringStatus()

	// Cari status monitoring berdasarkan nama gedung
	var monitoringStatus []map[string]string
	if status, exists := monitoringStatusMap[haos.NamaGedung]; exists {
		monitoringStatus = status
	} else {
		// Default status jika tidak ditemukan
		monitoringStatus = []map[string]string{
			{"monitoring air": "unknown"},
			{"monitoring listrik": "unknown"},
		}
	}

	response := entities.SettingResponse{
		ID:               haos.ID,
		NamaGedung:       haos.NamaGedung,
		HaosURL:          haos.HaosURL,
		HaosToken:        haos.HaosToken,
		Scheduler:        haos.Scheduler,
		HargaListrik:     haos.HargaListrik,
		JenisListrik:     haos.JenisListrik,
		MonitoringStatus: monitoringStatus,
	}

	return &response, nil
}

func (s *SettingServiceImpl) UpdateSetting(id int, request entities.CreateSettingRequest) (*entities.SettingResponse, error) {
	haos, err := s.haosRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	haos.NamaGedung = request.NamaGedung
	haos.HaosURL = request.HaosURL
	haos.HaosToken = request.HaosToken
	haos.Scheduler = request.Scheduler
	haos.HargaListrik = request.HargaListrik
	haos.JenisListrik = request.JenisListrik

	updatedSetting, err := s.haosRepo.Update(haos)
	if err != nil {
		return nil, err
	}

	response := entities.SettingResponse{
		ID:           updatedSetting.ID,
		NamaGedung:   updatedSetting.NamaGedung,
		HaosURL:      updatedSetting.HaosURL,
		HaosToken:    updatedSetting.HaosToken,
		Scheduler:    updatedSetting.Scheduler,
		HargaListrik: updatedSetting.HargaListrik,
		JenisListrik: updatedSetting.JenisListrik,
	}

	return &response, nil
}

func (s *SettingServiceImpl) DeleteSetting(id int) error {
	return s.haosRepo.Delete(id)
}
