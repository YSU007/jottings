package main

import "fmt"

// Component 接口定义了组合对象和叶子对象的公共方法
type Component interface {
	Search(keyword string)
}

// File 结构体实现了 Component 接口，代表叶子节点
type File struct {
	name string
}

// Search 方法实现了在 File 中搜索关键词的功能
func (f *File) Search(keyword string) {
	fmt.Printf("Searching for keyword %s in file %s\n", keyword, f.name)
}

// Folder 结构体也实现了 Component 接口，代表组合对象
type Folder struct {
	components []Component
	name       string
}

// Search 方法实现了在 Folder 中递归搜索关键词的功能
func (f *Folder) Search(keyword string) {
	fmt.Printf("Searching recursively for keyword %s in folder %s\n", keyword, f.name)
	for _, composite := range f.components {
		composite.Search(keyword)
	}
}

// Add 方法用于向 Folder 添加子组件
func (f *Folder) Add(c Component) {
	f.components = append(f.components, c)
}

func main() {
	file1 := &File{name: "File1"}
	file2 := &File{name: "File2"}
	file3 := &File{name: "File3"}

	folder1 := &Folder{name: "Folder1"}
	folder1.Add(file1)

	folder2 := &Folder{name: "Folder2"}
	folder2.Add(file2)
	folder2.Add(file3)
	folder2.Add(folder1)

	folder2.Search("rose")
}
