package repository

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DataLoader handles CSV data loading with worker pool
type DataLoader struct {
	repo       *MongoRepository
	workerSize int
}

// NewDataLoader creates a new data loader
func NewDataLoader(repo *MongoRepository, workerSize int) *DataLoader {
	return &DataLoader{
		repo:       repo,
		workerSize: workerSize,
	}
}

// LoadCSV loads CSV data into MongoDB using a worker pool
func (dl *DataLoader) LoadCSV(ctx context.Context, filepath string) error {
	startTime := time.Now()

	// Open CSV file
	file, err := os.Open(filepath)
	if err != nil {
		return dl.logRefresh(ctx, startTime, "failed", 0, fmt.Sprintf("failed to open file: %v", err))
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return dl.logRefresh(ctx, startTime, "failed", 0, fmt.Sprintf("failed to read header: %v", err))
	}

	log.Printf("CSV Header: %v", header)

	// // Clear existing data
	// if err := dl.clearCollections(ctx); err != nil {
	// 	return dl.logRefresh(ctx, startTime, "failed", 0, fmt.Sprintf("failed to clear collections: %v", err))
	// }

	// Create channels for worker pool
	recordChan := make(chan CSVRecord, dl.workerSize*2)
	errorChan := make(chan error, 1)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < dl.workerSize; i++ {
		wg.Add(1)
		go dl.worker(ctx, i+1, recordChan, errorChan, &wg)
	}

	// Read and send records to workers
	rowCount := 0
	go func() {
		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errorChan <- fmt.Errorf("error reading row: %w", err)
				break
			}

			record := dl.parseCSVRow(row)
			recordChan <- record
			rowCount++
		}
		close(recordChan)
	}()

	// Wait for workers to complete
	wg.Wait()
	close(errorChan)

	// Check for errors
	if err := <-errorChan; err != nil {
		return dl.logRefresh(ctx, startTime, "failed", rowCount, err.Error())
	}

	log.Printf("Successfully loaded %d rows in %v", rowCount, time.Since(startTime))
	return dl.logRefresh(ctx, startTime, "success", rowCount, "")
}

// worker processes CSV records
func (dl *DataLoader) worker(ctx context.Context, id int, records <-chan CSVRecord, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for record := range records {
		if err := dl.processRecord(ctx, record); err != nil {
			log.Printf("Worker %d: Error processing record %s: %v", id, record.OrderID, err)
			select {
			case errors <- err:
			default:
			}
			return
		}
	}
}

func (dl *DataLoader) processRecord(ctx context.Context, record CSVRecord) error {

	// ----------- CUSTOMER  -------------
	customerColl := dl.repo.GetCollection("customers")

	customer := Customer{
		CustomerID: record.CustomerID,
		Name:       record.CustomerName,
		Email:      record.CustomerEmail,
		Address:    record.CustomerAddr,
	}

	_, err := customerColl.UpdateOne( // Insert Only If Not Exists
		ctx,
		bson.M{"customer_id": customer.CustomerID},
		bson.M{"$setOnInsert": customer},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to upsert customer: %w", err)
	}

	// ----------- PRODUCT  -------------
	productColl := dl.repo.GetCollection("products")

	unitPrice, _ := strconv.ParseFloat(record.UnitPrice, 64)
	discount, _ := strconv.ParseFloat(record.Discount, 64)

	product := Product{
		ProductID: record.ProductID,
		Name:      record.ProductName,
		Category:  record.Category,
		UnitPrice: unitPrice,
		Discount:  discount,
	}

	_, err = productColl.UpdateOne( // Insert Only If Not Exists
		ctx,
		bson.M{"product_id": product.ProductID},
		bson.M{"$setOnInsert": product},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to upsert product: %w", err)
	}

	// ----------- ORDER  -------------
	orderColl := dl.repo.GetCollection("orders")

	dateOfSale, err := time.Parse("2006-01-02", record.DateOfSale)
	if err != nil {
		return fmt.Errorf("failed to parse order date: %w", err)
	}

	quantitySold, _ := strconv.Atoi(record.QuantitySold)
	shippingCost, _ := strconv.ParseFloat(record.ShippingCost, 64)

	order := Order{
		OrderID:       record.OrderID,
		ProductID:     record.ProductID,
		CustomerID:    record.CustomerID,
		Region:        record.Region,
		DateOfSale:    dateOfSale,
		QuantitySold:  quantitySold,
		ShippingCost:  shippingCost,
		PaymentMethod: record.PaymentMethod,
	}

	_, err = orderColl.UpdateOne( // Insert Only If Not Exists
		ctx,
		bson.M{"order_id": order.OrderID},
		bson.M{"$setOnInsert": order},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("failed to upsert order: %w", err)
	}

	return nil
}

// parseCSVRow parses a CSV row into a CSVRecord
func (dl *DataLoader) parseCSVRow(row []string) CSVRecord {
	return CSVRecord{
		OrderID:       row[0],
		ProductID:     row[1],
		CustomerID:    row[2],
		ProductName:   row[3],
		Category:      row[4],
		Region:        row[5],
		DateOfSale:    row[6],
		QuantitySold:  row[7],
		UnitPrice:     row[8],
		Discount:      row[9],
		ShippingCost:  row[10],
		PaymentMethod: row[11],
		CustomerName:  row[12],
		CustomerEmail: row[13],
		CustomerAddr:  row[14],
	}
}

// // clearCollections removes all documents from collections
// func (dl *DataLoader) clearCollections(ctx context.Context) error {
// 	collections := []string{"customers", "products", "orders"}

// 	for _, coll := range collections {
// 		if _, err := dl.repo.GetCollection(coll).DeleteMany(ctx, bson.M{}); err != nil {
// 			return err
// 		}
// 	}

// 	log.Println("Cleared all collections")
// 	return nil
// }

// logRefresh logs the data refresh operation
func (dl *DataLoader) logRefresh(ctx context.Context, startTime time.Time, status string, rowsLoaded int, errorMsg string) error {
	refreshLog := RefreshLog{
		StartTime:  startTime,
		EndTime:    time.Now(),
		Status:     status,
		RowsLoaded: rowsLoaded,
		ErrorMsg:   errorMsg,
	}

	_, err := dl.repo.GetCollection("refresh_logs").InsertOne(ctx, refreshLog)
	return err
}
