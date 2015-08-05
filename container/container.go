package container

// Container is a value object representing a running container
type Container struct {
	ID   string
	Name string
	Port int64
	Tags []string
}

// IsEqual is used for verifying that containers are the same
func (c Container) IsEqual(other Container) bool {
	return c.ID == other.ID &&
		c.Name == other.Name &&
		c.Port == other.Port &&
		sameTags(c.Tags, other.Tags)
}

func sameTags(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}

	return true
}
