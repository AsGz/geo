# geo
offline-country-reverse-geocoder with golang 

# use
```golang
package main

import (
	"fmt"

	"./georeverse"
)

func main() {
	dataPath := "./data/polygons.properties"
	c, err := georeverse.NewCountryReverser(dataPath)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.GetCountryCode(-89.234005, 41.645332)) //US
	fmt.Println(c.GetCountryCode(35.128458, 31.722628)) //PS
	fmt.Println(c.GetCountryCode(35.183690, 31.749162)) //IL
}
```
