package data

type Circle struct {
	ID        int
	Name      string
	ImagePath string
	Overview  string
	Owner     User
	Members   []User
}
