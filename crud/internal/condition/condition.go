package condition

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Condition string

func NewCondition(k string) Condition {
	for _, act := range DefaultActions {
		if strings.EqualFold(k, string(act)) {
			return act
		}
	}
	return Eq
}

func (c Condition) Action() func(db *gorm.DB, k string, v interface{}) *gorm.DB {
	switch c {
	case Eq:
		return EqAct
	case Ne:
		return NeAct
	case Gt:
		return GtAct
	case Gte:
		return GteAct
	case Lt:
		return LtAct
	case Lte:
		return LteAct
	case In:
		return InAct
	case NotIn:
		return NotInAct
	case Like:
		return LikeAct
	case NotLike:
		return NotLikeAct
	case LikeRight:
		return LikeRightAct
	case LikeLeft:
		return LikeLeftAct
	case IsNull:
		return IsNullAct
	case IsNotNull:
		return IsNotNullAct
	default:
		return EqAct
	}
}

var (
	Eq        = Condition("eq")
	Ne        = Condition("ne")
	Gt        = Condition("gt")
	Gte       = Condition("gte")
	Lt        = Condition("lt")
	Lte       = Condition("lte")
	In        = Condition("in")
	NotIn     = Condition("not_in")
	Like      = Condition("like")
	NotLike   = Condition("not_like")
	LikeRight = Condition("like_right")
	LikeLeft  = Condition("like_left")
	FK        = Condition("fk")
	IsNull    = Condition("is_null")
	IsNotNull = Condition("is_notnull")
	Sort      = Condition("sort")

	DefaultActions = []Condition{Eq, Ne, Gt, Gte, Lt, Lte, In, NotIn, Like, NotLike, LikeRight, LikeLeft, FK, IsNull, IsNotNull, Sort}
)

var (
	EqAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" = ?", v)
	}
	NeAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" <> ?", v)
	}
	GtAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" > ?", v)
	}
	GteAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" >= ?", v)
	}
	LtAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" < ?", v)
	}
	LteAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" <= ?", v)
	}
	InAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" IN ?", v)
	}
	NotInAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" NOT IN ?", v)
	}
	LikeAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" LIKE ?", "%"+v.(string)+"%")
	}
	NotLikeAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" NOT LIKE ?", "%"+v.(string)+"%")
	}
	LikeRightAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" LIKE ?", v.(string)+"%")
	}
	LikeLeftAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k+" LIKE ?", "%"+v.(string))
	}
	/**
	 * @param fk 关联表外键 author_id
	 * @param table 关联表名 author
	 * @param k 关联表主键	id
	 * @param v 关联表主键值 ?
	 */
	FKAct = func(db *gorm.DB, fk string, table string, k string, v interface{}) *gorm.DB {
		format := "%s.%s = %s.%s_%s"
		query := fmt.Sprintf(format, table, k, table, table, fk)
		return db.Joins(query).Where(table+"."+k+" = ?", v)
	}
	IsNullAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k + " IS NULL")
	}
	IsNotNullAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Where(k + " IS NOT NULL")
	}
	SortAct = func(db *gorm.DB, k string, v interface{}) *gorm.DB {
		return db.Order(k + " " + v.(string))
	}
)
