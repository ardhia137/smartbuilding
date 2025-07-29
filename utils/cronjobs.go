package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/usecases"
	"strconv"
	"strings"
	"sync"

	"github.com/robfig/cron/v3"

	//"strconv"
	"time"
)

type MonitoringStatus struct {
	MonitoringAir     string `json:"monitoring air"`
	MonitoringListrik string `json:"monitoring listrik"`
}

var (
	c                   *cron.Cron
	cronJobIDs          map[int]cron.EntryID
	lastSchedulers      map[int]int
	lastHaosURL         map[int]string
	lastHaosToken       map[int]string
	monitoringStatusMap sync.Map
)

func init() {
	cronJobIDs = make(map[int]cron.EntryID)
	lastSchedulers = make(map[int]int)
	lastHaosURL = make(map[int]string)
	lastHaosToken = make(map[int]string)
}

func StartMonitoringDataJob(useCase usecases.MonitoringDataUseCase, gedungUseCase usecases.GedungUseCase, monitoringDataRepo repositories.MonitoringDataRepository, gedungRepo repositories.GedungRepository) {
	c = cron.New()

	// Inisialisasi cron jobs untuk semua settings
	gedung, err := gedungUseCase.GetAllCornJobs()
	if err != nil {
		fmt.Println("Error fetching gedung:", err)
		return
	}

	for _, gedungItem := range gedung {
		cronExpression := fmt.Sprintf("@every %ds", gedungItem.Scheduler)
		cronJobID, err := c.AddFunc(cronExpression, func() {

			gedungEntity := entities.Gedung{
				ID:           gedungItem.ID,
				NamaGedung:   gedungItem.NamaGedung,
				HaosURL:      gedungItem.HaosURL,
				HaosToken:    gedungItem.HaosToken,
				Scheduler:    gedungItem.Scheduler,
				HargaListrik: gedungItem.HargaListrik,
				JenisListrik: gedungItem.JenisListrik,
			}
			runJob(useCase, gedungEntity)
		})
		if err != nil {
			fmt.Printf("Error adding cron job for gedung ID %d: %v\n", gedungItem.ID, err)
			continue
		}

		cronJobIDs[gedungItem.ID] = cronJobID
		lastSchedulers[gedungItem.ID] = gedungItem.Scheduler
		lastHaosURL[gedungItem.ID] = gedungItem.HaosURL
	}

	// Jadwalkan rekap harian
	_, err = c.AddFunc("59 23 * * *", func() {
		rekapHarian(monitoringDataRepo, gedungRepo)
	})
	if err != nil {
		fmt.Println("Error scheduling daily recap job:", err)
		return
	}

	// Mulai monitor perubahan scheduler
	go monitorSchedulerChanges(useCase, gedungUseCase)

	c.Start()
	select {}
}

func rekapHarian(monitoringDataRepo repositories.MonitoringDataRepository, gedungRepo repositories.GedungRepository) {
	fmt.Println("Starting daily recap at:", time.Now().Format("2006-01-02 15:04:05"))

	monitoringData, err := monitoringDataRepo.FindAll()
	if err != nil {
		fmt.Println("Error fetching monitoring data:", err)
		return
	}

	gedungDataMap := make(map[uint]map[string][]float64)

	for _, data := range monitoringData {
		cleanedValue := removeUnits(data.MonitoringValue)
		value, err := strconv.ParseFloat(cleanedValue, 64)
		if err != nil {
			fmt.Printf("Error parsing value for %s: %v\n", data.MonitoringName, err)
			continue
		}

		if _, ok := gedungDataMap[data.IDGedung]; !ok {
			gedungDataMap[data.IDGedung] = make(map[string][]float64)
		}

		// Simpan nilai berdasarkan MonitoringName
		gedungDataMap[data.IDGedung][data.MonitoringName] = append(gedungDataMap[data.IDGedung][data.MonitoringName], value)
	}

	for idGedung, monitoringMap := range gedungDataMap {
		for monitoringName, values := range monitoringMap {
			total := 0.0

			if strings.HasPrefix(monitoringName, "monitoring_air_total_water_flow_") {
				if len(values) > 0 {
					total = values[len(values)-1]
				}
			} else {
				for _, val := range values {
					total += val
				}
			}

			harianData := entities.MonitoringData{
				MonitoringName:  monitoringName,
				MonitoringValue: fmt.Sprintf("%.2f", total),
				IDGedung:        idGedung,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}

			_, err := monitoringDataRepo.SaveHarianData(harianData)
			if err != nil {
				fmt.Printf("Error saving harian data for %s (IDGedung %d): %v\n", monitoringName, idGedung, err)
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
func monitorSchedulerChanges(useCase usecases.MonitoringDataUseCase, gedungUseCase usecases.GedungUseCase) {
	for {
		gedung, err := gedungUseCase.GetAllCornJobs()
		if err != nil {
			fmt.Println("Error fetching settings:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, gedungItem := range gedung {
			schedulerChanged := gedungItem.Scheduler != lastSchedulers[gedungItem.ID]
			urlChanged := gedungItem.HaosURL != lastHaosURL[gedungItem.ID]
			tokenChanged := gedungItem.HaosToken != lastHaosToken[gedungItem.ID]

			if schedulerChanged || urlChanged || tokenChanged {
				// Simpan data terbaru
				lastSchedulers[gedungItem.ID] = gedungItem.Scheduler
				lastHaosURL[gedungItem.ID] = gedungItem.HaosURL
				lastHaosToken[gedungItem.ID] = gedungItem.HaosToken

				// Hapus cron job lama jika ada
				if cronJobID, exists := cronJobIDs[gedungItem.ID]; exists {
					c.Remove(cronJobID)
				}

				// Tambahkan cron job baru dengan data terbaru
				cronExpression := fmt.Sprintf("@every %ds", gedungItem.Scheduler)
				currentGedung := gedungItem // hindari closure bug
				newCronJobID, err := c.AddFunc(cronExpression, func() {
					// Konversi GedungResponse ke Gedung
					gedungEntity := entities.Gedung{
						ID:           currentGedung.ID,
						NamaGedung:   currentGedung.NamaGedung,
						HaosURL:      currentGedung.HaosURL,
						HaosToken:    currentGedung.HaosToken,
						Scheduler:    currentGedung.Scheduler,
						HargaListrik: currentGedung.HargaListrik,
						JenisListrik: currentGedung.JenisListrik,
					}
					runJob(useCase, gedungEntity)
				})
				if err != nil {
					fmt.Printf("Error adding cron job for gedung ID %d: %v\n", gedungItem.ID, err)
					continue
				}

				cronJobIDs[gedungItem.ID] = newCronJobID
				fmt.Printf("Updated job for gedung ID %d (Scheduler: %ds, URL: %s)\n", gedungItem.ID, gedungItem.Scheduler, gedungItem.HaosURL)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func runJob(useCase usecases.MonitoringDataUseCase, gedung entities.Gedung) {
	apiURL := gedung.HaosURL
	token := gedung.HaosToken
	namaGedung := gedung.NamaGedung
	fmt.Print(uint(gedung.ID))
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Error creating HTTP request for gedung ID %d: %v\n", gedung.ID, err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching monitoring data API for gedung ID %d: %v\n", gedung.ID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned non-200 status code for gedung ID %d: %d\n", gedung.ID, resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body for gedung ID %d: %v\n", gedung.ID, err)
		return
	}

	var apiResponse struct {
		EntityID   string                 `json:"entity_id"`
		State      string                 `json:"state"`
		Attributes map[string]interface{} `json:"attributes"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Printf("Error parsing monitoring data for gedung ID %d: %v\n", gedung.ID, err)
		return
	}

	// Inisialisasi status monitoring untuk gedung ini
	currentStatus := MonitoringStatus{
		MonitoringAir:     "online",
		MonitoringListrik: "online",
	}

	// Cek apakah gedung sudah ada di memory
	if existingData, exists := monitoringStatusMap.Load(namaGedung); exists {
		if status, ok := existingData.(MonitoringStatus); ok {
			currentStatus = status
		}
	}

	for key, value := range apiResponse.Attributes {
		if key == "friendly_name" {
			continue
		}

		valueStr := fmt.Sprintf("%v", value)

		// Cek status monitoring berdasarkan prefix dan value
		if strings.HasPrefix(key, "monitoring_air_water_flow_") && strings.Contains(valueStr, "unavailable") {
			currentStatus.MonitoringAir = "offline"
		} else if strings.HasPrefix(key, "monitoring_air_water_flow_") && !strings.Contains(valueStr, "unavailable") {
			currentStatus.MonitoringAir = "online"
		}

		if strings.HasPrefix(key, "monitoring_listrik") && strings.Contains(valueStr, "0.0 A") {
			currentStatus.MonitoringListrik = "offline"
		} else if strings.HasPrefix(key, "monitoring_listrik") && !strings.Contains(valueStr, "0.0 A") {
			currentStatus.MonitoringListrik = "online"
		}

		request := entities.CreateMonitoringDataRequest{
			MonitoringName:  key,
			MonitoringValue: valueStr,
			IDGedung:        uint(gedung.ID), // Tambahkan GedungID ke request
		}
		if !strings.Contains(request.MonitoringValue, "unavailable") {
			_, err := useCase.SaveMonitoringData(request)
			if err != nil {
				fmt.Printf("Error saving monitoring data (%s) for gedung ID %d: %v\n", key, gedung.ID, err)
			}
		}
	}

	// Update status monitoring di memory
	monitoringStatusMap.Store(namaGedung, currentStatus)

	fmt.Printf("Monitoring data saved for gedung ID %d (%s) at: %s\n", gedung.ID, namaGedung, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("Status monitoring %s - Air: %s, Listrik: %s\n", namaGedung, currentStatus.MonitoringAir, currentStatus.MonitoringListrik)
}

// GetMonitoringStatus returns the current monitoring status for all buildings
func GetMonitoringStatus() map[string][]map[string]string {
	result := make(map[string][]map[string]string)

	monitoringStatusMap.Range(func(key, value interface{}) bool {
		namaGedung := key.(string)
		status := value.(MonitoringStatus)

		result[namaGedung] = []map[string]string{
			{"monitoring air": status.MonitoringAir},
			{"monitoring listrik": status.MonitoringListrik},
		}
		return true
	})

	return result
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
