package models

type ApplicationInventory struct {
	IsActive        bool
	Application     Application
	ApplicationID   uint `gorm:"primary_key"`
	Inventory       Inventory
	InventoryID     uint   `gorm:"primary_key"`
	ApplicationUrls string `gorm:"type:text"`
}

func GetApplicationInventories() ([]*ApplicationInventory, error) {
	applicationInventories := make([]*ApplicationInventory, 0)
	err := GetDB().Preload("Application").Preload("Inventory").Find(&applicationInventories).Error
	if err != nil {
		return nil, err
	}
	return applicationInventories, nil
}

func GetApplicationInventory(id uint) *ApplicationInventory {
	var applicationInventory ApplicationInventory
	err := GetDB().Preload("Application").Preload("Inventory").First(&applicationInventory, id).Error
	if err != nil {
		return nil
	}
	return &applicationInventory
}

func SaveApplicationInventory(applicationInventory *ApplicationInventory) error {
	if GetDB().NewRecord(applicationInventory) {
		err := GetDB().Create(applicationInventory).Error
		if err != nil {
			return err
		}
	} else {
		err := GetDB().Omit("created_at").Save(applicationInventory).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteApplicationInventory(applicationInventory *ApplicationInventory) error {
	err := GetDB().Delete(applicationInventory).Error
	if err != nil {
		return err
	}
	return nil
}
