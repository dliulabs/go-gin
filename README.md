# Creating a new REST API

## build a skeleton

- message structure
- route
- route handler


## Test Drive

```
curl -X POST http://localhost:8080/recipes -d '{"name":"food","tags":["tag1", "tag2"], "ingredients":["a", "b"], "instructions":["step 1", "step 2"]}"'

curl --location --request POST 'http://localhost:8080/recipes' \
--header 'Content-Type: application/json' \
--data-raw '{
   "name":"Homemade Pizza",
   "tags":[
      "italian",
      "pizza",
      "dinner"
   ],
   "ingredients":[
      "1 1/2 cups (355 ml) warm water (105°F-115°F)",
      "1 package (2 1/4 teaspoons) of active dry yeast",
      "3 3/4 cups (490 g) bread flour",
      "feta cheese, firm mozzarella cheese, grated"
   ],
   "instructions":[
      "Step 1.",
      "Step 2.",
      "Step 3."
   ]
}' | jq -r
```

## List Recipes

```
curl -X GET http://localhost:8080/recipes | jq -r
```

```
curl -X GET http://localhost:8080/recipes | jq length
```

```
curl -X PUT http://localhost:8080/recipes/62c0a82c087a22434f0a255c -d '{"name":"food2","tags":["tag1", "tag2"], "ingredients":["a", "b"], "instructions":["step 1", "step 2"]}"'
```

```
curl -X GET http://localhost:8080/recipes/62bfbb4586d77039967c61cf | jq -r
```

```
curl -X GET 'http://localhost:8080/recipes/search?tag=main'
```