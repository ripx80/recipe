# Recipe

Library to create well structured brew recipes for brewman.
At the moment only one external source is implemented.

- M3: maische-malz-und-mehr.de

*they change the site and the json structure this repo handled with. only exists for blueprinting...*

## M3

Only normal recipies supported. Not scaled recipies on M3 because they change the type of values so the parser is confused.
Use the internal scale to change the water or the yield.

## m3c

Converter for M3 recipes into brewman format.

## m3d - maischemalzundmehr crawler

Get all recipes from m3 and convert it to brewman recipe json format. If you have downloaded recipes before they will be skipped if they are in the same dir.

## Broken json

- Parsing Amount error: Hopfen_VWH_1_Menge strconv.ParseFloat: parsing " 90": invalid syntax
- json: invalid use of ,string struct tag, trying to unmarshal "" into float64
- Rest value is empty: Infusion_Rastzeit4
- Parsing Amount error: Infusion_Rastzeit2 strconv.Atoi: parsing "7.5": invalid syntax

## Notices

### Elasticsearch Stuff

insert into recipes-m3 index. change your user:password arg

```bash
curl --user elastic:changeme -XPOST http://node1:9200/recipes-m3/doc -H "Content-Type: application/json" -d @400_Meraner_Weizen.json
for i in $(ls):; do curl --user elastic:changeme -XPOST http://node1:9200/recipes-m3/doc -H "Content-Type: application/json" -d @$i; done

#or with bulk
PUT hockey/_bulk?refresh
{"index":{"_id":1}}
{"first":"johnny","last":"gaudreau","goals":[9,27,1],"assists":[17,46,0],"gp":[26,82,1],"born":"1993/08/13"}
{"index":{"_id":2}}
```

Found a lot of interesting misspelled stuff in recipes. So this will be corrected.
Then convert the Data like: "Pilsener" and "Pilsener Malz" zu "Pilsener" in mapping from m3 to recipe struct

- convert Pilsner Malz, Pilsner, Pilsener, Pilsenermalz, "Pilsenermalz, hell" to Pilsener Malz with data pipelines in es
- Pale Ale, Pale ale Malz, Best Pale Ale Malz-> Pale Ale Malz
- Münchener, Münchner Malz -> Münchener Malz
- Wiener, Wienermalz -> Wiener Malz
- Weizenmalz hell, "Weizenmalz, hell", Weizenmalz Hell -> Weizenmalz
- Carahell, CaraHell, Cara Hell
- CaraMünch II, Cara II
- CaraAmber, Cara Amber
- Amber Malz, Amber Malt
- Best Röstmalz -> Röstmalz
- Best Red X -> Red X

find duplicates with elastic search and remove them: [elastic.co/de/blog/how-to-find-and-remove-duplicate-documents-in-elasticsearch](remove-duplicates-es)

## Calculations

[doc](http://www.mathe-fuer-hobbybrauer.de/bierrezepte/)
