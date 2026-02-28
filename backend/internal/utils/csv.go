package utils

import (
	"bytes"
	"encoding/csv"
)

func GenerateCSV(headers []string, rows [][]string) ([]byte, error) {
	var buf bytes.Buffer
	// Write UTF-8 BOM for Excel compatibility
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(&buf)

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
