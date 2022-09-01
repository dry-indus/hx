/**
 * @Author: lzw5399
 * @Date: 2020/8/3 22:41
 * @Desc: 模型映射 支持 struct->struct struct->map等等
 */
package util

import "github.com/jinzhu/copier"

func MapTo(from, to interface{}) error {
	return copier.Copy(to, from)
}
