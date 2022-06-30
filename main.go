package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type licenseType struct {
	start string
	size  string
	names []string
}

type licensesType = []licenseType

func main() {
	licenses := licensesType{}

	/**********************************
	* ファイル読み込みと解析
	***********************************/
	file_metadata, err_metadata := os.Open("third_party_license_metadata")
	if err_metadata != nil {
		log.Fatalf("Error when opening file: %s", err_metadata)
	}

	// 1行ずつ読み込み、解析
	fileScanner := bufio.NewScanner(file_metadata)
	for fileScanner.Scan() {
		arr1 := strings.Split(fileScanner.Text(), ":")
		arr2 := strings.SplitN(arr1[1], " ", 2)
		// start→arr1[0]  size→arr2[0]  name→arr2[1]
		fmt.Println("start=" + arr1[0] + ",size=" + arr2[0] + ",name=" + arr2[1])

		isFound := false
		for i := 0; i < len(licenses); i++ {
			//すでに同じOSSがある場合はライブラリ名を追加
			if licenses[i].start == arr1[0] {
				isFound = true
				//既に登録済のライセンス名がなければ登録
				if !isContains(licenses[i].names, arr2[1]) {
					licenses[i].names = append(licenses[i].names, arr2[1])
				}
				break
			}
		}
		//新規OSS追加
		if !isFound {
			license := licenseType{}
			license.start = arr1[0]
			license.size = arr2[0]
			license.names = append(license.names, arr2[1])
			licenses = append(licenses, license)
		}
	}

	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file_metadata.Close()

	/**********************************
	* テキストファイル書き込み
	***********************************/
	file_licenses, err_licenses := os.Open("third_party_licenses")
	if err_licenses != nil {
		log.Fatalf("Error when opening file: %s", err_licenses)
	}

	file_w_text, err_w_text := os.Create("oss-app-licenses.txt")
	if err_w_text != nil {
		log.Fatalf("Error when opening file: %s", err_w_text)
	}

	for i := 0; i < len(licenses); i++ {
		data := []byte("******************************************************************\n")
		file_w_text.Write(data)

		for j := 0; j < len(licenses[i].names); j++ {
			//fmt.Println(licenses[i].names[j])
			data = []byte(licenses[i].names[j] + "\n")
			file_w_text.Write(data)
		}

		data = []byte("\n")
		file_w_text.Write(data)

		var tmp int
		tmp, _ = strconv.Atoi(licenses[i].start)
		var start_num int64 = int64(tmp)
		tmp, _ = strconv.Atoi(licenses[i].size)
		var size_num int64 = int64(tmp)
		//読み込む位置へ移動
		_, err_licenses2 := file_licenses.Seek(start_num, 0)
		if err_licenses2 != nil {
			log.Fatalf("Error when opening file: %s", err_licenses2)
		}
		//読み込むバイト数だけメモリ確保
		read_data := make([]byte, size_num)
		//読み込み
		_, err_licenses3 := file_licenses.Read(read_data)
		if err_licenses3 != nil {
			log.Fatalf("Error when opening file: %s", err_licenses3)
		}

		data = []byte(string(read_data) + "\n")
		file_w_text.Write(data)
	}

	file_licenses.Close()
	file_w_text.Close()

	/**********************************
	* テキストファイル内の改行と<>をhtmlに変換
	***********************************/
	file_o_text, err_text := os.Open("oss-app-licenses.txt")
	if err_text != nil {
		log.Fatalf("Error when opening file: %s", err_text)
	}

	file_w_html, err_w_html := os.Create("oss-app-licenses.html")
	if err_w_html != nil {
		log.Fatalf("Error when opening file: %s", err_w_html)
	}

	br := bufio.NewReader(file_o_text)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		replaced := strings.Replace(string(line), `<`, `&#12296;`, -1)
		replaced = strings.Replace(string(replaced), `>`, `&#12297;`, -1)
		replaced += "<br />\n"
		data := []byte(string(replaced))
		file_w_html.Write(data)
	}

	file_o_text.Close()
	file_w_html.Close()
}

//配列の中に特定の文字列が含まれるかを返す
func isContains(arr []string, str string) bool {
	for _, val := range arr {
		if val == str {
			return true
		}
	}
	return false
}
