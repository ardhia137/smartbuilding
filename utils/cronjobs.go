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
	"time"
)

var (
	c             *cron.Cron
	cronJobID     cron.EntryID
	lastScheduler int
)

func StartMonitoringDataJob(useCase usecases.MonitoringDataUseCase, settingUseCase usecases.SettingUseCase, monitoringDataRepo repositories.MonitoringDataRepository) {
	c = cron.New()

	updateCronJob(useCase, settingUseCase)

	_, err := c.AddFunc("0 0 * * *", func() {
		rekapHarian(monitoringDataRepo)
	})
	if err != nil {
		fmt.Println("Error scheduling daily recap job:", err)
		return
	}

	go monitorSchedulerChanges(useCase, settingUseCase)

	select {}
}

func rekapHarian(monitoringDataRepo repositories.MonitoringDataRepository) {
	fmt.Println("Starting daily recap at:", time.Now().Format("2006-01-02 15:04:05"))

	monitoringData, err := monitoringDataRepo.FindAll()
	if err != nil {
		fmt.Println("Error fetching monitoring data:", err)
		return
	}

	totalMap := make(map[string]float64)
	countMap := make(map[string]int)

	for _, data := range monitoringData {
		cleanedValue := removeUnits(data.MonitoringValue)

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
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		_, err := monitoringDataRepo.SaveHarianData(harianData)
		if err != nil {
			fmt.Printf("Error saving harian data for %s: %v\n", monitoringName, err)
		}
	}

	fmt.Println("Daily recap completed at:", time.Now().Format("2006-01-02 15:04:05"))
}

func monitorSchedulerChanges(useCase usecases.MonitoringDataUseCase, settingUseCase usecases.SettingUseCase) {
	for {
		setting, err := settingUseCase.GetSettingByID(1)
		if err != nil {
			fmt.Println("Error fetching settings:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if setting == nil {
			fmt.Println("Setting with id 1 not found")
			time.Sleep(5 * time.Second)
			continue
		}

		if setting.Scheduler != lastScheduler {

			lastScheduler = setting.Scheduler
			updateCronJob(useCase, settingUseCase)
		}

		time.Sleep(5 * time.Second)
	}
}

func updateCronJob(useCase usecases.MonitoringDataUseCase, settingUseCase usecases.SettingUseCase) {
	setting, err := settingUseCase.GetSettingByID(1)
	if err != nil {
		fmt.Println("Error fetching settings:", err)
		return
	}
	if setting == nil {
		fmt.Println("Setting with id 1 not found")
		return
	}

	scheduler := setting.Scheduler

	if cronJobID != 0 {
		c.Remove(cronJobID)
	}

	cronExpression := fmt.Sprintf("@every %ds", scheduler)
	newCronJobID, err := c.AddFunc(cronExpression, func() {
		runJob(useCase, settingUseCase)
	})
	if err != nil {
		fmt.Println("Error adding cron job:", err)
		return
	}

	cronJobID = newCronJobID

	c.Start()
}

func runJob(useCase usecases.MonitoringDataUseCase, settingUseCase usecases.SettingUseCase) {
	setting, err := settingUseCase.GetSettingByID(1)
	if err != nil {
		fmt.Println("Error fetching settings:", err)
		return
	}

	if setting == nil {
		fmt.Println("Setting with id 1 not found")
		return
	}

	apiURL := setting.HaosURL
	token := setting.HaosToken
	globalScheduller := setting.Scheduler
	fmt.Println("scheduler atas :", globalScheduller)

	client := &http.Client{}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching monitoring data API:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned non-200 status code: %d\n", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var apiResponse struct {
		EntityID    string                 `json:"entity_id"`
		State       string                 `json:"state"`
		Attributes  map[string]interface{} `json:"attributes"`
		LastChanged string                 `json:"last_changed"`
		LastUpdated string                 `json:"last_updated"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Println("Error parsing monitoring data:", err)
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
		}

		_, err := useCase.SaveMonitoringData(request)
		if err != nil {
			fmt.Printf("Error saving monitoring data (%s): %v\n", key, err)
		}
	}
	fmt.Println("Monitoring data saved! At  : ", time.Now().Format("2006-01-02 15:04:05"))
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
