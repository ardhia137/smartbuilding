package services

import (
	"fmt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"strconv"
	"strings"
	"time"
)

type monitoringDataServiceImpl struct {
	monitoringDataRepository repositories.MonitoringDataRepository
}

func NewMonitoringDataService(monitorRepo repositories.MonitoringDataRepository) services.MonitoringDataService {
	return &monitoringDataServiceImpl{monitorRepo}
}

func (s *monitoringDataServiceImpl) SaveMonitoringData(request entities.CreateMonitoringDataRequest) (entities.MonitoringDataResponse, error) {
	monitoringData := entities.MonitoringData{
		MonitoringName:  request.MonitoringName,
		MonitoringValue: request.MonitoringValue,
	}

	createdData, err := s.monitoringDataRepository.SaveMonitoringData(monitoringData)
	if err != nil {
		return entities.MonitoringDataResponse{}, err
	}

	response := entities.MonitoringDataResponse{
		ID:              createdData.ID,
		MonitoringName:  createdData.MonitoringName,
		MonitoringValue: createdData.MonitoringValue,
		CreatedAt:       createdData.CreatedAt,
		UpdatedAt:       createdData.UpdatedAt,
	}

	return response, nil
}

func (s *monitoringDataServiceImpl) GetAirMonitoringData() ([]entities.GetAirDataResponse, error) {
	monitoringData, err := s.monitoringDataRepository.GetAirMonitoringData()
	if err != nil {
		return nil, err
	}

	monitoringDataHarian, err := s.monitoringDataRepository.GetAirMonitoringDataHarian()
	if err != nil {
		return nil, err
	}

	var kapasitasToren, airMasuk string
	var totalAirKeluar float64
	var createdAt, updatedAt time.Time

	dataPenggunaanHarian := make(map[string][]entities.PenggunaanAir)
	dataPenggunaanMingguan := make(map[string]map[string]float64)
	dataPenggunaanTahunan := make(map[string]map[string]float64)
	latestWaterFlow := make(map[string]float64)
	latestCreatedAt := make(map[string]time.Time)

	for _, data := range monitoringData {
		switch data.MonitoringName {
		case "monitoring_air_kapasitas_toren":
			kapasitasToren = data.MonitoringValue
		case "monitoring_air_total_water_flow_pipa_air_masuk":
			airMasuk = data.MonitoringValue
		default:
			if strings.HasPrefix(data.MonitoringName, "monitoring_air_total_water_flow_") {
				volume, _ := strconv.ParseFloat(strings.TrimSuffix(data.MonitoringValue, " L"), 64)

				// Ambil nama pipa dari monitoring name
				pipa := strings.TrimPrefix(data.MonitoringName, "monitoring_air_total_water_flow_")

				// Simpan nilai terbaru berdasarkan CreatedAt
				if lastTime, exists := latestCreatedAt[pipa]; !exists || data.CreatedAt.After(lastTime) {
					latestWaterFlow[pipa] = volume
					latestCreatedAt[pipa] = data.CreatedAt
				}
			}
		}

		createdAt = data.CreatedAt
		updatedAt = data.UpdatedAt
	}

	// Hitung total air keluar hanya dari data terbaru tiap pipa
	totalAirKeluar = 0
	for _, volume := range latestWaterFlow {
		totalAirKeluar += volume
	}
	indikatorLevel, _ := strconv.ParseFloat(strings.TrimSuffix(kapasitasToren, " %"), 64)
	volumeSensor := 5100 * (indikatorLevel / 100)

	now := time.Now()
	year, month, _ := now.Date()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	for _, harian := range monitoringDataHarian {
		if strings.HasPrefix(harian.MonitoringName, "monitoring_air_total_water_flow_") {
			pipa := strings.TrimPrefix(harian.MonitoringName, "monitoring_air_total_water_flow_")
			volume, _ := strconv.ParseFloat(strings.TrimSuffix(harian.MonitoringValue, " L"), 64)
			hari := getHariIndonesia(harian.CreatedAt.Weekday())

			if harian.CreatedAt.After(startOfWeek) && harian.CreatedAt.Before(endOfWeek) {
				dataPenggunaanHarian[hari] = append(dataPenggunaanHarian[hari], entities.PenggunaanAir{
					Pipa:   pipa,
					Volume: fmt.Sprintf("%.0f L", volume),
				})
			}

			_, minggu := harian.CreatedAt.ISOWeek()
			if harian.CreatedAt.Month() != time.Month(minggu) {
				minggu = 1
			}
			mingguanKey := fmt.Sprintf("Minggu %d", minggu)

			if harian.CreatedAt.Year() == year && harian.CreatedAt.Month() == month {
				if dataPenggunaanMingguan[mingguanKey] == nil {
					dataPenggunaanMingguan[mingguanKey] = make(map[string]float64)
				}
				dataPenggunaanMingguan[mingguanKey][pipa] += volume
			}

			bulan := getBulanIndonesia(harian.CreatedAt.Month())

			if harian.CreatedAt.Year() == year {
				if dataPenggunaanTahunan[bulan] == nil {
					dataPenggunaanTahunan[bulan] = make(map[string]float64)
				}
				dataPenggunaanTahunan[bulan][pipa] += volume
			}
		}
	}

	convertMingguanTahunan := func(data map[string]float64) []entities.PenggunaanAir {
		result := make([]entities.PenggunaanAir, 0, len(data))
		for pipa, volume := range data {
			result = append(result, entities.PenggunaanAir{
				Pipa:   pipa,
				Volume: fmt.Sprintf("%.0f L", volume),
			})
		}
		return result
	}

	response := entities.GetAirDataResponse{
		KapasitasToren:         kapasitasToren,
		AirMasuk:               airMasuk,
		AirKeluar:              fmt.Sprintf("%.0f L", totalAirKeluar),
		VolumeSensor:           fmt.Sprintf("%.0f L", volumeSensor),
		DataPenggunaanHarian:   dataPenggunaanHarian,
		DataPenggunaanMingguan: make(map[string][]entities.PenggunaanAir),
		DataPenggunaanTahunan:  make(map[string][]entities.PenggunaanAir),
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}

	for minggu, data := range dataPenggunaanMingguan {
		response.DataPenggunaanMingguan[minggu] = convertMingguanTahunan(data)
	}

	for bulan, data := range dataPenggunaanTahunan {
		response.DataPenggunaanTahunan[bulan] = convertMingguanTahunan(data)
	}

	return []entities.GetAirDataResponse{response}, nil
}
func (s *monitoringDataServiceImpl) GetListrikMonitoringData() (entities.GetListrikDataResponse, error) {
	monitoringData, err := s.monitoringDataRepository.GetListrikMonitoringData()
	if err != nil {
		return entities.GetListrikDataResponse{}, err
	}

	monitoringDataHarian, err := s.monitoringDataRepository.GetListrikMonitoringDataHarian()
	if err != nil {
		return entities.GetListrikDataResponse{}, err
	}

	var createdAt, updatedAt time.Time
	var totalWatt float64
	totalDaya := map[string]float64{
		"LT1": 0,
		"LT2": 0,
		"LT3": 0,
		"LT4": 0,
	}

	dataPenggunaanHarian := make(map[string]map[int]entities.PenggunaanListrik)
	dataBiayaHarian := make(map[string]map[int]entities.BiayaListrik)

	dataPenggunaanMingguan := make(map[string]map[int]entities.PenggunaanListrik)
	dataBiayaMingguan := make(map[string]map[int]entities.BiayaListrik)

	dataPenggunaanTahunan := make(map[string]map[int]entities.PenggunaanListrik)
	dataBiayaTahunan := make(map[string]map[int]entities.BiayaListrik)

	const (
		tegangan     = 220.0
		tarifListrik = 1900.0
	)

	for _, data := range monitoringData {
		arus, _ := strconv.ParseFloat(strings.TrimSuffix(data.MonitoringValue, " A"), 64)
		daya := (tegangan * arus) / 1000 // Konversi ke kW

		switch {
		case strings.Contains(data.MonitoringName, "l1"):
			totalDaya["LT1"] += daya
		case strings.Contains(data.MonitoringName, "l2"):
			totalDaya["LT2"] += daya
		case strings.Contains(data.MonitoringName, "l3"):
			totalDaya["LT3"] += daya
		case strings.Contains(data.MonitoringName, "lt4"):
			totalDaya["LT4"] += daya
		}

		totalWatt += daya
		createdAt = data.CreatedAt
		updatedAt = data.UpdatedAt
	}

	biayaLT1 := fmt.Sprintf("Rp. %.0f", totalDaya["LT1"]*tarifListrik)
	biayaLT2 := fmt.Sprintf("Rp. %.0f", totalDaya["LT2"]*tarifListrik)
	biayaLT3 := fmt.Sprintf("Rp. %.0f", totalDaya["LT3"]*tarifListrik)
	biayaLT4 := fmt.Sprintf("Rp. %.0f", totalDaya["LT4"]*tarifListrik)

	now := time.Now()
	year, month, _ := now.Date()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	for _, harian := range monitoringDataHarian {
		if strings.Contains(harian.MonitoringName, "arus_listrik") {
			arus, _ := strconv.ParseFloat(strings.TrimSuffix(harian.MonitoringValue, " A"), 64)
			daya := (tegangan * arus) / 1000 // Konversi ke kW

			var lantai int
			switch {
			case strings.Contains(harian.MonitoringName, "l1"):
				lantai = 1
			case strings.Contains(harian.MonitoringName, "l2"):
				lantai = 2
			case strings.Contains(harian.MonitoringName, "l3"):
				lantai = 3
			case strings.Contains(harian.MonitoringName, "lt4"):
				lantai = 4
			}

			hari := getHariIndonesia(harian.CreatedAt.Weekday())

			if harian.CreatedAt.After(startOfWeek) && harian.CreatedAt.Before(endOfWeek) {
				if dataPenggunaanHarian[hari] == nil {
					dataPenggunaanHarian[hari] = make(map[int]entities.PenggunaanListrik)
				}
				if dataBiayaHarian[hari] == nil {
					dataBiayaHarian[hari] = make(map[int]entities.BiayaListrik)
				}

				if existing, ok := dataPenggunaanHarian[hari][lantai]; ok {
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " kW"), 64)
					existingValue += daya
					existing.Value = fmt.Sprintf("%.2f kW", existingValue)
					dataPenggunaanHarian[hari][lantai] = existing
				} else {
					dataPenggunaanHarian[hari][lantai] = entities.PenggunaanListrik{
						Lantai: lantai,
						Value:  fmt.Sprintf("%.2f kW", daya),
					}
				}

				if existing, ok := dataBiayaHarian[hari][lantai]; ok {
					existingBiaya, _ := strconv.ParseFloat(strings.TrimPrefix(existing.Biaya, "Rp. "), 64)
					existingBiaya += daya * tarifListrik
					existing.Biaya = fmt.Sprintf("Rp. %.0f", existingBiaya)
					dataBiayaHarian[hari][lantai] = existing
				} else {
					dataBiayaHarian[hari][lantai] = entities.BiayaListrik{
						Lantai: lantai,
						Biaya:  fmt.Sprintf("Rp. %.0f", daya*tarifListrik),
					}
				}
			}

			_, minggu := harian.CreatedAt.ISOWeek()
			if harian.CreatedAt.Month() != time.Month(minggu) {
				minggu = 1
			}
			mingguanKey := fmt.Sprintf("Minggu %d", minggu)

			if harian.CreatedAt.Year() == year && harian.CreatedAt.Month() == month {
				if dataPenggunaanMingguan[mingguanKey] == nil {
					dataPenggunaanMingguan[mingguanKey] = make(map[int]entities.PenggunaanListrik)
				}
				if dataBiayaMingguan[mingguanKey] == nil {
					dataBiayaMingguan[mingguanKey] = make(map[int]entities.BiayaListrik)
				}

				if existing, ok := dataPenggunaanMingguan[mingguanKey][lantai]; ok {
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " kW"), 64)
					existingValue += daya
					existing.Value = fmt.Sprintf("%.2f kW", existingValue)
					dataPenggunaanMingguan[mingguanKey][lantai] = existing
				} else {
					dataPenggunaanMingguan[mingguanKey][lantai] = entities.PenggunaanListrik{
						Lantai: lantai,
						Value:  fmt.Sprintf("%.2f kW", daya),
					}
				}

				if existing, ok := dataBiayaMingguan[mingguanKey][lantai]; ok {
					existingBiaya, _ := strconv.ParseFloat(strings.TrimPrefix(existing.Biaya, "Rp. "), 64)
					existingBiaya += daya * tarifListrik
					existing.Biaya = fmt.Sprintf("Rp. %.0f", existingBiaya)
					dataBiayaMingguan[mingguanKey][lantai] = existing
				} else {
					dataBiayaMingguan[mingguanKey][lantai] = entities.BiayaListrik{
						Lantai: lantai,
						Biaya:  fmt.Sprintf("Rp. %.0f", daya*tarifListrik),
					}
				}
			}

			bulanHarian := getBulanIndonesia(harian.CreatedAt.Month())

			if harian.CreatedAt.Year() == year {
				if dataPenggunaanTahunan[bulanHarian] == nil {
					dataPenggunaanTahunan[bulanHarian] = make(map[int]entities.PenggunaanListrik)
				}
				if dataBiayaTahunan[bulanHarian] == nil {
					dataBiayaTahunan[bulanHarian] = make(map[int]entities.BiayaListrik)
				}

				if existing, ok := dataPenggunaanTahunan[bulanHarian][lantai]; ok {
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " kW"), 64)
					existingValue += daya
					existing.Value = fmt.Sprintf("%.2f kW", existingValue)
					dataPenggunaanTahunan[bulanHarian][lantai] = existing
				} else {
					dataPenggunaanTahunan[bulanHarian][lantai] = entities.PenggunaanListrik{
						Lantai: lantai,
						Value:  fmt.Sprintf("%.2f kW", daya),
					}
				}

				if existing, ok := dataBiayaTahunan[bulanHarian][lantai]; ok {
					existingBiaya, _ := strconv.ParseFloat(strings.TrimPrefix(existing.Biaya, "Rp. "), 64)
					existingBiaya += daya * tarifListrik
					existing.Biaya = fmt.Sprintf("Rp. %.0f", existingBiaya)
					dataBiayaTahunan[bulanHarian][lantai] = existing
				} else {
					dataBiayaTahunan[bulanHarian][lantai] = entities.BiayaListrik{
						Lantai: lantai,
						Biaya:  fmt.Sprintf("Rp. %.0f", daya*tarifListrik),
					}
				}
			}
		}
	}

	convertToSlice := func(data map[int]entities.PenggunaanListrik) []entities.PenggunaanListrik {
		result := make([]entities.PenggunaanListrik, 0, len(data))
		for _, v := range data {
			result = append(result, v)
		}
		return result
	}

	convertBiayaToSlice := func(data map[int]entities.BiayaListrik) []entities.BiayaListrik {
		result := make([]entities.BiayaListrik, 0, len(data))
		for _, v := range data {
			result = append(result, v)
		}
		return result
	}

	response := entities.GetListrikDataResponse{
		TotalWatt:                     fmt.Sprintf("%.2f kW", totalWatt),
		TotalDayaListrikLT1:           fmt.Sprintf("%.2f kW", totalDaya["LT1"]),
		TotalDayaListrikLT2:           fmt.Sprintf("%.2f kW", totalDaya["LT2"]),
		TotalDayaListrikLT3:           fmt.Sprintf("%.2f kW", totalDaya["LT3"]),
		TotalDayaListrikLT4:           fmt.Sprintf("%.2f kW", totalDaya["LT4"]),
		BiayaPemakaianLT1:             biayaLT1,
		BiayaPemakaianLT2:             biayaLT2,
		BiayaPemakaianLT3:             biayaLT3,
		BiayaPemakaianLT4:             biayaLT4,
		DataPenggunaanListrikHarian:   make(map[string][]entities.PenggunaanListrik),
		DataBiayaListrikHarian:        make(map[string][]entities.BiayaListrik),
		DataPenggunaanListrikMingguan: make(map[string][]entities.PenggunaanListrik),
		DataBiayaListrikMingguan:      make(map[string][]entities.BiayaListrik),
		DataPenggunaanListrikTahunan:  make(map[string][]entities.PenggunaanListrik),
		DataBiayaListrikTahunan:       make(map[string][]entities.BiayaListrik),
		CreatedAt:                     createdAt,
		UpdatedAt:                     updatedAt,
	}

	for hari, data := range dataPenggunaanHarian {
		response.DataPenggunaanListrikHarian[hari] = convertToSlice(data)
	}
	for hari, data := range dataBiayaHarian {
		response.DataBiayaListrikHarian[hari] = convertBiayaToSlice(data)
	}

	for minggu, data := range dataPenggunaanMingguan {
		response.DataPenggunaanListrikMingguan[minggu] = convertToSlice(data)
	}
	for minggu, data := range dataBiayaMingguan {
		response.DataBiayaListrikMingguan[minggu] = convertBiayaToSlice(data)
	}

	for bulan, data := range dataPenggunaanTahunan {
		response.DataPenggunaanListrikTahunan[bulan] = convertToSlice(data)
	}
	for bulan, data := range dataBiayaTahunan {
		response.DataBiayaListrikTahunan[bulan] = convertBiayaToSlice(data)
	}

	return response, nil
}

func getHariIndonesia(weekday time.Weekday) string {
	hari := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	return hari[weekday]
}

func getBulanIndonesia(month time.Month) string {
	bulan := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return bulan[month-1]
}
