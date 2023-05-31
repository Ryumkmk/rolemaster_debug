package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"RMV0.5/app/models"
)

// TopページのHtmlを生成
func top(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		//未ログイン情報ならloginページへ飛ぶ
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		//ログイン状態ならTopを作成
		generateHTML(w, nil, "layout", "top")
	}
}

// 初期アクセスページ
func index(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		//未ログイン情報ならloginページへ飛ぶ
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		//ログイン状態ならTopページへ飛ぶ
		http.Redirect(w, r, "/top", http.StatusFound)
	}
}

func allpjs(w http.ResponseWriter, r *http.Request) {
	y, _ := models.GetAllPjsByDB()
	generateHTML(w, y, "layout", "pjs")
}

// 入力ページのHtmlを生成
func typingpage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	//入力した日付を2023-03-12のフォーマットに変換し、dateに格納
	yy := r.PostFormValue("year")
	mm := r.PostFormValue("month")
	if len(mm) == 1 {
		mm = "0" + mm
	}
	dd := r.PostFormValue("day")
	if len(dd) == 1 {
		dd = "0" + dd
	}
	date := fmt.Sprintf("%s-%s-%s", yy, mm, dd)
	// fmt.Println(date)

	//日付から婚礼情報を取得
	wITPs := models.GetWeddingsByDateFromDB(date)
	if len(wITPs) == 0 {
		http.Redirect(w, r, "/top", http.StatusFound)
		return
	}

	//日付から出勤PJ情報を取得
	pLTIs, err := models.GetPjsByDateFromDB(date)
	if err != nil {
		log.Println(err)
	}
	//AMPMのそれぞれの役割を取得
	var rIITPsAM []models.RoleInfoInTypingPage
	var rIITPsPM []models.RoleInfoInTypingPage
	var wITPAM models.WeddingInTypingPage
	var wITPPM models.WeddingInTypingPage

	for _, wITP := range wITPs {
		if wITP.Ampm == "AM" {
			rIITPsAM, err = wITP.GetRoleInfoByDateFromDB()
			fmt.Println(rIITPsAM)
			wITPAM = wITP
			if err != nil {
				log.Println(err)
			}
		} else {
			rIITPsPM, err = wITP.GetRoleInfoByDateFromDB()
			fmt.Println(rIITPsPM)
			wITPPM = wITP
			if err != nil {
				log.Println(err)
			}
		}
	}
	//全ての情報をdITPに格納
	var dITP = models.DataInTypingPage{
		PLITs:    pLTIs,
		WITPAM:   wITPAM,
		WITPPM:   wITPPM,
		RIITPsAM: rIITPsAM,
		RIITPsPM: rIITPsPM,
	}
	// fmt.Println(dITP.WITPAM)
	// fmt.Println(dITP.WITPPM)
	// fmt.Println(dITP.RIITPsAM, dITP.RIITPsPM)
	// fmt.Println(dITP.PLITs)

	//ダブル、つまり婚礼情報を二つ持っていたらダブルの入力ページを生成
	if len(wITPs) == 2 {
		generateHTML(w, dITP, "layout", "doubleTypingPage")
	} else {
		//AMならAM入力ページを生成
		if dITP.WITPAM.Ampm == "AM" {
			generateHTML(w, dITP, "layout", "amTypingPage")
		} else {
			//PM or 試食会ならPM入力ページを生成
			generateHTML(w, dITP, "layout", "pmTypingPage")
		}
	}
}

// 確認ページのHtmlを生成
func cheakPj(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	var rIITPsAM []models.RoleInfoInTypingPage
	var rIITPsPM []models.RoleInfoInTypingPage
	var wITPAM models.WeddingInTypingPage
	var wITPPM models.WeddingInTypingPage

	//入力フォームの値を取得
	for n, v := range r.Form {
		if strings.Contains(n, "trainer") || strings.Contains(n, "trainee") {
			// models.WhosTrainee(n, v[0], &whatJob.Trainers)
		} else if strings.Contains(n, "date-form") {
			wITPAM.Date = v[0]
			wITPPM.Date = v[0]
		} else if strings.Contains(n, "ampm-form") {
			if v[0] == "AM" {
				wITPAM.Ampm = v[0]
			} else {
				wITPPM.Ampm = v[0]
			}
		} else {
			// fmt.Println(n)
			if len(v[0]) > 0 && v[0][len(v[0])-1] == 'P' {

				rIITPPM := models.RoleInfoInTypingPage{
					RoleName: n,
					PjName:   v[0],
				}
				rIITPsPM = append(rIITPsPM, rIITPPM)
			} else if len(v[0]) > 0 {
				rIITPAM := models.RoleInfoInTypingPage{
					RoleName: n,
					PjName:   v[0],
				}
				rIITPsPM = append(rIITPsAM, rIITPAM)
			}
		}
	}
	var dITP = models.DataInTypingPage{
		WITPAM:   wITPAM,
		WITPPM:   wITPPM,
		RIITPsAM: rIITPsAM,
		RIITPsPM: rIITPsPM,
	}
	//データベース更新
	// err = dITP.UpdateRoleInfoDB()
	if err != nil {
		log.Println(err)
	}
	// generateHTML(w, dITP, "layout", "checkpage")
	fmt.Println(dITP.WITPAM, dITP.RIITPsAM)
	// fmt.Println(dITP.WITPPM, dITP.RIITPsPM)

}

/*
func attendancelist(w http.ResponseWriter, r *http.Request) {
	var List models.AttendanceList
	List.PjName = r.FormValue("pjName")
	List.AttendanceDaysList, List.AttendanceTimeList = models.GetAttendanceList(List.PjName)
	// fmt.Println(List)
	generateHTML(w, List, "layout", "attendancelist")
}

*/
