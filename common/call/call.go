package call


func Retry(retries int, f func() error) error {
	err := f()
	if err == nil {
		return nil
	}
	if retries <=0 {
		return err
	}
	for i := 0; i< retries ; i++ {
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}