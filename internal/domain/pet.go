package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type ApiPetInfo struct {
	PetID        string `json:"pet_id,omitempty"`
	TypeOfAnimal string `json:"type_of_animal"`
	Name         string `json:"name,omitempty"`
	Info         string `json:"info,omitempty"`
	PetAvatar    string `json:"avatar,omitempty"`
}

type DBPetInfo struct {
	PetID        bson.ObjectID `bson:"_id,omitempty"`
	TypeOfAnimal string        `bson:"type"`
	Name         string        `bson:"name,omitempty"`
	Info         string        `bson:"info,omitempty"`
	PetAvatar    string        `bson:"avatar_url,omitempty"`
}

func (apiInfo *ApiPetInfo) ToDB() (*DBPetInfo, error) {
	dbInfo := &DBPetInfo{
		TypeOfAnimal: apiInfo.TypeOfAnimal,
		Name:         apiInfo.Name,
		Info:         apiInfo.Info,
		PetAvatar:    apiInfo.PetAvatar,
	}

	dbID, err := bson.ObjectIDFromHex(apiInfo.PetID)
	if err != nil {
		return nil, err
	}

	dbInfo.PetID = dbID

	return dbInfo, nil
}

func (dbInfo *DBPetInfo) ToApi() *ApiPetInfo {
	apiInfo := &ApiPetInfo{
		PetID:        dbInfo.PetID.Hex(),
		TypeOfAnimal: dbInfo.TypeOfAnimal,
		Name:         dbInfo.Name,
		Info:         dbInfo.Info,
		PetAvatar:    dbInfo.PetAvatar,
	}

	return apiInfo
}

type PetIDList struct {
	PetIDs []string `json:"pet_ids"`
}
