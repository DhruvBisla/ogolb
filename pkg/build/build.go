package build

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}
	return
}

func allFiles() ([]string, error) {
	fileinfo, err := ioutil.ReadDir("./content")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, file := range fileinfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func newCustomParser() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
}

type Post struct {
	Title      string
	Identifier string
	Summary    string
	Date       string
	Body       template.HTML
}

type Posts struct {
	Posts []Post
}

func Build() {
	var i bytes.Buffer
	fmt.Println(allFiles())
	md := newCustomParser()
	files, err := allFiles()
	var ps Posts
	if err != nil {
		panic(err)
	}
	_, err = os.Stat("public")
	if !os.IsNotExist(err) {
		os.RemoveAll(filepath.Join(".", "public"))
	}
	err = os.Mkdir("public", 0755)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		var buf bytes.Buffer
		f, err := ioutil.ReadFile(filepath.Join(".", "content", file))
		if err != nil {
			panic(err)
		}
		context := parser.NewContext()
		if err := md.Convert(f, &buf, parser.WithContext(context)); err != nil {
			panic(err)
		}
		data := meta.Get(context)
		tmpl, err := template.ParseFiles(filepath.Join(".", "templates", "post.html"), filepath.Join(".", "templates", "header.html"), filepath.Join(".", "templates", "footer.html"))
		if err != nil {
			panic(err)
		}
		var t bytes.Buffer
		var p Post = Post{
			Title:      fmt.Sprintf("%v", data["Title"]),
			Identifier: fmt.Sprintf("%v", data["Identifier"]),
			Summary:    fmt.Sprintf("%v", data["Summary"]),
			Date:       fmt.Sprintf("%v", data["Date"]),
			Body:       template.HTML(buf.String()),
		}
		ps.Posts = append(ps.Posts, p)
		tmpl.ExecuteTemplate(&t, "post.html", p)
		err = ioutil.WriteFile(filepath.Join(".", "public", (p.Identifier+".html")), t.Bytes(), 0744)
		if err != nil {
			panic(err)
		}
	}
	ind, err := template.ParseFiles(filepath.Join(".", "templates", "index.html"), filepath.Join(".", "templates", "header.html"), filepath.Join(".", "templates", "footer.html"))
	if err != nil {
		panic(err)
	}
	err = ind.ExecuteTemplate(&i, "index.html", ps)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join(".", "public", "index.html"), i.Bytes(), 0744)
	if err != nil {
		panic(err)
	}
	if err := CopyDir(filepath.Join(".", "static"), filepath.Join(".", "public", "static")); err != nil {
		panic(err)
	}
}
