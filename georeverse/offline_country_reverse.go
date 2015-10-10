//把一个php的离线gps定位到国家的程序翻译成go
//https://github.com/daveross/offline-country-reverse-geocoder
package georeverse

import (
	"bufio"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	POLYGON        = "POLYGON"
	MULTI_POLYGON  = "MULTIPOLYGON"
	POLYGONS_SPLIT = ")),(("
)

type CountryPolygon struct {
	CountryCode string
	PonitList   string
}

type CountryReverser struct {
	allCountryPolygonInfo []CountryPolygon
}

func myTrim(src string) string {
	src = strings.TrimSpace(src)
	src = strings.Trim(src, "(")
	src = strings.Trim(src, ")")
	return src
}

func NewCountryReverser(dataPath string) (*CountryReverser, error) {
	c := new(CountryReverser)
	err := c.load(dataPath)
	return c, err
}

func (c *CountryReverser) GetCountryCode(longitude, latitude float64) string {
	for _, poly := range c.allCountryPolygonInfo {
		if c.pointInPolygon(longitude, latitude, poly.PonitList) {
			return poly.CountryCode
		}
	}
	return ""
}

//这里可以转成json等可读性强的数据格式
//但是仍然要,采用跟之前同样的数据格式,虽然不好读取,但是方便今后数据更新
func (c *CountryReverser) load(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	var line string
	for {
		line, err = rd.ReadString('\n')
		line = strings.TrimSpace(line)
		info := strings.SplitN(line, "=", 2)
		countryCode := info[0]
		data := info[1]
		pointInfo := strings.SplitN(data, " ", 2)
		polygonType := pointInfo[0]
		data = myTrim(pointInfo[1])
		var polyInfo CountryPolygon
		if POLYGON == strings.ToUpper(polygonType) {
			polyInfo.CountryCode = countryCode
			polyInfo.PonitList = data
			c.allCountryPolygonInfo = append(c.allCountryPolygonInfo, polyInfo)
		} else if MULTI_POLYGON == strings.ToUpper(polygonType) {
			polygons := strings.Split(data, POLYGONS_SPLIT)
			for _, p := range polygons {
				polyInfo.CountryCode = countryCode
				polyInfo.PonitList = p
				c.allCountryPolygonInfo = append(c.allCountryPolygonInfo, polyInfo)
			}
		}
		if io.EOF == err {
			break
		}
	}
	return nil
}

//计算一个点是不是在多边形里面
//参考：http://alienryderflex.com/polygon/
func (c *CountryReverser) pointInPolygon(targetX, targetY float64, pointsList string) bool {
	var polyX, polyY []float64
	points := strings.Split(pointsList, ",")
	polyCorners := len(points)
	for _, p := range points {
		p = myTrim(p)
		info := strings.Split(p, " ")
		pX, _ := strconv.ParseFloat(info[0], 64)
		pY, _ := strconv.ParseFloat(info[1], 64)
		/**
		 * Performance optimization:
		 * If the first pair of coordinates are more than 90deg
		 * (1/4 of the Earth's circumference) in any direction,
		 * the answer is "no".
		 */
		if pX != 0 && int(math.Abs(pY-targetY)) > 90 {
			return false
		}
		if pY != 0 && int(math.Abs(pX-targetX)) > 90 {
			return false
		}
		polyX = append(polyX, pX)
		polyY = append(polyY, pY)
	}
	j := polyCorners - 1
	oddNodes := false
	for i := 0; i < polyCorners; i++ {
		if polyY[i] < targetY && polyY[j] >= targetY || polyY[j] < targetY && polyY[i] >= targetY {
			if polyX[i] <= targetX || polyX[j] <= targetX {
				oddNodes = (oddNodes != (polyX[i]+(targetY-polyY[i])/(polyY[j]-polyY[i])*(polyX[j]-polyX[i]) < targetX))
			}
		}
		j = i
	}
	return oddNodes
}
