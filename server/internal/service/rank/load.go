package rank

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/urlvet/urlvet/internal/constants"
)

var domainRankMap map[string]int

func LoadDomainRanks() error {
	filePath := constants.DOMAIN_RANK_FILE_PATH
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %v", err)
	}

	domainRankMap = make(map[string]int, len(records)-1)

	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		rankStr := record[0]
		domain := record[1]

		rank, err := strconv.Atoi(rankStr)
		if err != nil {
			continue // skip malformed rank
		}

		domainRankMap[domain] = rank
	}

	return nil
}
