package phone


import (
"fmt"


"github.com/nyaruka/phonenumbers"
)


func Lookup(input string) {
num, err := phonenumbers.Parse(input, "")
if err != nil {
fmt.Println("Invalid phone number")
return
}


fmt.Println("📞 Phone Metadata")
fmt.Println("Country:", phonenumbers.GetRegionCodeForNumber(num))
fmt.Println("International:", phonenumbers.Format(num, phonenumbers.INTERNATIONAL))
fmt.Println("Type:", phonenumbers.GetNumberType(num))
}
