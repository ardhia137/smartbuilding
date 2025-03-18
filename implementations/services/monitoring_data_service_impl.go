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
	dataTorenRepository      repositories.DataTorenRepository
	settingRepository        repositories.SettingRepository
}

func NewMonitoringDataService(monitorRepo repositories.MonitoringDataRepository,

	dataToren repositories.DataTorenRepository,
	settingRepo repositories.SettingRepository,
) services.MonitoringDataService {
	return &monitoringDataServiceImpl{monitorRepo, dataToren, settingRepo}

}

func (s *monitoringDataServiceImpl) SaveMonitoringData(request entities.CreateMonitoringDataRequest) (entities.MonitoringDataResponse, error) {
	monitoringData := entities.MonitoringData{
		MonitoringName:  request.MonitoringName,
		MonitoringValue: request.MonitoringValue,
		IDSetting:       request.IDSetting,
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

func (s *monitoringDataServiceImpl) GetAirMonitoringData(id int) ([]entities.GetAirDataResponse, error) {
	monitoringData, err := s.monitoringDataRepository.GetAirMonitoringData(id)
	if err != nil {
		return nil, err
	}
	torenData, err := s.dataTorenRepository.FindBySettingID(id)
	if err != nil {
		return nil, err
	}
	monitoringDataHarian, err := s.monitoringDataRepository.GetAirMonitoringDataHarian(id)
	if err != nil {
		return nil, err
	}

	settingRepo, err := s.settingRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	namaGedung := settingRepo.NamaGedung

	var totalAirKeluar, totalAirMasuk float64
	var createdAt, updatedAt time.Time

	dataPenggunaanHarian := make(map[string][]entities.PenggunaanAir)
	dataPenggunaanMingguan := make(map[string]map[string]float64)
	dataPenggunaanTahunan := make(map[string]map[string]float64)
	latestWaterFlowMasuk := make(map[string]float64)
	latestCreatedAtMasuk := make(map[string]time.Time)
	latestWaterFlowKeluar := make(map[string]float64)
	latestCreatedAtKeluar := make(map[string]time.Time)

	kapasitasTorenMap := make(map[string]entities.KapasitasTorenData)
	// Simpan data toren berdasarkan nama
	torenDataMap := make(map[string]entities.DataToren)
	for _, toren := range torenData {
		torenDataMap[toren.MonitoringName] = toren
	}

	for _, data := range monitoringData {
		switch {
		case strings.HasPrefix(data.MonitoringName, "monitoring_air_kapasitas_toren"):
			namaToren := data.MonitoringName
			kapasitas := data.MonitoringValue

			// Ambil kapasitas toren terbaru berdasarkan CreatedAt
			if lastData, exists := kapasitasTorenMap[namaToren]; !exists || data.CreatedAt.After(lastData.CreatedAt) {
				// Hitung volume sensor berdasarkan kapasitas toren dari monitoring
				kapasitasFloat, _ := strconv.ParseFloat(strings.TrimSuffix(kapasitas, " %"), 64)

				// Cek apakah ada data toren yang cocok
				kapasitasTorenFinal := kapasitas
				if toren, found := torenDataMap[namaToren]; found {
					kapasitasTorenFinal = strconv.Itoa(toren.KapasitasToren) // Gunakan kapasitas dari torenData
				}
				kapasitasTorenFloat, _ := strconv.ParseFloat(kapasitasTorenFinal, 64)
				volumeSensor := kapasitasTorenFloat * (kapasitasFloat / 100)
				namaTorenFormatted := strings.TrimPrefix(namaToren, "monitoring_air_")
				namaTorenFormatted = strings.ReplaceAll(namaTorenFormatted, "_", " ")
				kapasitasTorenMap[namaToren] = entities.KapasitasTorenData{
					Nama:           namaTorenFormatted,
					Kapasitas:      kapasitas,
					KapasitasToren: kapasitasTorenFinal, // Gunakan nilai yang telah diperbarui jika cocok
					VolumeSensor:   fmt.Sprintf("%.0f L", volumeSensor),
					CreatedAt:      data.CreatedAt,
				}
			}
		case strings.HasPrefix(data.MonitoringName, "monitoring_air_total_water_flow_air_masuk"):
			volumeMasuk, _ := strconv.ParseFloat(strings.TrimSuffix(data.MonitoringValue, " L"), 64)
			pipa := strings.TrimPrefix(data.MonitoringName, "monitoring_air_total_water_flow_air_masuk")
			if lastTime, exists := latestCreatedAtMasuk[pipa]; !exists || data.CreatedAt.After(lastTime) {
				latestWaterFlowMasuk[pipa] = volumeMasuk
				latestCreatedAtMasuk[pipa] = data.CreatedAt
			}
		default:
			if strings.HasPrefix(data.MonitoringName, "monitoring_air_total_water_flow_air_keluar_") {
				volume, _ := strconv.ParseFloat(strings.TrimSuffix(data.MonitoringValue, " L"), 64)

				// Ambil nama pipa dari monitoring name
				pipa := strings.TrimPrefix(data.MonitoringName, "monitoring_air_total_water_flow_air_keluar_")
				// Simpan nilai terbaru berdasarkan CreatedAt
				if lastTime, exists := latestCreatedAtKeluar[pipa]; !exists || data.CreatedAt.After(lastTime) {
					latestWaterFlowKeluar[pipa] = volume
					latestCreatedAtKeluar[pipa] = data.CreatedAt
				}
			}
		}

		createdAt = data.CreatedAt
		updatedAt = data.UpdatedAt
	}
	// Hitung total air keluar hanya dari data terbaru tiap pipa
	totalAirKeluar = 0
	for _, volume := range latestWaterFlowKeluar {
		totalAirKeluar += volume
	}
	totalAirMasuk = 0
	for _, volume := range latestWaterFlowMasuk {
		totalAirMasuk += volume
	}
	var kapasitasToren []entities.KapasitasTorenData
	for _, toren := range kapasitasTorenMap {
		kapasitasToren = append(kapasitasToren, toren)
	}
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
		NamaGedung:             namaGedung,
		KapasitasToren:         kapasitasToren,
		AirMasuk:               fmt.Sprintf("%.0f L", totalAirMasuk),
		AirKeluar:              fmt.Sprintf("%.0f L", totalAirKeluar),
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

func (s *monitoringDataServiceImpl) GetListrikMonitoringData(id int) (entities.GetListrikDataResponse, error) {
	monitoringData, err := s.monitoringDataRepository.GetListrikMonitoringData(id)
	if err != nil {
		return entities.GetListrikDataResponse{}, err
	}

	monitoringDataHarian, err := s.monitoringDataRepository.GetListrikMonitoringDataHarian(id)
	if err != nil {
		return entities.GetListrikDataResponse{}, err
	}

	setting, err := s.settingRepository.FindByID(id)
	if err != nil {
		return entities.GetListrikDataResponse{}, err
	}
	namaGedung := setting.NamaGedung
	jenisListrik := setting.JenisListrik
	var schadule = float64(setting.Scheduler)
	var tarifListrik = float64(setting.HargaListrik)
	var createdAt, updatedAt time.Time
	var totalWatt float64

	dataPenggunaanHarian := make(map[string]map[string]entities.PenggunaanListrik)
	dataBiayaHarian := make(map[string]map[string]entities.BiayaListrik)

	dataPenggunaanMingguan := make(map[string]map[string]entities.PenggunaanListrik)
	dataBiayaMingguan := make(map[string]map[string]entities.BiayaListrik)

	dataPenggunaanTahunan := make(map[string]map[string]entities.PenggunaanListrik)
	dataBiayaTahunan := make(map[string]map[string]entities.BiayaListrik)

	var totalDaya []entities.TotalDayaListrik
	totalBiaya := []entities.BiayaListrik{} // Map untuk menyimpan total biaya

	totalDayaMap := make(map[string]float64)  // Untuk menyimpan sementara sebelum konversi ke slice
	totalBiayaMap := make(map[string]float64) // Untuk menyimpan sementara sebelum konversi ke slice

	for _, data := range monitoringData {
		arus, _ := strconv.ParseFloat(strings.TrimSuffix(data.MonitoringValue, " A"), 64)
		var kw, kwh float64

		if jenisListrik == "1_phase" {
			kw = 220 * arus * 0.8 / 1000
			kwh = kw * (schadule / 3600.0)
		} else if jenisListrik == "3_phase" {
			kw = 1.732 * 380 * arus * 0.8 / 1000
			kwh = kw * (schadule / 3600.0)
		}

		if data.MonitoringName != "" {
			totalDayaMap[data.MonitoringName] += kw
			totalBiayaMap[data.MonitoringName] += kwh * float64(tarifListrik)
		}

		totalWatt += kw
		createdAt = data.CreatedAt
		updatedAt = data.UpdatedAt
	}

	// Konversi map ke slice
	for key, value := range totalDayaMap {
		totalDaya = append(totalDaya, entities.TotalDayaListrik{
			Nama:  strings.ReplaceAll(strings.TrimPrefix(key, "monitoring_listrik_arus_"), "_", " "),
			Value: fmt.Sprintf("%.2f kW", value),
		})
	}
	for key, value := range totalBiayaMap {
		totalBiaya = append(totalBiaya, entities.BiayaListrik{
			Nama:  strings.ReplaceAll(strings.TrimPrefix(key, "monitoring_listrik_arus_"), "_", " "),
			Biaya: fmt.Sprintf("Rp.%.2f", value),
		})
	}

	now := time.Now()
	year, month, _ := now.Date()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	for _, harian := range monitoringDataHarian {
		if strings.Contains(harian.MonitoringName, "arus_listrik") {
			arus, _ := strconv.ParseFloat(strings.TrimSuffix(harian.MonitoringValue, " A"), 64)
			var kw, kwh float64
			if jenisListrik == "1_phase" {
				kw = 220 * arus * 0.8 / 1000
				kwh = kw / 24
			} else if jenisListrik == "3_phase" {
				kw = 1.732 * 380 * arus * 0.8 / 1000
				kwh = kw / 24
			}

			hari := getHariIndonesia(harian.CreatedAt.Weekday())

			if harian.CreatedAt.After(startOfWeek) && harian.CreatedAt.Before(endOfWeek) {
				if dataPenggunaanHarian[hari] == nil {
					dataPenggunaanHarian[hari] = make(map[string]entities.PenggunaanListrik)
				}
				if dataBiayaHarian[hari] == nil {
					dataBiayaHarian[hari] = make(map[string]entities.BiayaListrik)
				}

				if existing, ok := dataPenggunaanHarian[hari][harian.MonitoringName]; ok {
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " kW"), 64)
					existingValue += kw
					existing.Value = fmt.Sprintf("%.2f kW", existingValue)
					dataPenggunaanHarian[hari][harian.MonitoringName] = existing
				} else {
					dataPenggunaanHarian[hari][harian.MonitoringName] = entities.PenggunaanListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Value: fmt.Sprintf("%.2f kW", kw),
					}
				}

				if existing, ok := dataBiayaHarian[hari][harian.MonitoringName]; ok {
					existingBiaya, _ := strconv.ParseFloat(strings.TrimPrefix(existing.Biaya, "Rp. "), 64)
					existingBiaya += kwh * tarifListrik
					existing.Biaya = fmt.Sprintf("Rp. %.0f", existingBiaya)
					dataBiayaHarian[hari][harian.MonitoringName] = existing
				} else {
					dataBiayaHarian[hari][harian.MonitoringName] = entities.BiayaListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Biaya: fmt.Sprintf("Rp. %.0f", kwh*tarifListrik),
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
					dataPenggunaanMingguan[mingguanKey] = make(map[string]entities.PenggunaanListrik)
				}
				if dataBiayaMingguan[mingguanKey] == nil {
					dataBiayaMingguan[mingguanKey] = make(map[string]entities.BiayaListrik)
				}

				if existing, ok := dataPenggunaanMingguan[mingguanKey][harian.MonitoringName]; ok {
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " kW"), 64)
					existingValue += kw
					existing.Value = fmt.Sprintf("%.2f kW", existingValue)
					dataPenggunaanMingguan[mingguanKey][harian.MonitoringName] = existing
				} else {
					dataPenggunaanMingguan[mingguanKey][harian.MonitoringName] = entities.PenggunaanListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Value: fmt.Sprintf("%.2f kW", kw),
					}
				}

				if existing, ok := dataBiayaMingguan[mingguanKey][harian.MonitoringName]; ok {
					existingBiaya, _ := strconv.ParseFloat(strings.TrimPrefix(existing.Biaya, "Rp. "), 64)
					existingBiaya += kwh * tarifListrik
					existing.Biaya = fmt.Sprintf("Rp. %.0f", existingBiaya)
					dataBiayaMingguan[mingguanKey][harian.MonitoringName] = existing
				} else {
					dataBiayaMingguan[mingguanKey][harian.MonitoringName] = entities.BiayaListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Biaya: fmt.Sprintf("Rp. %.0f", kwh*tarifListrik),
					}
				}
			}

			bulanHarian := getBulanIndonesia(harian.CreatedAt.Month())

			if harian.CreatedAt.Year() == year {
				if dataPenggunaanTahunan[bulanHarian] == nil {
					dataPenggunaanTahunan[bulanHarian] = make(map[string]entities.PenggunaanListrik)
				}
				if dataBiayaTahunan[bulanHarian] == nil {
					dataBiayaTahunan[bulanHarian] = make(map[string]entities.BiayaListrik)
				}

				if existing, ok := dataPenggunaanTahunan[bulanHarian][harian.MonitoringName]; ok {
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " kW"), 64)
					existingValue += kw
					existing.Value = fmt.Sprintf("%.2f kW", existingValue)
					dataPenggunaanTahunan[bulanHarian][harian.MonitoringName] = existing
				} else {
					dataPenggunaanTahunan[bulanHarian][harian.MonitoringName] = entities.PenggunaanListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Value: fmt.Sprintf("%.2f kW", kw),
					}
				}

				if existing, ok := dataBiayaTahunan[bulanHarian][harian.MonitoringName]; ok {
					existingBiaya, _ := strconv.ParseFloat(strings.TrimPrefix(existing.Biaya, "Rp. "), 64)
					existingBiaya += kwh * tarifListrik
					existing.Biaya = fmt.Sprintf("Rp. %.0f", existingBiaya)
					dataBiayaTahunan[bulanHarian][harian.MonitoringName] = existing
				} else {
					dataBiayaTahunan[bulanHarian][harian.MonitoringName] = entities.BiayaListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Biaya: fmt.Sprintf("Rp. %.0f", kwh*tarifListrik),
					}
				}
			}
		}
	}

	convertToSlice := func(data map[string]entities.PenggunaanListrik) []entities.PenggunaanListrik {
		result := make([]entities.PenggunaanListrik, 0, len(data))
		for _, v := range data {
			result = append(result, v)
		}
		return result
	}

	convertBiayaToSlice := func(data map[string]entities.BiayaListrik) []entities.BiayaListrik {
		result := make([]entities.BiayaListrik, 0, len(data))
		for _, v := range data {
			result = append(result, v)
		}
		return result
	}

	response := entities.GetListrikDataResponse{
		NamaGedung:                    namaGedung,
		TotalWatt:                     fmt.Sprintf("%.2f kW", totalWatt),
		TotalDayaListrik:              totalDaya,
		BiayaPemakaian:                totalBiaya,
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
