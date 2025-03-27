package mod

import (
	"github.com/gocroot/lite/mod/daftar"
	"github.com/gocroot/lite/mod/helpdesk"
	"github.com/gocroot/lite/mod/idgrup"
	"github.com/gocroot/lite/mod/kyc"
	"github.com/gocroot/lite/mod/lms"
	"github.com/gocroot/lite/mod/lmsdesa"
	"github.com/gocroot/lite/mod/pomokit"
	"github.com/gocroot/lite/mod/posint"
	"github.com/gocroot/lite/mod/presensi"
	"github.com/gocroot/lite/mod/siakad"
	"github.com/gocroot/lite/mod/strava"
	"github.com/gocroot/lite/mod/tasklist"
	"github.com/gocroot/lite/model"
	"go.mongodb.org/mongo-driver/mongo"
)

func Caller(Profile model.Profile, Modulename string, Pesan model.WAMessage, db *mongo.Database) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = idgrup.IDGroup(Pesan)
	case "feedbackhelpdesk":
		reply = helpdesk.FeedbackHelpdesk(Profile, Pesan, db)
	case "endhelpdesk":
		reply = helpdesk.EndHelpdesk(Profile, Pesan, db)
	case "helpdesk":
		reply = helpdesk.StartHelpdesk(Pesan, db)
	case "presensi-masuk":
		reply = presensi.PresensiMasuk(Pesan, db)
	case "presensi-pulang":
		reply = presensi.PresensiPulang(Pesan, db)
	case "upload-lmsdesa-file":
		reply = lmsdesa.ArsipFile(Pesan, db)
	case "upload-lmsdesa-gambar":
		reply = lmsdesa.ArsipGambar(Pesan, db)
	case "lms":
		reply = lms.ReplyRekapUsers(Profile, Pesan, db)
	case "cek-ktp":
		reply = kyc.CekKTP(Profile, Pesan, db)
	case "selfie-masuk":
		reply = presensi.CekSelfieMasuk(Profile, Pesan, db)
	case "selfie-pulang":
		reply = presensi.CekSelfiePulang(Profile, Pesan, db)
	case "tasklist-append":
		reply = tasklist.TaskListAppend(Pesan, db)
	case "tasklist-reset":
		reply = tasklist.TaskListReset(Pesan, db)
	case "tasklist-save":
		reply = tasklist.TaskListSave(Pesan, db)
	case "domyikado-user":
		reply = daftar.DaftarDomyikado(Profile, Pesan, db)
	case "panduan-siakad":
		reply = siakad.PanduanDosen(Pesan, db)
	case "login-siakad":
		reply = siakad.LoginSiakad(Pesan, db)
	case "approve-bap":
		reply = siakad.ApproveBAP(Pesan, db)
	case "cek-approval":
		reply = siakad.CekApprovalBAP(Pesan, db)
	case "cetak-bap":
		reply = siakad.CetakBAP(Pesan, db)
	case "approve-bimbingan":
		reply = siakad.ApproveBimbingan(Pesan, db)
	case "approve-bimbinganbypoin":
		reply = siakad.ApproveBimbinganbyPoin(Pesan, db)
	case "prohibited-items":
		reply = posint.GetProhibitedItems(Pesan, db)
	case "pomodoro-report":
		reply = pomokit.HandlePomodoroReport(Profile, Pesan, db)
	case "pomodoro-start":
		reply = pomokit.HandlePomodoroStart(Profile, Pesan, db)
	case "strava-identity":
		reply = strava.StravaIdentityHandler(Profile, Pesan, db)
	case "strava-activity":
		reply = strava.StravaActivityHandler(Profile, Pesan, db)
	case "strava-update-identity":
		reply = strava.StravaIdentityUpdateHandler(Profile, Pesan, db)
	case "strava-update-activity":
		reply = strava.StravaActivityUpdateIfEmptyDataHandler(Profile, Pesan, db)
	case "strava-poin":
		reply = strava.InisialisasiPoinDariAktivitasLama(Pesan, db)
	}

	return
}
