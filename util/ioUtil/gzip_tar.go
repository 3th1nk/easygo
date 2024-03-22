package ioUtil

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GZipTarToFile gzip 压缩打包到文件
//		filePath 打包的文件/文件夹
//		targetPath 打包后的文件路径
//		removeSrcFile 打包成功后删除打包的文件/文件夹 true.删除
func GZipTarToFile(filePath string, targetPath string, removeSrcFile ...bool) error {
	if IsFileExists(targetPath) {
		return os.ErrExist
	}

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	return GZipTar(filePath, targetFile, removeSrcFile...)
}

// GZipTar gzip 压缩打包
//		filename 打包的文件/文件夹
//		w 写入的文件流
//		removeSrcFile 打包成功后删除打包的文件/文件夹 true.删除
func GZipTar(filePath string, w io.Writer, removeSrcFile ...bool) (err error) {
	gw := gzip.NewWriter(w)
	tw := tar.NewWriter(gw)
	defer func() {
		tw.Close()
		gw.Close()

		// 打包成功，删除源文件/文件夹
		if err == nil && len(removeSrcFile) > 0 && removeSrcFile[0] {
			os.RemoveAll(filePath)
		}
	}()

	fi, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		err = tarDir(filePath, tw)
	} else {
		err = tarFile(filepath.Base(filePath), filePath, fi, tw)
	}
	return
}

// tarDir 打包文件夹
func tarDir(dir string, w *tar.Writer) error {
	// 打包进去的文件夹
	targetDir := filepath.Base(dir)

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}

		// 目标文件
		targetName := filepath.Join(targetDir, strings.TrimPrefix(path, dir))

		// 文件
		if !info.IsDir() {
			return tarFile(targetName, path, info, w)
		}

		// 文件夹
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Format = tar.FormatGNU
		header.Name = targetName
		if err = w.WriteHeader(header); err != nil {
			return err
		}
		return nil
	})
}

func tarFile(targetName, filePath string, fileInfo os.FileInfo, w *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}
	header.Format = tar.FormatGNU
	header.Name = targetName
	header.Mode = int64(fileInfo.Mode())

	if err = w.WriteHeader(header); err != nil {
		return err
	}

	if _, err = io.Copy(w, file); err != nil {
		return err
	}
	return nil
}

// UnGZipTarToDir gzip解压文件流到目录
//		r gzip压缩文件流
//		dir 解压到文件夹
func UnGZipTarToDir(r io.Reader, dir string) error {
	if err := MakeDir(dir); err != nil {
		return err
	}

	return UnGZipTarF(r, func(h *tar.Header, r io.Reader) error {
		if h.FileInfo().IsDir() {
			return MakeDir(filepath.Join(dir, h.Name))
		}

		// 打开文件
		file, err := os.OpenFile(filepath.Join(dir, h.Name), os.O_CREATE|os.O_WRONLY, os.FileMode(h.Mode))
		if err != nil {
			return err
		}
		defer file.Close()

		// 写文件
		_, err = io.Copy(file, r)
		return err
	})

}

// UnGZipTarF gzip 解压文件流
//		r gzip压缩文件流
//		f 回调解压后的文件流 如果是文件夹，r 为nil
func UnGZipTarF(r io.Reader, f func(h *tar.Header, r io.Reader) error) error {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	// 读取文件
	for {
		if h, err := tr.Next(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		} else {
			if h.FileInfo().IsDir() {
				err = f(h, nil)
			} else {
				err = f(h, tr)
			}
			if err != nil {
				return err
			}
		}
	}
}

// UnGZipTarFileF gzip 解压文件流
//		filePath gzip压缩文件
//		f 回调解压后的文件流 如果是文件夹，r 为nil
//		removeSrcFile 解压成功后删除压缩文件 true.删除
func UnGZipTarFileF(filePath string, f func(h *tar.Header, r io.Reader) error, removeSrcFile ...bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = UnGZipTarF(file, f)

	// 删除压缩文件
	if err == nil || len(removeSrcFile) > 0 && removeSrcFile[0] {
		os.RemoveAll(filePath)
	}
	return err
}

// UnGZipTarFileToDir gzip解压文件到目录
//		filePath gzip压缩文件
//		dir 解压到文件夹
//		removeSrcFile 解压成功后删除压缩文件 true.删除
func UnGZipTarFileToDir(filePath string, dir string, removeSrcFile ...bool) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = UnGZipTarToDir(file, dir)

	// 删除压缩文件
	if err == nil && len(removeSrcFile) > 0 && removeSrcFile[0] {
		os.RemoveAll(filePath)
	}
	return err
}
