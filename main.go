package main

import (
	"TASK-9/connection"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// STRUCT TEMMPLATE
type Project struct {
	Id                        int
	Project_name              string
	Project_start_date        time.Time
	Project_end_date          time.Time
	Project_start_date_string string
	Project_end_date_string   string
	Project_duration          string
	Project_description       string
	Project_technologies      []string
	Project_image             string
}

// LOCAL DATABASE
var Projects = []Project{}

// MAIN
func main() {
	route := mux.NewRouter()

	//CONNECT TO DATABASE
	connection.DatabaseConnect()

	//ROUTE PATH FOLDER PUBLIC
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	//ROUTING
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/addmyproject", addMyProject).Methods("GET")
	route.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	route.HandleFunc("/form-project", formAddProject).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")
	route.HandleFunc("/edit-project/{id}", editProject).Methods("GET")
	route.HandleFunc("/form-update-project", formUpdateProject).Methods("GET")
	route.HandleFunc("/update-project", updateProject).Methods("POST")

	fmt.Println("Server running on port 5000")
	http.ListenAndServe("localhost:5000", route)

	//127.0.0.1 (localhost)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")

	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {

		var renderData []Project
		item := Project{}

		// GET ALL PROJECTS FROM POSTGRESQL
		rows, _ := connection.Conn.Query(context.Background(), `SELECT "id", "project_name", "project_start_date", "project_end_date", "project_description", "project_technologies" FROM public.tb_projects`)

		//PARSE PROJECT
		for rows.Next() {
			err := rows.Scan(&item.Id, &item.Project_name, &item.Project_start_date, &item.Project_end_date, &item.Project_description, &item.Project_technologies)
			//ERROR HANDLING GET PROJECTS FROM POSTGRESQL
			if err != nil {
				fmt.Println(err.Error())
				return
			} else {
				//PARSING DATE
				item := Project{
					Id:                   item.Id,
					Project_name:         item.Project_name,
					Project_duration:     GetDuration(item.Project_start_date, item.Project_end_date),
					Project_description:  item.Project_description,
					Project_technologies: item.Project_technologies,
					Project_image:        item.Project_image,
				}

				renderData = append(renderData, item)
			}

		}

		response := map[string]interface{}{
			"Projects": renderData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conten-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/contact-me.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func addMyProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Println(Projects)

	var tmpl, err = template.ParseFiles("views/add-my-project.html")
	// ERROR HANDLING RENDER PROJECT TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {

		var renderData []Project
		item := Project{}

		// GET ALL PROJECTS FROM POSTGRESQL
		rows, _ := connection.Conn.Query(context.Background(), `SELECT "id", "project_name", "project_start_date", "project_end_date", "project_description", "project_technologies" FROM public.tb_projects`)

		//PARSE PROJECT
		for rows.Next() {
			err := rows.Scan(&item.Id, &item.Project_name, &item.Project_start_date, &item.Project_end_date, &item.Project_description, &item.Project_technologies)
			//ERROR HANDLING GET PROJECTS FROM POSTGRESQL
			if err != nil {
				fmt.Println(err.Error())
				return
			} else {
				//PARSING DATE
				item := Project{
					Id:                   item.Id,
					Project_name:         item.Project_name,
					Project_duration:     GetDuration(item.Project_start_date, item.Project_end_date),
					Project_description:  item.Project_description,
					Project_technologies: item.Project_technologies,
					Project_image:        item.Project_image,
				}

				renderData = append(renderData, item)
			}

		}

		response := map[string]interface{}{
			"Projects": renderData,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, response)
	}
}

func blogDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/blog-detail.html")

	//ERROR HANDLING BLOG ATAU PROJECT DETAIL
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		Id, _ := strconv.Atoi(mux.Vars(r)["Id"])
		renderDetail := Project{}

		//GET PROJECT BY ID FROM POSTGRESQL
		err = connection.Conn.QueryRow(context.Background(), `SELECT "id", "project_name", "project_start_date", "project_end_date", "project_description", "project_technologies", "project_image"
		FROM public.tb_projects WHERE "id" = $1`, Id).Scan(&renderDetail.Id, &renderDetail.Project_name, &renderDetail.Project_start_date, &renderDetail.Project_end_date, &renderDetail.Project_description, &renderDetail.Project_technologies, &renderDetail.Project_image)

		//ERROR HANDLING GET PROJECT DATA BY ID
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))

		} else {

			//PARSING DATE
			renderDetail := Project{
				Id:                        renderDetail.Id,
				Project_name:              renderDetail.Project_name,
				Project_start_date_string: FormatDate(renderDetail.Project_start_date),
				Project_end_date_string:   FormatDate(renderDetail.Project_end_date),
				Project_duration:          GetDuration(renderDetail.Project_start_date, renderDetail.Project_end_date),
				Project_description:       renderDetail.Project_description,
				Project_technologies:      renderDetail.Project_technologies,
				Project_image:             renderDetail.Project_image,
			}

			response := map[string]interface{}{
				"renderDetail": renderDetail,
			}

			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, response)
		}

	}

}

func formAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		projectName := r.PostForm.Get("projectName")
		projectStartDate := r.PostForm.Get("inputStartDate")
		projectEndDate := r.PostForm.Get("inputEndDate")
		projectDescription := r.PostForm.Get("inputContent")
		projectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("nextjs"), r.PostForm.Get("typescript")}
		projectImage := r.PostForm.Get("input-blog-image")

		//INSERT TO POSGRESQL
		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.tb_projects("project_name","project_start_date", "project_end_date", "project_description", "project_technologies", "project_image") VALUES ( $1, $2, $3, $4, $5, $6)`,
			projectName, projectStartDate, projectEndDate, projectDescription, projectTechnologies, projectImage)

		//ERROR HANDLING INSERT TO POSGRESQL
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message :" + err.Error()))
			return
		}

		http.Redirect(w, r, "/addmyproject", http.StatusMovedPermanently)

	}

}

func deleteBlog(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	//DELETE PROJECT BY ID AT POSTGRESQL
	_, err := connection.Conn.Exec(context.Background(), `DELETE FROM public.tb_projects WHERE "id" = $1`, id)

	//ERROR HANDLING DELETE PROJECT AT POSTGRESQL
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/addmyproject", http.StatusFound)
}

func editProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)

	Projects = append(Projects[:id], Projects[id+1:]...)
	fmt.Println(Projects)

	http.Redirect(w, r, "/form-update-project", http.StatusMovedPermanently)
}

func formUpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("views/update-project.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, nil)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {

		Id, _ := strconv.Atoi(mux.Vars(r)["id"])
		projectName := r.PostForm.Get("projectName")
		projectStartDate := r.PostForm.Get("inputStartDate")
		projectEndDate := r.PostForm.Get("inputEndDate")
		projectDescription := r.PostForm.Get("inputContent")
		projectTechnologies := []string{r.PostForm.Get("nodejs"), r.PostForm.Get("reactjs"), r.PostForm.Get("nextjs"), r.PostForm.Get("typescript")}
		projectImage := r.PostForm.Get("input-blog-image")

		//INSERT TO POSGRESQL
		_, err = connection.Conn.Exec(context.Background(), `UPDATE public.tb_projects SET "project_name"=$1, "project_start_date"=$2, "project_end_date"=$3, "project_description"=$4, "project_technologies"=$5, "project_image"=$6 WHERE "id"=$7`,
			projectName, projectStartDate, projectEndDate, projectDescription, projectTechnologies, projectImage, Id)

		//ERROR HANDLING INSERT TO POSGRESQL
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message :" + err.Error()))
			return
		}

		http.Redirect(w, r, "/addmyproject", http.StatusMovedPermanently)

	}

}

// ADDITIONAL FUNCTION

// GET DURATION
func GetDuration(startDate time.Time, endDate time.Time) string {

	margin := endDate.Sub(startDate).Hours() / 24
	var duration string

	if margin >= 30 {
		if (margin / 30) == 1 {
			duration = "1 Month"
		} else {
			duration = strconv.Itoa(int(margin/30)) + " Months"
		}
	} else {
		if margin <= 1 {
			duration = "1 Day"
		} else {
			duration = strconv.Itoa(int(margin)) + " Days"
		}
	}

	return duration
}

// CHANGE DATE FORMAT
func FormatDate(InputDate time.Time) string {

	Formated := InputDate.Format("02 January 2006")

	return Formated
}

// RETURN DATE FORMAT
func ReturnDate(InputDate time.Time) string {

	Formated := InputDate.Format("2006-01-02")

	return Formated
}
