package main

func main(){
	m := make(map[int]int)
	for i := 0 ; i < 10 ; i++{
		if _,ok := m[i];ok{
			m[i]=2
		}
	}
}
