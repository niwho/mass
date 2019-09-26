package utils

import (
	"bufio"
	"strings"

	"github.com/colinmarc/hdfs"
)

type Hadoop struct {
	client *hdfs.Client
}

func (hp *Hadoop) New(address string) error {
	var err error
	hp.client, err = hdfs.New(address)
	return err
}

// 协程安全问题
func (hp *Hadoop) Close() error {
	if hp.client != nil {
		return hp.client.Close()
	}
	return nil
}

func (hp *Hadoop) ReadTextFile(path string) ([]string, error) {
	var lines []string
	dat, err := hp.client.ReadDir(path)
	if err != nil {
		return lines, err
	}
	for _, fi := range dat {
		if !fi.IsDir() && strings.HasPrefix(fi.Name(), "part-") {
			part, err := hp.client.ReadFile(path + "/" + fi.Name())
			if err == nil {

				for {
					if ad, line, err := bufio.ScanLines(part, true); err == nil {
						if line == nil {
							break
						}
						part = part[ad:]
						lines = append(lines, string(line))

					} else {
						break
					}
				}
			} else {
				return lines, err
			}
		}
	}
	return lines, nil
}
