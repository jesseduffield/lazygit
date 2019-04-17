package collections

type (
	Collection interface {
		Do(func(interface{})bool)
	}
)

func GetRange(c Collection, start, length int) []interface{} {
	end := start + length
	items := make([]interface{}, length)
	i := 0
	j := 0
	c.Do(func(item interface{})bool{
		if i >= start {
			if i < end {
				items[j] = item
				j++
			} else {
				return false
			}
		}
		i++
		return true
	})
	return items[:j]
}
