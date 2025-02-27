package api

import (
	"producer-service/connectors"
	"producer-service/logs"
	"producer-service/rabbitmq"
	"producer-service/helpers"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// UploadCSV handles the upload of a CSV file and parses its content into JSON.
func UploadCSV(rabbitMQProducer *rabbitmq.Producer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form to retrieve the file
		err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB
		if err != nil {
			logs.Log.Errorln("Failed to parse form:", err)
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Retrieve the file from the form
		file, _, err := r.FormFile("file")
		if err != nil {
			logs.Log.Errorln("Failed to retrieve the file from form:", err)
			http.Error(w, "Unable to retrieve the file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save the file temporarily (optional)
		tempFile, err := os.CreateTemp("", "uploaded-*.csv")
		if err != nil {
			logs.Log.Errorln("Failed to save the file:", err)
			http.Error(w, "Unable to save the file", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		_, err = tempFile.ReadFrom(file)
		if err != nil {
			logs.Log.Errorln("Failed to save the file content:", err)
			http.Error(w, "Failed to save the file content", http.StatusInternalServerError)
			return
		}

		// Re-open the file for reading
		csvFile, err := os.Open(tempFile.Name())
		if err != nil {
			logs.Log.Errorln("Failed to read the uploaded file:", err)
			http.Error(w, "Failed to read the uploaded file", http.StatusInternalServerError)
			return
		}
		defer csvFile.Close()

		// Parse the CSV file
		reader := csv.NewReader(csvFile)
		records, err := reader.ReadAll()
		if err != nil {
			logs.Log.Errorln("Failed to parse CSV file:", err)
			http.Error(w, "Failed to parse CSV file", http.StatusInternalServerError)
			return
		}

		// Structure the CSV data into JSON
		if len(records) < 2 {
			logs.Log.Errorln("CSV file must have a header and at least one row")
			http.Error(w, "CSV file must have a header and at least one row", http.StatusBadRequest)
			return
		}

		headers := records[0]
		// data := []connectors.User{}

		go func() {
			// Publish each row of data to RabbitMQ
			for i, row := range records[1:] {
				if len(row) != len(headers) {
					logs.Log.Errorln("Row", i+2, "does not match the header length")
					http.Error(w, "Row "+strconv.Itoa(i+2)+" does not match the header length", http.StatusBadRequest)
					return
				}

				// Parse required fields
				id, err := strconv.Atoi(row[0])
				if err != nil {
					logs.Log.Infof("Invalid ID at row %d: %v\n", i+2, err)
					continue
				}

				// Handle optional parent_user_id
				var parentUserID *float64
				if row[7] != "-1" && row[7] != "" {
					val, err := strconv.ParseFloat(row[7], 64)
					if err == nil {
						parentUserID = &val
					}
				}

				// Create User instance
				user := connectors.User{
					ID:           id,
					FirstName:    row[1],
					LastName:     row[2],
					EmailAddress: row[3],
					CreatedAt:    row[4], // Use raw string from CSV
					DeletedAt:    row[5], // Use raw string from CSV
					MergedAt:     row[6], // Use raw string from CSV
					ParentUserID: parentUserID,
				}

				// Convert User to JSON
				userJSON, err := json.Marshal(user)
				if err != nil {
					logs.Log.Infof("Failed to marshal row %d to JSON: %v\n", i+2, err)
					continue
				}

				// Encrypt user data
				encryptionKey := viper.GetString("encryption_key")
				encryptedData, err := helpers.Encrypt(string(userJSON), encryptionKey)
				if err != nil {
					logs.Log.Infof("Failed to encrypt user data: %v\n", err)
					continue
				}

				// Publish to RabbitMQ
				err = rabbitMQProducer.Publish("csv_queue", encryptedData)
				if err != nil {
					logs.Log.Infof("Failed to publish row %d to RabbitMQ: %v\n", i+2, err)
					continue
				}
			}
		}()

		// Respond with a success message and the structured data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   records[1:],
		})
	}
}
