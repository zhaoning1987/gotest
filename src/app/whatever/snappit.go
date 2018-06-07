dirList, err := ioutil.ReadDir(dirPath)
if err != nil {
	fmt.Println(err)
	return
}
for _, file := range dirList {
	if !file.IsDir() {
		fmt.Println(filepath.Join(dirPath, file.Name()))
	}
}