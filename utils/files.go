package utils

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	doublestar "github.com/bmatcuk/doublestar/v4"
)

func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func Copy(from, to string) error {
	info, err := os.Stat(from)
	if err != nil {
		return err
	}

	if info.IsDir() {
		if err := os.MkdirAll(to, os.ModePerm); err != nil {
			return err
		}

		dir, err := os.Open(from)
		if err != nil {
			return err
		}
		defer dir.Close()

		// Read source directory contents
		fileInfos, err := dir.Readdir(-1)
		if err != nil {
			return err
		}

		for _, fileInfo := range fileInfos {
			srcPath := filepath.Join(from, fileInfo.Name())
			destPath := filepath.Join(to, fileInfo.Name())

			if fileInfo.IsDir() {
				// Copy subdirectory recursively
				if err := Copy(srcPath, destPath); err != nil {
					return err
				}
			} else {
				// Copy file
				if err := CopyFile(srcPath, destPath); err != nil {
					return err
				}
			}
		}
	} else {
		return CopyFile(from, to)
	}

	return nil
}

func CopyFile(from, to string) (err error) {
	if err = os.MkdirAll(filepath.Dir(to), os.ModePerm); err != nil {
		return fmt.Errorf("creating dest dir %s: %w", filepath.Dir(to), err)
	}

	src, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("opening src %s: %w", from, err)
	}
	defer src.Close()

	dest, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("opening dest %s: %w", to, err)
	}
	defer func() {
		if e := dest.Close(); e != nil {
			err = e
		}
	}()

	if _, err := io.Copy(dest, src); os.IsNotExist(err) {
		return fmt.Errorf("path \"%s\" does not exist", from)
	} else if err != nil {
		return fmt.Errorf("copying into cache: %w", err)
	}

	si, err := os.Stat(from)
	if err != nil {
		return fmt.Errorf("retrieving src permissions: %w", err)
	}
	err = os.Chmod(to, si.Mode())
	if err != nil {
		return fmt.Errorf("changing dest permissions: %w", err)
	}

	return err
}

// Helper function to read the exclusion file and return a list of excluded paths
func ReadExclusionFile(exclusionFile string) ([]string, error) {
	exclusions := []string{}

	file, err := os.Open(exclusionFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		exclusion := strings.TrimSpace(scanner.Text())
		if exclusion != "" {
			exclusions = append(exclusions, exclusion)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return exclusions, nil
}

// AbsoluteFilePath will check if target is an absolute path and if not, join it with root
// and return an absolute filepath
func AbsoluteFilePath(root, target string) (fullpath string) {
	if strings.HasPrefix(target, "/") {
		fullpath = target
	} else {
		fullpath, _ = filepath.Abs(filepath.Join(root, target))
	}

	return
}

func BringInsideRoot(root, path string) string {
	var relPath string
	if !strings.HasPrefix(path, "/") {
		path = filepath.Join(root, path)
	}

	path, _ = filepath.Abs(path)
	relPath, _ = filepath.Rel(root, path)
	for strings.HasPrefix(relPath, "..") {
		relPath = strings.TrimPrefix(relPath, "../")
	}
	return relPath
}

func GlobPath(root, path string) (map[string]string, error) {
	expanded := map[string]string{}
	fullpath := filepath.Join(root, path)

	fsys, patt := doublestar.SplitPattern(fullpath)

	nonVariablePath := strings.TrimPrefix(fsys, root)
	if err := doublestar.GlobWalk(os.DirFS(fsys), patt, func(p string, d fs.DirEntry) error {
		p = filepath.Join(nonVariablePath, p)
		fileInfo, _ := os.Stat(filepath.Join(root, p))

		if fileInfo.IsDir() {
			return nil
		}

		expanded[p] = filepath.Join(root, p)

		return nil
	}, doublestar.WithFailOnIOErrors()); err != nil {
		return nil, fmt.Errorf("glob walking: %w", err)
	}

	return expanded, nil
}
