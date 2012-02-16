package waitmap

type WaitMap struct {

}

func New() *WaitMap {
	return &WaitMap{

	}
}

func (m *WaitMap) Get(k interface{}) interface{} {
	return nil
}

func (m *WaitMap) Set(k interface{}, v interface{}) {

}

func (m *WaitMap) Check(k interface{}) bool {
	return false
}

