package cxlisthashtime


func AddHashValue(hashList *[]byte, newHashValue []byte, newTimestamp []byte, hashSize int, timeSize int, maxEntries int) {
	entrySize := hashSize + timeSize

	// Check if the hash value already exists in the list and replace its timestamp if needed
	for i := 0; i < len(*hashList)/entrySize; i++ {
		offset := i * entrySize
		if hashEquals((*hashList)[offset:offset+hashSize], newHashValue, hashSize) {
			copyToSlice((*hashList)[offset+hashSize:offset+entrySize], newTimestamp, timeSize)
			return
		}
	}

	// If the list is full, remove the oldest entry (by timestamp)
	if len(*hashList) == maxEntries*entrySize {
		oldestIndex := 0
		oldestTimestamp := (*hashList)[hashSize : hashSize+timeSize]
		for i := 1; i < maxEntries; i++ {
			offset := i * entrySize
			currentTimestamp := (*hashList)[offset+hashSize : offset+entrySize]
			if compareTimestamps(currentTimestamp, oldestTimestamp, timeSize) < 0 {
				oldestIndex = i
				oldestTimestamp = currentTimestamp
			}
		}
		// Shift the entries to the left
		copy((*hashList)[oldestIndex*entrySize:], (*hashList)[(oldestIndex+1)*entrySize:])
		*hashList = (*hashList)[:len(*hashList)-entrySize]
	}

	// Find the position to insert the new entry
	insertIndex := len(*hashList) / entrySize
	for i := 0; i < len(*hashList)/entrySize; i++ {
		offset := i * entrySize
		currentTimestamp := (*hashList)[offset+hashSize : offset+entrySize]
		if compareTimestamps(newTimestamp, currentTimestamp, timeSize) < 0 {
			insertIndex = i
			break
		}
	}

	// Shift the entries to the right
	*hashList = append(*hashList, make([]byte, entrySize)...)
	copy((*hashList)[(insertIndex+1)*entrySize:], (*hashList)[insertIndex*entrySize:])

	// Insert the new entry
	offset := insertIndex * entrySize
	copyToSlice((*hashList)[offset:offset+hashSize], newHashValue, hashSize)
	copyToSlice((*hashList)[offset+hashSize:offset+entrySize], newTimestamp, timeSize)
}


func DeleteHashValue(hashList *[]byte, hashValueToDelete []byte, hashSize int, timeSize int) {
	entrySize := hashSize + timeSize

	for i := 0; i < len(*hashList)/entrySize; i++ {
		offset := i * entrySize
		if hashEquals((*hashList)[offset:offset+hashSize], hashValueToDelete, hashSize) {
			// Shift the entries to the left
			copy((*hashList)[offset:], (*hashList)[offset+entrySize:])
			*hashList = (*hashList)[:len(*hashList)-entrySize]
			break
		}
	}
}


func hashEquals(a, b []byte, hashSize int) bool {
	for i := 0; i < hashSize; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func copyToSlice(dst, src []byte, size int) {
	for i := 0; i < size; i++ {
		dst[i] = src[i]
	}
}

func compareTimestamps(a, b []byte, timeSize int) int {
	for i := 0; i < timeSize; i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
