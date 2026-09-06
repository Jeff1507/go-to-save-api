package query

func calculateDiscount(originalPrice, finalPrice int64) int32 {
	if originalPrice <= 0 {
		return 0
	}

	return int32(100 - (finalPrice * 100 / originalPrice))
}
