package psql

import (
	"errors"
	"log"

	connection "github.com/Sakagam1/DBMS_TASK/internal/db/db_connection"
	"github.com/Sakagam1/DBMS_TASK/internal/models"
	"github.com/Sakagam1/DBMS_TASK/internal/repositories"
)

type JokeRepository struct {
	joke repositories.IJoke
}

func (j JokeRepository) AddToFavorite(user_id int, joke_id int) (err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return err
	}
	qry := `INSERT INTO public."Favorite jokes" (joke_id, user_id) values ($1, $2)`
	_, err = DB.Exec(qry, user_id, joke_id)
	if err != nil {
		log.Println("Adding to favorite error:", err)
		return err
	}
	return nil
}

func (j JokeRepository) DeleteFromFavorite(user_id int, joke_id int) (err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return err
	}
	qry := `DELETE FROM public."Favorite jokes" where user_id=$1 and joke_id=$2`
	_, err = DB.Exec(qry, user_id, joke_id)
	if err != nil {
		log.Println("Adding to favorite error:", err)
		return err
	}
	return nil
}

func (j JokeRepository) GetUserFavoriteJokes(user_id int) (jokes []models.Joke, err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return nil, err
	}
	qry := `select "Jokes".id, "Jokes".header, "Jokes".description, "Jokes".rating from public."Jokes", public."Users", public."Favorite jokes" where "Users".id="Favorite jokes".user_id and "Favorite jokes".joke_id="Jokes".id and "Users".id=$1`
	rows, err := DB.Query(qry, user_id)
	defer rows.Close()
	if err != nil {
		log.Println("Connection Error:", err)
		return nil, err
	}
	for rows.Next() {
		var id, rating int
		var header, description string
		err := rows.Scan(&id, &header, &description, &rating)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return nil, err
		}
		NewJoke := models.Joke{
			ID:          id,
			Header:      header,
			Description: description,
			Rating:      rating,
			AuthorId:    user_id,
		}
		jokes = append(jokes, NewJoke)
	}
	return jokes, nil
}

func (j JokeRepository) GetJokesByTag(tag_name string) (jokes []models.Joke, err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return nil, err
	}
	qry := `select "Jokes".id, "Jokes".header, "Jokes".description, "Jokes".rating from public."Jokes", public."TagsJokes", public."Tags" where "Jokes".id="TagsJokes".joke_id and "TagsJokes".tag_id="Tags".id and "Tags".name LIKE '$1'`
	rows, err := DB.Query(qry, tag_name)
	defer rows.Close()
	if err != nil {
		log.Println("Error while getting jokes by tag:", err)
		return nil, err
	}
	for rows.Next() {
		var id, rating, user_id int
		var header, description string
		err := rows.Scan(&id, &header, &description, &rating)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return nil, err
		}
		NewJoke := models.Joke{
			ID:          id,
			Header:      header,
			Description: description,
			Rating:      rating,
			AuthorId:    user_id,
		}
		jokes = append(jokes, NewJoke)
	}
	return jokes, nil
}

func (j JokeRepository) GetJokesByKeyword(keyword string) (jokes []models.Joke, err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return nil, err
	}
	qry := `select * from public."Jokes" where header LIKE '%$1%' or description LIKE '%$2%'`
	rows, err := DB.Query(qry, keyword, keyword)
	defer rows.Close()
	if err != nil {
		log.Println("Error while getting jokes by keyword:", err)
		return nil, err
	}
	for rows.Next() {
		var id, rating, user_id int
		var header, description string
		err := rows.Scan(&id, &header, &description, &rating)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return nil, err
		}
		NewJoke := models.Joke{
			ID:          id,
			Header:      header,
			Description: description,
			Rating:      rating,
			AuthorId:    user_id,
		}
		jokes = append(jokes, NewJoke)
	}
	return jokes, nil
}

func (j JokeRepository) GetUserJokes(user_id int, page int, per_page int, sort_mode string) (jokes []models.Joke, err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return nil, err
	}
	qry := ``
	if sort_mode == "no" {
		qry = `select "Jokes".id, "Jokes".header, "Jokes".description, "Jokes".rating, "Jokes".creation_date from public."Jokes", public."Users" where "Users".id="Jokes".author_id and "Users".id=$1 ORDERED BY creation_date DESC LIMIT $1 OFFSET $2`
	}
	if sort_mode == "all" {
		qry = `select "Jokes".id, "Jokes".header, "Jokes".description, "Jokes".rating, "Jokes".creation_date from public."Jokes", public."Users" where "Users".id="Jokes".author_id and "Users".id=$1 ORDERED BY rating DESC LIMIT $1 OFFSET $2`
	}
	if sort_mode == "hour" {
		qry = `select "Jokes".id, "Jokes".header, "Jokes".description, "Jokes".rating, "Jokes".creation_date from public."Jokes", public."Users" where "Users".id="Jokes".author_id and "Users".id=$1 and EXTRACT(HOUR from (CURRENT_TIMESTAMP - "Jokes".creation_date)) <= 1 ORDER BY rating DESC LIMIT $2 OFFSET $3`
	}
	if sort_mode == "day" {
		qry = `select "Jokes".id, "Jokes".header, "Jokes".description, "Jokes".rating, "Jokes".creation_date from public."Jokes", public."Users" where "Users".id="Jokes".author_id and "Users".id=$1 and EXTRACT(DAY from (CURRENT_TIMESTAMP - "Jokes" creation_date)) <= 1 ORDER BY rating DESC LIMIT $2 OFFSET $3`
	}
	if sort_mode == "week" {
		qry = `select "Jokes".id, "Jokes".header, "Jokes".description, "Jokes".rating, "Jokes".creation_date from public."Jokes", public."Users" where "Users".id="Jokes".author_id and "Users".id=$1 and EXTRACT(DAY from (CURRENT_TIMESTAMP - "Jokes".creation_date)) <= 7 ORDER BY rating DESC LIMIT $2 OFFSET $3`
	}
	if sort_mode == "month" {
		qry = `select "Jokes".id, "Jokes".header, "Jokes".description, "Jokes".rating, "Jokes".creation_date from public."Jokes", public."Users" where "Users".id="Jokes".author_id and "Users".id=$1 and EXTRACT(MONTH from (CURRENT_TIMESTAMP - "Jokes".creation_date)) <= 1 ORDER BY rating DESC LIMIT $2 OFFSET $3`
	}
	rows, err := DB.Query(qry, user_id, per_page, page*per_page)
	defer rows.Close()
	if err != nil {
		log.Println("Error while getting user jokes:", err)
		return nil, err
	}
	for rows.Next() {
		var id, rating int
		var header, description string
		err := rows.Scan(&id, &header, &description, &rating)
		if err != nil {
			log.Println("Error while scanning rows:", err)
		}
		NewJoke := models.Joke{
			ID:          id,
			Header:      header,
			Description: description,
			Rating:      rating,
			AuthorId:    user_id,
		}
		jokes = append(jokes, NewJoke)
	}
	return jokes, nil
}

func (j JokeRepository) GetJokeTags(joke_id int) (tags []models.Tag, err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return nil, err
	}
	qry := `select "Tags".id, "Tags".name from public."Jokes", public."TagsJokes", public."Tags" where "Jokes".id="TagsJokes".joke_id and "TagsJokes".tag_id="Tags".id and "Jokes".id=$1`
	rows, err := DB.Query(qry, joke_id)
	defer rows.Close()
	if err != nil {
		log.Println("Error while getting joke tags:", err)
		return nil, err
	}
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return nil, err
		}
		NewTag := models.Tag{
			ID:   id,
			Name: name,
		}
		tags = append(tags, NewTag)
	}
	return tags, nil
}

func (j JokeRepository) AddTagToJoke(joke_id int, tag_id int) (err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return err
	}
	qry := `INSERT INTO public."TagsJokes" (tag_id, joke_id) values ($1, $2)`
	_, err = DB.Exec(qry, tag_id, joke_id)
	if err != nil {
		log.Println("Error while trying to add tag to joke:", err)
		return err
	}
	return nil
}

func (j JokeRepository) DeleteTagFromJoke(joke_id int, tag_id int) (err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return err
	}
	qry := `DELETE FROM public."TagsJokes" where tag_id=$1 and joke_id=$2`
	_, err = DB.Exec(qry, tag_id, joke_id)
	if err != nil {
		log.Println("Error while trying to add tag to joke:", err)
		return err
	}
	return nil
}

func (j JokeRepository) GetJokeByID(JokeId int) (userOut *models.Joke, err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return nil, err
	}
	qry := `select * from public."Jokes" where id=$1`
	rows, err := DB.Query(qry, JokeId)
	defer rows.Close()
	if err != nil {
		log.Println("Error while searching joke by id:", err)
	}
	var id, rating, author_id int
	var header, description string
	id = -1
	for rows.Next() {
		err := rows.Scan(&id, &header, &description, &rating, &author_id)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return nil, err
		}
	}
	if id != -1 {
		return &models.Joke{
			ID:          id,
			Header:      header,
			Description: description,
			Rating:      rating,
			AuthorId:    author_id,
		}, nil
	}
	return &models.Joke{}, errors.New("Joke with this id does not exist!")
}

func (j JokeRepository) GetPageOfJokes(page int, per_page int, sort_mode string) (jokes []models.Joke, err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return nil, err
	}
	qry := ``
	if sort_mode == "no" {
		qry = `select * from public."Jokes" ORDERED BY creation_date DESC LIMIT $1 OFFSET $2`
	}
	if sort_mode == "all" {
		qry = `select * from public."Jokes" ORDERED BY rating DESC LIMIT 5 OFFSET 1`
	}
	if sort_mode == "hour" {
		qry = `select * from public."Jokes" where EXTRACT(HOUR from (CURRENT_TIMESTAMP - creation_date)) <= 1 ORDER BY rating DESC LIMIT 5 OFFSET 1`
	}
	if sort_mode == "day" {
		qry = `select * from public."Jokes" where EXTRACT(DAY from (CURRENT_TIMESTAMP - creation_date)) <= 1 ORDER BY rating DESC LIMIT 5 OFFSET 1`
	}
	if sort_mode == "week" {
		qry = `select * from public."Jokes" where EXTRACT(DAY from (CURRENT_TIMESTAMP - creation_date)) <= 7 ORDER BY rating DESC LIMIT 5 OFFSET 1`
	}
	if sort_mode == "month" {
		qry = `select * from public."Jokes" where EXTRACT(MONTH from (CURRENT_TIMESTAMP - creation_date)) <= 1 ORDER BY rating DESC LIMIT 5 OFFSET 1`
	}
	rows, err := DB.Query(qry, per_page, per_page*page)
	defer rows.Close()
	if err != nil {
		log.Println("Error while trying to get page of jokes:", err)
		return nil, err
	}
	for rows.Next() {
		var id, rating, author_id int
		var header, description string
		err := rows.Scan(&id, &header, &description, &rating, &author_id)
		if err != nil {
			log.Println("Error while scanning rows:", err)
			return nil, err
		}
		NewJoke := models.Joke{
			ID:          id,
			Header:      header,
			Description: description,
			Rating:      rating,
			AuthorId:    author_id,
		}
		jokes = append(jokes, NewJoke)
	}
	return jokes, nil
}

func (j JokeRepository) Create(joke *models.Joke) (err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return err
	}
	qry := `INSERT INTO public."Jokes" (id, header, description, rating, author_id) values ($1, $2, $3, $4, $5)`
	_, err = DB.Exec(qry, joke.ID, joke.Header, joke.Description, joke.Rating, joke.AuthorId)
	if err != nil {
		log.Println("Joke creation error:", err)
		return err
	}
	return nil
}

func (j JokeRepository) Delete(joke_id int) (err error) {
	DB, err := connection.GetConnectionToDB()
	if err != nil {
		log.Println("Connection error:", err)
		return err
	}
	qry := `DELETE FROM public."Jokes" where id=$1`
	_, err = DB.Exec(qry, joke_id)
	if err != nil {
		log.Println("Joke deletion error:", err)
		return err
	}
	return nil
}
