package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

var ShortUidFunc func(string) string = nil

func ShortUid(pool string) string {

	if ShortUidFunc != nil {
		return ShortUidFunc(pool)
	} else {
		newUUID, _ := uuid.NewRandom()
		return newUUID.String()
	}

}

func GetUUID4() string {
	//newUUID, _ := uuid.NewRandom()
	return GetUUID4_Shorter()
}

func GetUUID4_Shorter() string {
	newUUID, _ := uuid.NewRandom()
	return strings.Replace(newUUID.String(), "-", "", -1)
}

func ShortenUUIDs(strings []string) map[string]string {
	uidMap := make(map[string]string)
	uidSet := make(map[string]bool) // To track used UIDs

	for _, s := range strings {
		hash := md5.Sum([]byte(s))
		shortUID := hex.EncodeToString(hash[:])[:4]

		counter := 1
		uniqueUID := shortUID
		for uidSet[uniqueUID] {
			uniqueUID = fmt.Sprintf("%s%d", shortUID, counter)
			counter++
		}

		uidMap[s] = uniqueUID
		uidSet[uniqueUID] = true
	}

	return uidMap
}

func StringToObjectID(s string) (primitive.ObjectID, error) {
	//Generated from https://claude.ai/chat/c447750c-f01c-40c1-942c-b34634d26263
	// ObjectID needs exactly 24 hex characters (12 bytes)
	const objectIDLength = 24

	// Convert string to byte slice for manipulation
	bytes := []byte(s)

	// Create a 24-character hex string
	hexStr := ""

	// Convert each byte to hex and build the string
	for i := 0; i < len(bytes) && len(hexStr) < objectIDLength; i++ {
		hexStr += fmt.Sprintf("%02x", bytes[i])
	}

	// Pad with zeros if too short
	for len(hexStr) < objectIDLength {
		hexStr += "0"
	}

	// Truncate if too long
	if len(hexStr) > objectIDLength {
		hexStr = hexStr[:objectIDLength]
	}

	// Convert hex string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(hexStr)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create ObjectID: %w", err)
	}

	return objectID, nil
}

// Alternative simpler version that directly manipulates bytes
func StringToObjectIDSimple(s string) primitive.ObjectID {
	var objectID primitive.ObjectID

	// Convert string to bytes
	bytes := []byte(s)

	// Copy up to 12 bytes (ObjectID is 12 bytes)
	copy(objectID[:], bytes)

	// If string was shorter than 12 bytes, remaining bytes are already zero

	return objectID
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	// Generate hash with default cost (10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword checks if the provided password matches the hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
