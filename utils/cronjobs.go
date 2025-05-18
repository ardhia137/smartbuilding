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
	"strings"

	//"strconv"
	"time"
)

var (
	c              *cron.Cron
	cronJobIDs     map[int]cron.EntryID // Menyimpan ID cron job untuk setiap setting
	lastSchedulers map[int]int          // Menyimpan last scheduler untuk setiap setting
	lastHaosURL    map[int]string
	lastHaosToken  map[int]string
)

func init() {
	cronJobIDs = make(map[int]cron.EntryID)
	lastSchedulers = make(map[int]int)
	lastHaosURL = make(map[int]string)
	lastHaosToken = make(map[int]string)
}

func StartMonitoringDataJob(useCase usecases.MonitoringDataUseCase, settingUseCase usecases.SettingUseCase, monitoringDataRepo repositories.MonitoringDataRepository, settingRepo repositories.SettingRepository) {
	c = cron.New()

	// Inisialisasi cron jobs untuk semua settings
	settings, err := settingUseCase.GetAllCornJobs()
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
		lastHaosURL[setting.ID] = setting.HaosURL
	}

	// Jadwalkan rekap harian
	_, err = c.AddFunc("59 23 * * *", func() {
		rekapHarian(monitoringDataRepo, settingRepo)
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

func rekapHarian(monitoringDataRepo repositories.MonitoringDataRepository, settingRepo repositories.SettingRepository) {
	fmt.Println("Starting daily recap at:", time.Now().Format("2006-01-02 15:04:05"))

	monitoringData, err := monitoringDataRepo.FindAll()
	if err != nil {
		fmt.Println("Error fetching monitoring data:", err)
		return
	}

	// Buat map untuk menyimpan data per IDSetting
	settingDataMap := make(map[uint]map[string][]float64)

	for _, data := range monitoringData {
		cleanedValue := removeUnits(data.MonitoringValue)
		value, err := strconv.ParseFloat(cleanedValue, 64)
		if err != nil {
			fmt.Printf("Error parsing value for %s: %v\n", data.MonitoringName, err)
			continue
		}

		// Buat entri baru di map jika belum ada
		if _, ok := settingDataMap[data.IDSetting]; !ok {
			settingDataMap[data.IDSetting] = make(map[string][]float64)
		}

		// Simpan nilai berdasarkan MonitoringName
		settingDataMap[data.IDSetting][data.MonitoringName] = append(settingDataMap[data.IDSetting][data.MonitoringName], value)
	}

	// Looping setiap IDSetting
	for idSetting, monitoringMap := range settingDataMap {
		for monitoringName, values := range monitoringMap {
			total := 0.0

			if strings.HasPrefix(monitoringName, "monitoring_air_total_water_flow_") {
				// Gunakan nilai terakhir jika ada
				if len(values) > 0 {
					total = values[len(values)-1]
				}
			} else {
				// Hitung total dari semua nilai
				for _, val := range values {
					total += val
				}
			}

			harianData := entities.MonitoringData{
				MonitoringName:  monitoringName,
				MonitoringValue: fmt.Sprintf("%.2f", total),
				IDSetting:       idSetting,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}

			_, err := monitoringDataRepo.SaveHarianData(harianData)
			if err != nil {
				fmt.Printf("Error saving harian data for %s (IDSetting %d): %v\n", monitoringName, idSetting, err)
			}
		}
	}

	// Hapus data setelah direkap
	err = monitoringDataRepo.Truncate()
	if err != nil {
		fmt.Printf("Error truncating data: %v\n", err)
	}

	fmt.Println("Daily recap completed at:", time.Now().Format("2006-01-02 15:04:05"))
}
func monitorSchedulerChanges(useCase usecases.MonitoringDataUseCase, settingUseCase usecases.SettingUseCase) {
	for {
		settings, err := settingUseCase.GetAllCornJobs()
		if err != nil {
			fmt.Println("Error fetching settings:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, setting := range settings {
			schedulerChanged := setting.Scheduler != lastSchedulers[setting.ID]
			urlChanged := setting.HaosURL != lastHaosURL[setting.ID]
			tokenChanged := setting.HaosToken != lastHaosToken[setting.ID]

			if schedulerChanged || urlChanged || tokenChanged {
				// Simpan data terbaru
				lastSchedulers[setting.ID] = setting.Scheduler
				lastHaosURL[setting.ID] = setting.HaosURL
				lastHaosToken[setting.ID] = setting.HaosToken

				// Hapus cron job lama jika ada
				if cronJobID, exists := cronJobIDs[setting.ID]; exists {
					c.Remove(cronJobID)
				}

				// Tambahkan cron job baru dengan data terbaru
				cronExpression := fmt.Sprintf("@every %ds", setting.Scheduler)
				currentSetting := setting // hindari closure bug
				newCronJobID, err := c.AddFunc(cronExpression, func() {
					runJob(useCase, entities.Setting(currentSetting))
				})
				if err != nil {
					fmt.Printf("Error adding cron job for setting ID %d: %v\n", setting.ID, err)
					continue
				}

				cronJobIDs[setting.ID] = newCronJobID
				fmt.Printf("Updated job for setting ID %d (Scheduler: %ds, URL: %s)\n", setting.ID, setting.Scheduler, setting.HaosURL)
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
		if !strings.Contains(request.MonitoringValue, "Unavailable") {
			_, err := useCase.SaveMonitoringData(request)
			if err != nil {
				fmt.Printf("Error saving monitoring data (%s) for setting ID %d: %v\n", key, setting.ID, err)
			}
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
