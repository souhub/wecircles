package data

type Circle struct {
	ID        int
	Name      string
	ImagePath string
	Overview  string
	Category  string
	Owner     User
	CreatedAt string
}

func (circle *Circle) Create() (err error) {
	db := NewDB()
	defer db.Close()
	query := `INSERT INTO circles (name, image_path, overview, category, owner_id)
			  VALUES (?,?,?,?,?)`
	_, err = db.Exec(query, circle.Name, circle.ImagePath, circle.Overview, circle.Category, circle.Owner.Id)
	return
}

func (circle *Circle) Update() (err error) {
	return
}

func (circle *Circle) Delete() (err error) {
	return
}

// Delete all of the users
func ResetCircles() (err error) {
	db := NewDB()
	defer db.Close()
	query := `DELETE from circles`
	_, err = db.Exec(query)
	return
}
