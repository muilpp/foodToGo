## Food to Go

Notifies you every time there's an available item in any of your favourite shops. Developed with Go.

### Features

- Sends a notification (mail or telegram message) every time one of your favourite shop uploads a new item
- Get reports by shop, by day of week and by hour of day.

### Usage

Set up the .env file with your config and you'll be ready to go. 
The only mandatory fields are user, password, id, longitude and latitude. Everything else is optative. 
Please note that you have to fill either the mail or telegram properties to receive notifications. If the database fields are not given, all the information will be stored in a raw file.

To run the app, use the argument:
 - "getFood" to fetch new items.
 - "printGraph" to send the reports (only works with telegram).

### Examples
By store
![alt text](https://github.com/muilpp/foodToGo/blob/main/food-by-store.gif)

By day of week 
![alt text](https://github.com/muilpp/foodToGo/blob/main/food-by-day.gif)

By hour of day
![alt text](https://github.com/muilpp/foodToGo/blob/main/food-by-hour.gif)
