# Cards API
REST API to create and interact with decks of playing cards.

## Tech Stack
 - Go 1.16
 - MongoDB 

## Usage
 1.  Setup Go on your machine.
 2. Clone the repo.
 3. Setup and run a MongoDB instance locally. Either go through the manual setup or use the mongo Docker image
 ```docker run --name mongo -d  -p <port>:<port> mongo```
 Replace `<port>` with the port number you need.
 4. Copy the given `env_sample` file to new `.env` files
	```cp env_sample .env```
	```cp env_sample test.env``` for integration test env.
 5. Replace the values in `.env` and `test.env` with the appropriate database configuration values.
 6. Either run directly using ```go run .```
	or build and run using ```go build . && ./cards```
 7. Run all unit tests using `go test ./...`
 8. Run integration tests using `go test -tags=integration`

## API
### 1. Create new Deck
 `POST /deck` Creates a new deck according to the provided params.
#### Request Params  
| param | type | default | description|
| --- | --- | --- | --- |
| shuffle | boolean, optional | false | If true, the deck will be created in shuffled order |
|customDeck| boolean, optional| false | If true, the deck will be created using only cards provided in the `wantedCards` param|
| wantedCards| string array, optional| [] | If `customDeck` is true, this param _must_ be provided. The deck will be created using only the cards provided in this param. If `customDeck` is false, this param is ignored|

#### Response
| param | type | description|
| --- | --- | --- |
| deck_id | string | UUID of the created deck|
| shuffled | boolean | Indicates whether the deck was shuffled during creation |
| remaining | integer | The number of cards remaining in the deck |


### 2. Get Deck
 `GET /deck/{deck_uuid}` Returns the deck corresponding to the provided deck UUID
 
 #### Response
| param | type | description|
| --- | --- | --- |
| deck_id | string | UUID of the returned deck |
| shuffled | boolean | Indicates whether the deck was shuffled during creation |
| remaining | integer | The number of cards remaining in the deck |
| cards | array of card objects `{suit string, value string, code string}` | The cards in the deck |
 
 
 
 
 ### 3. Draw Cards
  `PATCH /deck/{deck_uuid}` Returns _n_ cards from the top of the deck corresponding to the provided deck UUID
  
  #### Request Params
  | param | type | default | description|
  | --- | --- | --- | --- |
  |numberOfCards| integer | N/A | The number of cards to draw. Must be greater than `0`.|
  
   #### Response
| param | type | description|
| --- | --- | --- |
| cards | array of card objects `{suit string, value string, code string}` | The drawn cards. |
