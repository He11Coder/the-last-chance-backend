package domain

import (
	"encoding/base64"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ApiPetInfo struct {
	PetID        string `json:"pet_id,omitempty"`
	TypeOfAnimal string `json:"type_of_animal,omitempty"`
	Name         string `json:"name,omitempty"`
	Info         string `json:"info,omitempty"`
	PetAvatar    string `json:"avatar"`
}

type DBPetInfo struct {
	PetID        bson.ObjectID `bson:"_id,omitempty"`
	TypeOfAnimal string        `bson:"type,omitempty"`
	Name         string        `bson:"name,omitempty"`
	Info         string        `bson:"info,omitempty"`
	PetAvatar    []byte        `bson:"avatar_url,omitempty"`
}

func (apiInfo *ApiPetInfo) ToDB() (*DBPetInfo, error) {
	dbInfo := &DBPetInfo{
		TypeOfAnimal: apiInfo.TypeOfAnimal,
		Name:         apiInfo.Name,
		Info:         apiInfo.Info,
	}

	if apiInfo.PetID != "" {
		dbID, err := bson.ObjectIDFromHex(apiInfo.PetID)
		if err != nil {
			return nil, err
		}

		dbInfo.PetID = dbID
	}

	if apiInfo.PetAvatar != "" {
		byteImage, err := base64.StdEncoding.DecodeString(apiInfo.PetAvatar)
		if err != nil {
			return nil, err
		}

		dbInfo.PetAvatar = byteImage
	}

	return dbInfo, nil
}

func (dbInfo *DBPetInfo) ToApi() *ApiPetInfo {
	apiInfo := &ApiPetInfo{
		PetID:        dbInfo.PetID.Hex(),
		TypeOfAnimal: dbInfo.TypeOfAnimal,
		Name:         dbInfo.Name,
		Info:         dbInfo.Info,
	}

	if len(dbInfo.PetAvatar) != 0 {
		apiInfo.PetAvatar = base64.StdEncoding.EncodeToString(dbInfo.PetAvatar)
	}

	return apiInfo
}

type PetIDList struct {
	PetIDs []string `json:"pet_ids"`
}

type ApiPetUpdate struct {
	TypeOfAnimal string `json:"type_of_animal,omitempty"`
	Name         string `json:"name,omitempty"`
	Info         string `json:"info,omitempty"`
	PetAvatar    string `json:"avatar,omitempty"`
}

type DBPetUpdate struct {
	TypeOfAnimal string `bson:"type,omitempty"`
	Name         string `bson:"name,omitempty"`
	Info         string `bson:"info,omitempty"`
	PetAvatar    []byte `bson:"avatar_url,omitempty"`
}

func (apiInfo *ApiPetUpdate) ToDB() (*DBPetUpdate, error) {
	dbInfo := &DBPetUpdate{
		TypeOfAnimal: apiInfo.TypeOfAnimal,
		Name:         apiInfo.Name,
		Info:         apiInfo.Info,
	}

	if apiInfo.PetAvatar != "" {
		byteImage, err := base64.StdEncoding.DecodeString(apiInfo.PetAvatar)
		if err != nil {
			return nil, err
		}

		dbInfo.PetAvatar = byteImage
	}

	return dbInfo, nil
}

type PetAdviceRequest struct {
	Animal string `json:"animal"`
	Prompt string `json:"prompt"`
}

type PetAdviceResponse struct {
	Advice string `json:"advice"`
	Animal string `json:"animal"`
	Prompt string `json:"prompt"`
}
