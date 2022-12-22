# Data loading 
data/FoodData_Central_branded_food_json_2022-10-28.zip:
	curl -Lo $@ https://fdc.nal.usda.gov/fdc-datasets/FoodData_Central_branded_food_json_2022-10-28.zip

data/FoodData_Central_foundation_food_json_2022-10-28.zip:
	curl -Lo $@ https://fdc.nal.usda.gov/fdc-datasets/FoodData_Central_foundation_food_json_2022-10-28.zip

data/branded-foods.json: data/FoodData_Central_branded_food_json_2022-10-28.zip
	unzip -p $? | jq -c .BrandedFoods[] > $@

data/foundation-foods.json: data/FoodData_Central_foundation_food_json_2022-10-28.zip
	unzip -p $? | jq -c .FoundationFoods[] > $@

.PHONY: crdb-trgrm
crdb-trgrm: data/foundation-foods.json data/branded-foods.json
	go build -o crdb-trgrm main.go
