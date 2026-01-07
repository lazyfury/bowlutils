package crud

import (
	"reflect"
	"strings"

	"github.com/lazyfury/bowlutils/crud/internal/condition"
	"gorm.io/gorm"
)

type TableName interface {
	TableName() string
}
type Model interface {
	GetID() uint
	TableName
	DeletedAtKey() string
}

type Repository[T Model] struct {
	db    *gorm.DB
	model T
}

func NewRepository[T Model](model T, db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db:    db,
		model: model,
	}
}

func (r *Repository[T]) FindByID(id uint) (T, error) {
	var model T
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return model, err
	}
	return model, nil
}

// query
func (r *Repository[T]) Query(kvs map[string]interface{}) *gorm.DB {
	return r.db.Table(r.model.TableName()).Where(kvs)
}

// db
func (r *Repository[T]) DB() *gorm.DB {
	return r.db.Table(r.model.TableName())
}

// tx
func (r *Repository[T]) Tx(fn func(db *gorm.DB) error) error {
	return r.db.Table(r.model.TableName()).Transaction(fn)
}

type QueryFunc func(db *gorm.DB) *gorm.DB

// list by deleted_at
func (r *Repository[T]) List(out any, opts ...QueryFunc) error {
	db := r.db.Table(r.model.TableName())
	for _, opt := range opts {
		db = opt(db)
	}
	if err := db.Find(out).Error; err != nil {
		return err
	}
	return nil
}

// page
func (r *Repository[T]) Page(out any, page, pageSize int, opts ...QueryFunc) (Page[T], error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	db := r.db.Table(r.model.TableName())
	for _, opt := range opts {
		db = opt(db)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return Page[T]{}, err
	}
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(out).Error; err != nil {
		return Page[T]{}, err
	}

	return Page[T]{
		PageNum:   int64(page),
		PageSize:  int64(pageSize),
		PageCount: (total + int64(pageSize) - 1) / int64(pageSize),
		Total:     total,
		Items:     (any)(out).(*[]T),
	}, nil
}

// exists
func (r *Repository[T]) Exists(id uint) (bool, error) {
	var model = r.model
	var count int64
	if err := r.db.Model(&model).Where("id = ?", id).Where(model.DeletedAtKey(), "IS NULL").Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// asset exists
func (r *Repository[T]) AssetExists(id uint) error {
	exists, err := r.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// create
func (r *Repository[T]) Create(model T) error {
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	return nil
}

// updates
func (r *Repository[T]) Updates(model T) error {
	if err := r.AssetExists(model.GetID()); err != nil {
		return err
	}
	if err := r.db.Updates(&model).Error; err != nil {
		return err
	}
	return nil
}

// update
func (r *Repository[T]) Update(key string, value interface{}) error {
	var model T
	if err := r.AssetExists(model.GetID()); err != nil {
		return err
	}
	if err := r.db.Model(&model).Where("id = ?", model.GetID()).Update(key, value).Error; err != nil {
		return err
	}
	return nil
}

// save
func (r *Repository[T]) Save(model T) error {
	if err := r.AssetExists(model.GetID()); err != nil {
		return err
	}
	if err := r.db.Save(&model).Error; err != nil {
		return err
	}
	return nil
}

// delete by id (soft delete if model has DeletedAt)
func (r *Repository[T]) DeleteByID(id uint) error {
	m := r.model
	if err := r.db.Table(m.TableName()).Where("id = ?", id).Delete(&m).Error; err != nil {
		return err
	}
	return nil
}

// reflect all model field key
func (r *Repository[T]) ReflectKeys() []string {
	var keys []string
	rType := reflect.TypeOf(r.model).Elem()
	fieldNum := rType.NumField()
	for i := 0; i < fieldNum; i++ {
		field := rType.Field(i)

		// if baseModel, skip
		if field.Type == reflect.TypeOf(&BaseModel{}) {
			brType := field.Type.Elem()
			fieldNum := brType.NumField()
			for j := 0; j < fieldNum; j++ {
				bf := brType.Field(j)
				if bf.Anonymous {
					continue
				}
				jsonTag := bf.Tag.Get("json")
				if jsonTag == "" || jsonTag == "-" {
					continue
				}
				keys = append(keys, jsonTag)
			}
			continue
		}

		if field.Anonymous {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		keys = append(keys, jsonTag)
	}
	return keys
}

// is valid key
func (r *Repository[T]) IsValidKey(key string) bool {
	keys := r.ReflectKeys()
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

// mapStringToMapInterface
func (r *Repository[T]) MapStringToMapInterface(params map[string]string) map[string]interface{} {
	var m map[string]interface{} = make(map[string]interface{})
	for k, v := range params {
		var val interface{} = v
		// test 1,2,3,4 transfrom to arr
		if strings.Contains(v, ",") {
			val = strings.Split(v, ",")
		} else {
			val = v
		}
		m[k] = val
	}
	return m
}

func (r *Repository[T]) QueryParamsToSearch(params map[string]string) []QueryFunc {
	return r.MapToSearch(r.MapStringToMapInterface(params))
}

// params map to list QueryFn
func (r *Repository[T]) MapToSearch(params map[string]interface{}) []QueryFunc {
	var fns []QueryFunc
	var keys = r.ReflectKeys()
	// logger.Debugf("keys: %v", keys)
	var isValid = func(key string) bool {
		for _, k := range keys {
			if k == key {
				return true
			}
		}
		return false
	}
	for k, v := range params {
		// logger.Debugf("key: %s, value: %v is valid: %v", k, v, isValid(k))
		// k 1 : name__action=?
		// k 2 : action_id__fk__eq=?
		// 特殊处理 : action 是 is_null is_notnull sort=asc/desc
		var key string
		var action condition.Condition
		// 解析 key action isFk fkAction
		if strings.Contains(k, "fk") {
			strs := strings.Split(k, "__fk__")[:2]
			key = strs[0]
			action = condition.NewCondition(strs[1])
			// logger.Attnf("key: %s, action: %s, isFk: %v", key, action, isFk)
		} else if strings.Contains(k, "__") {
			// split __
			strs := strings.Split(k, "__")[:2]
			key = strs[0]
			action = condition.NewCondition(strs[1])
		} else {
			key = k
			action = condition.Eq
		}
		// if fk
		if action == condition.FK && isValid(key) {
			var table, fKey string
			// split _
			strs := strings.Split(key, "_")
			table = strs[0]
			key = strings.Join(strs[1:], "_")
			fns = append(fns, func(db *gorm.DB) *gorm.DB {
				return condition.FKAct(db, fKey, table, key, v)
			})

			continue
		}

		if action == condition.Sort && isValid(key) {
			// 校验 sort 方向是否有效
			sortAction := v.(string)
			if sortAction != "asc" && sortAction != "desc" {
				continue
			}
			fns = append(fns, func(db *gorm.DB) *gorm.DB {
				return condition.SortAct(db, key, sortAction)
			})

			continue
		}
		// 校验 key 是否有效
		if isValid(key) {
			if action == condition.Sort {
				continue
			}
			fns = append(fns, func(db *gorm.DB) *gorm.DB {
				return action.Action()(db, key, v)
			})
		}
	}
	return fns
}
