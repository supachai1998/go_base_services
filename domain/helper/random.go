package domain

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/samber/lo"
	"gorm.io/datatypes"
)

func RandomOneOfLength[T any](el []T) T {
	elLen := len(el)
	if elLen == 0 {
		return el[0]
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	randInt := r.Intn(elLen)
	return el[randInt]
}

func RandomMinMaxInt(min, max int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return min + r.Intn(max-min)
}
func RandomMinMaxFloat64(min, max float64) float64 {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	if min == max {
		return min
	}
	if min > max {
		min, max = max, min
	}
	return min + r.Float64()*(max-min)
}
func RandomStringLen(length int, chars ...string) string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if len(chars) > 0 {
		char = chars[0]
	}
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, char[r.Intn(len(char))])
	}
	return string(result)
}
func RandomPhoneTH() string {
	phone := []string{"08", "09"}
	return RandomOneOfLength(phone) + RandomStringLen(8, "0123456789")
}

func RandomTime() time.Time {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	min := time.Now().Unix()
	max := time.Now().AddDate(365, 0, 0).Unix()
	delta := max - min
	sec := min + r.Int63n(delta)
	return time.Unix(sec, 0)
}

func RandomDatatypesJSONFromSliceString(str []string) datatypes.JSON {
	if len(str) == 0 {
		return nil
	}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	// we bias more than 0
	lenForSlice := r.Intn(len(str))
	if lenForSlice == 0 {
		lenForSlice = 1
	}
	var result []string
	var keepInd []int
	for i := 0; i < lenForSlice; i++ {
		// random index from str
		index := r.Intn(len(str))
		if len(keepInd) > 0 {
			for {
				if lo.Contains(keepInd, index) {
					index = r.Intn(len(str))
				} else {
					break
				}
			}
		}
		result = append(result, str[index])
		keepInd = append(keepInd, index)
	}
	if len(result) == 0 {
		// return empty array
		return datatypes.JSON([]byte(`[]`))
	}
	var b []byte
	b, _ = json.Marshal(result)
	return datatypes.JSON(b)

}
