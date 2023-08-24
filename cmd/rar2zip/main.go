package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/ishidaworks/rar2zip/internal/archive"
	"github.com/ishidaworks/rar2zip/internal/fileutils"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: rar2zip <search path> [-list] [-convert] [-convert-delete]")
		os.Exit(1)
	}

	searchPathArgs := os.Args[1]
	if _, err := os.Stat(searchPathArgs); err != nil {
		fmt.Printf("Error: %s\n", "search path is not found. please check path.")
		os.Exit(1)
	}
	mode := 0
	if len(os.Args) > 2 {
		modeArgs := os.Args[2]
		if modeArgs == "-list" {
			mode = 0
		} else if modeArgs == "-convert" {
			mode = 1
		} else if modeArgs == "-convert-delete" {
			mode = 2
		} else {
			fmt.Println("Usage: rar2zip <search path> [-list] [-convert] [-convert-delete]")
			os.Exit(1)
		}
	}

	searchPath, err := filepath.Abs(searchPathArgs)
	if err != nil {
		panic(err)
	}

	// RARファイル検索
	rarPathList, err := fileutils.SearchArcives(searchPath, "*.rar")
	if err != nil {
		panic(err)
	}
	if len(rarPathList) == 0 {
		fmt.Printf("%s\n", "rar file is not found. please check search path.")
		os.Exit(0)
	}

	for _, rarPath := range rarPathList {
		fmt.Printf("%s\n", rarPath)
		uuidObj, _ := uuid.NewRandom()
		fileName := fileutils.GetFileNameWithoutExt(rarPath)
		baseDir := filepath.Dir(rarPath)
		unRarDistPath := fmt.Sprintf("%s/%s", baseDir, uuidObj.String())
		zipDistPath := fmt.Sprintf("%s/%s.zip", baseDir, fileName)

		if _, err := os.Stat(zipDistPath); err == nil {
			fmt.Println("\t=> zip already exists")
			continue
		}
		if mode == 0 {
			continue
		}

		// RAR展開
		err := archive.UnRar(rarPath, unRarDistPath)
		if err != nil {
			panic(err)
		}
		// ZIP圧縮
		if err := archive.CompressZip(unRarDistPath, zipDistPath); err != nil {
			panic(err)
		}
		// RAR展開フォルダ削除
		if err := os.RemoveAll(unRarDistPath); err != nil {
			panic(err)
		}
		// RARファイル削除
		if mode == 2 {
			if err := os.Remove(rarPath); err != nil {
				panic(err)
			}
		}
		fmt.Println("\t=> success")
	}
}
