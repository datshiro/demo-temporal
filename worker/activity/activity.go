package activity

import (
	"app/shared"
	"context"
	"errors"
	"log"
	"math/rand"
)

func RefreshSegment(ctx context.Context, data shared.RequestDetails) (string, error) {
	log.Printf("Refreshing segment %d", data.SegmentID)
	isError := rand.Intn(2)
	if isError == 1 {
		return "", errors.New("refresh failed")
	}
	return generateTransactionID("REFRESH", data.SegmentID), nil
}

func ConvertSegment(ctx context.Context, data shared.RequestDetails) (string, error) {
	log.Printf("Converting segment %d", data.SegmentID)
	isError := rand.Intn(2)
	if isError == 1 {
		return "", errors.New("convert failed")
	}
	return generateTransactionID("CONVERT", data.SegmentID), nil
}

func ExportSegment(ctx context.Context, data shared.RequestDetails) (string, error) {
	log.Printf("Exporting segment %d", data.SegmentID)
	isError := rand.Intn(2)
	if isError == 1 {
		return "", errors.New("export failed")
	}
	return generateTransactionID("EXPORT", data.SegmentID), nil
}

func generateTransactionID(prefix string, length int) string {
	randChars := make([]byte, length)
	for i := range randChars {
		allowedChars := "0123456789"
		randChars[i] = allowedChars[rand.Intn(len(allowedChars))]
	}
	return prefix + string(randChars)
}
