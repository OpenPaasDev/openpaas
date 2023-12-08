package state

type Db struct {
	folder string
}

func Init(folder string) *Db {
	return &Db{folder: folder}
}
