package model

import "time"

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type Student struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	NIM         string     `json:"nim"`
	Nama        string     `json:"nama"`
	Prodi       string     `json:"prodi"`
	Angkatan    int        `json:"angkatan"`
	IPKTerakhir float64    `json:"ipk_terakhir"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Course struct {
	ID        int64  `json:"id"`
	KodeMK    string `json:"kode_mk"`
	NamaMK    string `json:"nama_mk"`
	SKS       int    `json:"sks"`
	Semester  int    `json:"semester"`
	Kuota     int    `json:"kuota"`
	Terisi    int    `json:"terisi"`
	SisaKuota int    `json:"sisa_kuota"`
}

type Enrollment struct {
	ID            int64     `json:"id"`
	StudentID     int64     `json:"student_id"`
	CourseID      int64     `json:"course_id"`
	TahunAkademik string    `json:"tahun_akademik"`
	CreatedAt     time.Time `json:"created_at"`
	KodeMK        string    `json:"kode_mk,omitempty"`
	NamaMK        string    `json:"nama_mk,omitempty"`
	SKS           int       `json:"sks,omitempty"`
}

type AuthUser struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}

type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type StudentCreateRequest struct {
	NIM         string  `json:"nim" validate:"required,len=12,numeric"`
	Nama        string  `json:"nama" validate:"required,min=2,max=100"`
	Email       string  `json:"email" validate:"required,email"`
	Prodi       string  `json:"prodi" validate:"required,min=2,max=100"`
	Angkatan    int     `json:"angkatan" validate:"required,gte=2000,lte=2100"`
	IPKTerakhir float64 `json:"ipk_terakhir" validate:"gte=0,lte=4"`
}

type StudentUpdateRequest struct {
	Nama        string  `json:"nama" validate:"required,min=2,max=100"`
	Prodi       string  `json:"prodi" validate:"required,min=2,max=100"`
	Angkatan    int     `json:"angkatan" validate:"required,gte=2000,lte=2100"`
	IPKTerakhir float64 `json:"ipk_terakhir" validate:"gte=0,lte=4"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type EnrollmentRequest struct {
	CourseID      int64  `json:"course_id" validate:"required,gt=0"`
	TahunAkademik string `json:"tahun_akademik" validate:"required"`
}

type StudentDetail struct {
	Student     Student      `json:"student"`
	Enrollments []Enrollment `json:"mata_kuliah"`
	TotalSKS    int          `json:"total_sks"`
	BatasSKS    int          `json:"batas_sks"`
}
