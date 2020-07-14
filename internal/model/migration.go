package model

//Migration auto create or upgrade table
type Migration struct{}

//Install create the tables
func (migration *Migration) Install() error {

	tables := []interface{}{
		&Device{},
		&Client{},
		&Application{},
		&Project{},
		&MQTTUser{},
	}
	Db.AutoMigrate(tables...)
	mock()
	return nil
}

//TODO: mock some data
func mock() error {
	return nil
}
