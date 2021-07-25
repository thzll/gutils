package myutils

import (
	"fmt"
	"go.uber.org/zap"
	"image"
	"image/color"
)

type Rect struct {
	X     int
	Y     int
	Count int //此点附近有几个相邻的点
}

type RectList struct {
	rect []Rect
}

func (self *RectList) Addpx(x, y int) int {
	if self.rect == nil {
		self.rect = make([]Rect, 0)
	}
	if n := self.existence(x, y); n < 0 {
		rect := Rect{X: x, Y: y}
		self.rect = append(self.rect, rect)
		return len(self.rect) - 1
	} else {
		return n
	}
}

func (self *RectList) existence(x, y int) int {
	if self.rect != nil {
		for i := 0; i < len(self.rect); i++ {
			v := &self.rect[i]
			if v.X == x && v.Y == y {
				return i
			}
		}
	} else {
		return -1
	}
	return -1
}

func (self *RectList) Count() int {
	if self.rect != nil {
		return len(self.rect)
	} else {
		return 0
	}
}

func (self *RectList) Clear() {
	if self.rect != nil {
		self.rect = self.rect[0:0]
	}
}

func (self *RectList) SetPex(img *image.RGBA, val uint8) {
	if self.rect != nil {
		for _, v := range self.rect {
			img.Set(v.X, v.Y, color.RGBA{A: 0xff, R: val, G: val, B: val})
		}

	} else {

	}
}

func (self *RectList) GetItem(i int) *Rect {
	if self.rect != nil {
		if i < len(self.rect) {
			return &self.rect[i]
		}
	}
	return nil
}

func (self *RectList) GetCenterofgravity() (int, int) {
	if self.rect != nil {
		xx, yy := 0, 0
		for _, v := range self.rect {
			xx += v.X
			yy += v.Y
		}
		return xx / self.Count(), yy / self.Count()
	} else {
		return 0, 0
	}
}

func (self *RectList) getCenterPos() (int, int) {
	if self.rect != nil {
		maxX, minX, maxY, minY := 0, 0, 0, 0
		for _, v := range self.rect {
			if maxX == 0 {
				maxX, minX, maxY, minY = v.X, v.X, v.Y, v.Y
			} else {
				if v.X > maxX {
					maxX = v.X
				}
				if v.X < minX {
					minX = v.X
				}
				if v.Y > maxY {
					maxY = v.Y
				}
				if minY > v.Y {
					minY = v.Y
				}
			}
		}
		return (maxX + minX) / 2, (maxY + minY) / 2
	} else {
		return 0, 0
	}
}

func (self *RectList) getCenterWH() (int, int) {
	if self.rect != nil {
		maxX, minX, maxY, minY := 0, 0, 0, 0
		for _, v := range self.rect {
			if maxX == 0 {
				maxX, minX, maxY, minY = v.X, v.X, v.Y, v.Y
			} else {
				if v.X > maxX {
					maxX = v.X
				}
				if v.X < minX {
					minX = v.X
				}
				if v.Y > maxY {
					maxY = v.Y
				}
				if minY > v.Y {
					minY = v.Y
				}
			}
		}
		return maxX - minX + 1, maxY - minY + 1
	} else {
		return 0, 0
	}
}

func GetImgNumByRect(img *image.RGBA, rect *Rect, rectList *RectList, v byte, strict bool) *RectList {
	rects := &RectList{rect: make([]Rect, 0)}
	for i := rect.X - 1; i <= rect.X+1; i++ {
		if i >= 0 && i < img.Rect.Dx() {
			for j := rect.Y - 1; j <= rect.Y+1; j++ {
				if j >= 0 && j < img.Rect.Dy() {
					//fmt.Println(num, x, y, "==>", i, j)
					rgba := img.RGBAAt(i, j)
					if rgba.R == v && rgba.G == v && rgba.B == v {
						if i == rect.X || j == rect.Y {
							//判断为直向的点才作为集团点继续遍历
							if rectList == nil || rectList.existence(i, j) < 0 {
								rects.Addpx(i, j)
							}
						} else {
							if !strict {
								if rectList == nil || rectList.existence(i, j) < 0 {
									rects.Addpx(i, j)
								}
							}
						}
						rect.Count += 1 //此点附近有多少个点相邻
					}
				}
			}
		}
	}
	return rects
}

func GetImgNum(img *image.RGBA, rectindex int, rectList *RectList, v byte, strict bool) {
	if rectList == nil {
		return
	}
	rect := rectList.GetItem(rectindex)
	rects := GetImgNumByRect(img, rect, rectList, v, strict)
	for i := 0; i < rects.Count(); i++ {
		rect := rects.GetItem(i)
		if rect != nil {
			index := rectList.Addpx(rect.X, rect.Y)
			GetImgNum(img, index, rectList, v, strict)
		}
	}
}

//去噪点 重心比较偏的噪点
func RemoveNoiseGravityH(img *image.RGBA, gravityH int, strict bool) int {
	rectList := &RectList{rect: make([]Rect, 0)}
	nRet := 0
	for i := 0; i < img.Rect.Dx(); i++ {
		for j := 0; j < img.Rect.Dy(); j++ {
			alphaAt := img.RGBAAt(i, j)
			if alphaAt.R == 0 {
				rectList.Clear()
				index := rectList.Addpx(i, j)
				GetImgNum(img, index, rectList, 0, strict)
				if w, h := rectList.getCenterWH(); w < 13 && h < 13 {
					if x, _ := rectList.getCenterPos(); x < 32 || x > 64 {
						//清除高度和宽度都比价小的数据
						rectList.SetPex(img, 0xff)
						nRet += rectList.Count()
					}
				}
				if _, y := rectList.GetCenterofgravity(); y > img.Rect.Dy()-gravityH || y < gravityH {
					n := rectList.Count()
					if n < 150 {
						rectList.SetPex(img, 0xff)
						nRet += rectList.Count()
					}
				}
			}
		}
	}
	zap.L().Info("RemoveNoiseGravityH",
		zap.Int("gravityH", gravityH),
		zap.Bool("strict", strict),
		zap.Int("rm", nRet),
	)
	return nRet
}

//去噪点 去除连续点比较少的点
func RemoveNoiseClutterSize(img *image.RGBA, clutterSize int, strict bool) int {
	rectList := &RectList{rect: make([]Rect, 0)}
	nRet := 0
	for i := 0; i < img.Rect.Dx(); i++ {
		for j := 0; j < img.Rect.Dy(); j++ {
			alphaAt := img.RGBAAt(i, j)
			if alphaAt.R == 0 {
				rectList.Clear()
				index := rectList.Addpx(i, j)
				if GetImgNum(img, index, rectList, 0, strict); rectList.Count() < clutterSize {
					rectList.SetPex(img, 0xff)
					nRet += rectList.Count()
				}
			}
		}
	}
	zap.L().Info(fmt.Sprintf("RemoveNoiseClutterSize clutterSize:%d strict:%v rm:%d", clutterSize, strict, nRet))
	return nRet
}

//去噪点 靠近上下边缘的点
func RemoveNoiseEdgewidth(img *image.RGBA, edgewidth int) int {
	RemoverectList := &RectList{rect: make([]Rect, 0)}
	for i := 0; i < img.Rect.Dx(); i++ {
		for j := 0; j < img.Rect.Dy(); j++ {
			alphaAt := img.RGBAAt(i, j)
			if alphaAt.R == 0 {
				if j >= img.Rect.Dy()-edgewidth || j < edgewidth {
					RemoverectList.Addpx(i, j)
				}
			}
		}
	}
	RemoverectList.SetPex(img, 0xff)
	zap.L().Info(fmt.Sprintf("RemoveNoiseEdgewidth edgewidth:%d rm:%d", edgewidth, RemoverectList.Count()))
	return RemoverectList.Count()
}

//去噪点严格的 清除比价零散的点
func RemoveNoiseStrict(img *image.RGBA, clutterSize int, ratio float32) int {
	rectList := &RectList{rect: make([]Rect, 0)}
	RemoverectList := &RectList{rect: make([]Rect, 0)}
	for i := 0; i < img.Rect.Dx(); i++ {
		for j := 0; j < img.Rect.Dy(); j++ {
			alphaAt := img.RGBAAt(i, j)
			if alphaAt.R == 0 {
				rectList.Clear()
				index := rectList.Addpx(i, j)
				GetImgNum(img, index, rectList, 0, true)
				n := rectList.Count()
				if n < clutterSize {
					w, h := rectList.getCenterWH()
					ra := float32(rectList.Count()) / float32(w*h)
					if ra < ratio {
						//清除高度和宽度都比价小的数据
						RemoverectList.Addpx(i, j)
					}
				}
			}
		}
	}
	RemoverectList.SetPex(img, 0xff)
	zap.L().Info(fmt.Sprintf("RemoveNoiseStrict clutterSize:%d ratio:%v rm:%d", clutterSize, ratio, RemoverectList.Count()))
	return RemoverectList.Count()
}

//去噪点严格的 清除比价零散的点
func RemoveNoiseByAroundNumLoop(img *image.RGBA, AroundNum int, strict bool) int {
	nRet := 0
	for {
		nn := 0
		rectList := &RectList{rect: make([]Rect, 0)}
		for i := 0; i < img.Rect.Dx(); i++ {
			for j := 0; j < img.Rect.Dy(); j++ {
				alphaAt := img.RGBAAt(i, j)
				if alphaAt.R == 0 {
					rectList.Clear()
					index := rectList.Addpx(i, j)
					rect := rectList.GetItem(index)
					rects := GetImgNumByRect(img, rect, rectList, 0, strict)
					if rects.Count() <= AroundNum {
						nn += 1
						rectList.SetPex(img, 0xff)
					}
				}
			}
		}
		nRet += nn
		if nn == 0 {
			break
		}
	}
	zap.L().Info(fmt.Sprintf("RemoveNoiseByAroundNumLoop AroundNum:%d strict:%v rm:%d", AroundNum, strict, nRet))
	return nRet
}

//去噪点严格的 清除比价零散的点
func RemoveNoiseByAroundNum(img *image.RGBA, AroundNum int, strict bool) int {
	rectList := &RectList{rect: make([]Rect, 0)}
	RemoverectList := &RectList{rect: make([]Rect, 0)}
	for i := 0; i < img.Rect.Dx(); i++ {
		for j := 0; j < img.Rect.Dy(); j++ {
			alphaAt := img.RGBAAt(i, j)
			if alphaAt.R == 0 {
				rectList.Clear()
				index := rectList.Addpx(i, j)
				rect := rectList.GetItem(index)
				rects := GetImgNumByRect(img, rect, rectList, 0, strict)
				if rects.Count() <= AroundNum {
					RemoverectList.Addpx(rect.X, rect.Y)
					//rectList.SetPex(img, 0xff)
				}
			}
		}
	}
	RemoverectList.SetPex(img, 0xff)
	zap.L().Info(fmt.Sprintf("RemoveNoiseByAroundNum AroundNum:%d strict:%v rm:%d", AroundNum, strict, RemoverectList.Count()))
	return RemoverectList.Count()
}

//横向和纵向单一方向如果只有一个点的清除掉
func RemovesignleW(img *image.RGBA, limit int) int {

	rectList := &RectList{}
	nRet := 0
	for j := 0; j < img.Rect.Dy(); j++ {
		rectList.Clear()
		for i := 0; i < img.Rect.Dx(); i++ {
			rgb := img.RGBAAt(i, j)
			if rgb.R == 0 {
				rectList.Addpx(i, j)
			}
		}
		if rectList.Count() <= limit {
			rectList.SetPex(img, 0xff)
			nRet += rectList.Count()
		}
	}
	zap.L().Info(fmt.Sprintf("RemovesignleW limit:%d  rm:%d", limit, nRet))
	return nRet
}

//横向和纵向单一方向如果只有一个点的清除掉
func RemovesignleH(img *image.RGBA, limit int) int {
	rectList := &RectList{}
	nRet := 0
	for i := 0; i < img.Rect.Dx(); i++ {
		rectList.Clear()
		for j := 0; j < img.Rect.Dy(); j++ {
			rgb := img.RGBAAt(i, j)
			if rgb.R == 0 {
				rectList.Addpx(i, j)
			}
		}
		if rectList.Count() <= limit {
			rectList.SetPex(img, 0xff)
			nRet += rectList.Count()
		}
	}
	zap.L().Info(fmt.Sprintf("RemovesignleH limit:%d  rm:%d", limit, nRet))
	return nRet
}

//func BmpToPng(fname string, gravityH int, clutterSize int, edgewidth int) error {
//	src := fname
//	dst := strings.Replace(src, ".bmp", ".bmp.png", 1)
//	fmt.Println("src=", src, " dst=", dst)
//	fIn, _ := os.Open(src)
//	defer fIn.Close()
//	if img, err := bmp.Decode(fIn); err == nil {
//		if fOut, err := os.Create(dst); err == nil {
//			defer fOut.Close()
//			switch img.(type) {
//			case *image.Alpha:
//
//			case *image.NRGBA:
//				img := img.(*image.NRGBA)
//				//subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
//				return png.Encode(fOut, img)
//			case *image.RGBA:
//				img := img.(*image.RGBA)
//				RemoveNoiseEdgewidth(img, edgewidth)            //去除上下边缘的点
//				RemoveNoiseByAroundNum(img, 1, true)            //清楚图像点周围少于或者等于指定点数的图像点
//				RemoveNoiseByAroundNumLoop(img, 1, true)        //清楚图像点周围少于或者等于指定点数的图像点
//				RemoveNoiseClutterSize(img, clutterSize, false) //去除连续点数量不大的澡点集合
//				//RemoveNoiseGravityH(img, gravityH, false) //去除垂直噪点集团重心比较靠边的点
//				//RemovesignleH(img, 1) //清楚纵向只有一个点的点
//				//RemoveNoiseStrict(img, 30, 0.6) //严格方法匹配噪点集合 去除占空比很大的噪点集合
//				//RemoveNoiseClutterSize(img, clutterSize, false) //去除连续点数量不大的澡点集合
//				//RemoveNoiseStrict(img, 10, 0.6) //严格方法匹配噪点集合 去除占空比很大的噪点集合
//				//RemoveNoiseClutterSize(img, 5, true) //去除连续点数量不大的澡点集合
//				//RemoveNoiseByAroundNum(img, 1, true) //清楚图像点周围少于或者等于指定点数的图像点
//				//RemoveNoiseStrict(img, 10, 0.6) //严格方法匹配噪点集合 去除占空比很大的噪点集合
//				return png.Encode(fOut, img)
//			case *image.Paletted:
//				img := img.(*image.Paletted)
//				//subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
//				return png.Encode(fOut, img)
//			}
//		} else {
//			return err
//		}
//	} else {
//		return err
//	}
//	return nil
//}
