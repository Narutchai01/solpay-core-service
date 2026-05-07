package main

import (
	"fmt"
	"os"

	"github.com/Narutchai01/solpay-core-service/internal/infra/supabase"
)

func main() {
	s := supabase.NewSupabaseStorage(os.Getenv("SUPABASE_PRIVATE_KEY"), os.Getenv("SUPABASE_URL"))
	url, err := s.UploadFile("slip", "test.txt", []byte("hello"))
	fmt.Println("Upload slip:", url, err)

	url, err = s.UploadFile("users", "test.txt", []byte("hello"))
	fmt.Println("Upload users:", url, err)
}
