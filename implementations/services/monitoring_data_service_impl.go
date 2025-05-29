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
	dataPenggunaanBulanan := make(map[string]map[string]float64)
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
	startOfWeek := getStartOfWeek(now)
	endOfWeek := getEndOfWeek(startOfWeek)

	for _, harian := range monitoringDataHarian {
		if strings.HasPrefix(harian.MonitoringName, "monitoring_air_total_water_flow_") {
			pipa := strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_air_total_water_flow_"), "_", " ")
			volume, _ := strconv.ParseFloat(strings.TrimSuffix(harian.MonitoringValue, " L"), 64)
			hari := getHariIndonesia(harian.CreatedAt.Weekday())

			if harian.CreatedAt.After(startOfWeek) && harian.CreatedAt.Before(endOfWeek) {
				dataPenggunaanHarian[hari] = append(dataPenggunaanHarian[hari], entities.PenggunaanAir{
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
				if dataPenggunaanMingguan[mingguanKey] == nil {
					dataPenggunaanMingguan[mingguanKey] = make(map[string]float64)
				}
				dataPenggunaanMingguan[mingguanKey][pipa] += volume
			}

			bulan := getBulanIndonesia(harian.CreatedAt.Month())

			if harian.CreatedAt.Year() == year {
				if dataPenggunaanBulanan[bulan] == nil {
					dataPenggunaanBulanan[bulan] = make(map[string]float64)
				}
				dataPenggunaanBulanan[bulan][pipa] += volume

				// Process yearly data
				tahun := strconv.Itoa(harian.CreatedAt.Year())
				if dataPenggunaanTahunan[tahun] == nil {
					dataPenggunaanTahunan[tahun] = make(map[string]float64)
				}
				dataPenggunaanTahunan[tahun][pipa] += volume
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
		DataPenggunaanBulanan:  make(map[string][]entities.PenggunaanAir),
		DataPenggunaanTahunan:  make(map[string][]entities.PenggunaanAir),
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}

	for minggu, data := range dataPenggunaanMingguan {
		response.DataPenggunaanMingguan[minggu] = convertMingguanTahunan(data)
	}

	for bulan, data := range dataPenggunaanBulanan {
		response.DataPenggunaanBulanan[bulan] = convertMingguanTahunan(data)
	}

	for tahun, data := range dataPenggunaanTahunan {
		response.DataPenggunaanTahunan[tahun] = convertMingguanTahunan(data)
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

	dataPenggunaanHarian := make(map[string]map[string]entities.PenggunaanListrik)
	dataBiayaHarian := make(map[string]map[string]entities.BiayaListrik)

	dataPenggunaanMingguan := make(map[string]map[string]entities.PenggunaanListrik)
	dataBiayaMingguan := make(map[string]map[string]entities.BiayaListrik)

	dataPenggunaanBulanan := make(map[string]map[string]entities.PenggunaanListrik)
	dataBiayaBulanan := make(map[string]map[string]entities.BiayaListrik)

	dataPenggunaanTahunan := make(map[string]map[string]entities.PenggunaanListrik)
	dataBiayaTahunan := make(map[string]map[string]entities.BiayaListrik)

	var totalDaya []entities.TotalDayaListrik
	totalBiaya := []entities.BiayaListrik{} // Map untuk menyimpan total biaya
	now := time.Now()
	//hour := now.Hour()
	//totalDayaMap := make(map[string]float64)    // Untuk menyimpan sementara sebelum konversi ke slice
	//totalBiayaMap := make(map[string]float64)    // Untuk menyimpan sementara sebelum konversi ke slice
	//totalJumlahData := make(map[string]float64)  // Untuk menyimpan sementara sebelum konversi ke slice
	totalArusPerName := make(map[string]float64) // Untuk menyimpan total arus per monitoring name
	jumlahDataPerName := make(map[string]int)    // Untuk menyimpan jumlah data per monitoring name
	for i, data := range monitoringData {
		arus, _ := strconv.ParseFloat(strings.TrimSuffix(data.MonitoringValue, " A"), 64)
		var _ float64

		// Tambahkan ke total arus untuk perhitungan rata-rata
		totalArus += arus
		jumlahData++

		// Tambahkan ke total arus per monitoring name
		if data.MonitoringName != "" {
			totalArusPerName[data.MonitoringName] += arus
			jumlahDataPerName[data.MonitoringName]++
		}
		// Hitung daya berdasarkan jenis listrik
		if jenisListrik == "1_phase" {
			_ = 220 * arus * 0.8 / 1000
		} else if jenisListrik == "3_phase" {
			_ = 1.732 * 380 * arus * 0.8 / 1000
		}

		// Tentukan createdAt (paling awal) dan updatedAt (paling akhir)
		if i == 0 || data.CreatedAt.Before(createdAt) {
			createdAt = data.CreatedAt
		}
		if i == 0 || data.UpdatedAt.After(updatedAt) {
			updatedAt = data.UpdatedAt
		}

		// Kita akan menghitung total daya nanti berdasarkan rata-rata arus per monitoring name
	}

	// Hitung selisih waktu dalam detik

	// Hitung daya (kW) dan energi (kWh) berdasarkan rata-rata arus
	if jumlahData > 0 {
		rataRataArus := totalArus / float64(jumlahData)
		currentHour := float64(now.Hour())
		if currentHour == 0 {
			currentHour = 1 // Hindari pembagian dengan 0
		}

		// Hitung daya dalam kW (tanpa mengalikan dengan jam)
		var dayaKW float64
		if jenisListrik == "1_phase" {
			dayaKW = 220 * rataRataArus * 0.8 / 1000
		} else if jenisListrik == "3_phase" {
			dayaKW = 1.732 * 380 * rataRataArus * 0.8 / 1000
		}

		// Hitung energi dalam kWh (daya dikali jam)
		energiKWh := dayaKW * currentHour

		// Set totalWatt ke nilai daya (kW)
		totalWatt = energiKWh

	}

	// Hindari pembagian dengan 0
	// Biaya akan dihitung berdasarkan energi (kWh) per monitoring name

	// Hitung daya berdasarkan rata-rata arus per monitoring name
	totalWatt = 0 // Reset totalWatt untuk menghitung ulang
	currentHour := float64(now.Hour())
	if currentHour == 0 {
		currentHour = 1 // Hindari pembagian dengan 0
	}

	for key, totalArus := range totalArusPerName {
		jumlahData := jumlahDataPerName[key]
		if jumlahData > 0 {
			// Hitung rata-rata arus untuk monitoring name ini
			rataRataArus := totalArus / float64(jumlahData)

			// Hitung daya dalam kW berdasarkan rata-rata arus
			var dayaKW float64
			if jenisListrik == "1_phase" {
				dayaKW = 220 * rataRataArus * 0.8 / 1000
			} else if jenisListrik == "3_phase" {
				dayaKW = 1.732 * 380 * rataRataArus * 0.8 / 1000
			}

			// Hitung energi dalam kWh (daya dikali jam)
			energiKWh := dayaKW * currentHour

			// Tambahkan ke totalWatt
			totalWatt += energiKWh

			// Tambahkan ke totalDaya
			totalDaya = append(totalDaya, entities.TotalDayaListrik{
				Nama:  strings.ReplaceAll(strings.TrimPrefix(key, "monitoring_listrik_arus_"), "_", " "),
				Value: fmt.Sprintf("%.2f kWh", energiKWh),
			})
		}
	}

	// Hitung biaya berdasarkan energi (kWh) yang telah dihitung sebelumnya
	for _, dayaItem := range totalDaya {
		// Ekstrak nama dan nilai kWh
		nama := dayaItem.Nama
		kwhStr := strings.TrimSuffix(dayaItem.Value, " kWh")
		kwh, _ := strconv.ParseFloat(kwhStr, 64)

		// Hitung biaya
		biaya := kwh * tarifListrik

		// Tambahkan ke totalBiaya
		totalBiaya = append(totalBiaya, entities.BiayaListrik{
			Nama:  nama,
			Biaya: fmt.Sprintf("Rp. %.0f", biaya),
		})
	}

	year, month, _ := now.Date()
	startOfWeek := getStartOfWeek(now)
	endOfWeek := getEndOfWeek(startOfWeek)

	for _, harian := range monitoringDataHarian {
		if strings.Contains(harian.MonitoringName, "arus_listrik") {
			arus, _ := strconv.ParseFloat(strings.TrimSuffix(harian.MonitoringValue, " A"), 64)
			var kw, kwh float64

			// Hitung daya (kW) dan energi (kWh) berdasarkan arus
			if jenisListrik == "1_phase" {
				// Hitung rata-rata arus per hari
				jumlahSampel := 86400 / schadule
				Rarus := arus / jumlahSampel
				// Hitung daya dalam kW
				kw = 220 * Rarus * 0.8 / 1000
				// Hitung energi dalam kWh (daya dikali 24 jam)
				kwh = kw * 24
			} else if jenisListrik == "3_phase" {
				// Hitung rata-rata arus per hari
				jumlahSampel := 86400 / schadule
				Rarus := arus / jumlahSampel
				// Hitung daya dalam kW
				kw = 1.732 * 380 * Rarus * 0.8 / 1000
				// Hitung energi dalam kWh (daya dikali 24 jam)
				kwh = kw * 24
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
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " Kwh"), 64)
					existingValue += kwh
					existing.Value = fmt.Sprintf("%.2f Kwh", existingValue)
					dataPenggunaanHarian[hari][harian.MonitoringName] = existing
				} else {
					dataPenggunaanHarian[hari][harian.MonitoringName] = entities.PenggunaanListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Value: fmt.Sprintf("%.2f Kwh", kwh),
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
				if dataPenggunaanMingguan[mingguanKey] == nil {
					dataPenggunaanMingguan[mingguanKey] = make(map[string]entities.PenggunaanListrik)
				}
				if dataBiayaMingguan[mingguanKey] == nil {
					dataBiayaMingguan[mingguanKey] = make(map[string]entities.BiayaListrik)
				}

				if existing, ok := dataPenggunaanMingguan[mingguanKey][harian.MonitoringName]; ok {
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " Kwh"), 64)
					// Akumulasikan energi dalam kWh
					existingValue += kwh
					existing.Value = fmt.Sprintf("%.2f Kwh", existingValue)
					dataPenggunaanMingguan[mingguanKey][harian.MonitoringName] = existing
				} else {
					dataPenggunaanMingguan[mingguanKey][harian.MonitoringName] = entities.PenggunaanListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Value: fmt.Sprintf("%.2f Kwh", kwh),
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
				if dataPenggunaanBulanan[bulanHarian] == nil {
					dataPenggunaanBulanan[bulanHarian] = make(map[string]entities.PenggunaanListrik)
				}
				if dataBiayaBulanan[bulanHarian] == nil {
					dataBiayaBulanan[bulanHarian] = make(map[string]entities.BiayaListrik)
				}

				if existing, ok := dataPenggunaanBulanan[bulanHarian][harian.MonitoringName]; ok {
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " Kwh"), 64)
					// Akumulasikan energi dalam kWh
					existingValue += kwh
					existing.Value = fmt.Sprintf("%.2f Kwh", existingValue)
					dataPenggunaanBulanan[bulanHarian][harian.MonitoringName] = existing
				} else {
					dataPenggunaanBulanan[bulanHarian][harian.MonitoringName] = entities.PenggunaanListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Value: fmt.Sprintf("%.2f Kwh", kwh),
					}
				}

				if existing, ok := dataBiayaBulanan[bulanHarian][harian.MonitoringName]; ok {
					existingBiaya, _ := strconv.ParseFloat(strings.TrimPrefix(existing.Biaya, "Rp. "), 64)
					existingBiaya += kwh * tarifListrik
					existing.Biaya = fmt.Sprintf("Rp. %.0f", existingBiaya)
					dataBiayaBulanan[bulanHarian][harian.MonitoringName] = existing
				} else {
					dataBiayaBulanan[bulanHarian][harian.MonitoringName] = entities.BiayaListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Biaya: fmt.Sprintf("Rp. %.0f", kwh*tarifListrik),
					}
				}

				// Process yearly data
				tahun := strconv.Itoa(harian.CreatedAt.Year())
				if dataPenggunaanTahunan[tahun] == nil {
					dataPenggunaanTahunan[tahun] = make(map[string]entities.PenggunaanListrik)
				}
				if dataBiayaTahunan[tahun] == nil {
					dataBiayaTahunan[tahun] = make(map[string]entities.BiayaListrik)
				}

				if existing, ok := dataPenggunaanTahunan[tahun][harian.MonitoringName]; ok {
					existingValue, _ := strconv.ParseFloat(strings.TrimSuffix(existing.Value, " Kwh"), 64)
					// Akumulasikan energi dalam kWh
					existingValue += kwh
					existing.Value = fmt.Sprintf("%.2f Kwh", existingValue)
					dataPenggunaanTahunan[tahun][harian.MonitoringName] = existing
				} else {
					dataPenggunaanTahunan[tahun][harian.MonitoringName] = entities.PenggunaanListrik{
						Nama:  strings.ReplaceAll(strings.TrimPrefix(harian.MonitoringName, "monitoring_listrik_arus_"), "_", " "),
						Value: fmt.Sprintf("%.2f Kwh", kwh),
					}
				}

				if existing, ok := dataBiayaTahunan[tahun][harian.MonitoringName]; ok {
					existingBiaya, _ := strconv.ParseFloat(strings.TrimPrefix(existing.Biaya, "Rp. "), 64)
					existingBiaya += kwh * tarifListrik
					existing.Biaya = fmt.Sprintf("Rp. %.0f", existingBiaya)
					dataBiayaTahunan[tahun][harian.MonitoringName] = existing
				} else {
					dataBiayaTahunan[tahun][harian.MonitoringName] = entities.BiayaListrik{
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
		TotalWatt:                     fmt.Sprintf("%.2f kWh", totalWatt),
		TotalDayaListrik:              totalDaya,
		BiayaPemakaian:                totalBiaya,
		DataPenggunaanListrikHarian:   make(map[string][]entities.PenggunaanListrik),
		DataBiayaListrikHarian:        make(map[string][]entities.BiayaListrik),
		DataPenggunaanListrikMingguan: make(map[string][]entities.PenggunaanListrik),
		DataBiayaListrikMingguan:      make(map[string][]entities.BiayaListrik),
		DataPenggunaanListrikBulanan:  make(map[string][]entities.PenggunaanListrik),
		DataBiayaListrikBulanan:       make(map[string][]entities.BiayaListrik),
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

	for bulan, data := range dataPenggunaanBulanan {
		response.DataPenggunaanListrikBulanan[bulan] = convertToSlice(data)
	}
	for bulan, data := range dataBiayaBulanan {
		response.DataBiayaListrikBulanan[bulan] = convertBiayaToSlice(data)
	}

	for tahun, data := range dataPenggunaanTahunan {
		response.DataPenggunaanListrikTahunan[tahun] = convertToSlice(data)
	}
	for tahun, data := range dataBiayaTahunan {
		response.DataBiayaListrikTahunan[tahun] = convertBiayaToSlice(data)
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
		weekday = 7 // Ubah Minggu jadi 7
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
