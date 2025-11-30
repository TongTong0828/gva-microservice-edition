package shop
import(
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type Product struct {
	global.GVA_MODEL
	Name    string `json:"name" gorm:"comment:商品名称"`
	Price   int    `json:"price" gorm:"comment:价格(分)"`
	Stock   int    `json:"stock" gorm:"comment:库存数量"`
	Version int    `json:"version" gorm:"default:0;comment:乐观锁版本号"` // 核心字段
}

func (Product) TableName() string {
	return "shop_products"
}