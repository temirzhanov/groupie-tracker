govendor init
go build

git add . 
git commit -m "deploy heroku"
git push heroku master
