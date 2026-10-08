package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"siakad-mini/app/model"
)

type EnrollmentRepository interface {
	ListByStudent(ctx context.Context, studentID int64, tx pgx.Tx) ([]model.Enrollment, error)
	TotalSKS(ctx context.Context, studentID int64, tx pgx.Tx) (int, error)
	FindDuplicate(ctx context.Context, studentID, courseID int64, year string, tx pgx.Tx) (bool, error)
	Create(ctx context.Context, tx pgx.Tx, e model.Enrollment) (model.Enrollment, error)
	DeleteOwned(ctx context.Context, tx pgx.Tx, enrollmentID, studentID int64) error
}

type enrollmentPostgresRepository struct{ pool *pgxpool.Pool }

func NewEnrollmentRepository(pool *pgxpool.Pool) EnrollmentRepository {
	return &enrollmentPostgresRepository{pool: pool}
}

func (r *enrollmentPostgresRepository) ListByStudent(ctx context.Context, studentID int64, tx pgx.Tx) ([]model.Enrollment, error) {
	rows, err := tx.Query(ctx, `SELECT e.id,e.student_id,e.course_id,e.tahun_akademik,e.created_at,c.kode_mk,c.nama_mk,c.sks
        FROM enrollments e JOIN courses c ON c.id=e.course_id WHERE e.student_id=$1 ORDER BY e.created_at DESC`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.Enrollment{}
	for rows.Next() {
		var e model.Enrollment
		if err := rows.Scan(&e.ID, &e.StudentID, &e.CourseID, &e.TahunAkademik, &e.CreatedAt, &e.KodeMK, &e.NamaMK, &e.SKS); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (r *enrollmentPostgresRepository) TotalSKS(ctx context.Context, studentID int64, tx pgx.Tx) (int, error) {
	var total int
	err := tx.QueryRow(ctx, `SELECT COALESCE(SUM(c.sks),0)::int FROM enrollments e JOIN courses c ON c.id=e.course_id WHERE e.student_id=$1`, studentID).Scan(&total)
	return total, err
}
func (r *enrollmentPostgresRepository) FindDuplicate(ctx context.Context, studentID, courseID int64, year string, tx pgx.Tx) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM enrollments WHERE student_id=$1 AND course_id=$2 AND tahun_akademik=$3)`, studentID, courseID, year).Scan(&exists)
	return exists, err
}
func (r *enrollmentPostgresRepository) Create(ctx context.Context, tx pgx.Tx, e model.Enrollment) (model.Enrollment, error) {
	err := tx.QueryRow(ctx, `INSERT INTO enrollments(student_id,course_id,tahun_akademik) VALUES($1,$2,$3)
        RETURNING id,created_at`, e.StudentID, e.CourseID, e.TahunAkademik).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return e, fmt.Errorf("create enrollment: %w", err)
	}
	return e, nil
}
func (r *enrollmentPostgresRepository) DeleteOwned(ctx context.Context, tx pgx.Tx, id, studentID int64) error {
	tag, err := tx.Exec(ctx, `DELETE FROM enrollments WHERE id=$1 AND student_id=$2`, id, studentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

var _ = errors.Is
