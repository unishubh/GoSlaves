package slaves

import "testing"

func TestOpen_Slave(t *testing.T) {
	s := &Slave{
		Work: func(obj interface{}) interface{} {
			println(obj.(string))
			return nil
		},
	}
	s.Open()
	defer s.Close()

	c := 0
	for _, str := range []string{
		"this", "is", "to", "check",
		"one", "slave", "hehe",
	} {
		s.SendWork(str)
		c++
	}
	if c != 7 {
		panic("error")
	}
}
