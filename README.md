# geo country reverse
[offline-country-reverse-geocoder]( https://github.com/daveross/offline-country-reverse-geocoder) with golang 


# how to use

```golang
go get github.com/AsGz/geo/georeverse
```

```golang
package main

import (
	"fmt"

	"github.com/AsGz/geo/georeverse"
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
