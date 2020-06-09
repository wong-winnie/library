	if path, err := helppath.GetWorkDir(); err != nil {
		log.Fatal(err.Error())
	} else {
		if err := helpconf.ViperReadInConfig("app", path+"/conf"); err != nil {
			log.Fatal(err.Error())
		}
	}

	for {
		fmt.Println(viper.GetString("pulicKey"))
		fmt.Println(viper.GetIntSlice("pulicKey"))
		time.Sleep(time.Second * 3)
	}
