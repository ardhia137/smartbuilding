package utils

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/usecases"
	"strconv"

	//"strconv"
	"time"
)

var (
	c              *cron.Cron
	cronJobIDs     map[int]cron.EntryID // Menyimpan ID cron job untuk setiap setting
	lastSchedulers map[int]int          // Menyimpan last scheduler untuk setiap setting
)

func init() {
	cronJobIDs = make(map[int]cron.EntryID)
	lastSchedulers = make(map[int]int)
}

func StartMonitoringDataJob(useCase usecases.MonitoringDataUseCase, settingUseCase usecases.SettingUseCase, monitoringDataRepo repositories.MonitoringDataRepository, settingRepo repositories.SettingRepository) {
	c = cron.New()

	// Inisialisasi cron jobs untuk semua settings
	settings, err := settingUseCase.GetAllSetting()
	if err != nil {
		fmt.Println("Error fetching settings:", err)
		return
	}

	for _, setting := range settings {
		cronExpression := fmt.Sprintf("@every %ds", setting.Scheduler)
		cronJobID, err := c.AddFunc(cronExpression, func() {
			runJob(useCase, entities.Setting(setting))
		})
		if err != nil {
			fmt.Printf("Error adding cron job for setting ID %d: %v\n", setting.ID, err)
			continue
		}

		cronJobIDs[setting.ID] = cronJobID
		lastSchedulers[setting.ID] = setting.Scheduler
	}

	// Jadwalkan rekap harian
	_, err = c.AddFunc("59 22 * * *", func() {
		rekapHarian(monitoringDataRepo)
	})
	if err != nil {
		fmt.Println("Error scheduling daily recap job:", err)
		return
	}

	// Mulai monitor perubahan scheduler
	go monitorSchedulerChanges(useCase, settingUseCase)

	c.Start()
	select {}
}

//
//func rekapHarian(monitoringDataRepo repositories.MonitoringDataRepository, settingRepo repositories.SettingRepository) {
//	fmt.Println("Starting daily recap at:", time.Now().Format("2006-01-02 15:04:05"))
//
//	// Ambil semua settings untuk rekap harian
//	settings, err := settingRepo.FindAll()
//	if err != nil {
//		fmt.Println("Error fetching settings:", err)
//		return
//	}
//
//	for _, setting := range settings {
//		// Ambil data monitoring berdasarkan SettingID
//		monitoringData, err := monitoringDataRepo.FindBySettingID(setting.ID)
//		if err != nil {
//			fmt.Printf("Error fetching monitoring data for setting ID %d: %v\n", setting.ID, err)
//			continue
//		}
//
//		totalMap := make(map[string]float64)
//		countMap := make(map[string]int)
//
//		for _, data := range monitoringData {
//			cleanedValue := removeUnits(data.MonitoringValue)
//
//			value, err := strconv.ParseFloat(cleanedValue, 64)
//			if err != nil {
//				fmt.Printf("Error parsing value for %s: %v\n", data.MonitoringName, err)
//				continue
//			}
//
//			totalMap[data.MonitoringName] += value
//			countMap[data.MonitoringName]++
//		}
//
//		for monitoringName, total := range totalMap {
//			count := countMap[monitoringName]
//			average := total / float64(count)
//
//			harianData := entities.MonitoringDataHarian{
//				MonitoringName:  monitoringName,
//				MonitoringValue: fmt.Sprintf("%.2f", average),
//				IDSetting:       uint(setting.ID),
//				CreatedAt:       time.Now(),
//				UpdatedAt:       time.Now(),
//			}
//
//			_, err := monitoringDataRepo.SaveHarianData(entities.MonitoringData(harianData))
//			if err != nil {
//				fmt.Printf("Error saving harian data for %s (Setting ID %d): %v\n", monitoringName, setting.ID, err)
//			}
//		}
//	}
//
//	fmt.Println("Daily recap completed at:", time.Now().Format("2006-01-02 15:04:05"))
//}

func rekapHarian(monitoringDataRepo repositories.MonitoringDataRepository) {
	fmt.Println("Starting daily recap at:", time.Now().Format("2006-01-02 15:04:05"))

	monitoringData, err := monitoringDataRepo.FindAll()
	if err != nil {
		fmt.Println("Error fetching monitoring data:", err)
		return
	}

	totalMap := make(map[string]float64)
	countMap := make(map[string]int)
	idSettingMap := make(map[string]int)
	for _, data := range monitoringData {
		cleanedValue := removeUnits(data.MonitoringValue)

		idSettingMap[data.MonitoringName] = int(data.IDSetting)
		value, err := strconv.ParseFloat(cleanedValue, 64)
		if err != nil {
			fmt.Printf("Error parsing value for %s: %v\n", data.MonitoringName, err)
			continue
		}

		totalMap[data.MonitoringName] += value
		countMap[data.MonitoringName]++
	}

	for monitoringName, total := range totalMap {
		count := countMap[monitoringName]
		average := total / float64(count)

		harianData := entities.MonitoringData{
			MonitoringName:  monitoringName,
			MonitoringValue: fmt.Sprintf("%.2f", average),
			IDSetting:       uint(idSettingMap[monitoringName]),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		_, err := monitoringDataRepo.SaveHarianData(harianData)
		if err != nil {
			fmt.Printf("Error saving harian data for %s: %v\n", monitoringName, err)
		}

	}
	err = monitoringDataRepo.Truncate()
	if err != nil {
		fmt.Printf("Error truncate data for %s: %v\n", err)
	}
	fmt.Println("Daily recap completed at:", time.Now().Format("2006-01-02 15:04:05"))
}

func monitorSchedulerChanges(useCase usecases.MonitoringDataUseCase, settingUseCase usecases.SettingUseCase) {
	for {
		settings, err := settingUseCase.GetAllSetting()
		if err != nil {
			fmt.Println("Error fetching settings:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, setting := range settings {
			if setting.Scheduler != lastSchedulers[setting.ID] {
				// Perbarui last scheduler
				lastSchedulers[setting.ID] = setting.Scheduler

				// Hapus cron job lama
				if cronJobID, exists := cronJobIDs[setting.ID]; exists {
					c.Remove(cronJobID)
				}

				// Buat cron job baru
				cronExpression := fmt.Sprintf("@every %ds", setting.Scheduler)
				newCronJobID, err := c.AddFunc(cronExpression, func() {
					runJob(useCase, entities.Setting(setting))
				})
				if err != nil {
					fmt.Printf("Error adding cron job for setting ID %d: %v\n", setting.ID, err)
					continue
				}

				// Simpan ID cron job baru
				cronJobIDs[setting.ID] = newCronJobID
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func runJob(useCase usecases.MonitoringDataUseCase, setting entities.Setting) {
	apiURL := setting.HaosURL
	token := setting.HaosToken
	fmt.Print(uint(setting.ID))
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Error creating HTTP request for setting ID %d: %v\n", setting.ID, err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching monitoring data API for setting ID %d: %v\n", setting.ID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned non-200 status code for setting ID %d: %d\n", setting.ID, resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body for setting ID %d: %v\n", setting.ID, err)
		return
	}

	var apiResponse struct {
		EntityID   string                 `json:"entity_id"`
		State      string                 `json:"state"`
		Attributes map[string]interface{} `json:"attributes"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Printf("Error parsing monitoring data for setting ID %d: %v\n", setting.ID, err)
		return
	}

	for key, value := range apiResponse.Attributes {
		if key == "friendly_name" {
			continue
		}

		valueStr := fmt.Sprintf("%v", value)

		request := entities.CreateMonitoringDataRequest{
			MonitoringName:  key,
			MonitoringValue: valueStr,
			IDSetting:       uint(setting.ID), // Tambahkan SettingID ke request
		}

		_, err := useCase.SaveMonitoringData(request)
		if err != nil {
			fmt.Printf("Error saving monitoring data (%s) for setting ID %d: %v\n", key, setting.ID, err)
		}
	}

	fmt.Printf("Monitoring data saved for setting ID %d at: %s\n", setting.ID, time.Now().Format("2006-01-02 15:04:05"))
}

func removeUnits(value string) string {
	var result []rune
	for _, r := range value {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '+' {
			result = append(result, r)
		}
	}
	return string(result)
}
