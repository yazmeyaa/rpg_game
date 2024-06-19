package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func (cm *ComponentsManager) Save(path string) error {
	out := make(map[string]string)

	cm.IterateOverStores(func(key string, cs ComponentStorer) {
		str, _ := cs.ToJSON()
		out[key] = string(str)
	})

	outStr, err := json.Marshal(out)
	if err != nil {
		return err
	}

	return cm.writeFile(path, outStr)
}

func (cm *ComponentsManager) writeFile(filename string, output []byte) error {
	return os.WriteFile(filename, output, 0644)
}

func (cm *ComponentsManager) Load(path string) error {
	loadJson := make(map[string]string)

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot open file")
	}

	json.Unmarshal(data, &loadJson)

	cm.IterateOverStores(func(key string, store ComponentStorer) {
		storeName := store.Name()
		v, exist := loadJson[storeName]
		fmt.Println(loadJson)
		if !exist {
			return
		}
		store.Load([]byte(v))

	})

	return nil
}
