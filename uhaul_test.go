package uhaul

import "testing"

func TestPack(t *testing.T) {
}

// func main() {
// b'\x01\x00\x02\x00\x03\x00\x00\x00'
// packed := Pack("hhl", 1, 2, 3)
// fmt.Printf("packed: %x\n", packed)
// fmt.Println("packed", packed)

// unpacked := Unpack("hhl", packed)
// fmt.Printf("unpacked: %x\n", unpacked)
// fmt.Println("unpacked", unpacked)

// 	meow := Pack("256s", "meow")
// 	fmt.Println("meow", meow)
// 	fmt.Println("meow:", Unpack("256s", meow))

// 	llama := Pack("I256s", 70, "llama")
// b'F\x00\x00\x00llama\x00\x00\x00\
// 70 == F?
// 	fmt.Println(llama)
// 	fmt.Println(Unpack("I256s", llama))
// }
