package BarbDB

import (
	"encoding/base64" // To encode and decode strings in the database
	"errors"          // To return custom errors
	"os"              // To perform CRUD operations on files
	"strings"         // To split and merge strings
)

// The struct that represents the database.
type barbDB struct {
	file os.File
}

// Opens a database at the given path.
func OpenDB(path string) (*barbDB, error) {
	// Open the file
	file, fileError := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if fileError != nil {
		return nil, fileError
	}

	// Return the database struct
	return &barbDB{
		file: *file,
	}, nil
}

// Helper function to read the file.
func (db barbDB) readFile() ([]string, error) {
	// Read the file
	data, readError := os.ReadFile(db.file.Name())
	if readError != nil {
		return nil, readError
	}

	// Split the data into lines and return it
	return strings.Split(string(data), "\n"), nil
}

// Returns the value of the given key.
func (db barbDB) Get(key string) (string, error) {
	// Base64 encode the key
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))

	// Read the file
	fileContent, fileReadError := db.readFile()
	if fileReadError != nil {
		return "", fileReadError
	}

	// Loop over the lines in the file
	for i := 0; i < len(fileContent); i++ {
		splitString := strings.Split(fileContent[i], "=")
		if splitString[0] == encodedKey {
			// Decode the value and return it
			toReturn, base64DecodeError := base64.RawStdEncoding.DecodeString(splitString[1])
			if base64DecodeError != nil {
				return "", base64DecodeError
			}
			return string(toReturn), nil
		}
	}

	// Return an error if the key doesn't exist
	return "", errors.New("key not found")
}

// Sets the value of the given key.
func (db barbDB) Set(key string, value string) error {
	// Base64 encode the key and value
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))
	encodedValue := base64.RawStdEncoding.EncodeToString([]byte(value))

	// Read the file
	fileContent, fileReadError := db.readFile()
	if fileReadError != nil {
		return fileReadError
	}

	// Check if the key already exists
	for i := 0; i < len(fileContent); i++ {
		splitString := strings.Split(fileContent[i], "=")
		if splitString[0] == encodedKey {
			// Delete the key if it does
			deleteError := db.Delete(key)
			if deleteError != nil {
				return deleteError
			}
		}
	}

	// Write the key and value to the file
	_, fileWriteError := db.file.WriteString(encodedKey + "=" + encodedValue + "\n")
	if fileWriteError != nil {
		return fileWriteError
	}
	return db.file.Sync()
}

// Deletes the given key from the database.
func (db barbDB) Delete(key string) error {
	// Base64 encode the key
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))

	// Read the file
	fileContent, fileReadError := db.readFile()
	if fileReadError != nil {
		return fileReadError
	}

	// Loop over the lines in the file
	for i := 0; i < len(fileContent); i++ {
		splitString := strings.Split(fileContent[i], "=")
		if splitString[0] == encodedKey {
			// Delete the key
			fileContent = append(fileContent[:i], fileContent[i+1:]...)
		}
	}

	// Remove the key and value from the file
	fileWriteError := os.WriteFile(db.file.Name(), []byte(strings.Join(fileContent, "\n")), 0600)
	if fileWriteError != nil {
		return fileWriteError
	}
	return db.file.Sync()
}

// Closes the database.
func (db barbDB) Close() error {
	return db.file.Close()
}
