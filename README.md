# Gator CLI program

### Command line program that reads RSS feeds

### Requires postgres in order to run. The cli saves posts and feeds to the db.

### With those prerequisites, after downloading the repo, you should be able to run go install .

### you do need to create a .gatorconfig.json in your linux home directory which needs to contain a postgres dsn
### Example: "{"db_url":"postgres://<dsn>"}"