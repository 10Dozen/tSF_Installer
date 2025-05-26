package application

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// https://github.com/10Dozen/tSF_Installer/blob/master/go/cmd/install.go

func downloadComponent(name, url string) string {
	// Downloads repo of given 'name' and given 'url'.
	// URL may be modified by CLI params, by default master branches will be used.
	// Returns:
	//   string - temp directory of the downloaded and uncompressed repo.

	fqUrl := getDownloadUrl(name, url)
	log.Printf("Repo [%s], fqUrl = %s\n", name, fqUrl)

	err := downloadFile(name, fqUrl)
	if err != nil {
		panic(err)
	}
	dir, err := os.MkdirTemp("", "temp")
	if err != nil {
		panic(err)
	}

	err = unzipFile(name, dir)
	if err != nil {
		panic(err)
	}

	os.Remove(name)

	return dir
}

func getDownloadUrl(name, url string) string {
	// Pre-format url with protocol, branches and path to actual file at GitHub
	// Basically converts GitHub web url of the branch into link a branch archive file
	// Return:
	//   string URL to selected repository/branch archive file
	fqUrl := completeUrl(url)

	//if fqUrl == completeUrl(DefaultOptions[name].Value) {
	//	fqUrl += fmt.Sprintf(GITHUB_DOWNLOAD_SUFFIX, GITHUB_DEFAULT_BRANCH)
	//	return fqUrl
	//}

	// Non-default repo path
	urlParts := strings.Split(fqUrl, "/")
	branchName := urlParts[6]

	fmt.Println(urlParts)
	fmt.Println(branchName)
	fmt.Println(strings.Join(urlParts[0:5], "/"))
	fmt.Println(fmt.Sprintf(GITHUB_DOWNLOAD_SUFFIX, branchName))

	return strings.Join(urlParts[0:5], "/") + "/" + fmt.Sprintf(GITHUB_DOWNLOAD_SUFFIX, branchName)
}

func completeUrl(url string) string {
	// Returns fully qualified URL (with protocol prefix and valid end char)
	fqUrl := url
	if !strings.HasPrefix(url, HTTPS_PREFIX) {
		fqUrl = HTTPS_PREFIX + fqUrl
	}
	lastChar := string(fqUrl[len(fqUrl)-1])
	if lastChar != "/" {
		fqUrl += "/"
	}
	return fqUrl
}

func downloadFile(filepath, url string) error {
	// Downloads content from given 'url' and writes to 'filepath' file on disk
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzipFile(source, dest string) error {
	// Uncompresses files/dirs from 'source' archive and puts them to 'dest' directory
	read, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer read.Close()

	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		name := path.Join(dest, file.Name)
		os.MkdirAll(path.Dir(name), os.ModeDir)
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
	}
	return nil
}

func zipDirectory(sourceDir string) error {
	zipFileName := filepath.Base(sourceDir) + ".BACKUP.zip"
	zipFilePath := filepath.Join(sourceDir, zipFileName)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Обходим все файлы и поддиректории
	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем сам созданный zip-файл
		if filePath == zipFilePath {
			return nil
		}

		// Получаем относительный путь файла от корневой директории
		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		// Если это директория, создаём её в архиве
		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		// Открываем файл для добавления в архив
		fileToZip, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		// Создаём файл в архиве
		zipFileEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Копируем содержимое файла в архив
		_, err = io.Copy(zipFileEntry, fileToZip)
		return err
	})

	return err
}

func composeComponentAtTempDirectory(source, dest string) {
	// Copies content of uncompressed file at 'source' to 'dest'
	// Files are stored in directory with archive name,
	// so it's neccessary to dive into 'source' directory and copy
	// from there.

	entries, err := os.ReadDir(source)
	if err != nil {
		panic(err)
	}

	var subDirName string
	for _, entry := range entries {
		if entry.IsDir() {
			subDirName = entry.Name()
			break
		}
	}

	if subDirName == "" {
		fmt.Println("Failed to copy", source)
	}

	err = copyDir(filepath.Join(source, subDirName), dest)
	if err != nil {
		panic(err)
	}
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("non-regular file: %s\n", src)
	}

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, sourceFileStat.Mode())
}

func PathExists(src string) bool {
	_, err := os.Stat(src)
	return !errors.Is(err, os.ErrNotExist)
}
