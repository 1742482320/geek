/*

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.

*/

package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/grafov/m3u8"
)

// const VERSION = "1.0.5"

// var USER_AGENT string

// var client = &http.Client{}

func doRequest(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func downloadSegment(fn string, tsList []string) error {
	out, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer out.Close()
	for _, v := range tsList {
		req, err := http.NewRequest("GET", v, nil)
		if err != nil {
			return err
		}
		resp, err := doRequest(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			log.Printf("Received HTTP %v for %v\n", resp.StatusCode, v)
			continue
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		log.Printf("Downloaded %v\n", v)
	}

	return nil
}

func getPlaylist(urlStr string) ([]string, error) {

	cache := map[string]struct{}{}
	playlistURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	urlList := []string{}
	for {
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			return nil, err
		}
		resp, err := doRequest(req)
		if err != nil {
			return nil, err
		}
		playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		if listType == m3u8.MEDIA {
			mpl := playlist.(*m3u8.MediaPlaylist)
			for _, v := range mpl.Segments {
				if v != nil {
					var msURI string
					if strings.HasPrefix(v.URI, "http") {
						msURI, err = url.QueryUnescape(v.URI)
						if err != nil {
							return nil, err
						}
					} else {
						msURL, err := playlistURL.Parse(v.URI)
						if err != nil {
							log.Print(err)
							continue
						}
						msURI, err = url.QueryUnescape(msURL.String())
						if err != nil {
							return nil, err
						}
					}
					_, hit := cache[msURI]
					if !hit {
						cache[msURI] = struct{}{}
						urlList = append(urlList, msURI)
					}
				}
			}
			if mpl.Closed {
				break
			} else {
				time.Sleep(time.Duration(int64(mpl.TargetDuration * 1000000000)))
			}
		} else {
			return nil, errors.New("Not a valid media playlist")
		}
	}

	return urlList, nil
}

// HLSdownload HLSdownload
func HLSdownload(url string, fpath string) error {

	tsList, err := getPlaylist(url)
	if err != nil {
		return err
	}

	return downloadSegment(fpath, tsList)
}
