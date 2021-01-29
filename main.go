package main

import (
  "fmt"
  "log"
  "net/http"
  "github.com/PuerkitoBio/goquery"
)

const SOURCE_URL = "https://japancars.ru/index.php?route=catalog/mitsubishi&cfg_id=1&area=M00&cat=A013A8A6A&mdl=PE8W&clf=HSEUF9"

func grable(url string, table *string) {
  response, error := http.Get(url)

  if error != nil {
    log.Fatalln(error)
  }

  defer func() {
    error := response.Body.Close()

    if error != nil {
      log.Fatalln(error)
    }
  }()

  if response.StatusCode != 200 {
    log.Fatalf("Status error: %d %s", response.StatusCode, response.Status)
  }

  document, error := goquery.NewDocumentFromReader(response.Body)

  if error != nil {
    log.Fatalln(error)
  }

  document.Find("table.list > tbody tr").Each(func(_ int, selection *goquery.Selection) {
    selection.Find("td").Each(func(index int, selection *goquery.Selection) {
      selection.Find("a").Each(func(_ int, selection *goquery.Selection) {
        href, exists := selection.Attr("href")

        if exists {
          name := selection.Text()

          if name != "Цена" {
            fmt.Println(name, href)
            grable(href, table)
          }
        }
      })

      if index != 0 {
        *table += ", "
      }
      *table += selection.Text()
    })
    *table += ", \n"
  })
}

func main() {
  table := ""

  grable(SOURCE_URL, &table)

  fmt.Println(table)
}