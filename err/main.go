package err

func ArrDo(functions []func() error) error {
	for _, function := range functions {
		err := function()
		if err != nil {
			return err
		}
	}
	return nil
}
