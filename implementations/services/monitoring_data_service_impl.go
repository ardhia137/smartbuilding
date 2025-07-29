package services

import (
	"fmt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"sort"
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
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	createdData, err := s.monitoringDataRepository.SaveMonitoringData(monitoringData)
	if err != nil {
		return entities.MonitoringDataResponse{}, err
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	response := entities.MonitoringDataResponse{
		ID:              createdData.ID,
		MonitoringName:  createdData.MonitoringName,
		MonitoringValue: createdData.MonitoringValue,
		CreatedAt:       createdData.CreatedAt,
		UpdatedAt:       createdData.UpdatedAt,
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	return response, nil
}

func (s *monitoringDataServiceImpl) GetAirMonitoringData(id int) ([]entities.GetAirDataResponse, error) {
	monitoringData, err := s.monitoringDataRepository.GetAirMonitoringData(id)
	if err != nil {
		return nil, err
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}
	torenData, err := s.dataTorenRepository.FindBySettingID(id)
	if err != nil {
		return nil, err
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}
	monitoringDataHarian, err := s.monitoringDataRepository.GetAirMonitoringDataHarian(id)
	if err != nil {
		return nil, err
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	settingRepo, err := s.settingRepository.FindByID(id)
	if err != nil {
		return nil, err
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	namaGedung := settingRepo.NamaGedung

	var totalAirKeluar, totalAirMasuk float64
	var createdAt, updatedAt time.Time

	dataPenggunaanHarian := make(map[string][]entities.PenggunaanAir)
	dataPenggunaanMingguan := make(map[string][]entities.PenggunaanAir)
	dataPenggunaanBulanan := make(map[string]map[string]float64)
	dataPenggunaanTahunan := make(map[string]map[string]float64)

	// --- Ambil data per jam dari monitoringData (hanya hari ini, hanya jam yang ada datanya) ---
	today := time.Now().Format("2006-01-02")
	// jam -> pipa -> volume
	jamPipaVolume := make(map[string]map[string]string)
	for _, data := range monitoringData {
		if strings.HasPrefix(data.MonitoringName, "monitoring_air_total_water_flow_") && data.CreatedAt.Format("2006-01-02") == today {
			pipa := strings.ReplaceAll(strings.TrimPrefix(data.MonitoringName, "monitoring_air_total_water_flow_"), "_", " ")
			volume := fmt.Sprintf("%.0f L", func() float64 {
				v, _ := strconv.ParseFloat(strings.TrimSuffix(data.MonitoringValue, " L"), 64)
				return v
			}())
			hour := data.CreatedAt.Hour()
			jamStr := fmt.Sprintf("%02d:00", hour)
			if jamPipaVolume[jamStr] == nil {
				jamPipaVolume[jamStr] = make(map[string]string)
			}
			jamPipaVolume[jamStr][pipa] = volume
		}
	}
	// Build DataPenggunaanHarian hanya untuk jam yang ada datanya
	for jamStr, pipaMap := range jamPipaVolume {
		list := make([]entities.PenggunaanAir, 0, len(pipaMap))
		for pipa, volume := range pipaMap {
			list = append(list, entities.PenggunaanAir{
				Pipa:   pipa,
				Volume: volume,
				// Hour bisa diambil dari jamStr jika dibutuhkan
			})
		}
		// Sort ascending berdasarkan nama pipa
		sort.Slice(list, func(i, j int) bool {
			return list[i].Pipa < list[j].Pipa
		})
		dataPenggunaanHarian[jamStr] = list
	}
	latestWaterFlowMasuk := make(map[string]float64)
	latestCreatedAtMasuk := make(map[string]time.Time)
	latestWaterFlowKeluar := make(map[string]float64)
	latestCreatedAtKeluar := make(map[string]time.Time)

	kapasitasTorenMap := make(map[string]entities.KapasitasTorenData)
	// Simpan data toren berdasarkan nama
	torenDataMap := make(map[string]entities.DataToren)
	for _, toren := range torenData {
		torenDataMap[toren.MonitoringName] = toren
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
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
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	// Hitung total air keluar hanya dari data terbaru tiap pipa
	totalAirKeluar = 0
	for _, volume := range latestWaterFlowKeluar {
		totalAirKeluar += volume
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}
	totalAirMasuk = 0
	for _, volume := range latestWaterFlowMasuk {
		totalAirMasuk += volume
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}
	var kapasitasToren []entities.KapasitasTorenData
	for _, toren := range kapasitasTorenMap {
		kapasitasToren = append(kapasitasToren, toren)
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	now := time.Now()
	//hour := now.Hour()
	//totalDayaMap := make(map[string]float64)    // Untuk menyimpan sementara sebelum konversi ke slice
	//totalBiayaMap := make(map[string]float64)    // Untuk menyimpan sementara sebelum konversi ke slice
	//totalJumlahData := make(map[string]float64)  // Untuk menyimpan sementara sebelum konversi ke slice
	year, month, _ := now.Date()
	startOfWeek := getStartOfWeek(now)
	endOfWeek := getEndOfWeek(startOfWeek)

	for _, harian := range monitoringDataHarian {
		if strings.HasPrefix(harian.MonitoringName, "monitoring_air_total_water_flow_") {
			pipa := strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_air_total_water_flow_"), "_", " ")
			volume, _ := strconv.ParseFloat(strings.TrimSuffix(harian.MonitoringValue, " L"), 64)
			hari := getHariIndonesia(harian.CreatedAt.Weekday())

			if harian.CreatedAt.After(startOfWeek) && harian.CreatedAt.Before(endOfWeek) {
				dataPenggunaanMingguan[hari] = append(dataPenggunaanMingguan[hari], entities.PenggunaanAir{
					Pipa:   pipa,
					Volume: fmt.Sprintf("%.0f L", volume),
				})
			}

			// Hitung minggu dalam bulan dengan mempertimbangkan bahwa tanggal 1 tidak selalu hari Senin
			// Dapatkan tanggal pertama dari bulan saat ini
			firstOfMonth := time.Date(harian.CreatedAt.Year(), harian.CreatedAt.Month(), 1, 0, 0, 0, 0, harian.CreatedAt.Location())

			// Hitung offset hari dalam minggu (0 = Minggu, 1 = Senin, ..., 6 = Sabtu)
			// Menggunakan standar ISO: Senin = 1, Minggu = 0/7
			firstDayOffset := int(firstOfMonth.Weekday())
			if firstDayOffset == 0 { // Jika Minggu, set ke 7 untuk perhitungan yang lebih mudah
				firstDayOffset = 7
			}

			// Hitung hari ke berapa dalam bulan
			dayOfMonth := harian.CreatedAt.Day()

			// Hitung hari ke berapa dalam minggu pertama (dengan offset)
			adjustedDay := dayOfMonth + firstDayOffset - 1

			// Hitung minggu (1-indexed)
			minggu := (adjustedDay-1)/7 + 1

			// Untuk memastikan data tanggal 12-18 masuk ke minggu 3
			if dayOfMonth >= 12 && dayOfMonth <= 18 {
				minggu = 3
			}
			mingguanKey := fmt.Sprintf("Minggu %d", minggu)

			if harian.CreatedAt.Year() == year && harian.CreatedAt.Month() == month {
				if dataPenggunaanBulanan[mingguanKey] == nil {
					dataPenggunaanBulanan[mingguanKey] = make(map[string]float64)
				}
				dataPenggunaanBulanan[mingguanKey][pipa] += volume
			}

			bulan := getBulanIndonesia(harian.CreatedAt.Month())

			if harian.CreatedAt.Year() == year {
				if dataPenggunaanTahunan[bulan] == nil {
					dataPenggunaanTahunan[bulan] = make(map[string]float64)
				}
				dataPenggunaanTahunan[bulan][pipa] += volume

			}
		}
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	convertMingguanTahunan := func(data map[string]float64) []entities.PenggunaanAir {
		result := make([]entities.PenggunaanAir, 0, len(data))
		for pipa, volume := range data {
			result = append(result, entities.PenggunaanAir{
				Pipa:   pipa,
				Volume: fmt.Sprintf("%.0f L", volume),
			})
		}
		// Sort ascending berdasarkan nama pipa
		sort.Slice(result, func(i, j int) bool {
			return result[i].Pipa < result[j].Pipa
		})
		return result
		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	response := entities.GetAirDataResponse{
		NamaGedung:             namaGedung,
		KapasitasToren:         kapasitasToren,
		AirMasuk:               fmt.Sprintf("%.0f L", totalAirMasuk),
		AirKeluar:              fmt.Sprintf("%.0f L", totalAirKeluar),
		DataPenggunaanHarian:   dataPenggunaanHarian,
		DataPenggunaanMingguan: make(map[string][]entities.PenggunaanAir),
		DataPenggunaanBulanan:  make(map[string][]entities.PenggunaanAir),
		DataPenggunaanTahunan:  make(map[string][]entities.PenggunaanAir),
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}

	for minggu, data := range dataPenggunaanMingguan {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Pipa < data[j].Pipa
		})
		response.DataPenggunaanMingguan[minggu] = data
	}

	for bulan, data := range dataPenggunaanBulanan {
		sortedData := convertMingguanTahunan(data)
		sort.Slice(sortedData, func(i, j int) bool {
			return sortedData[i].Pipa < sortedData[j].Pipa
		})
		response.DataPenggunaanBulanan[bulan] = sortedData
	}

	for tahun, data := range dataPenggunaanTahunan {
		sortedData := convertMingguanTahunan(data)
		sort.Slice(sortedData, func(i, j int) bool {
			return sortedData[i].Pipa < sortedData[j].Pipa
		})
		response.DataPenggunaanTahunan[tahun] = sortedData
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
	var totalArus float64
	var jumlahData int

	dataPenggunaanHarian := make(map[string][]entities.PenggunaanListrik) // Per jam
	dataBiayaHarian := make(map[string][]entities.BiayaListrik)           // Per jam

	dataPenggunaanMingguan := make(map[string][]entities.PenggunaanListrik) // Per hari
	dataBiayaMingguan := make(map[string][]entities.BiayaListrik)           // Per hari

	dataPenggunaanBulanan := make(map[string]map[string]float64) // Per minggu
	dataBiayaBulanan := make(map[string]map[string]float64)      // Per minggu

	dataPenggunaanTahunan := make(map[string]map[string]float64) // Per bulan
	dataBiayaTahunan := make(map[string]map[string]float64)      // Per bulan

	var totalDaya []entities.TotalDayaListrik
	totalBiaya := []entities.BiayaListrik{}
	now := time.Now()

	totalArusPerName := make(map[string]float64)
	jumlahDataPerName := make(map[string]int)

	today := time.Now().Format("2006-01-02")
	jamArusValue := make(map[string]map[string]float64)

	for i, data := range monitoringData {
		arus, _ := strconv.ParseFloat(strings.TrimSuffix(data.MonitoringValue, " A"), 64)

		totalArus += arus
		jumlahData++

		if data.MonitoringName != "" {
			totalArusPerName[data.MonitoringName] += arus
			jumlahDataPerName[data.MonitoringName]++
		}

		if strings.Contains(data.MonitoringName, "arus_listrik") && data.CreatedAt.Format("2006-01-02") == today {
			hour := data.CreatedAt.Hour()
			jamStr := fmt.Sprintf("%02d:00", hour)

			if jamArusValue[jamStr] == nil {
				jamArusValue[jamStr] = make(map[string]float64)
			}
			jamArusValue[jamStr][data.MonitoringName] = arus
		}

		if i == 0 || data.CreatedAt.Before(createdAt) {
			createdAt = data.CreatedAt
		}
		if i == 0 || data.UpdatedAt.After(updatedAt) {
			updatedAt = data.UpdatedAt
		}
	}

	for jamStr, arusMap := range jamArusValue {
		penggunaanList := make([]entities.PenggunaanListrik, 0, len(arusMap))
		biayaList := make([]entities.BiayaListrik, 0, len(arusMap))

		for monitoringName, arus := range arusMap {
			var kwh float64
			if jenisListrik == "1_phase" {
				kw := 220 * arus * 0.8 / 1000
				kwh = kw * 1 // 1 jam
			} else if jenisListrik == "3_phase" {
				kw := 1.732 * 380 * arus * 0.8 / 1000
				kwh = kw * 1 // 1 jam

			}

			nama := strings.ReplaceAll(strings.TrimPrefix(monitoringName, "monitoring_listrik_arus_"), "_", " ")
			biaya := kwh * tarifListrik

			penggunaanList = append(penggunaanList, entities.PenggunaanListrik{
				Nama:  nama,
				Value: fmt.Sprintf("%.2f kW", kwh),
			})

			biayaList = append(biayaList, entities.BiayaListrik{
				Nama:  nama,
				Biaya: fmt.Sprintf("Rp. %.0f", biaya),
			})
		}

		sort.Slice(penggunaanList, func(i, j int) bool {
			return penggunaanList[i].Nama < penggunaanList[j].Nama
		})
		sort.Slice(biayaList, func(i, j int) bool {
			return biayaList[i].Nama < biayaList[j].Nama
		})

		dataPenggunaanHarian[jamStr] = penggunaanList
		dataBiayaHarian[jamStr] = biayaList
	}

	if jumlahData > 0 {
		rataRataArus := totalArus / float64(jumlahData)
		currentHour := float64(now.Hour())
		if currentHour == 0 {
			currentHour = 1
		}

		var dayaKW float64
		if jenisListrik == "1_phase" {
			dayaKW = 220 * rataRataArus * 0.8 / 1000
		} else if jenisListrik == "3_phase" {
			dayaKW = 1.732 * 380 * rataRataArus * 0.8 / 1000
		}
		energiKWh := dayaKW * float64(len(dataPenggunaanHarian))
		totalWatt = energiKWh
	}

	totalWatt = 0
	currentHour := float64(now.Hour())
	if currentHour == 0 {
		currentHour = 1
	}

	for key, totalArus := range totalArusPerName {
		jumlahData := jumlahDataPerName[key]
		if jumlahData > 0 {
			rataRataArus := totalArus / float64(jumlahData)

			var dayaKW float64
			if jenisListrik == "1_phase" {
				dayaKW = 220 * rataRataArus * 0.8 / 1000
			} else if jenisListrik == "3_phase" {
				dayaKW = 1.732 * 380 * rataRataArus * 0.8 / 1000
			}

			energiKWh := dayaKW * float64(len(dataPenggunaanHarian))

			totalWatt += energiKWh

			totalDaya = append(totalDaya, entities.TotalDayaListrik{
				Nama:  strings.ReplaceAll(strings.TrimPrefix(key, "monitoring_listrik_arus_"), "_", " "),
				Value: fmt.Sprintf("%.2f kW", energiKWh),
			})
		}
	}

	for _, dayaItem := range totalDaya {
		nama := dayaItem.Nama
		kwhStr := strings.TrimSuffix(dayaItem.Value, " kW")
		kwh, _ := strconv.ParseFloat(kwhStr, 64)
		biaya := kwh * tarifListrik

		totalBiaya = append(totalBiaya, entities.BiayaListrik{
			Nama:  nama,
			Biaya: fmt.Sprintf("Rp. %.0f", biaya),
		})
	}

	year, month, _ := now.Date()
	startOfWeek := getStartOfWeek(now)
	endOfWeek := getEndOfWeek(startOfWeek)

	// Proses data harian untuk mingguan, bulanan, dan tahunan
	for _, harian := range monitoringDataHarian {
		if strings.Contains(harian.MonitoringName, "arus_listrik") {
			arus, _ := strconv.ParseFloat(strings.TrimSuffix(harian.MonitoringValue, " A"), 64)
			var kwh float64

			if jenisListrik == "1_phase" {
				jumlahSampel := 86400 / schadule
				Rarus := arus / jumlahSampel
				kw := 220 * Rarus * 0.8 / 1000
				kwh = kw * 24

			} else if jenisListrik == "3_phase" {
				jumlahSampel := 86400 / schadule
				Rarus := arus / jumlahSampel
				kw := 1.732 * 380 * Rarus * 0.8 / 1000
				kwh = kw * 24
			}

			nama := strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " ")
			biaya := kwh * tarifListrik

			hari := getHariIndonesia(harian.CreatedAt.Weekday())
			if harian.CreatedAt.After(startOfWeek) && harian.CreatedAt.Before(endOfWeek) {
				dataPenggunaanMingguan[hari] = append(dataPenggunaanMingguan[hari], entities.PenggunaanListrik{
					Nama:  nama,
					Value: fmt.Sprintf("%.2f kW", kwh),
				})

				dataBiayaMingguan[hari] = append(dataBiayaMingguan[hari], entities.BiayaListrik{
					Nama:  nama,
					Biaya: fmt.Sprintf("Rp. %.0f", biaya),
				})
			}

			firstOfMonth := time.Date(harian.CreatedAt.Year(), harian.CreatedAt.Month(), 1, 0, 0, 0, 0, harian.CreatedAt.Location())
			firstDayOffset := int(firstOfMonth.Weekday())
			if firstDayOffset == 0 {
				firstDayOffset = 7
			}

			dayOfMonth := harian.CreatedAt.Day()
			adjustedDay := dayOfMonth + firstDayOffset - 1
			minggu := (adjustedDay-1)/7 + 1

			if dayOfMonth >= 12 && dayOfMonth <= 18 {
				minggu = 3
			}
			mingguanKey := fmt.Sprintf("Minggu %d", minggu)

			if harian.CreatedAt.Year() == year && harian.CreatedAt.Month() == month {
				if dataPenggunaanBulanan[mingguanKey] == nil {
					dataPenggunaanBulanan[mingguanKey] = make(map[string]float64)
				}
				if dataBiayaBulanan[mingguanKey] == nil {
					dataBiayaBulanan[mingguanKey] = make(map[string]float64)
				}
				dataPenggunaanBulanan[mingguanKey][nama] += kwh
				dataBiayaBulanan[mingguanKey][nama] += biaya
			}

			bulanHarian := getBulanIndonesia(harian.CreatedAt.Month())
			if harian.CreatedAt.Year() == year {
				if dataPenggunaanTahunan[bulanHarian] == nil {
					dataPenggunaanTahunan[bulanHarian] = make(map[string]float64)
				}
				if dataBiayaTahunan[bulanHarian] == nil {
					dataBiayaTahunan[bulanHarian] = make(map[string]float64)
				}
				dataPenggunaanTahunan[bulanHarian][nama] += kwh
				dataBiayaTahunan[bulanHarian][nama] += biaya
			}
		}
	}

	convertToSliceFromMap := func(data map[string]float64, unit string) []entities.PenggunaanListrik {
		result := make([]entities.PenggunaanListrik, 0, len(data))
		for nama, value := range data {
			result = append(result, entities.PenggunaanListrik{
				Nama:  nama,
				Value: fmt.Sprintf("%.2f %s", value, unit),
			})
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].Nama < result[j].Nama
		})
		return result
	}

	convertBiayaToSliceFromMap := func(data map[string]float64) []entities.BiayaListrik {
		result := make([]entities.BiayaListrik, 0, len(data))
		for nama, biaya := range data {
			result = append(result, entities.BiayaListrik{
				Nama:  nama,
				Biaya: fmt.Sprintf("Rp. %.0f", biaya),
			})
		}
		// Sort ascending berdasarkan nama monitoring
		sort.Slice(result, func(i, j int) bool {
			return result[i].Nama < result[j].Nama
		})
		return result
	}

	response := entities.GetListrikDataResponse{
		NamaGedung:                    namaGedung,
		TotalWatt:                     fmt.Sprintf("%.2f kW", totalWatt),
		TotalDayaListrik:              totalDaya,
		BiayaPemakaian:                totalBiaya,
		DataPenggunaanListrikHarian:   dataPenggunaanHarian,   // Per jam
		DataBiayaListrikHarian:        dataBiayaHarian,        // Per jam
		DataPenggunaanListrikMingguan: dataPenggunaanMingguan, // Per hari
		DataBiayaListrikMingguan:      dataBiayaMingguan,      // Per hari
		DataPenggunaanListrikBulanan:  make(map[string][]entities.PenggunaanListrik),
		DataBiayaListrikBulanan:       make(map[string][]entities.BiayaListrik),
		DataPenggunaanListrikTahunan:  make(map[string][]entities.PenggunaanListrik),
		DataBiayaListrikTahunan:       make(map[string][]entities.BiayaListrik),
		CreatedAt:                     createdAt,
		UpdatedAt:                     updatedAt,
	}

	// Sort weekly data ascending berdasarkan nama monitoring
	for hari, data := range dataPenggunaanMingguan {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Nama < data[j].Nama
		})
		response.DataPenggunaanListrikMingguan[hari] = data
	}

	for hari, data := range dataBiayaMingguan {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Nama < data[j].Nama
		})
		response.DataBiayaListrikMingguan[hari] = data
	}

	// Sort total daya dan biaya ascending berdasarkan nama
	sort.Slice(totalDaya, func(i, j int) bool {
		return totalDaya[i].Nama < totalDaya[j].Nama
	})
	sort.Slice(totalBiaya, func(i, j int) bool {
		return totalBiaya[i].Nama < totalBiaya[j].Nama
	})

	// Konversi data bulanan (per minggu)
	for minggu, data := range dataPenggunaanBulanan {
		// Sort ascending berdasarkan nama monitoring
		sortedData := convertToSliceFromMap(data, "kW")
		sort.Slice(sortedData, func(i, j int) bool {
			return sortedData[i].Nama < sortedData[j].Nama
		})
		response.DataPenggunaanListrikBulanan[minggu] = sortedData
	}
	for minggu, data := range dataBiayaBulanan {
		// Sort ascending berdasarkan nama monitoring
		sortedData := convertBiayaToSliceFromMap(data)
		sort.Slice(sortedData, func(i, j int) bool {
			return sortedData[i].Nama < sortedData[j].Nama
		})
		response.DataBiayaListrikBulanan[minggu] = sortedData
	}

	// Konversi data tahunan (per bulan)
	for bulan, data := range dataPenggunaanTahunan {
		// Sort ascending berdasarkan nama monitoring
		sortedData := convertToSliceFromMap(data, "kW")
		sort.Slice(sortedData, func(i, j int) bool {
			return sortedData[i].Nama < sortedData[j].Nama
		})
		response.DataPenggunaanListrikTahunan[bulan] = sortedData
	}
	for bulan, data := range dataBiayaTahunan {
		// Sort ascending berdasarkan nama monitoring
		sortedData := convertBiayaToSliceFromMap(data)
		sort.Slice(sortedData, func(i, j int) bool {
			return sortedData[i].Nama < sortedData[j].Nama
		})
		response.DataBiayaListrikTahunan[bulan] = sortedData
	}

	return response, nil
}

func getHariIndonesia(weekday time.Weekday) string {
	hari := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	return hari[weekday]
}

func getStartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	// Mulai dari jam 00:00
	return time.Date(t.Year(), t.Month(), t.Day()-(weekday-1), 0, 0, 0, 0, t.Location())
}

func getEndOfWeek(start time.Time) time.Time {
	// Akhir minggu jam 23:59:59
	return time.Date(start.Year(), start.Month(), start.Day()+6, 23, 59, 59, 0, start.Location())
}

func getBulanIndonesia(month time.Month) string {
	bulan := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return bulan[month-1]
}
