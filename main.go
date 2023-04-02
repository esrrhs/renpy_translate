package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/bregydoc/gtranslate"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var input_path = flag.String("input", "", "input file path")

func main() {
	flag.Parse()

	if *input_path == "" {
		panic("input path is empty")
	}

	find_all_charater()
	find_all_key()
	save_all_key()
	translate_all_key()
	save_all_key()
	build_key_map()
	replace_all_file()
}

var g_key_map = make(map[string]string)

func build_key_map() {
	for _, v := range g_translate_infos.Keys {
		g_key_map[v.Src] = v.Dst
	}
}

func replace_all_file() {
	filepath.Walk(*input_path, func(path string, info os.FileInfo, err error) error {
		// not .rpy file?
		if filepath.Ext(path) != ".rpy" {
			return nil
		}

		replace_file(path)

		return nil
	})
}

func replace_file(path string) {

	// first copy src file to origin file
	origin_path := path + ".origin"

	// check if origin file exist
	if _, err := os.Stat(origin_path); os.IsNotExist(err) {
		// origin file not exist, copy src file to origin file
		input, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(origin_path, input, 0644)
		if err != nil {
			panic(err)
		}
	}

	// open origin file
	file, err := os.Open(origin_path)
	if err != nil {
		log.Printf("open file error: %s", err)
		panic(err)
	}
	defer file.Close()

	// open dst file
	dst_file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("open file error: %s", err)
		panic(err)
	}
	defer dst_file.Close()

	// read origin file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// check is match regex "  xx = \"yy\"  "
		r, err := regexp.Compile(`^\s*(\w+)\s+\"(.*)\"\s*$`)
		if err != nil {
			panic(err)
		}

		if r.MatchString(line) {
			names := r.FindAllStringSubmatch(line, -1)
			if len(names) > 0 {
				c := names[0][1]
				words := names[0][2]
				if _, ok := g_charater_map[c]; ok {
					if _, ok := g_key_map[words]; ok {
						dst_words := g_key_map[words]
						line = strings.Replace(line, words, dst_words, -1)
						log.Printf("replace key: %s -> %s", words, dst_words)
					}
				}
			}
		}

		// write to dst file
		dst_file.WriteString(line + "\n")
	}

	// delete .rpyc file
	rpyc_path := strings.Replace(path, ".rpy", ".rpyc", -1)
	os.Remove(rpyc_path)
}

func translate_all_key() {
	for index := range g_translate_infos.Keys {
		info := &g_translate_infos.Keys[index]
		info.Dst = translate_key(info.Src)
		log.Printf("translate key: %s -> %s", info.Src, info.Dst)
	}
}

func translate_key(src string) string {
	translated, err := gtranslate.TranslateWithParams(
		src,
		gtranslate.TranslationParams{
			From: "en",
			To:   "zh",
		},
	)

	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 100)
	return translated
}

func find_all_charater() {
	filepath.Walk(*input_path, func(path string, info os.FileInfo, err error) error {
		// not .rpy file?
		if filepath.Ext(path) != ".rpy" {
			return nil
		}

		find_file_charater(path)

		return nil
	})
}

var g_charater_map = make(map[string]int)
var g_charater_len = 0

func find_file_charater(path string) {

	// check origin file is exist
	origin_path := path + ".origin"
	if _, err := os.Stat(origin_path); os.IsNotExist(err) {
	} else {
		path = origin_path
	}

	// open file
	file, err := os.Open(path)
	if err != nil {
		log.Printf("open file error: %s", err)
		panic(err)
	}
	defer file.Close()

	// read file line by line
	scanner := bufio.NewScanner(file)
	lineno := 1
	for scanner.Scan() {
		line := scanner.Text()

		// check is match regex "define girl = Character("Girl", color="a9f2b4")" or "define girl = DynamicCharacter("Girl", color="a9f2b4")"
		r, err := regexp.Compile(`^\s*define\s+(\w+)\s*=\s*(Character|DynamicCharacter)\(.*\)\s*$`)
		if err != nil {
			panic(err)
		}

		if r.MatchString(line) {
			names := r.FindAllStringSubmatch(line, -1)
			if len(names) > 0 {
				c := names[0][1]
				g_charater_len++
				g_charater_map[c] = g_charater_len
				log.Printf("find charater: %s %d at line %d", c, g_charater_map[c], lineno)
			}
		}

		lineno++
	}
}

type TranslateInfo struct {
	Character string `json:"character"`
	Src       string `json:"src"`
	Dst       string `json:"dst"`
}

type TranslateInfos struct {
	Keys []TranslateInfo `json:"keys"`
}

var g_translate_infos = TranslateInfos{}

func find_all_key() {
	filepath.Walk(*input_path, func(path string, info os.FileInfo, err error) error {
		// not .rpy file?
		if filepath.Ext(path) != ".rpy" {
			return nil
		}

		find_file_key(path)

		return nil
	})
}

func save_all_key() {
	// write to json file
	file, err := os.Create("translate.json")
	if err != nil {
		log.Printf("create file error: %s", err)
		panic(err)
	}
	defer file.Close()

	// write json
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(g_translate_infos)
	if err != nil {
		log.Printf("write json error: %s", err)
		panic(err)
	}
}

func find_file_key(path string) {

	// check origin file exist
	origin_path := path + ".origin"
	if _, err := os.Stat(origin_path); os.IsNotExist(err) {
	} else {
		path = origin_path
	}

	// open file
	file, err := os.Open(path)
	if err != nil {
		log.Printf("open file error: %s", err)
		panic(err)
	}
	defer file.Close()

	// read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//log.Printf("line: %s", line)

		// check is match regex "  xx = \"yy\"  "
		r, err := regexp.Compile(`^\s*(\w+)\s+\"(.*)\"\s*$`)
		if err != nil {
			panic(err)
		}

		if r.MatchString(line) {
			names := r.FindAllStringSubmatch(line, -1)
			if len(names) > 0 {
				c := names[0][1]
				words := names[0][2]
				if _, ok := g_charater_map[c]; ok {
					log.Printf("find name: %s words: %s", c, words)
					g_translate_infos.Keys = append(g_translate_infos.Keys, TranslateInfo{Character: c, Src: words})
				}
			}
		}
	}
}
