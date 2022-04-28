package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/decadev/squad10/healthplus/db"
	"github.com/decadev/squad10/healthplus/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var Sessions *scs.SessionManager

//Indexhandler gets the homepage
func Indexhandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/index.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

// RegisterPatientHandler gets Patient's SignUp page
func RegisterPatientHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/registerPatient.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//PostRegisterPatientHandler successfully register's patient's name in the db if valid
func PostRegisterPatientHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Patient
	r.ParseForm()
	name := models.Capitalise(r.FormValue("name"))
	ageString := r.FormValue("age")
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	age, _ := strconv.Atoi(ageString)
	user.ID = uuid.NewString()
	user.Name = name
	user.Age = uint(age)
	user.Email = email
	user.Username = username
	user.Password = password
	_, err := db.FindUserByEmailandUserName(user.Email, user.Username)
	if err == nil {
		// this user already exists
		// return a message to the user
		t, e := template.ParseFiles("pages/registerPatient.html")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "User already exists. Check Email or Username")
		if e != nil {
			fmt.Println(e)
			return
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.PasswordHash = string(hashedPassword)
	db.CreatePatientInTable(user)
	file, err3 := template.ParseFiles("pages/registerPatient.html")
	if err3 != nil {
		fmt.Println(err3)
	}
	file.Execute(w, name+" "+" is now a Registered Patient. \n You can Login")

}

func PatientLoginHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/patientLogin.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//------------------------------PostPatientLoginHandler logs in doctor if valid-----------------------------------------------------
func PostLoginPatientdHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Doctor
	user.Username = strings.TrimSpace(r.FormValue("username"))
	user.Password = strings.TrimSpace(r.FormValue("password"))
	_, err := db.AuthenticatePatient(user.Username, user.Password)
	if err != nil {
		t, e := template.ParseFiles("pages/patientLogin.html")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "Check username or Password")
		if e != nil {
			fmt.Println(e)
			return
		}
		return
	}
	Sessions.Put(r.Context(), "username", user.Username)
	http.Redirect(w, r, "/patientDashboard", http.StatusFound)
}

//------------------------------PatientDashboardHandler gets Patient's Dashboard page-----------------------------------------------
func PatientHomeHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/patientDashboard.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	userName := Sessions.GetString(r.Context(), "username")
	patient, err := db.FindPatientByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}

	e = t.Execute(w, patient)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//------------------------------PatientLogoutHandler logsout ---------------------------------------------------------------------
func PatientLogoutHandler(w http.ResponseWriter, r *http.Request) {
	Sessions.Remove(r.Context(), "username")
	http.Redirect(w, r, "/", http.StatusFound)
}

//-------------------------RegisterDoctorHandler gets Doctor's SignUp page-----------------------------------------------
func RegisterDoctorHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/doctorRegister.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//-------------------PostRegisterDoctorHandler successfully registers doctor's name in the db if valid----------------------------
func PostRegisterDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Doctor
	ageString := r.FormValue("age")
	age, _ := strconv.Atoi(ageString)
	user.ID = uuid.NewString()
	user.Name = models.Capitalise(strings.TrimSpace(r.FormValue("name")))
	user.Age = uint(age)
	user.Email = strings.TrimSpace(r.FormValue("email"))
	user.Username = strings.TrimSpace(r.FormValue("username"))
	user.Password = strings.TrimSpace(r.FormValue("password"))
	user.Specialty = models.Capitalise(strings.TrimSpace(r.FormValue("specialty")))
	_, err := db.FindDocByEmailandUserName(user.Email, user.Username)
	if err == nil {
		t, e := template.ParseFiles("pages/doctorRegister.html")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "User already exists, confirm email or username")
		if e != nil {
			fmt.Println(e)
			return
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.PasswordHash = string(hashedPassword)
	db.CreateDocInTable(user)
	temp, errt := template.ParseFiles("pages/doctorRegister.html")
	if errt != nil {
		fmt.Println(errt)
	}
	temp.Execute(w, user.Name+" is now a registered Doctor. Login")

}

//------------------------------DoctorLoginHandler gets Doctor's Login page---------------------------------------------------------
func DoctorLoginHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/doctorLogin.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	e = t.Execute(w, nil)
	if e != nil {
		fmt.Println(e)
		return
	}
}

func PostLoginDoctordHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Doctor
	user.Username = strings.TrimSpace(r.FormValue("username"))
	user.Password = strings.TrimSpace(r.FormValue("password"))
	usa, err := db.Authenticate(user.Username, user.Password)
	if err != nil {
		t, e := template.ParseFiles("pages/doctorLogin.html")
		if e != nil {
			fmt.Println(e)
			return
		}
		e = t.Execute(w, "Check Username or Password")
		if e != nil {
			fmt.Println(e)
			return
		}
		return
	}
	Sessions.Put(r.Context(), "username", usa.Username)
	http.Redirect(w, r, "/doctorDashboard", http.StatusFound)
}

//------------------------------DoctorDashboardHandler gets Doctor's Dashboard page-----------------------------------------------
func DoctorHomeHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/doctorHome.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	userName := Sessions.GetString(r.Context(), "username")
	doc, err := db.FindDoctorByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	e = t.Execute(w, doc)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//------------------------------DoctorLogoutHandler logsout ---------------------------------------------------------------------
func DoctorLogoutHandler(w http.ResponseWriter, r *http.Request) {
	Sessions.Remove(r.Context(), "username")
	http.Redirect(w, r, "/", http.StatusFound)
}

//------------------------------List of Doctors for booking Appointments ---------------------------------------------------------------------
func DoctorListHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/doctorList.html")
	if e != nil {
		fmt.Println("now", e)
		return
	}

	e = t.Execute(w, db.GetAllDoctors())
	if e != nil {
		fmt.Println("no way", e)
		return
	}
}

func DoctorWorkingHoursHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/workinghours.html")
	if e != nil {
		fmt.Println("now", e)
		return
	}
	userName := Sessions.GetString(r.Context(), "username")
	doc, err := db.FindDoctorByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	e = t.Execute(w, doc)
	if e != nil {
		fmt.Println("no way", e)
		return
	}
}

func ChooseHoursHandler(w http.ResponseWriter, r *http.Request) {

	userName := Sessions.GetString(r.Context(), "username")
	doc, err := db.FindDoctorByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}

	e := r.ParseForm()
	if e != nil {
		fmt.Println(e)
	}

	starttime := r.PostForm.Get("start")
	closetime := r.PostForm.Get("end")
	startInt, _ := strconv.Atoi(starttime)
	closeInt, _ := strconv.Atoi(closetime)

	if startInt > closeInt {
		file, err := template.ParseFiles("pages/appointmentTimeError.html")
		if err != nil {
			fmt.Println(err)
		}
		file.Execute(w, "START TIME LATER THAN END TIME ")
		return
	}
	checkStart := startInt > 12
	noonStart := startInt == 12
	fmt.Println(closetime)
	if checkStart {
		starttime = strconv.Itoa(startInt-12) + ":" + "00" + "PM"
	} else if noonStart {
		starttime = strconv.Itoa(startInt) + ":" + "00" + "PM"
	} else if !checkStart {
		starttime = strconv.Itoa(startInt) + ":" + "00" + "AM"
	}
	checkEnd := closeInt > 12
	noonEnd := closeInt == 12
	if checkEnd {
		closetime = strconv.Itoa(closeInt-12) + ":" + "00" + "PM"
	} else if noonEnd {
		closetime = strconv.Itoa(closeInt) + ":" + "00" + "PM"
	} else {
		closetime = strconv.Itoa(closeInt) + ":" + "00" + "AM"
	}
	//Gorm command to update a field
	db.DB.Model(&doc).Update("string_start", starttime)
	db.DB.Model(&doc).Update("string_close", closetime)
	db.DB.Model(&doc).Update("start_time", startInt)
	db.DB.Model(&doc).Update("close_time", closeInt)

	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/doctorDashboard", http.StatusFound)
}

func BookByIdHandler(w http.ResponseWriter, r *http.Request) {

	//This points to the html location
	t, err := template.ParseFiles("pages/appointments.html")
	if err != nil {
		fmt.Println("now", err)
		return
	}

	params := mux.Vars(r)
	ID := params["ID"]

	doctor := db.FindDoctorByIDandUserName(ID)
	workinghrs := doctor.SetWorkingHours()

	err = t.Execute(w, workinghrs)
	if err != nil {
		fmt.Println("now", err)
		return
	}
}

func PostBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment

	userName := Sessions.GetString(r.Context(), "username")
	patient, err := db.FindPatientByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	e := r.ParseForm()
	if e != nil {
		fmt.Println(e)
	}

	appointment.ID = uuid.NewString()
	appointment.AppointmentHour = r.PostForm.Get("time")
	appointment.Purpose = r.PostForm.Get("purpose")
	params := mux.Vars(r)
	appointment.DoctorID = params["ID"]
	f := db.FindDoctorByID(appointment.DoctorID)
	fmt.Println(appointment.AppointmentHour)
	valid := db.CheckAppointmentIsValid(appointment.DoctorID, appointment.AppointmentHour)
	appointment.DoctorName = f.Name
	appointment.Date = time.Now().String()
	appointment.PatientName = patient.Name
	appointment.PatientID = patient.ID
	if valid {
		db.CreateAppointmentInTable(appointment)
		http.Redirect(w, r, "/patientDashboard", http.StatusFound)
	} else {
		file, err := template.ParseFiles("pages/appointmentError.html")
		if err != nil {
			fmt.Println(err)
		}
		file.Execute(w, "Appointment Time already taken")

	}

}

func CheckPatientAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/checkappointments.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	userName := Sessions.GetString(r.Context(), "username")
	patient, err := db.FindPatientByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	appointment := db.FindPatientAppointment(patient.ID)

	e = t.Execute(w, appointment)
	if e != nil {
		fmt.Println(e)
		return
	}
}

func DeletePatientAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["ID"]
	db.DeleteAppointmentbyID(ID)
	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/checkappointments", http.StatusFound)
}

func CheckDoctorAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	t, e := template.ParseFiles("pages/viewdocappointments.html")
	if e != nil {
		fmt.Println(e)
		return
	}
	userName := Sessions.GetString(r.Context(), "username")
	doctor, err := db.FindDoctorByUsername(userName)
	if err != nil {
		fmt.Println(err)
		return
	}
	appointment := db.FindDoctorAppointment(doctor.ID)

	e = t.Execute(w, appointment)
	if e != nil {
		fmt.Println(e)
		return
	}
}

func DeleteDoctorAppointmentHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["ID"]
	db.DeleteAppointmentbyID(ID)
	//redirect your page back to the index/home page when done (on a click)
	http.Redirect(w, r, "/viewdocappointments", http.StatusFound)
}
