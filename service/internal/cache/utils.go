package cache

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
