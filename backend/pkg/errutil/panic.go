package errutil

/**
* 如果err不为空,panic
*
* @param err - 错误对象
 */
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
