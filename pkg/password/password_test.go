package password

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	fmt.Println(Hash("123"))

	fmt.Println(len("$2a$14$EnuK03eXexv0jdiOYjvJVe1ZMSaj31No/CnzTdJaG/4br7NvnbbPe"))
}
