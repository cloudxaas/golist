package cxlisthashtime

type List struct {
	data      []byte
	size      int
	hashSize  int
	timeSize  int
	entrySize int
	maxEntries int
}

func New(hashSize, timeSize, maxEntries int) *HashList {
	entrySize := hashSize + timeSize
	data := make([]byte, entrySize*maxEntries)
	return &List{
		data:      data,
		hashSize:  hashSize,
		timeSize:  timeSize,
		entrySize: entrySize,
		maxEntries: maxEntries,
	}
}


func (hl *List) Add(newHashValue []byte, newTimestamp []byte) {
	// Check if the hash value already exists in the list and replace its timestamp if needed
	for i := 0; i < hl.size; i++ {
		offset := i * entrySize
		if hashEquals(hl.data[offset:offset+hashSize], newHashValue) {
			copyToSlice(hl.data[offset+hashSize:offset+entrySize], newTimestamp)
			return
		}
	}

	// If the list is full, remove the oldest entry (by timestamp)
	if hl.size == maxEntries {
		oldestIndex := 0
		oldestTimestamp := hl.data[hashSize : hashSize+timeSize]
		for i := 1; i < hl.size; i++ {
			offset := i * entrySize
			currentTimestamp := hl.data[offset+hashSize : offset+entrySize]
			if compareTimestamps(currentTimestamp, oldestTimestamp) < 0 {
				oldestIndex = i
				oldestTimestamp = currentTimestamp
			}
		}
		// Shift the entries to the left
		copy(hl.data[oldestIndex*entrySize:], hl.data[(oldestIndex+1)*entrySize:])
		hl.size--
	}

	// Find the position to insert the new entry
	insertIndex := hl.size
	for i := 0; i < hl.size; i++ {
		offset := i * entrySize
		currentTimestamp := hl.data[offset+hashSize : offset+entrySize]
		if compareTimestamps(newTimestamp, currentTimestamp) < 0 {
			insertIndex = i
			break
		}
	}

	// Shift the entries to the right
	copy(hl.data[(insertIndex+1)*entrySize:], hl.data[insertIndex*entrySize:])

	// Insert the new entry
	offset := insertIndex * entrySize
	copyToSlice(hl.data[offset:offset+hashSize], newHashValue)
	copyToSlice(hl.data[offset+hashSize:offset+entrySize], newTimestamp)

	hl.size++
}

func hashEquals(a, b []byte) bool {
	for i := 0; i < hashSize; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func copyToSlice(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		dst[i] = src[i]
	}
}

func compareTimestamps(a, b []byte) int {
	for i := 0; i < timeSize; i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
