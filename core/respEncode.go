package core
import(
	"fmt"
)
func EncodeString(str string)( []byte){
	eStr := (fmt.Sprintf("+%s\r\n", str))
	fmt.Println("Encoded String :: ", eStr)
	return []byte(eStr)
}

func EncodeBinaryString(str string)( []byte){
	eStr := fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)
	fmt.Println("EncodeBinaryString  :: ", eStr)
	return []byte(eStr)
}