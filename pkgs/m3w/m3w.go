package m3w

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/ripx80/brave/glue"
	"github.com/ripx80/recipe"
)

func getLastID(baseURL string) (int, error) {
	resp, err := http.Get(baseURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, err
	}
	link, ex := doc.Find(".rezeptlink").Eq(0).Attr("href")
	if !ex {
		return 0, fmt.Errorf("canot find the last id")
	}
	base, err := url.Parse(link)
	if err != nil {
		return 0, err
	}

	keys, ok := base.Query()["id"]
	if !ok {
		return 0, fmt.Errorf("no id extraction from index url")
	}
	lid, err := strconv.Atoi(keys[0])
	return lid, nil
}

func getUserComments(recipeID int) (*[]recipe.Comment, error) {
	var comments []recipe.Comment

	resp, err := http.Get(fmt.Sprintf("https://www.maischemalzundmehr.de/index.php?id=%d&inhaltmitte=rezept", recipeID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("div.nutzerKommentarBox").Each(func(i int, s *goquery.Selection) {
		child := s.Find("span").Contents()
		name := strip.StripTags(child.Eq(0).Text())
		datetime := strings.Split(glue.RmChar(strip.StripTags(child.Eq(1).Text()), " \n"), ",")
		comment := s.Find("div.kommentar_text").Contents().Text()
		comments = append(comments, recipe.Comment{
			Name:    name,
			Date:    datetime[0],
			Comment: comment,
		})
	})

	return &comments, nil
}

/*Down fetch all recipes form m3*/
func Down(m3url, outdir string) {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.Mkdir(outdir, 0755)
	}

	// simple check if some recipes downloaded before
	files, err := ioutil.ReadDir(outdir)
	if err != nil {
		log.Fatal(err)
	}

	var keymap []int
	for _, f := range files {
		id, err := strconv.Atoi(strings.Split(f.Name(), "_")[0])
		if err == nil {
			keymap = append(keymap, id)
		}
	}

	lid, err := getLastID(fmt.Sprintf("%s/index.php", m3url))
	if err != nil {
		log.Fatal(err)
	}

	cnt := 0
	errcnt := 0
	var body []byte

	//signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Printf("\n\nFetch: %d doc\nBroken docs: %d\n", cnt, errcnt)
		os.Exit(0)
		return
	}()

	for ; lid > 0; lid-- {
		if glue.InSlice(lid, keymap) {
			continue
		}
		//exclude broken recipes
		if lid == 786 || lid == 1240 || lid == 676 {
			continue
		}

		//get json export
		resp, err := http.Get(fmt.Sprintf("%s/export.php?id=%d", m3url, lid))
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)

		//empty json
		if len(body) <= 471 {
			continue
		}

		cnt++
		recipe, err := (&recipe.M3{}).Load(string(body))
		if err != nil {
			errcnt++
			broken := path.Join(outdir, "broken")
			if _, err := os.Stat(broken); os.IsNotExist(err) {
				os.Mkdir(broken, 0755)
			}
			ioutil.WriteFile(path.Join(broken, fmt.Sprintf("%d.json", lid)), body, 0644)
			log.Println(err)
			continue
		}

		fn := strings.ReplaceAll(recipe.Global.Name, " ", "_")
		fn = glue.RmChar(fn, "?%$&#-`'()./<>")

		//get comments
		c, err := getUserComments(lid)
		if err != nil {
			log.Fatal("usercomments: ", err)
		}
		recipe.Comments = *c

		recipe.SavePretty(path.Join(outdir, fmt.Sprintf("%d_%s.json", lid, fn)))
		log.Println(recipe.Global.Name)

	}
	log.Printf("\n\nFetch: %d doc\nBroken docs: %d\n", cnt, errcnt)
}
