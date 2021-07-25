package dbOrm

import (
	"fmt"
	"log"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type InTable interface {
	GetDbNameSpace() string
	TableName() string
}

type argInt []int

// get int by index from int slice
func (a argInt) Get(i int, args ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

type Page struct {
	Errcode  string
	PageSize int64
	Offset   int64
	Total    int64 //总页数
	Rows     interface{}
}
type TimeModel struct {
	CreatedAt time.Time  `gorm:"DEFAULT:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func (s *TimeModel) GetDbNameSpace() string {
	return "default"
}

// StrTo is the target string
type StrTo string

// Set string
func (f *StrTo) Set(v string) {
	if v != "" {
		*f = StrTo(v)
	} else {
		f.Clear()
	}
}

// Clear string
func (f *StrTo) Clear() {
	*f = StrTo(0x1E)
}

// Exist check string exist
func (f StrTo) Exist() bool {
	return string(f) != string(0x1E)
}

// Bool string to bool
func (f StrTo) Bool() (bool, error) {
	return strconv.ParseBool(f.String())
}

// Float32 string to float32
func (f StrTo) Float32() (float32, error) {
	v, err := strconv.ParseFloat(f.String(), 32)
	return float32(v), err
}

// Float64 string to float64
func (f StrTo) Float64() (float64, error) {
	return strconv.ParseFloat(f.String(), 64)
}

// Int string to int
func (f StrTo) Int() (int, error) {
	v, err := strconv.ParseInt(f.String(), 10, 32)
	return int(v), err
}

// Int8 string to int8
func (f StrTo) Int8() (int8, error) {
	v, err := strconv.ParseInt(f.String(), 10, 8)
	return int8(v), err
}

// Int16 string to int16
func (f StrTo) Int16() (int16, error) {
	v, err := strconv.ParseInt(f.String(), 10, 16)
	return int16(v), err
}

// Int32 string to int32
func (f StrTo) Int32() (int32, error) {
	v, err := strconv.ParseInt(f.String(), 10, 32)
	return int32(v), err
}

// Int64 string to int64
func (f StrTo) Int64() (int64, error) {
	v, err := strconv.ParseInt(f.String(), 10, 64)
	if err != nil {
		i := new(big.Int)
		ni, ok := i.SetString(f.String(), 10) // octal
		if !ok {
			return v, err
		}
		return ni.Int64(), nil
	}
	return v, err
}

// Uint string to uint
func (f StrTo) Uint() (uint, error) {
	v, err := strconv.ParseUint(f.String(), 10, 32)
	return uint(v), err
}

// Uint8 string to uint8
func (f StrTo) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(f.String(), 10, 8)
	return uint8(v), err
}

// Uint16 string to uint16
func (f StrTo) Uint16() (uint16, error) {
	v, err := strconv.ParseUint(f.String(), 10, 16)
	return uint16(v), err
}

// Uint32 string to uint32
func (f StrTo) Uint32() (uint32, error) {
	v, err := strconv.ParseUint(f.String(), 10, 32)
	return uint32(v), err
}

// Uint64 string to uint64
func (f StrTo) Uint64() (uint64, error) {
	v, err := strconv.ParseUint(f.String(), 10, 64)
	if err != nil {
		i := new(big.Int)
		ni, ok := i.SetString(f.String(), 10)
		if !ok {
			return v, err
		}
		return ni.Uint64(), nil
	}
	return v, err
}

// String string to string
func (f StrTo) String() string {
	if f.Exist() {
		return string(f)
	}
	return ""
}

func SetObjValue(obj interface{}, params map[string]interface{}) {
	val := reflect.ValueOf(obj)
	typ := reflect.Indirect(val).Type()
	if val.Kind() != reflect.Ptr {
		log.Println(fmt.Errorf("Err <ParseValue> cannot use non-ptr model struct `%s`", typ.Name()))
		return
	}
	ind := reflect.Indirect(val)
	for k, v := range params {
		for i := 0; i < ind.NumField(); i++ {
			fieldName := typ.Field(i).Name
			if k == fieldName || k == snakeString(fieldName) {
				setFieldValue(ind.Field(i), v)
			} else {
			}
		}
	}
}

func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
	case string:
		s = fmt.Sprintf("%s", v)
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("\"%v\"", v)
	}
	return s
}

func parseMdParam(obj interface{}, params map[string]interface{}) {
	if len(params) == 0 {
		return
	}
	val := reflect.ValueOf(obj)
	typ := reflect.Indirect(val).Type()
	if val.Kind() != reflect.Ptr {
		log.Println(fmt.Errorf("Err <ParseValue> cannot use non-ptr model struct `%s`", typ.Name()))
		return
	}
	ind := reflect.Indirect(val)
	for k, v := range params {
		for i := 0; i < ind.NumField(); i++ {
			if k == typ.Field(i).Name || k == snakeString(typ.Field(i).Name) {
				newvalue := parseValue(ind.Field(i), v)
				params[k] = newvalue
			} else {
			}
		}
	}
}
func parseValue(ind reflect.Value, value interface{}) interface{} {
	switch ind.Kind() {
	case reflect.Bool:
		if value == nil {
			return false
		} else if v, ok := value.(bool); ok {
			return v
		} else {
			v, _ := StrTo(ToStr(value)).Bool()
			return v
		}

	case reflect.String:
		if value == nil {
			return ""
		} else {
			return ToStr(value)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == nil {
			return 0
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return val.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return int64(val.Uint())
			default:
				v, _ := StrTo(ToStr(value)).Int64()
				return v
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == nil {
			return uint(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return uint64(val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return val.Uint()
			default:
				v, _ := StrTo(ToStr(value)).Uint64()
				return v
			}
		}
	case reflect.Float64, reflect.Float32:
		if value == nil {
			return float64(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Float64:
				return val.Float()
			default:
				v, _ := StrTo(ToStr(value)).Float64()
				return v
			}
		}

	case reflect.Struct:
		if value == nil {
			return nil
		} else if _, ok := ind.Interface().(time.Time); ok {
			var str string
			switch d := value.(type) {
			case time.Time:
				return reflect.ValueOf(d)
			case []byte:
				str = string(d)
			case string:
				str = d
			}
			if str != "" {
				if len(str) >= 19 {
					str = str[:19]
					t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
					if err == nil {
						t = t.In(time.Local)
						return reflect.ValueOf(t)
					}
				} else if len(str) >= 10 {
					str = str[:10]
					t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
					if err == nil {
						return reflect.ValueOf(t)
					}
				}
			}
		}
	}
	return nil
}

// set field value to row container
func setFieldValue(ind reflect.Value, value interface{}) {
	switch ind.Kind() {
	case reflect.Bool:
		if value == nil {
			ind.SetBool(false)
		} else if v, ok := value.(bool); ok {
			ind.SetBool(v)
		} else {
			v, _ := StrTo(ToStr(value)).Bool()
			ind.SetBool(v)
		}

	case reflect.String:
		if value == nil {
			ind.SetString("")
		} else {
			ind.SetString(ToStr(value))
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == nil {
			ind.SetInt(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ind.SetInt(val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ind.SetInt(int64(val.Uint()))
			default:
				v, _ := StrTo(ToStr(value)).Int64()
				ind.SetInt(v)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == nil {
			ind.SetUint(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				ind.SetUint(uint64(val.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ind.SetUint(val.Uint())
			default:
				v, _ := StrTo(ToStr(value)).Uint64()
				ind.SetUint(v)
			}
		}
	case reflect.Float64, reflect.Float32:
		if value == nil {
			ind.SetFloat(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Float64:
				ind.SetFloat(val.Float())
			default:
				v, _ := StrTo(ToStr(value)).Float64()
				ind.SetFloat(v)
			}
		}

	case reflect.Struct:
		if value == nil {
			ind.Set(reflect.Zero(ind.Type()))

		} else if _, ok := ind.Interface().(time.Time); ok {
			var str string
			switch d := value.(type) {
			case time.Time:
				ind.Set(reflect.ValueOf(d))
			case []byte:
				str = string(d)
			case string:
				str = d
			}
			if str != "" {
				if len(str) >= 19 {
					str = str[:19]
					t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
					if err == nil {
						t = t.In(time.Local)
						ind.Set(reflect.ValueOf(t))
					}
				} else if len(str) >= 10 {
					str = str[:10]
					t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
					if err == nil {
						ind.Set(reflect.ValueOf(t))
					}
				}
			}
		}
	}
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func snakeField(args map[string]interface{}) map[string]interface{} {
	newArgs := make(map[string]interface{})
	for k, v := range args {
		snakeK := snakeString(k)
		newArgs[snakeK] = v
	}
	return newArgs
}

func Add(table InTable) (int64, error) {
	db := getDb(table.GetDbNameSpace()).Debug()
	db = db.Model(table).Create(table)
	return db.RowsAffected, db.Error
}

func InsertOrUpdate(table InTable, args map[string]interface{}) (int64, error) {
	db := getDb(table.GetDbNameSpace()).Debug()
	//db = db.Create(table)
	db = db.Model(table)
	db = db.Updates(snakeField(args))
	if db.RowsAffected == 0 || db.Error != nil {
		db = db.Create(table)
	}
	return db.RowsAffected, db.Error
}

func UpdateOrInsert(table InTable, args map[string]interface{}, where string) (int64, error) {
	db := getDb(table.GetDbNameSpace()).Debug()
	//db = db.Create(table)
	db = db.Model(table)
	if where != "" {
		db = db.Where(where)
	}
	db = db.Updates(snakeField(args))
	if db.RowsAffected == 0 || db.Error != nil {
		db = db.Create(table)
	}
	return db.RowsAffected, db.Error
}

func ClearByUpdateLast(table InTable, last time.Duration) (int64, error) {
	db := getDb(table.GetDbNameSpace())
	db.Exec(fmt.Sprintf("DELETE FROM %s WHERE updated_at<'%s'",
		table.TableName(),
		time.Now().Add(-last).Format("2006-01-02 15:04:05")),
	)
	db.Exec(fmt.Sprintf("OPTIMIZE TABLE %s", table.TableName()))
	return db.RowsAffected, db.Error
}

func DelByWhere(table InTable, where string) (int64, error) {
	db := getDb(table.GetDbNameSpace())
	//db = db.Delete(table, where)
	db = db.Unscoped().Delete(table, where)

	return db.RowsAffected, db.Error
}

func GetCount(table InTable, where string) (int64, error) {
	db := getDb(table.GetDbNameSpace())
	var count int64
	db = db.Model(table).Where(where).Count(&count)
	return count, db.Error
}

func GetList(table InTable, vList interface{}, fields ...string) (int64, error) {
	db := getDb(table.GetDbNameSpace())
	db = db.Model(table)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Find(vList)
	return db.RowsAffected, db.Error
}

func GetListByWhere(table InTable, vList interface{}, where string, fields ...string) (int64, error) {
	db := getDb(table.GetDbNameSpace())
	db = db.Model(table).Debug()
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Where(where).Find(vList)
	return db.RowsAffected, db.Error
}

func GetByWhere(table InTable, where string, fields ...string) (int64, error) {
	db := getDb(table.GetDbNameSpace())
	db = db.Model(table)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Where(where).First(table)
	return db.RowsAffected, db.Error
}

//获取 通过更新时间
func GetListByUpdateTime(table InTable, vList interface{}, updateTime int64, fields ...string) (int64, error) {
	db := getDb(table.GetDbNameSpace())
	db = db.Model(table)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Where("? < unix_timestamp(updated_at)", updateTime).Find(vList)
	return db.RowsAffected, db.Error
}

func GetListWithPage(table InTable, vList interface{}, pageSize int, offset int, sort string, sortOrder string, where string) (page Page) {
	if sort != "" {
		switch sort[0] {
		case '-':
			sortOrder = "desc"
			sort = sort[1:]
		case '+':
			sort = sort[1:]
		default:

		}
	}
	return GetListWithPageEx(table, vList, pageSize, offset, sort, sortOrder, where, "")
}

func GetListWithPageEx(table InTable, vList interface{}, pageSize int, offset int, sort string, sortOrder string, where string, selectField string) (page Page) {
	// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。// 第二个返回值是错误对象，在这里略过
	db := getDb(table.GetDbNameSpace())
	db = db.Model(table).Debug().Where(where)
	db.Count(&page.Total)
	if selectField != "" {
		db = db.Select(selectField)
	}
	db = db.Offset(offset)
	if sort != "" {
		if sortOrder == "desc" {
			db = db.Order(fmt.Sprintf("%s %s", sort, sortOrder))
		}
	}
	db = db.Offset(offset).Limit(pageSize)
	db = db.Find(vList)
	page.Errcode = "0"
	page.Offset = int64(offset)
	page.Rows = vList
	page.PageSize = int64(pageSize)
	return page
}

func GetListWithPageJoins(table InTable, vList interface{}, pageSize int, offset int, sort string, sortOrder string, where string, selecte string, joins string) (page Page) {
	// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。// 第二个返回值是错误对象，在这里略过
	db := getDb(table.GetDbNameSpace())
	db = db.Table(table.TableName()).Debug().Where(where)
	db.Count(&page.Total)
	if selecte != "" {
		db = db.Select(selecte)
	}
	db = db.Offset(offset)
	if sort != "" {
		if sortOrder == "desc" {
			db = db.Order(fmt.Sprintf("%s %s", sort, sortOrder))
		}
	}
	db = db.Offset(offset).Limit(pageSize)
	if joins != "" {
		db = db.Joins(joins)
	}
	db = db.Find(vList)
	page.Errcode = "0"
	page.Offset = int64(offset)
	page.Rows = vList
	page.PageSize = int64(pageSize)
	return page
}

func Related(table InTable, list interface{}, foreignKeys ...string) (int64, error) {
	db := getDb(table.GetDbNameSpace()).Debug()
	db = db.Model(table).Related(list, foreignKeys...)
	return db.RowsAffected, db.Error
}

func UpdateByWhere(table InTable, selectField string, omitField string, where string, args map[string]interface{}) (int64, error) {
	parseMdParam(table, args)
	db := getDb(table.GetDbNameSpace())
	db = db.Model(table).Debug().Where(where)
	//db = db.Model(table).Debug().Where(where)
	if selectField != "" {
		db = db.Select(selectField)
	}
	if omitField != "" {
		db = db.Omit(omitField)
	}

	db = db.Updates(snakeField(args))
	return db.RowsAffected, db.Error
}

func Update(table InTable) (int64, error) {
	db := getDb(table.GetDbNameSpace())
	db = db.Model(table).Debug().Updates(table)
	return db.RowsAffected, db.Error
}
